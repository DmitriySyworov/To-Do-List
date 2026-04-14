package task

import (
	"time"
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/open_Db"

	"gorm.io/gorm"
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
func (r *RepositoryTask) UpdateTask(newTask *model.TaskForm, taskId, userId uint) error {
	res := r.DB.Where("task_id = ? AND user_id = ?", taskId, userId).Updates(&newTask)
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
func (r *RepositoryTask) GetAllTasks(userId uint, from, to time.Time) ([]model.TaskForm, error) {
	session := r.DB.Session(&gorm.Session{})
	var query *gorm.DB
	if !from.IsZero() && !to.IsZero() {
		query = session.Where("created_at >= ? AND created_at <= ? AND user_id = ?", from, to, userId)
	} else if !from.IsZero() && to.IsZero() {
		query = session.Where("created_at >= ? AND user_id = ?", from, userId)
	} else {
		query = session.Where("user_id = ?", userId)
	}
	var allTasks []model.TaskForm
	res := query.Find(&allTasks)
	if len(allTasks) == 0 {
		return nil, errors_custom.ErrRecordNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return allTasks, nil
}
