package model

import (
	"BTM-backend/internal/domain"
	"database/sql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BTMSumsub struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	CustomerId    uuid.UUID         `gorm:"index; not null"`
	ApplicantId   string            `gorm:"not null"`
	Info          domain.SumsubData `gorm:"type:json; not null"`
	IdNumber      string            `gorm:"uniqueIndex; not null"`
	BanExpireDate sql.NullInt64     `gorm:"index"`
}
