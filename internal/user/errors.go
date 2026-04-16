package user

import "errors"

var (
	ErrParamsUpdateUser = errors.New("incorrect parameters specified for user update")
	ErrUpdateUser       = errors.New("failed to update user")
	ErrDeleteUser       = errors.New("failed to delete user")
)
