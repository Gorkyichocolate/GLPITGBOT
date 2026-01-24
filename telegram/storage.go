package telegram

import "sync"

var (
	users = make(map[int64]*Ticket)
	mu    sync.Mutex
)

func getUser(chatID int64) *Ticket {
	mu.Lock()
	defer mu.Unlock()

	if user, ok := users[chatID]; ok {
		return user
	}

	user := &Ticket{
		State: StateIdle,
	}

	users[chatID] = user
	return user
}
