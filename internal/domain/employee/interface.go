package employee

import (
	"context"

	"github.com/hrz8/altalune/internal/shared/query"
)

type Repositor interface {
	Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[Employee], error)
	Create(ctx context.Context, input *CreateEmployeeInput) (*CreateEmployeeResult, error)
	GetByEmail(ctx context.Context, projectID int64, email string) (*Employee, error)
}
