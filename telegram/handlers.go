package telegram

import (
	"GLPITGBOT/i18n"

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

	user := getUser(telegramID)

	if user.Lang == "" {
		user.Lang = i18n.DetectLang(update.Message.From.LanguageCode)
		saveUser(user)
	}

	msg := tgbotapi.NewMessage(chatID, "")

	if text == "" {
		return
	}

	if text == "/start" {
		user.State = StateIdle
		user.ActiveTicket = nil

		msg.Text = i18n.T(user.Lang, "start")
		msg.ReplyMarkup = StartKeyboard(user.Lang)
		bot.Send(msg)
		return
	}

	if text == i18n.T(user.Lang, "btn_language") {
		msg.Text = i18n.T(user.Lang, "choose_language")
		msg.ReplyMarkup = LanguageKeyboard()
		bot.Send(msg)
		return
	}

	if text == i18n.T(user.Lang, "btn_notifications") {
		msg.Text = i18n.T(user.Lang, "notifications")
		msg.ReplyMarkup = NotificationsKeyboard(user.Lang)
		bot.Send(msg)
		return
	}

	if text == i18n.T(user.Lang, "btn_exit") {
		user.State = StateIdle
		user.ActiveTicket = nil
		msg.Text = i18n.T(user.Lang, "exit")
		msg.ReplyMarkup = StartKeyboard(user.Lang)
		bot.Send(msg)
		return
	}

	if text == i18n.T(user.Lang, "preferences") {
		user.State = StateIdle
		user.ActiveTicket = nil

		msg.Text = i18n.T(user.Lang, "preferences")
		msg.ReplyMarkup = PreferencesKeyboard(user.Lang)
		bot.Send(msg)
		return
	}

	createFSM(user, text, &msg, telegramID)

	if msg.Text != "" {
		bot.Send(msg)
	}
}
