package model

import "gorm.io/gorm"

type BTMLoginToken struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	UserID     uint   `gorm:"uniqueIndex:user_token,priority:1; uniqueIndex:user_unique"`
	LoginToken string `gorm:"uniqueIndex:user_token,priority:2"`
}
