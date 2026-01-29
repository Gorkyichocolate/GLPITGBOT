package db

import "time"

type User struct {
	Id         int
	TelegramId int64
	Username   string
	ApiToken   string
	Lang       string
}

type Ticket struct {
	Id          int
	Title       string
	Description string
	UserID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      string
}
