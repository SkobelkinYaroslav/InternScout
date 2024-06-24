package parser

import (
	"github.com/gocolly/colly"
	"log"
	"time"
	"yarl_intern_bot/internal/result"
)

type Parser struct {
	engine   *colly.Collector
	channels []string
}

func New(channels []string) Parser {
	return Parser{
		engine: colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"),
		),
		channels: channels,
	}
}

func (p Parser) Parse() []result.Result {
	responses := make([]result.Result, 0)

	p.engine.OnHTML(".tgme_widget_message", func(e *colly.HTMLElement) {
		url := e.ChildAttr("a", "href")
		text := e.DOM.Find(".tgme_widget_message_text").First().Text()
		dateTime := e.ChildAttr("time", "datetime")
		parsedDateTime, err := time.Parse(time.RFC3339, dateTime)
		if err != nil {
			log.Printf("Error parsing datetime: %v\n", err)
			return
		}

		curPost := result.Result{
			URL:  url,
			Text: text,
			Date: parsedDateTime,
		}

		responses = append(responses, curPost)
	})

	for _, channel := range p.channels {
		p.engine.Visit(channel)
	}

	return responses
}
