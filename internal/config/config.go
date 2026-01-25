package config

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hrz8/altalune"
)

var _ altalune.Config = (*AppConfig)(nil)

func (c *AppConfig) GetServerHost() string {
	return c.Server.Host
}

func (c *AppConfig) GetServerPort() int {
	return c.Server.Port
}

func (c *AppConfig) GetServerLogLevel() string {
	return c.Server.LogLevel
}

func (c *AppConfig) IsHTTPLoggingEnabled() bool {
	return c.Server.HTTPLogging
}

func (c *AppConfig) IsCORSEnabled() bool {
	return c.Server.EnableCORS
}

func (c *AppConfig) GetServerReadTimeout() time.Duration {
	return time.Duration(c.Server.ReadTimeout) * time.Second
}

func (c *AppConfig) GetServerWriteTimeout() time.Duration {
	return time.Duration(c.Server.WriteTimeout) * time.Second
}

func (c *AppConfig) GetServerIdleTimeout() time.Duration {
	return time.Duration(c.Server.IdleTimeout) * time.Second
}

func (c *AppConfig) GetServerCleanupTimeout() time.Duration {
	return time.Duration(c.Server.CleanupTimeout) * time.Second
}

func (c *AppConfig) GetDatabaseURL() string {
	return c.Database.URL
}

func (c *AppConfig) GetDatabaseMaxConnections() int {
	return c.Database.MaxConnections
}

func (c *AppConfig) GetDatabaseMaxIdleTime() time.Duration {
	return time.Duration(c.Database.MaxIdleTime) * time.Second
}

func (c *AppConfig) GetDatabaseConnectTimeout() time.Duration {
	return time.Duration(c.Database.ConnectTimeout) * time.Second
}

func (c *AppConfig) GetAllowedOrigins() []string {
	origins := make([]string, len(c.Security.AllowedOrigins))
	copy(origins, c.Security.AllowedOrigins)
	return origins
}

func (c *AppConfig) GetIAMEncryptionKey() []byte {
	// Return the encryption key from YAML config (no environment variable fallback)
	// The key is stored as base64-encoded string in config.yaml (44 chars)
	// and must be decoded to get the 32-byte key
	key, err := base64.StdEncoding.DecodeString(c.Security.IAMEncryptionKey)
	if err != nil {
		// Return empty slice if decode fails (will fail validation)
		return []byte{}
	}
	return key
}

// JWT configuration
func (c *AppConfig) GetJWTPrivateKeyPath() string {
	return c.Security.JWTPrivateKeyPath
}

func (c *AppConfig) GetJWTPublicKeyPath() string {
	return c.Security.JWTPublicKeyPath
}

func (c *AppConfig) GetJWKSKid() string {
	return c.Security.JWKSKid
}

// GetJWTIssuer returns the JWT issuer URL derived from auth server configuration.
// This URL should match the auth server's base URL for JWKS discovery.
func (c *AppConfig) GetJWTIssuer() string {
	// Use dashboardOauth.server as the issuer since it's the canonical auth server URL
	if c.DashboardOAuth != nil && c.DashboardOAuth.Server != "" {
		return c.DashboardOAuth.Server
	}
	// Fallback: construct from auth host/port
	return fmt.Sprintf("http://%s:%d", c.Auth.Host, c.Auth.Port)
}

// Auth configuration
func (c *AppConfig) GetAuthHost() string {
	return c.Auth.Host
}

func (c *AppConfig) GetAuthPort() int {
	return c.Auth.Port
}

func (c *AppConfig) GetSessionSecret() string {
	return c.Auth.SessionSecret
}

func (c *AppConfig) GetCodeExpiry() int {
	return c.Auth.CodeExpiry
}

func (c *AppConfig) GetAccessTokenExpiry() int {
	return c.Auth.AccessTokenExpiry
}

func (c *AppConfig) GetRefreshTokenExpiry() int {
	return c.Auth.RefreshTokenExpiry
}

func (c *AppConfig) IsAutoActivate() bool {
	return c.Auth.IsAutoActivate()
}

// Seeder configuration
func (c *AppConfig) GetSuperadminEmail() string {
	return c.Seeder.Superadmin.Email
}

