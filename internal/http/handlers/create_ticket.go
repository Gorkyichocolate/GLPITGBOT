package handlers

import (
	"net/http"

	"GLPITGBOT/internal/repository"
	"GLPITGBOT/internal/service"

	"github.com/gin-gonic/gin"
)

type CreateTicketRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	TicketID    string `json:"ticket_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func CreateTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.CreateTicket(repository.Ticket{
		UserID:      req.UserID,
		TicketID:    req.TicketID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "new",
		Active:      true,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ticket created"})
}
