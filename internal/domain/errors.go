package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user doesn't exists")
	ErrUserAlreadyExists = errors.New("user with such email already exists")

	ErrProjectNotFound = errors.New("project doesn't exists")

	ErrVerificationCodeInvalid = errors.New("verification code is invalid")
)
