package result

import "time"

type Result struct {
	URL  string
	Text string
	Date time.Time
}

func New(url, text string, date time.Time) Result {
	return Result{
		URL:  url,
		Text: text,
		Date: date,
	}
}
