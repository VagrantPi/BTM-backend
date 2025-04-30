package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BTMRiskControlCustomerLimitSetting struct {
	ID                   uint            `gorm:"primarykey;"`
	CustomerId           uuid.UUID       `gorm:"index; not null;"`
	Role                 uint8           `gorm:"not null; comment:'角色權限'"`
	DailyLimit           decimal.Decimal `gorm:"not null; comment:'日交易限額'"`
	MonthlyLimit         decimal.Decimal `gorm:"not null; comment:'月交易限額'"`
	Level1               decimal.Decimal `gorm:"default:0; comment:'等級門檻一'"`
	Level2               decimal.Decimal `gorm:"default:0; comment:'等級門檻二'"`
	Level1Days           uint32          `gorm:"comment:'等級門檻一的天數'"`
	Level2Days           uint32          `gorm:"comment:'等級門檻二的天數'"`
	VelocityDays         uint32          `gorm:"comment:'交易次數限制天數'"`
	VelocityTimes        uint32          `gorm:"comment:'交易次數限制次數'"`
	IsCustomized         bool            `gorm:"not null; default:false; comment:'是否有客製化限額'"`
	IsCustomizedEdd      bool            `gorm:"not null; default:false; comment:'是否有客製化EDD'"`
	IsCustomizedVelocity bool            `gorm:"not null; default:false; comment:'是否有客製化交易次數限制'"`
	EddAt                sql.NullTime    `gorm:"comment:'觸發EDD時間'"`
	IsEdd                bool            `gorm:"not null; default:false; comment:'是否為EDD'"`
	EddType              string          `gorm:"not null; default:''; comment:'EDD類型'"`
	LastBlackToNormalAt  sql.NullTime    `gorm:"comment:'最後一次黑名單切回正常名單的時間'"`
	LastRole             uint8           `gorm:"comment:'前一次角色權限'"`
	ChangeRoleReason     string          `gorm:"comment:'角色權限變更原因'"`
	ChangeLimitReason    string          `gorm:"comment:'限額變更原因'"`
	ChangeVelocityReason string          `gorm:"comment:'交易次數限制變更原因'"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
