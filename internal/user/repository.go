package user

import (
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/open_Db"
)

type RepositoryUsers struct {
	*open_Db.OpenPostgres
}

func NewRepositoryUsers(postgres *open_Db.OpenPostgres) *RepositoryUsers {
	return &RepositoryUsers{
		OpenPostgres: postgres,
	}
}
func (r *RepositoryUsers) DeleteUser(userId uint) error {
	res := r.DB.Where("user_id = ?", userId).Delete(&model.User{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryUsers) GetUserById(userId uint) (*model.User, error) {
	var user model.User
	res := r.DB.Where("user_id = ?", userId).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (r *RepositoryUsers) UpdateUser(user *model.User) (*model.User, error) {
	res := r.DB.Where("user_id = ?", user.UserId).Updates(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
func (r *RepositoryUsers) CreateUser(user *model.User) error {
	res := r.DB.Create(&user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryUsers) RestoreDeleteUser(userId uint) error {
	res := r.DB.Model(&model.User{}).Where("user_id = ?", userId).Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryUsers) GetUserByIdUnscoped(userId uint) (*model.User, error) {
	var user model.User
	res := r.DB.Unscoped().Where("user_id = ?", userId).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
func (r *RepositoryUsers) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	res := r.DB.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
func (r *RepositoryUsers) GetUserByEmailUnscoped(email string) error {
	var user model.User
	res := r.DB.Unscoped().Where("email = ?", email).First(&user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryUsers) GetUserByEmailDelete(email string) (*model.User, error) {
	var user model.User
	res := r.DB.Unscoped().Where("email = ? AND deleted_at not null", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
