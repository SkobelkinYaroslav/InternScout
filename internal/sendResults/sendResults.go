package sendResults

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"yarl_intern_bot/internal/user"
)

const (
	defaultTimeout  = time.Second * 5
	defaultMaxTries = 5
)

func Telegram(users []*user.User) {
	key := os.Getenv("API_KEY")
	if key == "" {
		log.Fatalln("API_KEY environment variable not set")
	}

	for _, user := range users {
		resultsString := "Сегодня ничего нет :("
		if len(user.Results) > 0 {
			arr := make([]string, 0, len(user.Results))
			for key := range user.Results {
				arr = append(arr, key)
			}
			resultsString = "Вот что найдено за сегодня: \n" + strings.Join(arr, "\n")
		}

		encodedResults := url.QueryEscape(resultsString)

		tries := 1
		url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", key, user.ID, encodedResults)
		resp, err := http.Get(url)
		for resp.StatusCode != http.StatusOK {
			time.Sleep(defaultTimeout * time.Duration(tries))
			resp, err = http.Get(url)
			if err != nil {
				log.Println("Error while retrying request:", err)
				break
			}
			if tries == defaultMaxTries {
				return
			}
			tries++

		}
		if err != nil {
			log.Println("Error making request:", err)
			continue
		}
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Println("Error decoding response:", err)
			continue
		}
		resp.Body.Close()

		fmt.Println(result)
	}
}
