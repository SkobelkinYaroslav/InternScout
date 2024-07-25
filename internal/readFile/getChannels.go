package readFile

import (
	"bufio"
	"os"
	"strings"
)

func GetChannels(channelFile string) []string {
	channels := make([]string, 0)
	file, err := os.Open(channelFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := strings.Replace(scanner.Text(), "https://t.me/", "https://t.me/s/", 1)
		channels = append(channels, str)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return channels
}
