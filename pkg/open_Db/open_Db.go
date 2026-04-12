package open_Db

import (
	"to-do-list/app/configs"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OpenPostgres struct {
	*gorm.DB
}
type OpenRedis struct {
	*redis.Client
}

func NewOpenPostgres(conf *configs.DbConf) *OpenPostgres {
	db, errOpen := gorm.Open(postgres.Open(conf.DSN))
	if errOpen != nil {
		panic(errOpen)
	}
	return &OpenPostgres{
		DB: db,
	}
}
func NewOpenRedis(conf *configs.DbConf) *OpenRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: conf.RedisPassword,
		DB:       0,
	})
	return &OpenRedis{
		Client: client,
	}
}
