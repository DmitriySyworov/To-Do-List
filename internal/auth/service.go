package auth

import (
	"log"
	"to-do-list/app/configs"
	"to-do-list/app/internal/model"
	"to-do-list/app/pkg/di"
	"to-do-list/app/pkg/errors_custom"
	"to-do-list/app/pkg/generate_rand"
	"to-do-list/app/pkg/jwt"
	"to-do-list/app/pkg/send_letter"

	"golang.org/x/crypto/bcrypt"
)

type ServiceAuth struct {
	Repo *RepositoryAuth
	*ServiceAuthDep
}
type ServiceAuthDep struct {
	di.IUserRepo
	*configs.Configs
}

func NewServiceAuth(repo *RepositoryAuth, dep *ServiceAuthDep) *ServiceAuth {
	return &ServiceAuth{
		Repo:           repo,
		ServiceAuthDep: dep,
	}
}

const (
	lengthUserId             = 11
	lengthHashId             = 9
	lengthTempCodeAndSession = 6
)

func (s *ServiceAuth) Register(body *RequestRegister) (*ResponseAuth, error) {
	errGet := s.IUserRepo.GetUserByEmailUnscoped(body.Email)
	if errGet == nil {
		return nil, ErrAlreadyExist
	}
	hashPassword, errCrypt := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if errCrypt != nil {
		return nil, errors_custom.ErrSecurityData
	}
	var userId uint
	for {
		userId = generate_rand.GenerateNumbers(lengthUserId)
		if _, errId := s.IUserRepo.GetUserByIdUnscoped(userId); errId != nil {
			break
		}
	}
	respAuth, errAuth := s.helperAuth(&model.TempUser{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashPassword),
		UserId:   userId,
	})
	if errAuth != nil {
		return nil, errAuth
	}
	return respAuth, nil
}
func (s *ServiceAuth) Login(body *RequestLoginAndRestore) (*ResponseAuth, error) {
	user, errGet := s.IUserRepo.GetUserByEmail(body.Email)
	if errGet != nil {
		return nil, errors_custom.ErrRecordNotFound
	}
	log.Println(user.Password, body.Password)

	errCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if errCompare != nil {
		return nil, errors_custom.ErrIncorrectPassword
	}
	respAuth, errAuth := s.helperAuth(&model.TempUser{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		UserId:   user.UserId,
	})
	if errAuth != nil {
		return nil, errAuth
	}
	return respAuth, nil
}
func (s *ServiceAuth) Restore(body *RequestLoginAndRestore, action string) (*ResponseAuth, error) {
	switch action {
	case "recoverDelete":
		deleteUser, errGetDelete := s.IUserRepo.GetUserByEmailDelete(body.Email)
		if errGetDelete != nil {
			return nil, errors_custom.ErrRecordNotFound
		}
		respAuth, errAuth := s.helperAuth(&model.TempUser{
			Name:     deleteUser.Name,
			Email:    deleteUser.Email,
			Password: deleteUser.Password,
			UserId:   deleteUser.UserId,
		})
		if errAuth != nil {
			return nil, errAuth
		}
		return respAuth, nil
	case "recoverLogin":
		user, errGet := s.IUserRepo.GetUserByEmail(body.Email)
		if errGet != nil {
			return nil, errors_custom.ErrRecordNotFound
		}
		respAuth, errAuth := s.helperAuth(&model.TempUser{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
			UserId:   user.UserId,
		})
		if errAuth != nil {
			return nil, errAuth
		}
		return respAuth, nil
	default:
		return nil, ErrIncorrectAction
	}
}
func (s *ServiceAuth) helperAuth(tempUser *model.TempUser) (*ResponseAuth, error) {
	idHash := generate_rand.GenerateNumbers(lengthHashId)
	sessionId := generate_rand.GenerateStr(lengthTempCodeAndSession)
	tempCode := generate_rand.GenerateNumbers(lengthTempCodeAndSession)
	j := jwt.NewJWT(s.Secret)
	token, errCreateTempJwt := j.CreateTemporaryJWT(float64(idHash), sessionId)
	if errCreateTempJwt != nil {
		return nil, errCreateTempJwt
	}
	errCreateTempUser := s.Repo.CreateTempUser(tempUser, idHash)
	if errCreateTempUser != nil {
		return nil, errors_custom.ErrWriteData
	}
	errCreateSession := s.Repo.CreateSession(&model.Session{
		SessionId:     sessionId,
		TemporaryCode: tempCode,
	}, idHash)
	if errCreateSession != nil {
		return nil, errors_custom.ErrWriteData
	}
	errSend := send_letter.SendByEmail(tempUser.Email, tempCode, s.Configs.SendEmail)
	if errSend != nil {
		return nil, errSend
	}
	return &ResponseAuth{
		Message: "We sent a letter with a verification code to the specified email: " + tempUser.Email,
		JWT:     token,
	}, nil
}
func (s *ServiceAuth) Confirm(body *RequestConfirm, hashId uint, sessionId, action string) (*ResponseConfirm, error) {
	tempUser, ErrGetTempUser := s.Repo.GetTempUser(hashId)
	if ErrGetTempUser != nil {
		return nil, errors_custom.ErrRecordNotFound
	}
	session, errSession := s.Repo.GetSession(hashId)
	if errSession != nil {
		return nil, errors_custom.ErrRecordNotFound
	}
	if session.TemporaryCode != body.TempCode {
		return nil, ErrIncorrectCode
	}
	if session.SessionId != sessionId {
		return nil, ErrValidSession
	}
	switch action {
	case "register":
		errGet := s.IUserRepo.GetUserByEmailUnscoped(tempUser.Email)
		if errGet == nil {
			return nil, ErrAlreadyExist
		}
		token, errJWT := jwt.NewJWT(s.Secret).CreateJWT(float64(tempUser.UserId))
		if errJWT != nil {
			return nil, errJWT
		}
		errCreate := s.IUserRepo.CreateUser(&model.User{
			Name:     tempUser.Name,
			Email:    tempUser.Email,
			Password: tempUser.Password,
			UserId:   tempUser.UserId,
		})
		if errCreate != nil {
			return nil, errCreate
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	case "login":
		user, errGet := s.IUserRepo.GetUserByEmail(tempUser.Email)
		if errGet != nil {
			return nil, errors_custom.ErrRecordNotFound
		}
		token, errJWT := jwt.NewJWT(s.Secret).CreateJWT(float64(user.UserId))
		if errJWT != nil {
			return nil, errJWT
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	case "recoverDelete":
		deleteUser, errGetDelete := s.IUserRepo.GetUserByEmailDelete(tempUser.Email)
		if errGetDelete != nil {
			return nil, errors_custom.ErrRecordNotFound
		}
		token, errJWT := jwt.NewJWT(s.Secret).CreateJWT(float64(deleteUser.UserId))
		if errJWT != nil {
			return nil, errJWT
		}
		errRestore := s.IUserRepo.RestoreDeleteUser(deleteUser.UserId)
		if errRestore != nil {
			return nil, ErrRestoreUser
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	case "recoverLogin":
		user, errGet := s.IUserRepo.GetUserByEmail(tempUser.Email)
		if errGet != nil {
			return nil, errors_custom.ErrRecordNotFound
		}
		newHashPass, errCrypt := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
		if errCrypt != nil {
			return nil, errors_custom.ErrSecurityData
		}
		token, errJWT := jwt.NewJWT(s.Secret).CreateJWT(float64(user.UserId))
		if errJWT != nil {
			return nil, errJWT
		}
		user.Password = string(newHashPass)
		_, errRestore := s.IUserRepo.UpdateUser(user)
		if errRestore != nil {
			return nil, ErrRestoreUser
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	default:
		return nil, ErrIncorrectAction
	}
}
