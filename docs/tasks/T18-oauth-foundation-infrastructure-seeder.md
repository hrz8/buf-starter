# Task T18: OAuth Server Foundation Infrastructure & Seeder

## Story Reference
**User Story**: US5-oauth-server-foundation.md
**Epic**: OAuth Server Infrastructure
**Priority**: P0 (Foundational)

## Task Overview

**Type**: Backend Foundation
**Priority**: High
**Estimated Effort**: 4-5 hours
**Complexity**: Medium

**Status**: ✅ Migration Already Complete
The database migration file `20260108000000_create_oauth_server_tables.sql` was accidentally created earlier but is **actually correct and complete**. This task focuses on the remaining infrastructure components needed to complete US5.

## Objective

Complete the OAuth server foundation by implementing the infrastructure components that support the database schema: partition registration, error handling, RSA key management, configuration extensions, and a custom seeder that reads from config.yaml to bootstrap the system with superadmin user, default OAuth client, and OAuth providers.

## Prerequisites

- ✅ Database migration `20260108000000_create_oauth_server_tables.sql` already exists and is correct
- ✅ Existing partitioned tables pattern in `internal/domain/project/repo.go`
- ✅ Existing error handling patterns in domain error files
- ✅ Existing AES encryption utility for secrets
- ✅ Existing nanoid utility for public IDs
- ✅ bcrypt library available for password hashing

## Acceptance Criteria

### 1. Partition Registration
- [ ] Add "altalune_oauth_clients" to partitionedTables list
- [ ] Add "altalune_oauth_client_scopes" to partitionedTables list
- [ ] New projects automatically create partitions for these tables
- [ ] Partition naming follows convention: `{table_name}_p{project_id}`

### 2. OAuth Error Codes
- [ ] Create `internal/domain/oauth_client/errors.go` with client-related errors
- [ ] Create `internal/domain/oauth_auth/errors.go` with auth flow errors
- [ ] Create `internal/domain/project_member/errors.go` with membership errors
- [ ] All errors follow established pattern (var Err... = errors.New("..."))
- [ ] Error messages are clear and actionable

### 3. RSA Key Utilities
- [ ] Create `internal/shared/jwt/keygen.go` package
- [ ] Implement RSA key pair generation (2048 or 4096 bit)
- [ ] Implement PEM file saving (private key, public key)
- [ ] Implement PEM file loading with validation
- [ ] Handle file permissions properly (0600 for private key)
- [ ] Comprehensive error handling for crypto operations

### 4. Config.yaml Extensions
- [ ] Add `auth` section (host, port, sessionSecret, code/token expiry)
- [ ] Add `security.jwtPrivateKeyPath` configuration
- [ ] Add `security.jwtPublicKeyPath` configuration
- [ ] Add `security.jwksKid` (key ID for JWKS)
- [ ] Add `seeder` section (superadmin, defaultOAuthClient, oauthProviders)
- [ ] All new config fields properly documented with comments

### 5. Custom Seeder Implementation
- [ ] Create `internal/domain/oauth_seeder/` package
- [ ] Read configuration from config.yaml
- [ ] Generate nanoid public IDs for records
- [ ] Hash default client secret with bcrypt (cost 12+)
- [ ] Encrypt OAuth provider secrets with existing AES utility
- [ ] Implement idempotent record creation (check existence before insert)
- [ ] Seed superadmin user (check by email)
- [ ] Seed default dashboard OAuth client with `pkce_required=true`
- [ ] Seed OAuth providers (Google, GitHub from config)
- [ ] Create user_identity for superadmin
- [ ] Create project_members entry (superadmin as owner)
- [ ] Comprehensive logging for seeding operations
- [ ] Proper error handling and rollback on failure

### 6. Migration Integration
- [ ] Hook seeder into existing migrate command workflow
- [ ] Seeder runs automatically after goose up migrations
- [ ] Seeder is optional via command flag (e.g., `--skip-seed`)
- [ ] Clear console output showing seeding progress

