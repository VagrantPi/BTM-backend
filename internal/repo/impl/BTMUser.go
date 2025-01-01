package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"

	"gorm.io/gorm"
)

func (repo *repository) GetBTMUserByAccount(db *gorm.DB, account string) (*domain.BTMUser, error) {
	info := model.BTMUser{}
	if err := db.Model(&model.BTMUser{}).Where("account = ?", account).Find(&info).Error; err != nil {
		return nil, err
	}
	resp := BTMUserModelToDomain(info)
	return &resp, nil
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
