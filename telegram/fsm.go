package telegram

import (
	"GLPITGBOT/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func createFSM(
	ticket *Ticket,
	text string,
	msg *tgbotapi.MessageConfig,
	chatID int64,
) {
	switch ticket.State {

	case StateIdle:
		switch text {

		case "Создание Заявки":
			msg.Text = "Введите Тему Заявки"
			ticket.State = StateWaitTitle

		case "Последние Заявки":
			result, err := db.GetLastUserTicketsText(chatID, 5)
			if err != nil {
				msg.Text = "Ошибка при получении заявок"
				return
			}

			msg.Text = result

		default:
			msg.Text = "Выберите пункт в меню"
		}

	case StateWaitTitle:
		ticket.Title = text
		msg.Text = "Введите Описание Заявки"
		ticket.State = StateWaitDescription

	case StateWaitDescription:
		ticket.Description = text

		err := db.CreateTicketByTelegramID(
			chatID,
			ticket.Title,
			ticket.Description,
		)

		if err != nil {
			msg.Text = "Ошибка при создании заявки"
			ticket.State = StateIdle
			return
		}

		msg.Text = "✅ Заявка успешно создана"
		ticket.State = StateIdle

		ticket.Title = ""
		ticket.Description = ""
	}
}
