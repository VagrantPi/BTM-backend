package main

import (
	"BTM-backend/configs"
	"BTM-backend/internal/di"
	"BTM-backend/internal/repo/model"
	"BTM-backend/third_party/db"
	"context"
)

func main() {
	db := db.ConnectToDatabase()
	if err := db.AutoMigrate(
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
		panic(err)
	}

	// Initialize the repository
	repo, err := di.NewRepo(configs.C.Mock)
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

	// migration
	// 2025_02_13_新增 udx 到 btm_whitelists
	// 	if err := db.Exec("DROP INDEX IF EXISTS idx_btm_whitelist_address;").Error; err != nil {
	// 		panic(err)
	// 	}
	// 	if err := db.Exec(`
	// DO $$
	// BEGIN
	//     IF NOT EXISTS (
	//         SELECT 1 FROM pg_indexes WHERE tablename = 'btm_whitelists' AND indexname = 'unique_address_idx'
	//     ) THEN
	//         CREATE UNIQUE INDEX unique_address_idx ON btm_whitelists (address) WHERE deleted_at IS NULL;
	//     END IF;
	// END $$;
	// `).Error; err != nil {
	// 		panic(err)
	// 	}

	// 2025_02_27_新增 idx 到 cash_in_txs
	// 	if err := db.Exec(`
	// DO $$
	// BEGIN
	//     IF NOT EXISTS (
	//         SELECT 1 FROM pg_indexes
	//         WHERE tablename = 'cash_in_txs'
	//         AND indexname = 'idx_cash_in_txs_fiat_nonzero'
	//     ) THEN
	//         CREATE INDEX idx_cash_in_txs_fiat_nonzero ON cash_in_txs (fiat) WHERE fiat != 0;
	//     END IF;
	// END $$;
	// `).Error; err != nil {
	// 		panic(err)
	// 	}

	// 2025_03_07_新增初始限額
	if err := db.Exec(`
	INSERT INTO "public"."btm_risk_control_limit_settings" ("role", "daily_limit", "monthly_limit", "created_at", "updated_at")
	VALUES
	    (1, '300000', '1000000', NOW(), NOW()),
	    (2, '250000', '700000', NOW(), NOW()),
	    (3, '0', '0', NOW(), NOW())
	ON CONFLICT ("role") DO NOTHING;
	`).Error; err != nil {
		panic(err)
	}

	// 2025_03_24_移除 btm_roles.role
	// if err := db.Exec(`ALTER TABLE IF EXISTS "public"."btm_roles" DROP COLUMN IF EXISTS "role";`).Error; err != nil {
	// 	panic(err)
	// }

	// 2025_03_25_新增 default no_role
	// if err := db.Exec(`
	// 	INSERT INTO "public"."btm_roles" ("role_name", "role_desc", "role_raw", "created_at")
	// 	VALUES ('no_role', 'default role', '[]', NOW())
	// 	ON CONFLICT ("role_name") DO NOTHING;
	// `).Error; err != nil {
	// 	panic(err)
	// }

	// 2025_03_26_加密 BTMSumsub
	// if err := db.Exec(`
	// 	ALTER TABLE IF EXISTS "public"."btm_sumsubs" DROP COLUMN IF EXISTS "email";
	// 	ALTER TABLE IF EXISTS "public"."btm_sumsubs" DROP COLUMN IF EXISTS "info";
	// `).Error; err != nil {
	// 	panic(err)
	// }

	// 2025_04_10_新增每種角色的等級門檻
	if err := db.Exec(`
		UPDATE "public"."btm_risk_control_limit_settings" SET "level1" = '500000', "level2" = '2000000', "level1_days" = 7, "level2_days" = 60 WHERE "role" = 1;
		UPDATE "public"."btm_risk_control_limit_settings" SET "level1" = '400000', "level2" = '1500000', "level1_days" = 7, "level2_days" = 60 WHERE "role" = 2;
		UPDATE "public"."btm_risk_control_limit_settings" SET "level1" = '0', "level2" = '0', "level1_days" = 7, "level2_days" = 60 WHERE "role" = 3;
	`).Error; err != nil {
		panic(err)
	}

	// 2025_04_24_新增每種角色交易次數限制
	if err := db.Exec(`
		UPDATE "public"."btm_risk_control_limit_settings" SET "velocity_days" = '1', "velocity_times" = 5 WHERE "role" = 1;
		UPDATE "public"."btm_risk_control_limit_settings" SET "velocity_days" = '1', "velocity_times" = 5 WHERE "role" = 2;
		UPDATE "public"."btm_risk_control_limit_settings" SET "velocity_days" = '0', "velocity_times" = 0 WHERE "role" = 3;
	`).Error; err != nil {
		panic(err)
	}
}
