package sendResults

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"yarl_intern_bot/internal/user"
)

func Telegram(users []*user.User) {
	key := os.Getenv("API_KEY")
	if key == "" {
		log.Fatalln("API_KEY environment variable not set")
	}

	for _, user := range users {
		url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", key, user.ID, user.Results)
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Error making request:", err)
			continue
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Println("Error decoding response:", err)
			continue
		}

		fmt.Println(result)
	}
}
