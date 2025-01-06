package model

import "gorm.io/gorm"

type BTMUser struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	Account  string `gorm:"uniqueIndex:user_account"`
	Password string
	Roles    int64
}
