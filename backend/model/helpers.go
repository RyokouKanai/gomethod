package model

import (
	"github.com/RyokouKanai/gomethod/database"
)

// UpdateTalkHistoryReplyPattern updates the reply pattern ID of a talk history.
func UpdateTalkHistoryReplyPattern(th *TalkHistory) error {
	return database.DB.Model(th).Update("reply_pattern_id", th.ReplyPatternID).Error
}

// GetExperienceArticles returns all experience articles.
func GetExperienceArticles() ([]Article, error) {
	var articles []Article
	err := database.DB.
		Joins("JOIN article_types ON article_types.id = articles.article_type_id").
		Where("article_types.name = ?", "experience").
		Find(&articles).Error
	return articles, err
}

// FindLessonByID finds a lesson by ID.
func FindLessonByID(id uint) *Lesson {
	var l Lesson
	if err := database.DB.First(&l, id).Error; err != nil {
		return nil
	}
	return &l
}

// UpdateWishS3URL updates the S3 object URL for a wish.
func UpdateWishS3URL(wishID uint, url string) error {
	return database.DB.Model(&Wish{}).Where("id = ?", wishID).Update("s3_object_url", url).Error
}

// FindWishByID finds a wish by ID.
func FindWishByID(id uint) *Wish {
	var w Wish
	if err := database.DB.First(&w, id).Error; err != nil {
		return nil
	}
	return &w
}

// FindHateByID finds a hate by ID.
func FindHateByID(id uint) *Hate {
	var h Hate
	if err := database.DB.First(&h, id).Error; err != nil {
		return nil
	}
	return &h
}

// FindHappinessByID finds a happiness by ID.
func FindHappinessByID(id uint) *Happiness {
	var h Happiness
	if err := database.DB.First(&h, id).Error; err != nil {
		return nil
	}
	return &h
}

// FindFeelingSettingByID finds a feeling setting by ID.
func FindFeelingSettingByID(id uint) *FeelingSetting {
	var fs FeelingSetting
	if err := database.DB.First(&fs, id).Error; err != nil {
		return nil
	}
	return &fs
}

// DeleteHatesByUserID deletes all hates for a user.
func DeleteHatesByUserID(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&Hate{}).Error
}

// FindGMessageByID finds a GMessage by ID.
func FindGMessageByID(id uint) *GMessage {
	var g GMessage
	if err := database.DB.First(&g, id).Error; err != nil {
		return nil
	}
	return &g
}

// UpdateGMessageContent updates a GMessage's encrypted content.
func UpdateGMessageContent(g *GMessage) error {
	if err := g.EncryptContent(); err != nil {
		return err
	}
	return database.DB.Save(g).Error
}

// DeleteGMessage deletes a GMessage.
func DeleteGMessage(id uint) error {
	return database.DB.Delete(&GMessage{}, id).Error
}

// UpdateWishContent updates a wish's content with encryption.
func UpdateWishContent(w *Wish) error {
	if err := w.EncryptContent(); err != nil {
		return err
	}
	return database.DB.Save(w).Error
}

// UpdateHateContent updates a hate's content with encryption.
func UpdateHateContent(h *Hate) error {
	if err := h.EncryptContent(); err != nil {
		return err
	}
	return database.DB.Save(h).Error
}

// UpdateFeelingSettingContent updates a feeling setting's content with encryption.
func UpdateFeelingSettingContent(fs *FeelingSetting) error {
	if err := fs.EncryptContent(); err != nil {
		return err
	}
	return database.DB.Save(fs).Error
}

// FindFeelingSettingByUserAndButton finds a feeling setting by user and button number.
func FindFeelingSettingByUserAndButton(userID uint, buttonNumber int) *FeelingSetting {
	var fs FeelingSetting
	err := database.DB.Where("user_id = ? AND button_number = ?", userID, buttonNumber).First(&fs).Error
	if err != nil {
		return nil
	}
	return &fs
}

// CheckBatchDuplicateExecution checks if a batch was already executed today.
func CheckBatchDuplicateExecution(batchName string) bool {
	var beh BatchExecutionHistory
	err := database.DB.Where("batch = ?", batchName).First(&beh).Error
	if err != nil {
		return false
	}
	if beh.IsToday() {
		return true
	}
	// Touch the record
	database.DB.Model(&beh).Update("updated_at", database.DB.NowFunc())
	return false
}
