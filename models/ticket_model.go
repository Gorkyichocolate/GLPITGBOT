package models

import "time"

type Ticket struct {
	Id          int
	Title       string
	Description string
	UserID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      string
}
