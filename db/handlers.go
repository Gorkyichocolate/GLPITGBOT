package db

import (
	"fmt"
	"strings"
	"time"
)

func GetLastUserTicketsText(telegramID int64, limit int) (string, error) {
	if DB == nil {
		return "DB is nil", nil
	}

	rows, err := DB.Query(`
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

	if counter == 1 {
		return "У вас пока нет заявок.", nil
	}

	return b.String(), nil
}
