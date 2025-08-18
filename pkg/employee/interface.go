package employee

import (
	"context"

	"github.com/hrz8/altalune"
)

type Repositor interface {
	Query(ctx context.Context, projectID int64, params *altalune.QueryParams) (*altalune.QueryResult[Employee], error)
}
