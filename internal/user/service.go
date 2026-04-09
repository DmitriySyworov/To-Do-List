package user

import (
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/errors_custom"

	"golang.org/x/crypto/bcrypt"
)

type ServiceUser struct {
	*RepositoryUsers
}

func NewServiceUsers(repo *RepositoryUsers) *ServiceUser {
	return &ServiceUser{
		RepositoryUsers: repo,
	}
}

func (s *ServiceUser) UpdateUser(body *RequestUpdateUser, userId uint) (*model.User, error) {
	user, errGet := s.GetUserById(userId)
	if errGet != nil {
		return nil, errors_custom.ErrRecordNotFound
	}
	hashedNewPassword := ""
	if body.NewPassword != "" {
		errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OriginalPassword))
		if errPassword != nil {
			return nil, errors_custom.ErrIncorrectPassword
		}
		newPass, errHashPass := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
		if errHashPass != nil {
			return nil, errors_custom.ErrSecurityData
		}
		hashedNewPassword = string(newPass)
	}
	user.Name = body.Name
	user.Email = body.Email
	user.Password = hashedNewPassword
	resUser, errUpdate := s.RepositoryUsers.UpdateUser(user)
	if errUpdate != nil {
		return nil, ErrUpdateUser
	}
	return resUser, nil
}
func (s *ServiceUser) DeleteUser(password string, userId uint) error {
	user, errGet := s.GetUserById(userId)
	if errGet != nil {
		return errors_custom.ErrRecordNotFound
	}
	errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errPassword != nil {
		return errors_custom.ErrIncorrectPassword
	}
	errDel := s.RepositoryUsers.DeleteUser(userId)

	if errDel != nil {
		return ErrDeleteUser
	}
	return nil
}
