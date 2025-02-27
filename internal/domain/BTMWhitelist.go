package domain

import (
	"time"

	"github.com/google/uuid"
)

type BTMWhitelist struct {
	ID         uint64    `json:"id"`
	CustomerID uuid.UUID `json:"customer_id" gorm:"index:customer-coin-address,priority:1"`
	CryptoCode string    `json:"crypto_code" gorm:"index:customer-coin-address,priority:2"`
	Address    string    `json:"address" gorm:"uniqueIndex; index:customer-coin-address,priority:3"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
