package handlerRequest

import (
	"encoding/json"
	"io"
	"net/http"
	"to-do-list/app/pkg/errorsCust"

	"github.com/go-playground/validator/v10"
)

func ResultRequest[T any](request *http.Request) (*T, error) {
	var payload T
	data, errRead := io.ReadAll(request.Body)
	if errRead != nil {
		return nil, errorsCust.ErrIncorrectFormatData
	}
	errUnmarshal := json.Unmarshal(data, &payload)
	if errUnmarshal != nil {
		return nil, errorsCust.ErrIncorrectFormatData
	}
	newValidator := validator.New()
	errValidate := newValidator.Struct(payload)
	if errValidate != nil {
		return nil, errorsCust.ErrIncorrectData
	}
	return &payload, nil
}
