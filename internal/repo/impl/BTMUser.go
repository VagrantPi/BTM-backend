package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *repository) GetBTMUserByAccount(db *gorm.DB, account string) (*domain.BTMUser, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	info := model.BTMUser{}
	if err := db.Model(&model.BTMUser{}).Where("account = ?", account).Find(&info).Error; err != nil {
		return nil, err
	}
	resp := BTMUserModelToDomain(info)
	return &resp, nil
}

func (repo *repository) InitAdmin(db *gorm.DB) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	user := domain.BTMUser{
		Account:  "admin",
		Password: "$2a$12$qzSY/1.YLuZ1FnYv4q8rlehgBA6nX/CQ9MDDwjoQeJJvDoUzkfVVO",
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
