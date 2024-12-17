package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/pkg/tools"
)

func (repo *repository) InitRawRole() error {
	for _, role := range tools.AllRoles {
		roleName := tools.RoleToString[role]
		item, err := repo.GetRawRoleByRoleName(roleName)
		if item.ID == 0 || err != nil {
			err := repo.db.Create(&domain.BTMRole{
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

func (repo *repository) GetRawRoles() (roles []domain.BTMRole, err error) {
	err = repo.db.Model(&roles).Find(&roles).Error
	return
}

func (repo *repository) GetRawRoleByRoleName(roleName string) (role domain.BTMRole, err error) {
	err = repo.db.Where("role_name = ?", roleName).First(&role).Error
	return
}
