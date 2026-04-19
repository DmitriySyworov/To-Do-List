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
}

func NewRepositoryAuth(redis *open_Db.OpenRedis) *RepositoryAuth {
	return &RepositoryAuth{
		OpenRedis: redis,
	}
}

const (
	keyName     = "name"
	keyEmail    = "email"
	keyPassword = "password"
	keyUserId   = "user_id"

	keySessionId = "session_id"
	keyTempCode  = "temporary_code"

	timeout = 30
)

func (r *RepositoryAuth) CreateTempUser(tempUser *model.TempUser, idHash uint) error {
	key := fmt.Sprintf("user:%d", idHash)
	redisCtx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)
	defer cancel()
	errHSet := r.Client.HSet(redisCtx, key, keyName, tempUser.Name, keyEmail, tempUser.Email, keyPassword, tempUser.Password, keyUserId, tempUser.UserId).Err()
	if errHSet != nil {
		return errHSet
	}
	errExpire := r.Client.Expire(redisCtx, key, time.Minute*5).Err()
	if errExpire != nil {
		return errExpire
	}
	return nil
}
func (r *RepositoryAuth) CreateSession(session *model.Session, idHash uint) error {
	key := fmt.Sprintf("session:%d", idHash)
	redisCtx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)
	defer cancel()
	errHSet := r.Client.HSet(redisCtx, key, keySessionId, session.SessionId, keyTempCode, session.TemporaryCode).Err()
	if errHSet != nil {
		return errHSet
	}
	errExpire := r.Client.Expire(redisCtx, key, time.Minute*5).Err()
	if errExpire != nil {
		return errExpire
	}
	return nil
}
func (r *RepositoryAuth) GetTempUser(idHash uint) (*model.TempUser, error) {
	key := fmt.Sprintf("user:%d", idHash)
	redisCtx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)
	defer cancel()
	mapValue, errHGetAll := r.Client.HGetAll(redisCtx, key).Result()
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
	redisCtx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)
	defer cancel()
	key := fmt.Sprintf("session:%d", idHash)
	mapValue, errHGetAll := r.Client.HGetAll(redisCtx, key).Result()
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
