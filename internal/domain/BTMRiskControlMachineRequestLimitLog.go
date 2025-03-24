package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BTMRiskControlMachineRequestLimitLog struct {
	ID             uint
	CustomerId     uuid.UUID
	TxId           uuid.UUID
	DailyAddTx     decimal.Decimal
	MonthlyAddTx   decimal.Decimal
	NowLimitConfig BTMRiskControlLimitSetting
	CreatedAt      time.Time
}
