package telegram

import "GLPITGBOT/telegram/i18n"

const (
	CmdCreateTicket = "CMD_CREATE_TICKET"
	CmdLastTickets  = "CMD_LAST_TICKETS"
)

func resolveCommand(lang, text string) string {
	switch text {
	case i18n.T(lang, "creating_tickets"):
		return CmdCreateTicket
	case i18n.T(lang, "last_tickets"):
		return CmdLastTickets
	default:
		return ""
	}
}
