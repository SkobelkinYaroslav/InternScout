package main

import (
	"yarl_intern_bot/internal/readFile"
	"yarl_intern_bot/internal/user"
)

func main() {
	// read channels list
	channels := readFile.GetChannels()

	//get users and their settings
	users := user.New("config.json")

	// parse tg

	// send results to users

	//for _, channel := range channels {
	//	fmt.Println(channel)
	//}
}
