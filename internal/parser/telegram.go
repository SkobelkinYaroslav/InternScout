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

func NewTelegramParser(channels []string) TelegramParser {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"),
	)

	// Устанавливаем задержку между запросами
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       5 * time.Second,
		RandomDelay: 1 * time.Second,
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

		html, err := e.DOM.Find(".tgme_widget_message_text").Html()
		if err != nil {
			log.Printf("Error getting HTML from %s: %v", e.Request.URL, err)
			return
		}

		text := re.ReplaceAllString(html, "\n")

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
		time.Sleep(10 * time.Second)
		if err := r.Request.Retry(); err != nil {
			log.Printf("Error while retrying request to %s: %v", r.Request.URL, err)
			return
		}
	})

	for _, channel := range p.channels {
		p.engine.Visit(channel)
	}

	return responses
}
