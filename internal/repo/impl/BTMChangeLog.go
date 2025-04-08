package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"strings"
	"time"

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

func (repo *repository) BatchCreateBTMChangeLog(db *gorm.DB, c []domain.BTMChangeLog) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	items := make([]model.BTMChangeLog, 0, len(c))
	for _, item := range c {
		items = append(items, BTMChangeLogDomainToModel(item))
	}
	return db.Create(&items).Error
}

func (repo *repository) GetBTMChangeLogs(db *gorm.DB, tableName, customerId string, startAt, endAt time.Time, limit int, page int) ([]domain.BTMChangeLog, int64, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []model.BTMChangeLog{}

	sql := db.Model(&model.BTMChangeLog{})
	if strings.TrimSpace(tableName) != "" {
		sql = sql.Where("table_name = ?", tableName)
	}

	if strings.TrimSpace(customerId) != "" {
		sql = sql.Where("customer_id = ?", customerId)
	}

	if !startAt.IsZero() {
		sql = sql.Where("created_at >= ?", startAt)
	}

	if !endAt.IsZero() {
		sql = sql.Where("created_at <= ?", endAt)
	}

	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.BTMChangeLog, 0, len(list))
	for _, logItem := range list {
		resp = append(resp, BTMChangeLogModelToDomain(logItem))
	}
	return resp, int64(total), nil
}

func (repo *repository) UpdateBTMChangeLog(db *gorm.DB, id uint, c domain.BTMChangeLog) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	item := BTMChangeLogDomainToModel(c)
	return db.
		Model(&model.BTMChangeLog{}).
		Where("id = ?", id).
		Updates(item).Error
}

func (repo *repository) AddressExistsInAfterValue(db *gorm.DB, address string) (domain.BTMChangeLog, error) {
	if db == nil {
		return domain.BTMChangeLog{}, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	query := `
			SELECT *
			FROM btm_change_logs 
			WHERE UPPER(convert_from(after_value, 'UTF8')::json->>'address') = UPPER($1)
	`

	var fetched model.BTMChangeLog
	// 執行查詢並將結果存入 fetched 變數
	err := db.Raw(query, address).Scan(&fetched).Error
	if err != nil {
		return domain.BTMChangeLog{}, err
	}

	return BTMChangeLogModelToDomain(fetched), nil
}

func BTMChangeLogDomainToModel(item domain.BTMChangeLog) model.BTMChangeLog {
	return model.BTMChangeLog{
		OperationUserId: item.OperationUserId,
		TableName:       string(item.TableName),
		OperationType:   uint8(item.OperationType),
		CustomerId:      item.CustomerId,
		BeforeValue:     item.BeforeValue,
		AfterValue:      item.AfterValue,
	}
}

func BTMChangeLogModelToDomain(item model.BTMChangeLog) domain.BTMChangeLog {
	return domain.BTMChangeLog{
		ID:              item.ID,
		OperationUserId: item.OperationUserId,
		TableName:       domain.BTMChangeLogTableName(item.TableName),
		OperationType:   domain.BTMChangeLogOperationType(item.OperationType),
		CustomerId:      item.CustomerId,
		BeforeValue:     item.BeforeValue,
		AfterValue:      item.AfterValue,
		CreatedAt:       item.CreatedAt,
	}
}
