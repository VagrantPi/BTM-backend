package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type BTMRiskControlLimitSetting struct {
	ID           uint            `gorm:"primarykey"`
	Role         uint8           `gorm:"uniqueIndex; not null"`
	DailyLimit   decimal.Decimal `gorm:"not null"`
	MonthlyLimit decimal.Decimal `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
