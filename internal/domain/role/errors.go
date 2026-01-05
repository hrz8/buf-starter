package role

import "errors"

var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role with this name already exists")
	ErrRoleInvalidName   = errors.New("invalid role name")
	ErrRoleInUse         = errors.New("role is in use and cannot be deleted")
)
