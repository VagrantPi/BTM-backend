package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) GetCashIns(db *gorm.DB, customerID, phone string, startAt, endAt time.Time, limit int, page int) ([]domain.CashInTxWithInfo, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []domain.CashInTxWithInfo{}

	sql := db.Model(&model.CashInTx{}).
		Select(
			"cash_in_txs.*",
			"devices.name AS device_name",
			"btm_invoices.invoice_no",
		).
		Joins("LEFT JOIN devices ON cash_in_txs.device_id = devices.device_id").
		Joins("LEFT JOIN btm_invoices ON cash_in_txs.id::TEXT = btm_invoices.tx_id").
		Where("cash_in_txs.fiat != 0")

	if customerID != "" {
		sql = sql.Where("cash_in_txs.customer_id::TEXT LIKE ?", "%"+customerID+"%")
	}
	if phone != "" {
		sql = sql.Where("cash_in_txs.phone LIKE ?", "%"+phone+"%")
	}
	if !startAt.IsZero() && !endAt.IsZero() {
		sql = sql.Where("cash_in_txs.created BETWEEN ? AND ?", startAt, endAt)
	}

	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Order("cash_in_txs.created DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, int(total), nil
}

func (repo *repository) GetCashInTxBySessionId(db *gorm.DB, sessionId string) (*domain.CashInTx, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var cashInTx domain.CashInTx
	err := db.Model(&model.CashInTx{}).
		Where("cash_in_txs.id::TEXT = ?", sessionId).
		First(&cashInTx).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &cashInTx, nil
}
