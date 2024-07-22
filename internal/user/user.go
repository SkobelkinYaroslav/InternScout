package user

import (
	"yarl_intern_bot/internal/result"
)

type User struct {
	ID         int                 `json:"ID"`
	Categories []string            `json:"Categories"`
	Results    map[string]struct{} `json:"Results,omitempty"`
}

func New(id int, categories []string) *User {
	return &User{
		ID:         id,
		Categories: categories,
	}
}

func (u *User) AddResults(result result.Result) {
	u.Results[result.URL] = struct{}{}
}
