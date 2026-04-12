package task

import (
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/open_Db"

	"gorm.io/gorm/clause"
)

type RepositoryTask struct {
	*open_Db.OpenPostgres
}

func NewRepositoryTask(postgres *open_Db.OpenPostgres) *RepositoryTask {
	return &RepositoryTask{
		OpenPostgres: postgres,
	}
}
func (r *RepositoryTask) CreateTask(taskForm *model.TaskForm) error {
	res := r.DB.Create(&taskForm)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryTask) GetTask(taskId, userId uint) (*model.TaskForm, error) {
	var taskForm model.TaskForm
	res := r.DB.Where("user_id = ? AND task_id = ?", userId, taskId).First(&taskForm)
	if res.Error != nil {
		return nil, res.Error
	}
	return &taskForm, nil
}
func (r *RepositoryTask) DeleteTask(taskId, userId uint) error {
	res := r.DB.Unscoped().Where("user_id = ? AND task_id = ?", userId, taskId).Delete(&model.TaskForm{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryTask) GetAllTasks(userId uint, format string) ([]model.TaskForm, error) {
	var allTasks []model.TaskForm
	res := r.DB.
		Model(&model.TaskForm{}).
		Select("to_char(created_at, ?) as period", format).
		Where("user_id = ?", userId).
		Group("period").
		Scan(&allTasks)
	if len(allTasks) == 0 {
		return nil, errors_custom.ErrRecordNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return allTasks, nil
}
func (r *RepositoryTask) UpdateTask(newTask *model.TaskForm, taskId, userId uint) error {
	res := r.DB.Clauses(clause.Clause{}).Where("task_id = ? AND user_id = ?", taskId, userId).Updates(&newTask)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
