# Task T2: API Key Backend Domain Implementation

**Story Reference:** US1-api-keys-crud.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 6-8 hours
**Prerequisites:** T1-api-key-protobuf-schema

## Objective

Implement the complete API Key domain following the established 7-file domain pattern with full CRUD operations.

## Acceptance Criteria

- [ ] Create all 7 domain files following established pattern
- [ ] Implement secure API key generation with proper entropy
- [ ] Support all CRUD operations with proper validation
- [ ] Handle project-specific partitioning correctly
- [ ] Implement proper error handling and logging
- [ ] Follow dual ID system (internal int64 + public nanoid)
- [ ] Include comprehensive unit tests
- [ ] Support query operations with filtering and pagination

## Technical Requirements

### Domain Structure (7-file pattern)

```
internal/domain/api_key/
├── model.go      # Domain models and conversion methods
├── interface.go  # Repository interface definition
├── repo.go       # Repository implementation
├── service.go    # Business logic and Connect-RPC service
├── handler.go    # Connect-RPC handlers
├── errors.go     # Domain-specific errors
└── mapper.go     # Protobuf ↔ Domain mapping
```

### Key Generation Requirements

- Use cryptographically secure random generator
- Format: `sk-` prefix + random string (following OpenAI pattern)
- Must be unique across all projects
- Should use `crypto/rand` for generation

### Repository Operations Required

```go
type Repositor interface {
    Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[ApiKey], error)
    Create(ctx context.Context, input *CreateApiKeyInput) (*CreateApiKeyResult, error)
    GetByID(ctx context.Context, projectID int64, publicID string) (*ApiKey, error)
    GetByKey(ctx context.Context, key string) (*ApiKey, error) // For authentication
    Update(ctx context.Context, input *UpdateApiKeyInput) (*UpdateApiKeyResult, error)
    Delete(ctx context.Context, input *DeleteApiKeyInput) error
}
```

### Service Layer Features

- Project validation using existing project domain
- Unique name validation within project scope
- Expiration date validation (future date, max 2 years)
- Secure key generation with uniqueness check
- Proper error handling with altalune.NewXXXError()
- Comprehensive logging for audit purposes

### Security Considerations

- API key value should never be returned after creation
- Store key securely (consider hashing for future)
- Implement proper access control (project-scoped)
- Audit logging for all operations

## Implementation Details

### Database Queries

- Use existing partitioned table `altalune_project_api_keys`
- Leverage partition key (project_id) in all queries
- Support filtering by name, expiration status
- Include proper indexing usage

### Error Handling

```go
var (
    ErrApiKeyNotFound      = errors.New("api key not found")
    ErrApiKeyAlreadyExists = errors.New("api key with this name already exists")
    ErrApiKeyExpired       = errors.New("api key has expired")
)
```

### Model Definitions

```go
type ApiKey struct {
    ID         string    // Public nanoid
    Name       string
    Expiration time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type CreateApiKeyInput struct {
    ProjectID  int64
    Name       string
    Expiration time.Time
}

type CreateApiKeyResult struct {
    ID         int64     // Internal database ID
    PublicID   string    // Public nanoid
    Name       string
    Key        string    // Generated API key (only in result)
    Expiration time.Time
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

## Files to Create

- `internal/domain/api_key/model.go`
- `internal/domain/api_key/interface.go`
- `internal/domain/api_key/repo.go`
- `internal/domain/api_key/service.go`
- `internal/domain/api_key/handler.go`
- `internal/domain/api_key/errors.go`
- `internal/domain/api_key/mapper.go`

## Files to Modify

- None (new domain)

## Testing Requirements

- Unit tests for all repository methods
- Unit tests for service business logic
- Integration tests with database
- Test key generation uniqueness
- Test project isolation
- Test error scenarios

## Commands to Run

```bash
go build -o ./tmp/test-app cmd/altalune/*.go
go test ./internal/domain/api_key/...
```

## Definition of Done

- [ ] All 7 domain files are implemented
- [ ] All CRUD operations work correctly
- [ ] API key generation is secure and unique
- [ ] Project isolation is enforced
- [ ] Error handling is comprehensive
- [ ] Unit tests have good coverage
- [ ] Code follows established patterns
- [ ] No security vulnerabilities in key handling
- [ ] Logging is comprehensive for audit purposes

## Dependencies

- T1: Protobuf schema must be completed
- Existing project domain for validation
- Existing query infrastructure
- Database partitioning system

## Risk Factors

- **Medium Risk**: Key generation security must be carefully implemented
- **Low Risk**: Following established domain patterns
- **Medium Risk**: Partition handling needs careful attention
- **Low Risk**: Business logic is straightforward CRUD
