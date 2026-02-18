package model

import (
	"time"

	"github.com/RyokouKanai/gomethod/database"
	"github.com/RyokouKanai/gomethod/encrypt"
)

// Wish represents a user's wish (dream or solution type).
type Wish struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"column:user_id" json:"user_id"`
	Content     *string   `gorm:"type:text" json:"content"`
	WishType    string    `gorm:"column:wish_type" json:"wish_type"`
	Salt        *string   `gorm:"column:salt" json:"salt"`
	S3ObjectURL *string   `gorm:"column:s3_object_url;type:text" json:"s3_object_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Wish) TableName() string { return "wishes" }

func (w *Wish) PlainContent() string {
	if w.Content == nil || w.Salt == nil {
		if w.Content != nil {
			return *w.Content
		}
		return ""
	}
	plain, err := encrypt.Decrypt(*w.Content, *w.Salt)
	if err != nil {
		if w.Content != nil {
			return *w.Content
		}
		return ""
	}
	return plain
}

func (w *Wish) EncryptContent() error {
	if w.Content == nil {
		return nil
	}
	enc, salt, err := encrypt.Encrypt(*w.Content)
	if err != nil {
		return err
	}
	w.Content = &enc
	w.Salt = &salt
	return nil
}

func (w *Wish) GetS3ObjectURL() string {
	if w.S3ObjectURL != nil {
		return *w.S3ObjectURL
	}
	return ""
}

// Hate represents something a user dislikes.
type Hate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	Content   *string   `gorm:"type:text" json:"content"`
	Salt      *string   `gorm:"column:salt" json:"salt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Hate) TableName() string { return "hates" }

func (h *Hate) PlainContent() string {
	if h.Content == nil || h.Salt == nil {
		if h.Content != nil {
			return *h.Content
		}
		return ""
	}
	plain, err := encrypt.Decrypt(*h.Content, *h.Salt)
	if err != nil {
		if h.Content != nil {
			return *h.Content
		}
		return ""
	}
	return plain
}

func (h *Hate) EncryptContent() error {
	if h.Content == nil {
		return nil
	}
	enc, salt, err := encrypt.Encrypt(*h.Content)
	if err != nil {
		return err
	}
	h.Content = &enc
	h.Salt = &salt
	return nil
}

// Happiness represents something that made a user happy.
type Happiness struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	Content   *string   `gorm:"type:text" json:"content"`
	Salt      *string   `gorm:"column:salt" json:"salt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Happiness) TableName() string { return "happiness" }

func (hp *Happiness) PlainContent() string {
	if hp.Content == nil || hp.Salt == nil {
		if hp.Content != nil {
			return *hp.Content
		}
		return ""
	}
	plain, err := encrypt.Decrypt(*hp.Content, *hp.Salt)
	if err != nil {
		if hp.Content != nil {
			return *hp.Content
		}
		return ""
	}
	return plain
}

func (hp *Happiness) EncryptContent() error {
	if hp.Content == nil {
		return nil
	}
	enc, salt, err := encrypt.Encrypt(*hp.Content)
	if err != nil {
		return err
	}
	hp.Content = &enc
	hp.Salt = &salt
	return nil
}

// CreateWish creates a new encrypted wish.
func CreateWish(userID uint, content, wishType string) (*Wish, error) {
	w := &Wish{UserID: userID, Content: &content, WishType: wishType}
	if err := w.EncryptContent(); err != nil {
		return nil, err
	}
	return w, database.DB.Create(w).Error
}

// CreateHate creates a new encrypted hate.
func CreateHate(userID uint, content string) (*Hate, error) {
	h := &Hate{UserID: userID, Content: &content}
	if err := h.EncryptContent(); err != nil {
		return nil, err
	}
	return h, database.DB.Create(h).Error
}

// CreateHappiness creates a new encrypted happiness.
func CreateHappiness(userID uint, content string) (*Happiness, error) {
	hp := &Happiness{UserID: userID, Content: &content}
	if err := hp.EncryptContent(); err != nil {
		return nil, err
	}
	return hp, database.DB.Create(hp).Error
}
