# Task T14: OAuth Crypto Utilities and Encryption Key Management

**Story Reference:** US4-oauth-provider-configuration.md
**Type:** Backend Foundation
**Priority:** High
**Estimated Effort:** 4-6 hours
**Prerequisites:** None (can run in parallel with T13)

## Objective

Create shared crypto package with AES-256-GCM encryption/decryption for OAuth client secrets, implement encryption key configuration loading from environment, and add startup validation that exits if the key is invalid.

## Acceptance Criteria

- [ ] Crypto package created with Encrypt, Decrypt, ValidateKey functions
- [ ] Unit tests for crypto package passing with 100% coverage
- [ ] Encryption key configuration added to Config interface
- [ ] Environment variable `IAM_ENCRYPTION_KEY` loaded on startup
- [ ] Startup validation exits application if key is missing or invalid
- [ ] Encryption key injected via dependency injection (not global variable)
- [ ] Documentation explains AES-256-GCM choice vs Argon2

## Technical Requirements

### Crypto Package

**Location:** `internal/shared/crypto/`

**Why AES-256-GCM:**
- **Purpose:** Symmetric encryption for OAuth client secrets (Google/Github/Microsoft/Apple)
- **Requirement:** Must be able to decrypt secrets later to communicate with OAuth providers
- **NOT password hashing:** Argon2/bcrypt are one-way functions (cannot decrypt)
- **Industry standard:** Used by AWS, Google Cloud, Azure, Stripe, payment processors
- **Features:** Authenticated encryption (confidentiality + integrity), hardware-accelerated, NIST-approved

### Functions to Implement

```go
package crypto

// Encrypt encrypts plaintext using AES-256-GCM with the provided 32-byte key
// Returns base64-encoded ciphertext with nonce prepended
func Encrypt(plaintext string, key []byte) (string, error)

// Decrypt decrypts base64-encoded ciphertext using AES-256-GCM with the provided 32-byte key
// Returns plaintext or error if decryption fails
func Decrypt(ciphertext string, key []byte) (string, error)

// ValidateKey validates that the key is exactly 32 bytes (256 bits)
func ValidateKey(key []byte) error
```

### Encryption Key Configuration

**Config Interface Extension:**
```go
// config.go
type Config interface {
    // ... existing methods ...

    // GetIAMEncryptionKey returns the 32-byte encryption key for IAM secrets
    GetIAMEncryptionKey() []byte
}
```

**Implementation:**
```go
// internal/config/config.go
type AppConfig struct {
    // ... existing fields ...
    Security struct {
        AllowedOrigins   []string `yaml:"allowed_origins"`
        IAMEncryptionKey string   `yaml:"iam_encryption_key"`  // NEW
    } `yaml:"security"`
}

func (c *AppConfig) GetIAMEncryptionKey() []byte {
    // Environment variable takes precedence over config file
    if envKey := os.Getenv("IAM_ENCRYPTION_KEY"); envKey != "" {
        return []byte(envKey)
    }
    return []byte(c.Security.IAMEncryptionKey)
}
```

### Startup Validation

**Location:** `cmd/altalune/serve.go`

```go
func runServe(cmd *cobra.Command, args []string) error {
    // ... existing config loading ...

    // Validate encryption key before starting server
    if err := validateEncryptionKey(cfg, logger); err != nil {
        return err  // Exit application
    }

    // ... continue with server startup ...
}

func validateEncryptionKey(cfg altalune.Config, logger altalune.Logger) error {
    key := cfg.GetIAMEncryptionKey()

    if len(key) == 0 {
        logger.Error("IAM_ENCRYPTION_KEY is not set")
        return fmt.Errorf("IAM_ENCRYPTION_KEY environment variable or config setting is required")
    }

    if err := crypto.ValidateKey(key); err != nil {
        logger.Error("invalid IAM encryption key", "error", err)
        return fmt.Errorf("IAM_ENCRYPTION_KEY validation failed: %w", err)
    }

    logger.Info("IAM encryption key validated successfully")
    return nil
}
```

## Implementation Details

### crypto.go - Core Encryption Logic

