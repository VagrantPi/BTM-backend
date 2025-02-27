package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CashInTx struct {
	ID                   uuid.UUID       `gorm:"type:uuid;primaryKey;not null"`
	DeviceID             string          `gorm:"type:text;not null"`
	ToAddress            string          `gorm:"type:text;not null"`
	CryptoAtoms          decimal.Decimal `gorm:"type:numeric(30,0);not null"`
	CryptoCode           string          `gorm:"type:text;not null"`
	Fiat                 decimal.Decimal `gorm:"type:numeric(14,5);not null"`
	FiatCode             string          `gorm:"type:text;not null"`
	Fee                  int64           `gorm:"type:int8"`
	TxHash               string          `gorm:"type:text"`
	Phone                string          `gorm:"type:text"`
	Error                string          `gorm:"type:text"`
	Created              time.Time       `gorm:"type:timestamptz;not null;default:now()"`
	Send                 bool            `gorm:"type:bool;not null;default:false"`
	SendConfirmed        bool            `gorm:"type:bool;not null;default:false"`
	TimedOut             bool            `gorm:"type:bool;not null;default:false"`
	SendTime             time.Time       `gorm:"type:timestamptz"`
	ErrorCode            string          `gorm:"type:text"`
	OperatorCompleted    bool            `gorm:"type:bool;not null;default:false"`
	SendPending          bool            `gorm:"type:bool;not null;default:false"`
	CashInFee            decimal.Decimal `gorm:"type:numeric(14,5);not null"`
	MinimumTx            int             `gorm:"type:int4;not null"`
	CustomerID           uuid.UUID       `gorm:"type:uuid;default:'47ac1184-8102-11e7-9079-8f13a7117867'"`
	TxVersion            int             `gorm:"type:int4;not null"`
	TermsAccepted        bool            `gorm:"type:bool;not null;default:false"`
	CommissionPercentage decimal.Decimal `gorm:"type:numeric(14,5)"`
	RawTickerPrice       decimal.Decimal `gorm:"type:numeric(14,5)"`
	IsPaperWallet        bool            `gorm:"type:bool;default:false"`
	Discount             int16           `gorm:"type:int2"`
	BatchID              uuid.UUID       `gorm:"type:uuid"`
	Batched              bool            `gorm:"type:bool;not null;default:false"`
	BatchTime            time.Time       `gorm:"type:timestamptz"`
	DiscountSource       string          `gorm:"type:discount_source"`
	TxCustomerPhotoAt    time.Time       `gorm:"type:timestamptz"`
	TxCustomerPhotoPath  string          `gorm:"type:text"`
	WalletScore          int16           `gorm:"type:int2"`
	Email                string          `gorm:"type:text"`
}

func (CashInTx) TableName() string {
	return "cash_in_txs"
}
