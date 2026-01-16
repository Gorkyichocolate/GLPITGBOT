// internal/http/handlers/ticket_status.go
package handlers

import (
	"net/http"

	"GLPITGBOT/internal/service"
	"github.com/gin-gonic/gin"
)

func TicketStatus(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket id required"})
		return
	}

	status, err := service.GetTicketStatus(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ticket_id": ticketID,
		"status":    status,
	})
}
