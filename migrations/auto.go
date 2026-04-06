package main

import (
	"to-do-list/app/configs"
	"to-do-list/app/internal/models"
	"to-do-list/app/pkg/openDb"
)

func main() {
	conf := configs.NewConfigs()
	db := openDb.NewOpenPostgres(conf.DbConf)
	errMigrate := db.AutoMigrate(&models.Users{})
	if errMigrate != nil {
		panic(errMigrate)
	}
}
