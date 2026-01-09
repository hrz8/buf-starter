# User Story US5: OAuth Server Foundation & Database Schema

## Story Overview

**As a** system architect
**I want** to establish the foundational database schema and infrastructure for the OAuth 2.0 authorization server
**So that** Altalune can function as an identity provider for both the dashboard and external applications

## Acceptance Criteria

### Core Functionality

#### OAuth Server Database Tables

- **Given** the need for OAuth server functionality
- **When** the migration is executed
- **Then** the following tables should be created:
  - `altalune_oauth_scopes` - Standard and custom OAuth scopes
  - `altalune_oauth_clients` - OAuth client applications (partitioned by project_id)
  - `altalune_oauth_client_scopes` - Client-scope mappings (partitioned by project_id)
  - `altalune_oauth_authorization_codes` - Short-lived authorization codes (global)
  - `altalune_oauth_refresh_tokens` - Long-lived refresh tokens (global)
  - `altalune_oauth_user_consents` - User consent tracking (global)
  - `altalune_project_members` - Project membership with roles (global)
- **And** all tables should have appropriate indexes for performance
- **And** partitioned tables should support automatic partition creation

#### User Identity Extensions

- **Given** the existing `altalune_user_identities` table
- **When** the migration is executed
- **Then** the table should have new columns:
  - `oauth_client_id` (UUID) - Links identity to OAuth client
  - `last_login_at` (TIMESTAMPTZ) - Tracks last login timestamp
- **And** appropriate indexes should be created on new columns

#### Standard OAuth Scopes Seeding

- **Given** the OAuth server needs standard scopes
- **When** the migration is executed
- **Then** the following standard scopes should be seeded:
  - `openid` - OpenID Connect authentication
  - `profile` - Access to user profile information
  - `email` - Access to user email address
  - `offline_access` - Request refresh token for offline access
- **And** these scopes should be marked as `is_standard = true`
- **And** each scope should have a unique public_id

#### Custom Migration Seeder

- **Given** the need to seed environment-specific OAuth data
- **When** the `migrate` command runs with custom seeder
- **Then** it should read configuration from `config.yaml`
- **And** it should seed:
  - Superadmin user (email from config)
  - Default dashboard OAuth client (credentials from config)
  - OAuth providers (Google, GitHub from config)
  - User identity for superadmin linked to dashboard client
  - Project membership for superadmin as owner
- **And** the seeder should be idempotent (safe to run multiple times)
- **And** it should check for existing records before inserting

#### Automatic Partition Creation

- **Given** `altalune_oauth_clients` and `altalune_oauth_client_scopes` are partitioned tables
- **When** a new project is created
- **Then** partitions should be automatically created for both tables
- **And** partition naming should follow: `{table_name}_p{project_id}`
- **And** these tables should be added to the `partitionedTables` list in `project` repo

### Security Requirements

#### Client Secret Encryption

- Client secrets for OAuth clients must be hashed with bcrypt (cost 12+)
- Client secrets should never be stored in plaintext
- Only the hash should be stored in the database
- Plaintext secret only shown once during client creation

#### RSA Key Pair for JWT Signing

- System must generate RSA-2048 or RSA-4096 key pair for JWT signing
- Private key must be securely stored (file system with restricted permissions)
- Public key path must be configured in `config.yaml`
- Keys should be in PEM format
- Utility should support key generation and loading

#### Soft-Delete Pattern

- Authorization codes must use soft-delete (exchange_at timestamp)
- Refresh tokens must use soft-delete (exchange_at timestamp)
- Deleted records preserved for audit trail
- Indexes optimized for active records only (WHERE exchange_at IS NULL)

### Data Validation

#### OAuth Clients Table

- `project_id` - Required, references altalune_projects
- `public_id` - Required, unique nanoid (14 chars)
- `name` - Required, 1-100 characters
- `client_id` - Required, unique UUID v4
- `client_secret_hash` - Required, bcrypt hash
- `redirect_uris` - Required, array with at least one URI
- `pkce_required` - Required boolean, default false
- `is_default` - Required boolean, default false (true for dashboard client)

#### OAuth Scopes Table

- `public_id` - Required, unique nanoid (14 chars)
- `name` - Required, unique, 1-50 characters
- `description` - Optional text
- `is_standard` - Required boolean, default false

#### Authorization Codes Table

- `code` - Required, unique UUID v4
- `client_id` - Required UUID
- `user_id` - Required, references altalune_users
- `redirect_uri` - Required, 1-500 characters
- `scope` - Optional text
- `nonce` - Optional, 1-100 characters
- `code_challenge` - Optional, 1-128 characters (for PKCE)
- `code_challenge_method` - Optional, enum ('S256', 'plain')
- `expires_at` - Required timestamp
- `exchange_at` - Optional timestamp (soft delete marker)

