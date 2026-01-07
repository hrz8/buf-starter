# Task T15: OAuth Backend Domain Layer (Proto + 7-File Pattern)

**Story Reference:** US4-oauth-provider-configuration.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 8-10 hours
**Prerequisites:** T13 (database schema), T14 (crypto utilities)

## Objective

Implement complete OAuth provider backend domain following the established 7-file pattern, including proto definitions with validation rules, repository with encryption integration, service layer with business logic, error handling, and container wiring.

## Acceptance Criteria

- [ ] Proto schema defined with 6 RPCs and validation rules
- [ ] `buf generate` runs successfully and generates Go/TypeScript code
- [ ] All 7 domain files implemented (model, interface, repo, service, handler, mapper, errors)
- [ ] Repository encrypts client secrets before INSERT
- [ ] Repository re-encrypts client secrets on UPDATE (if provided)
- [ ] Repository never returns actual client secret (only `client_secret_set` boolean)
- [ ] RevealClientSecret RPC decrypts and returns plaintext secret
- [ ] Service validates unique constraint on provider_type
- [ ] Service enforces provider_type immutability (cannot change in update)
- [ ] Error codes added in 608XX range
- [ ] Container wiring complete with encryption key injection
- [ ] Service registered in server
- [ ] All 6 RPCs tested and working

## Technical Requirements

### Proto Schema Design

**File:** `api/proto/altalune/v1/oauth_provider.proto`

**Key Design Decisions:**
1. **No Activate/Deactivate RPCs** - Use `enabled` boolean in Update instead
2. **client_secret_set flag** - OAuthProvider message never includes actual secret
3. **Separate RevealClientSecret RPC** - Security best practice, explicit action
4. **Optional client_secret in Update** - If empty, retain existing secret
5. **provider_type immutable** - Not included in UpdateRequest

**Service Methods (6 RPCs):**
```protobuf
service OAuthProviderService {
  rpc QueryOAuthProviders(QueryOAuthProvidersRequest) returns (QueryOAuthProvidersResponse) {}
  rpc CreateOAuthProvider(CreateOAuthProviderRequest) returns (CreateOAuthProviderResponse) {}
  rpc GetOAuthProvider(GetOAuthProviderRequest) returns (GetOAuthProviderResponse) {}
  rpc UpdateOAuthProvider(UpdateOAuthProviderRequest) returns (UpdateOAuthProviderResponse) {}
  rpc DeleteOAuthProvider(DeleteOAuthProviderRequest) returns (DeleteOAuthProviderResponse) {}
  rpc RevealClientSecret(RevealClientSecretRequest) returns (RevealClientSecretResponse) {}
}
```

**Validation Rules:**
- `provider_type`: Required, enum, defined_only
- `client_id`: Required, 1-500 chars
- `client_secret`: Required (create), optional (update), 1-500 chars
- `redirect_url`: Required, valid URI, max 500 chars
- `scopes`: Optional, max 1000 chars

### Domain Layer Architecture (7 Files)

**Location:** `internal/domain/oauth_provider/`

```
oauth_provider/
├── model.go       # Domain models, ProviderType enum, input/result types
├── interface.go   # Repository interface with 7 methods
├── mapper.go      # Proto ↔ domain conversions (ProviderType enum mapping)
├── errors.go      # Domain-specific errors (NotFound, DuplicateType, etc.)
├── repo.go        # PostgreSQL implementation with encryption/decryption
├── handler.go     # Connect-RPC HTTP handlers (thin wrappers)
└── service.go     # Business logic, validation, error handling
```

### Repository Encryption Integration

**Key Methods:**
```go
type Repository interface {
    Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[OAuthProvider], error)
    Create(ctx context.Context, input *CreateOAuthProviderInput) (*CreateOAuthProviderResult, error)
    GetByID(ctx context.Context, publicID string) (*OAuthProvider, error)
    GetByProviderType(ctx context.Context, providerType ProviderType) (*OAuthProvider, error)
    Update(ctx context.Context, input *UpdateOAuthProviderInput) (*UpdateOAuthProviderResult, error)
    Delete(ctx context.Context, input *DeleteOAuthProviderInput) error
    RevealClientSecret(ctx context.Context, publicID string) (string, error)
}
```

