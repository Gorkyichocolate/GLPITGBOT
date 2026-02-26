package notifications

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

func FormatCommentText(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}

	replacer := strings.NewReplacer(
		"<br>", "\n",
		"<br/>", "\n",
		"<br />", "\n",
		"</p>", "\n",
		"</div>", "\n",
	)
	s = replacer.Replace(s)

	s = htmlTagRe.ReplaceAllString(s, "")
	s = html.UnescapeString(s)

	s = strings.ReplaceAll(s, "\u00a0", " ")
	lines := strings.Split(s, "\n")
	cleanLines := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cleanLines = append(cleanLines, line)
	}

	return strings.TrimSpace(strings.Join(cleanLines, "\n"))
}

func normalizeTicketStatusRU(status string) string {
	s := strings.TrimSpace(strings.ToLower(status))

	switch s {
	case "1", "new", "новая", "новый":
		return "Новая"
	case "2", "in_progress", "processing", "в работе":
		return "В работе"
	case "3", "pending", "waiting", "ожидание", "в ожидании":
		return "Ожидание"
	case "4", "solved", "resolved", "решена", "решено":
		return "Решена"
	case "5", "closed", "закрыто", "закрыта":
		return "Закрыто"
	case "6", "cancelled", "canceled", "отменена", "отменено":
		return "Отменена"
	default:
		if status == "" {
			return "Неизвестно"
		}

		return status
	}
}

func BuildTicketCommentNotification(ticketName, commentText, status string) string {
	ticketName = strings.TrimSpace(ticketName)
	if ticketName == "" {
		ticketName = "Без названия"
	}

	commentText = strings.TrimSpace(commentText)
	if commentText == "" {
		commentText = "(пустой комментарий)"
	}

	return fmt.Sprintf(
		"💬 Комментарий к заявке\nЗаявка: %s\nКомментарий: %s\nСтатус: %s",
		ticketName,
		commentText,
		normalizeTicketStatusRU(status),
	)
}

func BuildTicketStatusChangedNotification(ticketName, newStatus string) string {
	ticketName = strings.TrimSpace(ticketName)
	if ticketName == "" {
		ticketName = "Без названия"
	}

	return fmt.Sprintf(
		"🔄 Статус заявки изменён\nЗаявка: %s\nНовый статус: %s",
		ticketName,
		normalizeTicketStatusRU(newStatus),
	)
}

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
