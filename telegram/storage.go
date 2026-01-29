package telegram

import (
	"GLPITGBOT/db"
	"sync"
)

var (
	users = make(map[int64]*User)
	mu    sync.Mutex
)

func getUser(telegramID int64) *User {
	user, err := db.GetUserByTelegramID(telegramID)
	if err == nil {
		return user
	}

	// если нет в БД — создаём
	user = &User{
		TelegramID: telegramID,
		Lang:       "",
		State:      StateIdle,
	}

	db.SaveUser(user)
	return user
}

func saveUser(user *User) {
	mu.Lock()
	defer mu.Unlock()
	users[user.ID] = user
}
