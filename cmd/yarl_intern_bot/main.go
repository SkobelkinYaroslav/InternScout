package main

import (
	"github.com/joho/godotenv"
	parse "yarl_intern_bot/internal/parser"
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
	tgParser := parse.NewTelegramParser(channels)
	results := tgParser.Telegram()

	// add results to users
	parse.Results(results, users)

	// send results to users
	sendResults.Telegram(users)

}
