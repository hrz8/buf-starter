package altalune

import "time"

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
}
