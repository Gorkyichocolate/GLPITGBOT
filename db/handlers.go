package db

import (
	"errors"
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

func CreateTicketByTelegramID(telegramID int64, title, description string) error {

	if DB == nil {
		return errors.New("DB is nil")
	}

	_, err := DB.Exec(`
		INSERT INTO tickets (user_id, title, description, status, created_at)
		SELECT u.id, $2, $3, 'open', NOW()
		FROM users u
		WHERE u.telegram_id = $1
	`, telegramID, title, description)

	return err
}

func SaveUser(user *User) error {
	if DB == nil {
		return errors.New("DB is nil")
	}

	_, err := DB.Exec(`
		INSERT INTO users (telegram_id, username, api_token, lang)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (telegram_id) DO UPDATE SET
			username = EXCLUDED.username,
			api_token = EXCLUDED.api_token,
			lang = EXCLUDED.lang
	`,
		user.TelegramId,
		user.Username,
		user.ApiToken,
		user.Lang,
	)

	return err
}

func GetUserByTelegramID(telegramID int64) (*User, error) {
	if DB == nil {
		return nil, errors.New("DB is nil")
	}

	u := &User{}

	err := DB.QueryRow(`
		SELECT id, telegram_id, username, api_token, lang
		FROM users
		WHERE telegram_id = $1
	`, telegramID).Scan(
		&u.Id,
		&u.TelegramId,
		&u.Username,
		&u.ApiToken,
		&u.Lang,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func EnsureUser(telegramID int64, username string) (*User, error) {
	user, err := GetUserByTelegramID(telegramID)
	if err == nil {
		return user, nil
	}

	user = &User{
		TelegramId: telegramID,
		Username:   username,
		Lang:       "",
	}

	err = SaveUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
