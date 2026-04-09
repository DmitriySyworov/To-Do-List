package user

import "errors"

var (
	ErrUpdateUser = errors.New("failed to update user")
	ErrDeleteUser = errors.New("failed to delete user")
)
