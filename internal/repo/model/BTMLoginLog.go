package model

import (
	"time"
)

type BTMLoginLog struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null"`
	UserName  string    `gorm:"not null"` // 反正規化
	IP        string    `gorm:"not null"`
	Browser   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"index; not null"`
}
