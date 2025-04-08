package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) CreateLoginLog(db *gorm.DB, log domain.BTMLoginLog) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	item := BTMLoginLogDomainToModel(log)
	return db.Create(&item).Error
}

func (repo *repository) GetLoginLogs(db *gorm.DB, limit int, page int) ([]domain.BTMLoginLog, int64, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	sql := db.Model(&model.BTMLoginLog{})
	var logs []model.BTMLoginLog

	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.BTMLoginLog, 0, len(logs))
	for _, logItem := range logs {
		resp = append(resp, BTMLoginLogModelToDomain(logItem))
	}
	return resp, int64(total), nil
}

func (repo *repository) IsAddressExistsInAfterValue(db *gorm.DB, address string) (bool, error) {
	if db == nil {
		return false, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM btm_change_logs 
			WHERE UPPER(convert_from(after_value, 'UTF8')::json->>'address') = UPPER($1)
		)
	`

	var exists bool
	// 執行查詢並將結果存入 exists 變數
	err := db.Raw(query, address).Scan(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}

func BTMLoginLogDomainToModel(item domain.BTMLoginLog) model.BTMLoginLog {
	return model.BTMLoginLog{
		UserID:   item.UserID,
		UserName: item.UserName,
		IP:       item.IP,
		Browser:  item.Browser,
	}
}

func BTMLoginLogModelToDomain(item model.BTMLoginLog) domain.BTMLoginLog {
	return domain.BTMLoginLog{
		ID:        item.ID,
		UserID:    item.UserID,
		UserName:  item.UserName,
		IP:        item.IP,
		Browser:   item.Browser,
		CreatedAt: item.CreatedAt,
	}
}
