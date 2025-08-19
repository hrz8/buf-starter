package project

import (
	"context"

	"github.com/hrz8/altalune/internal/query"
)

type Repositor interface {
	GetIDByPublicID(ctx context.Context, publicID string) (int64, error)
	Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[Project], error)
}
