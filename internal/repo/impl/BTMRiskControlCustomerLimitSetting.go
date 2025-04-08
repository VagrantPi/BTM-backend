package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func (repo *repository) CreateCustomerLimit(db *gorm.DB, customerID uuid.UUID) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	var defaultLimit model.BTMRiskControlLimitSetting
	if err := db.Where("role = ?", domain.RiskControlRoleWhite).First(&defaultLimit).Error; err != nil {
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	afterCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:   customerID,
		Role:         domain.RiskControlRoleWhite,
		DailyLimit:   defaultLimit.DailyLimit,
		MonthlyLimit: defaultLimit.MonthlyLimit,
	}
	afterCustomerLimitJsonData, err := json.Marshal(afterCustomerLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(afterCustomerLimit)").WithCause(err)
	}
	err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
		OperationUserId: domain.OperationUserIdSystem,
		TableName:       domain.BTMChangeLogTableNameBTMRiskControlCustomerLimitSetting,
		OperationType:   domain.BTMChangeLogOperationTypeCreate,
		CustomerId:      &customerID,
		BeforeValue:     nil,
		AfterValue:      afterCustomerLimitJsonData,
	})
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err)
	}

	if err := tx.Create(&model.BTMRiskControlCustomerLimitSetting{
		CustomerId:          customerID,
		Role:                domain.RiskControlRoleWhite.Uint8(), // 預設都為白名單
		DailyLimit:          defaultLimit.DailyLimit,
		MonthlyLimit:        defaultLimit.MonthlyLimit,
		IsCustomized:        false,
		LastBlackToNormalAt: sql.NullTime{},
		LastRole:            domain.RiskControlRoleInit.Uint8(), // 預設都為未設定
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *repository) UpdateCustomerLimit(db *gorm.DB, operationUserId uint, customerID uuid.UUID, newDailyLimit, newMonthlyLimit decimal.Decimal) error {
	var customerLimit model.BTMRiskControlCustomerLimitSetting
	if err := db.Where("customer_id = ?", customerID).First(&customerLimit).Error; err != nil {
		return err
	}

	// 如果為黑名單則不能調整限額
	if customerLimit.Role == domain.RiskControlRoleBlack.Uint8() {
		return errors.BadRequest(error_code.ErrRiskControlRoleIsBlack, "customer is black, cannot update limit")
	}
	if customerLimit.DailyLimit.Equal(newDailyLimit) && customerLimit.MonthlyLimit.Equal(newMonthlyLimit) {
		return errors.BadRequest(error_code.ErrInvalidRequest, "no limit update")
	}
	isUpdateToGray := false
	// 當為白名單，調整限額時，角色會切換成灰名單
	if customerLimit.Role == domain.RiskControlRoleWhite.Uint8() {
		isUpdateToGray = true
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	beforeCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:   customerID,
		Role:         domain.RiskControlRole(customerLimit.Role),
		DailyLimit:   customerLimit.DailyLimit,
		MonthlyLimit: customerLimit.MonthlyLimit,
	}
	beforeCustomerLimitJsonData, err := json.Marshal(beforeCustomerLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(beforeCustomerLimit)").WithCause(err)
	}
	afterCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:   customerID,
		Role:         domain.RiskControlRole(customerLimit.Role), // 固定不變
		DailyLimit:   newDailyLimit,
		MonthlyLimit: newMonthlyLimit,
	}
	// 當為白名單，調整限額時，角色會切換成灰名單
	if isUpdateToGray {
		afterCustomerLimit.Role = domain.RiskControlRoleGray
	}
	afterCustomerLimitJsonData, err := json.Marshal(afterCustomerLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(afterCustomerLimit)").WithCause(err)
	}
	err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
		OperationUserId: operationUserId,
		TableName:       domain.BTMChangeLogTableNameBTMRiskControlCustomerLimitSetting,
		OperationType:   domain.BTMChangeLogOperationTypeUpdate,
		CustomerId:      &customerID,
		BeforeValue:     beforeCustomerLimitJsonData,
		AfterValue:      afterCustomerLimitJsonData,
	})
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err)
	}

	// 更新用戶限額
	if isUpdateToGray {
		customerLimit.Role = domain.RiskControlRoleGray.Uint8()
	}
	customerLimit.DailyLimit = newDailyLimit
	customerLimit.MonthlyLimit = newMonthlyLimit
	customerLimit.IsCustomized = true
	customerLimit.UpdatedAt = time.Now()
	customerLimit.LastRole = beforeCustomerLimit.Role.Uint8()

	if err := tx.Save(&customerLimit).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *repository) ChangeCustomerRole(db *gorm.DB, operationUserId uint, customerID uuid.UUID, newRole domain.RiskControlRole) error {
	var customerLimit model.BTMRiskControlCustomerLimitSetting
	if err := db.Where("customer_id = ?", customerID).First(&customerLimit).Error; err != nil {
		return err
	}

	if customerLimit.Role == newRole.Uint8() {
		return errors.BadRequest(error_code.ErrInvalidRequest, "same role")
	}

	var newDefaultLimit model.BTMRiskControlLimitSetting
	if newRole == domain.RiskControlRoleBlack {
		// 如果設為黑名單，或從黑名單切換回原始，則用戶限額保留原始
		newDefaultLimit.DailyLimit = customerLimit.DailyLimit
		newDefaultLimit.MonthlyLimit = customerLimit.MonthlyLimit
	} else {
		// 取得新的預設限制
		if err := db.Where("role = ?", newRole.Uint8()).First(&newDefaultLimit).Error; err != nil {
			return err
		}
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	beforeCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:   customerID,
		Role:         domain.RiskControlRole(customerLimit.Role),
		DailyLimit:   customerLimit.DailyLimit,
		MonthlyLimit: customerLimit.MonthlyLimit,
	}
	beforeCustomerLimitJsonData, err := json.Marshal(beforeCustomerLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(beforeCustomerLimitChange)").WithCause(err)
	}
	afterCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:   customerID,
		Role:         newRole, // 固定不變
		DailyLimit:   newDefaultLimit.DailyLimit,
		MonthlyLimit: newDefaultLimit.MonthlyLimit,
	}
	afterCustomerLimitJsonData, err := json.Marshal(afterCustomerLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(afterCustomerLimitChange)").WithCause(err)
	}
	err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
		OperationUserId: operationUserId,
		TableName:       domain.BTMChangeLogTableNameBTMRiskControlCustomerLimitSetting,
		OperationType:   domain.BTMChangeLogOperationTypeUpdate,
		CustomerId:      &customerID,
		BeforeValue:     beforeCustomerLimitJsonData,
		AfterValue:      afterCustomerLimitJsonData,
	})
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err)
	}

	// 原始為黑名單切回白或黑需要記錄時間戳
	if customerLimit.Role == domain.RiskControlRoleBlack.Uint8() {
		customerLimit.LastBlackToNormalAt = sql.NullTime{Time: time.Now(), Valid: true}
	}
	customerLimit.Role = newRole.Uint8()
	customerLimit.DailyLimit = newDefaultLimit.DailyLimit
	customerLimit.MonthlyLimit = newDefaultLimit.MonthlyLimit
	customerLimit.IsCustomized = false
	customerLimit.UpdatedAt = time.Now()
	customerLimit.LastRole = beforeCustomerLimit.Role.Uint8() // 紀錄修改前的 role

	if err := tx.Save(&customerLimit).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *repository) GetRiskControlCustomerLimitSetting(db *gorm.DB, customerID uuid.UUID) (domain.BTMRiskControlCustomerLimitSetting, error) {
	if db == nil {
		return domain.BTMRiskControlCustomerLimitSetting{}, errors.BadRequest(error_code.ErrInvalidRequest, "db is nil")
	}

	var customerLimit model.BTMRiskControlCustomerLimitSetting
	err := db.Where("customer_id = ?", customerID).First(&customerLimit).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果用戶不存在，則初始化用戶風控限制
		if err := repo.CreateCustomerLimit(db, customerID); err != nil {
			return domain.BTMRiskControlCustomerLimitSetting{}, errors.InternalServer(error_code.ErrBTMSumsubGetItem, "ChangeCustomerRole not found, init limit err").WithCause(err).
				WithMetadata(map[string]string{
					"customerId": customerID.String(),
				})
		}

		err = db.Where("customer_id = ?", customerID).First(&customerLimit).Error
		return BTMRiskControlCustomerLimitSettingModelToDomain(customerLimit), err
	} else if err != nil {
		err = errors.InternalServer(error_code.ErrDBError, "GetRiskControlCustomerLimitSetting err").WithCause(err)
		return domain.BTMRiskControlCustomerLimitSetting{}, err
	}

	return BTMRiskControlCustomerLimitSettingModelToDomain(customerLimit), nil
}

