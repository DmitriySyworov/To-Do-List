package task

import (
	"net/http"
	"to-do-list/app/configs"
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/handler_request"
	"to-do-list/app/pkg/handler_response"
	"to-do-list/app/pkg/middleware"
)

type HandlerTask struct {
	ResponseAllTasksPeriod
	model.TaskForm
	*HandlerTaskDep
}
type HandlerTaskDep struct {
	*ServiceTask
	*configs.Configs
}

func NewHandlerTask(router *http.ServeMux, dep *HandlerTaskDep) {
	task := HandlerTask{
		HandlerTaskDep: dep,
	}
	router.Handle("POST /user/task", middleware.IsAuthUserId(task.CreateTask(), dep.Configs.Secret))
	router.Handle("PATCH /user/task/{id}", middleware.IsAuthUserId(task.UpdateTask(), dep.Configs.Secret))
	router.Handle("GET /user/task/{id}", middleware.IsAuthUserId(task.GetTask(), dep.Configs.Secret))
	router.Handle("DELETE /user/task/{id}", middleware.IsAuthUserId(task.DeleteTask(), dep.Configs.Secret))
	router.Handle("GET /user/my-tasks", middleware.IsAuthUserId(task.GetAllTasks(), dep.Configs.Secret))
}

func (hl *HandlerTask) CreateTask() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handler_request.ResultRequest[RequestCreateTaskForm](request)
		if errBody != nil {
			hl.TaskForm.Error = errBody.Error()
			handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusBadRequest)
			return
		}
		userId, ok := request.Context().Value(middleware.KeyCtxUserId).(uint)
		if !ok {
			hl.TaskForm.Error = errors_custom.ErrToken.Error()
			handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusUnauthorized)
			return
		}
		taskForm, errCreate := hl.ServiceTask.CreateTask(body, userId)
		if errCreate != nil {
			hl.TaskForm.Error = errCreate.Error()
			switch errCreate {
			case errors_custom.ErrNoExistUser:
				handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusUnauthorized)
			case ErrIncorrectDeadline:
				handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusBadRequest)
			case ErrCreateTask:
				handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusInternalServerError)
			}
			return
		}
		handler_response.HandlerResponse(writer, taskForm, http.StatusCreated)
	}
}
func (hl *HandlerTask) UpdateTask() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := handler_request.ResultRequest[RequestUpdateTaskForm](request)
		if errBody != nil {
			hl.TaskForm.Error = errBody.Error()
			handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusBadRequest)
			return
		}
		if !body.StatusDone && body.Task == "" && body.Header == "" && body.Deadline == "" {
			hl.TaskForm.Error = errors_custom.ErrIncorrectData.Error()
			handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusBadRequest)
			return
		}
		taskIdStr := request.PathValue("id")
		if len(taskIdStr) != lengthTaskId {
			hl.TaskForm.Error = ErrIncorrectTaskId.Error()
			handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusBadRequest)
			return
		}
		userId, ok := request.Context().Value(middleware.KeyCtxUserId).(uint)
		if !ok {
			hl.TaskForm.Error = errors_custom.ErrToken.Error()
			handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusUnauthorized)
			return
		}
		newTaskForm, errUpdate := hl.ServiceTask.UpdateTask(body, userId, taskIdStr)
		if errUpdate != nil {
			hl.TaskForm.Error = ErrUpdateTask.Error()
			switch errUpdate {
			case ErrTaskNotFound:
				handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusNotFound)
			case ErrIncorrectDeadline, ErrIncorrectTaskId, ErrDoneUpdate, ErrImpossibleParams:
				handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusBadRequest)
			case ErrUpdateTask:
				handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusInternalServerError)
			}
			return
		}
		handler_response.HandlerResponse(writer, newTaskForm, http.StatusOK)
	}
}
func (hl *HandlerTask) GetTask() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId, ok := request.Context().Value(middleware.KeyCtxUserId).(uint)
		if !ok {
			hl.TaskForm.Error = errors_custom.ErrToken.Error()
			handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusUnauthorized)
			return
		}
		taskIdStr := request.PathValue("id")
		if len(taskIdStr) != lengthTaskId {
			hl.TaskForm.Error = ErrIncorrectTaskId.Error()
			handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusBadRequest)
			return
		}
		taskForm, errGet := hl.ServiceTask.GetTask(taskIdStr, userId)
		if errGet != nil {
			hl.TaskForm.Error = errGet.Error()
			switch errGet {
			case ErrTaskNotFound:
				handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusNotFound)
			case ErrIncorrectTaskId:
				handler_response.HandlerResponse(writer, hl.TaskForm, http.StatusBadRequest)
			}
			return
		}
		handler_response.HandlerResponse(writer, taskForm, http.StatusOK)
	}
}
func (hl *HandlerTask) DeleteTask() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId, ok := request.Context().Value(middleware.KeyCtxUserId).(uint)
		if !ok {
			http.Error(writer, errors_custom.ErrToken.Error(), http.StatusUnauthorized)
			return
		}
		taskIdStr := request.PathValue("id")
		if len(taskIdStr) != lengthTaskId {
			http.Error(writer, ErrIncorrectTaskId.Error(), http.StatusBadRequest)
			return
		}
		errDel := hl.ServiceTask.DeleteTask(taskIdStr, userId)
		if errDel != nil {
			switch errDel {
			case ErrTaskNotFound:
				http.Error(writer, ErrTaskNotFound.Error(), http.StatusNotFound)
			case ErrIncorrectTaskId:
				http.Error(writer, ErrIncorrectTaskId.Error(), http.StatusBadRequest)
			case ErrDeleteTask:
				http.Error(writer, ErrDeleteTask.Error(), http.StatusInternalServerError)
			}
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}
func (hl *HandlerTask) GetAllTasks() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId, ok := request.Context().Value(middleware.KeyCtxUserId).(uint)
		if !ok {
			hl.ResponseAllTasksPeriod.Error = errors_custom.ErrToken.Error()
			handler_response.HandlerResponse(writer, hl.ResponseAllTasksPeriod, http.StatusUnauthorized)
			return
		}
		from := request.URL.Query().Get("from")
		to := request.URL.Query().Get("to")
		respAllTasks, errGetAll := hl.ServiceTask.GetAllTasks(userId, from, to)
		if errGetAll != nil {
			hl.ResponseAllTasksPeriod.Error = errors_custom.ErrRecordNotFound.Error()
			handler_response.HandlerResponse(writer, hl.ResponseAllTasksPeriod, http.StatusNotFound)
			return
		}
		handler_response.HandlerResponse(writer, respAllTasks, http.StatusOK)
	}
}