func (c *AppConfig) GetDefaultOAuthClientName() string {
	return c.DashboardOAuth.Name
}

func (c *AppConfig) GetDefaultOAuthClientID() string {
	return c.DashboardOAuth.ClientID
}

func (c *AppConfig) GetDefaultOAuthClientSecret() string {
	return c.DashboardOAuth.ClientSecret
}

func (c *AppConfig) GetDefaultOAuthClientRedirectURIs() []string {
	uris := make([]string, len(c.DashboardOAuth.RedirectURIs))
	copy(uris, c.DashboardOAuth.RedirectURIs)
	return uris
}

func (c *AppConfig) GetDefaultOAuthClientPKCERequired() bool {
	return c.DashboardOAuth.PKCERequired
}

// Dashboard OAuth configuration (for token exchange proxy)
func (c *AppConfig) IsDashboardOAuthExternalServer() bool {
	return c.DashboardOAuth.ExternalServer
}

func (c *AppConfig) GetDashboardOAuthServerURL() string {
	return c.DashboardOAuth.Server
}

func (c *AppConfig) GetOAuthProviders() []altalune.OAuthProviderConfig {
	providers := make([]altalune.OAuthProviderConfig, len(c.Seeder.OAuthProviders))
	for i, p := range c.Seeder.OAuthProviders {
		providers[i] = altalune.OAuthProviderConfig{
			Provider:     p.Provider,
			ClientID:     p.ClientID,
			ClientSecret: p.ClientSecret,
			RedirectURL:  p.RedirectURL,
			Scopes:       p.Scopes,
			Enabled:      p.Enabled,
		}
	}
	return providers
}

// Notification configuration
func (c *AppConfig) GetNotificationEmailProvider() string {
	if c.Notification == nil || c.Notification.Email == nil {
		return ""
	}
	return c.Notification.Email.Provider
}

func (c *AppConfig) GetNotificationEmailFromEmail() string {
	if c.Notification == nil || c.Notification.Email == nil {
		return ""
	}
	if c.Notification.Email.Resend != nil && c.Notification.Email.Resend.FromEmail != "" {
		return c.Notification.Email.Resend.FromEmail
	}
	if c.Notification.Email.SES != nil && c.Notification.Email.SES.FromEmail != "" {
		return c.Notification.Email.SES.FromEmail
	}
	return ""
}

func (c *AppConfig) GetNotificationEmailFromName() string {
	if c.Notification == nil || c.Notification.Email == nil || c.Notification.Email.Resend == nil {
		return ""
	}
	return c.Notification.Email.Resend.FromName
}

func (c *AppConfig) GetNotificationResendAPIKey() string {
	if c.Notification == nil || c.Notification.Email == nil || c.Notification.Email.Resend == nil {
		return ""
	}
	return c.Notification.Email.Resend.APIKey
}

func (c *AppConfig) GetNotificationSESRegion() string {
	if c.Notification == nil || c.Notification.Email == nil || c.Notification.Email.SES == nil {
		return ""
	}
	return c.Notification.Email.SES.Region
}

// GetNotificationAuthBaseURL returns the base URL for verification links in emails.
func (c *AppConfig) GetNotificationAuthBaseURL() string {
	if c.Notification == nil {
		return ""
	}
	return c.Notification.AuthBaseURL
}

// OTP configuration
func (c *AppConfig) GetOTPExpirySeconds() int {
	if c.Notification == nil || c.Notification.OTP == nil {
		return 300 // 5 minutes default
	}
	return c.Notification.OTP.ExpirySeconds
}

func (c *AppConfig) GetOTPRateLimit() int {
	if c.Notification == nil || c.Notification.OTP == nil {
		return 3 // default
	}
	return c.Notification.OTP.RateLimit
}

func (c *AppConfig) GetOTPRateLimitWindowMins() int {
	if c.Notification == nil || c.Notification.OTP == nil {
		return 15 // 15 minutes default
	}
	return c.Notification.OTP.RateLimitWindowMins
}

// Email verification configuration
func (c *AppConfig) GetVerificationTokenExpiryHours() int {
	if c.Notification == nil || c.Notification.Verification == nil {
		return 24 // 24 hours default
	}
	return c.Notification.Verification.TokenExpiryHours
}
