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
