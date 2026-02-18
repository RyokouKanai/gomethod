package model

import (
	"github.com/RyokouKanai/gomethod/database"
)

type ReplyPattern struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	SentMessageID   uint   `gorm:"column:sent_message_id" json:"sent_message_id"`
	Position        *int   `gorm:"column:position" json:"position"`
	NextMessageID   uint   `gorm:"column:next_message_id" json:"next_message_id"`
	ExecutionMethod string `gorm:"column:execution_method;default:base" json:"execution_method"`
}

func (ReplyPattern) TableName() string { return "reply_patterns" }

// GetNextMessage returns the next message for this reply pattern.
func (rp *ReplyPattern) GetNextMessage() *Message {
	msg, _ := FindMessageByID(rp.NextMessageID)
	return msg
}

// GetSentMessage returns the sent message for this reply pattern.
func (rp *ReplyPattern) GetSentMessage() *Message {
	msg, _ := FindMessageByID(rp.SentMessageID)
	return msg
}

// FindReplyPatternByID finds a reply pattern by ID.
func FindReplyPatternByID(id int) *ReplyPattern {
	var rp ReplyPattern
	if err := database.DB.First(&rp, id).Error; err != nil {
		return nil
	}
	return &rp
}

// FindReplyPatternByMessageAndPosition finds a reply pattern by sent message ID and position.
func FindReplyPatternByMessageAndPosition(sentMessageID uint, position int) *ReplyPattern {
	var rp ReplyPattern
	err := database.DB.Where("sent_message_id = ? AND position = ?", sentMessageID, position).First(&rp).Error
	if err != nil {
		return nil
	}
	return &rp
}

// FindFirstReplyPatternByMessage finds the first reply pattern for a sent message.
func FindFirstReplyPatternByMessage(sentMessageID uint) *ReplyPattern {
	var rp ReplyPattern
	err := database.DB.Where("sent_message_id = ?", sentMessageID).First(&rp).Error
	if err != nil {
		return nil
	}
	return &rp
}
