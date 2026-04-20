package main

import (
	"os"
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/open_Db"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("GO_TEST") == "test" {
		errEnv := godotenv.Load(".env.test")
		if errEnv != nil {
			panic(errEnv)
		}
		db := open_Db.NewOpenPostgres(os.Getenv("DSN_TEST"))
		errMigrate := db.AutoMigrate(&model.User{}, &model.TaskForm{})
		if errMigrate != nil {
			panic(errMigrate)
		}
	} else {
		errEnv := godotenv.Load(".env")
		if errEnv != nil {
			panic(errEnv)
		}
		db := open_Db.NewOpenPostgres(os.Getenv("DSN"))
		errMigrate := db.AutoMigrate(&model.User{}, &model.TaskForm{})
		if errMigrate != nil {
			panic(errMigrate)
		}
	}
}
