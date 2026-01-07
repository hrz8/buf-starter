package crypto_test

import (
	"testing"

	"github.com/hrz8/altalune/internal/config"
	"github.com/hrz8/altalune/internal/shared/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegrationWithConfig tests that the crypto package works correctly
// with the encryption key loaded from config.yaml
func TestIntegrationWithConfig(t *testing.T) {
	// Load config
	cfg, err := config.Load("../../../config.yaml")
	require.NoError(t, err, "Should load config successfully")

	// Get encryption key from config
	key := cfg.GetIAMEncryptionKey()
	require.NotNil(t, key, "Encryption key should not be nil")

	// Validate the key
	err = crypto.ValidateKey(key)
	require.NoError(t, err, "Encryption key from config should be valid")

	// Test encryption/decryption round-trip
	plaintext := "my-oauth-client-secret-12345"

	// Encrypt
	ciphertext, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err, "Should encrypt successfully")
	assert.NotEmpty(t, ciphertext, "Ciphertext should not be empty")
	assert.NotEqual(t, plaintext, ciphertext, "Ciphertext should differ from plaintext")

	// Decrypt
	decrypted, err := crypto.Decrypt(ciphertext, key)
	require.NoError(t, err, "Should decrypt successfully")
	assert.Equal(t, plaintext, decrypted, "Decrypted text should match original")
}

// TestIntegrationWithRealOAuthSecret tests with a realistic OAuth client secret
func TestIntegrationWithRealOAuthSecret(t *testing.T) {
	// Load config
	cfg, err := config.Load("../../../config.yaml")
	require.NoError(t, err)

	key := cfg.GetIAMEncryptionKey()

	// Test with various OAuth client secret formats (using fake/test patterns)
	// These are NOT real secrets - just testing different string formats
	secrets := []string{
		"FAKE-GOOGLE-xxxxxxxxxxxxxxxxxxxxxxxx",                    // Google-like format
		"test_github_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",              // GitHub-like format
		"8Q~FAKE-MICROSOFT-xxxxxxxxxxxxxxxxxxxxx",                 // Microsoft-like format
		"fake.apple.secret.xxxxxxxxxxxxxxxxx.0.yyyyyyyyyyyyyyy",   // Apple-like format
		"test_fake_key_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // Generic long format
		"short-secret",                                            // Short format
		"a-very-long-oauth-client-secret-that-spans-many-characters-to-test-encryption-of-lengthy-strings", // Long format
	}

	for _, secret := range secrets {
		// Encrypt
		ciphertext, err := crypto.Encrypt(secret, key)
		require.NoError(t, err, "Should encrypt %s", secret)

		// Decrypt
		decrypted, err := crypto.Decrypt(ciphertext, key)
		require.NoError(t, err, "Should decrypt %s", secret)

		// Verify
		assert.Equal(t, secret, decrypted, "Round-trip should work for %s", secret)
	}
}
