package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
)

func (repo *repository) GetLatestConfData() (resp domain.UserConfigJSON, err error) {
	var info model.UserConfig
	err = repo.db.Where("type = 'config'").Order("created DESC").Limit(1).Find(&info).Error
	if err != nil {
		return resp, err
	}

	return info.Data, nil
}
