package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) GetCashIns(db *gorm.DB, customerID string, startAt, endAt time.Time, limit int, page int) ([]domain.CashInTx, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []model.CashInTx{}

	sql := db.Model(&model.CashInTx{}).Where("fiat != 0")

	if customerID != "" {
		sql = sql.Where("customer_id::TEXT LIKE ?", "%"+customerID+"%")
	}
	if !startAt.IsZero() && !endAt.IsZero() {
		sql = sql.Where("created BETWEEN ? AND ?", startAt, endAt)
	}

	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Order("created DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.CashInTx, 0, len(list))
	for _, tx := range list {
		resp = append(resp, CashInTxModelToDomain(tx))
	}
	return resp, int(total), nil
}

func CashInTxModelToDomain(model model.CashInTx) domain.CashInTx {
	return domain.CashInTx{
		ID:                   model.ID,
		DeviceID:             model.DeviceID,
		ToAddress:            model.ToAddress,
		CryptoAtoms:          model.CryptoAtoms,
		CryptoCode:           model.CryptoCode,
		Fiat:                 model.Fiat,
		FiatCode:             model.FiatCode,
		Fee:                  model.Fee,
		TxHash:               model.TxHash,
		Phone:                model.Phone,
		Error:                model.Error,
		Created:              model.Created,
		Send:                 model.Send,
		SendConfirmed:        model.SendConfirmed,
		TimedOut:             model.TimedOut,
		SendTime:             model.SendTime,
		ErrorCode:            model.ErrorCode,
		OperatorCompleted:    model.OperatorCompleted,
		SendPending:          model.SendPending,
		CashInFee:            model.CashInFee,
		MinimumTx:            model.MinimumTx,
		CustomerID:           model.CustomerID,
		TxVersion:            model.TxVersion,
		TermsAccepted:        model.TermsAccepted,
		CommissionPercentage: model.CommissionPercentage,
		RawTickerPrice:       model.RawTickerPrice,
		IsPaperWallet:        model.IsPaperWallet,
		Discount:             model.Discount,
		BatchID:              model.BatchID,
		Batched:              model.Batched,
		BatchTime:            model.BatchTime,
		DiscountSource:       model.DiscountSource,
		TxCustomerPhotoAt:    model.TxCustomerPhotoAt,
		TxCustomerPhotoPath:  model.TxCustomerPhotoPath,
		WalletScore:          model.WalletScore,
		Email:                model.Email,
	}
}

func CashInTxDomainToModel(tx domain.CashInTx) model.CashInTx {
	return model.CashInTx{
		ID:                   tx.ID,
		DeviceID:             tx.DeviceID,
		ToAddress:            tx.ToAddress,
		CryptoAtoms:          tx.CryptoAtoms,
		CryptoCode:           tx.CryptoCode,
		Fiat:                 tx.Fiat,
		FiatCode:             tx.FiatCode,
		Fee:                  tx.Fee,
		TxHash:               tx.TxHash,
		Phone:                tx.Phone,
		Error:                tx.Error,
		Created:              tx.Created,
		Send:                 tx.Send,
		SendConfirmed:        tx.SendConfirmed,
		TimedOut:             tx.TimedOut,
		SendTime:             tx.SendTime,
		ErrorCode:            tx.ErrorCode,
		OperatorCompleted:    tx.OperatorCompleted,
		SendPending:          tx.SendPending,
		CashInFee:            tx.CashInFee,
		MinimumTx:            tx.MinimumTx,
		CustomerID:           tx.CustomerID,
		TxVersion:            tx.TxVersion,
		TermsAccepted:        tx.TermsAccepted,
		CommissionPercentage: tx.CommissionPercentage,
		RawTickerPrice:       tx.RawTickerPrice,
		IsPaperWallet:        tx.IsPaperWallet,
		Discount:             tx.Discount,
		BatchID:              tx.BatchID,
		Batched:              tx.Batched,
		BatchTime:            tx.BatchTime,
		DiscountSource:       tx.DiscountSource,
		TxCustomerPhotoAt:    tx.TxCustomerPhotoAt,
		TxCustomerPhotoPath:  tx.TxCustomerPhotoPath,
		WalletScore:          tx.WalletScore,
		Email:                tx.Email,
	}
}
