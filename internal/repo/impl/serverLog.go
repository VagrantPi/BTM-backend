package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) GetServerLogs(db *gorm.DB, limit, page int, startAt, endAt time.Time) ([]domain.ServerLog, int64, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	sql := db.Model(&model.ServerLog{})
	if !startAt.IsZero() && !endAt.IsZero() {
		sql = sql.Where("timestamp BETWEEN ? AND ?", startAt, endAt)
	}

	offset := (page - 1) * limit

	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	list := []model.ServerLog{}
	if err := sql.Order("timestamp desc").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.ServerLog, 0, len(list))
	for _, item := range list {
		resp = append(resp, ServerLogModelToDomain(item))
	}
	return resp, int64(total), nil
}

func ServerLogModelToDomain(item model.ServerLog) domain.ServerLog {
	return domain.ServerLog{
		ID:        item.ID,
		DeviceID:  item.DeviceID,
		LogLevel:  item.LogLevel,
		Timestamp: item.Timestamp,
		Message:   item.Message,
		Meta:      item.Meta,
	}
}
