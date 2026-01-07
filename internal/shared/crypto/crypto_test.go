package crypto_test

import (
	"crypto/rand"
	"testing"

	"github.com/hrz8/altalune/internal/shared/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEncryptDecrypt verifies basic encrypt/decrypt round-trip with a valid key
func TestEncryptDecrypt(t *testing.T) {
	// Generate valid 32-byte key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := "my-secret-oauth-client-secret"

	// Encrypt
	ciphertext, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err)
	assert.NotEmpty(t, ciphertext)
	assert.NotEqual(t, plaintext, ciphertext)

	// Decrypt
	decrypted, err := crypto.Decrypt(ciphertext, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

// TestEncryptDecrypt_LongPlaintext tests encryption with 1KB+ plaintext
func TestEncryptDecrypt_LongPlaintext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	// Test with 1KB plaintext
	plaintextBytes := make([]byte, 1024)
	_, err = rand.Read(plaintextBytes)
	require.NoError(t, err)
	plaintext := string(plaintextBytes)

	ciphertext, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err)

	decrypted, err := crypto.Decrypt(ciphertext, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

// TestEncryptDecrypt_EmptyPlaintext tests encryption with empty string
func TestEncryptDecrypt_EmptyPlaintext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := ""

	ciphertext, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err)

	decrypted, err := crypto.Decrypt(ciphertext, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

// TestEncrypt_WrongKey verifies that decryption fails with wrong key
func TestEncrypt_WrongKey(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := "secret"
	ciphertext, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err)

	// Try to decrypt with different key
	wrongKey := make([]byte, 32)
	_, err = rand.Read(wrongKey)
	require.NoError(t, err)

	_, err = crypto.Decrypt(ciphertext, wrongKey)
	assert.Error(t, err) // Should fail GCM authentication
}

// TestValidateKey tests key validation with various key sizes
func TestValidateKey(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
		wantErr bool
	}{
		{"Valid 32 bytes", 32, false},
		{"Invalid 16 bytes", 16, true},
		{"Invalid 64 bytes", 64, true},
		{"Invalid 0 bytes", 0, true},
		{"Invalid 31 bytes", 31, true},
		{"Invalid 33 bytes", 33, true},
		{"Invalid 1 byte", 1, true},
		{"Invalid 24 bytes", 24, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := make([]byte, tt.keySize)
			err := crypto.ValidateKey(key)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "32 bytes")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestEncrypt_InvalidKey tests encryption with invalid key size
func TestEncrypt_InvalidKey(t *testing.T) {
	invalidKey := make([]byte, 16) // Only 16 bytes

	_, err := crypto.Encrypt("secret", invalidKey)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")
}

// TestDecrypt_InvalidKey tests decryption with invalid key size
func TestDecrypt_InvalidKey(t *testing.T) {
	invalidKey := make([]byte, 16) // Only 16 bytes

	_, err := crypto.Decrypt("some-ciphertext", invalidKey)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")
}

// TestDecrypt_InvalidBase64 tests decryption with invalid base64 string
func TestDecrypt_InvalidBase64(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	_, err = crypto.Decrypt("not-valid-base64!!!", key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode base64")
}

// TestDecrypt_TamperedCiphertext tests that GCM authentication detects tampering
func TestDecrypt_TamperedCiphertext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := "secret"
	ciphertext, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err)

	// Tamper with ciphertext (change last character)
	tampered := ciphertext[:len(ciphertext)-1] + "X"

	_, err = crypto.Decrypt(tampered, key)
	assert.Error(t, err) // GCM authentication should fail
}

// TestDecrypt_ShortCiphertext tests decryption with ciphertext shorter than nonce
func TestDecrypt_ShortCiphertext(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	// Create a very short base64 string (will decode to less than 12 bytes)
	shortCiphertext := "YWJj" // "abc" in base64 (3 bytes)

	_, err = crypto.Decrypt(shortCiphertext, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ciphertext too short")
}

// TestEncrypt_DifferentCiphertexts verifies each encryption produces unique ciphertext
// This is important because nonce should be random for each encryption
func TestEncrypt_DifferentCiphertexts(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := "same-plaintext"

	// Encrypt same plaintext twice
	ciphertext1, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err)

	ciphertext2, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err)

	// Ciphertexts should be different (different random nonces)
	assert.NotEqual(t, ciphertext1, ciphertext2)

	// Both should decrypt to same plaintext
	decrypted1, err := crypto.Decrypt(ciphertext1, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted1)

	decrypted2, err := crypto.Decrypt(ciphertext2, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted2)
}

// TestEncrypt_SpecialCharacters tests encryption with special characters
func TestEncrypt_SpecialCharacters(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := "🔐 special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?"

	ciphertext, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err)

	decrypted, err := crypto.Decrypt(ciphertext, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

// TestEncrypt_Unicode tests encryption with Unicode characters
func TestEncrypt_Unicode(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := "こんにちは世界 مرحبا بالعالم Привет мир"

	ciphertext, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err)

	decrypted, err := crypto.Decrypt(ciphertext, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

// TestEncrypt_Newlines tests encryption with newlines and whitespace
func TestEncrypt_Newlines(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	plaintext := "line1\nline2\r\nline3\ttabbed"

	ciphertext, err := crypto.Encrypt(plaintext, key)
	require.NoError(t, err)

	decrypted, err := crypto.Decrypt(ciphertext, key)
	require.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

// TestValidateKey_NilKey tests validation with nil key
func TestValidateKey_NilKey(t *testing.T) {
	var key []byte
	err := crypto.ValidateKey(key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")
}

// TestEncrypt_NilKey tests encryption with nil key
func TestEncrypt_NilKey(t *testing.T) {
	var key []byte
	_, err := crypto.Encrypt("secret", key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")
}

// TestDecrypt_NilKey tests decryption with nil key
func TestDecrypt_NilKey(t *testing.T) {
	var key []byte
	_, err := crypto.Decrypt("some-ciphertext", key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")
}
