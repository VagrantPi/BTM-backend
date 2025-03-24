package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type BTMRiskControlThreshold struct {
	ID            uint
	Role          RiskControlRole
	Threshold     decimal.Decimal
	ThresholdDays uint8
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
