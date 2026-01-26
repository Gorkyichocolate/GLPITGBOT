package telegram

import (
	"GLPITGBOT/i18n"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartKeyboard(lang string) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.T(lang, "btn_create_ticket")),
			tgbotapi.NewKeyboardButton(i18n.T(lang, "btn_last_tickets")),
		),
	)
}

func NotificationsKeyboard(lang string) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.T(lang, "btn_notifications")),
			tgbotapi.NewKeyboardButton(i18n.T(lang, "btn_language")),
			tgbotapi.NewKeyboardButton(i18n.T(lang, "btn_exit")),
		),
	)
}
