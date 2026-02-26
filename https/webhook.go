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
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const webhookDebugBodyPreviewLimit = 1200

func webhookDebugEnabled() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("WEBHOOK_DEBUG")))
	return v == "1" || v == "true" || v == "yes"
}

func debugWebhookRequest(c *gin.Context, body []byte) {
	if !webhookDebugEnabled() {
		return
	}

	bodyPreview := string(body)
	if len(bodyPreview) > webhookDebugBodyPreviewLimit {
		bodyPreview = bodyPreview[:webhookDebugBodyPreviewLimit] + "...(truncated)"
	}

	log.Printf(
		"webhook debug: remote=%s method=%s path=%s content_type=%s user_agent=%s headers=%v body=%s",
		c.ClientIP(),
		c.Request.Method,
		c.Request.URL.Path,
		c.GetHeader("Content-Type"),
		c.GetHeader("User-Agent"),
		c.Request.Header,
		bodyPreview,
	)
}

func readIncomingWebhookKey(c *gin.Context) string {
	for _, value := range []string{
		c.GetHeader("X-Webhook-Key"),
		c.GetHeader("X-Webhook-Secret"),
		c.GetHeader("X-GLPI-Webhook-Key"),
		c.GetHeader("X-GLPI-Webhook-Secret"),
		c.GetHeader("Authorization"),
		c.Query("key"),
		c.Query("secret"),
	} {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		if strings.HasPrefix(strings.ToLower(value), "bearer ") {
			return strings.TrimSpace(value[7:])
		}

		return value
	}

	return ""
}

func uniqueNonEmpty(values ...string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		if _, ok := seen[value]; ok {
			continue
		}

		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}

func verifyHMACSHA256(secret string, message []byte, signatureHex string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signatureHex))
}

func verifyGLPISignature(secret string, body []byte, timestamp string, signature string) bool {
	if strings.TrimSpace(signature) == "" {
		return false
	}

	signature = strings.TrimSpace(strings.TrimPrefix(strings.ToLower(strings.TrimSpace(signature)), "sha256="))

	if verifyHMACSHA256(secret, body, signature) {
		return true
	}

	timestamp = strings.TrimSpace(timestamp)
	if timestamp == "" {
		return false
	}

	if verifyHMACSHA256(secret, append(append([]byte{}, body...), []byte(timestamp)...), signature) {
		return true
	}

	return verifyHMACSHA256(secret, append(append([]byte{}, []byte(timestamp)...), body...), signature)
}

func verifyWebhookAuth(c *gin.Context, body []byte, allowedKeys ...string) bool {
	incomingKey := readIncomingWebhookKey(c)
	validKeys := uniqueNonEmpty(allowedKeys...)
	for _, key := range validKeys {
		if incomingKey != "" && incomingKey == key {
			return true
		}
	}

	signature := strings.TrimSpace(c.GetHeader("X-GLPI-Signature"))
	timestamp := strings.TrimSpace(c.GetHeader("X-GLPI-Timestamp"))
	if signature == "" {
		return false
	}

	secrets := uniqueNonEmpty(
		strings.TrimSpace(os.Getenv("WEBHOOK_SECRET")),
		strings.TrimSpace(os.Getenv("WEBHOOK_KEY")),
		strings.TrimSpace(os.Getenv("WEBHOOK_TICKET")),
		strings.TrimSpace(os.Getenv("WEBHOOK_COMMENT")),
	)

	for _, secret := range secrets {
		if verifyGLPISignature(secret, body, timestamp, signature) {
			return true
		}
	}

	return false
}

func WebhookGLPI(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read body"})
		return
	}

	debugWebhookRequest(c, body)

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
		ItemsID json.RawMessage `json:"items_id"`
	}
	_ = json.Unmarshal(base.Item, &itemProbe)

	event := strings.ToLower(strings.TrimSpace(base.Event))
	itemType := strings.ToLower(strings.TrimSpace(probe.ItemType))
	hasParentItem := len(bytes.TrimSpace(parentProbe.ParentItem)) > 0 && string(bytes.TrimSpace(parentProbe.ParentItem)) != "null"
	hasItemsID := notifications.UnmarshalStringOrNumber(itemProbe.ItemsID) != ""
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
	ticketID := notifications.UnmarshalStringOrNumber(payload.ParentItem.ID)
	if ticketID == "" {
		ticketID = notifications.UnmarshalStringOrNumber(payload.Item.ItemsID)
	}

	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is missing in webhook payload"})
		return
	}

	notificationText := notifications.BuildTicketCommentNotification(
		payload.ParentItem.Name,
		notifications.FormatCommentText(payload.Item.Content),
		notifications.UnmarshalStatusValue(payload.ParentItem.Status),
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
	ticketID := notifications.UnmarshalStringOrNumber(payload.Item.ID)
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID is missing in webhook payload"})
		return
	}

	statusValue := payload.Item.StatusName
	if strings.TrimSpace(statusValue) == "" {
		statusValue = notifications.UnmarshalStatusValue(payload.Item.Status)
	}

	notificationText := notifications.BuildTicketStatusChangedNotification(
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
