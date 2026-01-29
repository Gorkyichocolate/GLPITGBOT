package telegram

import (
	"GLPITGBOT/db"
	"GLPITGBOT/models"
	"GLPITGBOT/telegram/i18n"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	if update.CallbackQuery != nil {
		handleCallback(bot, update)
		return
	}

	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	telegramID := update.Message.From.ID
	text := update.Message.Text

	user, err := db.EnsureUser(
		update.Message.From.ID,
		update.Message.From.UserName,
	)
	if err != nil {
		log.Println("ensure user error:", err)
		return
	}

	if user.Lang == "" {
		user.Lang = i18n.DetectLang(update.Message.From.LanguageCode)
		_ = db.SaveUser(user)
	}

	msg := tgbotapi.NewMessage(chatID, "")

	switch text {

	case "/start":
		user.State = models.StateIdle
		user.ActiveTicket = nil

		msg.Text = i18n.T(user.Lang, "start")
		msg.ReplyMarkup = StartKeyboard(user.Lang)

	case i18n.T(user.Lang, "btn_language"):
		msg.Text = i18n.T(user.Lang, "choose_language")
		msg.ReplyMarkup = LanguageKeyboard()

	case i18n.T(user.Lang, "btn_notifications"):
		msg.Text = i18n.T(user.Lang, "notifications")
		msg.ReplyMarkup = NotificationsKeyboard(user.Lang)

	case i18n.T(user.Lang, "btn_exit"):
		user.State = models.StateIdle
		user.ActiveTicket = nil

		msg.Text = i18n.T(user.Lang, "exit")
		msg.ReplyMarkup = StartKeyboard(user.Lang)

	case i18n.T(user.Lang, "preferences"):
		user.State = models.StateIdle
		user.ActiveTicket = nil

		msg.Text = i18n.T(user.Lang, "preferences")
		msg.ReplyMarkup = PreferencesKeyboard(user.Lang)

	default:
		command := resolveCommand(user.Lang, text)
		createFSM(user, command, text, &msg, telegramID)
	}

	if msg.Text == "" {
		msg.Text = "⚠️ Неизвестная команда"
		msg.ReplyMarkup = StartKeyboard(user.Lang)
	}

	_ = db.SaveUser(user)

	if _, err := bot.Send(msg); err != nil {
		log.Println("SEND ERROR:", err)
	}
}
