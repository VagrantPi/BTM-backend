package model

import (
	"BTM-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BTMSumsub struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	CustomerId  uuid.UUID         `gorm:"not null"`
	ApplicantId string            `gorm:"not null"`
	Info        domain.SumsubData `gorm:"type:json; not null"`
	IdNumber    string            `gorm:"uniqueIndex; not null"`
}
