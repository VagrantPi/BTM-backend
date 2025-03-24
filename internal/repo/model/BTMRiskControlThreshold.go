package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type BTMRiskControlThreshold struct {
	ID            uint            `gorm:"primarykey"`
	Role          uint8           `gorm:"uniqueIndex:role_threshold,priority:1; not null; comment:'角色權限'"`
	Threshold     decimal.Decimal `gorm:"uniqueIndex:role_threshold,priority:2; not null; comment:'風控門檻閾值'"`
	ThresholdDays uint8           `gorm:"not null; comment:'風控門檻累計n天內'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
