package impl

import (
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/error_code"
	"database/sql"
	"encoding/json"
	"fmt"
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
		Level1:       defaultLimit.Level1,
		Level2:       defaultLimit.Level2,
		Level1Days:   defaultLimit.Level1Days,
		Level2Days:   defaultLimit.Level2Days,
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
	fmt.Println("defaultLimit.Level1Days:", defaultLimit.Level1Days)

	if err := tx.Create(&model.BTMRiskControlCustomerLimitSetting{
		CustomerId:          customerID,
		Role:                domain.RiskControlRoleWhite.Uint8(), // 預設都為白名單
		DailyLimit:          defaultLimit.DailyLimit,
		MonthlyLimit:        defaultLimit.MonthlyLimit,
		Level1:              defaultLimit.Level1,
		Level2:              defaultLimit.Level2,
		Level1Days:          defaultLimit.Level1Days,
		Level2Days:          defaultLimit.Level2Days,
		IsCustomized:        false,
		LastBlackToNormalAt: sql.NullTime{},
		LastRole:            domain.RiskControlRoleInit.Uint8(), // 預設都為未設定
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *repository) UpdateCustomerLimit(db *gorm.DB, operationUserId int64, customerID uuid.UUID, newDailyLimit, newMonthlyLimit decimal.Decimal, reason string) error {
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

	if reason == "" {
		return errors.BadRequest(error_code.ErrInvalidRequest, "reason is empty")
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	afterChangeRoleReason := customerLimit.ChangeRoleReason
	if isUpdateToGray {
		afterChangeRoleReason = "系統自動切換"
	}
	beforeCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:        customerID,
		Role:              domain.RiskControlRole(customerLimit.Role),
		DailyLimit:        customerLimit.DailyLimit,
		MonthlyLimit:      customerLimit.MonthlyLimit,
		Level1:            customerLimit.Level1,
		Level2:            customerLimit.Level2,
		ChangeRoleReason:  customerLimit.ChangeRoleReason,
		ChangeLimitReason: customerLimit.ChangeLimitReason,
	}
	beforeCustomerLimitJsonData, err := json.Marshal(beforeCustomerLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(beforeCustomerLimit)").WithCause(err)
	}
	afterCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:        customerID,
		Role:              domain.RiskControlRole(customerLimit.Role), // 固定不變
		DailyLimit:        newDailyLimit,
		MonthlyLimit:      newMonthlyLimit,
		Level1:            customerLimit.Level1,
		Level2:            customerLimit.Level2,
		ChangeRoleReason:  afterChangeRoleReason, // 如果是白名單調整限額，則會更改角色，並帶入 系統自動切換 原因，否則該欄位固定不變
		ChangeLimitReason: reason,
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

		// 取得灰名單的等級門檻
		var grayDefaultLimit model.BTMRiskControlLimitSetting
		if err := db.Where("role = ?", domain.RiskControlRoleGray).First(&grayDefaultLimit).Error; err != nil {
			return errors.InternalServer(error_code.ErrDBError, "GetRiskControlCustomerLimitSetting err").WithCause(err)
		}
		customerLimit.Level1 = grayDefaultLimit.Level1
		customerLimit.Level2 = grayDefaultLimit.Level2
	}
	customerLimit.DailyLimit = newDailyLimit
	customerLimit.MonthlyLimit = newMonthlyLimit
	customerLimit.IsCustomized = true
	customerLimit.UpdatedAt = time.Now()
	customerLimit.LastRole = beforeCustomerLimit.Role.Uint8()
	customerLimit.ChangeLimitReason = reason
	customerLimit.ChangeRoleReason = "X"

	if err := tx.Save(&customerLimit).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *repository) ResetCustomerRole(db *gorm.DB, operationUserId int64, customerID uuid.UUID) error {
	var customerLimit model.BTMRiskControlCustomerLimitSetting
	if err := db.Where("customer_id = ?", customerID).First(&customerLimit).Error; err != nil {
		return err
	}

	if customerLimit.Role != domain.RiskControlRoleBlack.Uint8() {
		return errors.BadRequest(error_code.ErrInvalidRequest, "customer role need black")
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	beforeCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:        customerID,
		Role:              domain.RiskControlRole(customerLimit.Role),
		DailyLimit:        customerLimit.DailyLimit,
		MonthlyLimit:      customerLimit.MonthlyLimit,
		Level1:            customerLimit.Level1,
		Level2:            customerLimit.Level2,
		Level1Days:        customerLimit.Level1Days,
		Level2Days:        customerLimit.Level2Days,
		ChangeRoleReason:  customerLimit.ChangeRoleReason,
		ChangeLimitReason: customerLimit.ChangeLimitReason,
	}
	beforeCustomerLimitJsonData, err := json.Marshal(beforeCustomerLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(beforeCustomerLimitChange)").WithCause(err)
	}
	afterCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:        customerID,
		Role:              domain.RiskControlRole(customerLimit.LastRole),
		DailyLimit:        customerLimit.DailyLimit,
		MonthlyLimit:      customerLimit.MonthlyLimit,
		Level1:            customerLimit.Level1,
		Level2:            customerLimit.Level2,
		Level1Days:        customerLimit.Level1Days,
		Level2Days:        customerLimit.Level2Days,
		ChangeRoleReason:  "系統回復原始角色",
		ChangeLimitReason: "",
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

	// customers.authorized_override 切換成 verified
	err = repo.UpdateCustomerAuthorizedOverride(tx, customerID, domain.CustomerAuthorizedOverrideVerified)
	if err != nil {
		return errors.InternalServer(error_code.ErrUserUpdate, "UpdateCustomerAuthorizedOverride").WithCause(err)
	}

	customerLimit.Role = customerLimit.LastRole // 返回原本角色
	customerLimit.UpdatedAt = time.Now()
	customerLimit.LastRole = customerLimit.Role
	customerLimit.ChangeRoleReason = "系統回復原始角色"
	customerLimit.ChangeLimitReason = "X"
	customerLimit.LastBlackToNormalAt = sql.NullTime{Time: time.Now(), Valid: true}
	customerLimit.IsEdd = false

	if err := tx.Save(&customerLimit).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *repository) ChangeCustomerRole(db *gorm.DB, operationUserId int64, customerID uuid.UUID, newRole domain.RiskControlRole, reason string) error {
	var customerLimit model.BTMRiskControlCustomerLimitSetting
	if err := db.Where("customer_id = ?", customerID).First(&customerLimit).Error; err != nil {
		return err
	}

	if customerLimit.Role == newRole.Uint8() {
		return errors.BadRequest(error_code.ErrInvalidRequest, "same role")
	}

	if reason == "" {
		return errors.BadRequest(error_code.ErrInvalidRequest, "reason is empty")
	}

	var newDefaultLimit model.BTMRiskControlLimitSetting
	if newRole == domain.RiskControlRoleBlack {
		// 如果設為黑名單，或從黑名單切換回原始，則用戶限額保留原始
		newDefaultLimit.DailyLimit = customerLimit.DailyLimit
		newDefaultLimit.MonthlyLimit = customerLimit.MonthlyLimit
		newDefaultLimit.Level1 = customerLimit.Level1
		newDefaultLimit.Level2 = customerLimit.Level2
		newDefaultLimit.Level1Days = customerLimit.Level1Days
		newDefaultLimit.Level2Days = customerLimit.Level2Days

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
		CustomerId:        customerID,
		Role:              domain.RiskControlRole(customerLimit.Role),
		DailyLimit:        customerLimit.DailyLimit,
		MonthlyLimit:      customerLimit.MonthlyLimit,
		Level1:            customerLimit.Level1,
		Level2:            customerLimit.Level2,
		Level1Days:        customerLimit.Level1Days,
		Level2Days:        customerLimit.Level2Days,
		ChangeRoleReason:  customerLimit.ChangeRoleReason,
		ChangeLimitReason: customerLimit.ChangeLimitReason,
	}
	beforeCustomerLimitJsonData, err := json.Marshal(beforeCustomerLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(beforeCustomerLimitChange)").WithCause(err)
	}
	afterCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:        customerID,
		Role:              newRole, // 固定不變
		DailyLimit:        newDefaultLimit.DailyLimit,
		MonthlyLimit:      newDefaultLimit.MonthlyLimit,
		Level1:            newDefaultLimit.Level1,
		Level2:            newDefaultLimit.Level2,
		Level1Days:        newDefaultLimit.Level1Days,
		Level2Days:        newDefaultLimit.Level2Days,
		ChangeRoleReason:  reason,
		ChangeLimitReason: "",
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

	if customerLimit.Role == domain.RiskControlRoleBlack.Uint8() {
		// 原始為黑名單切回白或黑需要記錄時間戳
		customerLimit.LastBlackToNormalAt = sql.NullTime{Time: time.Now(), Valid: true}
		customerLimit.EddType = ""

		// customers.authorized_override 切換回 verified
		err = repo.UpdateCustomerAuthorizedOverride(tx, customerID, domain.CustomerAuthorizedOverrideVerified)
		if err != nil {
			return errors.InternalServer(error_code.ErrUserUpdate, "UpdateCustomerAuthorizedOverride").WithCause(err)
		}
	}

	if newRole == domain.RiskControlRoleBlack {
		// customers.authorized_override 切換成 blocked
		err = repo.UpdateCustomerAuthorizedOverride(tx, customerID, domain.CustomerAuthorizedOverrideBlocked)
		if err != nil {
			return errors.InternalServer(error_code.ErrUserUpdate, "UpdateCustomerAuthorizedOverride").WithCause(err)
		}
	}

	customerLimit.Role = newRole.Uint8()
	customerLimit.DailyLimit = newDefaultLimit.DailyLimit
	customerLimit.MonthlyLimit = newDefaultLimit.MonthlyLimit
	customerLimit.Level1 = newDefaultLimit.Level1
	customerLimit.Level2 = newDefaultLimit.Level2
	customerLimit.Level1Days = newDefaultLimit.Level1Days
	customerLimit.Level2Days = newDefaultLimit.Level2Days
	customerLimit.IsCustomized = false
	customerLimit.UpdatedAt = time.Now()
	customerLimit.LastRole = beforeCustomerLimit.Role.Uint8() // 紀錄修改前的 role
	customerLimit.ChangeRoleReason = reason
	customerLimit.ChangeLimitReason = "X"

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

	// edd level1, level2 2025-04-10 才加入，因此如果沒有則更新預設
	if customerLimit.Level1.IsZero() || customerLimit.Level2.IsZero() {
		var defaultLimit model.BTMRiskControlLimitSetting
		// 取得新的預設限制
		if err := db.Where("role = ?", customerLimit.Role).First(&defaultLimit).Error; err != nil {
			return domain.BTMRiskControlCustomerLimitSetting{}, errors.InternalServer(error_code.ErrDBError, "GetRiskControlCustomerLimitSetting err").WithCause(err)
		}
		customerLimit.Level1 = defaultLimit.Level1
		customerLimit.Level2 = defaultLimit.Level2

		if err := db.Save(&customerLimit).Error; err != nil {
			return domain.BTMRiskControlCustomerLimitSetting{}, errors.InternalServer(error_code.ErrDBError, "GetRiskControlCustomerLimitSetting err").WithCause(err)
		}
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

func (repo *repository) UpdateCustomerEddSetting(db *gorm.DB, operationUserId int64, customerID uuid.UUID, newLevel1, newLevel2 decimal.Decimal, newLevel1Days, newLevel2Days uint32) error {
	var customerLimit model.BTMRiskControlCustomerLimitSetting
	if err := db.Where("customer_id = ?", customerID).First(&customerLimit).Error; err != nil {
		return err
	}

	// 如果為黑名單則不能調整限額
	if customerLimit.Role == domain.RiskControlRoleBlack.Uint8() {
		return errors.BadRequest(error_code.ErrRiskControlRoleIsBlack, "customer is black, cannot update edd")
	}
	if customerLimit.Level1.Equal(newLevel1) &&
		customerLimit.Level2.Equal(newLevel2) &&
		customerLimit.Level1Days == newLevel1Days &&
		customerLimit.Level2Days == newLevel2Days {
		return errors.BadRequest(error_code.ErrInvalidRequest, "no edd limit or days update")
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	beforeCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:        customerID,
		Role:              domain.RiskControlRole(customerLimit.Role),
		DailyLimit:        customerLimit.DailyLimit,
		MonthlyLimit:      customerLimit.MonthlyLimit,
		Level1:            customerLimit.Level1,
		Level2:            customerLimit.Level2,
		Level1Days:        customerLimit.Level1Days,
		Level2Days:        customerLimit.Level2Days,
		ChangeRoleReason:  customerLimit.ChangeRoleReason,
		ChangeLimitReason: customerLimit.ChangeLimitReason,
	}
	beforeCustomerLimitJsonData, err := json.Marshal(beforeCustomerLimit)
	if err != nil {
		return errors.InternalServer(error_code.ErrDBError, "json.Marshal(beforeCustomerLimit)").WithCause(err)
	}
	afterCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
		CustomerId:        customerID,
		Role:              domain.RiskControlRole(customerLimit.Role), // 固定不變
		DailyLimit:        customerLimit.DailyLimit,
		MonthlyLimit:      customerLimit.MonthlyLimit,
		Level1:            newLevel1,
		Level2:            newLevel2,
		Level1Days:        newLevel1Days,
		Level2Days:        newLevel2Days,
		ChangeRoleReason:  customerLimit.ChangeRoleReason, // 如果是白名單調整限額，則會更改角色，並帶入 系統自動切換 原因，否則該欄位固定不變
		ChangeLimitReason: customerLimit.ChangeLimitReason,
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

	customerLimit.Level1 = newLevel1
	customerLimit.Level2 = newLevel2
	customerLimit.Level1Days = newLevel1Days
	customerLimit.Level2Days = newLevel2Days
	customerLimit.UpdatedAt = time.Now()
	customerLimit.LastRole = beforeCustomerLimit.Role.Uint8()
	customerLimit.IsCustomizedEdd = true
	if err := tx.Save(&customerLimit).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *repository) UpdateAllCustomerLimitSettingWithoutCustomized(db *gorm.DB, operationUserId int64, newSetting domain.BTMRiskControlLimitSetting, reason string) error {
	if db == nil {
		return errors.InternalServer(error_code.ErrDBError, "db is nil")
	}

	// 只修改未客製化過的
	var customerLimits []model.BTMRiskControlCustomerLimitSetting
	if err := db.Where("is_customized = False and is_customized_edd = False").Find(&customerLimits).Error; err != nil {
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range customerLimits {
		beforeCustomerLimitJsonData, err := json.Marshal(BTMRiskControlCustomerLimitSettingModelToDomain(item))
		if err != nil {
			return errors.InternalServer(error_code.ErrDBError, "json.Marshal(beforeCustomerLimit)").WithCause(err)
		}

		afterCustomerLimit := domain.BTMRiskControlCustomerLimitSetting{
			Role:              domain.RiskControlRole(item.Role),
			CustomerId:        item.CustomerId,
			DailyLimit:        newSetting.DailyLimit,
			MonthlyLimit:      newSetting.MonthlyLimit,
			Level1:            newSetting.Level1,
			Level2:            newSetting.Level2,
			Level1Days:        newSetting.Level1Days,
			Level2Days:        newSetting.Level2Days,
			IsCustomized:      item.IsCustomized,
			IsCustomizedEdd:   item.IsCustomizedEdd,
			EddType:           item.EddType,
			ChangeRoleReason:  reason,
			ChangeLimitReason: reason,
		}
		afterCustomerLimitJsonData, err := json.Marshal(afterCustomerLimit)
		if err != nil {
			return errors.InternalServer(error_code.ErrDBError, "json.Marshal(afterCustomerLimit)").WithCause(err)
		}

		// 建立系統修改預設設定的 change log
		err = repo.CreateBTMChangeLog(tx, domain.BTMChangeLog{
			OperationUserId: operationUserId,
			TableName:       domain.BTMChangeLogTableNameBTMRiskControlCustomerLimitSetting,
			OperationType:   domain.BTMChangeLogOperationTypeUpdate,
			CustomerId:      nil,
			BeforeValue:     beforeCustomerLimitJsonData,
			AfterValue:      afterCustomerLimitJsonData,
		})
		if err != nil {
			return errors.InternalServer(error_code.ErrDBError, "CreateBTMChangeLog").WithCause(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	updateSettingIds := make([]uint, len(customerLimits))
	for i := range customerLimits {
		updateSettingIds[i] = customerLimits[i].ID
	}

	// 批量更新
	return db.Model(model.BTMRiskControlCustomerLimitSetting{}).
		Where("id IN (?)", updateSettingIds).
		Updates(model.BTMRiskControlCustomerLimitSetting{
			DailyLimit:        newSetting.DailyLimit,
			MonthlyLimit:      newSetting.MonthlyLimit,
			Level1:            newSetting.Level1,
			Level2:            newSetting.Level2,
			Level1Days:        newSetting.Level1Days,
			Level2Days:        newSetting.Level2Days,
			ChangeRoleReason:  reason,
			ChangeLimitReason: reason,
		}).Error
}

func BTMRiskControlCustomerLimitSettingDomainToModel(item domain.BTMRiskControlCustomerLimitSetting) model.BTMRiskControlCustomerLimitSetting {
	return model.BTMRiskControlCustomerLimitSetting{
		Role:            item.Role.Uint8(),
		CustomerId:      item.CustomerId,
		DailyLimit:      item.DailyLimit,
		MonthlyLimit:    item.MonthlyLimit,
		Level1:          item.Level1,
		Level2:          item.Level2,
		Level1Days:      item.Level1Days,
		Level2Days:      item.Level2Days,
		IsCustomized:    item.IsCustomized,
		IsCustomizedEdd: item.IsCustomizedEdd,
		EddType:         item.EddType,
	}
}

func BTMRiskControlCustomerLimitSettingModelToDomain(item model.BTMRiskControlCustomerLimitSetting) domain.BTMRiskControlCustomerLimitSetting {
	return domain.BTMRiskControlCustomerLimitSetting{
		ID:              item.ID,
		Role:            domain.RiskControlRole(item.Role),
		CustomerId:      item.CustomerId,
		DailyLimit:      item.DailyLimit,
		MonthlyLimit:    item.MonthlyLimit,
		Level1:          item.Level1,
		Level2:          item.Level2,
		Level1Days:      item.Level1Days,
		Level2Days:      item.Level2Days,
		IsCustomized:    item.IsCustomized,
		IsCustomizedEdd: item.IsCustomizedEdd,
		EddType:         item.EddType,
	}
}
