package repository

import (
	"context"
	"time"
)

type Ticket struct {
	ID          int64
	TicketID    string
	UserID      string
	Title       string
	Description string
	Status      string
	Active      bool
	CreatedAt   time.Time
}
type User struct {
	ID         int64
	UserID     string // ← ОБЯЗАТЕЛЬНО
	TelegramID int64  // ← ОБЯЗАТЕЛЬНО
	Name       string
	Surname    string
	Company    string
	Status     string
}

// Добавление тикета
func AddTicket(ticket Ticket) error {
	_, err := Pool.Exec(
		context.Background(),
		`INSERT INTO tickets
		 (ticket_id, user_id, title, description, status, active)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		ticket.TicketID,
		ticket.UserID,
		ticket.Title,
		ticket.Description,
		ticket.Status,
		ticket.Active,
	)
	return err
}

// Получение последних 5 тикетов пользователя
func GetLastFiveTickets(userID string) ([]Ticket, error) {
	rows, err := Pool.Query(
		context.Background(),
		`SELECT id, ticket_id, user_id, title, description, status, active, created_at
		 FROM tickets
		 WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT 5`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []Ticket

	for rows.Next() {
		var t Ticket
		if err := rows.Scan(
			&t.ID,
			&t.TicketID,
			&t.UserID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Active,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, nil
}

// Смена статуса заявки
func UpdateTicketStatus(ticketID, status string) error {
	_, err := Pool.Exec(
		context.Background(),
		`UPDATE tickets SET status = $1 WHERE ticket_id = $2`,
		status, ticketID,
	)
	return err
}

// Получить тикет по айдишке
func GetTicketByTicketID(ticketID string) (*Ticket, error) {
	var t Ticket

	err := Pool.QueryRow(
		context.Background(),
		`SELECT id, ticket_id, user_id, title, description, status, active, created_at
		 FROM tickets
		 WHERE ticket_id = $1`,
		ticketID,
	).Scan(
		&t.ID,
		&t.TicketID,
		&t.UserID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Active,
		&t.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Получить Юзера
func GetUserByUserID(userID string) (*User, error) {
	var u User

	err := Pool.QueryRow(
		context.Background(),
		`SELECT id, user_id, telegram_id, name, surname, company, status
		 FROM users
		 WHERE user_id = $1`,
		userID,
	).Scan(
		&u.ID,
		&u.UserID,
		&u.TelegramID,
		&u.Name,
		&u.Surname,
		&u.Company,
		&u.Status,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
