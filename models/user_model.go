package models

type User struct {
	ID         int
	TelegramID int64
	Username   string
	ApiToken   string
	Lang       string
	SessionToken string

	State        State
	ActiveTicket *Ticket
}
