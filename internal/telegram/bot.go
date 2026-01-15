package telegram

import (
	"GLPITGBOT/internal/config"
	"GLPITGBOT/internal/telegram/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func StartBot(cfg *config.Config) {
	bot, err := tgbotapi.NewBotAPI(cfg.TGBotToken)
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = "Выберите действие"
				msg.ReplyMarkup = keyboards.StartKeyboard
			}
		} else {
			switch update.Message.Text {
			case "Создание Заявки":
				msg.Text = "Создание заявки..."
			case "Мои Заявки":
				msg.Text = "Ваши заявки"
			}
		}

		bot.Send(msg)
	}
}
