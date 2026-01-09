# Task T19: Argon2 Password Hashing Integration + CLI Command

**Story Reference:** US6-oauth-client-management.md
**Type:** Backend Foundation
**Priority:** High (P0 - Must be first)
**Estimated Effort:** 3-4 hours
**Prerequisites:** None (foundational task)

## Objective

Replace bcrypt with Argon2id for OAuth client secret hashing and provide a CLI utility for manual secret hashing.

## Acceptance Criteria

- [ ] Argon2 password hashing package created at `internal/shared/password/`
- [ ] HashPassword function implemented with production parameters
- [ ] VerifyPassword function implemented with constant-time comparison
- [ ] CLI command `altalune hash <secret>` created for manual hashing
- [ ] Unit tests written for password hashing functions
- [ ] CLI command registered in main.go
- [ ] Output format follows PHC string standard

## Technical Requirements

### Argon2 Package Structure

Create `internal/shared/password/` package with 4 files:

**1. password.go** - Types, constants, errors:
```go
package password

import "errors"

var (
    ErrInvalidHashedString = errors.New("the encoded hash is not in the correct format")
    ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

const (
    saltLen = 16  // 16 bytes for salt
    threads = 4   // Parallelism
    splitN  = 6   // PHC string format parts
)

type HashOption struct {
    Iterations uint32  // Time cost
    Memory     uint32  // Memory cost in KB
    Threads    uint8   // Parallelism
    Len        uint32  // Output hash length
}
```

**2. hash.go** - Password hashing:
```go
func generateRandomBytes(n uint32) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)  // crypto/rand for secure randomness
    return b, err
}

func HashPassword(plainText string, opt HashOption) (string, error) {
    // 1. Generate 16-byte cryptographically secure salt
    // 2. Hash with argon2.IDKey()
    // 3. Base64 encode salt and hash
    // 4. Return PHC format: $argon2id$v=19$m=65536,t=2,p=4$<salt>$<hash>
}
```

**3. verify.go** - Password verification:
```go
func decodeHash(encodedHash string) (*HashOption, []byte, []byte, error) {
    // Parse PHC string format
    // Extract parameters, salt, and hash
    // Validate version compatibility
}

func VerifyPassword(password, encodedHash string) (bool, error) {
    // 1. Decode PHC string to get parameters, salt, hash
    // 2. Re-hash password with same parameters and salt
    // 3. Constant-time comparison using crypto/subtle
}
```

**4. password_test.go** - Unit tests:
```go
func TestHashPassword(t *testing.T) {
    // Test with production parameters
    // Validate output format
    // Ensure hash length correct
}

func TestVerifyPassword(t *testing.T) {
    // Test matching password
    // Test non-matching password
    // Test invalid hash format
}
```

### Production Parameters

**Recommended Settings** (balanced security and performance):
```go
HashOption{
    Iterations: 2,          // Time cost (t)
    Memory:     64 * 1024,  // 64MB memory cost (m)
    Threads:    4,          // Parallelism (p)
    Len:        32,         // 32-byte output hash
}
```

**Performance**: ~50-100ms on modern hardware (acceptable for client creation)

**Security**: OWASP recommended, memory-hard algorithm resistant to GPU/ASIC attacks

### CLI Command Structure

**File: `cmd/altalune/hash.go`**

