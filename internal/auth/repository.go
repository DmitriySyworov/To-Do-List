package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/open_Db"
)

type RepositoryAuth struct {
	*open_Db.OpenRedis
	RedisCtx context.Context
}

func NewRepositoryAuth(redis *open_Db.OpenRedis, redisCtx context.Context) *RepositoryAuth {
	return &RepositoryAuth{
		OpenRedis: redis,
		RedisCtx:  redisCtx,
	}
}

const (
	keyName     = "name"
	keyEmail    = "email"
	keyPassword = "password"
	keyUserId   = "user_id"

	keySessionId = "session_id"
	keyTempCode  = "temporary_code"
)

func (r *RepositoryAuth) CreateTempUser(tempUser *model.TempUser, idHash uint) error {
	key := fmt.Sprintf("user:%d", idHash)
	errHSet := r.Client.HSet(r.RedisCtx, key, keyName, tempUser.Name, keyEmail, tempUser.Email, keyPassword, tempUser.Password, keyUserId, tempUser.UserId).Err()
	if errHSet != nil {
		return errHSet
	}
	errExpire := r.Client.Expire(r.RedisCtx, key, time.Minute*5).Err()
	if errExpire != nil {
		return errExpire
	}
	return nil
}
func (r *RepositoryAuth) CreateSession(session *model.Session, idHash uint) error {
	key := fmt.Sprintf("session:%d", idHash)
	errHSet := r.Client.HSet(r.RedisCtx, key, keySessionId, session.SessionId, keyTempCode, session.TemporaryCode).Err()
	if errHSet != nil {
		return errHSet
	}
	errExpire := r.Client.Expire(r.RedisCtx, key, time.Minute*5).Err()
	if errExpire != nil {
		return errExpire
	}
	return nil
}
func (r *RepositoryAuth) GetTempUser(idHash uint) (*model.TempUser, error) {
	key := fmt.Sprintf("user:%d", idHash)
	mapValue, errHGetAll := r.Client.HGetAll(r.RedisCtx, key).Result()
	if errHGetAll != nil {
		return nil, errHGetAll
	}
	userId, errCode := strconv.Atoi(mapValue[keyUserId])
	if errCode != nil {
		return nil, errCode
	}
	return &model.TempUser{
		Name:     mapValue[keyName],
		Email:    mapValue[keyEmail],
		Password: mapValue[keyPassword],
		UserId:   uint(userId),
	}, nil
}
func (r *RepositoryAuth) GetSession(idHash uint) (*model.Session, error) {
	key := fmt.Sprintf("session:%d", idHash)
	mapValue, errHGetAll := r.Client.HGetAll(r.RedisCtx, key).Result()
	if errHGetAll != nil {
		return nil, errHGetAll
	}
	tempCode, errCode := strconv.Atoi(mapValue[keyTempCode])
	if errCode != nil {
		return nil, errCode
	}
	return &model.Session{
		SessionId:     mapValue[keySessionId],
		TemporaryCode: uint(tempCode),
	}, nil
}
