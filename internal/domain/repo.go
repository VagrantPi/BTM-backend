package domain

import "github.com/google/uuid"

type Repository interface {
	/**
	 * BTM
	 **/

	// BTMUser
	GetBTMUserByAccount(account string) (*BTMUser, error)

	// BTMWhitelist
	CreateWhitelist(whitelist *BTMWhitelist) error
	GetWhiteListByCustomerId(customerID uuid.UUID, limit int, page int) (list []BTMWhitelist, total int64, err error)
	UpdateWhitelist(whitelist *BTMWhitelist) error
	DeleteWhitelist(id int64) error
	SearchWhitelistByAddress(customerID uuid.UUID, address string, limit int, page int) (list []BTMWhitelist, total int64, err error)

	// BTMLoginToken
	IsLastLoginToken(userID uint, loginToken string) (bool, error)
	CreateOrUpdateLastLoginToken(userID uint, loginToken string) error
	DeleteLastLoginToken(userID uint) error

	/**
	 * lamassu original
	 **/

	// customers
	GetCustomers(limit int, page int) ([]Customer, int, error)
	GetCustomerById(id uuid.UUID) (*Customer, error)
	SearchCustomersByPhone(phone string, limit int, page int) ([]Customer, int, error)

	// userConfig
	GetLatestConfData() (UserConfigJSON, error)
}
