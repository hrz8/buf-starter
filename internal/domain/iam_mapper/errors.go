package iam_mapper

import "errors"

// Domain-specific errors for IAM Mapper operations
var (
	// ErrMappingNotFound indicates a mapping record was not found
	ErrMappingNotFound = errors.New("mapping not found")

	// ErrMappingAlreadyExists indicates a mapping already exists (can be treated as warning with ON CONFLICT)
	ErrMappingAlreadyExists = errors.New("mapping already exists")

	// ErrInvalidProjectRole indicates an invalid project role was provided
	ErrInvalidProjectRole = errors.New("invalid project role")

	// ErrCannotRemoveLastOwner indicates attempt to remove the last owner from a project
	ErrCannotRemoveLastOwner = errors.New("cannot remove last owner from project")

	// ErrUserNotFound indicates user was not found during mapping operation
	ErrUserNotFound = errors.New("user not found for mapping")

	// ErrRoleNotFound indicates role was not found during mapping operation
	ErrRoleNotFound = errors.New("role not found for mapping")

	// ErrPermissionNotFound indicates permission was not found during mapping operation
	ErrPermissionNotFound = errors.New("permission not found for mapping")

	// ErrProjectNotFound indicates project was not found during mapping operation
	ErrProjectNotFound = errors.New("project not found for mapping")
)

// Error codes for IAM Mapper domain (608XX range)
const (
	CodeMappingNotFound        = 60800
	CodeMappingAlreadyExists   = 60801
	CodeInvalidProjectRole     = 60802
	CodeCannotRemoveLastOwner  = 60803
	CodeUserNotFoundForMapping = 60804
	CodeRoleNotFoundForMapping = 60805
	CodePermissionNotFound     = 60806
	CodeProjectNotFound        = 60807
)
