package model

import (
	"time"

	"github.com/google/uuid"
)

type BTMInvoice struct {
	ID           uint      `gorm:"primarykey"`
	CustomerId   uuid.UUID `gorm:"index; not null"`
	TxId         string    `gorm:"index; not null"`
	InvoiceNo    string
	InvoiceDate  time.Time
	RandomNumber string
	RawResp      string `gorm:"not null"`
	CreatedAt    time.Time
}
