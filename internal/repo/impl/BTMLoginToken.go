package impl

import (
	"BTM-backend/internal/repo/model"

	"gorm.io/gorm/clause"
)

func (repo *repository) IsLastLoginToken(userID uint, loginToken string) (bool, error) {
	var count int64
	if err := repo.db.Model(&model.BTMLoginToken{}).
		Where("user_id = ? AND login_token = ? AND deleted_at IS NULL", userID, loginToken).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *repository) CreateOrUpdateLastLoginToken(userID uint, loginToken string) error {
	return repo.db.Clauses(clause.OnConflict{
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

func (repo *repository) DeleteLastLoginToken(userID uint) error {
	return repo.db.Delete(&model.BTMLoginToken{}, "user_id = ? ", userID).Error
}
