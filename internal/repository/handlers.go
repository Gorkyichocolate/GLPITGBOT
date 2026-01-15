package repository

import (
	"context"
)

type User struct {
	ID       int64
	Username string
	Language string
}

func SaveUser(user User) error {
	_, err := Pool.Exec(
		context.Background(),
		`INSERT INTO users (id, username, language)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (id) DO NOTHING`,
		user.ID,
		user.Username,
		user.Language,
	)
	return err
}
