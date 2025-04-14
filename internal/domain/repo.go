package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	GetBTMUserById(db *gorm.DB, id uint) (*BTMUser, error)
	InitAdmin(db *gorm.DB) error
	CreateBTMUser(db *gorm.DB, user BTMUser) error
	GetUsers(db *gorm.DB, limit int, page int) (list []BTMUserWithRoles, total int64, err error)
	UpdateUserNameRoles(db *gorm.DB, id uint, account string, roles uint) error

	// BTMWhitelist
	CreateWhitelist(db *gorm.DB, whitelist *BTMWhitelist) error
	UpdateWhitelistSoftDelete(db *gorm.DB, whitelist *BTMWhitelist) error
	GetWhiteListById(db *gorm.DB, id int64) (data BTMWhitelist, err error)
	GetWhiteListByAddress(db *gorm.DB, address string) (data BTMWhitelist, err error)
	GetWhiteListByCustomerId(db *gorm.DB, customerID uuid.UUID, limit int, page int) (list []BTMWhitelist, total int64, err error)
	GetWhiteLists(db *gorm.DB) (list []BTMWhitelist, err error)
	CheckExistWhitelist(db *gorm.DB, customerID uuid.UUID, cryptoCode string, address string, isUnscoped bool) (bool, bool, error)
	UpdateWhitelist(db *gorm.DB, whitelist *BTMWhitelist) error
	DeleteWhitelist(db *gorm.DB, id int64) error
	SearchWhitelistByAddress(db *gorm.DB, customerID uuid.UUID, address string, limit int, page int) (list []BTMWhitelist, total int64, err error)

	// BTMLoginToken
	IsLastLoginToken(db *gorm.DB, userID uint, loginToken string) (bool, error)
	CreateOrUpdateLastLoginToken(db *gorm.DB, userID uint, loginToken string) error
	DeleteLastLoginToken(db *gorm.DB, userID uint) error

	// BTMLoginLog
	CreateLoginLog(db *gorm.DB, log BTMLoginLog) error
	GetLoginLogs(db *gorm.DB, limit int, page int) ([]BTMLoginLog, int64, error)

	// BTMRole
	InitRawRole(db *gorm.DB) error
	GetRawRoleByRoleName(db *gorm.DB, roleName string) (role BTMRole, err error)
	GetRawRoleById(db *gorm.DB, id int64) (role BTMRole, err error)
	GetRawRoles(db *gorm.DB) (roles []BTMRole, err error)
	CreateRole(db *gorm.DB, role BTMRole) error
	UpdateRole(db *gorm.DB, role BTMRole) error

	// BTMCIB
	UpsertBTMCIB(db *gorm.DB, cib BTMCIB) error
	DeleteBTMCIB(db *gorm.DB, pid string) error
	GetBTMCIBs(db *gorm.DB, id string, limit int, page int) (list []BTMCIB, total int64, err error)
	IsBTMCIBExist(db *gorm.DB, pid string) (bool, int64, error)

	// BTMSumsub
	UpsertBTMSumsub(db *gorm.DB, btmsumsub BTMSumsub) error
	GetBTMSumsub(db *gorm.DB, customerId string) (*BTMSumsub, error)
	UpdateBTMSumsubBanExpireDate(db *gorm.DB, customerId string, banExpireDate int64) error
	DeleteBTMSumsub(db *gorm.DB, customerId string) error
	GetUnCompletedSumsubCustomerIds(db *gorm.DB) ([]string, error)

	// BTMChangeLog
	CreateBTMChangeLog(db *gorm.DB, c BTMChangeLog) error
	UpdateBTMChangeLog(db *gorm.DB, id uint, c BTMChangeLog) error
	BatchCreateBTMChangeLog(db *gorm.DB, c []BTMChangeLog) error
	GetBTMChangeLogs(db *gorm.DB, tableName, customerId string, startAt, endAt time.Time, limit int, page int) (list []BTMChangeLog, total int64, err error)
	AddressExistsInAfterValue(db *gorm.DB, address string) (BTMChangeLog, error)

	// BTMRiskControlLimitSetting
	GetRiskControlCustomerLimitSetting(db *gorm.DB, customerID uuid.UUID) (BTMRiskControlCustomerLimitSetting, error)
	CreateCustomerLimit(db *gorm.DB, customerID uuid.UUID) error
	UpdateCustomerLimit(db *gorm.DB, operationUserId uint, customerID uuid.UUID, newDailyLimit, newMonthlyLimit decimal.Decimal, reason string) error
	UpdateCustomerEdd(db *gorm.DB, operationUserId uint, customerID uuid.UUID, newLevel1, newLevel2 decimal.Decimal, newLevel1Days, newLevel2Days uint32) error
	ChangeCustomerRole(db *gorm.DB, operationUserId uint, customerID uuid.UUID, newRole RiskControlRole, reason string) error
	GetRiskControlRoles() ([]RiskControlRoleKeyValue, error)

	// BTMCustomerNote
	CreateCustomerNote(db *gorm.DB, note BTMCustomerNote) error
	GetCustomerNotes(db *gorm.DB, customerId uuid.UUID, noteType CustomerNoteType, limit int, page int) ([]BTMCustomerNote, int64, error)

	/**
	 * lamassu original
	 **/

	// customers
	GetCustomerById(db *gorm.DB, id uuid.UUID) (*Customer, error)
	SearchCustomers(db *gorm.DB, phone, customerId, address, emailHash, name string,
		whitelistCreatedStartAt, whitelistCreatedEndAt, customerCreatedStartAt, customerCreatedEndAt time.Time,
		customerType CustomerType, active bool,
		limit int, page int) ([]CustomerWithWhiteListCreated, int, error)

	// userConfig
	GetLatestConfData(db *gorm.DB) (UserConfigJSON, error)

	// cashInTx
	GetCashIns(db *gorm.DB, customerID, phone string, startAt, endAt time.Time, limit int, page int) ([]CashInTx, int, error)
}

// type Count struct {
// 	Count int `json:"count"`
// }
