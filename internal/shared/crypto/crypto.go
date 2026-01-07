package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

// Why AES-256-GCM:
//
// AES-256-GCM (Galois/Counter Mode) is used for symmetric encryption of OAuth client secrets.
// This is the correct choice because:
//
// 1. **Reversible Encryption**: We need to decrypt OAuth client secrets later to communicate
//    with OAuth providers (Google, Github, Microsoft, Apple). This rules out one-way hash
//    functions like Argon2 or bcrypt, which are designed for passwords that never need decryption.
//
// 2. **Authenticated Encryption**: GCM mode provides both confidentiality (encryption) and
//    integrity (authentication). It detects tampering - if ciphertext is modified, decryption
//    will fail with an authentication error.
//
// 3. **Industry Standard**: AES-256-GCM is used by AWS, Google Cloud, Azure, Stripe, and
//    payment processors for encrypting sensitive data at rest. It's NIST-approved and
//    considered the gold standard for symmetric encryption.
//
// 4. **Hardware Acceleration**: Modern CPUs have AES-NI instructions that make AES-GCM
//    extremely fast (hardware-accelerated), making it suitable for production use.
//
// Comparison with Argon2:
// - Argon2: Password hashing (one-way, cannot decrypt) - used for user passwords
// - AES-256-GCM: Symmetric encryption (two-way, can decrypt) - used for OAuth secrets
//
// For OAuth client secrets, we MUST use AES-256-GCM because we need to retrieve the plaintext
// secret to authenticate with OAuth providers during the OAuth flow.

// Encrypt encrypts plaintext using AES-256-GCM with the provided 32-byte key.
// Returns base64-encoded ciphertext with nonce prepended.
//
// The encryption process:
// 1. Validates the key is exactly 32 bytes (256 bits)
// 2. Creates an AES cipher block
// 3. Wraps it in GCM mode for authenticated encryption
// 4. Generates a random 12-byte nonce
// 5. Encrypts and authenticates the plaintext
// 6. Prepends the nonce to the ciphertext (needed for decryption)
// 7. Base64-encodes the result for safe storage
//
// Example:
//   key := make([]byte, 32) // 32-byte encryption key from config
//   encrypted, err := Encrypt("my-oauth-secret", key)
//   if err != nil {
//       log.Fatal(err)
//   }
//   // Store encrypted in database
func Encrypt(plaintext string, key []byte) (string, error) {
	// Validate key first
	if err := ValidateKey(key); err != nil {
		return "", err
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	// Create GCM mode (Galois/Counter Mode for authenticated encryption)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	// Generate random nonce (12 bytes for GCM)
	// Nonce must be unique for each encryption with the same key
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	// Encrypt and authenticate
	// GCM.Seal prepends nonce to ciphertext
	// nil means we allocate a new slice for the result
	// nonce is used for encryption and will be prepended to output
	// []byte(plaintext) is the data to encrypt
	// nil means no additional authenticated data (AAD)
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return base64-encoded result for safe storage in database
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64-encoded ciphertext using AES-256-GCM with the provided 32-byte key.
// Returns plaintext or error if decryption/authentication fails.
//
// The decryption process:
// 1. Validates the key is exactly 32 bytes (256 bits)
// 2. Base64-decodes the ciphertext
// 3. Creates an AES cipher block
// 4. Wraps it in GCM mode
// 5. Extracts the nonce from the beginning of the ciphertext
// 6. Decrypts and verifies authentication
// 7. Returns the plaintext
//
// Security notes:
// - If the ciphertext has been tampered with, authentication will fail
// - If the wrong key is used, decryption will fail
// - If the ciphertext is corrupted, decryption will fail
//
// Example:
//   key := make([]byte, 32) // Same 32-byte key used for encryption
//   plaintext, err := Decrypt(encryptedFromDB, key)
//   if err != nil {
//       log.Fatal(err) // Decryption failed
//   }
//   // Use plaintext for OAuth provider communication
func Decrypt(ciphertext string, key []byte) (string, error) {
	// Validate key first
	if err := ValidateKey(key); err != nil {
		return "", err
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	// Verify nonce size
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	// Nonce was prepended during encryption
	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]

	// Decrypt and verify authentication
	// nil means we allocate a new slice for the result
	// nonce is extracted from the ciphertext
	// ciphertextBytes is the encrypted data
	// nil means no additional authenticated data (AAD)
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(plaintext), nil
}

// ValidateKey validates that the key is exactly 32 bytes (256 bits).
// AES-256 requires a 256-bit (32-byte) key.
//
// This is called:
// - On application startup to validate the encryption key from config
// - Before every encryption/decryption operation
//
// The key must be:
// - Exactly 32 bytes (256 bits) - no more, no less
// - Generated with a cryptographically secure random number generator
// - Stored securely (never committed to git, never logged)
//
// To generate a valid key:
//   openssl rand -base64 32
//
// Example output (44 characters when base64-encoded):
//   vK8s2R7pN4jF9mT3xQ1wL6hY0dC5aE8b2Z9vM4nG7rJ=
//
// This 44-character base64 string decodes to exactly 32 bytes.
func ValidateKey(key []byte) error {
	if len(key) != 32 {
		return fmt.Errorf("encryption key must be exactly 32 bytes (256 bits), got %d bytes", len(key))
	}
	return nil
}
