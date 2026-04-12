package main

import (
	"to-do-list/app/configs"
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/open_Db"
)

func main() {
	conf := configs.NewConfigs()
	db := open_Db.NewOpenPostgres(conf.DbConf)
	errMigrate := db.AutoMigrate(&model.User{}, &model.TaskForm{})
	if errMigrate != nil {
		panic(errMigrate)
	}
}
