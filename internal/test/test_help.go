package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"to-do-list/app/internal/model"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbTest struct {
	*gorm.DB
	*redis.Client
}

func OpenAllTestDb() *DbTest {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	db, errOpen := gorm.Open(postgres.Open(os.Getenv("DSN_TEST")))
	if errOpen != nil {
		panic(errOpen)
	}
	return &DbTest{
		DB: db,
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6380",
			Password: os.Getenv("REDIS_TEST"),
			DB:       0,
		})}
}

var BaseTestUser = model.User{
	Name:     UserTestName,
	Email:    UserTestEmail,
	Password: HashTestUserPassword,
	UserId:   UserTestID,
}

const (
	UserTestName         = "Bob"
	UserTestID           = 33013910811
	UserTestEmail        = "xzinkes1@gmail.com"
	OrigTestUserPassword = "znxlczxlcnopiqos;mk"
	HashTestUserPassword = "$2a$10$SXdiX9TwSVQ3/rWgz9UNEu8d.r3HgdM3HcnWHe1BbU1KmLgj2n44S"
	JWTTestToken         = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozMzAxMzkxMDgxMX0.v5S50cSssvdgHhjAmP3uRds1-PmDsLMf19L0O4hArfk"
)

func (db *DbTest) CreateTestUser(user *model.User) {
	res := db.DB.Create(&user)
	if res.Error != nil {
		log.Println(res.Error)
	}
}

var redisTestCtx = context.Background()

func (db *DbTest) CleanAllDb() {
	res := db.DB.Exec("TRUNCATE TABLE task_forms, users")
	if res.Error != nil {
		log.Println(res.Error)
	}
	errCleanRedis := db.Client.FlushAll(redisTestCtx).Err()
	if errCleanRedis != nil {
		log.Println(errCleanRedis)
	}
}

const (
	fieldCreate = "Create_Task"
	fieldDone   = "Done_task"
	fieldDelete = "Delete_Task"
	fieldName   = "Name"
)

func (db *DbTest) CreateStat(userId uint, name string) {
	key := fmt.Sprintf("task:%d", userId)
	errHSet := db.Client.HSet(redisTestCtx, key, fieldCreate, 1, fieldDelete, 0, fieldDone, 0, fieldName, name).Err()
	if errHSet != nil {
		log.Println(errHSet)
	}
}
