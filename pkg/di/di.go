package di

import "to-do-list/app/internal/models"

type IUserRepo interface {
	GetUserByIdUnscoped(uint) (*models.Users, error)
	GetUserByEmailUnscoped(string) error
	GetUserByEmail(string) (*models.Users, error)
	GetUserByEmailDelete(string) (*models.Users, error)
	CreateUser(*models.Users) error
	RestoreUser(uint) error
}
