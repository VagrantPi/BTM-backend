package domain

import "github.com/google/uuid"

type Customer struct {
	ID    uuid.UUID
	Phone string
}
