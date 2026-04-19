package middleware

import (
	"log"
	"net/http"
)

func RecoveryPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if errPanic := recover(); errPanic != nil {
				log.Println(errPanic)
				http.Error(writer, "critical error on the server", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
