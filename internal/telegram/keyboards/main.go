package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var StartKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Создание Заявки"),
		tgbotapi.NewKeyboardButton("Мои Заявки"),
	),
)
