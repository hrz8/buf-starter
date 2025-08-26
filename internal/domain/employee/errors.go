package employee

import "errors"

var (
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrEmployeeAlreadyExists = errors.New("employee with this email already exists")
)
