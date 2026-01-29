package telegram

import (
	"GLPITGBOT/i18n"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	q := update.CallbackQuery
	user := getUser(q.From.ID)

	switch q.Data {
	case "lang_ru":
		user.Lang = "ru"
	case "lang_en":
		user.Lang = "en"
	case "lang_kk":
		user.Lang = "kk"
	default:
		return
	}

	saveUser(user)
	_, err := bot.Request(tgbotapi.NewCallback(q.ID, ""))
	if err != nil {
		log.Println("callback error:", err)
	}

	msg := tgbotapi.NewMessage(
		q.Message.Chat.ID,
		i18n.T(user.Lang, "preferences"),
	)
	msg.ReplyMarkup = PreferencesKeyboard(user.Lang)

	bot.Send(msg)
}
