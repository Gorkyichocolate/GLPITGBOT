package handlers

import (
	"net/http"

	"GLPITGBOT/internal/service"

	"github.com/gin-gonic/gin"
)

type GLPIWebhookRequest struct {
	TicketID string `json:"ticket_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

func GLPIWebhook(c *gin.Context) {
	var req GLPIWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1️⃣ Обновляем тикет в БД
	if err := service.UpdateTicketStatus(req.TicketID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2️⃣ Уведомляем Telegram (асинхронно)
	go service.NotifyUserAboutStatusChange(req.TicketID, req.Status)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
