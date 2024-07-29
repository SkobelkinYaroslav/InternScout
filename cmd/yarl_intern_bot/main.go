package main

import (
	"context"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"time"
	"yarl_intern_bot/internal/parser"
	"yarl_intern_bot/internal/readFile"
	"yarl_intern_bot/internal/telegram"
)

func main() {
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	execDir := filepath.Dir(execPath)

	err = godotenv.Load(execDir + "/.env")
	if err != nil {
		panic(err)
	}

	fr := readFile.NewFileManager(filepath.Join(execDir, "config.json"), filepath.Join(execDir, "channels.txt"))

	// read channels list
	channels, err := fr.GetChannels()
	if err != nil {
		panic(err)
	}
	//get users and their settings
	users, err := fr.GetUsers()
	if err != nil {
		panic(err)
	}

	chanData := make(chan any)
	timeString := "15:04"
	parsedTime, err := time.Parse("15:04", timeString)
	if err != nil {
		panic(err)
	}

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Delay:      5 * time.Second,
	})

	p := parser.NewParser(c, users, channels, parsedTime, chanData, fr)
	go p.Run()

	apiKey := os.Getenv("API_KEY")
	tg := telegram.New(context.Background(), apiKey, chanData)
	go tg.Run()

	select {}

}
