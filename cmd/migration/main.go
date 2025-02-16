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
		&model.BTMChangeLog{},
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
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_btm_whitelist_address ON btm_whitelists (address);").Error; err != nil {
		panic(err)
	}
}
