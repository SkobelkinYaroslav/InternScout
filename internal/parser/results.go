package parse

import (
	"strings"
	"yarl_intern_bot/internal/result"
	"yarl_intern_bot/internal/user"
)

func Results(parsedResults []result.Result, appUsers []*user.User) {
	for _, parsedResult := range parsedResults {
		for _, appUser := range appUsers {
			for _, keyword := range appUser.Categories {
				if strings.Contains(parsedResult.Text, keyword) {
					appUser.AddResults(parsedResult)
				}
			}

		}
	}
}
