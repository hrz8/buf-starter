package config

import (
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
