package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) GetRiskControlLimitSetting(db *gorm.DB) ([]domain.BTMRiskControlLimitSetting, error) {
	if db == nil {
		return nil, errors.BadRequest(error_code.ErrInvalidRequest, "db is nil")
	}

	results := []model.BTMRiskControlLimitSetting{}
	if err := db.Model(&model.BTMRiskControlLimitSetting{}).Scan(&results).Error; err != nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "GetRiskControlLimitSetting").WithCause(err)
	}

	reps := []domain.BTMRiskControlLimitSetting{}
	for _, v := range results {
		reps = append(reps, BTMRiskControlLimitSettingModelToDomain(v))
	}

	return reps, nil
}

func (repo *repository) UpdateRiskControlLimitSetting(db *gorm.DB, operationUserId int64, newSetting domain.BTMRiskControlLimitSetting, reason string) error {
	if db == nil {
		return errors.BadRequest(error_code.ErrInvalidRequest, "db is nil")
	}

	var beforeLimit model.BTMRiskControlLimitSetting
	if err := db.Where("role = ?", domain.RiskControlRoleWhite).First(&beforeLimit).Error; err != nil {
		return err
	}
	beforeLimitJsonData, err := json.Marshal(beforeLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(beforeLimit)").WithCause(err)
	}
	afterLimit := model.BTMRiskControlLimitSetting{
		ID:            beforeLimit.ID,
		Role:          beforeLimit.Role,
		DailyLimit:    newSetting.DailyLimit,
		MonthlyLimit:  newSetting.MonthlyLimit,
		Level1:        newSetting.Level1,
		Level2:        newSetting.Level2,
		Level1Days:    newSetting.Level1Days,
		Level2Days:    newSetting.Level2Days,
		VelocityDays:  newSetting.VelocityDays,
		VelocityTimes: newSetting.VelocityTimes,
		ChangeReason:  reason,
	}
	afterLimitJsonData, err := json.Marshal(afterLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(afterLimit)").WithCause(err)
	}
	err = repo.CreateBTMChangeLog(db, domain.BTMChangeLog{
		OperationUserId: operationUserId,
		TableName:       domain.BTMChangeLogTableNameBTMRiskControlLimitSetting,
		OperationType:   domain.BTMChangeLogOperationTypeUpdate,
		CustomerId:      nil,
		BeforeValue:     beforeLimitJsonData,
		AfterValue:      afterLimitJsonData,
	})
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err)
	}

	if err := db.Save(&afterLimit).Error; err != nil {
		return errors.InternalServer(error_code.ErrDBError, "UpdateRiskControlLimitSetting").WithCause(err)
	}

	return nil
}

func BTMRiskControlLimitSettingModelToDomain(m model.BTMRiskControlLimitSetting) domain.BTMRiskControlLimitSetting {
	return domain.BTMRiskControlLimitSetting{
		ID:            m.ID,
		Role:          domain.RiskControlRole(m.Role),
		DailyLimit:    m.DailyLimit,
		MonthlyLimit:  m.MonthlyLimit,
		Level1:        m.Level1,
		Level2:        m.Level2,
		Level1Days:    m.Level1Days,
		Level2Days:    m.Level2Days,
		VelocityDays:  m.VelocityDays,
		VelocityTimes: m.VelocityTimes,
	}
}
