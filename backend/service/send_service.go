package service

import (
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

// SendService handles sending messages via LINE Bot API.
type SendService struct {
	bot *messaging_api.MessagingApiAPI
}

// NewSendService creates a new SendService.
func NewSendService() *SendService {
	bot, err := messaging_api.NewMessagingApiAPI(os.Getenv("LINE_CHANNEL_TOKEN"))
	if err != nil {
		log.Printf("Error creating LINE bot: %v", err)
		return &SendService{}
	}
	return &SendService{bot: bot}
}

// Reply sends a reply message to the given reply token.
func (s *SendService) Reply(messages interface{}, replyToken string) {
	if s.bot == nil {
		return
	}

	var lineMessages []messaging_api.MessageInterface

	switch v := messages.(type) {
	case string:
		chunks := splitMessage(v, 4500)
		for _, chunk := range chunks {
			lineMessages = append(lineMessages, &messaging_api.TextMessage{Text: chunk})
		}
	case []string:
		for _, msg := range v {
			lineMessages = append(lineMessages, &messaging_api.TextMessage{Text: msg})
		}
	}

	if len(lineMessages) == 0 {
		return
	}

	_, err := s.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages:   lineMessages,
	})
	if err != nil {
		log.Printf("Error replying message: %v", err)
	}
}

// ReplyImage sends an image reply.
func (s *SendService) ReplyImage(imageURL, replyToken string) {
	if s.bot == nil {
		return
	}

	_, err := s.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.ImageMessage{
				OriginalContentUrl: imageURL,
				PreviewImageUrl:    imageURL,
			},
		},
	})
	if err != nil {
		log.Printf("Error replying image: %v", err)
	}
}

// ReplyImageAndMessages sends mixed image and text messages.
func (s *SendService) ReplyImageAndMessages(contents []map[string]string, replyToken string) {
	if s.bot == nil {
		return
	}

	var lineMessages []messaging_api.MessageInterface
	for _, content := range contents {
		switch content["type"] {
		case "image":
			lineMessages = append(lineMessages, &messaging_api.ImageMessage{
				OriginalContentUrl: content["content"],
				PreviewImageUrl:    content["content"],
			})
		case "text":
			lineMessages = append(lineMessages, &messaging_api.TextMessage{Text: content["content"]})
		}
	}

	_, err := s.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages:   lineMessages,
	})
	if err != nil {
		log.Printf("Error replying image and messages: %v", err)
	}
}

// Broadcast sends a message to all users.
func (s *SendService) Broadcast(message string) {
	if s.bot == nil {
		return
	}

	_, err := s.bot.Broadcast(&messaging_api.BroadcastRequest{
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{Text: message},
		},
	}, "")
	if err != nil {
		log.Printf("Error broadcasting message: %v", err)
	}
}

// BroadcastToShik sends a message to shik users via push.
func (s *SendService) BroadcastToShik(message string, lineUserIDs []string) {
	for _, id := range lineUserIDs {
		s.Unicast(id, message)
	}
}

// Unicast sends a message to a specific user.
func (s *SendService) Unicast(lineUserID, message string) {
	if s.bot == nil {
		return
	}

	_, err := s.bot.PushMessage(&messaging_api.PushMessageRequest{
		To: lineUserID,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{Text: message},
		},
	}, "")
	if err != nil {
		log.Printf("Error sending unicast message: %v", err)
	}
}

// splitMessage splits a message into chunks of the given max size.
func splitMessage(msg string, maxLen int) []string {
	if len([]rune(msg)) <= maxLen {
		return []string{msg}
	}

	var chunks []string
	runes := []rune(msg)
	for i := 0; i < len(runes); i += maxLen {
		end := i + maxLen
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}
