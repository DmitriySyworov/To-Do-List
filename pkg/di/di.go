package di

import "to-do-list/app/internal/model"

type IUserRepo interface {
	GetUserByIdUnscoped(uint) (*model.User, error)
	GetUserByEmailUnscoped(string) error
	GetUserByEmail(string) (*model.User, error)
	GetUserById(uint) (*model.User, error)
	GetUserByEmailDelete(string) (*model.User, error)
	CreateUser(*model.User) error
	RestoreDeleteUser(uint) error
	UpdateUser(*model.User) (*model.User, error)
}
