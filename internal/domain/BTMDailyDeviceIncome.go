package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type BTMDailyDeviceIncome struct {
	StatDate           time.Time       `json:"stat_date"`
	DeviceId           string          `json:"device_id"`
	TotalFiat          decimal.Decimal `json:"total_fiat"`
	AllDeviceTotalFiat decimal.Decimal `json:"all_device_total_fiat"`
}

type DeviceData struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}
