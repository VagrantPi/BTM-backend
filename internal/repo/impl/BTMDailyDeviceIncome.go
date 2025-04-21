package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) SnapshotYesterday(db *gorm.DB) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	// 呼叫時會是當日 00:30
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return err
	}

	now := time.Now().In(loc)
	fmt.Println("now", now)
	yesterday := now.AddDate(0, 0, -1)
	fmt.Println("yesterday", yesterday)
	dateStr := yesterday.Format("2006-01-02")
	fmt.Println("dateStr", dateStr)

	return repo.SnapshotByDate(db, dateStr)
}

func (repo *repository) SnapshotByDate(db *gorm.DB, dateStr string) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	query := `
		WITH all_devices AS (
			SELECT DISTINCT device_id FROM cash_in_txs
		),
		device_totals AS (
			SELECT
				device_id,
				SUM(fiat) AS total_fiat
			FROM cash_in_txs
			WHERE created >= ($1::date AT TIME ZONE 'UTC+8')
				AND created < (($1::date + INTERVAL '1 day') AT TIME ZONE 'UTC+8')
				AND fiat > 0
				AND NOT (
					(NOT send_confirmed)
					AND (created <= now() - INTERVAL '60 minutes')
				)
			GROUP BY device_id
		),
		final_result AS (
			SELECT
				d.device_id,
				COALESCE(t.total_fiat, 0) AS total_fiat
			FROM all_devices d
			LEFT JOIN device_totals t ON d.device_id = t.device_id
		)
		INSERT INTO btm_daily_device_incomes (stat_date, device_id, total_fiat, all_device_total_fiat)
		SELECT
			$1::date AS stat_date,
			device_id,
			total_fiat,
			(SELECT SUM(total_fiat) FROM final_result) AS all_device_total_fiat
		FROM final_result
		ON CONFLICT (stat_date, device_id)
		DO UPDATE SET 
			total_fiat = EXCLUDED.total_fiat,
			all_device_total_fiat = EXCLUDED.all_device_total_fiat;
	`

	return db.Exec(query, dateStr).Error
}

func (repo *repository) SnapshotRange(db *gorm.DB, startDateStr, endDateStr string) error {
	const layout = "2006-01-02"

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return fmt.Errorf("無法解析 startDate: %w", err)
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return fmt.Errorf("無法解析 endDate: %w", err)
	}

	if endDate.Before(startDate) {
		return fmt.Errorf("endDate 不可早於 startDate")
	}

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format(layout)
		if err := repo.SnapshotByDate(db, dateStr); err != nil {
			return fmt.Errorf("快照 %s 失敗: %w", dateStr, err)
		}
	}

	return nil
}

func (repo *repository) FetchByStatDate(db *gorm.DB, startDate, endDate string) ([]domain.BTMDailyDeviceIncome, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	fmt.Println("startDate", startDate)
	fmt.Println("endDate", endDate)
	results := []model.BTMDailyDeviceIncome{}
	err := db.
		Where("stat_date >= ? AND stat_date <= ?", startDate, endDate).
		Order("stat_date, device_id").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	resp := make([]domain.BTMDailyDeviceIncome, 0, len(results))
	for _, item := range results {
		resp = append(resp, BTMDailyDeviceIncomeModelToDomain(item))
	}
	return resp, nil
}

func BTMDailyDeviceIncomeDomainToModel(item domain.BTMDailyDeviceIncome) model.BTMDailyDeviceIncome {
	return model.BTMDailyDeviceIncome{
		StatDate:           item.StatDate,
		DeviceId:           item.DeviceId,
		TotalFiat:          item.TotalFiat,
		AllDeviceTotalFiat: item.AllDeviceTotalFiat,
	}
}

func BTMDailyDeviceIncomeModelToDomain(item model.BTMDailyDeviceIncome) domain.BTMDailyDeviceIncome {
	return domain.BTMDailyDeviceIncome{
		StatDate:           item.StatDate,
		DeviceId:           item.DeviceId,
		TotalFiat:          item.TotalFiat,
		AllDeviceTotalFiat: item.AllDeviceTotalFiat,
	}
}
