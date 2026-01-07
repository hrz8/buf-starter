package oauth_provider

import (
	"time"

	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProviderType represents the OAuth provider type
type ProviderType string

const (
	ProviderTypeGoogle    ProviderType = "google"
	ProviderTypeGithub    ProviderType = "github"
	ProviderTypeMicrosoft ProviderType = "microsoft"
	ProviderTypeApple     ProviderType = "apple"
)

// OAuthProvider represents an OAuth provider configuration
// CRITICAL: Never exposes actual client_secret (use ClientSecretSet instead)
type OAuthProvider struct {
	ID              string       // Public nanoid
	ProviderType    ProviderType // OAuth provider type
	ClientID        string       // OAuth client ID (public)
	ClientSecretSet bool         // True if secret exists (NEVER actual secret)
	RedirectURL     string       // OAuth redirect/callback URL
	Scopes          string       // Comma-separated OAuth scopes
	Enabled         bool         // Whether provider is enabled
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (m *OAuthProvider) ToOAuthProviderProto() *altalunev1.OAuthProvider {
	return &altalunev1.OAuthProvider{
		Id:              m.ID,
		ProviderType:    ProviderTypeToProto(m.ProviderType),
		ClientId:        m.ClientID,
		ClientSecretSet: m.ClientSecretSet,
		RedirectUrl:     m.RedirectURL,
		Scopes:          m.Scopes,
		Enabled:         m.Enabled,
		CreatedAt:       timestamppb.New(m.CreatedAt),
		UpdatedAt:       timestamppb.New(m.UpdatedAt),
	}
}

// OAuthProviderQueryResult represents a single OAuth provider query result
type OAuthProviderQueryResult struct {
	ID           int64     // Internal ID
	PublicID     string    // Public nanoid
	ProviderType string    // OAuth provider type (stored as string in DB)
	ClientID     string
	RedirectURL  string
	Scopes       string
	Enabled      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r *OAuthProviderQueryResult) ToOAuthProvider() *OAuthProvider {
	return &OAuthProvider{
		ID:              r.PublicID,
		ProviderType:    ProviderType(r.ProviderType),
		ClientID:        r.ClientID,
		ClientSecretSet: true, // If record exists, secret is set
		RedirectURL:     r.RedirectURL,
		Scopes:          r.Scopes,
		Enabled:         r.Enabled,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

// CreateOAuthProviderInput contains data for creating a new OAuth provider
type CreateOAuthProviderInput struct {
	ProviderType ProviderType
	ClientID     string
	ClientSecret string // Plaintext (encrypted in repo)
	RedirectURL  string
	Scopes       string
	Enabled      bool
}

// CreateOAuthProviderResult represents the result of creating an OAuth provider
type CreateOAuthProviderResult struct {
	ID           int64
	PublicID     string
	ProviderType ProviderType
	ClientID     string
	RedirectURL  string
	Scopes       string
	Enabled      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r *CreateOAuthProviderResult) ToOAuthProvider() *OAuthProvider {
	return &OAuthProvider{
		ID:              r.PublicID,
		ProviderType:    r.ProviderType,
		ClientID:        r.ClientID,
		ClientSecretSet: true, // Secret was just set during creation
		RedirectURL:     r.RedirectURL,
		Scopes:          r.Scopes,
		Enabled:         r.Enabled,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

// UpdateOAuthProviderInput contains data for updating an OAuth provider
type UpdateOAuthProviderInput struct {
	PublicID     string
	ClientID     string
	ClientSecret string // Optional - if empty, retain existing secret
	RedirectURL  string
	Scopes       string
	Enabled      bool
}

// UpdateOAuthProviderResult represents the result of updating an OAuth provider
type UpdateOAuthProviderResult struct {
	ID          int64
	PublicID    string
	ClientID    string
	RedirectURL string
	Scopes      string
	Enabled     bool
	UpdatedAt   time.Time
}

// ToOAuthProvider converts result to OAuthProvider with preserved provider_type
func (r *UpdateOAuthProviderResult) ToOAuthProvider(providerType ProviderType, createdAt time.Time) *OAuthProvider {
	return &OAuthProvider{
		ID:              r.PublicID,
		ProviderType:    providerType, // Preserved from existing record
		ClientID:        r.ClientID,
		ClientSecretSet: true,
		RedirectURL:     r.RedirectURL,
		Scopes:          r.Scopes,
		Enabled:         r.Enabled,
		CreatedAt:       createdAt, // Preserved from existing record
		UpdatedAt:       r.UpdatedAt,
	}
}

// DeleteOAuthProviderInput contains data for deleting an OAuth provider
type DeleteOAuthProviderInput struct {
	PublicID string
}
