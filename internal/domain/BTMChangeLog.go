package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type BTMChangeLog struct {
	ID              uint                      `json:"id"`
	OperationUserId int64                     `json:"operation_user_id"`
	TableName       BTMChangeLogTableName     `json:"table_name"`
	OperationType   BTMChangeLogOperationType `json:"operation_type"`
	CustomerId      *uuid.UUID                `json:"customer_id"`
	BeforeValue     json.RawMessage           `json:"before_value"`
	AfterValue      json.RawMessage           `json:"after_value"`
	CreatedAt       time.Time                 `json:"created_at"`
}

const OperationUserIdSystem int64 = 0

type BTMChangeLogOperationType uint8

const (
	BTMChangeLogOperationTypeCreate BTMChangeLogOperationType = iota + 1
	BTMChangeLogOperationTypeUpdate
	BTMChangeLogOperationTypeDelete
)

type BTMChangeLogTableName string

const (
	BTMChangeLogTableNameBTMUsers                           BTMChangeLogTableName = "btm_users"
	BTMChangeLogTableNameBTMWhitelist                       BTMChangeLogTableName = "btm_whitelists"
	BTMChangeLogTableNameBTMRoles                           BTMChangeLogTableName = "btm_roles"
	BTMChangeLogTableNameBTMRiskControlCustomerLimitSetting BTMChangeLogTableName = "btm_risk_control_customer_limit_settings"
	BTMChangeLogTableNameBTMRiskControlLimitSetting         BTMChangeLogTableName = "btm_risk_control_limit_settings"
	BTMChangeLogTableNameBTMCustomerNotes                   BTMChangeLogTableName = "btm_customer_notes"
)
