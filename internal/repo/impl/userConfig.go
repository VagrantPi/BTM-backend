package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"

	"gorm.io/gorm"
)

func (repo *repository) GetLatestConfData(db *gorm.DB) (resp domain.UserConfigJSON, err error) {
	var info model.UserConfig
	err = db.Where("type = 'config'").Order("created DESC").Limit(1).Find(&info).Error
	if err != nil {
		return resp, err
	}

	return info.Data, nil
}
