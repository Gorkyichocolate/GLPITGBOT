package https

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type authResponse struct {
	SessionToken string `json:"session_token"`
}

func AuthByUserToken(userToken string) (string, error) {
	baseURL := strings.TrimSpace(os.Getenv("GLPI_URL"))
	appToken := strings.TrimSpace(os.Getenv("APP_TOKEN"))
	userToken = strings.TrimSpace(userToken)

	if baseURL == "" {
		return "", fmt.Errorf("GLPI_URL is empty")
	}

	if appToken == "" {
		return "", fmt.Errorf("APP_TOKEN is empty")
	}

	if userToken == "" {
		return "", fmt.Errorf("user token is empty")
	}

	endpoint := strings.TrimRight(baseURL, "/") + "/apirest.php/initSession"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "user_token "+userToken)
	req.Header.Set("App-Token", appToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var result authResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("invalid auth response: %w", err)
	}

	if strings.TrimSpace(result.SessionToken) == "" {
		return "", fmt.Errorf("session token is empty in auth response")
	}

	return result.SessionToken, nil
}

func CheckAuth(c *gin.Context) {
	token := c.Query("token")
	sessionToken, err := AuthByUserToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":            true,
		"session_token": sessionToken,
	})
}
