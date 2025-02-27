package domain

import "github.com/google/uuid"

type BTMSumsub struct {
	CustomerId uuid.UUID
	Info       SumsubData
	IdNumber   string
}
