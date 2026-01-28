package https

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {

	c.JSON(http.StatusOK, true)
}
