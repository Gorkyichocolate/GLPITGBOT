package service

import "GLPITGBOT/internal/repository"

// Получить последние 5 тикетов пользователя
func GetLastTickets(userID string) ([]repository.Ticket, error) {
	return repository.GetLastFiveTickets(userID)
}

// Добавить тикет
func CreateTicket(ticket repository.Ticket) error {
	return repository.AddTicket(ticket)
}
