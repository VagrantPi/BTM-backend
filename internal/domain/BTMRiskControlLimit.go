package domain

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

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
	ID            uint            `json:"id"`
	Role          RiskControlRole `json:"role"`
	DailyLimit    decimal.Decimal `json:"daily_limit"`
	MonthlyLimit  decimal.Decimal `json:"monthly_limit"`
	Level1        decimal.Decimal `json:"level1"`
	Level2        decimal.Decimal `json:"level2"`
	Level1Days    uint32          `json:"level1_days"`
	Level2Days    uint32          `json:"level2_days"`
	VelocityDays  uint32          `json:"velocity_days"`
	VelocityTimes uint32          `json:"velocity_times"`
	ChangeReason  string          `json:"change_reason"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (s *BTMRiskControlLimitSetting) Scan(value interface{}) error {
	if value == nil {
		*s = BTMRiskControlLimitSetting{}
		return nil
	}

	bytesValue, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan BTMRiskControlLimitSetting")
	}

	data := BTMRiskControlLimitSetting{}
	if err := json.Unmarshal(bytesValue, &data); err != nil {
		return errors.New("invalid scan BTMRiskControlLimitSetting unmarshal")
	}

	*s = data
	return nil
}

// Value return json value, implement driver.Valuer interface
func (s BTMRiskControlLimitSetting) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type BTMRiskControlCustomerLimitSetting struct {
	ID                   uint            `json:"id"`
	Role                 RiskControlRole `json:"role"`
	CustomerId           uuid.UUID       `json:"customer_id"`
	DailyLimit           decimal.Decimal `json:"daily_limit"`
	MonthlyLimit         decimal.Decimal `json:"monthly_limit"`
	Level1               decimal.Decimal `json:"level1"`
	Level2               decimal.Decimal `json:"level2"`
	Level1Days           uint32          `json:"level1_days"`
	Level2Days           uint32          `json:"level2_days"`
	VelocityDays         uint32          `json:"velocity_days"`
	VelocityTimes        uint32          `json:"velocity_times"`
	IsCustomized         bool            `json:"is_customized"`
	IsCustomizedEdd      bool            `json:"is_customized_edd"`
	IsCustomizedVelocity bool            `json:"is_customized_velocity"`
	EddAt                sql.NullTime    `json:"edd_at"`
	IsEdd                bool            `json:"is_edd"`
	EddType              string          `json:"edd_type"`
	ChangeRoleReason     string          `json:"change_role_reason"`
	ChangeLimitReason    string          `json:"change_limit_reason"`
	ChangeVelocityReason string          `json:"change_velocity_reason"`
}
