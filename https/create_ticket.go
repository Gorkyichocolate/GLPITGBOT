package https

import (
	"GLPITGBOT/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type createTicketResponse struct {
	ID json.RawMessage `json:"id"`
}

func ExtractCreatedTicketID(body []byte) (string, error) {
	var resp createTicketResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("invalid create ticket response: %w", err)
	}

	if len(resp.ID) == 0 {
		return "", fmt.Errorf("id field is missing in create ticket response")
	}

	var idAsString string
	if err := json.Unmarshal(resp.ID, &idAsString); err == nil {
		idAsString = strings.TrimSpace(idAsString)
		if idAsString != "" {
			return idAsString, nil
		}
	}

	var idAsNumber int64
	if err := json.Unmarshal(resp.ID, &idAsNumber); err == nil {
		return fmt.Sprintf("%d", idAsNumber), nil
	}

	return "", fmt.Errorf("unsupported id format in create ticket response")
}

type createTicketRequest struct {
	Input models.CreateTicketInput `json:"input"`
}

func CreateTicketWithSession(sessionToken string, input models.CreateTicketInput) ([]byte, int, error) {
	baseURL := strings.TrimSpace(os.Getenv("GLPI_URL"))
	if baseURL == "" {
		baseURL = strings.TrimSpace(os.Getenv("GLPI_API"))
	}
	appToken := strings.TrimSpace(os.Getenv("APP_TOKEN"))
	sessionToken = strings.TrimSpace(sessionToken)

	if baseURL == "" {
		return nil, 0, fmt.Errorf("GLPI_URL is empty")
	}

	if appToken == "" {
		return nil, 0, fmt.Errorf("APP_TOKEN is empty")
	}

	if sessionToken == "" {
		return nil, 0, fmt.Errorf("session token is empty")
	}

	payload, err := json.Marshal(createTicketRequest{Input: input})
	if err != nil {
		return nil, 0, err
	}

	endpoint := strings.TrimRight(baseURL, "/") + "/apirest.php/Ticket"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-Token", appToken)
	req.Header.Set("Session-Token", sessionToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return body, resp.StatusCode, fmt.Errorf("create ticket failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return body, resp.StatusCode, nil
}

func CreateTicket(c *gin.Context) {
	sessionToken := c.Query("session_token")
	if strings.TrimSpace(sessionToken) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session_token"})
		return
	}

	var reqBody createTicketRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body, statusCode, err := CreateTicketWithSession(sessionToken, reqBody.Input)
	if err != nil {
		if statusCode > 0 {
			c.Data(statusCode, "application/json", body)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(statusCode, "application/json", body)
}
