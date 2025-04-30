package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) RemoveExtraMockTxHistoryLog(db *gorm.DB) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	return db.Exec(`
		DELETE FROM btm_mock_tx_history_logs
		WHERE id NOT IN (
			SELECT MIN(id)
			FROM btm_mock_tx_history_logs
			GROUP BY
				customer_id,
				device_id,
				default_daily_limit,
				default_monthly_limit,
				limit_daily_limit,
				limit_monthly_limit,
				day_limit,
				month_limit,
				start_at
		);
	`).Error
}

func (repo *repository) GetMockTxHistoryLogs(db *gorm.DB, limit int, page int, startAt, endAt time.Time) ([]domain.BTMMockTxHistoryLog, int64, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	sql := db.Model(&model.BTMMockTxHistoryLog{})

	if !startAt.IsZero() && !endAt.IsZero() {
		sql = sql.Where("created_at BETWEEN ? AND ?", startAt, endAt)
	}

	var total int64
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var list []model.BTMMockTxHistoryLog
	err := sql.
		Limit(limit).
		Offset(offset).
		Find(&list).
		Error
	if err != nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, err.Error())
	}

	resp := make([]domain.BTMMockTxHistoryLog, 0, len(list))
	for _, log := range list {
		resp = append(resp, BTMMockTxHistoryLogModelToDomain(log))
	}
	return resp, total, nil
}

func BTMMockTxHistoryLogModelToDomain(log model.BTMMockTxHistoryLog) domain.BTMMockTxHistoryLog {
	return domain.BTMMockTxHistoryLog{
		Id:                  log.Id,
		CustomerID:          log.CustomerID,
		DeviceId:            log.DeviceId,
		DefaultDailyLimit:   log.DefaultDailyLimit,
		DefaultMonthlyLimit: log.DefaultMonthlyLimit,
		LimitDailyLimit:     log.LimitDailyLimit,
		LimitMonthlyLimit:   log.LimitMonthlyLimit,
		DayLimit:            log.DayLimit,
		MonthLimit:          log.MonthLimit,
		StartAt:             log.StartAt,
		CreatedAt:           log.CreatedAt,
	}
}
