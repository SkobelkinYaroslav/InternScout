package readFile

import (
	"bufio"
	"log"
	"os"
)

func GetChannels() []string {
	channels := make([]string, 0)
	file, err := os.Open("channels.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		channels = append(channels, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return channels
}
