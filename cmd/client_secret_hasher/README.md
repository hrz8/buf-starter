# Client Secret Hasher

A standalone utility for hashing OAuth client secrets using Argon2id.

## Overview

This utility hashes client secrets using the Argon2id algorithm (winner of the Password Hashing Competition 2015). It's used for OAuth client secret management in the Altalune project.

## Security

- **Algorithm**: Argon2id (memory-hard, GPU/ASIC resistant)
- **Default Parameters**:
  - Iterations (time cost): 2
  - Memory: 64MB (65536 KB)
  - Parallelism: 4 threads
  - Hash length: 32 bytes
- **Output Format**: PHC string format (`$argon2id$v=19$m=65536,t=2,p=4$<salt>$<hash>`)
- **Minimum Secret Length**: 32 characters (enforced)

## Building

```bash
# Build just this utility
make client-secret-hasher

# Build all utilities (publicid, secret_encrypter, client_secret_hasher)
make build-utils

# The binary will be in ./bin/client_secret_hasher
```

## Usage

### Basic Usage (Command Line Flag)

```bash
./bin/client_secret_hasher --secret "your-32-character-minimum-secret-here"
```

### Stdin Pipe

```bash
echo "your-32-character-minimum-secret-here" | ./bin/client_secret_hasher
```

### Interactive Mode

```bash
./bin/client_secret_hasher
# Prompts: Enter client secret (min 32 chars):
```

### Custom Parameters (Higher Security)

```bash
./bin/client_secret_hasher \
  --secret "your-secret" \
  --iterations 3 \
  --memory 131072 \
  --threads 8 \
  --length 64
```

## Flags

- `--secret` (string): Client secret to hash (optional, can use stdin)
- `--iterations` (uint): Time cost, default: 2
- `--memory` (uint): Memory cost in KB, default: 65536 (64MB)
- `--threads` (uint): Parallelism, default: 4
- `--length` (uint): Hash length in bytes, default: 32

## Examples

### Generate hash for database insertion

```bash
./bin/client_secret_hasher --secret "my-oauth-client-secret-32chars-min"
# Output: $argon2id$v=19$m=65536,t=2,p=4$<salt>$<hash>
```

### Pipe from password generator

```bash
openssl rand -base64 32 | ./bin/client_secret_hasher
```

### Store hash in variable

```bash
HASH=$(./bin/client_secret_hasher --secret "my-secret-32chars-minimum-length")
echo "INSERT INTO oauth_clients (client_secret_hash) VALUES ('$HASH');"
```

## Output Format

The output is a PHC (Password Hashing Competition) string format:

```
$argon2id$v=19$m=65536,t=2,p=4$CO5hu/iRl5ey1rr8h4FbRQ$qgk/PEQzuAdh4b06CmxTS/djb7F7Fojdhubl0QEKWQw
^         ^    ^              ^                      ^
algorithm vers params         base64(salt)           base64(hash)
```

Components:
- **$argon2id**: Algorithm identifier (Argon2id variant)
- **v=19**: Argon2 version
- **m=65536,t=2,p=4**: Memory (64MB), Time cost (2), Parallelism (4)
- **CO5hu...**: Base64-encoded salt (16 bytes)
- **qgk/P...**: Base64-encoded hash (32 bytes)

## Error Handling

### Secret too short

```bash
./bin/client_secret_hasher --secret "short"
# Error: client secret must be at least 32 characters, got 5
```

### Empty secret

```bash
echo "" | ./bin/client_secret_hasher
# Error: stdin provided empty client secret
```

## Performance

With default parameters (t=2, m=64MB, p=4):
- **Hashing time**: ~50-100ms on modern hardware
- **Memory usage**: 64MB during hashing
- **CPU cores**: Uses 4 threads

## Use Cases

1. **Manual database updates**: Hash secrets for existing OAuth clients
2. **Testing**: Generate test hashes for development
3. **Scripts**: Automate client creation with pre-hashed secrets
4. **Debugging**: Verify hash generation and format
5. **Migration**: Update secrets from old hashing methods

## Comparison with bcrypt

| Feature | bcrypt (old) | Argon2id (new) |
|---------|-------------|----------------|
| Algorithm | Blowfish-based | Memory-hard |
| Memory | ~4KB | 64MB (configurable) |
| Resistance | CPU attacks | CPU + GPU + ASIC |
| Standard | 1999 | 2015 (PHC winner) |
| Performance | ~250ms | ~50-100ms |
| OWASP | Acceptable | Recommended |

## Integration with Password Package

The binary uses `internal/shared/password` package:

```go
import "github.com/hrz8/altalune/internal/shared/password"

// Hash a password
hash, err := password.HashPassword(secret, password.HashOption{
    Iterations: 2,
    Memory:     64 * 1024,
    Threads:    4,
    Len:        32,
})

// Verify a password
match, err := password.VerifyPassword(plaintext, hash)
```

## Security Best Practices

1. **Minimum Length**: Always enforce 32+ character secrets
2. **Unique Salts**: Each hash uses a unique random salt (automatic)
3. **Secure Storage**: Store hashes, never plaintext secrets
4. **Audit Logging**: Log when secrets are revealed (not when hashed)
5. **Production Parameters**: Use default or higher for production

## Related Utilities

- **publicid**: Generate NanoID-based public IDs
- **secret_encrypter**: Encrypt OAuth provider secrets (AES-256-GCM)
- **client_secret_hasher**: Hash OAuth client secrets (Argon2id) - this utility

## References

- [Argon2 Specification](https://github.com/P-H-C/phc-winner-argon2)
- [OWASP Password Storage](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [PHC String Format](https://github.com/P-H-C/phc-string-format)
