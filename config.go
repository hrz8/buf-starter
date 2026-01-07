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

	// IAM encryption configuration
	// GetIAMEncryptionKey returns the 32-byte encryption key for IAM secrets
	// This key is used to encrypt/decrypt OAuth client secrets
	GetIAMEncryptionKey() []byte
}
