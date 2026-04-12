package stat

import (
	"context"
	"fmt"
	"log"
	"strings"
	"to-do-list/app/pkg/open_Db"
)

type RepositoryStat struct {
	*open_Db.OpenRedis
	RedisCtx context.Context
}

func NewRepositoryStat(redis *open_Db.OpenRedis, redisCtx context.Context) *RepositoryStat {
	return &RepositoryStat{
		OpenRedis: redis,
		RedisCtx:  redisCtx,
	}
}

const (
	fieldCreate = "Create_Task"
	fieldDone   = "Done_task"
	fieldDelete = "Delete_Task"
	fieldName   = "Name"
)

func (r *RepositoryStat) GetStatUser(userId uint) (*ResponseMyStat, error) {
	key := fmt.Sprintf("task:%d", userId)
	mapFields, errHGetAll := r.Client.HGetAll(r.RedisCtx, key).Result()
	if errHGetAll != nil {
		return nil, errHGetAll
	}
	return &ResponseMyStat{
		QuantityActiveTask: mapFields[fieldCreate],
		QuantityDoneTask:   mapFields[fieldDone],
		QuantityDeleteTask: mapFields[fieldDelete],
	}, nil
}
func (r *RepositoryStat) GetLeaderboard() ([]UserStat, error) {
	AllKeys, errKey := r.Client.Keys(r.RedisCtx, "*").Result()
	if errKey != nil {
		return nil, errKey
	}
	var leaderboard []UserStat
	for _, key := range AllKeys {
		var tempLeaderboard UserStat
		if strings.Contains(key, "task:") {
			mapFields, errHGetAll := r.Client.HGetAll(r.RedisCtx, key).Result()
			if errHGetAll != nil {
				log.Println(errHGetAll)
			}
			tempLeaderboard.QuantityDoneTask = mapFields[fieldDone]
			tempLeaderboard.Name = mapFields[fieldName]
			leaderboard = append(leaderboard, tempLeaderboard)
		}
	}
	return leaderboard, nil
}
func (r *RepositoryStat) AddCreateTask(userId uint, name string) error {
	key := fmt.Sprintf("task:%d", userId)
	errKey := r.Client.Keys(r.RedisCtx, key).Err()
	if errKey != nil {
		errHSet := r.Client.HSet(r.RedisCtx, key, fieldCreate, 1, fieldDelete, 0, fieldDone, 0, fieldName, name).Err()
		if errHSet != nil {
			return errHSet
		}
	} else {
		errIncr := r.Client.HIncrBy(r.RedisCtx, key, fieldCreate, 1).Err()
		if errIncr != nil {
			return errIncr
		}
	}
	return nil
}
func (r *RepositoryStat) AddDoneTask(userId uint) error {
	key := fmt.Sprintf("task:%d", userId)
	errIncr := r.Client.HIncrBy(r.RedisCtx, key, fieldDone, 1).Err()
	if errIncr != nil {
		return errIncr
	}
	errDecr := r.Client.HIncrBy(r.RedisCtx, key, fieldCreate, -1).Err()
	if errDecr != nil {
		return errDecr
	}
	return nil
}
func (r *RepositoryStat) AddDeleteTask(userId uint, choiceField string) error {
	key := fmt.Sprintf("task:%d", userId)
	errIncr := r.Client.HIncrBy(r.RedisCtx, key, fieldDelete, 1).Err()
	if errIncr != nil {
		return errIncr
	}
	errDecr := r.Client.HIncrBy(r.RedisCtx, key, choiceField, -1).Err()
	if errDecr != nil {
		return errDecr
	}
	return nil
}
