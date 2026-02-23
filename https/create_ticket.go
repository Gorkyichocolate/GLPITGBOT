package https

import (
	"time"
	"net/http"
	"io"
	"github.com/gin-gonic/gin"
)

func CreateTicket(c *gin.Context) {
	token := c.Query("token")
	app_token := c.Query("app_token")

	if token == "" || app_token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token or app_token"})
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/apirest.php/Ticket", nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization","user-token "+token)
	req.Header.Set("App-Token", app_token)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.Data(resp.StatusCode, "application/json", body)
	if resp.StatusCode != http.StatusOK {
		return
	}

	
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
