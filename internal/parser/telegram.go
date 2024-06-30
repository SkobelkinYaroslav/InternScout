package parse

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
	return TelegramParser{
		engine: colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"),
		),
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
			log.Printf("Error getting HTML: %v\n", err)
			return
		}

		text := re.ReplaceAllString(html, "\n")

		dateTime := e.ChildAttr(".time", "datetime")

		parsedDateTime, err := time.Parse(time.RFC3339, dateTime)

		localTime := parsedDateTime.Local()

		if err != nil {
			log.Printf("Error parsing datetime: %v\n", err)
			return
		}

		if localTime.Before(startOfToday) {
			log.Println("Post is outdated:", localTime)
			return
		}

		curPost := result.New(url, text, parsedDateTime)

		responses = append(responses, curPost)
	})

	for _, channel := range p.channels {
		p.engine.Visit(channel)
		time.Sleep(10)
	}

	return responses
}
