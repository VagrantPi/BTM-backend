package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BTMChangeLog struct {
	Db *gorm.DB `gorm:"-"`
	gorm.Model

	OperationUserId uint
	TableName       string     `gorm:"index:operation_table,priority:1; index:operation_table_customer,priority:1"`
	OperationType   uint8      `gorm:"index:operation_table,priority:2"`
	CustomerId      *uuid.UUID `gorm:"index:operation_table_customer,priority:2"`
	BeforeValue     json.RawMessage
	AfterValue      json.RawMessage
}
