package readFile

import (
	"bufio"
	"os"
	"strings"
)

func GetChannels(channelFile string) ([]string, error) {
	channels := make([]string, 0)
	file, err := os.Open(channelFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		str := strings.Replace(scanner.Text(), "https://t.me/", "https://t.me/s/", 1)
		channels = append(channels, str)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return channels, nil
}

func AddChannels(channels []string, channelFile string) error {
	file, err := os.OpenFile(channelFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, channel := range channels {
		_, err = file.WriteString(channel + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