#### Refresh Tokens Table

- `token` - Required, unique UUID v4
- `client_id` - Required UUID
- `user_id` - Required, references altalune_users
- `scope` - Required text
- `nonce` - Optional, 1-100 characters
- `expires_at` - Required timestamp
- `exchange_at` - Optional timestamp (soft delete marker)

#### Project Members Table

- `public_id` - Required, unique nanoid (14 chars)
- `project_id` - Required, references altalune_projects
- `user_id` - Required, references altalune_users
- `role` - Required, enum ('owner', 'admin', 'member', 'user')
- Unique constraint on (project_id, user_id)

### Configuration Requirements

#### config.yaml OAuth Section

- **Given** the need for OAuth server configuration
- **When** config.yaml is updated
- **Then** it should include the following sections:

```yaml
auth:
  host: localhost
  port: 3101
  sessionSecret: "session-encryption-key"
  codeExpiry: 600                    # 10 minutes
  accessTokenExpiry: 3600            # 1 hour
  refreshTokenExpiry: 2592000        # 30 days

security:
  jwtPrivateKeyPath: "path/to/rsa-private.pem"
  jwtPublicKeyPath: "path/to/rsa-public.pem"
  jwksKid: "altalune-oauth-2024"

seeder:
  superadmin:
    email: "admin@altalune.com"
    firstName: "Super"
    lastName: "Admin"

  defaultOAuthClient:
    name: "Altalune Dashboard"
    clientId: "00000000-0000-0000-0000-000000000001"
    clientSecret: "your-secure-dashboard-client-secret"
    pkceRequired: true                # CRITICAL: Dashboard is public client
    redirectUris:
      - "http://localhost:3000/auth/callback"
      - "https://dashboard.altalune.com/auth/callback"

  oauthProviders:
    - provider: "google"
      clientId: "your-google-client-id"
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

### User Experience

#### Migration Execution

- Migration should execute cleanly without errors
- Migration should be reversible (goose down works correctly)
- Clear error messages if migration fails
- Progress indication during long operations

#### Seeder Execution

- Seeder should provide clear feedback on what's being created
- Should skip existing records with informative messages
- Should report success/failure for each seeded entity
- Should be safe to run multiple times (idempotent)

## Technical Requirements

### Backend Architecture

- **Migration Framework**: Use Goose for database migrations
- **Migration File**: `database/migrations/20260108000000_create_oauth_server_tables.sql`
- **Partitioned Tables**: Add `altalune_oauth_clients` and `altalune_oauth_client_scopes` to `partitionedTables` in `internal/domain/project/repo.go`
- **Seeder Implementation**: Create `internal/domain/oauth_seeder/` package
  - Read config from config.yaml
  - Generate nanoid public IDs
  - Hash client secret with bcrypt
  - Encrypt OAuth provider secrets with AES-256-GCM (existing utility)
  - Create records with proper error handling
- **RSA Key Generation**: Create utility in `internal/shared/jwt/keygen.go`
  - Generate RSA key pair
  - Save to PEM files
  - Load keys from files
  - Validate key format and size

### Database Considerations

- **Partitioning Strategy**:
  - OAuth clients partitioned by project_id (multi-tenant isolation)
  - Authorization codes and refresh tokens global (cross-project identity)
- **Indexes**:
  - Unique indexes on all ID fields (public_id, client_id, code, token)
  - Partial indexes on active records (WHERE exchange_at IS NULL)
  - Foreign key indexes for joins
  - Composite indexes for common query patterns
- **Constraints**:
  - Foreign keys with CASCADE DELETE where appropriate
  - CHECK constraints for enum validation
  - UNIQUE constraints for business rules
  - Array length validation for redirect_uris

### Error Codes

Add OAuth-specific error codes to `errors.go`:

```go
// OAuth Client errors (608xx)
CodeOAuthClientNotFound          = "60810"
CodeOAuthClientAlreadyExists     = "60811"
CodeOAuthClientInvalidRedirectURI = "60812"
CodeOAuthClientDeleteDefault     = "60813"  // Cannot delete default client

// OAuth Authorization Code errors (609xx)
CodeOAuthAuthCodeNotFound        = "60910"
CodeOAuthAuthCodeExpired         = "60911"
CodeOAuthAuthCodeAlreadyUsed     = "60912"
CodeOAuthAuthCodeInvalidPKCE     = "60913"

