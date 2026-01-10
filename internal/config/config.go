package config

import (
	"encoding/base64"
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

// Seeder configuration
func (c *AppConfig) GetSuperadminEmail() string {
	return c.Seeder.Superadmin.Email
}

func (c *AppConfig) GetDefaultOAuthClientName() string {
	return c.Seeder.DefaultOAuthClient.Name
}

func (c *AppConfig) GetDefaultOAuthClientID() string {
	return c.Seeder.DefaultOAuthClient.ClientID
}

func (c *AppConfig) GetDefaultOAuthClientSecret() string {
	return c.Seeder.DefaultOAuthClient.ClientSecret
}

func (c *AppConfig) GetDefaultOAuthClientRedirectURIs() []string {
	uris := make([]string, len(c.Seeder.DefaultOAuthClient.RedirectURIs))
	copy(uris, c.Seeder.DefaultOAuthClient.RedirectURIs)
	return uris
}

func (c *AppConfig) GetDefaultOAuthClientPKCERequired() bool {
	return c.Seeder.DefaultOAuthClient.PKCERequired
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
