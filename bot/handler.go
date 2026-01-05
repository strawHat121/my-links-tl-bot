package bot

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"my-links-bot/db"
	"my-links-bot/models"
	"my-links-bot/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	urlRegex = regexp.MustCompile(`https?://[^\s]+`)
	tagRegex = regexp.MustCompile(`#\w+`)
	cache    = map[int64][]models.Resource{}
)

func Handle(body []byte) {
	var update tgbotapi.Update
	json.Unmarshal(body, &update)

	if update.Message == nil {
		return
	}

	text := update.Message.Text

	if strings.HasPrefix(text, "/list") {
		handleList(update)
		return
	}

	if strings.HasPrefix(text, "/done") {
		handleDone(update)
		return
	}

	handleSave(update)
}

func handleSave(update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	urls := urlRegex.FindAllString(text, -1)
	if len(urls) == 0 {
		return
	}

	tags := []string{}
	for _, t := range tagRegex.FindAllString(text, -1) {
		tags = append(tags, strings.TrimPrefix(strings.ToLower(t), "#"))
	}

	resType := "article"
	parts := strings.Fields(text)
	if strings.HasPrefix(text, "/save") && len(parts) > 1 && !strings.HasPrefix(parts[1], "http") {
		resType = parts[1]
	}

	for _, url := range urls {
		title := util.ExtractTitle(url)
		r := models.Resource{
			UserID: userID,
			Type:   resType,
			Title:  title,
			URL:    url,
			Status: "to_read",
			Tags:   tags,
			Notes:  text,
		}
		db.SaveResource(r)
	}

	send(chatID, "✅ Saved")
}

func handleList(update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	items, _ := db.ListResources(userID, 5)
	cache[userID] = items

	var b strings.Builder
	b.WriteString("📚 <b>Your saved items</b>\n\n")

	for i, r := range items {
		b.WriteString(fmt.Sprintf(
			"%d. <a href=\"%s\">%s</a>\n",
			i+1, r.URL, htmlEscape(r.Title),
		))
	}

	msg := tgbotapi.NewMessage(chatID, b.String())
	msg.ParseMode = "HTML"
	bot, _ := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	bot.Send(msg)
}

func handleDone(update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	parts := strings.Fields(update.Message.Text)

	if len(parts) != 2 {
		send(chatID, "Usage: /done <number>")
		return
	}

	i, _ := strconv.Atoi(parts[1])
	items := cache[userID]
	if i < 1 || i > len(items) {
		send(chatID, "Invalid item")
		return
	}

	db.MarkDone(userID, items[i-1].SK)
	send(chatID, "✅ Marked as completed")
}

func send(chatID int64, text string) {
	bot, _ := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	bot.Send(tgbotapi.NewMessage(chatID, text))
}

func htmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return r.Replace(s)
}
