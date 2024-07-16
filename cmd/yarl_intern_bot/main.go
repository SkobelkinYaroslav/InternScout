package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"time"
	"yarl_intern_bot/internal/parser"
	"yarl_intern_bot/internal/readFile"
	"yarl_intern_bot/internal/sendResults"
	"yarl_intern_bot/internal/user"
)

func main() {
	now := time.Now()
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	// read channels list
	channels := readFile.GetChannels("channels.txt")

	//get users and their settings
	users := user.New("config.json")

	// parse tg
	telegramParser := parser.NewTelegramParser(channels)
	results := telegramParser.Telegram()

	// add results to users
	parser.InsertResults(results, users)

	// send results to users
	sendResults.Telegram(users)

	fmt.Printf("%d posts were processed in %.3f", len(results), time.Since(now).Minutes())

}
