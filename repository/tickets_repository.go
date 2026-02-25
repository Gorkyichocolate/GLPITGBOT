package repository

import (
	"GLPITGBOT/db"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	TicketStatusWaiting   = "Ожидании"
	TicketStatusReview    = "В рассмотрении"
	TicketStatusResolved  = "Решена"
	TicketStatusCompleted = "Завершена"
)

func normalizeTicketStatus(status string) string {
	status = strings.TrimSpace(status)

	switch status {
	case TicketStatusWaiting, TicketStatusReview, TicketStatusResolved, TicketStatusCompleted:
		return status
	case "open", "pending", "pending_glpi", "glpi_failed", "new", "":
		return TicketStatusWaiting
	case "in_progress":
		return TicketStatusReview
	case "closed", "done":
		return TicketStatusCompleted
	default:
		return TicketStatusWaiting
	}
}

func GetLastUserTicketsText(telegramID int64, limit int) (string, error) {
	if err := ensureDB(); err != nil {
		return "", err
	}

	if telegramID == 0 {
		return "", fmt.Errorf("invalid telegram id: %d", telegramID)
	}

	if limit <= 0 {
		limit = 5
	}

	rows, err := db.DB.Query(`
		SELECT t.id, t.title, t.description, t.created_at, t.status
		FROM tickets t
		JOIN users u ON u.id = t.user_id
		WHERE u.telegram_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2
	`, telegramID, limit)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	var b strings.Builder
	counter := 1

	for rows.Next() {
		var (
			id          int
			title       string
			description string
			createdAt   time.Time
			status      string
		)

		if err := rows.Scan(
			&id,
			&title,
			&description,
			&createdAt,
			&status,
		); err != nil {
			return "", err
		}

		b.WriteString(fmt.Sprintf(
			"%d. Тема: %s\nОписание: %s\nВремя создания: %s\nСтатус: %s\n\n",
			counter,
			title,
			description,
			createdAt.Format("02.01.2006 15:04:05"),
			status,
		))

		counter++
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	if counter == 1 {
		return "У вас пока нет заявок.", nil
	}

	return b.String(), nil
}

func CreateTicketByTelegramID(telegramID int64, title, description string) error {
	_, err := CreateTicketByTelegramIDWithStatus(telegramID, title, description, TicketStatusWaiting)
	return err
}

func CreateTicketByTelegramIDWithStatus(telegramID int64, title, description, status string) (int, error) {
	if err := ensureDB(); err != nil {
		return 0, err
	}

	if telegramID == 0 {
		return 0, fmt.Errorf("invalid telegram id: %d", telegramID)
	}

	if strings.TrimSpace(title) == "" {
		return 0, errors.New("title is empty")
	}

	if strings.TrimSpace(description) == "" {
		return 0, errors.New("description is empty")
	}

	status = normalizeTicketStatus(status)

	var ticketID int
	err := db.DB.QueryRow(`
		INSERT INTO tickets (user_id, title, description, status, created_at, updated_at)
		SELECT u.id, $2, $3, $4, NOW(), NOW()
		FROM users u
		WHERE u.telegram_id = $1
		RETURNING id
	`, telegramID, title, description, status).Scan(&ticketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user with telegram_id=%d not found", telegramID)
		}

		return 0, err
	}

	return ticketID, nil
}

func UpdateTicketStatus(ticketID int, status string) error {
	if err := ensureDB(); err != nil {
		return err
	}

	if ticketID <= 0 {
		return fmt.Errorf("invalid ticket id: %d", ticketID)
	}

	status = normalizeTicketStatus(status)

	result, err := db.DB.Exec(`
		UPDATE tickets
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`, ticketID, status)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("ticket with id=%d not found", ticketID)
	}

	return nil
}

func BindExternalTicketID(localTicketID int, externalTicketID string) error {
	if err := ensureDB(); err != nil {
		return err
	}

	if localTicketID <= 0 {
		return fmt.Errorf("invalid local ticket id: %d", localTicketID)
	}

	externalTicketID = strings.TrimSpace(externalTicketID)
	if externalTicketID == "" {
		return errors.New("external ticket id is empty")
	}

	result, err := db.DB.Exec(`
		UPDATE tickets
		SET external_ticket_id = $2, updated_at = NOW()
		WHERE id = $1
	`, localTicketID, externalTicketID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("ticket with id=%d not found", localTicketID)
	}

	return nil
}

func ListTelegramIDsByExternalTicketID(externalTicketID string) ([]int64, error) {
	if err := ensureDB(); err != nil {
		return nil, err
	}

	externalTicketID = strings.TrimSpace(externalTicketID)
	if externalTicketID == "" {
		return nil, errors.New("external ticket id is empty")
	}

	rows, err := db.DB.Query(`
		SELECT DISTINCT u.telegram_id
		FROM tickets t
		JOIN users u ON u.id = t.user_id
		WHERE t.external_ticket_id = $1
	`, externalTicketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]int64, 0)
	for rows.Next() {
		var telegramID int64
		if err := rows.Scan(&telegramID); err != nil {
			return nil, err
		}

		if telegramID != 0 {
			result = append(result, telegramID)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
