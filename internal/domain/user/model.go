package user

import (
	"time"

	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// User represents a system user with OAuth-only authentication
type User struct {
	ID        string // Public nanoid
	Email     string // Unique, lowercase
	FirstName string // Optional
	LastName  string // Optional
	IsActive  bool   // User activation status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *User) ToUserProto() *altalunev1.User {
	return &altalunev1.User{
		Id:        m.ID,
		Email:     m.Email,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		IsActive:  m.IsActive,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}

// UserQueryResult represents a single user query result
type UserQueryResult struct {
	ID        int64  // Internal ID
	PublicID  string // Public nanoid
	Email     string
	FirstName string
	LastName  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *UserQueryResult) ToUser() *User {
	return &User{
		ID:        r.PublicID,
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		IsActive:  r.IsActive,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// CreateUserInput contains data for creating a new user
type CreateUserInput struct {
	Email     string
	FirstName string
	LastName  string
}

// CreateUserResult represents the result of creating a user
type CreateUserResult struct {
	ID        int64
	PublicID  string
	Email     string
	FirstName string
	LastName  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *CreateUserResult) ToUser() *User {
	return &User{
		ID:        r.PublicID,
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		IsActive:  r.IsActive,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// UpdateUserInput contains data for updating a user
type UpdateUserInput struct {
	ID        int64  // Internal ID
	PublicID  string // Public ID
	Email     string
	FirstName string
	LastName  string
}

// UpdateUserResult represents the result of updating a user
type UpdateUserResult struct {
	ID        int64
	PublicID  string
	Email     string
	FirstName string
	LastName  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *UpdateUserResult) ToUser() *User {
	return &User{
		ID:        r.PublicID,
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		IsActive:  r.IsActive,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// ActivateUserInput contains data for activating a user
type ActivateUserInput struct {
	PublicID string
}

// DeactivateUserInput contains data for deactivating a user
type DeactivateUserInput struct {
	PublicID string
}

// UserIdentity represents an OAuth provider identity linked to a user
type UserIdentity struct {
	ID                    int64
	PublicID              string
	UserID                int64
	Provider              string
	ProviderUserID        string
	Email                 string
	FirstName             string
	LastName              string
	OAuthClientID         *string // UUID of OAuth client (nullable)
	OriginOAuthClientName *string // Client name at signup time (historical snapshot, nullable)
	LastLoginAt           *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// CreateUserIdentityInput contains data for creating a user identity
type CreateUserIdentityInput struct {
	UserID                int64
	Provider              string
	ProviderUserID        string
	Email                 string
	FirstName             string
	LastName              string
	OAuthClientID         *string // UUID of OAuth client (nullable)
	OriginOAuthClientName *string // Client name at signup time (historical snapshot, nullable)
}
