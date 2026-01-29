package telegram

import (
	"GLPITGBOT/telegram/i18n"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartKeyboard(lang string) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.T(lang, "creating_tickets")),
			tgbotapi.NewKeyboardButton(i18n.T(lang, "last_tickets")),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.T(lang, "preferences")),
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

func PreferencesKeyboard(lang string) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.T(lang, "btn_notifications")),
			tgbotapi.NewKeyboardButton(i18n.T(lang, "btn_exit")),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.T(lang, "btn_language")),
		),
	)
}

func LanguageKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang_ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang_en"),
			tgbotapi.NewInlineKeyboardButtonData("🇰🇿 Қазақша", "lang_kk"),
		),
	)
}

func RemoveKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
	}
}
