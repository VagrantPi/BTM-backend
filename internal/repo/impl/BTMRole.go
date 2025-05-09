package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (repo *repository) InitRawRole(db *gorm.DB) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	roleName := domain.RoleAdminName
	item, err := repo.GetRawRoleByRoleName(db, roleName)
	if item.ID == 0 || err != nil {
		// 不存在則建立
		err := db.Create(&model.BTMRole{
			RoleName: domain.RoleAdminName,
			RoleDesc: domain.RoleAdminName,
			RoleRaw:  domain.DefaultRoleRaw,
		}).Error
		if err != nil {
			return err
		}
	} else {
		// 更新最新 DefaultRoleRaw
		err := db.Model(&model.BTMRole{}).Where("role_name = ?", roleName).
			Updates(map[string]interface{}{
				"role_raw":   domain.DefaultRoleRaw,
				"updated_at": time.Now(),
			}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *repository) GetRawRoles(db *gorm.DB) (roles []domain.BTMRole, err error) {
	if db == nil {
		return nil, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	err = db.Model(&roles).Order("id ASC").Find(&roles).Error
	return
}

func (repo *repository) GetRawRoleByRoleName(db *gorm.DB, roleName string) (role domain.BTMRole, err error) {
	if db == nil {
		return domain.BTMRole{}, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	err = db.Where("role_name = ?", roleName).First(&role).Error
	return
}

func (repo *repository) GetRawRoleById(db *gorm.DB, id int64) (role domain.BTMRole, err error) {
	if db == nil {
		return domain.BTMRole{}, errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	err = db.Where("id = ?", id).First(&role).Error
	return
}

func (repo *repository) CreateRole(db *gorm.DB, role domain.BTMRole) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	item := BTMRoleDomainToModel(role)
	return db.Model(&model.BTMRole{}).Create(&item).Error
}

func (repo *repository) UpdateRole(db *gorm.DB, role domain.BTMRole) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	return db.Model(&model.BTMRole{}).Where("role_name = ?", role.RoleName).
		Updates(map[string]interface{}{
			"role_desc":  role.RoleDesc,
			"role_raw":   role.RoleRaw,
			"updated_at": time.Now(),
		}).Error
}

func BTMRoleDomainToModel(role domain.BTMRole) model.BTMRole {
	return model.BTMRole{
		RoleName: role.RoleName,
		RoleDesc: role.RoleDesc,
		RoleRaw:  role.RoleRaw,
	}
}

func BTMRoleModelToDomain(role model.BTMRole) domain.BTMRole {
	return domain.BTMRole{
		ID:       role.ID,
		RoleName: role.RoleName,
		RoleDesc: role.RoleDesc,
		RoleRaw:  role.RoleRaw,
	}
}
