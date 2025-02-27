package domain

import "encoding/json"

type BTMChangeLog struct {
	OperationUserId uint
	TableName       BTMChangeLogTableName
	OperationType   BTMChangeLogOperationType
	BeforeValue     json.RawMessage
	AfterValue      json.RawMessage
}

type BTMChangeLogOperationType uint8

const (
	BTMChangeLogOperationTypeCreate BTMChangeLogOperationType = iota + 1
	BTMChangeLogOperationTypeUpdate
	BTMChangeLogOperationTypeDelete
)

type BTMChangeLogTableName string

const (
	BTMChangeLogTableNameBTMWhitelist BTMChangeLogTableName = "btm_whitelists"
)
