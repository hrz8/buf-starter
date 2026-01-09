package role

import (
	"time"

	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Role represents a role in the system
type Role struct {
	ID          string // Public nanoid
	Name        string // Unique
	Description string // Optional
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *Role) ToRoleProto() *altalunev1.Role {
	return &altalunev1.Role{
		Id:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
}

// RoleQueryResult represents a single role query result
type RoleQueryResult struct {
	ID          int64  // Internal ID
	PublicID    string // Public nanoid
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *RoleQueryResult) ToRole() *Role {
	return &Role{
		ID:          r.PublicID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// CreateRoleInput contains data for creating a new role
type CreateRoleInput struct {
	Name        string
	Description string
}

// CreateRoleResult represents the result of creating a role
type CreateRoleResult struct {
	ID          int64
	PublicID    string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *CreateRoleResult) ToRole() *Role {
	return &Role{
		ID:          r.PublicID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// UpdateRoleInput contains data for updating a role
type UpdateRoleInput struct {
	ID          int64  // Internal ID
	PublicID    string // Public ID
	Name        string
	Description string
}

// UpdateRoleResult represents the result of updating a role
type UpdateRoleResult struct {
	ID          int64
	PublicID    string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *UpdateRoleResult) ToRole() *Role {
	return &Role{
		ID:          r.PublicID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
