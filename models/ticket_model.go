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
	AppToken    string
	DueDate     time.Time
	Priority    int
	EntitiesId  int
}

type CreateTicketInput struct {
	Name       string `json:"name"`
	Content    string `json:"content"`
	EntitiesID int    `json:"entities_id"`
	Status     int    `json:"status"`
	Priority   int    `json:"priority"`
	DueDate    string `json:"due_date"`
}
