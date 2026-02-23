package telegram

import (
	"GLPITGBOT/models"
	"sync"
)

type sessionData struct {
	Step         string
	ActiveTicket *models.Ticket
}

var userSessions sync.Map

const (
	StepIdle               = ""
	StepWaitApiToken       = "wait_api_token"
	StepWaitTicketTitle    = "wait_ticket_title"
	StepWaitTicketDesc     = "wait_ticket_description"
	StepWaitTicketPriority = "wait_ticket_priority"
	StepWaitTicketDueDate  = "wait_ticket_due_date"
)

func getSession(telegramID int64) sessionData {
	v, ok := userSessions.Load(telegramID)
	if !ok {
		return sessionData{Step: StepIdle}
	}

	s, ok := v.(sessionData)
	if !ok {
		return sessionData{Step: StepIdle}
	}

	return s
}

func setSession(telegramID int64, data sessionData) {
	userSessions.Store(telegramID, data)
}

func resetSession(telegramID int64) {
	userSessions.Delete(telegramID)
}
