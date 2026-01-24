package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var start = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Создание Заявки"),
		tgbotapi.NewKeyboardButton("Последние Заявки"),
	),
)
