package user

import (
	"context"

	"github.com/hrz8/altalune/internal/shared/query"
)

type Repository interface {
	GetIDByPublicID(ctx context.Context, publicID string) (int64, error)
	Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[User], error)
	Create(ctx context.Context, input *CreateUserInput) (*CreateUserResult, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, publicID string) (*User, error)
	Update(ctx context.Context, input *UpdateUserInput) (*UpdateUserResult, error)
	Delete(ctx context.Context, publicID string) error
	Activate(ctx context.Context, publicID string) (*User, error)
	Deactivate(ctx context.Context, publicID string) (*User, error)
}
