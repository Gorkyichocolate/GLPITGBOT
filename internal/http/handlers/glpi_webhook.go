package handlers

import (
	"GLPITGBOT/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TicketStatus(c *gin.Context) {
	if err := repository.Pool.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "db down",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
