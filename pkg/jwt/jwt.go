package jwt

import (
	"time"
	"to-do-list/app/pkg/errors_custom"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret []byte
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: []byte(secret),
	}
}

const (
	jHashId    = "hash_id"
	jSessionId = "session_id"
	jUserId    = "user_id"
)

func (j *JWT) CreateTemporaryJWT(hashId float64, session string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		jHashId:      hashId,
		jSessionId:   session,
		"expires_at": time.Now().Add(time.Minute * 5).Unix(),
	})
	token, errToken := claims.SignedString(j.Secret)
	if errToken != nil {
		return "", errors_custom.ErrWriteData
	}
	return token, nil
}
func (j *JWT) CreateJWT(userId float64) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		jUserId: userId,
	})
	token, errToken := claims.SignedString(j.Secret)
	if errToken != nil {
		return "", errors_custom.ErrWriteData
	}
	return token, nil
}
func (j *JWT) ParseTemporaryJWt(token string) (float64, string, error) {
	value, errParse := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return j.Secret, nil
	})
	if errParse != nil || !value.Valid {
		return 0, "", errors_custom.ErrToken
	}
	hashId, ok := value.Claims.(jwt.MapClaims)[jHashId].(float64)
	if !ok {
		return 0, "", errors_custom.ErrToken
	}
	sessionId, ok := value.Claims.(jwt.MapClaims)[jSessionId].(string)
	if !ok {
		return 0, "", errors_custom.ErrToken
	}
	return hashId, sessionId, nil
}
func (j *JWT) ParseJWt(token string) (float64, error) {
	value, errParse := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return j.Secret, nil
	})
	if errParse != nil || !value.Valid {
		return 0, errors_custom.ErrToken
	}
	userId, ok := value.Claims.(jwt.MapClaims)[jUserId].(float64)
	if !ok {
		return 0, errors_custom.ErrToken
	}
	return userId, nil
}
