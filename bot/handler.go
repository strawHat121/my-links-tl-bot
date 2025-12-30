package bot

import (
	"encoding/json"
	"log"
	"my-links-bot/db"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

func Handle(body []byte) {
	var update tgbotapi.Update

	if err := json.Unmarshal(body, &update); err != nil {
		log.Println("Invalid update:", err)
		return
	}

	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	text := update.Message.Text

	urls := urlRegex.FindAllString(text, -1)

	if len(urls) == 0 {
		return
	}

	for _, url := range urls {
		err := db.SaveResource(userID, url, text)

		if err != nil {
			log.Println("DB error:", err)
		}
	}
}
