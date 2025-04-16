package domain

import (
	"time"

	"github.com/google/uuid"
)

type BTMCustomerNote struct {
	ID                uint             `json:"id"`
	CreatedAt         time.Time        `json:"created_at"`
	CustomerId        uuid.UUID        `json:"customer_id"`
	Note              string           `json:"note"`
	OperationUserId   int64            `json:"operation_user_id"`
	OperationUserName string           `json:"operation_user_name"`
	NoteType          CustomerNoteType `json:"note_type"`
}
