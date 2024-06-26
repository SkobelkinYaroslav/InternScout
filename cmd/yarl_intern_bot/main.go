package main

import (
	"fmt"
	"yarl_intern_bot/internal/readFile"
)

func main() {
	// read channels list
	channels := readFile.GetChannels()

	//get users and their settings

	// parse tg

	// send results to users

	for _, channel := range channels {
		fmt.Println(channel)
	}
}
