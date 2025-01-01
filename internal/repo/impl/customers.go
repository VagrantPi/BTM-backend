package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (repo *repository) GetCustomers(db *gorm.DB, limit int, page int) ([]domain.Customer, int, error) {
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
	modelCustomer := model.Customer{}
	if err := db.Where("id = ?", id).First(&modelCustomer).Error; err != nil {
		return nil, err
	}
	customer := CustomerModelToDomain(modelCustomer)
	return &customer, nil
}

func (repo *repository) SearchCustomersByPhone(db *gorm.DB, phone string, limit int, page int) ([]domain.Customer, int, error) {
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

func CustomerModelToDomain(customer model.Customer) domain.Customer {
	return domain.Customer{
		ID:    customer.ID,
		Phone: customer.Phone,
	}
}

func CustomerDomainToModel(customer domain.Customer) model.Customer {
	return model.Customer{
		ID:    customer.ID,
		Phone: customer.Phone,
	}
}
