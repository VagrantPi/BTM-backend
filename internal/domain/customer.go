package domain

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID      uuid.UUID
	Phone   string
	Created time.Time
}
