package impl

import (
	"BTM-backend/configs"
	"BTM-backend/internal/domain"

	"gorm.io/gorm"
)

type repository struct {
	db      *gorm.DB
	configs configs.Config
}

func NewRepository(db *gorm.DB, configs configs.Config) domain.Repository {
	return &repository{
		db:      db,
		configs: configs,
	}
}
