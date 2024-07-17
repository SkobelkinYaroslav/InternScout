package parser

import (
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"time"
	"yarl_intern_bot/internal/result"
)

type TelegramParser struct {
	engine   *colly.Collector
	channels []string
}

const (
	defaultTimeout  = time.Second * 5
	defaultMaxTries = 5
)

func NewTelegramParser(channels []string) TelegramParser {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Delay:      defaultTimeout,
	})

	return TelegramParser{
		engine:   c,
		channels: channels,
	}
}
func (p TelegramParser) Telegram() []result.Result {
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
