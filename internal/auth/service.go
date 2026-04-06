package auth

import (
	"to-do-list/app/configs"
	"to-do-list/app/internal/models"
	"to-do-list/app/pkg/di"
	"to-do-list/app/pkg/errorsCust"
	"to-do-list/app/pkg/generateRand"
	"to-do-list/app/pkg/jwt"
	"to-do-list/app/pkg/sendLetter"

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
		return nil, errorsCust.ErrSecurityData
	}
	var userId uint
	for {
		userId = generateRand.GenerateNumbers(lengthUserId)
		if _, errId := s.IUserRepo.GetUserByIdUnscoped(userId); errId != nil {
			break
		}
	}
	respAuth, errAuth := s.helperAuth(&models.TempUser{
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
		return nil, errorsCust.ErrRecordNotFound
	}
	errCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if errCompare != nil {
		return nil, ErrIncorrectPassword
	}
	respAuth, errAuth := s.helperAuth(&models.TempUser{
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
func (s *ServiceAuth) Restore(body *RequestLoginAndRestore) (*ResponseAuth, error) {
	deleteUser, errGetDelete := s.IUserRepo.GetUserByEmailDelete(body.Email)
	if errGetDelete != nil {
		return nil, errorsCust.ErrRecordNotFound
	}
	respAuth, errAuth := s.helperAuth(&models.TempUser{
		Name:     deleteUser.Name,
		Email:    deleteUser.Email,
		Password: deleteUser.Password,
		UserId:   deleteUser.UserId,
	})
	if errAuth != nil {
		return nil, errAuth
	}
	return respAuth, nil
}
func (s *ServiceAuth) helperAuth(tempUser *models.TempUser) (*ResponseAuth, error) {
	idHash := generateRand.GenerateNumbers(lengthHashId)
	sessionId := generateRand.GenerateStr(lengthTempCodeAndSession)
	tempCode := generateRand.GenerateNumbers(lengthTempCodeAndSession)
	j := jwt.NewJWT(s.Secret)
	token, errCreateTempJwt := j.CreateTemporaryJWT(float64(idHash), sessionId)
	if errCreateTempJwt != nil {
		return nil, errCreateTempJwt
	}
	errCreateTempUser := s.Repo.CreateTempUser(tempUser, idHash)
	if errCreateTempUser != nil {
		return nil, errorsCust.ErrWriteData
	}
	errCreateSession := s.Repo.CreateSession(&models.Session{
		SessionId:     sessionId,
		TemporaryCode: tempCode,
	}, idHash)
	if errCreateSession != nil {
		return nil, errorsCust.ErrWriteData
	}
	errSend := sendLetter.SendByEmail(tempUser.Email, tempCode, s.Configs.SendEmail)
	if errSend != nil {
		return nil, errSend
	}
	return &ResponseAuth{
		Message:   "We sent a letter with a verification code to the specified email: " + tempUser.Email,
		SessionId: sessionId,
		JWT:       token,
	}, nil
}
func (s *ServiceAuth) Confirm(hashId, tempCode uint, sessionId, action string) (*ResponseConfirm, error) {
	tempUser, ErrGetTempUser := s.Repo.GetTempUser(hashId)
	if ErrGetTempUser != nil {
		return nil, errorsCust.ErrRecordNotFound
	}
	session, errSession := s.Repo.GetSession(hashId)
	if errSession != nil {
		return nil, errorsCust.ErrRecordNotFound
	}
	if session.TemporaryCode != tempCode {
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
		errCreate := s.IUserRepo.CreateUser(&models.Users{
			Name:     tempUser.Name,
			Email:    tempUser.Email,
			Password: tempUser.Password,
			UserId:   tempUser.UserId,
		})
		if errCreate != nil {
			return nil, errCreate
		}
		token, errJWT := jwt.NewJWT(s.Secret).CreateJWT(float64(tempUser.UserId))
		if errJWT != nil {
			return nil, errJWT
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	case "login":
		user, errGet := s.IUserRepo.GetUserByEmail(tempUser.Email)
		if errGet != nil {
			return nil, errorsCust.ErrRecordNotFound
		}
		token, errJWT := jwt.NewJWT(s.Secret).CreateJWT(float64(user.UserId))
		if errJWT != nil {
			return nil, errJWT
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	case "restore":
		deleteUser, errGetDelete := s.IUserRepo.GetUserByEmailDelete(tempUser.Email)
		if errGetDelete != nil {
			return nil, errorsCust.ErrRecordNotFound
		}
		errRestore := s.IUserRepo.RestoreUser(deleteUser.UserId)
		if errRestore != nil {
			return nil, ErrRestoreUser
		}
		token, errJWT := jwt.NewJWT(s.Secret).CreateJWT(float64(deleteUser.UserId))
		if errJWT != nil {
			return nil, errJWT
		}
		return &ResponseConfirm{
			JWT: token,
		}, nil
	default:
		return nil, ErrIncorrectAction
	}
}
