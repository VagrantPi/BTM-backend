package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) GetLatestConfData(db *gorm.DB) (resp domain.UserConfigJSON, err error) {
	if db == nil {
		return domain.UserConfigJSON{}, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var info model.UserConfig
	err = db.Where("type = 'config'").Order("created DESC").Limit(1).Find(&info).Error
	if err != nil {
		return resp, err
	}

	return info.Data, nil
}
