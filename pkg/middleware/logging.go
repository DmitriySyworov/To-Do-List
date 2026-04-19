package middleware

import (
	"log"
	"net/http"
)

type WrapperWriter struct {
	http.ResponseWriter
	statusCode int
}

func (wrapper *WrapperWriter) WriteHeader(status int) {
	wrapper.statusCode = status
	wrapper.ResponseWriter.WriteHeader(status)
}
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		wrapperWriter := &WrapperWriter{
			ResponseWriter: writer,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapperWriter, request)
		if wrapperWriter.statusCode >= 500 {
			log.Printf("Method: %s UrlPath: %s StatusCode: %d\n", request.Method, request.URL.Path, wrapperWriter.statusCode)
		}
	})
}
