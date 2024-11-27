package model

import "gorm.io/gorm"

type BTMUser struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	Account  string
	Password string
	Roles    int64
}
