package domain

import "github.com/google/uuid"

type BTMSumsub struct {
	CustomerId  uuid.UUID
	ApplicantId string
	Info        SumsubData
	IdNumber    string
}
