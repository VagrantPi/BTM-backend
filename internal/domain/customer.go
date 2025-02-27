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

type CustomerWithWhiteListCreated struct {
	ID                    uuid.UUID `json:"id"`
	Phone                 string    `json:"phone"`
	Created               time.Time `json:"created_at"`
	FirstWhiteListCreated time.Time `json:"first_white_list_created"`
}
