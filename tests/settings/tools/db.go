package tools

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/internal/domain"
	"BTM-backend/internal/repo/model"
	"BTM-backend/pkg/tools"
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func MigrateDB(db *gorm.DB) {
	if err := db.AutoMigrate(
		// lamassu
		// &model.Customer{},

		// BTM
		&model.BTMUser{},
		&model.BTMRole{},
		&model.BTMWhitelist{},
		&model.BTMLoginToken{},
		&model.BTM_CIB{},
		&model.BTMSumsub{},
		&model.BTMChangeLog{},

		// 2025_03_21_發票紀錄
		&model.BTMInvoice{},

		// 2025_03_21_新增限額功能
		&model.BTMRiskControlCustomerLimitSetting{},
		&model.BTMRiskControlLimitSetting{},
		// &model.BTMRiskControlThreshold{}, -- delete
		// &model.BTMRiskControlMachineRequestLimitLog{}, -- delete

		// 2025_03_24_新增後台登入日誌
		&model.BTMLoginLog{},

		// 2025_03_31_新增用戶備註
		&model.BTMCustomerNote{},

		// 2025_04_02_新增 新增限額塞入假資料 log
		&model.BTMMockTxHistoryLog{},

		// 2025_04_18_新增日交易快照
		&model.BTMDailyDeviceIncome{},
	); err != nil {
		fmt.Println("err1", err)
		panic(err)
	}

	// Initialize the repository
	repo, err := di.NewRepo(true)
	if err != nil {
		panic(err)
	}

	tx := repo.GetDb(context.Background())
	// Initialize all roles
	if err := repo.InitRawRole(tx); err != nil {
		panic(err)
	}

	// Initialize the admin
	if err := repo.InitAdmin(tx); err != nil {
		panic(err)
	}

	// init role
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.BTMRiskControlLimitSetting{
		ID:            1,
		Role:          1,
		DailyLimit:    decimal.NewFromInt(300000),
		MonthlyLimit:  decimal.NewFromInt(1000000),
		Level1:        decimal.NewFromInt(500000),
		Level2:        decimal.NewFromInt(2000000),
		Level1Days:    7,
		Level2Days:    60,
		VelocityDays:  1,
		VelocityTimes: 5,
	}).Error; err != nil {
		panic(err)
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.BTMRiskControlLimitSetting{
		ID:            2,
		Role:          2,
		DailyLimit:    decimal.NewFromInt(250000),
		MonthlyLimit:  decimal.NewFromInt(700000),
		Level1:        decimal.NewFromInt(400000),
		Level2:        decimal.NewFromInt(1500000),
		Level1Days:    7,
		Level2Days:    60,
		VelocityDays:  1,
		VelocityTimes: 5,
	}).Error; err != nil {
		panic(err)
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.BTMRiskControlLimitSetting{
		ID:            3,
		Role:          3,
		DailyLimit:    decimal.NewFromInt(0),
		MonthlyLimit:  decimal.NewFromInt(0),
		Level1:        decimal.NewFromInt(0),
		Level2:        decimal.NewFromInt(0),
		Level1Days:    7,
		Level2Days:    60,
		VelocityDays:  0,
		VelocityTimes: 0,
	}).Error; err != nil {
		panic(err)
	}
}

func InitApiToken(db *gorm.DB) string {
	token, err := tools.GenerateJWT(domain.UserJwt{
		Account: "admin",
		Role:    1,
		Id:      1,
	}, configs.C.JWT.Secret)
	if err != nil {
		panic(err)
	}

	err = db.Unscoped().Where("user_id = ?", 1).Delete(&model.BTMLoginToken{}).Error
	if err != nil {
		panic(err)
	}

	// init token
	err = db.Model(&model.BTMLoginToken{}).Create(&model.BTMLoginToken{
		UserID:     1,
		LoginToken: token,
	}).Error
	if err != nil {
		panic(err)
	}

	return token
}