func (repo *repository) GetRiskControlRoles() ([]domain.RiskControlRoleKeyValue, error) {
	return []domain.RiskControlRoleKeyValue{
		{Id: domain.RiskControlRoleWhite.Uint8(), Name: domain.RiskControlRoleWhite.String()},
		{Id: domain.RiskControlRoleGray.Uint8(), Name: domain.RiskControlRoleGray.String()},
		{Id: domain.RiskControlRoleBlack.Uint8(), Name: domain.RiskControlRoleBlack.String()},
	}, nil
}

func BTMRiskControlCustomerLimitSettingDomainToModel(item domain.BTMRiskControlCustomerLimitSetting) model.BTMRiskControlCustomerLimitSetting {
	return model.BTMRiskControlCustomerLimitSetting{
		Role:         item.Role.Uint8(),
		CustomerId:   item.CustomerId,
		DailyLimit:   item.DailyLimit,
		MonthlyLimit: item.MonthlyLimit,
		IsCustomized: item.IsCustomized,
	}
}

func BTMRiskControlCustomerLimitSettingModelToDomain(item model.BTMRiskControlCustomerLimitSetting) domain.BTMRiskControlCustomerLimitSetting {
	return domain.BTMRiskControlCustomerLimitSetting{
		ID:           item.ID,
		Role:         domain.RiskControlRole(item.Role),
		CustomerId:   item.CustomerId,
		DailyLimit:   item.DailyLimit,
		MonthlyLimit: item.MonthlyLimit,
		IsCustomized: item.IsCustomized,
	}
}
