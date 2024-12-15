package model

import (
	"BTM-backend/internal/domain"
	"time"
)

type UserConfig struct {
	ID      int64                 `gorm:"type:int4;primary_key"`
	Type    string                `gorm:"type:text"`
	Data    domain.UserConfigJSON `gorm:"type:json"`
	Created time.Time             `gorm:"type:timestamptz;not null"`
}

func (UserConfig) TableName() string {
	return "user_config"
}
