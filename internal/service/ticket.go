package service

import (
	"GLPITGBOT/internal/repository"
	"GLPITGBOT/internal/telegram"
)

// Получить последние 5 тикетов пользователя
func GetLastTickets(userID string) ([]repository.Ticket, error) {
	return repository.GetLastFiveTickets(userID)
}

// Добавить тикет
func CreateTicket(ticket repository.Ticket) error {
	return repository.AddTicket(ticket)
}

// Смена статуса заявки
func UpdateTicketStatus(ticketID, status string) error {
	return repository.UpdateTicketStatus(ticketID, status)
}

// Уведомление о смене статуса заявки
func NotifyUserAboutStatusChange(ticketID, status string) error {

	// 1 Получаем тикет
	ticket, err := repository.GetTicketByTicketID(ticketID)
	if err != nil {
		return err
	}

	// 2 Получаем пользователя
	user, err := repository.GetUserByUserID(ticket.UserID)
	if err != nil {
		return err
	}

	// 3 Отправляем сообщение
	return telegram.SendMessage(
		user.TelegramID,
		"Статус вашей заявки "+ticket.TicketID+" изменён: "+status,
	)
}

// GetTicketStatus — получить статус тикета по ticket_id
func GetTicketStatus(ticketID string) (string, error) {
	ticket, err := repository.GetTicketByTicketID(ticketID)
	if err != nil {
		return "", err
	}
	return ticket.Status, nil
}
