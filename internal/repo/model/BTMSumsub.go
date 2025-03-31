package model

import (
	"database/sql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BTMSumsub struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	CustomerId       uuid.UUID     `gorm:"index; not null"`
	ApplicantId      string        `gorm:"not null"`
	IdNumber         string        `gorm:"uniqueIndex; not null"`
	BanExpireDate    sql.NullInt64 `gorm:"index"`
	Phone            string
	InspectionId     string
	IdCardFrontImgId string
	IdCardBackImgId  string
	SelfieImgId      string
	Name             string
	EmailHash        string `gorm:"index"`
	InfoHash         string
	Status           string `gorm:"index"`
}
