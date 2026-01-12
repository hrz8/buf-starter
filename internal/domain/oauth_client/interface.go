package oauth_client

import (
	"context"

	"github.com/hrz8/altalune/internal/shared/query"
)

// Repositor defines the interface for OAuth client repository operations
// OAuth clients are GLOBAL entities (not project-scoped)
type Repositor interface {
	// Create creates a new OAuth client with generated client_id and hashed secret
	Create(ctx context.Context, input *CreateOAuthClientInput) (*CreateOAuthClientResult, error)

	// Query returns a paginated list of all OAuth clients
	Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[OAuthClient], error)

	// GetByPublicID retrieves an OAuth client by its public nanoid
	GetByPublicID(ctx context.Context, publicID string) (*OAuthClient, error)

	// GetByClientID retrieves an OAuth client by its UUID client_id (for OAuth flows)
	GetByClientID(ctx context.Context, clientID string) (*OAuthClient, error)

	// Update updates an existing OAuth client
	Update(ctx context.Context, input *UpdateOAuthClientInput) (*OAuthClient, error)

	// Delete deletes an OAuth client (with default client protection)
	Delete(ctx context.Context, publicID string) error

	// RevealClientSecret retrieves the hashed client secret (with audit logging)
	RevealClientSecret(ctx context.Context, publicID string) (string, error)
}
