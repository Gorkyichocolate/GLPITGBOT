package telegram

import "time"

type Ticket struct {
	Title       string
	Description string
	CreatedAt   time.Time
	State       State
}

type User struct {
	ID           int64
	Lang         string
	State        State
	ActiveTicket *Ticket
}
