package task

import (
	"strconv"
	"time"
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/di"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/generate_rand"
)

type ServiceTask struct {
	*RepositoryTask
	*ServiceTaskDep
}
type ServiceTaskDep struct {
	di.IUserRepo
}

func NewServiceTask(repo *RepositoryTask, dep *ServiceTaskDep) *ServiceTask {
	return &ServiceTask{
		RepositoryTask: repo,
		ServiceTaskDep: dep,
	}
}

const lengthTaskId = 7

func (s *ServiceTask) CreateTask(body *RequestCreateTaskForm, userId uint) (*model.TaskForm, error) {
	_, errUser := s.IUserRepo.GetUserById(userId)
	if errUser != nil {
		return nil, errors_custom.ErrNoExistUser
	}
	var deadline time.Time
	if body.Deadline != "" {
		date, errDate := time.Parse(time.DateOnly, body.Deadline)
		if errDate != nil {
			return nil, ErrIncorrectDeadline
		}
		deadline = date
	}
	taskForm := &model.TaskForm{
		Header:     body.Header,
		Task:       body.Task,
		Deadline:   deadline,
		StatusDone: false,
		TaskId:     generate_rand.GenerateNumbers(lengthTaskId),
		UserId:     userId,
	}
	errCreate := s.RepositoryTask.Create(taskForm)
	if errCreate != nil {
		return nil, ErrCreateTask
	}
	return taskForm, nil
}
func (s *ServiceTask) parseIdAndGetTask(taskIdStr string, userId uint) (*model.TaskForm, uint, error) {
	taskId, errParseId := strconv.Atoi(taskIdStr)
	if errParseId != nil {
		return nil, 0, ErrIncorrectTaskId
	}
	taskForm, errGet := s.RepositoryTask.GetTask(uint(taskId), userId)
	if errGet != nil {
		return nil, 0, ErrTaskNotFound
	}
	return taskForm, uint(taskId), nil
}
func (s *ServiceTask) UpdateTask(body *RequestUpdateTaskForm, userId uint, taskIdStr string) (*model.TaskForm, error) {
	_, taskId, errParseAndGet := s.parseIdAndGetTask(taskIdStr, userId)
	if errParseAndGet != nil {
		return nil, errParseAndGet
	}
	var deadline time.Time
	if body.Deadline != "" {
		date, errDate := time.Parse(time.DateOnly, body.Deadline)
		if errDate != nil {
			return nil, ErrIncorrectDeadline
		}
		deadline = date
	}
	newTaskForm := &model.TaskForm{
		Header:     body.Header,
		Task:       body.Task,
		Deadline:   deadline,
		StatusDone: body.StatusDone,
	}
	errUpdate := s.RepositoryTask.UpdateTask(newTaskForm, taskId, userId)
	if errUpdate != nil {
		return nil, ErrUpdateTask
	}
	return newTaskForm, nil
}
func (s *ServiceTask) GetTask(taskIdStr string, userId uint) (*model.TaskForm, error) {
	taskForm, _, errParseAndGet := s.parseIdAndGetTask(taskIdStr, userId)
	if errParseAndGet != nil {
		return nil, errParseAndGet
	}
	return taskForm, nil
}
func (s *ServiceTask) DeleteTask(taskIdStr string, userId uint) error {
	_, taskId, errParseAndGet := s.parseIdAndGetTask(taskIdStr, userId)
	if errParseAndGet != nil {
		return errParseAndGet
	}
	errDel := s.RepositoryTask.DeleteTask(taskId, userId)
	if errDel != nil {
		return ErrDeleteTask
	}
	return nil
}
func (s *ServiceTask) GetAllTasks(userId uint, period string) (*ResponseAllTasksPeriod, error) {
	format := ""
	if period == "year" {
		format = "YYYY"
	} else if period == "month" {
		format = "YYYY-MM"
	} else if period == "day" {
		format = "YYYY-MM-DD"
	}
	AllTasks, errGet := s.RepositoryTask.GetAllTasks(userId, format)
	if errGet != nil {
		return nil, errors_custom.ErrRecordNotFound
	}
	sliceActiveTasks := make([]model.TaskForm, len(AllTasks))
	sliceDoneTasks := make([]model.TaskForm, len(AllTasks))
	for _, task := range AllTasks {
		if task.StatusDone {
			sliceDoneTasks = append(sliceDoneTasks, task)
		} else {
			sliceActiveTasks = append(sliceActiveTasks, task)
		}
	}
	return &ResponseAllTasksPeriod{
		ActiveTasks: sliceActiveTasks,
		DoneTasks:   sliceDoneTasks,
	}, nil
}
