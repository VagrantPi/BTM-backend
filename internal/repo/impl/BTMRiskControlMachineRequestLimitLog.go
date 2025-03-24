package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) CreateMachineRequestLimitLog(db *gorm.DB, c domain.BTMRiskControlMachineRequestLimitLog) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	item := BTMRiskControlMachineRequestLimitLogDomainToModel(c)
	return db.Create(&item).Error
}

func BTMRiskControlMachineRequestLimitLogDomainToModel(item domain.BTMRiskControlMachineRequestLimitLog) model.BTMRiskControlMachineRequestLimitLog {
	return model.BTMRiskControlMachineRequestLimitLog{
		CustomerId:     item.CustomerId,
		TxId:           item.TxId,
		DailyAddTx:     item.DailyAddTx,
		MonthlyAddTx:   item.MonthlyAddTx,
		NowLimitConfig: item.NowLimitConfig,
	}
}

func BTMRiskControlMachineRequestLimitLogModelToDomain(item model.BTMRiskControlMachineRequestLimitLog) domain.BTMRiskControlMachineRequestLimitLog {
	return domain.BTMRiskControlMachineRequestLimitLog{
		CustomerId:     item.CustomerId,
		TxId:           item.TxId,
		DailyAddTx:     item.DailyAddTx,
		MonthlyAddTx:   item.MonthlyAddTx,
		NowLimitConfig: item.NowLimitConfig,
	}
}
