package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type BTMChangeLog struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	OperationUserId uint
	TableName       string `gorm:"index:operation_table,priority:1"`
	OperationType   uint8  `gorm:"index:operation_table,priority:2"`
	BeforeValue     json.RawMessage
	AfterValue      json.RawMessage
}
