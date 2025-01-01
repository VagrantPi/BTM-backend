package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/tools"

	"gorm.io/gorm"
)

func (repo *repository) InitRawRole(db *gorm.DB) error {
	for _, role := range tools.AllRoles {
		roleName := tools.RoleToString[role]
		item, err := repo.GetRawRoleByRoleName(db, roleName)
		if item.ID == 0 || err != nil {
			err := db.Create(&domain.BTMRole{
				RoleName: roleName,
				RoleDesc: roleName,
				Role:     int64(role),
				RoleRaw:  domain.DefaultRoleRaw,
			}).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *repository) GetRawRoles(db *gorm.DB) (roles []domain.BTMRole, err error) {
	err = db.Model(&roles).Find(&roles).Error
	return
}

func (repo *repository) GetRawRoleByRoleName(db *gorm.DB, roleName string) (role domain.BTMRole, err error) {
	err = db.Where("role_name = ?", roleName).First(&role).Error
	return
}
