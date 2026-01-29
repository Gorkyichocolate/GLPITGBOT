package telegram

import (
	"GLPITGBOT/db"
	"GLPITGBOT/telegram/i18n"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	q := update.CallbackQuery

	user, err := db.EnsureUser(
		q.From.ID,
		q.From.UserName,
	)
	if err != nil {
		log.Println("ensure user error:", err)
		return
	}

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

	if err := db.SaveUser(user); err != nil {
		log.Println("save user error:", err)
	}

	bot.Send(tgbotapi.NewCallback(q.ID, ""))

	msg := tgbotapi.NewMessage(
		q.Message.Chat.ID,
		i18n.T(user.Lang, "preferences"),
	)
	msg.ReplyMarkup = PreferencesKeyboard(user.Lang)

	if _, err := bot.Send(msg); err != nil {
		log.Println("send msg error:", err)
	}
}
