package api_key

import (
	"context"

	"github.com/hrz8/altalune/internal/shared/query"
)

type Repositor interface {
	Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[ApiKey], error)
	Create(ctx context.Context, input *CreateApiKeyInput) (*CreateApiKeyResult, error)
	GetByID(ctx context.Context, projectID int64, publicID string) (*ApiKey, error)
	GetByKey(ctx context.Context, key string) (*ApiKey, error) // For authentication
	Update(ctx context.Context, input *UpdateApiKeyInput) (*UpdateApiKeyResult, error)
	Delete(ctx context.Context, input *DeleteApiKeyInput) error
	Activate(ctx context.Context, input *ActivateApiKeyInput) (*ActivateApiKeyResult, error)
	Deactivate(ctx context.Context, input *DeactivateApiKeyInput) (*DeactivateApiKeyResult, error)
}