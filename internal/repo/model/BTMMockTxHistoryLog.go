package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BTMMockTxHistoryLog struct {
	Id        int64     `gorm:"column:id;primarykey;autoIncrement"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;index"`

	CustomerID          uuid.UUID       `gorm:"column:customer_id;not null;index"`
	DeviceId            string          `gorm:"column:device_id;not null"`
	DefaultDailyLimit   decimal.Decimal `gorm:"column:default_daily_limit;not null; comment:'預設日限额'"`
	DefaultMonthlyLimit decimal.Decimal `gorm:"column:default_monthly_limit;not null; comment:'預設月限额'"`
	LimitDailyLimit     decimal.Decimal `gorm:"column:limit_daily_limit;not null; comment:'客製日限额'"`
	LimitMonthlyLimit   decimal.Decimal `gorm:"column:limit_monthly_limit;not null; comment:'客製月限额'"`
	DayLimit            decimal.Decimal `gorm:"column:day_limit;not null; comment:'日限额塞入的假交易'"`
	MonthLimit          decimal.Decimal `gorm:"column:month_limit;not null; comment:'月限额塞入的假交易'"`
	StartAt             string          `gorm:"column:start_at;not null; comment:'撈取交易紀錄開始時間，當解禁時間比較晚，則會是解禁時間'"`
	BanExpireDateRaw    string          `gorm:"column:ban_expire_date_raw; comment:'用戶解禁時間（民國）'"`
	BanExpireDate       string          `gorm:"column:ban_expire_date; comment:'用戶解禁時間'"`
}
