package domain

import "errors"

var (
	ErrBadRequest         = errors.New("invalid request payload")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrTokenExpired   = errors.New("token has expired")
	ErrTokenInvalid   = errors.New("token is invalid")
	ErrTokenMalformed = errors.New("token format is incorrect")

	ErrInternal = errors.New("internal server error")
	ErrTimeout  = errors.New("operation timed out")
	ErrDatabase = errors.New("database provider error")
)
