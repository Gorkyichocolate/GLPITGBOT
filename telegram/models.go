package telegram

import "time"

type Ticket struct {
	State State

	Title       string
	Description string
	CreatedAt   time.Time
}
