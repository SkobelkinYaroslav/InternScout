package user

import (
	"encoding/json"
	"io"
	"os"
	"yarl_intern_bot/internal/result"
)

type User struct {
	ID         int      `json:"ID"`
	Categories []string `json:"Categories"`
	Results    []string `json:"Results,omitempty"`
}

func New(userFile string) []*User {
	file, err := os.Open(userFile)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)

	if err != nil {
		panic(err)
	}

	var users []*User
	err = json.Unmarshal(byteValue, &users)

	if err != nil {
		panic(err)
	}

	return users

}

func (u *User) AddResults(result result.Result) {
	u.Results = append(u.Results, result.URL)
}
