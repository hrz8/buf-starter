package permission

import (
	"context"

	"github.com/hrz8/altalune/internal/shared/query"
)

type Repository interface {
	GetIDByPublicID(ctx context.Context, publicID string) (int64, error)
	Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[Permission], error)
	Create(ctx context.Context, input *CreatePermissionInput) (*CreatePermissionResult, error)
	GetByName(ctx context.Context, name string) (*Permission, error)
	GetByID(ctx context.Context, publicID string) (*Permission, error)
	Update(ctx context.Context, input *UpdatePermissionInput) (*UpdatePermissionResult, error)
	Delete(ctx context.Context, publicID string) error
}
