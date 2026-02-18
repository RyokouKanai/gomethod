package model

import (
	"time"

	"github.com/RyokouKanai/gomethod/database"
	"github.com/RyokouKanai/gomethod/encrypt"
)

type GMessage struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	Content *string `gorm:"type:text" json:"content"`
	Salt    *string `gorm:"column:salt" json:"salt"`
	Period  string  `gorm:"column:period;default:daily" json:"period"`
}

func (GMessage) TableName() string { return "g_messages" }

// GetContent safely returns the content.
func (g *GMessage) GetContent() string {
	if g.Content != nil {
		return *g.Content
	}
	return ""
}

// PlainContent decrypts and returns the plain text content.
func (g *GMessage) PlainContent() string {
	if g.Content == nil || g.Salt == nil {
		return g.GetContent()
	}
	plain, err := encrypt.Decrypt(*g.Content, *g.Salt)
	if err != nil {
		return g.GetContent()
	}
	return plain
}

// EncryptContent encrypts the content before save.
func (g *GMessage) EncryptContent() error {
	if g.Content == nil {
		return nil
	}
	enc, salt, err := encrypt.Encrypt(*g.Content)
	if err != nil {
		return err
	}
	g.Content = &enc
	g.Salt = &salt
	return nil
}

// CreateGMessage creates a new GMessage with encryption.
func CreateGMessage(content, period string) (*GMessage, error) {
	g := &GMessage{Content: &content, Period: period}
	if err := g.EncryptContent(); err != nil {
		return nil, err
	}
	if err := database.DB.Create(g).Error; err != nil {
		return nil, err
	}
	return g, nil
}

// GetGMessagesByPeriod returns all g_messages for a given period.
func GetGMessagesByPeriod(period string) ([]GMessage, error) {
	var messages []GMessage
	err := database.DB.Where("period = ?", period).Find(&messages).Error
	return messages, err
}

// GMessageHistory tracks which g_messages have been sent to which users.
type GMessageHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"column:user_id" json:"user_id"`
	GMessageID uint      `gorm:"column:g_message_id" json:"g_message_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (GMessageHistory) TableName() string { return "g_message_histories" }
