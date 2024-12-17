package main

import (
	_ "BTM-backend/configs"
	"BTM-backend/internal/repo/model"
	"BTM-backend/third_party/db"
)

func main() {
	db := db.ConnectToDatabase()
	db.AutoMigrate(&model.BTMUser{}, &model.BTMWhitelist{}, &model.BTMLoginToken{})
}
