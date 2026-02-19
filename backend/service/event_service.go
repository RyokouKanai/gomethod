package service

import (
	"log"

	"github.com/RyokouKanai/gomethod/action"
	"github.com/RyokouKanai/gomethod/model"
)

// EventService handles LINE webhook events.
type EventService struct {
	sendService    *SendService
	actionRegistry *action.Registry
}

// NewEventService creates a new EventService.
func NewEventService() *EventService {
	ss := NewSendService()
	return &EventService{
		sendService:    ss,
		actionRegistry: action.NewRegistry(ss),
	}
}

// HandleFollow handles a follow event (new user).
func (es *EventService) HandleFollow(lineUserID string) {
	user, err := model.FindOrCreateByLineUserID(lineUserID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return
	}
	if err := user.SaveProfile(); err != nil {
		log.Printf("Error saving profile: %v", err)
	}
}

// HandleMessage handles a message event.
func (es *EventService) HandleMessage(lineUserID, receivedMessage, replyToken string) {
	user, err := model.FindOrCreateByLineUserID(lineUserID)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return
	}

	// ReplyPatternService にアクションレジストリを接続
	rps := NewReplyPatternService(user, receivedMessage, replyToken, es.sendService)
	rps.SetActionExecutor(es.actionRegistry)

	// Chain of responsibility - same order as Rails
	services := []ServiceHandler{
		NewAvailableService(user, receivedMessage, replyToken, es.sendService),
		NewThanksCountService(user, receivedMessage, replyToken, es.sendService),
		NewTopBackService(user, receivedMessage, replyToken, es.sendService),
		NewBackService(user, receivedMessage, replyToken, es.sendService),
		NewAdminLoginService(user, receivedMessage, replyToken, es.sendService),
		rps,
	}

	for _, svc := range services {
		if svc.Executed() {
			return
		}
	}

	// Default: send top message
	NewTopMessageSendService(user, receivedMessage, replyToken, es.sendService).Execute()
}

// ServiceHandler interface for chain of responsibility pattern.
type ServiceHandler interface {
	Executed() bool
	Execute()
}
