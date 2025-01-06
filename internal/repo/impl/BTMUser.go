package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *repository) GetBTMUserByAccount(db *gorm.DB, account string) (*domain.BTMUser, error) {
	info := model.BTMUser{}
	if err := db.Model(&model.BTMUser{}).Where("account = ?", account).Find(&info).Error; err != nil {
		return nil, err
	}
	resp := BTMUserModelToDomain(info)
	return &resp, nil
}

func (repo *repository) InitAdmin(db *gorm.DB) error {
	user := domain.BTMUser{
		Account:  "admin",
		Password: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50IjoiYWRtaW4iLCJyYW5kb20iOjE3MzMwMzkxNzg5ODc2NjgsInJvbGUiOjB9.dkefRyMNOKhLYTAd70lQ-QVfpFYrcB3KfsDY3HugLfs",
		Roles:    1,
	}
	item := BTMUserDomainToModel(user)
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account"}},
		DoNothing: true,
	}).Create(&item).Error
}

func BTMUserDomainToModel(user domain.BTMUser) model.BTMUser {
	return model.BTMUser{
		Account:  user.Account,
		Password: user.Password,
		Roles:    user.Roles,
	}
}

func BTMUserModelToDomain(user model.BTMUser) domain.BTMUser {
	return domain.BTMUser{
		Id:       user.ID,
		Account:  user.Account,
		Password: user.Password,
		Roles:    user.Roles,
	}
}
