package domain

import (
	"time"

	"github.com/google/uuid"
)

type CustomerType int64

const (
	CustomerTypeNone CustomerType = iota + 1
	CustomerTypeWhiteList
	CustomerTypeGrayList
	CustomerTypeBlackList
)

type Customer struct {
	ID             uuid.UUID
	Phone          string
	Created        time.Time
	SuspendedUntil string
}

type CustomerWithWhiteListCreated struct {
	ID                    uuid.UUID `json:"id"`
	Phone                 string    `json:"phone"`
	EmailHash             string    `json:"email_hash"`
	InfoHash              string    `json:"info_hash"`
	Name                  string    `json:"name"`
	Created               time.Time `json:"created_at"`
	FirstWhiteListCreated time.Time `json:"first_white_list_created"`
	IsCibBlock            bool      `json:"is_cib_block"`
	EddType               string    `json:"edd_type"`
	ChangeRoleReason      string    `json:"change_role_reason"`
}

type CustomerAuthorizedOverride string

func (c CustomerAuthorizedOverride) String() string {
	return string(c)
}

const (
	CustomerAuthorizedOverrideVerified  CustomerAuthorizedOverride = "verified"
	CustomerAuthorizedOverrideBlocked   CustomerAuthorizedOverride = "blocked"
	CustomerAuthorizedOverrideAutomatic CustomerAuthorizedOverride = "automatic"
)
