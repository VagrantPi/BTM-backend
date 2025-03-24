package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BTMRiskControlCustomerLimitSetting struct {
	ID                  uint            `gorm:"primarykey;"`
	CustomerId          uuid.UUID       `gorm:"index; not null;"`
	Role                uint8           `gorm:"not null; comment:'角色權限'"`
	DailyLimit          decimal.Decimal `gorm:"not null; comment:'日交易限額'"`
	MonthlyLimit        decimal.Decimal `gorm:"not null; comment:'月交易限額'"`
	IsCustomized        bool            `gorm:"not null; default:false; comment:'是否有客製化限額'"`
	LastBlackToNormalAt sql.NullTime    `gorm:"comment:'最後一次黑名單切回正常名單的時間'"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
