package readFile

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"
	"yarl_intern_bot/internal/user"
)

type FileManager struct {
	userFile    string
	channelFile string
}

func NewFileManager(userFile, channelFile string) *FileManager {
	return &FileManager{
		userFile:    userFile,
		channelFile: channelFile,
	}
}

func (f *FileManager) GetChannels() ([]string, error) {
	channels := make([]string, 0)
	file, err := os.Open(f.channelFile)
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

func (f *FileManager) AddChannels(channels []string) error {
	file, err := os.OpenFile(f.channelFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, channel := range channels {
		if _, err := file.WriteString(channel + "\n"); err != nil {
			return err
		}
	}
	return nil
}
func (f *FileManager) GetUsers() ([]*user.User, error) {
	file, err := os.Open(f.userFile)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	var users []*user.User
	err = json.Unmarshal(byteValue, &users)

	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].Results = make(map[string]struct{})
	}

	return users, nil
}

// TODO: adding user to file is not valid
func (f *FileManager) AddUsers(users []*user.User) error {
	for _, user := range users {
		err := f.AddUser(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileManager) AddUser(usr *user.User) error {
	file, err := os.OpenFile(f.userFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := json.Marshal(usr)

	if err != nil {
		return err
	}

	err = os.WriteFile(f.userFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
