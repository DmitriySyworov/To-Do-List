package task

import (
	"to-do-list/app/internal/model"
)

type RequestCreateTaskForm struct {
	Header   string `json:"header"`
	Task     string `json:"task" validate:"required"`
	Deadline string `json:"deadline" validate:"omitempty,datetime=2006-01-02"`
}
type RequestUpdateTaskForm struct {
	Header     string `json:"header" validate:"excluded_with=StatusDone"`
	Task       string `json:"task" validate:"excluded_with=StatusDone"`
	Deadline   string `json:"deadline" validate:"omitempty,datetime=2006-01-02,excluded_with=StatusDone"`
	StatusDone bool   `json:"status_done"`
}
type ResponseAllTasksPeriod struct {
	ActiveTasks []model.TaskForm `json:"active_tasks"`
	DoneTasks   []model.TaskForm `json:"done_tasks"`
	Error       string           `json:"error"`
}
