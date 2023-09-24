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
	ErrTaskAlreadyExist    = errors.New("task already exist")
	ErrTaskDoesNotExist    = errors.New("task does not exist")
	ErrInvalidStatus       = errors.New("invalid status (status should be 'done' or 'not done'")
	ErrInvalidData         = errors.New("invalid data (should be like `2006-Jan-02`")
)
