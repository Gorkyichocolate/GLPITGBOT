package telegram

import "sync"

var (
	users = make(map[int64]*User)
	mu    sync.Mutex
)

func getUser(chatID int64) *User {
	mu.Lock()
	defer mu.Unlock()

	if user, ok := users[chatID]; ok {
		return user
	}

	user := &User{
		ID:    chatID,
		State: StateIdle,
		Lang:  "",
	}

	users[chatID] = user
	return user
}

func saveUser(user *User) {
	mu.Lock()
	defer mu.Unlock()
	users[user.ID] = user
}
