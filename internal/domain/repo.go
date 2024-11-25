package domain

import "BTM-backend/internal/repo/model"

type Repository interface {
	// customers
	GetCustomerByPhone(phone string) (*model.Customer, error)
}
