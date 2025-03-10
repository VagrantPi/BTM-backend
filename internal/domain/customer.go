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
	ID      uuid.UUID
	Phone   string
	Created time.Time
}

type CustomerWithWhiteListCreated struct {
	ID                    uuid.UUID `json:"id"`
	Phone                 string    `json:"phone"`
	Created               time.Time `json:"created_at"`
	FirstWhiteListCreated time.Time `json:"first_white_list_created"`
	IsLamassuBlock        bool      `json:"is_lamassu_block"`
	IsAdminBlock          bool      `json:"is_admin_block"`
	IsCibBlock            bool      `json:"is_cib_block"`
}
