package user

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user with this email already exists")
	ErrUserInvalidEmail     = errors.New("invalid email format")
	ErrUserAlreadyActive    = errors.New("user is already active")
	ErrUserAlreadyInactive  = errors.New("user is already inactive")
	ErrUserCannotDeleteSelf = errors.New("cannot delete your own user account")
)
