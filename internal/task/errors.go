package task

import "errors"

var (
	ErrIncorrectDeadline = errors.New("the deadline is incorrect")
	ErrCreateTask        = errors.New("failed to create task")
	ErrIncorrectTaskId   = errors.New("incorrect task_id")
	ErrTaskNotFound      = errors.New("such task not found")
	ErrDeleteTask        = errors.New("failed to delete task")
	ErrIncorrectPeriod   = errors.New("incorrect search period")
	ErrUpdateTask        = errors.New("failed to update task")
	ErrDoneUpdate        = errors.New("it is not possible to update a task that has already been completed")
	ErrImpossibleParams  = errors.New("impossible parameters were passed to update the task")
)
