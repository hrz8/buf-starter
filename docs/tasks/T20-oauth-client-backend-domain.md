# Task T20: OAuth Client Backend Domain Implementation

**Story Reference:** US6-oauth-client-management.md
**Type:** Backend Implementation
**Priority:** High (P0)
**Estimated Effort:** 6-8 hours
**Prerequisites:** T19-argon2-password-hashing-cli (uses Argon2 package)

## Objective

Implement complete backend domain for OAuth client management following the 7-file domain pattern, with Argon2-hashed client secrets and full CRUD operations.

## Acceptance Criteria

- [ ] Protobuf schema defined for OAuth client service (6 RPCs)
- [ ] 7-file domain pattern implemented
- [ ] CreateClient operation with Argon2 secret hashing
- [ ] QueryClients operation with pagination, filtering, sorting
- [ ] GetClient operation with scope assignment
- [ ] UpdateClient operation with optional secret re-hash
- [ ] DeleteClient operation with default client protection
- [ ] RevealClientSecret operation with audit logging
- [ ] Default client cannot be deleted (enforced)
- [ ] Default client must have PKCE required (enforced)
- [ ] Client name unique within project (validated)
- [ ] Redirect URIs validated (at least one, valid HTTP/HTTPS)
- [ ] Role-based permissions implemented (owner/admin/member/user)
- [ ] Code generated with `buf generate`
- [ ] Service registered in container
- [ ] Handlers registered in server

## Technical Requirements

### Protobuf Schema

**File: `api/proto/altalune/v1/oauth_client.proto`**

```protobuf
syntax = "proto3";

package altalune.v1;

import "buf/validate/validate.proto";
import "google/protobuf/timestamp.proto";
import "altalune/v1/query.proto";

// OAuth Client Service
service OAuthClientService {
  rpc CreateClient(CreateClientRequest) returns (CreateClientResponse) {}
  rpc QueryClients(QueryClientsRequest) returns (QueryClientsResponse) {}
  rpc GetClient(GetClientRequest) returns (GetClientResponse) {}
  rpc UpdateClient(UpdateClientRequest) returns (UpdateClientResponse) {}
  rpc DeleteClient(DeleteClientRequest) returns (DeleteClientResponse) {}
  rpc RevealClientSecret(RevealClientSecretRequest) returns (RevealClientSecretResponse) {}
}

// OAuth Client Message
message OAuthClient {
  string id = 1;                          // Public nanoid
  string project_id = 2;                  // Project public_id
  string name = 3;
  string client_id = 4;                   // UUID
  repeated string redirect_uris = 5;
  bool pkce_required = 6;
  bool is_default = 7;
  bool client_secret_set = 8;             // Boolean flag, NOT actual secret
  repeated string allowed_scopes = 9;     // Scope names
  google.protobuf.Timestamp created_at = 98;
  google.protobuf.Timestamp updated_at = 99;
}

// Create Client Request
message CreateClientRequest {
  string project_id = 1 [(buf.validate.field).required = true];
  string name = 2 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
  repeated string redirect_uris = 3 [(buf.validate.field).repeated = {min_items: 1}];
  bool pkce_required = 4;
  repeated string allowed_scopes = 5;     // Optional scope names
}

message CreateClientResponse {
  OAuthClient client = 1;
  string client_secret = 2;               // ONLY returned during creation
  string message = 3;
}

// Query Clients Request
message QueryClientsRequest {
  QueryParams params = 1;
  string project_id = 2 [(buf.validate.field).required = true];
}

message QueryClientsResponse {
  repeated OAuthClient clients = 1;
  QueryMeta meta = 2;
  string message = 3;
}

// Get Client Request
message GetClientRequest {
  string id = 1 [(buf.validate.field).required = true];
  string project_id = 2 [(buf.validate.field).required = true];
}

message GetClientResponse {
  OAuthClient client = 1;
  string message = 2;
}

// Update Client Request
message UpdateClientRequest {
  string id = 1 [(buf.validate.field).required = true];
  string project_id = 2 [(buf.validate.field).required = true];
  optional string name = 3 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
  repeated string redirect_uris = 4;
  optional bool pkce_required = 5;
  repeated string allowed_scopes = 6;
}

message UpdateClientResponse {
  OAuthClient client = 1;
  string message = 2;
}

// Delete Client Request
message DeleteClientRequest {
  string id = 1 [(buf.validate.field).required = true];
  string project_id = 2 [(buf.validate.field).required = true];
}

message DeleteClientResponse {
  string message = 1;
}

// Reveal Client Secret Request
message RevealClientSecretRequest {
  string id = 1 [(buf.validate.field).required = true];
  string project_id = 2 [(buf.validate.field).required = true];
}

message RevealClientSecretResponse {
  string client_secret = 1;
  string message = 2;
}
```

