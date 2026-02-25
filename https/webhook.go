package https

import (
	"GLPITGBOT/models"
	"GLPITGBOT/notifications"
	"GLPITGBOT/repository"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}

	return ""
}

func readWebhookSecret(c *gin.Context) string {
	return strings.TrimSpace(firstNonEmpty(
		c.GetHeader("X-Webhook-Key"),
		c.GetHeader("X-Webhook-Secret"),
		c.GetHeader("X-GLPI-Webhook-Key"),
		c.GetHeader("X-GLPI-Webhook-Secret"),
		c.GetHeader("Authorization"),
		c.Query("key"),
		c.Query("secret"),
	))
}

func normalizeBearerToken(value string) string {
	v := strings.TrimSpace(value)
	if strings.HasPrefix(strings.ToLower(v), "bearer ") {
		return strings.TrimSpace(v[7:])
	}

	return v
}

func verifyGLPISignature(secret string, body []byte, ts int64, sigHeader string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	mac.Write([]byte(fmt.Sprintf("%d", ts)))
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(strings.TrimSpace(sigHeader)))
}

func verifyWebhookAuth(c *gin.Context, body []byte, allowedKeys ...string) bool {
	incomingKey := normalizeBearerToken(readWebhookSecret(c))
	if incomingKey != "" {
		for _, key := range allowedKeys {
			if strings.TrimSpace(key) != "" && incomingKey == strings.TrimSpace(key) {
				return true
			}
		}
	}

	signSecret := strings.TrimSpace(os.Getenv("WEBHOOK_SECRET"))
	if signSecret == "" {
		return false
	}

	sig := strings.TrimSpace(c.GetHeader("X-GLPI-Signature"))
	tsHeader := strings.TrimSpace(c.GetHeader("X-GLPI-Timestamp"))
	if sig == "" || tsHeader == "" {
		return false
	}

	ts, err := strconv.ParseInt(tsHeader, 10, 64)
	if err != nil {
		return false
	}

	return verifyGLPISignature(signSecret, body, ts, sig)
}

func WebhookGLPI(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body"})
		return
	}

	if !verifyWebhookAuth(
		c,
		body,
		os.Getenv("WEBHOOK_KEY"),
		os.Getenv("WEBHOOK_TICKET"),
		os.Getenv("WEBHOOK_COMMENT"),
	) {
		log.Printf("webhook unauthorized: remote=%s path=%s", c.ClientIP(), c.Request.URL.Path)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var base models.BaseWebhookPayload
	if err := json.Unmarshal(body, &base); err != nil {
		log.Printf("webhook invalid json: remote=%s err=%v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	var probe models.WebhookItemTypeProbe
	_ = json.Unmarshal(base.Item, &probe)

	var parentProbe struct {
		ParentItem json.RawMessage `json:"parent_item"`
	}
	_ = json.Unmarshal(body, &parentProbe)

	var itemProbe struct {
		ItemsID models.StringOrNumber `json:"items_id"`
	}
	_ = json.Unmarshal(base.Item, &itemProbe)

	event := strings.ToLower(strings.TrimSpace(base.Event))
	itemType := strings.ToLower(strings.TrimSpace(probe.ItemType))
	hasParentItem := len(bytes.TrimSpace(parentProbe.ParentItem)) > 0 && string(bytes.TrimSpace(parentProbe.ParentItem)) != "null"
	hasItemsID := itemProbe.ItemsID.String() != ""
	isFollowupPayload := itemType == "itilfollowup" || event == "add_followup" || (hasParentItem && hasItemsID)

	switch {
	case isFollowupPayload:
		var payload models.FullFollowupWebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid followup payload"})
			return
		}

		handleFollowupWebhook(c, payload)
		return

	case itemType == "ticket" || event == "update" || event == "new":
		var payload models.FullTicketWebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket payload"})
			return
		}

		handleTicketWebhook(c, payload)
		return

	default:
		log.Printf("webhook unknown payload: event=%s itemtype=%s", base.Event, probe.ItemType)
		c.JSON(http.StatusOK, gin.H{"message": "Ignored unknown webhook payload"})
	}
}

func handleFollowupWebhook(c *gin.Context, payload models.FullFollowupWebhookPayload) {
	ticketID := payload.ParentItem.ID.String()
	if ticketID == "" {
		ticketID = payload.Item.ItemsID.String()
	}

	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is missing in webhook payload"})
		return
	}

	notificationText := models.BuildTicketCommentNotification(
		payload.ParentItem.Name,
		payload.Item.Content,
		payload.ParentItem.Status.String(),
	)

	recipients, err := repository.ListTelegramIDsByExternalTicketID(ticketID)
	if err != nil {
		log.Println("followup webhook recipients error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve recipients"})
		return
	}

	if len(recipients) == 0 {
		log.Printf("followup webhook no recipients: event=%s external_ticket_id=%s", payload.Event, ticketID)
		c.JSON(http.StatusOK, gin.H{"message": "Followup received, no recipients"})
		return
	}

	if err := notifications.SendTelegramMessageToMany(recipients, notificationText); err != nil {
		log.Println("followup webhook send error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send telegram notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followup received"})
}

func handleTicketWebhook(c *gin.Context, payload models.FullTicketWebhookPayload) {
	ticketID := payload.Item.ID.String()
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is missing in webhook payload"})
		return
	}

	statusValue := payload.Item.StatusName
	if strings.TrimSpace(statusValue) == "" {
		statusValue = payload.Item.Status.String()
	}

	notificationText := models.BuildTicketStatusChangedNotification(
		payload.Item.Name,
		statusValue,
	)

	recipients, err := repository.ListTelegramIDsByExternalTicketID(ticketID)
	if err != nil {
		log.Println("ticket webhook recipients error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve recipients"})
		return
	}

	if len(recipients) == 0 {
		log.Printf("ticket webhook no recipients: event=%s external_ticket_id=%s", payload.Event, ticketID)
		c.JSON(http.StatusOK, gin.H{"message": "Ticket update received, no recipients"})
		return
	}

	if err := notifications.SendTelegramMessageToMany(recipients, notificationText); err != nil {
		log.Println("ticket webhook send error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send telegram notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket update received"})
}
