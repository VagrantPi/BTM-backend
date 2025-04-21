package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type BTMDailyDeviceIncome struct {
	ID uint `gorm:"primarykey"`

	StatDate           time.Time       `gorm:"uniqueIndex:idx_stat_device; not null; comment:'統計時間'"`
	DeviceId           string          `gorm:"uniqueIndex:idx_stat_device; not null; comment:'裝置 ID'"`
	TotalFiat          decimal.Decimal `gorm:"not null; default: 0; comment:'總收入'"`
	AllDeviceTotalFiat decimal.Decimal `gorm:"not null; default: 0; comment:'所有裝置統計日總收入'"`
}
