package iam_mapper

import (
	"context"

	"github.com/hrz8/altalune/internal/domain/permission"
	"github.com/hrz8/altalune/internal/domain/role"
)

// Repository defines the interface for IAM mapping operations
type Repository interface {
	// User-Role Mappings
	AssignUserRoles(ctx context.Context, userID int64, roleIDs []int64) error
	RemoveUserRoles(ctx context.Context, userID int64, roleIDs []int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]*role.Role, error)

	// Role-Permission Mappings
	AssignRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
	RemoveRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]*permission.Permission, error)

	// User-Permission Mappings (Direct Assignments)
	AssignUserPermissions(ctx context.Context, userID int64, permissionIDs []int64) error
	RemoveUserPermissions(ctx context.Context, userID int64, permissionIDs []int64) error
	GetUserPermissions(ctx context.Context, userID int64) ([]*permission.Permission, error)

	// Project Members
	AssignProjectMembers(ctx context.Context, projectID int64, members []ProjectMemberInput) error
	RemoveProjectMembers(ctx context.Context, projectID int64, userIDs []int64) error
	GetProjectMembers(ctx context.Context, projectID int64) ([]*ProjectMemberWithUser, error)

	// User Projects (reverse lookup - projects a user belongs to)
	GetUserProjects(ctx context.Context, userID int64) ([]*UserProjectMembership, error)
}
