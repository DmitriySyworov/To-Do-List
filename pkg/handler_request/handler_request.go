package handler_request

import (
	"encoding/json"
	"io"
	"net/http"
	"to-do-list/app/pkg/errors_custom"

	"github.com/go-playground/validator/v10"
)

func ResultRequest[T any](request *http.Request) (*T, error) {
	var payload T
	data, errRead := io.ReadAll(request.Body)
	if errRead != nil {
		return nil, errors_custom.ErrIncorrectFormatData
	}
	errUnmarshal := json.Unmarshal(data, &payload)
	if errUnmarshal != nil {
		return nil, errors_custom.ErrIncorrectFormatData
	}
	newValidator := validator.New()
	errValidate := newValidator.Struct(payload)
	if errValidate != nil {
		return nil, errors_custom.ErrIncorrectData
	}
	return &payload, nil
}
