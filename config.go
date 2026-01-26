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
	GetJWTIssuer() string // JWT issuer URL (derived from auth server URL)

	// Auth configuration
	GetAuthHost() string
	GetAuthPort() int
	GetSessionSecret() string
	GetCodeExpiry() int
	GetAccessTokenExpiry() int
	GetRefreshTokenExpiry() int
	IsAutoActivate() bool // Whether new users are automatically activated (default: true)

	// Seeder configuration
	GetSuperadminEmail() string
	GetOAuthProviders() []OAuthProviderConfig

	// Dashboard OAuth configuration (from dashboardOauth config section)
	IsDashboardOAuthExternalServer() bool
	GetDashboardOAuthServerURL() string
	GetDefaultOAuthClientName() string
	GetDefaultOAuthClientID() string
	GetDefaultOAuthClientSecret() string
	GetDefaultOAuthClientRedirectURIs() []string
	GetDefaultOAuthClientPKCERequired() bool

	// Notification configuration
	GetNotificationAuthBaseURL() string   // Base URL for verification links in emails
	GetNotificationEmailProvider() string // "resend" or "ses"
	GetNotificationEmailFromEmail() string
	GetNotificationEmailFromName() string
	GetNotificationResendAPIKey() string
	GetNotificationSESRegion() string

	// OTP configuration
	GetOTPExpirySeconds() int
	GetOTPRateLimit() int
	GetOTPRateLimitWindowMins() int

	// Email verification configuration
	GetVerificationTokenExpiryHours() int

	// Branding configuration
	GetDashboardBrandingName() string
	GetAuthServerBrandingName() string
}
