package auth

import "errors"

var (
	ErrAlreadyExist    = errors.New("user with this email already exists")
	ErrIncorrectAction = errors.New("such action is not provided")
	ErrIncorrectCode   = errors.New("incorrect code")
	ErrValidSession    = errors.New("session is not valid")
	ErrRestoreUser     = errors.New("failed to restore user")
	ErrNewPassword     = errors.New("you did not specify a new password")
)
