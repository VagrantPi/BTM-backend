package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) CreateThreshold(db *gorm.DB, c domain.BTMRiskControlThreshold) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	item := BTMRiskControlThresholdDomainToModel(c)
	return db.Create(&item).Error
}

func BTMRiskControlThresholdDomainToModel(item domain.BTMRiskControlThreshold) model.BTMRiskControlThreshold {
	return model.BTMRiskControlThreshold{
		Role:          item.Role.Uint8(),
		Threshold:     item.Threshold,
		ThresholdDays: item.ThresholdDays,
	}
}

func BTMRiskControlThresholdModelToDomain(item model.BTMRiskControlThreshold) domain.BTMRiskControlThreshold {
	return domain.BTMRiskControlThreshold{
		Role:          domain.RiskControlRole(item.Role),
		Threshold:     item.Threshold,
		ThresholdDays: item.ThresholdDays,
	}
}
