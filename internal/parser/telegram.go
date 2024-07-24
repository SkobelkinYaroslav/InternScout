package parser

import (
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strings"
	"time"
	"yarl_intern_bot/internal/readFile"
	"yarl_intern_bot/internal/result"
	"yarl_intern_bot/internal/user"
	"yarl_intern_bot/internal/utils"
)

const (
	defaultTimeout  = time.Second * 5
	defaultMaxTries = 5
)

type Parser struct {
	engine    *colly.Collector
	channels  map[string]struct{}
	users     []*user.User
	parseTime time.Time
	chanData  chan any
	manager   *readFile.FileManager
}

func NewParser(
	c *colly.Collector,
	users []*user.User,
	channels []string,
	parseTime time.Time,
	chanData chan any,
	manager *readFile.FileManager,
) *Parser {
	mpChannels := utils.ArrayToMapStruct(channels)
	p := &Parser{
		engine:    c,
		users:     users,
		channels:  mpChannels,
		parseTime: parseTime,
		chanData:  chanData,
		manager:   manager,
	}
	return p
}

func (p *Parser) Run() {
	interval := p.calculateInterval()
	timer := time.NewTimer(interval)

	for {
		select {
		case msg := <-p.chanData:
			switch msg.(type) {
			case time.Time:
				p.setParseTime(msg.(time.Time))
				timer.Stop()
				interval = p.calculateInterval()
				timer = time.NewTimer(interval)
			case []string:
				log.Printf("Run Adding channels %v\n", msg.([]string))
				p.addChannels(msg.([]string))
			case *user.User:
				p.addUser(msg.(*user.User))
			}
		case <-timer.C:
			results := p.parse()
			p.insertResults(results)
			p.chanData <- p.users
			timer.Stop()
			interval = 24 * time.Hour
			timer = time.NewTimer(interval)
		}
	}
}

func (p *Parser) setParseTime(t time.Time) {
	p.parseTime = t
}
func (p *Parser) getParseTime() time.Time {
	return p.parseTime.UTC()
}

func (p *Parser) addChannels(channels []string) {
	log.Printf("p.channels %v\n", p.channels)
	uniqueChannels := make([]string, 0, len(channels))
	for _, channel := range channels {
		if _, ok := p.channels[channel]; !ok {
			p.channels[channel] = struct{}{}
			uniqueChannels = append(uniqueChannels, channel)
		}
	}
	err := p.manager.AddChannels(uniqueChannels)
	if err != nil {
		log.Println(err)
	}
}

func (p *Parser) addUser(usr *user.User) {
	p.users = append(p.users, usr)
}

func (p *Parser) calculateInterval() time.Duration {
	now := time.Now()
	nextParseTime := time.Date(now.Year(), now.Month(), now.Day(), p.parseTime.Hour(), p.parseTime.Minute(), p.parseTime.Second(), 0, now.Location())
	if nextParseTime.Before(now) {
		nextParseTime = nextParseTime.Add(24 * time.Hour)
	}
	return nextParseTime.Sub(now)
}

func (p *Parser) insertResults(results []result.Result) {
	for _, parsedResult := range results {
		for _, appUser := range p.users {
			if appUser.IsInterested(parsedResult) {
				appUser.AddResults(parsedResult)
			}
		}
	}
}

func (p *Parser) parse() []result.Result {
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	re := regexp.MustCompile(`(?i)<br\s*/?>`)
	responses := make([]result.Result, 0)

	p.engine.OnHTML(".tgme_widget_message_wrap", func(e *colly.HTMLElement) {
		url := e.ChildAttr(".tgme_widget_message_date", "href")

		rawText := e.ChildText(".tgme_widget_message_text")
		text := re.ReplaceAllString(rawText, "\n")

		dateTime := e.ChildAttr(".time", "datetime")

		parsedDateTime, err := time.Parse(time.RFC3339, dateTime)

		localTime := parsedDateTime.Local()

		if err != nil {
			log.Printf("Error parsing datetime %v:%v\n", localTime, err)
			return
		}

		if localTime.Before(startOfToday) {
			return
		}

		curPost := result.New(url, strings.ToLower(text), parsedDateTime)

		responses = append(responses, curPost)
	})

	p.engine.OnError(func(r *colly.Response, err error) {
		log.Printf("Error while sending request to %s: %v", r.Request.URL, err)
		tries := 1
		for err != nil {
			time.Sleep(defaultTimeout * time.Duration(tries))
			if err = r.Request.Retry(); err != nil {
				log.Printf("Error while retrying request to %s: %v", r.Request.URL, err)
			} else {
				log.Printf("Successfully retried request to %s", r.Request.URL)
			}
			if tries == defaultMaxTries {
				return
			}
			tries++
		}
	})

	for channel := range p.channels {
		p.engine.Visit(channel)
	}

	return responses
}
