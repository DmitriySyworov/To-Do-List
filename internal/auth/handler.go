package auth

import (
	"net/http"
	"to-do-list/app/configs"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/handler_request"
	"to-do-list/app/pkg/handler_response"
	"to-do-list/app/pkg/middleware"
)

type HandlerAuth struct {
	ResponseAuth
	ResponseConfirm
	*HandlerAuthDep
}
type HandlerAuthDep struct {
	*ServiceAuth
	*configs.Configs
}

func NewHandlerAuth(router *http.ServeMux, dep *HandlerAuthDep) {
	auth := &HandlerAuth{
		HandlerAuthDep: dep,
	}
	router.HandleFunc("POST /auth/register", auth.Register())
	router.HandleFunc("POST /auth/login", auth.Login())
	router.HandleFunc("POST /auth/restore", auth.Restore())
	router.Handle("POST /auth/confirm", middleware.IsUserToken(auth.Confirm(), dep.Secret))

}
func (hl *HandlerAuth) Register() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handler_request.ResultRequest[RequestRegister](request)
		if errBody != nil {
			if errBody != nil {
				hl.ResponseAuth.Error = errBody.Error()
				handler_response.HandlerResponse(writer, hl.ResponseAuth, http.StatusBadRequest)
				return
			}
		}
		respAuth, errAuth := hl.ServiceAuth.Register(body)
		if errAuth != nil {
			hl.ResponseAuth.Error = errAuth.Error()
			switch errAuth {
			case ErrAlreadyExist:
				handler_response.HandlerResponse(writer, hl.ResponseAuth, http.StatusUnauthorized)
			default:
				handler_response.HandlerResponse(writer, hl.ResponseAuth, http.StatusInternalServerError)
			}
			return
		}
		handler_response.HandlerResponse(writer, respAuth, http.StatusOK)
	}
}
func (hl *HandlerAuth) Login() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handler_request.ResultRequest[RequestLoginAndRestore](request)
		if errBody != nil {
			if errBody != nil {
				hl.ResponseAuth.Error = errBody.Error()
				handler_response.HandlerResponse(writer, hl.ResponseAuth, http.StatusBadRequest)
				return
			}
		}
		respAuth, errAuth := hl.ServiceAuth.Login(body)
		if errAuth != nil {
			hl.ResponseAuth.Error = errAuth.Error()
			switch errAuth {
			case errors_custom.ErrRecordNotFound, errors_custom.ErrIncorrectPassword:
				handler_response.HandlerResponse(writer, hl.ResponseAuth, http.StatusUnauthorized)
			default:
				handler_response.HandlerResponse(writer, hl.ResponseAuth, http.StatusInternalServerError)
			}
			return
		}
		handler_response.HandlerResponse(writer, respAuth, http.StatusOK)
	}
}
func (hl *HandlerAuth) Restore() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handler_request.ResultRequest[RequestLoginAndRestore](request)
		if errBody != nil {
			hl.ResponseAuth.Error = errBody.Error()
			handler_response.HandlerResponse(writer, hl.ResponseAuth, http.StatusBadRequest)
			return
		}
		action := request.URL.Query().Get("action")
		if action != "recoverLogin" && action != "recoverDelete" {
			respAuth, errAuth := hl.ServiceAuth.Restore(body, action)
			if errAuth != nil {
				hl.ResponseAuth.Error = errAuth.Error()
				switch errAuth {
				case errors_custom.ErrRecordNotFound:
					handler_response.HandlerResponse(writer, hl.ResponseAuth, http.StatusUnauthorized)
				default:
					handler_response.HandlerResponse(writer, hl.ResponseAuth, http.StatusInternalServerError)
				}
				return
			}
			handler_response.HandlerResponse(writer, respAuth, http.StatusOK)
		}
	}
}
func (hl *HandlerAuth) Confirm() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handler_request.ResultRequest[RequestConfirm](request)
		if errBody != nil {
			hl.ResponseConfirm.Error = errBody.Error()
			handler_response.HandlerResponse(writer, hl.ResponseConfirm, http.StatusBadRequest)
			return
		}
		hashId, okHashId := request.Context().Value(middleware.KeyCtxHashId).(uint)
		sessionId, okSession := request.Context().Value(middleware.KeyCtxSessionId).(string)
		if !okHashId || !okSession {
			hl.ResponseConfirm.Error = errors_custom.ErrToken.Error()
			handler_response.HandlerResponse(writer, hl.ResponseConfirm, http.StatusBadRequest)
			return
		}
		action := request.URL.Query().Get("action")
		if action != "restore" && action != "register" && action != "login" && action != "recoverLogin" && action != "recoverDelete" {
			hl.ResponseConfirm.Error = ErrIncorrectAction.Error()
			handler_response.HandlerResponse(writer, hl.ResponseConfirm, http.StatusBadRequest)
			return
		}
		if action == "recoverLogin" && body.NewPassword == "" {
			hl.ResponseConfirm.Error = ErrNewPassword.Error()
			handler_response.HandlerResponse(writer, hl.ResponseConfirm, http.StatusBadRequest)
			return
		}
		respConfirm, errConfirm := hl.ServiceAuth.Confirm(body, hashId, sessionId, action)
		if errConfirm != nil {
			hl.ResponseConfirm.Error = errConfirm.Error()
			switch errConfirm {
			case errors_custom.ErrRecordNotFound, ErrValidSession, ErrIncorrectCode:
				handler_response.HandlerResponse(writer, hl.ResponseConfirm, http.StatusUnauthorized)
			default:
				handler_response.HandlerResponse(writer, hl.ResponseConfirm, http.StatusInternalServerError)
			}
			return
		}
		if action == "register" {
			handler_response.HandlerResponse(writer, respConfirm, http.StatusCreated)
		} else {
			handler_response.HandlerResponse(writer, respConfirm, http.StatusOK)
		}
	}
}
