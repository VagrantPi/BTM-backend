package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type BTMRiskControlLimitSetting struct {
	ID            uint            `gorm:"primarykey"`
	Role          uint8           `gorm:"uniqueIndex; not null; comment:'角色權限'"`
	DailyLimit    decimal.Decimal `gorm:"not null; comment:'預設日交易限額'"`
	MonthlyLimit  decimal.Decimal `gorm:"not null; comment:'預設月交易限額'"`
	Level1        decimal.Decimal `gorm:"comment:'預設等級門檻一'"`
	Level2        decimal.Decimal `gorm:"comment:'預設等級門檻二'"`
	Level1Days    uint32          `gorm:"comment:'預設等級門檻一天數'"`
	Level2Days    uint32          `gorm:"comment:'預設等級門檻二天數'"`
	VelocityDays  uint32          `gorm:"comment:'預設交易次數限制天數'"`
	VelocityTimes uint32          `gorm:"comment:'預設交易次數限制次數'"`
	ChangeReason  string          `gorm:"comment:'變更原因'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
