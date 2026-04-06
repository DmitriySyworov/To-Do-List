package auth

import (
	"net/http"
	"to-do-list/app/pkg/errorsCust"
	"to-do-list/app/pkg/handlerRequest"
	"to-do-list/app/pkg/handlerResponse"
)

type HandlerAuth struct {
	ResponseAuth
	ResponseConfirm
	*HandlerAuthDep
}
type HandlerAuthDep struct {
	*ServiceAuth
}

func NewHandlerAuth(router *http.ServeMux, dep *HandlerAuthDep) {
	auth := &HandlerAuth{
		HandlerAuthDep: dep,
	}
	router.HandleFunc("POST /auth/register", auth.Register())
	router.HandleFunc("POST /auth/login", auth.Login())
	router.HandleFunc("POST /auth/restore", auth.Restore())
	router.Handle("POST /auth/confirm", auth.Confirm())

}
func (hl *HandlerAuth) Register() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handlerRequest.ResultRequest[RequestRegister](request)
		if errBody != nil {
			if errBody != nil {
				hl.ResponseAuth.Error = errBody.Error()
				handlerResponse.HandlerResponse(writer, hl.ResponseAuth, http.StatusBadRequest)
				return
			}
		}
		respAuth, errAuth := hl.ServiceAuth.Register(body)
		if errAuth != nil {
			hl.ResponseAuth.Error = errAuth.Error()
			switch errAuth {
			case ErrAlreadyExist:
				handlerResponse.HandlerResponse(writer, hl.ResponseAuth, http.StatusUnauthorized)
			default:
				handlerResponse.HandlerResponse(writer, hl.ResponseAuth, http.StatusInternalServerError)
			}
			return
		}
		handlerResponse.HandlerResponse(writer, respAuth, http.StatusOK)
	}
}
func (hl *HandlerAuth) Login() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handlerRequest.ResultRequest[RequestLoginAndRestore](request)
		if errBody != nil {
			if errBody != nil {
				hl.ResponseAuth.Error = errBody.Error()
				handlerResponse.HandlerResponse(writer, hl.ResponseAuth, http.StatusBadRequest)
				return
			}
		}
		respAuth, errAuth := hl.ServiceAuth.Login(body)
		if errAuth != nil {
			hl.ResponseAuth.Error = errAuth.Error()
			switch errAuth {
			case errorsCust.ErrRecordNotFound, ErrIncorrectPassword:
				handlerResponse.HandlerResponse(writer, hl.ResponseAuth, http.StatusUnauthorized)
			default:
				handlerResponse.HandlerResponse(writer, hl.ResponseAuth, http.StatusInternalServerError)
			}
			return
		}
		handlerResponse.HandlerResponse(writer, respAuth, http.StatusOK)
	}
}
func (hl *HandlerAuth) Restore() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handlerRequest.ResultRequest[RequestLoginAndRestore](request)
		if errBody != nil {
			hl.ResponseAuth.Error = errBody.Error()
			handlerResponse.HandlerResponse(writer, hl.ResponseAuth, http.StatusBadRequest)
			return
		}
		respAuth, errAuth := hl.ServiceAuth.Restore(body)
		if errAuth != nil {
			hl.ResponseAuth.Error = errAuth.Error()
			switch errAuth {
			case errorsCust.ErrRecordNotFound:
				handlerResponse.HandlerResponse(writer, hl.ResponseAuth, http.StatusUnauthorized)
			default:
				handlerResponse.HandlerResponse(writer, hl.ResponseAuth, http.StatusInternalServerError)
			}
			return
		}
		handlerResponse.HandlerResponse(writer, respAuth, http.StatusOK)
	}
}
func (hl *HandlerAuth) Confirm() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handlerRequest.ResultRequest[RequestConfirm](request)
		if errBody != nil {
			hl.ResponseConfirm.Error = errBody.Error()
			handlerResponse.HandlerResponse(writer, hl.ResponseConfirm, http.StatusBadRequest)
			return
		}
		//

		//тут добавить авторизацию

		//
		action := request.URL.Query().Get("action")
		if action != "restore" && action != "register" && action != "login" {
			hl.ResponseConfirm.Error = ErrIncorrectAction.Error()
			handlerResponse.HandlerResponse(writer, hl.ResponseConfirm, http.StatusBadRequest)
			return
		}
		respConfirm, errConfirm := hl.ServiceAuth.Confirm(0, body.TempCode, "", action)
		if errConfirm != nil {
			hl.ResponseConfirm.Error = errConfirm.Error()
			switch errConfirm {
			case errorsCust.ErrRecordNotFound, ErrValidSession, ErrIncorrectCode:
				handlerResponse.HandlerResponse(writer, hl.ResponseConfirm, http.StatusUnauthorized)
			default:
				handlerResponse.HandlerResponse(writer, hl.ResponseConfirm, http.StatusInternalServerError)
			}
			return
		}
		if action == "register" {
			handlerResponse.HandlerResponse(writer, respConfirm, http.StatusCreated)
		} else {
			handlerResponse.HandlerResponse(writer, respConfirm, http.StatusOK)
		}
	}
}
