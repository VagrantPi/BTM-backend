package model

import (
	"BTM-backend/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// 用來存放每次限額功能多塞入的假交易紀錄
type BTMRiskControlMachineRequestLimitLog struct {
	ID             uint                              `gorm:"primarykey"`
	CustomerId     uuid.UUID                         `gorm:"index; not null"`
	TxId           uuid.UUID                         `gorm:"index; not null; comment:'交易ID, 非 txhash'"`
	DailyAddTx     decimal.Decimal                   `gorm:"not null; comment:'日交易額插入的假交易'"`
	MonthlyAddTx   decimal.Decimal                   `gorm:"not null; comment:'月交易額插入的假交易'"`
	NowLimitConfig domain.BTMRiskControlLimitSetting `gorm:"not null; comment:'當前限額設定'"`
	CreatedAt      time.Time
}
