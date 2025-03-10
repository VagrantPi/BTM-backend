package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RiskControlRole uint8

const (
	RiskControlRoleInit RiskControlRole = iota
	RiskControlRoleWhite
	RiskControlRoleGray
	RiskControlRoleBlack
)

func (e RiskControlRole) String() string {
	switch e {
	case RiskControlRoleInit:
		return "未設定"
	case RiskControlRoleWhite:
		return "白名單"
	case RiskControlRoleGray:
		return "灰名單"
	case RiskControlRoleBlack:
		return "黑名單"
	default:
		return "未知權限"
	}
}

func (e RiskControlRole) Uint8() uint8 {
	return uint8(e)
}

type RiskControlRoleKeyValue struct {
	Id   uint8  `json:"id"`
	Name string `json:"name"`
}

type BTMRiskControlLimitSetting struct {
	ID           uint            `json:"id"`
	Role         RiskControlRole `json:"role"`
	DailyLimit   decimal.Decimal `json:"daily_limit"`
	MonthlyLimit decimal.Decimal `json:"monthly_limit"`
}

type BTMRiskControlLimitSettingChange struct {
	ID              uint            `json:"id"`
	OldRole         RiskControlRole `json:"old_role"`
	OldDailyLimit   decimal.Decimal `json:"old_daily_limit"`
	OldMonthlyLimit decimal.Decimal `json:"old_monthly_limit"`
	NewRole         RiskControlRole `json:"new_role"`
	NewDailyLimit   decimal.Decimal `json:"new_daily_limit"`
	NewMonthlyLimit decimal.Decimal `json:"new_monthly_limit"`
	CreatedAt       time.Time       `json:"created_at"`
}

type BTMRiskControlCustomerLimitSetting struct {
	ID           uint            `json:"id"`
	Role         RiskControlRole `json:"role"`
	CustomerId   uuid.UUID       `json:"customer_id"`
	DailyLimit   decimal.Decimal `json:"daily_limit"`
	MonthlyLimit decimal.Decimal `json:"monthly_limit"`
	IsCustomized bool            `json:"is_customized"`
}

type BTMRiskControlCustomerLimitSettingChange struct {
	ID              uint            `json:"id"`
	CustomerId      uuid.UUID       `json:"customer_id"`
	OldRole         RiskControlRole `json:"old_role"`
	OldDailyLimit   decimal.Decimal `json:"old_daily_limit"`
	OldMonthlyLimit decimal.Decimal `json:"old_monthly_limit"`
	NewRole         RiskControlRole `json:"new_role"`
	NewDailyLimit   decimal.Decimal `json:"new_daily_limit"`
	NewMonthlyLimit decimal.Decimal `json:"new_monthly_limit"`
	CreatedAt       time.Time       `json:"created_at"`
}