### Domain Implementation (7-File Pattern)

#### 1. model.go - Domain Models

```go
package oauth_client

import (
    "time"
    "github.com/google/uuid"
)

// OAuthClient domain model
type OAuthClient struct {
    ID              int64       // Internal ID (database)
    PublicID        string      // Public nanoid
    ProjectID       int64
    ProjectPublicID string
    Name            string
    ClientID        uuid.UUID   // UUID for OAuth flow
    RedirectURIs    []string
    PKCERequired    bool
    IsDefault       bool
    AllowedScopes   []string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// CreateOAuthClientInput
type CreateOAuthClientInput struct {
    ProjectID      int64
    Name           string
    RedirectURIs   []string
    PKCERequired   bool
    AllowedScopes  []string
}

// CreateOAuthClientResult
type CreateOAuthClientResult struct {
    Client       *OAuthClient
    ClientSecret string  // Plaintext secret (ONLY for creation response)
}

// UpdateOAuthClientInput
type UpdateOAuthClientInput struct {
    PublicID       string
    ProjectID      int64
    Name           *string
    RedirectURIs   []string
    PKCERequired   *bool
    AllowedScopes  []string
}

// QueryOAuthClientsResult
type QueryOAuthClientsResult struct {
    Clients    []*OAuthClient
    TotalRows  int64
}
```

#### 2. interface.go - Repository Interface

```go
package oauth_client

import (
    "context"
    "github.com/hrz8/altalune/internal/shared/query"
)

type Repository interface {
    Create(ctx context.Context, input *CreateOAuthClientInput) (*CreateOAuthClientResult, error)
    Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[OAuthClient], error)
    GetByPublicID(ctx context.Context, projectID int64, publicID string) (*OAuthClient, error)
    GetByClientID(ctx context.Context, clientID string) (*OAuthClient, error)  // For OAuth flows
    Update(ctx context.Context, input *UpdateOAuthClientInput) (*OAuthClient, error)
    Delete(ctx context.Context, projectID int64, publicID string) error
    RevealClientSecret(ctx context.Context, projectID int64, publicID string) (string, error)
}
```

#### 3. repo.go - PostgreSQL Implementation

**Key Implementation Points**:

```go
package oauth_client

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "github.com/hrz8/altalune/internal/shared/nanoid"
    "github.com/hrz8/altalune/internal/shared/password"
    "github.com/hrz8/altalune/internal/shared/postgres"
    "github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
    db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
    return &repo{db: db}
}

// Create - Generate UUID client_id, hash secret with Argon2
func (r *repo) Create(ctx context.Context, input *CreateOAuthClientInput) (*CreateOAuthClientResult, error) {
    // 1. Generate public ID
    publicID, _ := nanoid.GeneratePublicID()

    // 2. Generate UUID client_id
    clientID := uuid.New()

    // 3. Generate secure random secret (minimum 32 characters)
    clientSecret := generateSecureRandom(32)

    // 4. Hash secret with Argon2id
    hashedSecret, err := password.HashPassword(clientSecret, password.HashOption{
        Iterations: 2,
        Memory:     64 * 1024,
        Threads:    4,
        Len:        32,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to hash client secret: %w", err)
    }

    // 5. Insert into partitioned table
    query := `
        INSERT INTO altalune_oauth_clients (
            project_id, public_id, name, client_id,
            client_secret_hash, redirect_uris, pkce_required
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `

    var id int64
    var createdAt, updatedAt time.Time

    err = r.db.QueryRow(ctx, query,
        input.ProjectID,
        publicID,
        input.Name,
        clientID,
        hashedSecret,
        input.RedirectURIs,
        input.PKCERequired,
    ).Scan(&id, &createdAt, &updatedAt)

    if err != nil {
        if postgres.IsUniqueViolation(err) {
            return nil, ErrOAuthClientAlreadyExists
        }
        return nil, err
    }

    // 6. Return client with PLAINTEXT secret (ONLY time it's returned)
    client := &OAuthClient{
        ID:           id,
        PublicID:     publicID,
        ProjectID:    input.ProjectID,
        Name:         input.Name,
        ClientID:     clientID,
        RedirectURIs: input.RedirectURIs,
        PKCERequired: input.PKCERequired,
        IsDefault:    false,
        CreatedAt:    createdAt,
        UpdatedAt:    updatedAt,
    }

    return &CreateOAuthClientResult{
        Client:       client,
        ClientSecret: clientSecret,  // Plaintext secret
    }, nil
}

// Query - NEVER select client_secret_hash
func (r *repo) Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[OAuthClient], error) {
    // Base query WITHOUT client_secret_hash
    baseQuery := `
        FROM altalune_oauth_clients
        WHERE project_id = $1
    `

    // Count query
    var totalRows int64
    countQuery := "SELECT COUNT(*) " + baseQuery
    _ = r.db.QueryRow(ctx, countQuery, projectID).Scan(&totalRows)

    // Data query with pagination, sorting
    dataQuery := `
        SELECT id, public_id, project_id, name, client_id,
               redirect_uris, pkce_required, is_default,
               created_at, updated_at
    ` + baseQuery + `
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

    // Execute and scan rows...
}

