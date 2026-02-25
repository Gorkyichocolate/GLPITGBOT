package repository

import (
	"GLPITGBOT/db"
	"GLPITGBOT/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func ensureDB() error {
	if db.DB == nil {
		return errors.New("database is not connected")
	}

	return nil
}

func GetUserByTelegramID(telegramID int64) (*models.User, error) {
	if err := ensureDB(); err != nil {
		return nil, err
	}

	if telegramID == 0 {
		return nil, fmt.Errorf("invalid telegram id: %d", telegramID)
	}

	u := &models.User{}

	err := db.DB.QueryRow(`
		SELECT id, telegram_id, username, api_token, session_token, lang
		FROM users
		WHERE telegram_id = $1
	`, telegramID).Scan(
		&u.ID,
		&u.TelegramID,
		&u.Username,
		&u.ApiToken,
		&u.SessionToken,
		&u.Lang,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func SaveUser(user *models.User) error {
	if err := ensureDB(); err != nil {
		return err
	}

	if user == nil {
		return errors.New("user is nil")
	}

	if user.TelegramID == 0 {
		return fmt.Errorf("invalid telegram id: %d", user.TelegramID)
	}

	err := db.DB.QueryRow(`
		INSERT INTO users (telegram_id, username, api_token, session_token, lang)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (telegram_id) DO UPDATE SET
			username = EXCLUDED.username,
			api_token = EXCLUDED.api_token,
			session_token = EXCLUDED.session_token,
			lang = EXCLUDED.lang
		RETURNING id, telegram_id, username, api_token, session_token, lang
	`,
		user.TelegramID,
		user.Username,
		user.ApiToken,
		user.SessionToken,
		user.Lang,
	).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.ApiToken,
		&user.SessionToken,
		&user.Lang,
	)

	return err
}

func EnsureUser(telegramID int64, username string) (*models.User, error) {
	if telegramID == 0 {
		return nil, fmt.Errorf("invalid telegram id: %d", telegramID)
	}

	user, err := GetUserByTelegramID(telegramID)

	if err == nil {
		return user, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	user = &models.User{
		TelegramID: telegramID,
		Username:   username,
		Lang:       "",
	}

	if err := SaveUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func IsUserAuthorized(user *models.User) bool {
	if user == nil {
		return false
	}

	return strings.TrimSpace(user.ApiToken) != "" && strings.TrimSpace(user.SessionToken) != ""
}

func ListAuthorizedTelegramIDs() ([]int64, error) {
	if err := ensureDB(); err != nil {
		return nil, err
	}

	rows, err := db.DB.Query(`
		SELECT telegram_id
		FROM users
		WHERE COALESCE(TRIM(api_token), '') <> ''
		  AND COALESCE(TRIM(session_token), '') <> ''
	`)
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
