package model

import (
	"gorm.io/gorm"
)

type BTMRole struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	RoleName string
	RoleDesc string
	Role     int64  `gorm:"uniqueIndex:role"`
	RoleRaw  string `gorm:"type:json"`
}
