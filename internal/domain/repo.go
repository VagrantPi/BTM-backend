package domain

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	NewTransactionBegin(ctx context.Context) (*gorm.DB, error)
	Close(tx *gorm.DB)
	Rollback(tx *gorm.DB)
	TransactionCommit(tx *gorm.DB) error
	GetDb(ctx context.Context) *gorm.DB
	NewTxWithContext(ctx context.Context) (tx *gorm.DB, err error)

	/**
	 * BTM
	 **/

	// BTMUser
	GetBTMUserByAccount(db *gorm.DB, account string) (*BTMUser, error)
	InitAdmin(db *gorm.DB) error

	// BTMWhitelist
	CreateWhitelist(db *gorm.DB, whitelist *BTMWhitelist) error
	GetWhiteListByCustomerId(db *gorm.DB, customerID uuid.UUID, limit int, page int) (list []BTMWhitelist, total int64, err error)
	UpdateWhitelist(db *gorm.DB, whitelist *BTMWhitelist) error
	DeleteWhitelist(db *gorm.DB, id int64) error
	SearchWhitelistByAddress(db *gorm.DB, customerID uuid.UUID, address string, limit int, page int) (list []BTMWhitelist, total int64, err error)

	// BTMLoginToken
	IsLastLoginToken(db *gorm.DB, userID uint, loginToken string) (bool, error)
	CreateOrUpdateLastLoginToken(db *gorm.DB, userID uint, loginToken string) error
	DeleteLastLoginToken(db *gorm.DB, userID uint) error

	// BTMRole
	InitRawRole(db *gorm.DB) error
	GetRawRoleByRoleName(db *gorm.DB, roleName string) (role BTMRole, err error)
	GetRawRoles(db *gorm.DB) (roles []BTMRole, err error)

	// BTMCIB
	UpsertBTMCIB(db *gorm.DB, cib BTMCIB) error

	// BTMSumsub
	CreateBTMSumsub(db *gorm.DB, btmsumsub BTMSumsub) error
	GetBTMSumsub(db *gorm.DB, customerId string) (*BTMSumsub, error)

	/**
	 * lamassu original
	 **/

	// customers
	GetCustomers(db *gorm.DB, limit int, page int) ([]Customer, int, error)
	GetCustomerById(db *gorm.DB, id uuid.UUID) (*Customer, error)
	SearchCustomersByPhone(db *gorm.DB, phone string, limit int, page int) ([]Customer, int, error)
	GetSumsubExternalId(db *gorm.DB, customerId string) (string, error)

	// userConfig
	GetLatestConfData(db *gorm.DB) (UserConfigJSON, error)
}
