package project_member

import "errors"

var (
	ErrProjectMemberNotFound      = errors.New("project member not found")
	ErrProjectMemberAlreadyExists = errors.New("user is already a member of this project")
	ErrInvalidRole                = errors.New("invalid role specified")
	ErrCannotModifyOwnerRole      = errors.New("owner role is reserved for superadmin")
	ErrCannotRemoveLastOwner      = errors.New("cannot remove the last owner from project")
	ErrInsufficientPermissions    = errors.New("insufficient permissions for this operation")
)
