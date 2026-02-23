package https

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/joho/godotenv"
)

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
}