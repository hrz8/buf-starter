package altalune

import "time"

// OAuthProviderConfig represents OAuth provider configuration
// This is a data transfer object used by the seeder
type OAuthProviderConfig struct {
	Provider     string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       string
	Enabled      bool
}

type Config interface {
	// Server configuration
	GetServerHost() string
	GetServerPort() int
	GetServerLogLevel() string
	IsHTTPLoggingEnabled() bool
	IsCORSEnabled() bool
	GetServerReadTimeout() time.Duration
	GetServerWriteTimeout() time.Duration
	GetServerIdleTimeout() time.Duration
	GetServerCleanupTimeout() time.Duration

	// Database configuration
	GetDatabaseURL() string
	GetDatabaseMaxConnections() int
	GetDatabaseMaxIdleTime() time.Duration
	GetDatabaseConnectTimeout() time.Duration

	// Security configuration
	GetAllowedOrigins() []string

	// IAM encryption configuration
	// GetIAMEncryptionKey returns the 32-byte encryption key for IAM secrets
	// This key is used to encrypt/decrypt OAuth client secrets
	GetIAMEncryptionKey() []byte

	// JWT configuration
	GetJWTPrivateKeyPath() string
	GetJWTPublicKeyPath() string
	GetJWKSKid() string

	// Auth configuration
	GetAuthHost() string
	GetAuthPort() int
	GetSessionSecret() string
	GetCodeExpiry() int
	GetAccessTokenExpiry() int
	GetRefreshTokenExpiry() int

	// Seeder configuration
	GetSuperadminEmail() string
	GetDefaultOAuthClientName() string
	GetDefaultOAuthClientID() string
	GetDefaultOAuthClientSecret() string
	GetDefaultOAuthClientRedirectURIs() []string
	GetDefaultOAuthClientPKCERequired() bool
	GetOAuthProviders() []OAuthProviderConfig
}
