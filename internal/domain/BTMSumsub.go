package domain

import (
	"database/sql"

	"github.com/google/uuid"
)

type BTMSumsub struct {
	CustomerId    uuid.UUID
	ApplicantId   string
	Info          SumsubData
	IdNumber      string
	BanExpireDate sql.NullInt64
	Email         string
	Phone         string
}
