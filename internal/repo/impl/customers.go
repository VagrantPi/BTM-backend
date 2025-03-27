package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func _searchCustomers(db *gorm.DB,
	phone, customerId, address, emailHash, name string,
	whitelistCreatedStartAt, whitelistCreatedEndAt, customerCreatedStartAt, customerCreatedEndAt time.Time,
	customerType domain.CustomerType, limit, page int) (*gorm.DB, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	sql := db.Model(&model.Customer{}).
		Select(
			"DISTINCT ON (customers.id) customers.id",
			"customers.phone",
			"btm_sumsubs.name",
			"btm_sumsubs.email_hash",
			"btm_sumsubs.info_hash",
			"customers.created",
			"btm_whitelists.created_at AS first_white_list_created",
		).
		Joins("INNER JOIN btm_whitelists ON customers.id = btm_whitelists.customer_id").
		Joins("LEFT JOIN btm_sumsubs ON btm_sumsubs.customer_id = customers.id::text")

	switch {
	case !whitelistCreatedStartAt.IsZero() && !whitelistCreatedEndAt.IsZero():
		sql = sql.Where("btm_whitelists.created_at BETWEEN ? AND ? AND btm_whitelists.deleted_at ISNULL", whitelistCreatedStartAt, whitelistCreatedEndAt)
	case !customerCreatedStartAt.IsZero() && !customerCreatedEndAt.IsZero():
		sql = sql.Where("customers.created BETWEEN ? AND ?", customerCreatedStartAt, customerCreatedEndAt)
	case strings.TrimSpace(address) != "":
		sql = sql.Where("btm_whitelists.address = ?", strings.TrimSpace(address))
	case strings.TrimSpace(phone) != "":
		sql = sql.Where("customers.phone LIKE ?", "%"+strings.TrimSpace(phone)+"%")
	case strings.TrimSpace(customerId) != "":
		sql = sql.Where("customers.id::TEXT LIKE ?", "%"+strings.TrimSpace(customerId)+"%")
	case strings.TrimSpace(name) != "":
		sql = sql.Where("btm_sumsubs.name = ?", strings.TrimSpace(name))
	case strings.TrimSpace(emailHash) != "":
		sql = sql.Where("btm_sumsubs.email_hash = ?", strings.TrimSpace(emailHash))
	default:
		sql = sql.Where("customers.phone != ''")
	}

	// 取得現在的中華民國年日期
	now := time.Now()
	twYear := now.Year() - 1911
	today := fmt.Sprintf("%03d%02d%02d", twYear, int(now.Month()), now.Day())
	switch customerType {
	case domain.CustomerTypeWhiteList:
		sql = sql.
			Joins("LEFT JOIN btm_risk_control_customer_limit_settings ON customers.id::TEXT = btm_risk_control_customer_limit_settings.customer_id").
			Where("btm_sumsubs.ban_expire_date IS NULL OR btm_sumsubs.ban_expire_date > ?", today).
			Where("btm_risk_control_customer_limit_settings.role IS NULL OR btm_risk_control_customer_limit_settings.role = ?", domain.RiskControlRoleWhite)
	case domain.CustomerTypeGrayList:
		sql = sql.
			Joins("LEFT JOIN btm_risk_control_customer_limit_settings ON customers.id::TEXT = btm_risk_control_customer_limit_settings.customer_id").
			Where("btm_sumsubs.ban_expire_date IS NULL OR btm_sumsubs.ban_expire_date > ?", today).
			Where("btm_risk_control_customer_limit_settings.role = ?", domain.RiskControlRoleGray)
	case domain.CustomerTypeBlackList:
		sql = sql.Select(
			"DISTINCT ON (customers.id) customers.id",
			"customers.phone",
			"customers.created",
			"btm_whitelists.created_at AS first_white_list_created",
			"customers.authorized_override = 'blocked' AS is_lamassu_block",
			"btm_risk_control_customer_limit_settings.role = 3 AS is_admin_block",
			"UPPER(TRIM(btm_sumsubs.id_number)) = UPPER(TRIM(btm_cibs.pid)) AS is_cib_block",
		).
			Joins("LEFT JOIN btm_sumsubs ON customers.id::TEXT = btm_sumsubs.customer_id").
			Joins("LEFT JOIN btm_risk_control_customer_limit_settings ON customers.id::TEXT = btm_risk_control_customer_limit_settings.customer_id").
			Joins("LEFT JOIN btm_cibs ON btm_sumsubs.id_number = btm_cibs.pid").
			Where("(btm_sumsubs.ban_expire_date IS NOT NULL AND btm_sumsubs.ban_expire_date < ?) OR btm_risk_control_customer_limit_settings.role = ? OR customers.authorized_override = 'blocked' OR UPPER(TRIM(btm_sumsubs.id_number)) = UPPER(TRIM(btm_cibs.pid))", today, domain.RiskControlRoleBlack)
	}

	return sql, nil
}

func (repo *repository) SearchCustomers(db *gorm.DB,
	phone, customerId, address, emailHash, name string,
	whitelistCreatedStartAt, whitelistCreatedEndAt, customerCreatedStartAt, customerCreatedEndAt time.Time,
	customerType domain.CustomerType,
	limit, page int) ([]domain.CustomerWithWhiteListCreated, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []domain.CustomerWithWhiteListCreated{}

	sql, err := _searchCustomers(db, phone, customerId, address, emailHash, name,
		whitelistCreatedStartAt, whitelistCreatedEndAt, customerCreatedStartAt, customerCreatedEndAt,
		customerType, limit, page)
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := sql.Distinct("customers.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sql, err = _searchCustomers(db, phone, customerId, address, emailHash, name,
		whitelistCreatedStartAt, whitelistCreatedEndAt, customerCreatedStartAt, customerCreatedEndAt,
		customerType, limit, page)
	if err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).
		Offset(offset).
		Order("customers.id ASC").
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, int(total), nil
}

func (repo *repository) GetCustomerById(db *gorm.DB, id uuid.UUID) (*domain.Customer, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	modelCustomer := model.Customer{}
	if err := db.Where("id = ?", id).First(&modelCustomer).Error; err != nil {
		return nil, err
	}
	customer := CustomerModelToDomain(modelCustomer)
	return &customer, nil
}

func CustomerModelToDomain(customer model.Customer) domain.Customer {
	return domain.Customer{
		ID:      customer.ID,
		Phone:   customer.Phone,
		Created: customer.Created,
	}
}

func CustomerDomainToModel(customer domain.Customer) model.Customer {
	return model.Customer{
		ID:      customer.ID,
		Phone:   customer.Phone,
		Created: customer.Created,
	}
}
