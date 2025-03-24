package main

import (
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
		// &model.BTMRiskControlMachineRequestLimitLog{},
		// &model.BTMRiskControlThreshold{},
	); err != nil {
		panic(err)
	}

	// Initialize the repository
	repo, err := di.NewRepo()
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
	if err := db.Exec("DROP INDEX IF EXISTS idx_btm_whitelist_address;").Error; err != nil {
		panic(err)
	}
	if err := db.Exec(`
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes WHERE tablename = 'btm_whitelists' AND indexname = 'unique_address_idx'
    ) THEN
        CREATE UNIQUE INDEX unique_address_idx ON btm_whitelists (address) WHERE deleted_at IS NULL;
    END IF;
END $$;
`).Error; err != nil {
		panic(err)
	}

	// 2025_02_27_新增 idx 到 cash_in_txs
	if err := db.Exec(`
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE tablename = 'cash_in_txs' 
        AND indexname = 'idx_cash_in_txs_fiat_nonzero'
    ) THEN
        CREATE INDEX idx_cash_in_txs_fiat_nonzero ON cash_in_txs (fiat) WHERE fiat != 0;
    END IF;
END $$;
`).Error; err != nil {
		panic(err)
	}

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

	// 2025_03_21_新增初始白名單與灰名單門檻
	if err := db.Exec(`
		INSERT INTO "public"."btm_risk_control_thresholds" ("role", "threshold", "threshold_days", "created_at") VALUES (1, '500000', 7, NOW()) ON CONFLICT ("role", "threshold") DO NOTHING;
		INSERT INTO "public"."btm_risk_control_thresholds" ("role", "threshold", "threshold_days", "created_at") VALUES (1, '2000000', 60, NOW()) ON CONFLICT ("role", "threshold") DO NOTHING;
		INSERT INTO "public"."btm_risk_control_thresholds" ("role", "threshold", "threshold_days", "created_at") VALUES (2, '400000', 7, NOW()) ON CONFLICT ("role", "threshold") DO NOTHING;
		INSERT INTO "public"."btm_risk_control_thresholds" ("role", "threshold", "threshold_days", "created_at") VALUES (2, '1500000', 60, NOW()) ON CONFLICT ("role", "threshold") DO NOTHING;
	`).Error; err != nil {
		panic(err)
	}
}
