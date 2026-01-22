package role

import (
	"context"

	"github.com/hrz8/altalune/internal/shared/query"
)

type Repository interface {
	GetIDByPublicID(ctx context.Context, publicID string) (int64, error)
	GetInternalIDByName(ctx context.Context, name string) (int64, error)
	Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[Role], error)
	Create(ctx context.Context, input *CreateRoleInput) (*CreateRoleResult, error)
	GetByName(ctx context.Context, name string) (*Role, error)
	GetByID(ctx context.Context, publicID string) (*Role, error)
	Update(ctx context.Context, input *UpdateRoleInput) (*UpdateRoleResult, error)
	Delete(ctx context.Context, publicID string) error
}
