package oauth_provider

import (
	"context"

	"github.com/hrz8/altalune/internal/shared/query"
)

type Repository interface {
	// Query returns a paginated list of OAuth providers
	Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[OAuthProvider], error)

	// Create creates a new OAuth provider (encrypts client_secret)
	Create(ctx context.Context, input *CreateOAuthProviderInput) (*CreateOAuthProviderResult, error)

	// GetByID retrieves an OAuth provider by public ID
	GetByID(ctx context.Context, publicID string) (*OAuthProvider, error)

	// GetByProviderType retrieves an OAuth provider by provider type
	GetByProviderType(ctx context.Context, providerType ProviderType) (*OAuthProvider, error)

	// Update updates an OAuth provider (re-encrypts client_secret if provided)
	Update(ctx context.Context, input *UpdateOAuthProviderInput) (*UpdateOAuthProviderResult, error)

	// Delete deletes an OAuth provider by public ID
	Delete(ctx context.Context, input *DeleteOAuthProviderInput) error

	// RevealClientSecret decrypts and returns the plaintext client secret
	RevealClientSecret(ctx context.Context, publicID string) (string, error)
}
