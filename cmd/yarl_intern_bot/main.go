package main

import (
	parse "yarl_intern_bot/internal/parser"
	"yarl_intern_bot/internal/readFile"
	"yarl_intern_bot/internal/user"
)

func main() {
	// read channels list
	channels := readFile.GetChannels()

	//get users and their settings
	users := user.New("config.json")

	// parse tg
	tgParser := parse.NewTelegramParser(channels)
	results := tgParser.Telegram()

	// add results to users
	parse.ParseResults(results, users)

	// send results to users

}
