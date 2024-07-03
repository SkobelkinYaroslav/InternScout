package main

import (
	"github.com/joho/godotenv"
	"yarl_intern_bot/internal/parser"
	"yarl_intern_bot/internal/readFile"
	"yarl_intern_bot/internal/sendResults"
	"yarl_intern_bot/internal/user"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	// read channels list
	channels := readFile.GetChannels()

	//get users and their settings
	users := user.New("config.json")

	// parse tg
	telegramParser := parser.NewTelegramParser(channels)
	results := telegramParser.Telegram()

	// add results to users
	parser.InsertResults(results, users)

	// send results to users
	sendResults.Telegram(users)

}
