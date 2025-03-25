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

func (repo *repository) GetBTMUserById(db *gorm.DB, id uint) (*domain.BTMUser, error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	info := model.BTMUser{}
	if err := db.Model(&model.BTMUser{}).Where("id = ?", id).Find(&info).Error; err != nil {
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

func (repo *repository) CreateBTMUser(db *gorm.DB, user domain.BTMUser) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	item := BTMUserDomainToModel(user)
	return db.Create(&item).Error
}

func (repo *repository) GetUsers(db *gorm.DB, limit int, page int) (list []domain.BTMUserWithRoles, total int64, err error) {
	if db == nil {
		return nil, 0, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	offset := (page - 1) * limit
	list = []domain.BTMUserWithRoles{}

	sql := db.Model(&model.BTMUser{})

	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := sql.Select(
		"btm_users.id",
		"btm_users.account",
		"btm_roles.id as role_id",
		"btm_roles.role_name",
	).
		Joins("LEFT JOIN btm_roles ON btm_users.roles = btm_roles.id").
		Offset(offset).
		Limit(limit).
		Order("id").
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (repo *repository) UpdateUserNameRoles(db *gorm.DB, id uint, account string, roles uint) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	return db.Model(&model.BTMUser{}).Where("id = ?", id).Updates(map[string]interface{}{
		"account": account,
		"roles":   roles,
	}).Error
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
