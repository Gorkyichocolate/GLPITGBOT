package telegram

import (
	"GLPITGBOT/db"
	"GLPITGBOT/i18n"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func createFSM(
	user *User,
	text string,
	msg *tgbotapi.MessageConfig,
	telegramID int64,
) {
	lang := user.Lang

	switch user.State {

	case StateIdle:
		switch text {

		case i18n.T(lang, "creating_tickets"):
			user.ActiveTicket = &Ticket{}
			user.State = StateWaitTitle
			msg.Text = i18n.T(lang, "enter_ticket_title")

		case i18n.T(lang, "last_tickets"):
			result, err := db.GetLastUserTicketsText(telegramID, 5)
			if err != nil {
				msg.Text = i18n.T(lang, "error_get_tickets")
				return
			}
			msg.Text = result

		default:
			msg.Text = i18n.T(lang, "choose_menu")
		}

	case StateWaitTitle:
		user.ActiveTicket.Title = text
		user.State = StateWaitDescription
		msg.Text = i18n.T(lang, "enter_ticket_description")

	case StateWaitDescription:
		user.ActiveTicket.Description = text

		err := db.CreateTicketByTelegramID(
			telegramID,
			user.ActiveTicket.Title,
			user.ActiveTicket.Description,
		)

		if err != nil {
			msg.Text = i18n.T(lang, "error_create_ticket")
		} else {
			msg.Text = i18n.T(lang, "ticket_created")
		}

		user.State = StateIdle
		user.ActiveTicket = nil
	}
}
