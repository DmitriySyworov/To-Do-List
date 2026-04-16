package user

import (
	"net/http"
	"to-do-list/app/configs"
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/handler_request"
	"to-do-list/app/pkg/handler_response"
	"to-do-list/app/pkg/middleware"
)

type HandlerUser struct {
	model.User
	*HandlerUserDep
}
type HandlerUserDep struct {
	*ServiceUser
	*configs.Configs
}

func NewHandlerUser(router *http.ServeMux, dep *HandlerUserDep) {
	user := &HandlerUser{
		HandlerUserDep: dep,
	}
	router.Handle("GET /my-user", middleware.IsAuthUserId(user.GetUser(), dep.Configs.Secret))
	router.Handle("PATCH /my-user", middleware.IsAuthUserId(user.UpdateUser(), dep.Configs.Secret))
	router.Handle("DELETE /my-user", middleware.IsAuthUserId(user.DeleteUser(), dep.Configs.Secret))
}
func (hl *HandlerUser) GetUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId, ok := request.Context().Value(middleware.KeyCtxUserId).(uint)
		if !ok {
			hl.User.Error = errors_custom.ErrToken.Error()
			handler_response.HandlerResponse(writer, hl.User, http.StatusUnauthorized)
			return
		}
		user, errGet := hl.RepositoryUsers.GetUserById(userId)
		if errGet != nil {
			hl.User.Error = errors_custom.ErrRecordNotFound.Error()
			handler_response.HandlerResponse(writer, hl.User, http.StatusNotFound)
			return
		}
		handler_response.HandlerResponse(writer, user, http.StatusOK)
	}
}
func (hl *HandlerUser) UpdateUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handler_request.ResultRequest[RequestUpdateUser](request)
		if errBody != nil {
			hl.User.Error = errBody.Error()
			handler_response.HandlerResponse(writer, hl.User, http.StatusBadRequest)
			return
		}
		userId, ok := request.Context().Value(middleware.KeyCtxUserId).(uint)
		if !ok {
			hl.User.Error = errors_custom.ErrToken.Error()
			handler_response.HandlerResponse(writer, hl.User, http.StatusUnauthorized)
			return
		}
		updateUser, errUpdate := hl.ServiceUser.UpdateUser(body, userId)
		if errUpdate != nil {
			hl.User.Error = errUpdate.Error()
			switch errUpdate {
			case errors_custom.ErrRecordNotFound:
				handler_response.HandlerResponse(writer, hl.User, http.StatusNotFound)
			case errors_custom.ErrIncorrectPassword, ErrParamsUpdateUser:
				handler_response.HandlerResponse(writer, hl.User, http.StatusBadRequest)
			default:
				handler_response.HandlerResponse(writer, hl.User, http.StatusInternalServerError)
			}
			return
		}
		handler_response.HandlerResponse(writer, updateUser, http.StatusOK)
	}
}
func (hl *HandlerUser) DeleteUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handler_request.ResultRequest[RequestDeleteUser](request)
		if errBody != nil {
			http.Error(writer, errBody.Error(), http.StatusBadRequest)
			return
		}
		userId, ok := request.Context().Value(middleware.KeyCtxUserId).(uint)
		if !ok {
			http.Error(writer, errors_custom.ErrToken.Error(), http.StatusUnauthorized)
			return
		}
		errDel := hl.ServiceUser.DeleteUser(body.Password, userId)
		if errDel != nil {
			switch errDel {
			case errors_custom.ErrRecordNotFound:
				http.Error(writer, errDel.Error(), http.StatusNotFound)
			case errors_custom.ErrIncorrectPassword:
				http.Error(writer, errDel.Error(), http.StatusBadRequest)
			case ErrDeleteUser:
				http.Error(writer, errDel.Error(), http.StatusInternalServerError)
			}
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}
