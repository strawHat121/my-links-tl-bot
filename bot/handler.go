package bot

import (
	"encoding/json"
	"log"
	"my-links-bot/db"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	urlRegex = regexp.MustCompile(`https?://[^\s]+`)
	tagRegex = regexp.MustCompile(`#\w+`)
)

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
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	urls := urlRegex.FindAllString(text, -1)

	if len(urls) == 0 {
		return
	}

	rawTags := tagRegex.FindAllString(text, -1)
	tags := []string{}

	for _, tag := range rawTags {
		tags = append(tags, strings.TrimPrefix(strings.ToLower(tag), "#"))
	}

	resourceType := "article"

	if strings.HasPrefix(text, "/save") {
		parts := strings.Fields(text)
		if len(parts) >= 2 && !strings.HasPrefix(parts[1], "http") {
			resourceType = strings.ToLower(parts[1])
		}
	}

	for _, url := range urls {
		log.Printf("Extracted tags: %#v\n", tags)
		err := db.SaveResource(userID, url, resourceType, text, tags)

		if err != nil {
			log.Println("DB error:", err)
			sendMessage(chatID, "Failed to save the link")
		}
	}

	sendMessage(chatID, "Link(s) saved successfully!")
}

func sendMessage(chatID int64, text string) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Println("Failed to send message:", err)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
