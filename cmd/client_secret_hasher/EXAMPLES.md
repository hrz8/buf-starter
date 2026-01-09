# Client Secret Hasher - Usage Examples

## Quick Start

### Basic Usage

```bash
# Hash a client secret
./bin/client_secret_hasher --secret "your-32-character-minimum-secret-here"

# Output:
# $argon2id$v=19$m=65536,t=2,p=4$CO5hu/iRl5ey1rr8h4FbRQ$qgk/PEQzuAdh4b06CmxTS/djb7F7Fojdhubl0QEKWQw
```

## Common Workflows

### 1. Generate and Store Hash for Database

```bash
# Generate hash
HASH=$(./bin/client_secret_hasher --secret "my-oauth-client-secret-32chars-minimum")

# Insert into database
psql -d altalune << EOF
INSERT INTO altalune_oauth_clients (
    project_id, public_id, name, client_id,
    client_secret_hash, redirect_uris, pkce_required
) VALUES (
    1,
    'abc123',
    'My OAuth Client',
    '550e8400-e29b-41d4-a716-446655440000',
    '$HASH',
    ARRAY['http://localhost:3000/callback'],
    true
);
EOF
```

### 2. Update Existing Client Secret

```bash
# Generate new hash
NEW_HASH=$(./bin/client_secret_hasher --secret "new-secret-32-characters-minimum-length")

# Update database
psql -d altalune -c "
UPDATE altalune_oauth_clients
SET client_secret_hash = '$NEW_HASH',
    updated_at = NOW()
WHERE public_id = 'client123';
"
```

### 3. Generate Multiple Hashes for Batch Operations

```bash
# Create a list of secrets
cat > secrets.txt << EOF
client1-secret-32chars-minimum-length-here
client2-secret-32chars-minimum-length-here
client3-secret-32chars-minimum-length-here
EOF

# Hash each secret
while IFS= read -r secret; do
    echo "Secret: $secret"
    echo "Hash: $(./bin/client_secret_hasher --secret "$secret")"
    echo "---"
done < secrets.txt
```

### 4. Use with Password Generator

```bash
# Generate random secret and hash it
SECRET=$(openssl rand -base64 48 | cut -c1-32)
echo "Generated secret: $SECRET"
HASH=$(./bin/client_secret_hasher --secret "$SECRET")
echo "Generated hash: $HASH"

# Store both (secret shown once, hash stored in DB)
echo "$SECRET" > /secure/client-secret-backup.txt
chmod 600 /secure/client-secret-backup.txt
```

### 5. Script for OAuth Client Creation

```bash
#!/bin/bash

# generate-oauth-client.sh
CLIENT_NAME="${1:-Default Client}"
SECRET=$(openssl rand -base64 48 | cut -c1-32)
CLIENT_ID=$(uuidgen | tr '[:upper:]' '[:lower:]')
HASH=$(./bin/client_secret_hasher --secret "$SECRET")

echo "=== OAuth Client Created ==="
echo "Name: $CLIENT_NAME"
echo "Client ID: $CLIENT_ID"
echo "Client Secret (SAVE THIS): $SECRET"
echo "Client Secret Hash: $HASH"
echo
echo "Database insertion:"
echo "INSERT INTO altalune_oauth_clients (name, client_id, client_secret_hash)"
echo "VALUES ('$CLIENT_NAME', '$CLIENT_ID', '$HASH');"
```

## Advanced Usage

### Higher Security Parameters

```bash
# More secure (slower, more memory)
./bin/client_secret_hasher \
    --secret "my-secret-32chars-minimum-length" \
    --iterations 4 \
    --memory 262144 \
    --threads 8 \
    --length 64

# Output includes higher parameters:
# $argon2id$v=19$m=262144,t=4,p=8$...
```

### Testing Hash Verification

