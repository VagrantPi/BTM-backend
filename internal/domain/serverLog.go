package domain

import (
	"time"

	"github.com/google/uuid"
)

type ServerLog struct {
	ID        uuid.UUID `json:"id"`
	DeviceID  *string   `json:"device_id"`
	LogLevel  *string   `json:"log_level"`
	Timestamp time.Time `json:"timestamp"`
	Message   *string   `json:"message"`
	Meta      JSON      `json:"meta"`
}