```go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "fmt"
)

// Encrypt encrypts plaintext using AES-256-GCM with the provided key
// Returns base64-encoded ciphertext with nonce prepended
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
    nonce := make([]byte, gcm.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return "", fmt.Errorf("generate nonce: %w", err)
    }

    // Encrypt and authenticate
    // GCM.Seal prepends nonce to ciphertext
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

    // Return base64-encoded result for storage
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64-encoded ciphertext using AES-256-GCM
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
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]

    // Decrypt and verify
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", fmt.Errorf("decrypt: %w", err)
    }

    return string(plaintext), nil
}

// ValidateKey validates that the key is exactly 32 bytes (256 bits)
func ValidateKey(key []byte) error {
    if len(key) != 32 {
        return fmt.Errorf("encryption key must be exactly 32 bytes (256 bits), got %d bytes", len(key))
    }
    return nil
}
```

### crypto_test.go - Unit Tests

```go
package crypto_test

import (
    "crypto/rand"
    "testing"

    "github.com/hrz8/altalune/internal/shared/crypto"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

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

func TestEncryptDecrypt_LongPlaintext(t *testing.T) {
    key := make([]byte, 32)
    _, err := rand.Read(key)
    require.NoError(t, err)

    // Test with 1KB plaintext
    plaintext := string(make([]byte, 1024))

    ciphertext, err := crypto.Encrypt(plaintext, key)
    require.NoError(t, err)

    decrypted, err := crypto.Decrypt(ciphertext, key)
    require.NoError(t, err)
    assert.Equal(t, plaintext, decrypted)
}

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
    assert.Error(t, err) // Should fail authentication
}

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
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            key := make([]byte, tt.keySize)
            err := crypto.ValidateKey(key)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestEncrypt_InvalidKey(t *testing.T) {
    invalidKey := make([]byte, 16) // Only 16 bytes

    _, err := crypto.Encrypt("secret", invalidKey)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "32 bytes")
}

func TestDecrypt_InvalidBase64(t *testing.T) {
    key := make([]byte, 32)
    _, err := rand.Read(key)
    require.NoError(t, err)

    _, err = crypto.Decrypt("not-valid-base64!!!", key)
    assert.Error(t, err)
}

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
```

## Files to Create

- `internal/shared/crypto/crypto.go` - Core encryption/decryption functions
- `internal/shared/crypto/crypto_test.go` - Comprehensive unit tests

## Files to Modify

- `config.go` - Add `GetIAMEncryptionKey()` method to Config interface
- `internal/config/config.go` - Implement encryption key loading from env/config
- `cmd/altalune/serve.go` - Add startup validation for encryption key

## Testing Requirements

### Unit Tests

**Run tests:**
```bash
go test -v ./internal/shared/crypto
go test -v ./internal/shared/crypto -cover
```

**Coverage target:** 100%

**Test cases:**
- [x] Encrypt/decrypt round-trip with valid key
- [x] Long plaintext (1KB+)
- [x] Empty plaintext
- [x] Wrong key decryption fails
- [x] Invalid key sizes (16, 64, 0, 31, 33 bytes)
- [x] Invalid base64 ciphertext
- [x] Tampered ciphertext (GCM authentication)

### Integration Tests

**Test encryption key loading:**
```bash
# Test with environment variable
export IAM_ENCRYPTION_KEY="$(openssl rand -base64 32)"
./bin/app serve -c config.yaml
# Should start successfully

# Test with missing key
unset IAM_ENCRYPTION_KEY
./bin/app serve -c config.yaml
# Should exit with error

# Test with invalid key (wrong length)
export IAM_ENCRYPTION_KEY="short"
./bin/app serve -c config.yaml
# Should exit with error
```

## Commands to Run

```bash
# 1. Create crypto package directory
mkdir -p internal/shared/crypto

# 2. Generate test encryption key
openssl rand -base64 32

# 3. Set in environment
export IAM_ENCRYPTION_KEY="<generated-key>"

# 4. Run unit tests
go test -v ./internal/shared/crypto

# 5. Run with coverage
go test -v ./internal/shared/crypto -coverprofile=coverage.out
go tool cover -html=coverage.out

# 6. Test startup validation
./bin/app serve -c config.yaml
```

