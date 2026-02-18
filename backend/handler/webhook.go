package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/RyokouKanai/gomethod/service"
	"github.com/gin-gonic/gin"
)

// WebhookHandler handles LINE webhook callbacks.
func WebhookHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}

	signature := c.GetHeader("X-Line-Signature")
	if !validateSignature(body, signature) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}

	var webhook LineWebhookBody
	if err := json.Unmarshal(body, &webhook); err != nil {
		log.Printf("Error parsing events: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse events"})
		return
	}

	es := service.NewEventService()
	for _, event := range webhook.Events {
		switch event.Type {
		case "follow":
			es.HandleFollow(event.Source.UserID)
		case "message":
			text := ""
			if event.Message.Text != "" {
				text = event.Message.Text
			} else if event.Message.ID != "" {
				text = event.Message.ID
			}
			es.HandleMessage(event.Source.UserID, text, event.ReplyToken)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func validateSignature(body []byte, signature string) bool {
	secret := os.Getenv("LINE_CHANNEL_SECRET")
	if secret == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// LineEvent represents a simplified LINE webhook event.
type LineEvent struct {
	Type       string      `json:"type"`
	ReplyToken string      `json:"replyToken"`
	Source     LineSource  `json:"source"`
	Message    LineMessage `json:"message"`
}

type LineSource struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

type LineMessage struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

type LineWebhookBody struct {
	Events []LineEvent `json:"events"`
}

// MCJHandler handles the MCJ registration endpoint.
func MCJHandler(c *gin.Context) {
	// Placeholder for MCJ::RegistProductAndRelateAuthRule
	log.Println("MCJ regist_product_and_relate_auth_rule called")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
