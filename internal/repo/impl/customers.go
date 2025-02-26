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

func (repo *repository) GetCustomers(db *gorm.DB, limit int, page int) ([]domain.Customer, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []model.Customer{}

	sql := db.Model(&model.Customer{}).Where("phone != ''")
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.Customer, 0, len(list))
	for _, customer := range list {
		resp = append(resp, CustomerModelToDomain(customer))
	}
	return resp, int(total), nil
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

func (repo *repository) SearchCustomersByCustomerId(db *gorm.DB, customerId string, limit int, page int) ([]domain.Customer, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []model.Customer{}

	sql := db.Model(&model.Customer{}).Where("id::TEXT LIKE ?", "%"+customerId+"%")
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.Customer, 0, len(list))
	for _, customer := range list {
		resp = append(resp, CustomerModelToDomain(customer))
	}
	return resp, int(total), nil
}

func (repo *repository) SearchCustomersByPhone(db *gorm.DB, phone string, limit int, page int) ([]domain.Customer, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []model.Customer{}

	sql := db.Model(&model.Customer{}).Where("phone LIKE ?", "%"+phone+"%")
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.Customer, 0, len(list))
	for _, customer := range list {
		resp = append(resp, CustomerModelToDomain(customer))
	}
	return resp, int(total), nil
}

func (repo *repository) SearchCustomersByAddress(db *gorm.DB, address string, limit int, page int) ([]domain.Customer, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []model.Customer{}

	sql := db.Model(&model.Customer{}).
		Select("DISTINCT ON (customers.id) customers.*").
		Joins("INNER JOIN btm_whitelists ON customers.id = btm_whitelists.customer_id").
		Where("btm_whitelists.address = ? AND btm_whitelists.deleted_at ISNULL", address)
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.Customer, 0, len(list))
	for _, customer := range list {
		resp = append(resp, CustomerModelToDomain(customer))
	}
	return resp, int(total), nil
}

func (repo *repository) SearchCustomersByWhitelistCreatedAt(db *gorm.DB, startAt, endAt time.Time, limit int, page int) ([]domain.Customer, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []model.Customer{}

	sql := db.Model(&model.Customer{}).
		Select("DISTINCT ON (customers.id) customers.*").
		Joins("INNER JOIN btm_whitelists ON customers.id = btm_whitelists.customer_id").
		Where("btm_whitelists.created_at BETWEEN ? AND ? AND btm_whitelists.deleted_at ISNULL", startAt, endAt)
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.Customer, 0, len(list))
	for _, customer := range list {
		resp = append(resp, CustomerModelToDomain(customer))
	}
	return resp, int(total), nil
}

func (repo *repository) SearchCustomersByTxCreatedAt(db *gorm.DB, startAt, endAt time.Time, limit int, page int) ([]domain.Customer, int, error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list := []model.Customer{}

	sql := db.Model(&model.Customer{}).
		Select("DISTINCT ON (customers.id) customers.*").
		Joins("INNER JOIN cash_in_txs ON cash_in_txs.customer_id = customers.id").
		Where("cash_in_txs.created BETWEEN ? AND ?", startAt, endAt)
	var total int64 = 0
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	resp := make([]domain.Customer, 0, len(list))
	for _, customer := range list {
		resp = append(resp, CustomerModelToDomain(customer))
	}
	return resp, int(total), nil
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
