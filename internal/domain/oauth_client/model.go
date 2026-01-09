package oauth_client

import (
	"time"

	"github.com/google/uuid"
)

// OAuthClientQueryResult represents query result with internal database ID
type OAuthClientQueryResult struct {
	ID              int64  // Internal database ID
	PublicID        string // Public nanoid
	ProjectID       int64  // Internal project ID
	ProjectPublicID string // Project public nanoid
	Name            string
	ClientID        uuid.UUID // OAuth client_id (UUID)
	RedirectURIs    []string
	PKCERequired    bool
	IsDefault       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// OAuthClient represents the domain model with public IDs only
type OAuthClient struct {
	ID           string // Public nanoid
	ProjectID    string // Project public nanoid
	Name         string
	ClientID     uuid.UUID // OAuth client_id (UUID)
	RedirectURIs []string
	PKCERequired bool
	IsDefault    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CreateOAuthClientInput represents input for creating an OAuth client
type CreateOAuthClientInput struct {
	ProjectID       int64  // Internal project ID
	ProjectPublicID string // Project public nanoid (for response)
	Name            string
	RedirectURIs    []string
	PKCERequired    bool
	AllowedScopes   []string
}

// CreateOAuthClientResult represents the result of creating an OAuth client
type CreateOAuthClientResult struct {
	Client       *OAuthClient
	ClientSecret string // Plaintext secret (ONLY for creation response)
}

// UpdateOAuthClientInput represents input for updating an OAuth client
type UpdateOAuthClientInput struct {
	PublicID      string
	ProjectID     int64
	Name          *string
	RedirectURIs  []string
	PKCERequired  *bool
	AllowedScopes []string
}

// ToOAuthClient converts query result to domain model (hides internal IDs)
func (r *OAuthClientQueryResult) ToOAuthClient() *OAuthClient {
	return &OAuthClient{
		ID:           r.PublicID,
		ProjectID:    r.ProjectPublicID,
		Name:         r.Name,
		ClientID:     r.ClientID,
		RedirectURIs: r.RedirectURIs,
		PKCERequired: r.PKCERequired,
		IsDefault:    r.IsDefault,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
