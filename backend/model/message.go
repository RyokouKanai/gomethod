package model

import (
	"fmt"
	"strings"

	"github.com/RyokouKanai/gomethod/database"
)

type Message struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	Content *string `gorm:"type:text" json:"content"`
}

func (Message) TableName() string { return "messages" }

// GetContent safely returns the content string.
func (m *Message) GetContent() string {
	if m.Content != nil {
		return *m.Content
	}
	return ""
}

// GetOptions returns options for this message.
func (m *Message) GetOptions() ([]Option, error) {
	var options []Option
	err := database.DB.Where("message_id = ?", m.ID).Order("position ASC").Find(&options).Error
	return options, err
}

// ToFormattedText returns formatted text including options if present.
func (m *Message) ToFormattedText() string {
	options, _ := m.GetOptions()
	content := m.GetContent()

	if len(options) > 0 {
		var optTexts []string
		for _, o := range options {
			optTexts = append(optTexts, fmt.Sprintf("%d: %s", o.Position, o.GetContent()))
		}
		selectNum := GetMessageByScope("select_number")
		selectText := ""
		if selectNum != nil {
			selectText = selectNum.GetContent()
		}
		content = strings.Join([]string{content, strings.Join(optTexts, "\n"), selectText}, "\n\n")
	}
	return content
}

// GetReplyPatterns returns reply patterns where this message is the sent message.
func (m *Message) GetReplyPatterns() ([]ReplyPattern, error) {
	var patterns []ReplyPattern
	err := database.DB.Where("sent_message_id = ?", m.ID).Find(&patterns).Error
	return patterns, err
}

// GetNextReplyPattern returns the first reply pattern for this message.
func (m *Message) GetNextReplyPattern() *ReplyPattern {
	var rp ReplyPattern
	err := database.DB.Where("sent_message_id = ?", m.ID).First(&rp).Error
	if err != nil {
		return nil
	}
	return &rp
}

// Message scopes - equivalent to Rails scopes
var messageScopeIDs = map[string]uint{
	"default":                    110,
	"maintenance":                10,
	"validation_error":           16,
	"select_number":              18,
	"bad_talk_response":          21,
	"admin_default":              23,
	"new_moon_tomorrow":          27,
	"new_moon_today":             28,
	"full_moon_tomorrow":         29,
	"full_moon_today":            30,
	"duplicate_send":             31,
	"no_wishes":                  62,
	"todays_g_message":           63,
	"todays_weekly_g_message":    93,
	"todays_experience_g_message": 109,
	"select_broadcast_range":     121,
	"over_post_capacity":         122,
	"unavailable":                125,
	"lets_customize_feeling_button": 127,
	"todays_weekly_blog_g_message": 130,
}

// GetMessageByScope returns a message by its scope name.
func GetMessageByScope(scope string) *Message {
	id, ok := messageScopeIDs[scope]
	if !ok {
		return nil
	}
	var msg Message
	if err := database.DB.First(&msg, id).Error; err != nil {
		return nil
	}
	return &msg
}

// FindMessageByID finds a message by ID.
func FindMessageByID(id uint) (*Message, error) {
	var msg Message
	if err := database.DB.First(&msg, id).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}
