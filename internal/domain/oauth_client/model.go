package oauth_client

import (
	"time"

	"github.com/google/uuid"
)

// OAuthClientQueryResult represents query result with internal database ID
// OAuth clients are GLOBAL entities (infrastructure-level, like Auth0 Applications)
type OAuthClientQueryResult struct {
	ID           int64  // Internal database ID
	PublicID     string // Public nanoid
	Name         string
	ClientID     uuid.UUID // OAuth client_id (UUID)
	RedirectURIs []string
	PKCERequired bool
	IsDefault    bool
	Confidential bool // true = requires secret (confidential), false = public/SPA
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// OAuthClient represents the domain model with public IDs only
// OAuth clients are GLOBAL entities (infrastructure-level, like Auth0 Applications)
type OAuthClient struct {
	ID           string // Public nanoid
	Name         string
	ClientID     uuid.UUID // OAuth client_id (UUID)
	RedirectURIs []string
	PKCERequired bool
	IsDefault    bool
	Confidential bool // true = requires secret (confidential), false = public/SPA
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CreateOAuthClientInput represents input for creating an OAuth client
type CreateOAuthClientInput struct {
	Name          string
	RedirectURIs  []string
	PKCERequired  bool
	AllowedScopes []string
	Confidential  bool // true = requires secret (confidential), false = public/SPA
}

// CreateOAuthClientResult represents the result of creating an OAuth client
type CreateOAuthClientResult struct {
	Client       *OAuthClient
	ClientSecret string // Plaintext secret (ONLY for creation response)
}

// UpdateOAuthClientInput represents input for updating an OAuth client
type UpdateOAuthClientInput struct {
	PublicID      string
	Name          *string
	RedirectURIs  []string
	PKCERequired  *bool
	AllowedScopes []string
}

// ToOAuthClient converts query result to domain model (hides internal IDs)
func (r *OAuthClientQueryResult) ToOAuthClient() *OAuthClient {
	return &OAuthClient{
		ID:           r.PublicID,
		Name:         r.Name,
		ClientID:     r.ClientID,
		RedirectURIs: r.RedirectURIs,
		PKCERequired: r.PKCERequired,
		IsDefault:    r.IsDefault,
		Confidential: r.Confidential,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
