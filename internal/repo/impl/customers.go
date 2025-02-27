package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (repo *repository) GetCustomers(db *gorm.DB, limit int, page int) ([]domain.CustomerWithWhiteListCreated, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []domain.CustomerWithWhiteListCreated{}

	sql := db.Model(&model.Customer{}).
		Select(
			"DISTINCT ON (customers.id) customers.*",
			"btm_whitelists.created_at AS first_white_list_created",
		).
		Joins("LEFT JOIN btm_whitelists ON btm_whitelists.customer_id = customers.id").
		Where("customers.phone != '' AND btm_whitelists.deleted_at ISNULL").
		Order("customers.id, btm_whitelists.created_at ASC")
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
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

func (repo *repository) SearchCustomersByCustomerId(db *gorm.DB, customerId string, limit int, page int) ([]domain.CustomerWithWhiteListCreated, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []domain.CustomerWithWhiteListCreated{}

	sql := db.Model(&model.Customer{}).
		Select(
			"DISTINCT ON (customers.id) customers.*",
			"btm_whitelists.created_at AS first_white_list_created",
		).
		Joins("LEFT JOIN btm_whitelists ON btm_whitelists.customer_id = customers.id").
		Where("customers.id::TEXT LIKE ? AND btm_whitelists.deleted_at ISNULL", "%"+customerId+"%").
		Order("customers.id, btm_whitelists.created_at ASC")
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, int(total), nil
}

func (repo *repository) SearchCustomersByPhone(db *gorm.DB, phone string, limit int, page int) ([]domain.CustomerWithWhiteListCreated, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []domain.CustomerWithWhiteListCreated{}

	sql := db.Model(&model.Customer{}).
		Select(
			"DISTINCT ON (customers.id) customers.*",
			"btm_whitelists.created_at AS first_white_list_created",
		).
		Joins("LEFT JOIN btm_whitelists ON btm_whitelists.customer_id = customers.id").
		Where("customers.phone LIKE ? AND btm_whitelists.deleted_at ISNULL", "%"+phone+"%").
		Order("customers.id, btm_whitelists.created_at ASC")
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, int(total), nil
}

func (repo *repository) SearchCustomersByAddress(db *gorm.DB, address string, limit int, page int) ([]domain.CustomerWithWhiteListCreated, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []domain.CustomerWithWhiteListCreated{}

	sql := db.Model(&model.Customer{}).
		Select(
			"DISTINCT ON (customers.id) customers.*",
			"btm_whitelists.created_at AS first_white_list_created",
		).
		Joins("INNER JOIN btm_whitelists ON customers.id = btm_whitelists.customer_id").
		Where("btm_whitelists.address = ? AND btm_whitelists.deleted_at ISNULL", address).
		Order("customers.id, btm_whitelists.created_at ASC")
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, int(total), nil
}

func (repo *repository) SearchCustomersByWhitelistCreatedAt(db *gorm.DB, startAt, endAt time.Time, limit int, page int) ([]domain.CustomerWithWhiteListCreated, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []domain.CustomerWithWhiteListCreated{}

	sql := db.Model(&model.Customer{}).
		Select(
			"DISTINCT ON (customers.id) customers.*",
			"btm_whitelists.created_at AS first_white_list_created",
		).
		Joins("INNER JOIN btm_whitelists ON customers.id = btm_whitelists.customer_id").
		Where("btm_whitelists.created_at BETWEEN ? AND ? AND btm_whitelists.deleted_at ISNULL", startAt, endAt).
		Order("customers.id, btm_whitelists.created_at ASC")
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, int(total), nil
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
