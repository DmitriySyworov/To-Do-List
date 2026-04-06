package handlerResponse

import (
	"encoding/json"
	"net/http"
)

func HandlerResponse(writer http.ResponseWriter, resp any, status int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	errEncode := json.NewEncoder(writer).Encode(resp)
	if errEncode != nil {
		panic(errEncode)
	}
}
