package open_Db

import (
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

func NewOpenPostgres(DSN string) *OpenPostgres {
	db, errOpen := gorm.Open(postgres.Open(DSN))
	if errOpen != nil {
		panic(errOpen)
	}
	return &OpenPostgres{
		DB: db,
	}
}
func NewOpenRedis(redisPass string) *OpenRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: redisPass,
		DB:       0,
	})
	return &OpenRedis{
		Client: client,
	}
}