## Validation Checklist

- [ ] `crypto.go` implements Encrypt, Decrypt, ValidateKey
- [ ] Encryption uses AES-256-GCM (not ECB or CBC)
- [ ] Nonce is randomly generated (12 bytes for GCM)
- [ ] Ciphertext is base64-encoded for storage
- [ ] All unit tests pass
- [ ] Test coverage is 100%
- [ ] Config interface has GetIAMEncryptionKey method
- [ ] Environment variable takes precedence over config file
- [ ] Startup validation exits if key missing
- [ ] Startup validation exits if key invalid (not 32 bytes)
- [ ] Logger messages are clear and helpful

## Definition of Done

- [ ] Crypto package created with all functions
- [ ] Unit tests written and passing (100% coverage)
- [ ] Config interface extended
- [ ] Encryption key loading implemented (env + config)
- [ ] Startup validation implemented
- [ ] Application exits gracefully with clear error if key invalid
- [ ] Documentation in code explains AES-256-GCM choice
- [ ] No global variables (key passed via dependency injection)

## Dependencies

**Upstream:** None (can run in parallel with T13)

**Downstream:** T15 (Backend Domain) requires this crypto package

## Risk Factors

- **Medium Risk**: Encryption key management in production
  - **Mitigation**: Document production setup with secrets manager
  - **Mitigation**: Validate key on startup, fail fast

- **Low Risk**: Incorrect AES mode selection
  - **Mitigation**: Use GCM mode (authenticated encryption)
  - **Mitigation**: Unit tests verify encryption/decryption

- **Low Risk**: Key rotation not implemented
  - **Mitigation**: Document key rotation process for future
  - **Mitigation**: Note: Out of scope for US4

## Notes

### Encryption Key Format

**What it is:**
- 32-byte (256-bit) symmetric key
- Generated with: `openssl rand -base64 32`
- Output: `vK8s2R7pN4jF9mT3xQ1wL6hY0dC5aE8b2Z9vM4nG7rJ=` (44 chars)
- **NOT a .pem file** - this is symmetric encryption, not certificates
- **NOT stored in database** - only encrypted data in DB

**Where it lives:**
- ✅ Environment variable (production): `IAM_ENCRYPTION_KEY`
- ✅ config.yaml (development only): `security.iam_encryption_key`
- ❌ NOT in git, NOT in database, NOT in logs

**Production Setup:**
```bash
# Store in AWS Secrets Manager / Google Secret Manager / Azure Key Vault
# Application reads from environment variable
# Rotate periodically (requires re-encryption migration)
```

### Why AES-256-GCM vs Argon2

| Feature | AES-256-GCM | Argon2 |
|---------|-------------|--------|
| **Type** | Symmetric encryption | Password hashing |
| **Direction** | Two-way (encrypt + decrypt) | One-way (hash only) |
| **Use case** | OAuth secrets (need to decrypt) | User passwords (never decrypt) |
| **Output** | Ciphertext (reversible) | Hash (irreversible) |
| **For OAuth** | ✅ Correct | ❌ Wrong (can't decrypt) |

### GCM Mode Benefits

**Authenticated Encryption:**
- Provides confidentiality (encryption)
- Provides integrity (authentication tag)
- Detects tampering (GCM.Open fails if modified)

**Standard:**
- NIST approved
- Used by TLS 1.3
- Hardware-accelerated (AES-NI)

### Security Best Practices

1. **Key Storage:**
   - Never commit to git
   - Never log the key
   - Use secrets manager in production
   - Rotate periodically

2. **Key Validation:**
   - Validate on startup (fail fast)
   - Exact 32 bytes required
   - Clear error messages

3. **Encryption:**
   - Random nonce per encryption
   - GCM for authenticated encryption
   - Base64 encode for storage

4. **Decryption:**
   - Only on explicit request (RevealClientSecret)
   - Never log decrypted secrets
   - Auto-hide in frontend (30 seconds)

### Future Enhancements (Out of Scope)

- Key rotation (re-encrypt all secrets with new key)
- Multiple key versions (during rotation)
- Audit logging (who decrypted what when)
- Key derivation (KDF from master key)
