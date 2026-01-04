package project

import (
	"context"

	"github.com/hrz8/altalune/internal/shared/query"
)

type Repositor interface {
	GetIDByPublicID(ctx context.Context, publicID string) (int64, error)
	Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[Project], error)
	Create(ctx context.Context, input *CreateProjectInput) (*CreateProjectResult, error)
	GetByName(ctx context.Context, name string) (*Project, error)
	GetByID(ctx context.Context, publicID string) (*Project, error)
	Update(ctx context.Context, input *UpdateProjectInput) (*UpdateProjectResult, error)
	Delete(ctx context.Context, publicID string) error
}
