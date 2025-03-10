package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BTMRiskControlCustomerLimitSetting struct {
	ID                  uint            `gorm:"primarykey"`
	CustomerId          uuid.UUID       `gorm:"not null"`
	Role                uint8           `gorm:"not null"`
	DailyLimit          decimal.Decimal `gorm:"not null"`
	MonthlyLimit        decimal.Decimal `gorm:"not null"`
	IsCustomized        bool            `gorm:"not null; default:false"`
	LastBlackToNormalAt sql.NullTime
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
