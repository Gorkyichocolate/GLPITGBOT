package https

import (
	"io"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {
	token := c.Query("token")

	client = &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/apirest.php/initSession", nil)
	if err!=nil{
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization","user-token "+token)

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
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
