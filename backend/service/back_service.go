package service

import (
	"github.com/RyokouKanai/gomethod/model"
)

// BackService handles "戻る" (back) messages.
type BackService struct {
	BaseService
}

func NewBackService(user *model.User, msg, token string, ss *SendService) *BackService {
	return &BackService{BaseService: newBaseService(user, msg, token, ss)}
}

func (s *BackService) Executed() bool {
	return s.ReceivedMessage == "戻る" && s.execute()
}

func (s *BackService) Execute() {
	s.execute()
}

func (s *BackService) execute() bool {
	histories, err := s.User.GetRecentTalkHistories()
	if err != nil {
		return false
	}

	// Force back if not enough history or missing reply patterns
	if len(histories) < 2 || s.forceBack(histories) {
		return s.executeForceBack()
	}

	// Get reply pattern from second-to-last history
	th := histories[1]
	rp := th.GetReplyPattern()
	if rp == nil {
		return s.executeForceBack()
	}

	nextMsg := rp.GetNextMessage()
	if nextMsg == nil {
		return s.executeForceBack()
	}

	s.sendService.Reply(nextMsg.ToFormattedText(), s.ReplyToken)
	newTH, err := s.createTalkHistory(nextMsg)
	if err == nil && newTH != nil {
		rpID := int(rp.ID)
		newTH.ReplyPatternID = &rpID
		model.UpdateTalkHistoryReplyPattern(newTH)
	}
	return true
}

func (s *BackService) forceBack(histories []model.TalkHistory) bool {
	for _, h := range histories[:2] {
		if h.ReplyPatternID == nil {
			return true
		}
	}
	return false
}

func (s *BackService) executeForceBack() bool {
	topMsg := model.GetMessageByScope("default")
	if topMsg == nil {
		return false
	}
	s.sendService.Reply(topMsg.ToFormattedText(), s.ReplyToken)
	s.createTalkHistory(topMsg)
	return true
}
