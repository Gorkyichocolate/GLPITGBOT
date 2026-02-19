package repository

import (
	"GLPITGBOT/db"
	"GLPITGBOT/models"
	"database/sql"
	"errors"
)

func GetUserByTelegramID(telegramID int64) (*models.User, error) {
	u := &models.User{}

	err := db.DB.QueryRow(`
		SELECT id, telegram_id, username, api_token, lang
		FROM users
		WHERE telegram_id = $1
	`, telegramID).Scan(
		&u.ID,
		&u.TelegramID,
		&u.Username,
		&u.ApiToken,
		&u.Lang,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func SaveUser(user *models.User) error {
	_, err := db.DB.Exec(`
		INSERT INTO users (telegram_id, username, api_token, lang)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (telegram_id) DO UPDATE SET
			username = EXCLUDED.username,
			api_token = EXCLUDED.api_token,
			lang = EXCLUDED.lang
	`,
		user.TelegramID,
		user.Username,
		user.ApiToken,
		user.Lang,
	)

	return err
}

func EnsureUser(telegramID int64, username string) (*models.User, error) {
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
		State:      models.StateIdle,
	}

	if err := SaveUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
