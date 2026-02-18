package model

import (
	"time"

	"github.com/RyokouKanai/gomethod/database"
	"github.com/RyokouKanai/gomethod/encrypt"
)

// LastMessage stores the user's last message for context tracking.
type LastMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	Content   string    `gorm:"type:text" json:"content"`
	Salt      *string   `gorm:"column:salt" json:"salt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (LastMessage) TableName() string { return "last_messages" }

func (lm *LastMessage) PlainContent() string {
	if lm.Salt == nil || lm.Content == "" {
		return lm.Content
	}
	plain, err := encrypt.Decrypt(lm.Content, *lm.Salt)
	if err != nil {
		return lm.Content
	}
	return plain
}

// ActionRecord tracks user actions like thanks count.
type ActionRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"column:user_id" json:"user_id"`
	ThanksCount int       `gorm:"column:thanks_count;default:0" json:"thanks_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ActionRecord) TableName() string { return "action_records" }

// TalkHistory tracks conversation history between a user and the bot.
type TalkHistory struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"column:user_id" json:"user_id"`
	MessageID      uint      `gorm:"column:message_id" json:"message_id"`
	ReplyPatternID *int      `gorm:"column:reply_pattern_id" json:"reply_pattern_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (TalkHistory) TableName() string { return "talk_histories" }

// GetMessage returns the message associated with this talk history.
func (th *TalkHistory) GetMessage() *Message {
	msg, _ := FindMessageByID(th.MessageID)
	return msg
}

// GetReplyPattern returns the reply pattern for this talk history.
func (th *TalkHistory) GetReplyPattern() *ReplyPattern {
	if th.ReplyPatternID == nil {
		return nil
	}
	return FindReplyPatternByID(*th.ReplyPatternID)
}

// MoonPhase represents lunar phase data.
type MoonPhase struct {
	ID    uint      `gorm:"primaryKey" json:"id"`
	Phase string    `gorm:"column:phase" json:"phase"`
	Date  time.Time `gorm:"column:date;type:date" json:"date"`
}

func (MoonPhase) TableName() string { return "moon_phases" }

// GetMoonPhaseToday returns today's moon phase, if any.
func GetMoonPhaseToday() *MoonPhase {
	var mp MoonPhase
	today := time.Now().Format("2006-01-02")
	if err := database.DB.Where("date = ?", today).First(&mp).Error; err != nil {
		return nil
	}
	return &mp
}

// GetMoonPhaseTomorrow returns tomorrow's moon phase, if any.
func GetMoonPhaseTomorrow() *MoonPhase {
	var mp MoonPhase
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	if err := database.DB.Where("date = ?", tomorrow).First(&mp).Error; err != nil {
		return nil
	}
	return &mp
}

// BatchExecutionHistory tracks batch execution for deduplication.
type BatchExecutionHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Batch     string    `gorm:"column:batch" json:"batch"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BatchExecutionHistory) TableName() string { return "batch_execution_histories" }

// IsToday checks if the batch was already executed today.
func (b *BatchExecutionHistory) IsToday() bool {
	now := time.Now()
	bod := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	eod := bod.Add(24 * time.Hour)
	return b.UpdatedAt.After(bod) && b.UpdatedAt.Before(eod)
}

// Option represents selectable options for a message.
type Option struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	MessageID uint    `gorm:"column:message_id" json:"message_id"`
	Position  int     `gorm:"column:position" json:"position"`
	Content   *string `gorm:"type:text" json:"content"`
}

func (Option) TableName() string { return "options" }

func (o *Option) GetContent() string {
	if o.Content != nil {
		return *o.Content
	}
	return ""
}

// FindOptionByMessageAndPosition finds an option by message ID and position.
func FindOptionByMessageAndPosition(messageID uint, position int) *Option {
	var o Option
	err := database.DB.Where("message_id = ? AND position = ?", messageID, position).First(&o).Error
	if err != nil {
		return nil
	}
	return &o
}

