package db

import "time"

type User struct {
	Id         int
	ApiToken   string
	Username   string
	TelegramId int64
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