### 7. Testing & Validation
- [ ] Migration up/down runs successfully
- [ ] Seeder creates all expected records
- [ ] Seeder is idempotent (safe to run multiple times)
- [ ] Partitions auto-created for new projects
- [ ] Default client has `pkce_required=true`
- [ ] All constraints and indexes work correctly
- [ ] RSA key generation and loading tested
- [ ] Error codes compile and are usable

## Technical Requirements

### 1. Partition Registration

**File**: `internal/domain/project/repo.go` (Line ~580)

**Change**:
```go
var partitionedTables = []string{
	"altalune_example_employees",
	"altalune_project_api_keys",
	"altalune_oauth_clients",        // Add this
	"altalune_oauth_client_scopes",  // Add this
}
```

**Why**: Ensures partitions are automatically created when new projects are created.

---

### 2. OAuth Error Codes

**File**: `internal/domain/oauth_client/errors.go` (NEW)

```go
package oauth_client

import "errors"

var (
	ErrOAuthClientNotFound          = errors.New("oauth client not found")
	ErrOAuthClientAlreadyExists     = errors.New("oauth client already exists")
	ErrInvalidRedirectURI           = errors.New("invalid redirect URI")
	ErrClientSecretMismatch         = errors.New("client secret does not match")
	ErrDefaultClientCannotBeDeleted = errors.New("default dashboard client cannot be deleted")
	ErrClientBelongsToOtherProject  = errors.New("oauth client belongs to another project")
)
```

**File**: `internal/domain/oauth_auth/errors.go` (NEW)

```go
package oauth_auth

import "errors"

var (
	ErrInvalidAuthorizationCode     = errors.New("invalid authorization code")
	ErrAuthorizationCodeExpired     = errors.New("authorization code has expired")
	ErrAuthorizationCodeAlreadyUsed = errors.New("authorization code already used")
	ErrInvalidRefreshToken          = errors.New("invalid refresh token")
	ErrRefreshTokenExpired          = errors.New("refresh token has expired")
	ErrRefreshTokenAlreadyUsed      = errors.New("refresh token already used")
	ErrInvalidPKCEVerifier          = errors.New("invalid PKCE code verifier")
	ErrPKCERequired                 = errors.New("PKCE is required for this client")
	ErrInvalidScope                 = errors.New("invalid scope requested")
	ErrUserConsentRequired          = errors.New("user consent required")
	ErrInvalidClientCredentials     = errors.New("invalid client credentials")
)
```

**File**: `internal/domain/project_member/errors.go` (NEW)

```go
package project_member

import "errors"

var (
	ErrProjectMemberNotFound      = errors.New("project member not found")
	ErrProjectMemberAlreadyExists = errors.New("user is already a member of this project")
	ErrInvalidRole                = errors.New("invalid role specified")
	ErrCannotModifyOwnerRole      = errors.New("owner role is reserved for superadmin")
	ErrCannotRemoveLastOwner      = errors.New("cannot remove the last owner from project")
	ErrInsufficientPermissions    = errors.New("insufficient permissions for this operation")
)
```

---

### 3. RSA Key Utilities

**File**: `internal/shared/jwt/keygen.go` (NEW)

**Functions**:
- `GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, error)` - Generate RSA key pair
- `SavePrivateKeyPEM(key *rsa.PrivateKey, filepath string) error` - Save private key
- `SavePublicKeyPEM(key *rsa.PublicKey, filepath string) error` - Save public key
- `LoadPrivateKeyPEM(filepath string) (*rsa.PrivateKey, error)` - Load private key
- `LoadPublicKeyPEM(filepath string) (*rsa.PublicKey, error)` - Load public key

**Implementation Requirements**:
- Use `crypto/rsa`, `crypto/x509`, `encoding/pem` packages
- Default key size: 2048 bits (configurable to 4096)
- Private key file permissions: 0600
- Public key file permissions: 0644
- PKCS#1 format for PEM encoding
- Comprehensive validation on loading

