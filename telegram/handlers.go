package telegram

import (
	"GLPITGBOT/i18n"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	user := getUser(chatID)
	text := update.Message.Text

	lang := user.Lang

	user.Lang = "ru"
	user.State = StateCreatingTicket
	user.ActiveTicket = &Ticket{}

	msg := tgbotapi.NewMessage(chatID, "")

	if text == "/start" {

		if user.Lang == "" {
			user.Lang = i18n.DetectLang(update.Message.From.LanguageCode)
			saveUser(user)
		}

		msg.Text = i18n.T(user.Lang, "start")
		msg.ReplyMarkup = StartKeyboard(user.Lang)
		bot.Send(msg)
		return
	}

	if text == "/auth" {
		msg.Text = i18n.T(lang, "enter_api_key")
		bot.Send(msg)
		return
	}

	if text == i18n.T(lang, "btn_language") {
		msg.Text = i18n.T(lang, "choose_language")
		msg.ReplyMarkup = LanguageKeyboard()
		bot.Send(msg)
		return
	}

	if text == "/notifications" {
		msg.Text = i18n.T(lang, "notifications")
		msg.ReplyMarkup = NotificationsKeyboard(lang)
		bot.Send(msg)
		return
	}

	if text == "/preferences" {
		msg.Text = i18n.T(lang, "preferences")
		msg.ReplyMarkup = PreferencesKeyboard(lang)
		bot.Send(msg)
	}

	createFSM(user, text, &msg, chatID)

	if msg.Text != "" {
		bot.Send(msg)
	}
}
