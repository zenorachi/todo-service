package entity

import "errors"

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrEmptyAuthHeader     = errors.New("empty authorization header")
	ErrInvalidAuthHeader   = errors.New("invalid authorization header")
	ErrUserAlreadyExists   = errors.New("user with such login/email already exists")
	ErrUserDoesNotExist    = errors.New("user does not exist")
	ErrIncorrectPassword   = errors.New("incorrect password")
	ErrSessionDoesNotExist = errors.New("session does not exist")
)
