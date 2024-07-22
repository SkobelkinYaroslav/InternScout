package readFile

import (
	"encoding/json"
	"io"
	"os"
	"yarl_intern_bot/internal/user"
)

func GetUsers(userFile string) ([]*user.User, error) {
	file, err := os.Open(userFile)

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

func AddUsers(users []*user.User, userFile string) error {
	for _, user := range users {
		err := AddUser(user, userFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddUser(usr *user.User, userFile string) error {
	file, err := os.OpenFile(userFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := json.Marshal(usr)

	if err != nil {
		return err
	}

	err = os.WriteFile(userFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
