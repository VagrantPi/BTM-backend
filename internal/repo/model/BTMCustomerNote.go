package model

import (
	"time"

	"github.com/google/uuid"
)

type BTMCustomerNote struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index; not null"`

	CustomerId        uuid.UUID `gorm:"index; not null"`
	Note              string    `gorm:"not null"`
	OperationUserId   uint      `gorm:"not null"`
	OperationUserName string    `gorm:"not null"`
}
