package middleware

import "net/http"

func CORS(next http.Handler)http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	origin := request.Header.Get("Origin")
	if origin != "" {
		writer.Header().Set("Access-Control-Allow-Origin", origin)
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	if request.Method == http.MethodOptions {
		writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD,PATCH")
		writer.Header().Set("Access-Control-Allow-Headers", "authorization,content-type,content-length")
		writer.Header().Set("Access-Control-Max-Age", "86400")
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	next.ServeHTTP(writer, request)
})
}