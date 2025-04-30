package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BTMMockTxHistoryLog struct {
	Id                  int64           `json:"id"`
	CreatedAt           time.Time       `json:"created_at"`
	CustomerID          uuid.UUID       `json:"customer_id"`
	DeviceId            string          `json:"device_id"`
	DefaultDailyLimit   decimal.Decimal `json:"default_daily_limit"`
	DefaultMonthlyLimit decimal.Decimal `json:"default_monthly_limit"`
	LimitDailyLimit     decimal.Decimal `json:"limit_daily_limit"`
	LimitMonthlyLimit   decimal.Decimal `json:"limit_monthly_limit"`
	DayLimit            decimal.Decimal `json:"day_limit"`
	MonthLimit          decimal.Decimal `json:"month_limit"`
	StartAt             string          `json:"start_at"`
	BanExpireDateRaw    string          `json:"ban_expire_date_raw"`
	BanExpireDate       string          `json:"ban_expire_date"`
}
