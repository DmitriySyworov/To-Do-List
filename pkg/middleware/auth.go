package middleware

import (
	"context"
	"net/http"
	"strings"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/jwt"
)

const (
	KeyCtxUserId    = "keyCtxUserId"
	KeyCtxSessionId = "keyCtxSessionId"
	KeyCtxHashId    = "KeyCtxHashId"
)

var CtxAuth = context.Background()

func validateFormatToken(header []string) (string, error) {
	if len(header) != 2 {
		return "", errors_custom.ErrToken
	}
	if strings.Count(header[1], ".") != 2 {
		return "", errors_custom.ErrToken
	}
	return header[1], nil
}
func IsAuthUserId(handler http.Handler, signature string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		header := request.Header.Get("Authorization")
		validToken, errValidateToken := validateFormatToken(strings.Split(header, "Bearer "))
		if errValidateToken != nil {
			http.Error(writer, errValidateToken.Error(), http.StatusUnauthorized)
			return
		}
		userId, errParse := jwt.NewJWT(signature).ParseJWt(validToken)
		if errParse != nil {
			http.Error(writer, errParse.Error(), http.StatusUnauthorized)
			return
		}
		CtxAuth = context.WithValue(CtxAuth, KeyCtxUserId, uint(userId))
		ctxRequest := request.WithContext(CtxAuth)
		handler.ServeHTTP(writer, ctxRequest)
	})
}
func IsUserToken(handler http.Handler, signature string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		header := request.Header.Get("X-User-Token")
		validToken, errValidateToken := validateFormatToken(strings.Split(header, "Bearer "))
		if errValidateToken != nil {
			http.Error(writer, errValidateToken.Error(), http.StatusUnauthorized)
			return
		}
		hashId, sessionId, errParse := jwt.NewJWT(signature).ParseTemporaryJWt(validToken)
		if errParse != nil {
			http.Error(writer, errParse.Error(), http.StatusUnauthorized)
			return
		}
		CtxAuth = context.WithValue(CtxAuth, KeyCtxHashId, uint(hashId))
		CtxAuth = context.WithValue(CtxAuth, KeyCtxSessionId, sessionId)
		ctxRequest := request.WithContext(CtxAuth)
		handler.ServeHTTP(writer, ctxRequest)
	})
}
