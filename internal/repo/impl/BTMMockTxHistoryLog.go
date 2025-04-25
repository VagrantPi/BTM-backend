package impl

import (
	"BTM-backend/pkg/error_code"

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