**Encryption Integration Points:**
1. **Create:** Encrypt client_secret before INSERT
2. **Update:** Re-encrypt client_secret if provided (optional)
3. **Query/Get:** Never select client_secret column, return client_secret_set=true
4. **RevealClientSecret:** SELECT client_secret, decrypt, return plaintext

### Service Layer Validation

**Business Rules:**
1. **Unique provider_type:** Check before create, return DuplicateProviderType error
2. **Provider type immutability:** Preserve existing provider_type on update
3. **Optional secret update:** If client_secret empty, don't update in database
4. **Trim inputs:** client_id, redirect_url, scopes

**Error Handling:**
- Encryption failures → EncryptionFailed error
- Decryption failures → DecryptionFailed error
- Not found → OAuthProviderNotFound error
- Duplicate → DuplicateProviderType error

### Error Codes (608XX Range)

**Location:** `errors.go`

```go
const (
    CodeOAuthProviderNotFound        = "60810"
    CodeOAuthProviderDuplicateType   = "60811"
    CodeOAuthProviderEncryptionError = "60812"
    CodeOAuthProviderDecryptionError = "60813"
)
```

## Implementation Details

### 1. Proto Schema (`oauth_provider.proto`)

See detailed proto schema in US4 technical requirements. Key sections:

**Enum:**
```protobuf
enum ProviderType {
  PROVIDER_TYPE_UNSPECIFIED = 0;
  PROVIDER_TYPE_GOOGLE = 1;
  PROVIDER_TYPE_GITHUB = 2;
  PROVIDER_TYPE_MICROSOFT = 3;
  PROVIDER_TYPE_APPLE = 4;
}
```

**Message (Never includes actual secret):**
```protobuf
message OAuthProvider {
  string id = 1;
  ProviderType provider_type = 2;
  string client_id = 3;
  bool client_secret_set = 4;  // True if secret exists, NEVER actual secret
  string redirect_url = 5;
  string scopes = 6;
  bool enabled = 7;
  google.protobuf.Timestamp created_at = 98;
  google.protobuf.Timestamp updated_at = 99;
}
```

### 2. Model (`model.go`)

**Domain Types:**
```go
type ProviderType string

const (
    ProviderTypeGoogle    ProviderType = "google"
    ProviderTypeGithub    ProviderType = "github"
    ProviderTypeMicrosoft ProviderType = "microsoft"
    ProviderTypeApple     ProviderType = "apple"
)

type OAuthProvider struct {
    ID              string
    ProviderType    ProviderType
    ClientID        string
    ClientSecretSet bool  // Never expose actual secret
    RedirectURL     string
    Scopes          string
    Enabled         bool
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type CreateOAuthProviderInput struct {
    ProviderType ProviderType
    ClientID     string
    ClientSecret string  // Plaintext (encrypted in repo)
    RedirectURL  string
    Scopes       string
    Enabled      bool
}

type UpdateOAuthProviderInput struct {
    PublicID     string
    ClientID     string
    ClientSecret string  // Optional, if empty retain existing
    RedirectURL  string
    Scopes       string
    Enabled      bool
}
```

### 3. Repository (`repo.go`) - Critical Encryption Integration

**Constructor with encryption key:**
```go
type Repo struct {
    db            postgres.DB
    encryptionKey []byte  // Injected from config
}

func NewRepo(db postgres.DB, encryptionKey []byte) *Repo {
    return &Repo{
        db:            db,
        encryptionKey: encryptionKey,
    }
}
```