// GetByPublicID - NEVER return client_secret_hash
func (r *repo) GetByPublicID(ctx context.Context, projectID int64, publicID string) (*OAuthClient, error) {
    query := `
        SELECT id, public_id, project_id, name, client_id,
               redirect_uris, pkce_required, is_default,
               created_at, updated_at
        FROM altalune_oauth_clients
        WHERE project_id = $1 AND public_id = $2
    `
    // Scan and return (no secret)
}

// RevealClientSecret - Separate method with audit logging
func (r *repo) RevealClientSecret(ctx context.Context, projectID int64, publicID string) (string, error) {
    query := `
        SELECT client_secret_hash
        FROM altalune_oauth_clients
        WHERE project_id = $1 AND public_id = $2
    `

    var hashedSecret string
    err := r.db.QueryRow(ctx, query, projectID, publicID).Scan(&hashedSecret)
    if err != nil {
        return "", ErrOAuthClientNotFound
    }

    // NOTE: In production, add audit logging here
    // logger.Info("client_secret_revealed", "project_id", projectID, "client_id", publicID)

    // Return hashed secret (frontend will display as-is)
    // Alternative: If we stored encrypted secret, decrypt here
    return hashedSecret, nil
}

// Helper: Generate secure random string
func generateSecureRandom(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    rand.Read(b)
    for i := range b {
        b[i] = charset[b[i]%byte(len(charset))]
    }
    return string(b)
}
```

#### 4. service.go - Business Logic

```go
package oauth_client

import (
    "context"
    protovalidate "github.com/bufbuild/protovalidate-go"
    "github.com/hrz8/altalune/internal/shared/altalune"
    "go.uber.org/zap"
)

type Service struct {
    repo      Repository
    validator *protovalidate.Validator
    logger    *zap.Logger
}

func NewService(repo Repository, validator *protovalidate.Validator, logger *zap.Logger) *Service {
    return &Service{
        repo:      repo,
        validator: validator,
        logger:    logger,
    }
}

// CreateClient
func (s *Service) CreateClient(ctx context.Context, req *CreateClientRequest) (*CreateClientResponse, error) {
    // 1. Validate request
    if err := s.validator.Validate(req); err != nil {
        return nil, altalune.NewInvalidInputError(err.Error())
    }

    // 2. Validate redirect URIs
    if len(req.RedirectUris) == 0 {
        return nil, altalune.NewInvalidInputError("at least one redirect URI required")
    }
    for _, uri := range req.RedirectUris {
        if !isValidRedirectURI(uri) {
            return nil, altalune.NewInvalidInputError(fmt.Sprintf("invalid redirect URI: %s", uri))
        }
    }

    // 3. Create client
    result, err := s.repo.Create(ctx, &CreateOAuthClientInput{
        ProjectID:     projectID,  // From context
        Name:          strings.TrimSpace(req.Name),
        RedirectURIs:  req.RedirectUris,
        PKCERequired:  req.PkceRequired,
        AllowedScopes: req.AllowedScopes,
    })
    if err != nil {
        return nil, err
    }

    // 4. Return client + plaintext secret
    return &CreateClientResponse{
        Client:       result.Client.ToProto(),
        ClientSecret: result.ClientSecret,  // ONLY time secret is returned
        Message:      "OAuth client created successfully",
    }, nil
}

