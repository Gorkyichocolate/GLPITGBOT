package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	user := getUser(chatID)
	text := update.Message.Text

	msg := tgbotapi.NewMessage(chatID, "")

	if text == "/start" {
		user.State = StateIdle
		msg.Text = "startuem"
		msg.ReplyMarkup = start
		bot.Send(msg)
		return
	}

	createFSM(user, text, &msg, chatID)

	if msg.Text != "" {
		bot.Send(msg)
	}
}
