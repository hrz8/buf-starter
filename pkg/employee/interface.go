package employee

import (
	"context"

	"github.com/hrz8/altalune/internal/query"
)

type Repositor interface {
	Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[Employee], error)
}
