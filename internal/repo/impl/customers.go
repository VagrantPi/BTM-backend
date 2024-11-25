package impl

import (
	"BTM-backend/internal/repo/model"
)

func (repo *repository) GetCustomerByPhone(phone string) (*model.Customer, error) {
	info := model.Customer{}
	if err := repo.db.Where("phone = ?", phone).Find(&info).Error; err != nil {
		return nil, err
	}
	return &info, nil
}
