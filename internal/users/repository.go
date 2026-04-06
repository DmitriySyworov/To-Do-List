package users

import (
	"to-do-list/app/internal/models"
	"to-do-list/app/pkg/openDb"
)

type RepositoryUsers struct {
	*openDb.OpenPostgres
}

func NewRepositoryUsers(postgres *openDb.OpenPostgres) *RepositoryUsers {
	return &RepositoryUsers{
		OpenPostgres: postgres,
	}
}
func (r *RepositoryUsers) CreateUser(user *models.Users) error {
	res := r.DB.Create(&user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryUsers) RestoreUser(userId uint) error {
	res := r.DB.Unscoped().Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryUsers) GetUserByIdUnscoped(userId uint) (*models.Users, error) {
	var user models.Users
	res := r.DB.Unscoped().Where("user_id = ?", userId).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
func (r *RepositoryUsers) GetUserByEmail(email string) (*models.Users, error) {
	var user models.Users
	res := r.DB.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
func (r *RepositoryUsers) GetUserByEmailUnscoped(email string) error {
	var user models.Users
	res := r.DB.Unscoped().Where("email = ?", email).First(&user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *RepositoryUsers) GetUserByEmailDelete(email string) (*models.Users, error) {
	var user models.Users
	res := r.DB.Unscoped().Where("email = ? AND deleted_at not null", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
