package errors_custom

import "errors"

var (
	ErrIncorrectFormatData = errors.New("incorrect format of transmitted data")
	ErrIncorrectData       = errors.New("incorrect data")

	ErrSecurityData = errors.New("it was not possible to ensure the security of data storage")

	ErrWriteData = errors.New("failed to write your data")
	ErrToken     = errors.New("the token is invalid or has expired")

	ErrRecordNotFound = errors.New("record not found")

	ErrIncorrectPassword = errors.New("the password is incorrect")
	ErrNoExistUser       = errors.New("such user does not exist")
)
