package oauth_seeder

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/hrz8/altalune/internal/shared/crypto"
)

const (
	// BcryptCost is the cost factor for bcrypt hashing (higher = more secure but slower)
	// Cost 12 takes ~250ms on modern hardware, good balance between security and performance
	BcryptCost = 12
)

// HashClientSecret hashes an OAuth client secret using bcrypt
// This is used for dashboard client secrets that need to be verified during authentication
func HashClientSecret(secret string) (string, error) {
	if len(secret) < 32 {
		return "", fmt.Errorf("client secret must be at least 32 characters, got %d", len(secret))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(secret), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash client secret: %w", err)
	}

	return string(hash), nil
}

// EncryptProviderSecret encrypts an OAuth provider secret using AES-256-GCM
// This is used for Google/GitHub client secrets that need to be retrieved during OAuth flows
func EncryptProviderSecret(secret string, encryptionKey string) (string, error) {
	// Decode the base64-encoded encryption key
	keyBytes, err := base64.StdEncoding.DecodeString(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("decode encryption key: %w", err)
	}

	// Validate key size (must be 32 bytes for AES-256)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes, got %d", len(keyBytes))
	}

	// Encrypt using the shared crypto package
	encrypted, err := crypto.Encrypt(secret, keyBytes)
	if err != nil {
		return "", fmt.Errorf("encrypt secret: %w", err)
	}

	return encrypted, nil
}
