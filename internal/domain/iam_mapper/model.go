package iam_mapper

import (
	"time"

	"github.com/hrz8/altalune/internal/domain/permission"
	"github.com/hrz8/altalune/internal/domain/role"
	"github.com/hrz8/altalune/internal/domain/user"
)

// UserRole represents the junction table for user-role assignments
type UserRole struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	RoleID    int64     `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
}

// RolePermission represents the junction table for role-permission assignments
type RolePermission struct {
	ID           int64     `db:"id"`
	RoleID       int64     `db:"role_id"`
	PermissionID int64     `db:"permission_id"`
	CreatedAt    time.Time `db:"created_at"`
}

// UserPermission represents the junction table for direct user-permission assignments
type UserPermission struct {
	ID           int64     `db:"id"`
	UserID       int64     `db:"user_id"`
	PermissionID int64     `db:"permission_id"`
	CreatedAt    time.Time `db:"created_at"`
}

// ProjectMemberDB represents the project members table record
type ProjectMemberDB struct {
	ID        int64     `db:"id"`
	ProjectID int64     `db:"project_id"`
	UserID    int64     `db:"user_id"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ProjectMemberInput contains data for creating a project member
type ProjectMemberInput struct {
	UserID int64
	Role   string
}

// ProjectMemberWithUser contains user details and their role in the project
type ProjectMemberWithUser struct {
	User      *user.User
	Role      string
	CreatedAt time.Time
}

// ProjectMemberQueryResult represents the query result from database
type ProjectMemberQueryResult struct {
	// User fields
	UserID        int64     `db:"user_id"`
	UserPublicID  string    `db:"user_public_id"`
	Email         string    `db:"email"`
	FirstName     string    `db:"first_name"`
	LastName      string    `db:"last_name"`
	IsActive      bool      `db:"is_active"`
	UserCreatedAt time.Time `db:"user_created_at"`
	UserUpdatedAt time.Time `db:"user_updated_at"`
	// Project member fields
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

// ToProjectMemberWithUser converts query result to domain model
func (r *ProjectMemberQueryResult) ToProjectMemberWithUser() *ProjectMemberWithUser {
	return &ProjectMemberWithUser{
		User: &user.User{
			ID:        r.UserPublicID,
			Email:     r.Email,
			FirstName: r.FirstName,
			LastName:  r.LastName,
			IsActive:  r.IsActive,
			CreatedAt: r.UserCreatedAt,
			UpdatedAt: r.UserUpdatedAt,
		},
		Role:      r.Role,
		CreatedAt: r.CreatedAt,
	}
}

// RoleQueryResult represents role query result from JOIN query
type RoleQueryResult struct {
	ID          int64     `db:"id"`
	PublicID    string    `db:"public_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ToRole converts query result to role domain model
func (r *RoleQueryResult) ToRole() *role.Role {
	return &role.Role{
		ID:          r.PublicID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// PermissionQueryResult represents permission query result from JOIN query
type PermissionQueryResult struct {
	ID          int64     `db:"id"`
	PublicID    string    `db:"public_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ToPermission converts query result to permission domain model
func (r *PermissionQueryResult) ToPermission() *permission.Permission {
	return &permission.Permission{
		ID:          r.PublicID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// UserProjectMembership contains project details and user's role in it
type UserProjectMembership struct {
	ProjectID   string
	ProjectName string
	Role        string
	JoinedAt    time.Time
}

// UserProjectQueryResult represents the query result for user projects from database
type UserProjectQueryResult struct {
	ProjectPublicID string    `db:"project_public_id"`
	ProjectName     string    `db:"project_name"`
	Role            string    `db:"role"`
	JoinedAt        time.Time `db:"joined_at"`
}

// ToUserProjectMembership converts query result to domain model
func (r *UserProjectQueryResult) ToUserProjectMembership() *UserProjectMembership {
	return &UserProjectMembership{
		ProjectID:   r.ProjectPublicID,
		ProjectName: r.ProjectName,
		Role:        r.Role,
		JoinedAt:    r.JoinedAt,
	}
}
