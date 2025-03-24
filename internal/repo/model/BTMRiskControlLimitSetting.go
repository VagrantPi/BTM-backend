package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type BTMRiskControlLimitSetting struct {
	ID           uint            `gorm:"primarykey"`
	Role         uint8           `gorm:"uniqueIndex; not null; comment:'角色權限'"`
	DailyLimit   decimal.Decimal `gorm:"not null; comment:'日交易限額'"`
	MonthlyLimit decimal.Decimal `gorm:"not null; comment:'月交易限額'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
