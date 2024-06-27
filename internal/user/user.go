package user

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type User struct {
	ID         int      `json:"ID"`
	Categories []string `json:"Categories"`
	Results    []string `json:"Results,omitempty"`
}

func New(userFile string) []User {
	file, err := os.Open(userFile)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)

	if err != nil {
		panic(err)
	}

	var users []User
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		log.Fatalf("failed to unmarshal json: %s", err)
	}
	for _, user := range users {
		fmt.Printf("ID: %d, Categories: %v, Results: %v\n", user.ID, user.Categories, user.Results)
	}

	return users

}
