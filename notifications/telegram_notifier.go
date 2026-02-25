package notifications

import (
	"fmt"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendTelegramMessageToMany(chatIDs []int64, text string) error {
	token := strings.TrimSpace(os.Getenv("TG_BOT_TOKEN"))
	if token == "" {
		return fmt.Errorf("TG_BOT_TOKEN is empty")
	}

	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("message text is empty")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	for _, chatID := range chatIDs {
		if chatID == 0 {
			continue
		}

		msg := tgbotapi.NewMessage(chatID, text)
		if _, err := bot.Send(msg); err != nil {
			return fmt.Errorf("send to chat_id=%d failed: %w", chatID, err)
		}
	}

	return nil
}
