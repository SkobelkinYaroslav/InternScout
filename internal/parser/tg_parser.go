package parser

import (
	"github.com/gocolly/colly"
	"log"
	"strings"
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

func isToday(postDate time.Time) bool {
	now := time.Now().In(postDate.Location())

	postYear, postMonth, postDay := postDate.Date()
	currentYear, currentMonth, currentDay := now.Date()

	return postYear == currentYear && postMonth == currentMonth && postDay == currentDay
}

func (p Parser) Parse() []result.Result {
	responses := make([]result.Result, 0)

	//FixMe: <br> не добавляется в текст. Просто идёт подряд
	p.engine.OnHTML(".tgme_widget_message_wrap", func(e *colly.HTMLElement) {
		url := e.ChildAttr(".tgme_widget_message_date", "href")

		// .First() ?
		var sb strings.Builder
		//text := e.ChildText(".tgme_widget_message_text")
		e.ForEach(".tgme_widget_message_text", func(_ int, elem *colly.HTMLElement) {
			elem.ForEach("br", func(_ int, br *colly.HTMLElement) {
				sb.WriteString("\n")
			})
			text := strings.TrimSpace(elem.Text)
			sb.WriteString(text)
			sb.WriteString("\n") // Добавляем перенос строки после каждого элемента .tgme_widget_message_text
		})

		text := sb.String()

		dateTime := e.ChildAttr(".time", "datetime")

		parsedDateTime, err := time.Parse(time.RFC3339, dateTime)
		if err != nil {
			log.Printf("Error parsing datetime: %v\n", err)
			return
		}

		if !isToday(parsedDateTime) {
			log.Println("Post is outdated", parsedDateTime)
			return
		}

		curPost := result.New(url, text, parsedDateTime)

		responses = append(responses, curPost)
	})

	for _, channel := range p.channels {
		p.engine.Visit(channel)
	}

	return responses
}