**Create with encryption:**
```go
func (r *Repo) Create(ctx context.Context, input *CreateOAuthProviderInput) (*CreateOAuthProviderResult, error) {
    publicID, _ := nanoid.GeneratePublicID()

    // Check duplicate provider_type
    existing, err := r.GetByProviderType(ctx, input.ProviderType)
    if err != nil && err != ErrOAuthProviderNotFound {
        return nil, fmt.Errorf("check duplicate: %w", err)
    }
    if existing != nil {
        return nil, ErrDuplicateProviderType
    }

    // CRITICAL: Encrypt client secret before INSERT
    encryptedSecret, err := crypto.Encrypt(input.ClientSecret, r.encryptionKey)
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrEncryptionFailed, err)
    }

    query := `
        INSERT INTO altalune_oauth_providers (
            public_id, provider_type, client_id, client_secret,
            redirect_url, scopes, enabled, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at, updated_at
    `

    // Use encryptedSecret in query, not plaintext
    // ...
}
```

**Query (never returns secret):**
```go
func (r *Repo) Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[OAuthProvider], error) {
    baseQuery := `
        SELECT
            id, public_id, provider_type, client_id,
            redirect_url, scopes, enabled,
            created_at, updated_at
            -- NOTE: client_secret column NOT selected!
        FROM altalune_oauth_providers
        WHERE 1=1
    `
    // ... filtering, sorting, pagination logic ...
}
```

**RevealClientSecret (explicit decryption):**
```go
func (r *Repo) RevealClientSecret(ctx context.Context, publicID string) (string, error) {
    query := `
        SELECT client_secret
        FROM altalune_oauth_providers
        WHERE public_id = $1
    `

    var encryptedSecret string
    err := r.db.QueryRowContext(ctx, query, publicID).Scan(&encryptedSecret)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", ErrOAuthProviderNotFound
        }
        return "", fmt.Errorf("get encrypted secret: %w", err)
    }

    // Decrypt and return plaintext
    plaintext, err := crypto.Decrypt(encryptedSecret, r.encryptionKey)
    if err != nil {
        return "", fmt.Errorf("%w: %v", ErrDecryptionFailed, err)
    }

    return plaintext, nil
}
```

### 4. Service (`service.go`) - Business Logic

**Provider type immutability:**
```go
func (s *Service) UpdateOAuthProvider(ctx context.Context, req *altalunev1.UpdateOAuthProviderRequest) (*altalunev1.UpdateOAuthProviderResponse, error) {
    // ... validation ...

    // Get existing provider to preserve provider_type
    existingProvider, err := s.repo.GetByID(ctx, req.ProviderId)
    if err != nil {
        if err == ErrOAuthProviderNotFound {
            return nil, altalune.NewOAuthProviderNotFoundError(req.ProviderId)
        }
        return nil, altalune.NewUnexpectedError("get existing provider", err)
    }

    // Update (provider_type is NOT changeable)
    input := &UpdateOAuthProviderInput{
        PublicID:     req.ProviderId,
        ClientID:     strings.TrimSpace(req.ClientId),
        ClientSecret: req.ClientSecret,  // Optional
        RedirectURL:  strings.TrimSpace(req.RedirectUrl),
        Scopes:       strings.TrimSpace(req.Scopes),
        Enabled:      req.Enabled,
    }

    result, err := s.repo.Update(ctx, input)
    // ...

    // Return with preserved provider_type
    return &altalunev1.UpdateOAuthProviderResponse{
        Provider: result.ToOAuthProvider(existingProvider.ProviderType).ToProto(),
        Message:  "OAuth provider updated successfully",
    }, nil
}
```

**Duplicate check:**
```go
func (s *Service) CreateOAuthProvider(ctx context.Context, req *altalunev1.CreateOAuthProviderRequest) (*altalunev1.CreateOAuthProviderResponse, error) {
    // ... validation ...

    providerType := ProviderTypeFromProto(req.ProviderType)

    // Check duplicate provider_type
    existingProvider, err := s.repo.GetByProviderType(ctx, providerType)
    if err != nil && err != ErrOAuthProviderNotFound {
        return nil, altalune.NewUnexpectedError("check duplicate", err)
    }
    if existingProvider != nil {
        return nil, altalune.NewOAuthProviderDuplicateTypeError(string(providerType))
    }

    // ... create ...
}
```

### 5. Container Wiring (`internal/container/container.go`)

**Add to Container struct:**
```go
type Container struct {
    // ... existing fields ...

    oauthProviderRepo    oauth_provider.Repository
    oauthProviderService altalunev1.OAuthProviderServiceServer
}
```

**Initialize repository with encryption key:**
```go
func (c *Container) initRepositories() error {
    // ... existing repos ...

    // OAuth Provider Repository with encryption key from config
    c.oauthProviderRepo = oauth_provider.NewRepo(
        c.db,
        c.config.GetIAMEncryptionKey(),  // Inject encryption key
    )

    return nil
}
```

**Initialize service:**
```go
func (c *Container) initServices() error {
    validator, err := protovalidate.New()
    if err != nil {
        return fmt.Errorf("create validator: %w", err)
    }

    // ... existing services ...

    // OAuth Provider Service
    c.oauthProviderService = oauth_provider.NewService(
        validator,
        c.logger,
        c.oauthProviderRepo,
    )

    return nil
}
```

**Getter:**
```go
func (c *Container) GetOAuthProviderService() altalunev1.OAuthProviderServiceServer {
    return c.oauthProviderService
}
```

### 6. Server Registration

**File:** `internal/server/server.go` (or equivalent HTTP server setup)

```go
// Register OAuth Provider Service
oauthProviderPath, oauthProviderHandler := altalunev1connect.NewOAuthProviderServiceHandler(
    oauth_provider.NewHandler(container.GetOAuthProviderService()),
    connect.WithInterceptors(middleware.ErrorInterceptor()),
)
mux.Handle(oauthProviderPath, oauthProviderHandler)
```

## Files to Create

- `api/proto/altalune/v1/oauth_provider.proto` - Proto schema with 6 RPCs
- `internal/domain/oauth_provider/model.go` - Domain models and types
- `internal/domain/oauth_provider/interface.go` - Repository interface
- `internal/domain/oauth_provider/mapper.go` - Proto ↔ domain conversions
- `internal/domain/oauth_provider/errors.go` - Domain errors
- `internal/domain/oauth_provider/repo.go` - PostgreSQL implementation
- `internal/domain/oauth_provider/service.go` - Business logic
- `internal/domain/oauth_provider/handler.go` - Connect-RPC handlers

## Files to Modify

- `errors.go` - Add error codes (608XX) and constructor functions
- `internal/container/container.go` - Add oauth_provider repo and service
- `internal/server/server.go` - Register OAuthProviderService
- `buf.gen.yaml` - Verify configuration for code generation

## Testing Requirements

### Proto Generation

```bash
# Generate proto code
buf generate

# Verify generated files exist
ls gen/altalune/v1/oauth_pb.go
ls frontend/gen/altalune/v1/oauth_pb.ts
```

### Repository Tests

**Manual testing with psql:**
```sql
-- Test encryption (create)
-- Client secret should be base64-encoded ciphertext, not plaintext

-- Test unique constraint (duplicate provider_type)
-- Should return DuplicateProviderType error

-- Test reveal (decrypt)
-- Should return plaintext secret
```

### API Tests (via curl)

**Create:**
```bash
curl -X POST http://localhost:8080/altalune.v1.OAuthProviderService/CreateOAuthProvider \
  -H "Content-Type: application/json" \
  -d '{
    "provider_type": "PROVIDER_TYPE_GOOGLE",
    "client_id": "test-client-id",
    "client_secret": "test-secret-123",
    "redirect_url": "http://localhost:3000/auth/callback",
    "scopes": "openid,email,profile",
    "enabled": true
  }'
```

**Query:**
```bash
curl http://localhost:8080/altalune.v1.OAuthProviderService/QueryOAuthProviders \
  -H "Content-Type: application/json" \
  -d '{"query": {"pagination": {"page": 1, "page_size": 10}}}'
```

**Reveal Secret:**
```bash
curl -X POST http://localhost:8080/altalune.v1.OAuthProviderService/RevealClientSecret \
  -H "Content-Type: application/json" \
  -d '{"provider_id": "<provider-id>"}'
```

**Update (without changing secret):**
```bash
curl -X POST http://localhost:8080/altalune.v1.OAuthProviderService/UpdateOAuthProvider \
  -H "Content-Type: application/json" \
  -d '{
    "provider_id": "<provider-id>",
    "client_id": "updated-client-id",
    "client_secret": "",
    "redirect_url": "http://localhost:3000/auth/callback",
    "scopes": "openid,email",
    "enabled": true
  }'
```

## Commands to Run

```bash
# 1. Create proto file
touch api/proto/altalune/v1/oauth_provider.proto

# 2. Generate proto code
buf generate

# 3. Create domain directory
mkdir -p internal/domain/oauth_provider

# 4. Create all 7 domain files
touch internal/domain/oauth_provider/{model,interface,mapper,errors,repo,service,handler}.go

# 5. Build backend
make build

# 6. Run database migration (if not done in T13)
./bin/app migrate -c config.yaml

# 7. Set encryption key
export IAM_ENCRYPTION_KEY="$(openssl rand -base64 32)"

# 8. Run server
./bin/app serve -c config.yaml

# 9. Test APIs
curl http://localhost:8080/altalune.v1.OAuthProviderService/QueryOAuthProviders
```

## Validation Checklist

### Proto Schema
- [ ] 6 RPCs defined (Query, Create, Get, Update, Delete, Reveal)
- [ ] ProviderType enum with 4 values
- [ ] OAuthProvider message has client_secret_set (NOT client_secret)
- [ ] Validation rules with buf.validate
- [ ] `buf generate` runs without errors

### Domain Layer
- [ ] All 7 files created
- [ ] Repository interface matches implementation
- [ ] Encryption key injected via constructor (not global)
- [ ] Create encrypts client_secret before INSERT
- [ ] Update re-encrypts if client_secret provided
- [ ] Query never selects client_secret column
- [ ] RevealClientSecret decrypts and returns plaintext

### Business Logic
- [ ] Duplicate provider_type check before create
- [ ] Provider type immutability enforced in update
- [ ] Optional client_secret in update (empty = keep existing)
- [ ] Trim whitespace from inputs
- [ ] Error handling for encryption/decryption failures

### Container & Server
- [ ] OAuth provider repo initialized with encryption key
- [ ] OAuth provider service initialized
- [ ] Service registered in HTTP server
- [ ] Error codes added (608XX)

## Definition of Done

- [ ] Proto schema complete and generates code
- [ ] All 7 domain files implemented
- [ ] Repository encrypts/decrypts correctly
- [ ] Service validates business rules
- [ ] Error codes added and used
- [ ] Container wiring complete
- [ ] Service registered in server
- [ ] All 6 RPCs tested via curl
- [ ] Encryption verified in database (ciphertext, not plaintext)
- [ ] Reveal returns plaintext secret
- [ ] Duplicate provider_type returns error
- [ ] Provider type immutability works

## Dependencies

**Upstream:**
- T13 (Database) - Required for testing
- T14 (Crypto) - Required for encryption/decryption

**Downstream:**
- T16 (Frontend Foundation) - Needs generated proto types

## Risk Factors

- **High Risk**: Encryption integration errors
  - **Mitigation**: Thorough testing of encrypt/decrypt flow
  - **Mitigation**: Unit tests for crypto package (T14)

- **Medium Risk**: Provider type immutability enforcement
  - **Mitigation**: Preserve provider_type in service layer
  - **Mitigation**: Test update doesn't change provider_type

- **Medium Risk**: Secret exposure in logs/responses
  - **Mitigation**: Never select client_secret in Query/Get
  - **Mitigation**: Separate RevealClientSecret RPC
  - **Mitigation**: Code review focusing on security

## Notes

### Query Implementation Pattern

Follow existing pattern from `internal/domain/user/repo.go`:
- Build base query with WHERE conditions
- Keyword search with LOWER() + LIKE
- Column filters with IN clauses
- Count total rows before pagination
- Add ORDER BY clause
- Add LIMIT/OFFSET for pagination
- Get distinct values for filters

### Service Pattern

Follow existing pattern from `internal/domain/user/service.go`:
- Validate request with protovalidate
- Additional business logic validation
- Log all operations
- Return structured errors
- Convert domain results to proto messages

### Handler Pattern

Follow existing pattern from `internal/domain/user/handler.go`:
- Thin wrapper around service
- Convert Connect request to service call
- Return Connect response

### Mapper Pattern

Follow existing pattern from `internal/domain/user/mapper.go`:
- ProviderType enum conversion (domain ↔ proto)
- Filters map conversion
- Slice conversions (domain slice to proto slice)

### Testing Encrypted Data

**Verify in database:**
```sql
SELECT public_id, provider_type, client_secret
FROM altalune_oauth_providers;

-- client_secret should look like: "xK3jF9mL2pQ8rT1vW4aZ7bC0dE5gH9iN..."
-- NOT: "my-actual-secret"
```

### Security Reminders

- ❌ Never log client secrets (encrypted or plaintext)
- ❌ Never return actual secret in Query/Get responses
- ✅ Only RevealClientSecret returns plaintext
- ✅ Encryption key never exposed in logs
- ✅ Use client_secret_set boolean flag in responses
