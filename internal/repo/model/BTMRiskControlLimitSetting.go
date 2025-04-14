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
	Level1       decimal.Decimal `gorm:"comment:'等級門檻一'"`
	Level2       decimal.Decimal `gorm:"comment:'等級門檻二'"`
	Level1Days   uint32          `gorm:"comment:'等級門檻一天數'"`
	Level2Days   uint32          `gorm:"comment:'等級門檻二天數'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
