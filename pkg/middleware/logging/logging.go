package logging

import "net/http"

func Logging(handler http.Handler){


	handler.ServeHTTP()
}
