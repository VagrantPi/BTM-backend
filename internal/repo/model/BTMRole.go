package model

import (
	"gorm.io/gorm"
)

type BTMRole struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	RoleName string `gorm:"uniqueIndex:role_name"`
	RoleDesc string
	RoleRaw  string `gorm:"type:json"`
}
