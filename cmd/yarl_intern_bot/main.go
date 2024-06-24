package main

import (
	"fmt"
	"yarl_intern_bot/internal/parser"
)

func main() {
	p := parser.New(
		[]string{"https://t.me/s/dvachannel"},
	)

	for _, item := range p.Parse() {
		fmt.Println(item.URL, item.Text, item.Date)
	}
}
