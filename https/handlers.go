package https

import (
	"GLPITGBOT/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, false)
		return
	}

	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		c.JSON(http.StatusUnauthorized, false)
		return
	}

	token := authHeader[len(prefix):]

	var exists bool
	err := db.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM users WHERE api_token = $1
		)
	`, token).Scan(&exists)

	if err != nil || !exists {
		c.JSON(http.StatusUnauthorized, false)
		return
	}

	c.JSON(http.StatusOK, true)
}