```go
package main

import (
    "fmt"
    "github.com/hrz8/altalune/internal/shared/password"
    "github.com/spf13/cobra"
)

func NewHashCommand(rootCmd *cobra.Command) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "hash [secret]",
        Short: "Hash a secret using Argon2id",
        Long:  "Generate an Argon2id hash for OAuth client secrets",
        Args:  cobra.ExactArgs(1),
        RunE:  runHash,
    }

    // Configurable parameters (optional flags)
    cmd.Flags().Uint32("iterations", 2, "Time cost (iterations)")
    cmd.Flags().Uint32("memory", 64*1024, "Memory cost in KB")
    cmd.Flags().Uint8("threads", 4, "Parallelism (threads)")
    cmd.Flags().Uint32("length", 32, "Hash length in bytes")

    return cmd
}

func runHash(cmd *cobra.Command, args []string) error {
    secret := args[0]

    // Validate minimum length
    if len(secret) < 32 {
        return fmt.Errorf("secret must be at least 32 characters, got %d", len(secret))
    }

    // Get flags
    iterations, _ := cmd.Flags().GetUint32("iterations")
    memory, _ := cmd.Flags().GetUint32("memory")
    threads, _ := cmd.Flags().GetUint8("threads")
    length, _ := cmd.Flags().GetUint32("length")

    // Hash the secret
    hash, err := password.HashPassword(secret, password.HashOption{
        Iterations: iterations,
        Memory:     memory,
        Threads:    threads,
        Len:        length,
    })
    if err != nil {
        return fmt.Errorf("failed to hash secret: %w", err)
    }

    fmt.Println("Hashed secret:")
    fmt.Println(hash)

    return nil
}
```

**Register in `cmd/altalune/main.go`**:
```go
func registerCommands(cmd *cobra.Command) {
    cmd.AddCommand(
        NewServeCommand(cmd),
        NewMigrateCommand(cmd),
        NewHashCommand(cmd),  // Add this line
    )
}
```

## Implementation Details

### Copy from Reference Implementation

Source: `/Users/hirzi/src/altalune-id/hak/password/`

Files to copy (with adaptation):
- `password.go` → `internal/shared/password/password.go`
- `hash.go` → `internal/shared/password/hash.go`
- `verify.go` → `internal/shared/password/verify.go`
- `password_test.go` → `internal/shared/password/password_test.go`

**Key Changes**:
- Update package name to `password`
- Update import paths for altalune project
- Keep production parameters (iterations=2, memory=64KB)

### PHC String Format

**Output Example**:
```
$argon2id$v=19$m=65536,t=2,p=4$CO5hu/iRl5ey1rr8h4FbRQ$qgk/PEQzuAdh4b06CmxTS/djb7F7Fojdhubl0QEKWQw
^         ^    ^              ^                      ^
algorithm vers params         base64(salt)           base64(hash)
```

**Format Breakdown**:
- `$argon2id` - Algorithm identifier
- `$v=19` - Argon2 version
- `$m=65536,t=2,p=4` - Parameters (memory, iterations, parallelism)
- `$<base64-salt>` - Base64-encoded salt (16 bytes)
- `$<base64-hash>` - Base64-encoded hash (32 bytes)

## Files to Create

- `internal/shared/password/password.go`
- `internal/shared/password/hash.go`
- `internal/shared/password/verify.go`
- `internal/shared/password/password_test.go`
- `cmd/altalune/hash.go`

## Files to Modify

- `cmd/altalune/main.go` - Register hash command

## Testing Requirements

### Unit Tests

```go
// Test hash generation
func TestHashPassword(t *testing.T) {
    opt := password.HashOption{
        Iterations: 2,
        Memory:     64 * 1024,
        Threads:    4,
        Len:        32,
    }

    hash, err := password.HashPassword("test-secret-32-characters-minimum", opt)
    assert.NoError(t, err)
    assert.Contains(t, hash, "$argon2id$v=19$")
    assert.Contains(t, hash, "m=65536,t=2,p=4")
}

// Test verification (matching)
func TestVerifyPassword_Match(t *testing.T) {
    secret := "test-secret-32-characters-minimum"
    hash, _ := password.HashPassword(secret, opt)

    match, err := password.VerifyPassword(secret, hash)
    assert.NoError(t, err)
    assert.True(t, match)
}

// Test verification (non-matching)
func TestVerifyPassword_NoMatch(t *testing.T) {
    secret := "test-secret-32-characters-minimum"
    hash, _ := password.HashPassword(secret, opt)

    match, err := password.VerifyPassword("wrong-secret", hash)
    assert.NoError(t, err)
    assert.False(t, match)
}

// Test invalid hash format
func TestVerifyPassword_InvalidFormat(t *testing.T) {
    _, err := password.VerifyPassword("test", "invalid-hash")
    assert.Error(t, err)
    assert.Equal(t, password.ErrInvalidHashedString, err)
}
```