**Example Usage**:
```go
// Generate new key pair
privateKey, err := GenerateRSAKeyPair(2048)
if err != nil {
    return err
}

// Save keys
err = SavePrivateKeyPEM(privateKey, "keys/jwt-private.pem")
err = SavePublicKeyPEM(&privateKey.PublicKey, "keys/jwt-public.pem")

// Load keys later
privateKey, err := LoadPrivateKeyPEM("keys/jwt-private.pem")
publicKey, err := LoadPublicKeyPEM("keys/jwt-public.pem")
```

---

### 4. Config.yaml Extensions

**Add these sections** to `config.yaml`:

```yaml
# Existing sections...
server:
  host: localhost
  port: 3100
  # ... existing server config

database:
  # ... existing database config

# NEW: Authentication server configuration
auth:
  host: localhost
  port: 3101
  sessionSecret: "change-me-in-production-min-32-chars"
  codeExpiry: 600              # Authorization code expiry (seconds, 10 minutes)
  accessTokenExpiry: 3600      # Access token expiry (seconds, 1 hour)
  refreshTokenExpiry: 2592000  # Refresh token expiry (seconds, 30 days)

# EXTEND: Security configuration
security:
  # Existing
  iamEncryptionKey: "your-base64-encoded-32-byte-key-here"

  # NEW: JWT signing keys
  jwtPrivateKeyPath: "keys/jwt-private.pem"
  jwtPublicKeyPath: "keys/jwt-public.pem"
  jwksKid: "altalune-oauth-2026"  # Key ID for JWKS endpoint

# NEW: Database seeder configuration
seeder:
  # Superadmin user (created on first migration)
  superadmin:
    email: "admin@altalune.com"
    firstName: "Super"
    lastName: "Admin"

  # Default OAuth client (Altalune Dashboard)
  defaultOAuthClient:
    name: "Altalune Dashboard"
    clientId: "e730207a-0fce-495d-bac3-6211963ac423"  # Fixed UUID
    clientSecret: "change-me-dashboard-secret-min-32-chars"
    redirectUris:
      - "http://localhost:3000/auth/callback"
      - "https://dashboard.altalune.com/auth/callback"
    pkceRequired: true  # CRITICAL: Must be true for public client (dashboard)

  # OAuth providers for login
  oauthProviders:
    - provider: "google"
      clientId: "your-google-client-id.apps.googleusercontent.com"
      clientSecret: "your-google-client-secret"
      redirectUrl: "http://localhost:3101/auth/callback"
      scopes: "openid,profile,email"
      enabled: true

    - provider: "github"
      clientId: "your-github-client-id"
      clientSecret: "your-github-client-secret"
      redirectUrl: "http://localhost:3101/auth/callback"
      scopes: "read:user,user:email"
      enabled: true
```

**Notes**:
- sessionSecret must be at least 32 characters (for cookie encryption)
- clientSecret must be strong (will be hashed with bcrypt)
- RSA key paths are relative to project root or absolute
- Provider secrets will be encrypted with AES before storage

---

### 5. Custom Seeder Implementation

**Package**: `internal/domain/oauth_seeder/`

**Files**:
- `seeder.go` - Main seeder logic
- `config.go` - Config reading utilities
- `crypto.go` - Hashing and encryption utilities

**Seeder Interface**:
```go
type Seeder struct {
    db     *sql.DB
    config *Config  // Parsed from config.yaml
    logger *slog.Logger
}

func NewSeeder(db *sql.DB, configPath string) (*Seeder, error)
func (s *Seeder) Seed(ctx context.Context) error
```

**Seeding Steps** (in order):

1. **Parse config.yaml** - Read seeder section
2. **Check superadmin exists** - SELECT by email
3. **Create superadmin user** (if not exists)
   - Email from config
   - First/Last name from config
   - Generate nanoid for public_id
   - Created via INSERT INTO altalune_users
