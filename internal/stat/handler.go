package stat

import (
	"net/http"
	"to-do-list/app/configs"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/handler_response"
	"to-do-list/app/pkg/middleware"
)

type HandlerStat struct {
	*HandlerStatDep
	ResponseMyStat
	ResponseLeaderboard
}
type HandlerStatDep struct {
	*ServiceStat
	*configs.Configs
}

func NewHandlerStat(router *http.ServeMux, dep *HandlerStatDep) {
	stat := &HandlerStat{
		HandlerStatDep: dep,
	}
	router.Handle("/user/my-stat", middleware.IsAuthUserId(stat.GetMyStat(), dep.Secret))
	router.Handle("/user/leaderboard", middleware.IsAuthUserId(stat.GetLeaderboard(), dep.Secret))
}

func (hl *HandlerStat) GetMyStat() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId, ok := request.Context().Value(middleware.KeyCtxUserId).(uint)
		if !ok {
			hl.ResponseMyStat.Error = errors_custom.ErrToken.Error()
			handler_response.HandlerResponse(writer, hl.ResponseMyStat, http.StatusUnauthorized)
			return
		}
		respStat, errGetStat := hl.RepositoryStat.GetStatUser(userId)
		if errGetStat != nil {
			hl.ResponseMyStat.Error = ErrNotFoundStat.Error()
			handler_response.HandlerResponse(writer, hl.ResponseMyStat, http.StatusNotFound)
			return
		}
		handler_response.HandlerResponse(writer, respStat, http.StatusOK)
	}
}
func (hl *HandlerStat) GetLeaderboard() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId, ok := request.Context().Value(middleware.KeyCtxUserId).(uint)
		if !ok {
			hl.ResponseLeaderboard.Error = errors_custom.ErrToken.Error()
			handler_response.HandlerResponse(writer, hl.ResponseLeaderboard, http.StatusUnauthorized)
			return
		}
		limit := request.URL.Query().Get("limit")
		if limit == "" || limit == "0" {
			hl.ResponseLeaderboard.Error = ErrLimit.Error()
			handler_response.HandlerResponse(writer, hl.ResponseLeaderboard, http.StatusBadRequest)
			return
		}
		respLeaderboard, errResp := hl.ServiceStat.GetLeaderBoard(userId, limit)
		if errResp != nil {
			hl.ResponseLeaderboard.Error = respLeaderboard.Error
			switch errResp {
			case ErrLimit:
				handler_response.HandlerResponse(writer, hl.ResponseLeaderboard, http.StatusBadRequest)
			case errors_custom.ErrNoExistUser:
				handler_response.HandlerResponse(writer, hl.ResponseLeaderboard, http.StatusUnauthorized)
			case ErrLeaderboard:
				handler_response.HandlerResponse(writer, hl.ResponseLeaderboard, http.StatusInternalServerError)
			}
			return
		}
		handler_response.HandlerResponse(writer, respLeaderboard, http.StatusOK)
	}
}
