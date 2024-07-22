package parser

import (
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"time"
	"yarl_intern_bot/internal/result"
	"yarl_intern_bot/internal/user"
)

const (
	defaultTimeout  = time.Second * 5
	defaultMaxTries = 5
)

type Parser struct {
	engine    *colly.Collector
	channels  []string
	users     []*user.User
	parseTime time.Time
	chanData  chan any
}

func NewParser(users []*user.User, channels []string, parseTime time.Time, chanData chan any) *Parser {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Delay:      defaultTimeout,
	})

	p := &Parser{
		engine:    c,
		users:     users,
		channels:  channels,
		parseTime: parseTime,
		chanData:  chanData,
	}
	return p
}

func (p *Parser) setParseTime(t time.Time) {
	p.parseTime = t
}
func (p *Parser) getParseTime() time.Time {
	return p.parseTime.UTC()
}

func (p *Parser) addChannels(channels []string) {
	p.channels = append(p.channels, channels...)
}

func (p *Parser) addUser(usr *user.User) {
	p.users = append(p.users, usr)
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
				p.addChannels(msg.([]string))
			case *user.User:
				p.addUser(msg.(*user.User))
			}
		case <-timer.C:
			p.Parse()
			timer.Stop()
			interval = 24 * time.Hour
			timer = time.NewTimer(interval)
		}
	}
}

func (p *Parser) calculateInterval() time.Duration {
	now := time.Now()
	nextParseTime := time.Date(now.Year(), now.Month(), now.Day(), p.parseTime.Hour(), p.parseTime.Minute(), p.parseTime.Second(), 0, now.Location())
	if nextParseTime.Before(now) {
		nextParseTime = nextParseTime.Add(24 * time.Hour)
	}
	return nextParseTime.Sub(now)
}

func (p *Parser) Parse() []result.Result {
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

		curPost := result.New(url, text, parsedDateTime)

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

	for _, channel := range p.channels {
		p.engine.Visit(channel)
	}

	return responses
}
