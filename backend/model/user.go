package model

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/RyokouKanai/gomethod/database"
)

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LineUserID  string    `gorm:"column:line_user_id;not null" json:"line_user_id"`
	MemberType  string    `gorm:"column:member_type;default:basic" json:"member_type"`
	PlanID      int64     `gorm:"column:plan_id;default:1" json:"plan_id"`
	DisplayName *string   `gorm:"column:display_name" json:"display_name"`
	PictureURL  *string   `gorm:"column:picture_url" json:"picture_url"`
	IsActive    bool      `gorm:"column:is_active;default:true" json:"is_active"`
	IsShik      bool      `gorm:"column:is_shik;default:false" json:"is_shik"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

func (u *User) IsAdmin() bool {
	return u.MemberType == "admin"
}

// FindOrCreateByLineUserID finds or creates a user by LINE user ID.
func FindOrCreateByLineUserID(lineUserID string) (*User, error) {
	var user User
	result := database.DB.Where("line_user_id = ?", lineUserID).First(&user)
	if result.Error != nil {
		user = User{LineUserID: lineUserID}
		if err := database.DB.Create(&user).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil
}

// GetMasterUser returns the admin user.
func GetMasterUser() (*User, error) {
	var user User
	if err := database.DB.Where("member_type = ?", "admin").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetActiveUsers returns all active users.
func GetActiveUsers() ([]User, error) {
	var users []User
	if err := database.DB.Where("is_active = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetShikUsers returns all shik users.
func GetShikUsers() ([]User, error) {
	var users []User
	if err := database.DB.Where("is_shik = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetLastMessage returns the user's last message.
func (u *User) GetLastMessage() (*LastMessage, error) {
	var lm LastMessage
	err := database.DB.Where("user_id = ?", u.ID).First(&lm).Error
	if err != nil {
		return nil, err
	}
	return &lm, nil
}

// UpsertLastMessage creates or updates the user's last message.
func (u *User) UpsertLastMessage(message string) error {
	var lm LastMessage
	result := database.DB.Where("user_id = ?", u.ID).First(&lm)
	if result.Error != nil {
		lm = LastMessage{UserID: u.ID, Content: message}
		return database.DB.Create(&lm).Error
	}
	lm.Content = message
	return database.DB.Save(&lm).Error
}

// GetActionRecord returns the user's action record.
func (u *User) GetActionRecord() (*ActionRecord, error) {
	var ar ActionRecord
	err := database.DB.Where("user_id = ?", u.ID).First(&ar).Error
	if err != nil {
		return nil, err
	}
	return &ar, nil
}

// UpsertActionRecord increments the specified column.
func (u *User) UpsertActionRecord(column string, point int) error {
	var ar ActionRecord
	result := database.DB.Where("user_id = ?", u.ID).First(&ar)
	if result.Error != nil {
		ar = ActionRecord{UserID: u.ID, ThanksCount: point}
		return database.DB.Create(&ar).Error
	}
	return database.DB.Model(&ar).Update(column, ar.ThanksCount+point).Error
}

// GetThanksCount returns the thanks count for this user.
func (u *User) GetThanksCount() int {
	ar, err := u.GetActionRecord()
	if err != nil {
		return 0
	}
	return ar.ThanksCount
}

// ResetThanksCount resets the user's thanks count to 0.
func (u *User) ResetThanksCount() error {
	return database.DB.Model(&ActionRecord{}).Where("user_id = ?", u.ID).Update("thanks_count", 0).Error
}

// CreateTalkHistory creates a new talk history entry.
func (u *User) CreateTalkHistory(message *Message) (*TalkHistory, error) {
	th := TalkHistory{
		UserID:    u.ID,
		MessageID: message.ID,
	}
	if err := database.DB.Create(&th).Error; err != nil {
		return nil, err
	}
	// Destroy oldest if >= 4
	var count int64
	database.DB.Model(&TalkHistory{}).Where("user_id = ?", u.ID).Count(&count)
	if count >= 4 {
		var oldest TalkHistory
		database.DB.Where("user_id = ?", u.ID).Order("created_at ASC").First(&oldest)
		database.DB.Delete(&oldest)
	}
	return &th, nil
}

// GetRecentTalkHistories returns recent talk histories ordered desc.
func (u *User) GetRecentTalkHistories() ([]TalkHistory, error) {
	var histories []TalkHistory
	err := database.DB.Where("user_id = ?", u.ID).Order("created_at DESC").Find(&histories).Error
	return histories, err
}

// GetLatestTalkHistory returns the most recent talk history.
func (u *User) GetLatestTalkHistory() (*TalkHistory, error) {
	var th TalkHistory
	err := database.DB.Where("user_id = ?", u.ID).Order("created_at DESC").First(&th).Error
	if err != nil {
		return nil, err
	}
	return &th, nil
}

// GetDreamWishes returns all dream-type wishes for the user.
func (u *User) GetDreamWishes() ([]Wish, error) {
	var wishes []Wish
	err := database.DB.Where("user_id = ? AND wish_type = ?", u.ID, "dream").Find(&wishes).Error
	return wishes, err
}

// GetSolutionWishes returns all solution-type wishes for the user.
func (u *User) GetSolutionWishes() ([]Wish, error) {
	var wishes []Wish
	err := database.DB.Where("user_id = ? AND wish_type = ?", u.ID, "solution").Find(&wishes).Error
	return wishes, err
}

// GetHates returns all hates for the user.
func (u *User) GetHates() ([]Hate, error) {
	var hates []Hate
	err := database.DB.Where("user_id = ?", u.ID).Find(&hates).Error
	return hates, err
}

// GetHappiness returns all happiness for the user.
func (u *User) GetHappiness() ([]Happiness, error) {
	var happiness []Happiness
	err := database.DB.Where("user_id = ?", u.ID).Find(&happiness).Error
	return happiness, err
}

// GetFeelingSettings returns the user's feeling settings.
func (u *User) GetFeelingSettings() ([]FeelingSetting, error) {
	var settings []FeelingSetting
	err := database.DB.Where("user_id = ?", u.ID).Find(&settings).Error
	return settings, err
}

// CreateFeelingSettings creates default feeling settings for the user.
func (u *User) CreateFeelingSettings() error {
	for _, d := range DefaultFeelingSettings {
		fs := FeelingSetting{
			UserID:       u.ID,
			ButtonNumber: d.ButtonNumber,
			Content:      d.Content,
		}
		if err := database.DB.Create(&fs).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetGMessageHistories returns g_message_histories for the user.
func (u *User) GetGMessageHistories() ([]GMessageHistory, error) {
	var histories []GMessageHistory
	err := database.DB.Where("user_id = ?", u.ID).Find(&histories).Error
	return histories, err
}

// GetGMessageHistoriesByPeriod returns g_message_histories filtered by period.
func (u *User) GetGMessageHistoriesByPeriod(period string) ([]GMessageHistory, error) {
	var histories []GMessageHistory
	err := database.DB.
		Joins("JOIN g_messages ON g_messages.id = g_message_histories.g_message_id").
		Where("g_message_histories.user_id = ? AND g_messages.period = ?", u.ID, period).
		Find(&histories).Error
	return histories, err
}

// CreateGMessageHistory creates a new g_message_history entry.
func (u *User) CreateGMessageHistory(gMessage *GMessage) error {
	h := GMessageHistory{
		UserID:     u.ID,
		GMessageID: gMessage.ID,
	}
	return database.DB.Create(&h).Error
}

// FetchGMessageByPeriod fetches a random unsent g_message of the given period.
func (u *User) FetchGMessageByPeriod(period string) (*GMessage, error) {
	histories, err := u.GetGMessageHistoriesByPeriod(period)
	if err != nil {
		return nil, err
	}

	sentIDs := make([]uint, 0, len(histories))
	for _, h := range histories {
		sentIDs = append(sentIDs, h.GMessageID)
	}

	var leftMessages []GMessage
	query := database.DB.Where("period = ?", period)
	if len(sentIDs) > 0 {
		query = query.Where("id NOT IN ?", sentIDs)
	}
	query.Find(&leftMessages)

	if len(leftMessages) == 0 {
		// Reset histories for this period
		if len(sentIDs) > 0 {
			database.DB.Where("user_id = ? AND g_message_id IN ?", u.ID, sentIDs).Delete(&GMessageHistory{})
		}
		database.DB.Where("period = ?", period).Find(&leftMessages)
	}

	if len(leftMessages) == 0 {
		return nil, fmt.Errorf("no g_messages found for period: %s", period)
	}

	return &leftMessages[rand.Intn(len(leftMessages))], nil
}

// SaveProfile fetches and saves the user's LINE profile.
func (u *User) SaveProfile() error {
	token := os.Getenv("LINE_CHANNEL_TOKEN")
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.line.me/v2/bot/profile/%s", u.LineUserID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var profile struct {
		DisplayName string `json:"displayName"`
		PictureURL  string `json:"pictureUrl"`
	}
	if err := json.Unmarshal(body, &profile); err != nil {
		return err
	}

	u.DisplayName = &profile.DisplayName
	u.PictureURL = &profile.PictureURL
	return database.DB.Save(u).Error
}
