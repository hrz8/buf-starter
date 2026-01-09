package oauth_seeder

import (
	"fmt"

	"github.com/hrz8/altalune/internal/shared/crypto"
	"github.com/hrz8/altalune/internal/shared/password"
)

const (
	// Production parameters for Argon2id hashing
	// These parameters provide strong security while maintaining acceptable performance (~50-100ms)
	argon2Iterations = 2         // Time cost (t)
	argon2Memory     = 64 * 1024 // Memory cost in KB (m) = 64MB
	argon2Threads    = 4         // Parallelism (p)
	argon2KeyLength  = 32        // Hash length in bytes
)

// HashClientSecret hashes an OAuth client secret using Argon2id
// This is used for dashboard client secrets that need to be verified during authentication
// Replaces bcrypt with modern Argon2id (PHC 2015 winner, OWASP recommended)
func HashClientSecret(secret string) (string, error) {
	if len(secret) < 32 {
		return "", fmt.Errorf("client secret must be at least 32 characters, got %d", len(secret))
	}

	hash, err := password.HashPassword(secret, password.HashOption{
		Iterations: argon2Iterations,
		Memory:     argon2Memory,
		Threads:    argon2Threads,
		Len:        argon2KeyLength,
	})
	if err != nil {
		return "", fmt.Errorf("hash client secret with argon2: %w", err)
	}

	return hash, nil
}

// VerifyClientSecret verifies an OAuth client secret against a stored Argon2id hash
// Uses constant-time comparison to prevent timing attacks
// Returns true if the secret matches the hash, false otherwise
func VerifyClientSecret(secret, hash string) (bool, error) {
	match, err := password.VerifyPassword(secret, hash)
	if err != nil {
		return false, fmt.Errorf("verify client secret: %w", err)
	}
	return match, nil
}

// EncryptProviderSecret encrypts an OAuth provider secret using AES-256-GCM
// This is used for Google/GitHub client secrets that need to be retrieved during OAuth flows
func EncryptProviderSecret(secret string, encryptionKey []byte) (string, error) {
	// Validate key size (must be 32 bytes for AES-256)
	if len(encryptionKey) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes, got %d", len(encryptionKey))
	}

	// Encrypt using the shared crypto package
	encrypted, err := crypto.Encrypt(secret, encryptionKey)
	if err != nil {
		return "", fmt.Errorf("encrypt secret: %w", err)
	}

	return encrypted, nil
}