// DeleteClient
func (s *Service) DeleteClient(ctx context.Context, req *DeleteClientRequest) (*DeleteClientResponse, error) {
    // 1. Validate request
    if err := s.validator.Validate(req); err != nil {
        return nil, altalune.NewInvalidInputError(err.Error())
    }

    // 2. Get client to check if default
    client, err := s.repo.GetByPublicID(ctx, projectID, req.Id)
    if err != nil {
        return nil, err
    }

    // 3. Protect default client from deletion
    if client.IsDefault {
        return nil, altalune.NewInvalidInputError("cannot delete default dashboard client")
    }

    // 4. Delete
    if err := s.repo.Delete(ctx, projectID, req.Id); err != nil {
        return nil, err
    }

    return &DeleteClientResponse{
        Message: "OAuth client deleted successfully",
    }, nil
}

// RevealClientSecret
func (s *Service) RevealClientSecret(ctx context.Context, req *RevealClientSecretRequest) (*RevealClientSecretResponse, error) {
    // 1. Validate request
    if err := s.validator.Validate(req); err != nil {
        return nil, altalune.NewInvalidInputError(err.Error())
    }

    // 2. Reveal secret (with audit logging in repo)
    secret, err := s.repo.RevealClientSecret(ctx, projectID, req.Id)
    if err != nil {
        return nil, err
    }

    // 3. Log action (additional service-level logging)
    s.logger.Info("oauth_client_secret_revealed",
        zap.Int64("project_id", projectID),
        zap.String("client_id", req.Id),
        zap.String("user_id", userID),  // From context
    )

    return &RevealClientSecretResponse{
        ClientSecret: secret,
        Message:      "Client secret revealed. This action has been logged.",
    }, nil
}

// Helper: Validate redirect URI
func isValidRedirectURI(uri string) bool {
    // Must be valid HTTP/HTTPS URL
    // Must not contain wildcards
    // localhost allowed for development
    parsed, err := url.Parse(uri)
    if err != nil {
        return false
    }
    if parsed.Scheme != "http" && parsed.Scheme != "https" {
        return false
    }
    if strings.Contains(uri, "*") || strings.Contains(uri, "?") {
        return false
    }
    return true
}
```

#### 5. handler.go - Connect-RPC Handlers

```go
package oauth_client

