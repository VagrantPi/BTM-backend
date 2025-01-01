package impl

import (
	"BTM-backend/internal/repo/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *repository) IsLastLoginToken(db *gorm.DB, userID uint, loginToken string) (bool, error) {
	var count int64
	if err := db.Model(&model.BTMLoginToken{}).
		Where("user_id = ? AND login_token = ? AND deleted_at IS NULL", userID, loginToken).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *repository) CreateOrUpdateLastLoginToken(db *gorm.DB, userID uint, loginToken string) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"login_token": loginToken,
			"deleted_at":  nil,
		}),
	}).Create(&model.BTMLoginToken{
		UserID:     userID,
		LoginToken: loginToken,
	}).Error
}

func (repo *repository) DeleteLastLoginToken(db *gorm.DB, userID uint) error {
	return db.Delete(&model.BTMLoginToken{}, "user_id = ? ", userID).Error
}
