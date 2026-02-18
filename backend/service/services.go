package service

import (
	"github.com/RyokouKanai/gomethod/model"
)

// BaseService provides common fields and methods for all services.
type BaseService struct {
	User            *model.User
	ReceivedMessage string
	ReplyToken      string
	sendService     *SendService
}

func newBaseService(user *model.User, receivedMessage, replyToken string, ss *SendService) BaseService {
	return BaseService{
		User:            user,
		ReceivedMessage: receivedMessage,
		ReplyToken:      replyToken,
		sendService:     ss,
	}
}

func (bs *BaseService) createTalkHistory(message *model.Message) (*model.TalkHistory, error) {
	return bs.User.CreateTalkHistory(message)
}

// AvailableService checks if the user is active.
type AvailableService struct {
	BaseService
}

func NewAvailableService(user *model.User, msg, token string, ss *SendService) *AvailableService {
	return &AvailableService{BaseService: newBaseService(user, msg, token, ss)}
}

func (s *AvailableService) Executed() bool {
	return !s.User.IsActive && s.execute()
}

func (s *AvailableService) Execute() {
	s.execute()
}

func (s *AvailableService) execute() bool {
	unavailableMsg := model.GetMessageByScope("unavailable")
	if unavailableMsg == nil {
		return false
	}
	s.sendService.Reply(unavailableMsg.ToFormattedText(), s.ReplyToken)
	s.createTalkHistory(unavailableMsg)
	return true
}

// ThanksCountService handles thanks counting.
type ThanksCountService struct {
	BaseService
}

func NewThanksCountService(user *model.User, msg, token string, ss *SendService) *ThanksCountService {
	return &ThanksCountService{BaseService: newBaseService(user, msg, token, ss)}
}

func (s *ThanksCountService) Executed() bool {
	return s.ReceivedMessage == "ありがとう、感謝します" && s.execute()
}

func (s *ThanksCountService) Execute() {
	s.execute()
}

func (s *ThanksCountService) execute() bool {
	s.User.UpsertActionRecord("thanks_count", 10)
	thanksCount := s.User.GetThanksCount()

	// 100の倍数(x10=20の倍数)ごとにお知らせ
	if thanksCount%20 == 0 {
		message := "あり感ツイート" + itoa(thanksCount) + "回達成おめでとう！"

		// 1000の倍数(x10=50の倍数)ごとに応援メッセージ
		if thanksCount%50 == 0 {
			tl := model.FindThanksLevelByCount(thanksCount)
			if tl != nil && tl.Cheering != nil {
				message += "\n\n " + *tl.Cheering
			}
		}

		s.sendService.Reply(message, s.ReplyToken)
	}
	return true
}

// TopBackService sends the top (default) message when user sends "TOP".
type TopBackService struct {
	BaseService
}

func NewTopBackService(user *model.User, msg, token string, ss *SendService) *TopBackService {
	return &TopBackService{BaseService: newBaseService(user, msg, token, ss)}
}

func (s *TopBackService) Executed() bool {
	return s.ReceivedMessage == "TOP" && s.execute()
}

func (s *TopBackService) Execute() {
	s.execute()
}

func (s *TopBackService) execute() bool {
	topMsg := model.GetMessageByScope("default")
	if topMsg == nil {
		return false
	}
	s.sendService.Reply(topMsg.ToFormattedText(), s.ReplyToken)
	s.createTalkHistory(topMsg)
	return true
}

// TopMessageSendService always sends the top message (fallback).
type TopMessageSendService struct {
	BaseService
}

func NewTopMessageSendService(user *model.User, msg, token string, ss *SendService) *TopMessageSendService {
	return &TopMessageSendService{BaseService: newBaseService(user, msg, token, ss)}
}

func (s *TopMessageSendService) Executed() bool {
	return s.execute()
}

func (s *TopMessageSendService) Execute() {
	s.execute()
}

func (s *TopMessageSendService) execute() bool {
	topMsg := model.GetMessageByScope("default")
	if topMsg == nil {
		return false
	}
	s.sendService.Reply(topMsg.ToFormattedText(), s.ReplyToken)
	s.createTalkHistory(topMsg)
	return true
}

// AdminLoginService handles admin login via LINE.
type AdminLoginService struct {
	BaseService
}

func NewAdminLoginService(user *model.User, msg, token string, ss *SendService) *AdminLoginService {
	return &AdminLoginService{BaseService: newBaseService(user, msg, token, ss)}
}

func (s *AdminLoginService) Executed() bool {
	return s.ReceivedMessage == "ログイン" && s.User.IsAdmin() && s.execute()
}

func (s *AdminLoginService) Execute() {
	s.execute()
}

func (s *AdminLoginService) execute() bool {
	adminMsg := model.GetMessageByScope("admin_default")
	if adminMsg == nil {
		return false
	}
	s.sendService.Reply(adminMsg.ToFormattedText(), s.ReplyToken)
	s.createTalkHistory(adminMsg)
	return true
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	if neg {
		result = "-" + result
	}
	return result
}
