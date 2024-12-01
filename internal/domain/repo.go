package domain

import "github.com/google/uuid"

type Repository interface {
	// customers
	GetCustomers(limit int, page int) ([]Customer, int, error)
	GetCustomerById(id uuid.UUID) (*Customer, error)

	// BTMUser
	GetBTMUserByAccount(account string) (*BTMUser, error)

	// Whitelist
	CreateWhitelist(whitelist *BTMWhitelist) error
	GetWhiteListByCustomerId(customerID uuid.UUID, limit int, page int) (list []BTMWhitelist, total int64, err error)
	UpdateWhitelist(whitelist *BTMWhitelist) error
	DeleteWhitelist(id int64) error
}
