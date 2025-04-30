package model

import (
	"BTM-backend/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ServerLog struct {
	ID        uuid.UUID   `gorm:"column:id;type:uuid;primaryKey;not null"`
	DeviceID  *string     `gorm:"column:device_id;type:text"`
	LogLevel  *string     `gorm:"column:log_level;type:text"`
	Timestamp time.Time   `gorm:"column:timestamp;type:timestamptz;default:now()"`
	Message   *string     `gorm:"column:message;type:text"`
	Meta      domain.JSON `gorm:"column:meta;type:json"`
}