### Manual Testing

```bash
# Build binary
make build

# Test CLI with minimum length secret
./bin/app hash "test-secret-32-characters-minimum"

# Expected output:
# Hashed secret:
# $argon2id$v=19$m=65536,t=2,p=4$...

# Test with custom parameters
./bin/app hash "test-secret" --iterations 3 --memory 131072

# Test with too-short secret (should fail)
./bin/app hash "short"
# Expected error: secret must be at least 32 characters
```

## Commands to Run

```bash
# Create directory
mkdir -p internal/shared/password

# Copy files from reference implementation
cp /Users/hirzi/src/altalune-id/hak/password/*.go internal/shared/password/

# Update package name and imports
# (Manual editing required)

# Run tests
go test ./internal/shared/password/

# Build binary
make build

# Test CLI command
./bin/app hash "test-secret-32-characters-minimum-length"
```

## Validation Checklist

- [ ] Argon2 package compiles without errors
- [ ] Unit tests pass for HashPassword
- [ ] Unit tests pass for VerifyPassword
- [ ] CLI command registered in main.go
- [ ] CLI command accepts secret argument
- [ ] CLI command validates minimum 32 characters
- [ ] CLI command outputs valid PHC string format
- [ ] CLI flags work (--iterations, --memory, --threads, --length)
- [ ] Hash verification works with constant-time comparison
- [ ] Performance acceptable (<100ms for hashing)

## Definition of Done

- [ ] Password package created with all 4 files
- [ ] HashPassword implemented with Argon2id
- [ ] VerifyPassword implemented with constant-time comparison
- [ ] CLI command `altalune hash` implemented
- [ ] CLI command registered and working
- [ ] Unit tests written and passing
- [ ] Manual testing completed successfully
- [ ] Output format follows PHC string standard
- [ ] Code documented with inline comments
- [ ] Performance tested (<100ms hashing time)

## Dependencies

**External**:
- `golang.org/x/crypto/argon2` - Argon2 implementation
- `crypto/rand` - Secure random number generation
- `crypto/subtle` - Constant-time comparison
- `github.com/spf13/cobra` - CLI framework (existing)

**Internal**:
- None (this is a foundational task)

## Risk Factors

- **Medium Risk**: Performance impact if parameters too aggressive
  - **Mitigation**: Use tested parameters (iterations=2, memory=64KB)
- **Low Risk**: PHC string format compatibility
  - **Mitigation**: Follow standard format, test parsing

## Notes

### Why Argon2id over bcrypt?

1. **Modern Standard**: Won Password Hashing Competition 2015
2. **Memory-Hard**: Better GPU/ASIC resistance than bcrypt
3. **OWASP Recommended**: Current best practice
4. **Performance**: Tunable parameters for security/speed balance
5. **User Request**: Explicitly requested in breakdown instructions

### Why CLI Command?

- **Manual Migration**: Update existing bcrypt hashes to Argon2
- **Database Updates**: Hash secrets for manual database insertion
- **Debugging**: Verify hash generation and format
- **Testing**: Generate test hashes for development

### Migration from bcrypt

**No backward compatibility needed**:
- Product not yet released
- Can fully deprecate bcrypt
- All new secrets use Argon2id from day 1

**Future Consideration**:
- If backward compatibility needed later, detect hash type:
  ```go
  if strings.HasPrefix(hash, "$argon2") {
      return password.VerifyPassword(plaintext, hash)
  }
  // Fallback to bcrypt for legacy hashes
  return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
  ```

### Security Considerations

- **Salt Generation**: Uses `crypto/rand` (cryptographically secure)
- **Constant-Time Comparison**: Prevents timing attacks
- **Memory-Hard**: Resistant to parallel cracking attempts
- **Configurable**: Can increase security if needed in future

### Reference Implementation

Source: `/Users/hirzi/src/altalune-id/hak/password`

All code adapted from this proven implementation. Only changes:
- Package name
- Import paths
- Production parameter defaults
