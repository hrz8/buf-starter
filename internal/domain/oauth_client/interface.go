package oauth_client

import (
	"context"

	"github.com/hrz8/altalune/internal/shared/query"
)

// Repositor defines the interface for OAuth client repository operations
type Repositor interface {
	// Create creates a new OAuth client with generated client_id and hashed secret
	Create(ctx context.Context, input *CreateOAuthClientInput) (*CreateOAuthClientResult, error)

	// Query returns a paginated list of OAuth clients for a project
	Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[OAuthClient], error)

	// GetByPublicID retrieves an OAuth client by its public nanoid
	GetByPublicID(ctx context.Context, projectID int64, publicID string) (*OAuthClient, error)

	// GetByClientID retrieves an OAuth client by its UUID client_id (for OAuth flows)
	GetByClientID(ctx context.Context, clientID string) (*OAuthClient, error)

	// Update updates an existing OAuth client
	Update(ctx context.Context, input *UpdateOAuthClientInput) (*OAuthClient, error)

	// Delete deletes an OAuth client (with default client protection)
	Delete(ctx context.Context, projectID int64, publicID string) error

	// RevealClientSecret retrieves the hashed client secret (with audit logging)
	RevealClientSecret(ctx context.Context, projectID int64, publicID string) (string, error)
}
