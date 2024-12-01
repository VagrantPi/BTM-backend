package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BTMWhitelist struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	CustomerID uuid.UUID `gorm:"type:uuid;uniqueIndex:customer-coin-address,priority:1"`
	CryptoCode string    `gorm:"uniqueIndex:customer-coin-address,priority:2"`
	Address    string    `gorm:"uniqueIndex:customer-coin-address,priority:3"`
}
