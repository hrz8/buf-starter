package user

import (
	"context"

	"github.com/hrz8/altalune/internal/shared/query"
)

type Repository interface {
	GetIDByPublicID(ctx context.Context, publicID string) (int64, error)
	GetInternalIDByEmail(ctx context.Context, email string) (int64, error)
	Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[User], error)
	Create(ctx context.Context, input *CreateUserInput) (*CreateUserResult, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, publicID string) (*User, error)
	GetByInternalID(ctx context.Context, internalID int64) (*User, error)
	Update(ctx context.Context, input *UpdateUserInput) (*UpdateUserResult, error)
	UpdateProfileByInternalID(ctx context.Context, internalID int64, firstName, lastName string) (*User, error)
	Delete(ctx context.Context, publicID string) error
	Activate(ctx context.Context, publicID string) (*User, error)
	Deactivate(ctx context.Context, publicID string) (*User, error)

	// User Identity operations for OAuth authentication
	GetUserIdentityByProvider(ctx context.Context, provider, providerUserID string) (*UserIdentity, error)
	GetUserIdentities(ctx context.Context, userID int64) ([]*UserIdentity, error)
	CreateUserIdentity(ctx context.Context, input *CreateUserIdentityInput) error
	UpdateUserIdentityLastLogin(ctx context.Context, userID int64, provider string) error

	// Project membership for OAuth user onboarding
	AddProjectMember(ctx context.Context, projectID, userID int64, role string) error
}
