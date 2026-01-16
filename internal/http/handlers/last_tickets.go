package handlers

import (
	"net/http"

	"GLPITGBOT/internal/service"

	"github.com/gin-gonic/gin"
)

func GetLastTickets(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	tickets, err := service.GetLastTickets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}
