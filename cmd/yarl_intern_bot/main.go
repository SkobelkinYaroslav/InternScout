package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"time"
	"yarl_intern_bot/internal/parser"
	"yarl_intern_bot/internal/readFile"
	"yarl_intern_bot/internal/sendResults"
	"yarl_intern_bot/internal/user"
)

func main() {
	now := time.Now()

	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	execDir := filepath.Dir(execPath)

	err = godotenv.Load(execDir + "/.env")

	if err != nil {
		panic(err)
	}

	// read channels list
	channels := readFile.GetChannels(execDir + "/channels.txt")

	//get users and their settings
	users := user.New(execDir + "/config.json")

	// parse tg
	telegramParser := parser.NewTelegramParser(channels)
	results := telegramParser.Telegram()

	// add results to users
	parser.InsertResults(results, users)

	// send results to users
	sendResults.Telegram(users)

	fmt.Printf("%d posts were processed in %.3f", len(results), time.Since(now).Minutes())

}