import (
    "context"
    "connectrpc.com/connect"
    pb "github.com/hrz8/altalune/gen/altalune/v1"
    "github.com/hrz8/altalune/internal/shared/altalune"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) CreateClient(ctx context.Context, req *connect.Request[pb.CreateClientRequest]) (*connect.Response[pb.CreateClientResponse], error) {
    resp, err := h.service.CreateClient(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(resp), nil
}

// Similar handlers for Query, Get, Update, Delete, RevealSecret...
```

#### 6. mapper.go - Proto ↔ Domain Conversions

```go
package oauth_client

import (
    pb "github.com/hrz8/altalune/gen/altalune/v1"
    "google.golang.org/protobuf/types/known/timestamppb"
)

func (c *OAuthClient) ToProto() *pb.OAuthClient {
    return &pb.OAuthClient{
        Id:               c.PublicID,
        ProjectId:        c.ProjectPublicID,
        Name:             c.Name,
        ClientId:         c.ClientID.String(),
        RedirectUris:     c.RedirectURIs,
        PkceRequired:     c.PKCERequired,
        IsDefault:        c.IsDefault,
        ClientSecretSet:  true,  // Always true (never show actual secret)
        AllowedScopes:    c.AllowedScopes,
        CreatedAt:        timestamppb.New(c.CreatedAt),
        UpdatedAt:        timestamppb.New(c.UpdatedAt),
    }
}
```

#### 7. errors.go - Domain Errors (Enhance Existing)

```go
package oauth_client

import (
    "errors"
)

var (
    ErrOAuthClientNotFound          = errors.New("oauth client not found")
    ErrOAuthClientAlreadyExists     = errors.New("oauth client already exists")
    ErrInvalidRedirectURI           = errors.New("invalid redirect URI")
    ErrClientSecretMismatch         = errors.New("client secret mismatch")
    ErrDefaultClientCannotBeDeleted = errors.New("default client cannot be deleted")
    ErrClientBelongsToOtherProject  = errors.New("client belongs to other project")
    ErrPKCECannotBeDisabled         = errors.New("PKCE cannot be disabled for default client")
)
```

## Files to Create

- `api/proto/altalune/v1/oauth_client.proto`
- `internal/domain/oauth_client/model.go`
- `internal/domain/oauth_client/interface.go`
- `internal/domain/oauth_client/repo.go`
- `internal/domain/oauth_client/service.go`
- `internal/domain/oauth_client/handler.go`
- `internal/domain/oauth_client/mapper.go`

## Files to Modify

- `internal/domain/oauth_client/errors.go` (enhance if needed)
- `internal/container/container.go` (register service)
- `internal/server/server.go` (register handlers)

## Commands to Run

```bash
# Generate code from proto
buf generate

# Build backend
make build

# Run server
./bin/app serve -c config.yaml
```

## Validation Checklist

- [ ] Proto schema compiles with `buf generate`
- [ ] Generated Go code exists in `gen/altalune/v1/`
- [ ] All 7 domain files created
- [ ] Repository implements all interface methods
- [ ] Service validates all inputs with protovalidate
- [ ] Default client deletion blocked
- [ ] PKCE enforcement for default client
- [ ] Client secret hashed with Argon2 on creation
- [ ] Client secret NEVER returned in Query/GetByID
- [ ] RevealClientSecret separate method
- [ ] Redirect URIs validated (HTTP/HTTPS only)
- [ ] Service registered in container
- [ ] Handlers registered in server
- [ ] Backend compiles without errors

## Definition of Done

- [ ] Protobuf schema defined and generated
- [ ] 7-file domain pattern implemented
- [ ] All CRUD operations functional
- [ ] Default client protection enforced
- [ ] Argon2 secret hashing working
- [ ] RevealClientSecret with audit logging
- [ ] Redirect URI validation working
- [ ] Role-based permissions checked
- [ ] Code follows established patterns
- [ ] Service registered in container
- [ ] Handlers registered in server
- [ ] Manual testing with curl/Bruno successful

## Dependencies

**External**:
- `github.com/google/uuid` - UUID generation
- `golang.org/x/crypto/argon2` - Via T19 password package
- `github.com/bufbuild/protovalidate-go` - Request validation
- Existing: `pgx`, `connect-rpc`, `nanoid`

**Internal**:
- T19: `internal/shared/password` package (Argon2 hashing)
- Existing: `internal/shared/query` (pagination)
- Existing: `internal/shared/nanoid` (public ID generation)
- Existing: `internal/shared/postgres` (DB utilities)

## Risk Factors

- **High Risk**: Partitioned table queries without project_id
  - **Mitigation**: Always include project_id in WHERE clauses
- **Medium Risk**: Accidentally returning client_secret_hash
  - **Mitigation**: Never SELECT it in Query/GetByID
- **Medium Risk**: Default client deletion
  - **Mitigation**: Check is_default flag in service layer

## Notes

### Client Secret Security

**CRITICAL**:
1. Hash with Argon2id immediately on creation
2. Return plaintext ONLY in CreateClientResponse
3. NEVER select client_secret_hash in Query/GetByID
4. RevealClientSecret is separate method with audit logging
5. Frontend displays hashed secret (no way to recover plaintext)

### Partitioned Table Considerations

**Always include project_id**:
```sql
-- CORRECT: Partition routing works
SELECT * FROM altalune_oauth_clients WHERE project_id = 1 AND public_id = 'abc123'

-- WRONG: Partition routing fails
SELECT * FROM altalune_oauth_clients WHERE public_id = 'abc123'
```

### Default Client Handling

**Database Level**:
- `is_default BOOLEAN DEFAULT false`
- Only ONE default client per project

**Service Level**:
- Check `is_default` before deletion
- Enforce `pkce_required = true` for default
- Return specific error if delete attempted

**UI Level** (T22):
- Disabled delete button
- Tooltip explanation
- Cannot toggle PKCE off
