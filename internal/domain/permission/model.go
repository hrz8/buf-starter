package permission

import (
	"time"

	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Permission represents a permission in the system
type Permission struct {
	ID          string    // Public nanoid
	Name        string    // Machine-readable: "project:read"
	Effect      string    // "allow" or "deny"
	Description string    // Human-readable (optional)
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *Permission) ToPermissionProto() *altalunev1.Permission {
	return &altalunev1.Permission{
		Id:          m.ID,
		Name:        m.Name,
		Effect:      m.Effect,
		Description: m.Description,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
}

// PermissionQueryResult represents a single permission query result
type PermissionQueryResult struct {
	ID          int64  // Internal ID
	PublicID    string // Public nanoid
	Name        string
	Effect      string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *PermissionQueryResult) ToPermission() *Permission {
	return &Permission{
		ID:          r.PublicID,
		Name:        r.Name,
		Effect:      r.Effect,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// CreatePermissionInput contains data for creating a new permission
type CreatePermissionInput struct {
	Name        string
	Effect      string
	Description string
}

// CreatePermissionResult represents the result of creating a permission
type CreatePermissionResult struct {
	ID          int64
	PublicID    string
	Name        string
	Effect      string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *CreatePermissionResult) ToPermission() *Permission {
	return &Permission{
		ID:          r.PublicID,
		Name:        r.Name,
		Effect:      r.Effect,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// UpdatePermissionInput contains data for updating a permission
type UpdatePermissionInput struct {
	ID          int64  // Internal ID
	PublicID    string // Public ID
	Name        string
	Effect      string
	Description string
}

// UpdatePermissionResult represents the result of updating a permission
type UpdatePermissionResult struct {
	ID          int64
	PublicID    string
	Name        string
	Effect      string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *UpdatePermissionResult) ToPermission() *Permission {
	return &Permission{
		ID:          r.PublicID,
		Name:        r.Name,
		Effect:      r.Effect,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
