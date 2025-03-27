package domain

import (
	"database/sql"

	"github.com/google/uuid"
)

type BTMSumsub struct {
	CustomerId       uuid.UUID
	ApplicantId      string
	IdNumber         string
	BanExpireDate    sql.NullInt64
	Phone            string
	InspectionId     string
	IdCardFrontImgId string
	IdCardBackImgId  string
	SelfieImgId      string
	Name             string
	EmailHash        string
	InfoHash         string
}