4. **Create user_identity for superadmin** (if not exists)
   - Link to superadmin user_id
   - Provider: "system" or "N/A"
   - oauth_client_id: NULL (superadmin doesn't use OAuth)
5. **Check default project exists** (assume project_id=1 or query first project)
6. **Create project_members entry** (if not exists)
   - Link superadmin to project
   - Role: "owner"
   - Generate nanoid for public_id
7. **Check default OAuth client exists** - SELECT by client_id UUID
8. **Create default OAuth client** (if not exists)
   - Name: "Altalune Dashboard"
   - client_id: Fixed UUID from config
   - client_secret_hash: bcrypt hash of plaintext secret from config
   - redirect_uris: Array from config
   - **pkce_required: TRUE** (critical requirement)
   - is_default: TRUE
   - project_id: Default project ID
   - Generate nanoid for public_id
9. **Seed OAuth providers** (Google, GitHub)
   - Check if exists by provider name
   - Encrypt client_secret with AES
   - Store in altalune_user_identities table or separate oauth_providers table (TBD based on existing schema)

**Idempotency Strategy**:
- Use SELECT before INSERT for all records
- Use UPSERT (ON CONFLICT DO UPDATE) where appropriate
- Log skipped records: "Superadmin already exists, skipping..."
- Safe to run seeder multiple times without duplicating data

**Error Handling**:
- Transaction wrapper for all operations
- Rollback on any error
- Detailed logging for each step
- Return aggregated error if seeding fails

**Security Requirements**:
- NEVER log plaintext secrets
- Hash client_secret with bcrypt.GenerateFromPassword (cost 12)
- Encrypt provider secrets with existing AES utility
- Validate config values before using (e.g., email format, UUID format)

---

### 6. Migration Integration

**File**: `cmd/altalune/migrate.go` (MODIFY)

**Changes**:
- After goose up completes successfully, call seeder
- Add `--skip-seed` flag to bypass seeding
- Log seeding progress clearly

**Example Integration**:
```go
// After goose up migrations
if !skipSeed {
    logger.Info("Running database seeder...")
    seeder, err := oauth_seeder.NewSeeder(db, configPath)
    if err != nil {
        return fmt.Errorf("seeder initialization failed: %w", err)
    }

    if err := seeder.Seed(ctx); err != nil {
        return fmt.Errorf("seeding failed: %w", err)
    }

    logger.Info("Database seeding completed successfully")
}
```

**Output Example**:
```
INFO: Running goose up migrations...
INFO: Migration 20260108000000_create_oauth_server_tables.sql applied
INFO: Running database seeder...
INFO: Checking superadmin user...
INFO: Superadmin user already exists, skipping
INFO: Checking default OAuth client...
INFO: Creating default OAuth client 'Altalune Dashboard'
INFO: Seeding OAuth provider: google
INFO: Seeding OAuth provider: github
INFO: Database seeding completed successfully
```

---

## Implementation Details

### Partition Registration (5 minutes)

1. Open `internal/domain/project/repo.go`
2. Locate `partitionedTables` variable (line ~580)
3. Add two new entries to the slice
4. No additional code changes needed (partition creation is automatic)

### Error Codes (10 minutes)

1. Create three new error files in respective domain packages
2. Follow existing error.go pattern from api_key domain
3. Use descriptive error messages
4. Ensure errors are exported (start with Err prefix)

### RSA Key Utilities (30 minutes)

1. Create new package `internal/shared/jwt/`
2. Implement 5 functions (generate, save private, save public, load private, load public)
3. Use standard crypto libraries
4. Add comprehensive error handling
5. Write unit tests for key round-trip (generate → save → load)

### Config Extensions (15 minutes)

1. Open `config.yaml`
2. Add `auth` section with 6 fields
3. Extend `security` section with 3 JWT fields
4. Add `seeder` section with nested structure
5. Add comments explaining each field
6. Validate YAML syntax

### Custom Seeder (2-3 hours)

**Phase 1: Setup** (30 min)
- Create package structure
- Define Seeder struct
- Implement config parsing
- Set up database connection

**Phase 2: Core Logic** (1.5 hours)
- Implement each seeding step as separate method
- Superadmin seeding
- Default client seeding
- OAuth provider seeding
- User identity and project member linkage

**Phase 3: Error Handling & Logging** (30 min)
- Transaction wrapper
- Comprehensive logging
- Idempotency checks
- Security validation

**Phase 4: Integration** (30 min)
- Hook into migrate command
- Add CLI flag
- Test end-to-end

### Testing (30 minutes)

1. **Migration Test**:
   ```bash
   ./bin/app migrate -c config.yaml
   # Check: All tables created, seeder runs successfully
   ```

2. **Idempotency Test**:
   ```bash
   ./bin/app migrate -c config.yaml  # Run again
   # Check: No duplicate records, "already exists" logs appear
   ```

3. **Partition Test**:
   ```sql
   -- Create new project
   INSERT INTO altalune_projects (public_id, name) VALUES ('test123', 'Test Project') RETURNING id;

   -- Check partitions created
   SELECT tablename FROM pg_tables WHERE tablename LIKE 'altalune_oauth_clients_p%';
   SELECT tablename FROM pg_tables WHERE tablename LIKE 'altalune_oauth_client_scopes_p%';
   ```

4. **Seeder Validation**:
   ```sql
   -- Check superadmin
   SELECT * FROM altalune_users WHERE email = 'admin@altalune.com';

   -- Check default client (pkce_required must be TRUE)
   SELECT client_id, name, pkce_required, is_default
   FROM altalune_oauth_clients
   WHERE is_default = true;

   -- Check project membership
   SELECT u.email, pm.role
   FROM project_members pm
   JOIN altalune_users u ON pm.user_id = u.id
   WHERE pm.role = 'owner';
   ```

5. **RSA Key Test**:
   ```go
   // Unit test
   func TestRSAKeyGeneration(t *testing.T) {
       key, err := GenerateRSAKeyPair(2048)
       require.NoError(t, err)

       err = SavePrivateKeyPEM(key, "/tmp/test-private.pem")
       require.NoError(t, err)

       loadedKey, err := LoadPrivateKeyPEM("/tmp/test-private.pem")
       require.NoError(t, err)
       require.Equal(t, key.D, loadedKey.D)
   }
   ```

---

## Files to Create

### New Files
1. `internal/domain/oauth_client/errors.go` - OAuth client error definitions
2. `internal/domain/oauth_auth/errors.go` - OAuth auth flow error definitions
3. `internal/domain/project_member/errors.go` - Project membership error definitions
4. `internal/shared/jwt/keygen.go` - RSA key generation and loading utilities
5. `internal/domain/oauth_seeder/seeder.go` - Main seeder logic
6. `internal/domain/oauth_seeder/config.go` - Config parsing utilities
7. `internal/domain/oauth_seeder/crypto.go` - Hashing and encryption utilities
8. `keys/.gitkeep` - Placeholder for keys directory (keys themselves in .gitignore)

### Files to Modify
1. `internal/domain/project/repo.go` - Add partition table names (line ~580)
2. `cmd/altalune/migrate.go` - Integrate seeder call after migrations
3. `config.yaml` - Add auth, security.jwt*, seeder sections
4. `.gitignore` - Add `keys/*.pem` to ignore RSA private keys

---

## Testing Requirements

### Unit Tests
- [ ] RSA key generation creates valid keys
- [ ] RSA key save/load round-trip works
- [ ] Config parsing reads all seeder fields correctly
- [ ] bcrypt hashing produces valid hashes
- [ ] Errors are properly defined and usable

### Integration Tests
- [ ] Migration creates all tables successfully
- [ ] Seeder creates superadmin user
- [ ] Seeder creates default OAuth client with pkce_required=true
- [ ] Seeder creates OAuth providers
- [ ] Seeder creates project membership with owner role
- [ ] Seeder is idempotent (no duplicates on re-run)
- [ ] Partitions auto-created for new projects
- [ ] All indexes and constraints work

### Manual Validation
- [ ] Run `./bin/app migrate -c config.yaml`
- [ ] Check database for seeded records
- [ ] Verify default client has `pkce_required=true`
- [ ] Verify superadmin has owner role
- [ ] Create new project and verify partitions exist
- [ ] Run seeder twice and verify no duplicate records

---

## Validation Checklist

### Configuration
- [ ] config.yaml has `auth` section with all 6 fields
- [ ] config.yaml has `security.jwt*` fields (3 new fields)
- [ ] config.yaml has `seeder` section with superadmin, client, providers
- [ ] sessionSecret is at least 32 characters
- [ ] clientSecret is strong and unique
- [ ] redirectUris include localhost and production URLs
- [ ] pkceRequired is explicitly set to `true`

### Database
- [ ] Migration file exists and is correct (already validated)
- [ ] All 7 tables created: oauth_scopes, oauth_clients, oauth_client_scopes, oauth_authorization_codes, oauth_refresh_tokens, oauth_user_consents, project_members
- [ ] user_identities has oauth_client_id and last_login_at columns
- [ ] Partitioned tables registered in project repo
- [ ] Partitions auto-created for oauth_clients and oauth_client_scopes

### Seeder
- [ ] Superadmin user created with email from config
- [ ] Default OAuth client created with fixed UUID
- [ ] Client secret hashed with bcrypt (cost 12+)
- [ ] Default client has `pkce_required=true` ⚠️ CRITICAL
- [ ] Default client has `is_default=true`
- [ ] OAuth providers seeded (Google, GitHub)
- [ ] Provider secrets encrypted with AES
- [ ] User identity created for superadmin
- [ ] Project membership created (role=owner)
- [ ] Seeder is idempotent (check before insert)

### Code Quality
- [ ] Error codes follow established pattern
- [ ] RSA utilities have proper error handling
- [ ] Seeder has comprehensive logging
- [ ] No plaintext secrets in logs
- [ ] File permissions correct (0600 for private key)
- [ ] All exported functions have godoc comments
- [ ] Code passes `make format`

---

## Definition of Done

- [ ] All 7 components implemented (partitions, errors, RSA, config, seeder, integration, tests)
- [ ] Partition registration complete (2 tables added to list)
- [ ] Error codes defined in 3 new domain packages
- [ ] RSA key utilities implemented and tested
- [ ] config.yaml extended with all required sections
- [ ] Custom seeder implemented with idempotency
- [ ] Seeder integrated into migrate command
- [ ] Migration runs successfully (up and down)
- [ ] Seeder creates all expected records
- [ ] Default OAuth client has `pkce_required=true` ⚠️ CRITICAL
- [ ] Partitions auto-created for new projects
- [ ] All tests passing (unit + integration)
- [ ] Manual validation completed
- [ ] Code reviewed and approved
- [ ] No security issues (secrets properly handled)

---

## Dependencies

### Existing Code
- Migration file: `database/migrations/20260108000000_create_oauth_server_tables.sql` ✅
- Partition pattern: `internal/domain/project/repo.go` partitionedTables
- Error pattern: `internal/domain/api_key/errors.go`
- Nanoid utility: Existing in codebase
- AES encryption: `internal/shared/iam/encryption.go` (assumed)

### External Libraries
- `golang.org/x/crypto/bcrypt` - Client secret hashing
- `crypto/rsa`, `crypto/x509`, `encoding/pem` - Standard library (RSA keys)
- `database/sql` - Database operations
- `gopkg.in/yaml.v3` - Config parsing (assumed existing)

### System Requirements
- PostgreSQL 14+ (partitioned table support)
- Go 1.21+ (existing requirement)
- File system write access for keys directory

---

## Risk Factors

### Security Risks (High Priority)

1. **Client Secret Exposure**
   - Risk: Plaintext secret in config.yaml could be committed to git
   - Mitigation: Add warning comment in config.yaml, document in README
   - Severity: High

2. **Private Key Security**
   - Risk: RSA private key file permissions too permissive
   - Mitigation: Enforce 0600 permissions in SavePrivateKeyPEM
   - Severity: High

3. **Seeder Secret Logging**
   - Risk: Accidentally logging plaintext secrets during seeding
   - Mitigation: Never log raw secrets, only hash/encrypt first
   - Severity: High

### Implementation Risks (Medium Priority)

1. **Config Parsing Errors**
   - Risk: Malformed YAML breaks seeder
   - Mitigation: Validate config structure before using values
   - Severity: Medium

2. **Idempotency Failures**
   - Risk: Duplicate records created on re-run
   - Mitigation: Check existence before INSERT, use transactions
   - Severity: Medium

3. **Partition Creation Timing**
   - Risk: Seeder runs before partitions exist
   - Mitigation: Migration creates base tables first, partitions added dynamically
   - Severity: Low (migrations run before seeder)

---

## Notes

### Critical Implementation Details

1. **PKCE Requirement**:
   - The default dashboard client MUST have `pkce_required=true`
   - This is non-negotiable for public client security (dashboard is a SPA)
   - Verify in seeder and tests

2. **Fixed Client ID**:
   - Dashboard client uses fixed UUID: `e730207a-0fce-495d-bac3-6211963ac423`
   - This ensures consistent reference across environments
   - Seeder should use this exact UUID from config

3. **Owner Role Restriction**:
   - Only superadmin can have "owner" role
   - Seeder creates this mapping automatically
   - Future user registrations get "user" role by default

4. **Partition Registration Order**:
   - Must register partitioned tables BEFORE creating projects
   - Our case: Migration creates tables, registration happens in code
   - New projects will auto-create partitions

5. **Seeder Transaction Scope**:
   - Entire seeding process should be in ONE transaction
   - Rollback if ANY step fails
   - Prevents partial seeding state

### Config.yaml Security Notes

Add these comments to config.yaml:

```yaml
seeder:
  # ⚠️ SECURITY WARNING:
  # - Never commit real client secrets to version control
  # - Use strong secrets (min 32 characters, random)
  # - Rotate secrets regularly in production
  # - Consider using environment variable substitution

  defaultOAuthClient:
    clientSecret: "change-me-dashboard-secret-min-32-chars"  # ⚠️ CHANGE IN PRODUCTION
```

### Seeder Execution Context

The seeder should:
- Run automatically after `goose up` migrations
- Be skippable via `--skip-seed` flag
- Log clearly what it's creating vs skipping
- Exit with error code if seeding fails
- Not affect `goose down` rollback (seeding is data, not schema)

### Testing in Development

Recommended test workflow:
```bash
# Clean slate
dropdb altalune_dev
createdb altalune_dev

# Run migration + seeder
./bin/app migrate -c config.yaml

# Verify seeding
psql altalune_dev -c "SELECT email, first_name FROM altalune_users;"
psql altalune_dev -c "SELECT name, pkce_required, is_default FROM altalune_oauth_clients;"

# Test idempotency
./bin/app migrate -c config.yaml  # Run again, should skip existing records
```

### Future Enhancements

- [ ] Environment variable substitution in config.yaml (e.g., `${DASHBOARD_CLIENT_SECRET}`)
- [ ] Seeder dry-run mode (`--seed-dry-run`)
- [ ] Seeder force mode (`--seed-force` to recreate records)
- [ ] Multiple default clients per environment (dev, staging, prod)
- [ ] Automated RSA key rotation workflow
- [ ] Seeder validation step (verify all required records exist)

---

## Related Documentation

- **User Story**: `docs/stories/US5-oauth-server-foundation.md`
- **Planning Doc**: `oauth_server_prepare/plan.md`
- **Migration File**: `database/migrations/20260108000000_create_oauth_server_tables.sql`
- **Partition Pattern**: `internal/domain/project/repo.go` (line 580)
- **Error Pattern**: `internal/domain/api_key/errors.go`

---

**Task Created**: 2026-01-08
**Estimated Completion**: 4-5 hours
**Assigned To**: Backend Team
**Status**: Ready for Implementation
