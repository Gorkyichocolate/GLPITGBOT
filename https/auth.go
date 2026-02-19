package https

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"ok":      false,
			"message": "token is required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
