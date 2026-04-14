package task

import (
	"strconv"
	"time"
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/di"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/event_bus"
	"to-do-list/app/pkg/generate_rand"
)

type ServiceTask struct {
	*RepositoryTask
	*ServiceTaskDep
}
type ServiceTaskDep struct {
	di.IUserRepo
	*event_bus.EventBus
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
	taskForm := model.TaskForm{
		Header:     body.Header,
		Task:       body.Task,
		Deadline:   deadline,
		StatusDone: false,
		TaskId:     generate_rand.GenerateNumbers(lengthTaskId),
		UserId:     userId,
	}
	errCreate := s.RepositoryTask.CreateTask(&taskForm)
	if errCreate != nil {
		return nil, ErrCreateTask
	}
	go s.EventBus.Publish(&event_bus.Event{
		Name: event_bus.EventCreateTask,
		Data: userId,
	})
	return &taskForm, nil
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
	taskForm, taskId, errParseAndGet := s.parseIdAndGetTask(taskIdStr, userId)
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
	if taskForm.StatusDone {
		return nil, ErrDoneUpdate
	} else if body.Header != "" && body.Task == "" && body.Deadline == "" && !body.StatusDone {
		taskForm.Header = body.Header
	} else if body.Header == "" && body.Task != "" && body.Deadline == "" && !body.StatusDone {
		taskForm.Task = body.Task
	} else if body.Header == "" && body.Task == "" && body.Deadline != "" && !body.StatusDone {
		taskForm.Deadline = deadline
	} else if body.Header == "" && body.Task != "" && body.Deadline != "" && !body.StatusDone {
		taskForm.Task = body.Task
		taskForm.Deadline = deadline
	} else if body.Header != "" && body.Task == "" && body.Deadline != "" && !body.StatusDone {
		taskForm.Header = body.Header
		taskForm.Deadline = deadline
	} else if body.Header != "" && body.Task != "" && body.Deadline != "" && !body.StatusDone {
		taskForm.Header = body.Header
		taskForm.Task = body.Task
		taskForm.Deadline = deadline
	} else if body.Header != "" && body.Task != "" && body.Deadline == "" && !body.StatusDone {
		taskForm.Header = body.Header
		taskForm.Task = body.Task
	} else if body.Header == "" && body.Task == "" && body.Deadline == "" && body.StatusDone {
		taskForm.StatusDone = body.StatusDone
	} else {
		return nil, ErrImpossibleParams
	}
	errUpdate := s.RepositoryTask.UpdateTask(taskForm, taskId, userId)
	if errUpdate != nil {
		return nil, ErrUpdateTask
	}
	if body.StatusDone {
		go s.EventBus.Publish(&event_bus.Event{
			Name: event_bus.EventDoneTask,
			Data: userId,
		})
	}
	return taskForm, nil
}
func (s *ServiceTask) GetTask(taskIdStr string, userId uint) (*model.TaskForm, error) {
	taskForm, _, errParseAndGet := s.parseIdAndGetTask(taskIdStr, userId)
	if errParseAndGet != nil {
		return nil, errParseAndGet
	}
	return taskForm, nil
}
func (s *ServiceTask) DeleteTask(taskIdStr string, userId uint) error {
	taskForm, taskId, errParseAndGet := s.parseIdAndGetTask(taskIdStr, userId)
	if errParseAndGet != nil {
		return errParseAndGet
	}
	errDel := s.RepositoryTask.DeleteTask(taskId, userId)
	if errDel != nil {
		return ErrDeleteTask
	}
	if taskForm.StatusDone {
		go s.EventBus.Publish(&event_bus.Event{
			Name: event_bus.EventDeleteDoneTask,
			Data: userId,
		})
	} else {
		go s.EventBus.Publish(&event_bus.Event{
			Name: event_bus.EventDeleteActiveTask,
			Data: userId,
		})
	}
	return nil
}
func (s *ServiceTask) GetAllTasks(userId uint, fromStr, toStr string) (*ResponseAllTasksPeriod, error) {
	var from, to time.Time
	if fromStr != "" {
		date, errDate := time.Parse(time.DateOnly, fromStr)
		if errDate != nil {
			return nil, ErrIncorrectPeriod
		}
		from = date
	}
	if toStr != "" {
		date, errDate := time.Parse(time.DateOnly, toStr)
		if errDate != nil {
			return nil, ErrIncorrectPeriod
		}
		to = date
	}
	allTasks, errGetAll := s.RepositoryTask.GetAllTasks(userId, from, to)
	if errGetAll != nil {
		return nil, errors_custom.ErrRecordNotFound
	}
	var sliceActiveTasks, sliceDoneTasks []model.TaskForm
	for _, task := range allTasks {
		if task.StatusDone && task.ID != 0 {
			sliceDoneTasks = append(sliceDoneTasks, task)
		} else if !task.StatusDone && task.ID != 0 {
			sliceActiveTasks = append(sliceActiveTasks, task)
		}
	}
	return &ResponseAllTasksPeriod{
		ActiveTasks: sliceActiveTasks,
		DoneTasks:   sliceDoneTasks,
	}, nil
}
