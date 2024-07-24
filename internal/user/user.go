package user

import (
	"encoding/json"
	"io"
	"os"
	"yarl_intern_bot/internal/result"
)

type User struct {
	ID         int                 `json:"ID"`
	Categories []string            `json:"Categories"`
	Results    map[string]struct{} `json:"Results,omitempty"`
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

	for i := range users {
		users[i].Results = make(map[string]struct{})
	}

	return users

}

func (u *User) AddResults(result result.Result) {
	u.Results[result.URL] = struct{}{}
}

func (u *User) IsInterested(result result.Result) bool {
	for _, category := range u.Categories {
		if category == result.Text {
			return true
		}
	}
	return false
}
