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
	if body.Email != "" || body.NewPassword != "" {
		errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OriginalPassword))
		if errPassword != nil {
			return nil, errors_custom.ErrIncorrectPassword
		}
	}
	if body.NewPassword != "" {
		newPass, errHashPass := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
		if errHashPass != nil {
			return nil, errors_custom.ErrSecurityData
		}
		hashedNewPassword = string(newPass)
	}
	if body.Name != "" && body.OriginalPassword == "" && body.NewPassword == "" && body.Email == "" {
		user.Name = body.Name
	} else if body.Name == "" && body.OriginalPassword != "" && body.NewPassword != "" && body.Email == "" {
		user.Password = body.NewPassword
	} else if body.Name == "" && body.OriginalPassword != "" && body.NewPassword == "" && body.Email != "" {
		user.Error = body.Email
	} else if body.Name == "" && body.OriginalPassword != "" && body.NewPassword != "" && body.Email != "" {
		user.Email = body.Email
		user.Password = hashedNewPassword
	} else {
		return nil, ErrParamsUpdateUser
	}
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
