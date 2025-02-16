package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) CreateBTMChangeLog(db *gorm.DB, c domain.BTMChangeLog) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	item := BTMChangeLogDomainToModel(c)
	return db.Create(&item).Error
}

func (repo *repository) GetBTMChangeLogs(db *gorm.DB, limit int, page int) ([]domain.BTMChangeLog, int64, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []model.BTMChangeLog{}

	sql := db.Model(&model.BTMChangeLog{})
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.BTMChangeLog, 0, len(list))
	for _, logItem := range list {
		resp = append(resp, BTMChangeLogModelToDomain(logItem))
	}
	return resp, int64(total), nil
}

func BTMChangeLogDomainToModel(item domain.BTMChangeLog) model.BTMChangeLog {
	return model.BTMChangeLog{
		OperationUserId: item.OperationUserId,
		TableName:       string(item.TableName),
		OperationType:   uint8(item.OperationType),
		BeforeValue:     item.BeforeValue,
		AfterValue:      item.AfterValue,
	}
}

func BTMChangeLogModelToDomain(item model.BTMChangeLog) domain.BTMChangeLog {
	return domain.BTMChangeLog{
		OperationUserId: item.OperationUserId,
		TableName:       domain.BTMChangeLogTableName(item.TableName),
		OperationType:   domain.BTMChangeLogOperationType(item.OperationType),
		BeforeValue:     item.BeforeValue,
		AfterValue:      item.AfterValue,
	}
}