// FeelingSetting represents customizable feeling buttons.
type FeelingSetting struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"column:user_id" json:"user_id"`
	ButtonNumber int       `gorm:"column:button_number" json:"button_number"`
	Content      string    `gorm:"column:content" json:"content"`
	Salt         *string   `gorm:"column:salt" json:"salt"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (FeelingSetting) TableName() string { return "feeling_settings" }

func (fs *FeelingSetting) PlainContent() string {
	if fs.Salt == nil || fs.Content == "" {
		return fs.Content
	}
	plain, err := encrypt.Decrypt(fs.Content, *fs.Salt)
	if err != nil {
		return fs.Content
	}
	return plain
}

func (fs *FeelingSetting) EncryptContent() error {
	enc, salt, err := encrypt.Encrypt(fs.Content)
	if err != nil {
		return err
	}
	fs.Content = enc
	fs.Salt = &salt
	return nil
}

// DefaultFeelingSettings holds the default feeling button configurations.
var DefaultFeelingSettings = []struct {
	ButtonNumber int
	Content      string
}{
	{1, "嫌だ！"},
	{2, "ムカつく！"},
	{3, "悔しい！"},
	{4, "クソ！"},
	{5, "辛いよ"},
}

// ThanksLevel represents milestones for thanks counts.
type ThanksLevel struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Count    int     `gorm:"column:count" json:"count"`
	Cheering *string `gorm:"type:text" json:"cheering"`
}

func (ThanksLevel) TableName() string { return "thanks_levels" }

// FindThanksLevelByCount finds a thanks level by count.
func FindThanksLevelByCount(count int) *ThanksLevel {
	var tl ThanksLevel
	if err := database.DB.Where("count = ?", count).First(&tl).Error; err != nil {
		return nil
	}
	return &tl
}

// Article represents content articles.
type Article struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	ArticleTypeID uint   `gorm:"column:article_type_id" json:"article_type_id"`
	Title         string `gorm:"column:title" json:"title"`
}

func (Article) TableName() string { return "articles" }

// GetSections returns the sections of this article.
func (a *Article) GetSections() ([]Section, error) {
	var sections []Section
	err := database.DB.Where("article_id = ?", a.ID).Order("position ASC").Find(&sections).Error
	return sections, err
}

// ArticleType represents article categories.
type ArticleType struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (ArticleType) TableName() string { return "article_types" }

// Section represents a section within an article.
type Section struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	ArticleID uint    `gorm:"column:article_id" json:"article_id"`
	Position  int     `gorm:"column:position;default:1" json:"position"`
	Content   *string `gorm:"type:text" json:"content"`
}

func (Section) TableName() string { return "sections" }

func (s *Section) GetContent() string {
	if s.Content != nil {
		return *s.Content
	}
	return ""
}

// Lesson represents a learning lesson.
type Lesson struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Position int    `gorm:"column:position" json:"position"`
	Title    string `gorm:"column:title" json:"title"`
}

func (Lesson) TableName() string { return "lessons" }

// GetArticles returns experience articles for this lesson.
func (l *Lesson) GetArticles() ([]Article, error) {
	var articles []Article
	err := database.DB.
		Joins("JOIN lesson_articles ON lesson_articles.article_id = articles.id").
		Joins("JOIN article_types ON article_types.id = articles.article_type_id").
		Where("lesson_articles.lesson_id = ? AND article_types.name = ?", l.ID, "experience").
		Find(&articles).Error
	return articles, err
}

// LessonArticle is a join table between lessons and articles.
type LessonArticle struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	LessonID  uint `gorm:"column:lesson_id" json:"lesson_id"`
	ArticleID uint `gorm:"column:article_id" json:"article_id"`
}

func (LessonArticle) TableName() string { return "lesson_articles" }

// Plan represents a subscription plan.
type Plan struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Identifier string `gorm:"column:identifier" json:"identifier"`
	Name       string `gorm:"column:name" json:"name"`
}

func (Plan) TableName() string { return "plans" }
