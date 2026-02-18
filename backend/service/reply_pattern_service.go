package service

import (
	"log"
	"strconv"

	"github.com/RyokouKanai/gomethod/model"
)

// ReplyPatternService handles message reply patterns.
type ReplyPatternService struct {
	BaseService
	actionExecutor ActionExecutor
}

// ActionExecutor is an interface for executing reply pattern actions.
type ActionExecutor interface {
	Execute(method string, user *model.User, receivedMessage string, replyToken string, nextMessage *model.Message) interface{}
}

func NewReplyPatternService(user *model.User, msg, token string, ss *SendService) *ReplyPatternService {
	return &ReplyPatternService{
		BaseService: newBaseService(user, msg, token, ss),
	}
}

// SetActionExecutor sets the action executor (called from main to avoid circular deps).
func (s *ReplyPatternService) SetActionExecutor(ae ActionExecutor) {
	s.actionExecutor = ae
}

func (s *ReplyPatternService) Executed() bool {
	return s.enablePatternReply() && s.execute()
}

func (s *ReplyPatternService) Execute() {
	s.execute()
}

func (s *ReplyPatternService) execute() bool {
	rp := s.replyPattern()
	if rp == nil {
		return false
	}

	nextMsg := rp.GetNextMessage()
	if nextMsg == nil {
		return false
	}

	// Execute the method, get reply content
	content := s.nextMessageContents(rp, nextMsg)

	s.sendService.Reply(content, s.ReplyToken)
	th, err := s.createTalkHistory(nextMsg)
	if err == nil && th != nil {
		rpID := int(rp.ID)
		th.ReplyPatternID = &rpID
		model.UpdateTalkHistoryReplyPattern(th)
	}
	return true
}

func (s *ReplyPatternService) enablePatternReply() bool {
	return s.lastSentMessage() != nil && s.replyPattern() != nil
}

func (s *ReplyPatternService) nextMessageContents(rp *model.ReplyPattern, nextMsg *model.Message) interface{} {
	if s.actionExecutor != nil && rp.ExecutionMethod != "base" {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Error executing action %s: %v", rp.ExecutionMethod, r)
			}
		}()

		result := s.actionExecutor.Execute(rp.ExecutionMethod, s.User, s.ReceivedMessage, s.ReplyToken, nextMsg)
		if result != nil {
			return result
		}
	}

	// Default: base method - just return formatted text
	return nextMsg.ToFormattedText()
}

func (s *ReplyPatternService) replyPattern() *model.ReplyPattern {
	lastMsg := s.lastSentMessage()
	if lastMsg == nil {
		return nil
	}

	if s.receivedOption(lastMsg) {
		pos, _ := strconv.Atoi(s.ReceivedMessage)
		return model.FindReplyPatternByMessageAndPosition(lastMsg.ID, pos)
	}
	return model.FindFirstReplyPatternByMessage(lastMsg.ID)
}

func (s *ReplyPatternService) lastSentMessage() *model.Message {
	th, err := s.User.GetLatestTalkHistory()
	if err != nil || th == nil {
		return nil
	}
	return th.GetMessage()
}

func (s *ReplyPatternService) receivedOption(lastMsg *model.Message) bool {
	options, _ := lastMsg.GetOptions()
	if len(options) == 0 {
		return false
	}
	if s.ReceivedMessage == "0" {
		return false
	}
	n, err := strconv.Atoi(s.ReceivedMessage)
	return err == nil && n != 0
}
