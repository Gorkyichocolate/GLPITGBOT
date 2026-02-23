package https

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Notification struct {
	TicketID int    `json:"ticket_id"`
	Message  string `json:"message"`
}

func NotifyTicketUpdate(c *gin.Context) {
	if c.GetHeader("X-API-KEY") != os.Getenv("X-API-KEY") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}
	var data Notification

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
