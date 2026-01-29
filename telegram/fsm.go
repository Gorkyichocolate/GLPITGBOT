package telegram

import (
	"GLPITGBOT/db"
	"GLPITGBOT/models"
	"GLPITGBOT/telegram/i18n"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func createFSM(
	user *models.User,
	command string,
	text string,
	msg *tgbotapi.MessageConfig,
	telegramID int64,
) {
	lang := user.Lang

	switch user.State {

	case models.StateIdle:

		switch command {

		case CmdCreateTicket:
			user.ActiveTicket = &models.Ticket{}
			user.State = models.StateWaitTitle
			msg.Text = i18n.T(lang, "enter_ticket_title")
			msg.ReplyMarkup = RemoveKeyboard()

		case CmdLastTickets:
			result, _ := db.GetLastUserTicketsText(telegramID, 5)
			msg.Text = result
			msg.ReplyMarkup = StartKeyboard(lang)

		default:
			msg.Text = i18n.T(lang, "choose_menu")
			msg.ReplyMarkup = StartKeyboard(lang)
		}

	case models.StateWaitTitle:
		user.ActiveTicket.Title = text
		user.State = models.StateWaitDescription
		msg.Text = i18n.T(lang, "enter_ticket_description")

	case models.StateWaitDescription:
		user.ActiveTicket.Description = text
		db.CreateTicketByTelegramID(
			telegramID,
			user.ActiveTicket.Title,
			user.ActiveTicket.Description,
		)
		msg.Text = i18n.T(lang, "ticket_created")
		user.State = models.StateIdle
		user.ActiveTicket = nil
		msg.ReplyMarkup = StartKeyboard(lang)
	}
}