// OAuth Refresh Token errors (610xx)
CodeOAuthRefreshTokenNotFound    = "61010"
CodeOAuthRefreshTokenExpired     = "61011"
CodeOAuthRefreshTokenAlreadyUsed = "61012"

// OAuth Scope errors (611xx)
CodeOAuthScopeNotFound           = "61110"
CodeOAuthScopeInvalid            = "61111"

// Project Member errors (612xx)
CodeProjectMemberNotFound        = "61210"
CodeProjectMemberAlreadyExists   = "61211"
CodeProjectMemberInvalidRole     = "61212"
CodeProjectMemberCannotRemoveOwner = "61213"
```

## Out of Scope

- OAuth flow implementation (covered in US7)
- OAuth client CRUD UI (covered in US6)
- Token generation and validation logic (covered in US7)
- serve-auth command implementation (covered in US7)
- Dashboard OAuth integration (covered in US8)
- Frontend components for OAuth (covered in US6, US8)
- Project member management UI (covered in US8)
- User role transitions (covered in US8)

## Dependencies

- Existing project management functionality (`altalune_projects` table)
- Existing user management functionality (`altalune_users` table)
- Existing IAM tables (`altalune_user_identities`)
- Existing OAuth provider table (`altalune_oauth_providers`)
- Existing encryption utilities (AES-256-GCM for OAuth provider secrets)
- Existing nanoid generation utility
- Goose migration framework
- PostgreSQL 14+ with partition support
- bcrypt library for password hashing

## Definition of Done

- [ ] Migration file created with all OAuth server tables
- [ ] Migration includes user_identities extensions
- [ ] Standard OAuth scopes seeded in migration
- [ ] Migration tested with `goose up` and `goose down`
- [ ] Partitioned tables added to project repo's partitionedTables list
- [ ] Custom seeder implemented reading from config.yaml
- [ ] Seeder creates superadmin user
- [ ] Seeder creates default dashboard OAuth client (with pkce_required=true)
- [ ] Seeder creates OAuth providers from config
- [ ] Seeder creates user_identity for superadmin
- [ ] Seeder creates project_member for superadmin as owner
- [ ] Seeder is idempotent (can run multiple times safely)
- [ ] RSA key generation utility created
- [ ] RSA keys can be loaded from config paths
- [ ] OAuth-specific error codes added to errors.go
- [ ] config.yaml updated with all OAuth sections
- [ ] config.yaml includes example values with comments
- [ ] All database constraints properly defined
- [ ] All indexes created for performance
- [ ] Foreign key relationships correctly defined
- [ ] Migration documentation updated
- [ ] Code follows established patterns and guidelines
- [ ] Migration tested in development environment
- [ ] Partition creation tested for new projects

## Notes

### Critical Implementation Details

1. **Dashboard Client PKCE**: The default dashboard OAuth client MUST have `pkce_required = true` because the dashboard is a public client (SPA). This is a security requirement for OAuth 2.1 compliance.

2. **Soft-Delete Pattern**: Authorization codes and refresh tokens use `exchange_at` timestamp instead of hard delete. This preserves audit trail while maintaining query performance through partial indexes.

3. **Role Hierarchy**:
   - `owner`: Reserved for superadmin only, auto-assigned to all projects
   - `admin`: Can manage project and members (except owner role)
   - `member`: Can access project data, read-only settings
   - `user`: OAuth user with no dashboard access unless upgraded

4. **Partition Strategy**:
   - OAuth clients are project-specific (partitioned)
   - Authorization codes/refresh tokens are user-specific (global)
   - This allows cross-project user identity while maintaining project isolation

5. **Seeder Idempotency**: The seeder should check for existing records using:
   - Superadmin: Check by email
   - Default client: Check by client_id UUID
   - OAuth providers: Check by provider type
   - Use `ON CONFLICT DO NOTHING` or explicit checks

### Future Enhancements

- Custom scope management UI (future story)
- OAuth provider auto-registration on first use
- Key rotation mechanisms
- Token cleanup cron job (remove expired codes/tokens)
- Project member role expiration (temporary access)
- Audit logging for seeder operations

### Related Stories

- US3: IAM Core Entities and Mappings (provides user/identity tables)
- US4: OAuth Provider Configuration (provides OAuth provider login)
- US6: OAuth Client Management (uses these tables for CRUD)
- US7: OAuth Authorization Server (generates codes/tokens)
- US8: Dashboard OAuth Integration (uses project_members for access control)

### Security Considerations

- Client secrets hashed with bcrypt (cost 12+), never plaintext
- RSA private key must have restricted file permissions (600)
- OAuth provider secrets encrypted with AES-256-GCM
- Soft-delete preserves forensic evidence for security audits
- Indexes optimized to prevent timing attacks on token lookups
