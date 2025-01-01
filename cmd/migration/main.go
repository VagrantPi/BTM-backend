package main

import (
	"BTM-backend/internal/di"
	"BTM-backend/internal/repo/model"
	"BTM-backend/third_party/db"
)

func main() {
	db := db.ConnectToDatabase()
	db.AutoMigrate(
		&model.BTMUser{},
		&model.BTMRole{},
		&model.BTMWhitelist{},
		&model.BTMLoginToken{},
		&model.BTM_CIB{},
	)

	// Initialize the repository
	repo, err := di.NewRepo()
	if err != nil {
		panic(err)
	}
	// Now you can call InitRawRole
	if err := repo.InitRawRole(); err != nil {
		panic(err)
	}
}
