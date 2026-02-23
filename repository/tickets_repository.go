package repository

import (
	"GLPITGBOT/db"
	"errors"
	"fmt"
	"strings"
	"time"
)

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
	if err := ensureDB(); err != nil {
		return err
	}

	if telegramID == 0 {
		return fmt.Errorf("invalid telegram id: %d", telegramID)
	}

	if strings.TrimSpace(title) == "" {
		return errors.New("title is empty")
	}

	if strings.TrimSpace(description) == "" {
		return errors.New("description is empty")
	}

	result, err := db.DB.Exec(`
		INSERT INTO tickets (user_id, title, description, status, created_at)
		SELECT u.id, $2, $3, 'open', NOW()
		FROM users u
		WHERE u.telegram_id = $1
	`, telegramID, title, description)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user with telegram_id=%d not found", telegramID)
	}

	return nil
}