```bash
# Generate hash
SECRET="test-secret-32-characters-minimum"
HASH=$(./bin/client_secret_hasher --secret "$SECRET")

# Write test verification script
cat > test_verify.go << 'GOEOF'
package main

import (
    "fmt"
    "os"
    "github.com/hrz8/altalune/internal/shared/password"
)

func main() {
    secret := os.Args[1]
    hash := os.Args[2]

    match, err := password.VerifyPassword(secret, hash)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    if match {
        fmt.Println("✓ Secret matches hash")
        os.Exit(0)
    }

    fmt.Println("✗ Secret does not match hash")
    os.Exit(1)
}
GOEOF

# Run verification
go run test_verify.go "$SECRET" "$HASH"
# Output: ✓ Secret matches hash

rm test_verify.go
```

### Pipe from Environment Variable

```bash
# Store secret in environment
export OAUTH_CLIENT_SECRET="my-secure-32-character-minimum-secret"

# Hash from environment
echo "$OAUTH_CLIENT_SECRET" | ./bin/client_secret_hasher

# Or directly
./bin/client_secret_hasher --secret "$OAUTH_CLIENT_SECRET"
```

### Interactive Mode (Secure Input)

```bash
# Run without arguments
./bin/client_secret_hasher

# Prompts:
# Enter client secret (min 32 chars): [type here]
# [secret not echoed in terminal]
# Output: $argon2id$v=19$m=65536,t=2,p=4$...
```

## Integration Examples

### With Docker

```bash
# Dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN make client-secret-hasher

# Use in entrypoint
docker run myapp ./bin/client_secret_hasher --secret "$CLIENT_SECRET"
```

### With Kubernetes Secret

```bash
# Create Kubernetes secret
kubectl create secret generic oauth-client-secret \
  --from-literal=plaintext="$SECRET" \
  --from-literal=hash="$(./bin/client_secret_hasher --secret "$SECRET")"

# Use in deployment
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    env:
    - name: OAUTH_CLIENT_SECRET_HASH
      valueFrom:
        secretKeyRef:
          name: oauth-client-secret
          key: hash
```

### With Ansible

```yaml
# ansible playbook
- name: Generate OAuth client secret hash
  shell: |
    ./bin/client_secret_hasher --secret "{{ oauth_client_secret }}"
  register: secret_hash

- name: Insert OAuth client to database
  postgresql_query:
    query: |
      INSERT INTO altalune_oauth_clients (client_secret_hash)
      VALUES ('{{ secret_hash.stdout }}');
```

## Troubleshooting

### Secret Too Short

```bash
./bin/client_secret_hasher --secret "short"
# Error: client secret must be at least 32 characters, got 5

# Fix: Use minimum 32 characters
./bin/client_secret_hasher --secret "$(openssl rand -base64 48 | cut -c1-32)"
```

### Binary Not Found

```bash
# Build the binary first
make client-secret-hasher

# Or build all utilities
make build-utils

# Check binary exists
ls -lh ./bin/client_secret_hasher
```

### Permission Denied

```bash
# Make binary executable
chmod +x ./bin/client_secret_hasher

# Run with explicit path
./bin/client_secret_hasher --secret "..."
```

## Performance Benchmarks

### Default Parameters (t=2, m=64MB)

```bash
# Test 100 hashes
time for i in {1..100}; do
    ./bin/client_secret_hasher --secret "test-secret-32chars-minimum-length" > /dev/null
done

# Typical results:
# real    0m5.200s (average ~52ms per hash)
# user    0m3.800s
# sys     0m1.200s
```

### Higher Security (t=4, m=256MB)

```bash
# Single hash with high security
time ./bin/client_secret_hasher \
    --secret "test-secret-32chars-minimum-length" \
    --iterations 4 \
    --memory 262144 > /dev/null

# Typical results:
# real    0m0.420s (~420ms)
```

## Security Recommendations

1. **Never commit secrets**: Use environment variables or secret managers
2. **Rotate regularly**: Re-hash secrets periodically (quarterly/yearly)
3. **Audit access**: Log when secrets are revealed or re-hashed
4. **Secure transmission**: Use TLS/HTTPS when transmitting hashes
5. **Backup safely**: Encrypt backup files containing secrets
6. **Production parameters**: Use default or higher for production systems

## See Also

- [README.md](./README.md) - Full documentation
- [Password Package](../../internal/shared/password/) - Go package reference
- [OAuth Client Domain](../../internal/domain/oauth_client/) - Integration example
