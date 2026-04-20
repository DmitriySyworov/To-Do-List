package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"to-do-list/app/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbTest struct {
	*gorm.DB
	*redis.Client
	Secret []byte
}

func OpenAllTestDb() *DbTest {
	errEnv := godotenv.Load(".env.test")
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
			Addr:     os.Getenv("REDIS_PORT_TEST"),
			Password: os.Getenv("REDIS_TEST"),
			DB:       0,
		}),
		Secret: []byte(os.Getenv("SECRET_TEST")),
	}
}
func (db *DbTest) CreateTemporaryJWTTest(hashId float64, session string) string {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"hash_id":    hashId,
		"session_id": session,
		"expires_at": time.Now().Add(time.Minute * 5).Unix(),
	})
	token, errToken := claims.SignedString(db.Secret)
	if errToken != nil {
		log.Println(errToken)
	}
	return token
}

var BaseUserTest = model.User{
	Name:     UserNameTest,
	Email:    UserEmailTest,
	Password: HashUserPasswordTest,
	UserId:   UserIDTest,
}

const (
	UserNameTest              = "Bob"
	UserIDTest                = 33013910811
	UserEmailTest             = "xzinkes1@gmail.com"
	OrigUserPasswordTest      = "znxlczxlcnopiqos;mk"
	HashUserPasswordTest      = "$2a$10$SXdiX9TwSVQ3/rWgz9UNEu8d.r3HgdM3HcnWHe1BbU1KmLgj2n44S"
	IdJWTTest                 = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozMzAxMzkxMDgxMX0.v5S50cSssvdgHhjAmP3uRds1-PmDsLMf19L0O4hArfk"
	IdHashTest           uint = 645830281
	TempCodeTest              = 345123
	SessionIdTest             = "ashWxm"
)

func (db *DbTest) CreateUserTest(user *model.User) {
	res := db.DB.Create(&user)
	if res.Error != nil {
		log.Println(res.Error)
	}
}

var BaseTaskTest = model.TaskForm{
	Header:   HeaderTest,
	Task:     TaskTest,
	Deadline: dateTest(),
	TaskId:   TaskIdTest,
	UserId:   UserIDTest,
}

func dateTest() time.Time {
	date, errDate := time.Parse(time.DateOnly, DateTest)
	if errDate != nil {
		log.Println(errDate)
	}
	return date
}

const (
	HeaderTest = "gotta to relax"
	TaskTest   = "one come to home"
	DateTest   = "2026-10-03"
	TaskIdTest = 3452859
)

func (db *DbTest) CreateTaskTest(task *model.TaskForm) {
	res := db.DB.Create(&task)
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
	fieldCreate  = "Create_Task"
	fieldDone    = "Done_task"
	fieldDelete  = "Delete_Task"
	fieldName    = "Name"
	keyName      = "name"
	keyEmail     = "email"
	keyPassword  = "password"
	keyUserId    = "user_id"
	keySessionId = "session_id"
	keyTempCode  = "temporary_code"
)

func (db *DbTest) CreateStatTest(userId uint, name string) {
	key := fmt.Sprintf("task:%d", userId)
	errHSet := db.Client.HSet(redisTestCtx, key, fieldCreate, 1, fieldDelete, 0, fieldDone, 1, fieldName, name).Err()
	if errHSet != nil {
		log.Println(errHSet)
	}
}
func (db *DbTest) CreateDoneStatTest(userId, quantityDone uint, name string) {
	key := fmt.Sprintf("task:%d", userId)
	errHSet := db.Client.HSet(redisTestCtx, key, fieldCreate, 1, fieldDelete, 0, fieldDone, quantityDone, fieldName, name).Err()
	if errHSet != nil {
		log.Println(errHSet)
	}
}
func (db *DbTest) CreateTempUserTest(tempUser *model.TempUser, idHash uint) {
	key := fmt.Sprintf("user:%d", idHash)
	errHSet := db.Client.HSet(redisTestCtx, key, keyName, tempUser.Name, keyEmail, tempUser.Email, keyPassword, tempUser.Password, keyUserId, tempUser.UserId).Err()
	if errHSet != nil {
		log.Println(errHSet)
	}
	errExpire := db.Client.Expire(redisTestCtx, key, time.Minute*5).Err()
	if errExpire != nil {
		log.Println(errExpire)
	}
}
func (db *DbTest) CreateSessionTest(session *model.Session, idHash uint) {
	key := fmt.Sprintf("session:%d", idHash)
	errHSet := db.Client.HSet(redisTestCtx, key, keySessionId, session.SessionId, keyTempCode, session.TemporaryCode).Err()
	if errHSet != nil {
		log.Println(errHSet)
	}
	errExpire := db.Client.Expire(redisTestCtx, key, time.Minute*5).Err()
	if errExpire != nil {
		log.Println(errExpire)
	}
}
