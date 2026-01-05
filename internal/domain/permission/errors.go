package permission

import "errors"

var (
	ErrPermissionNotFound      = errors.New("permission not found")
	ErrPermissionAlreadyExists = errors.New("permission with this name already exists")
	ErrPermissionInvalidName   = errors.New("invalid permission name (must match ^[a-zA-Z0-9_:]+$)")
	ErrPermissionInvalidEffect = errors.New("invalid permission effect (must be 'allow' or 'deny')")
	ErrPermissionInUse         = errors.New("permission is in use and cannot be deleted")
)
