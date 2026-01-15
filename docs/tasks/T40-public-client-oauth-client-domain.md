# Task T40: Public Client OAuth Client Domain Updates

**Story Reference:** US10-public-oauth-clients.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 2-3 hours
**Prerequisites:** T39-public-client-database-proto

## Objective

Update the `oauth_client` domain layer to handle the new `confidential` field, including model updates, repository changes for nullable secret, service logic for PKCE enforcement, and mapper updates.

## Acceptance Criteria

- [ ] Model structs include `Confidential` field
- [ ] Repository handles nullable `client_secret_hash` for public clients
- [ ] Repository skips secret generation for public clients
- [ ] Service enforces `pkce_required=true` for public clients
- [ ] Service returns empty `client_secret` for public clients
- [ ] Mapper converts `confidential` field to/from proto
- [ ] Update operations prevent changing client type (immutable)

## Technical Requirements

### Model Updates

File: `internal/domain/oauth_client/model.go`

Add `Confidential bool` to:
- `OAuthClientQueryResult`
- `OAuthClient`
- `CreateOAuthClientInput`

```go
type OAuthClientQueryResult struct {
    // ... existing fields
    Confidential bool
}

type OAuthClient struct {
    // ... existing fields
    Confidential bool
}

type CreateOAuthClientInput struct {
    // ... existing fields
    Confidential bool
}
```

### Repository Updates

File: `internal/domain/oauth_client/repo.go`

#### Create Method Changes

```go
func (r *repo) Create(ctx context.Context, input *CreateOAuthClientInput) (*CreateOAuthClientResult, error) {
    // ... existing setup

    var clientSecret string
    var hashedSecret *string // Now pointer (nullable)

    // Only generate secret for confidential clients
    if input.Confidential {
        clientSecret = generateSecureRandom(32)
        hash, err := password.HashPassword(clientSecret, password.HashOption{
            Iterations: 2,
            Memory:     64 * 1024,
            Threads:    4,
            Len:        32,
        })
        if err != nil {
            return nil, fmt.Errorf("hash client secret: %w", err)
        }
        hashedSecret = &hash
    }

    // Force PKCE for public clients
    pkceRequired := input.PKCERequired
    if !input.Confidential {
        pkceRequired = true
    }

    // Update INSERT query to include confidential, use hashedSecret (can be nil)
    insertQuery := `
        INSERT INTO altalune_oauth_clients (
            public_id, name, client_id, client_secret_hash,
            redirect_uris, pkce_required, is_default, confidential
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at
    `
    // ... execute with hashedSecret
}
```

#### Query Methods

Update all SELECT queries to include `confidential` column:
- `Query()` - Add to SELECT and scan
- `GetByPublicID()` - Add to SELECT and scan
- `GetByClientID()` - Add to SELECT and scan

### Service Updates

File: `internal/domain/oauth_client/service.go`

#### CreateOAuthClient

```go
func (s *Service) CreateOAuthClient(ctx context.Context, req *altalunev1.CreateOAuthClientRequest) (*altalunev1.CreateOAuthClientResponse, error) {
    // ... validation

    // Force PKCE for public clients
    pkceRequired := req.PkceRequired
    if !req.Confidential {
        pkceRequired = true
    }

    input := &CreateOAuthClientInput{
        Name:          strings.TrimSpace(req.Name),
        RedirectURIs:  req.RedirectUris,
        PKCERequired:  pkceRequired,
        AllowedScopes: req.AllowedScopes,
        Confidential:  req.Confidential,
    }

    result, err := s.oauthClientRepo.Create(ctx, input)
    if err != nil {
        // ... error handling
    }

    // client_secret is empty for public clients
    return &altalunev1.CreateOAuthClientResponse{
        Client:       result.Client.ToOAuthClientProto(),
        ClientSecret: result.ClientSecret,
        Message:      "OAuth client created successfully",
    }, nil
}
```

#### UpdateOAuthClient

Add check to prevent changing client type:

```go
func (s *Service) UpdateOAuthClient(ctx context.Context, req *altalunev1.UpdateOAuthClientRequest) (*altalunev1.UpdateOAuthClientResponse, error) {
    // Get existing client
    existing, err := s.oauthClientRepo.GetByPublicID(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    // Client type is immutable - don't allow changes via update
    // (Proto update request should not include confidential field)

    // For public clients, enforce PKCE stays enabled
    if !existing.Confidential && req.PkceRequired != nil && !*req.PkceRequired {
        return nil, altalune.NewInvalidArgumentError("PKCE cannot be disabled for public clients")
    }

    // ... rest of update logic
}
```

### Mapper Updates

File: `internal/domain/oauth_client/mapper.go`

```go
func (c *OAuthClient) ToOAuthClientProto() *altalunev1.OAuthClient {
    return &altalunev1.OAuthClient{
        Id:              c.ID,
        Name:            c.Name,
        ClientId:        c.ClientID.String(),
        RedirectUris:    c.RedirectURIs,
        PkceRequired:    c.PKCERequired,
        IsDefault:       c.IsDefault,
        ClientSecretSet: c.Confidential, // true only if confidential client
        Confidential:    c.Confidential,
        AllowedScopes:   []string{},
        CreatedAt:       timestamppb.New(c.CreatedAt),
        UpdatedAt:       timestamppb.New(c.UpdatedAt),
    }
}
```

### Error Updates

File: `internal/domain/oauth_client/errors.go`

```go
var (
    // ... existing errors
    ErrClientTypeImmutable = errors.New("client type cannot be changed after creation")
    ErrPKCERequiredForPublicClient = errors.New("PKCE cannot be disabled for public clients")
)
```

## Implementation Details

### Nullable Secret Handling

The repository must handle nullable `client_secret_hash`:

```go
// In Query/Get methods, use sql.NullString for scanning
var secretHash sql.NullString
err := row.Scan(&secretHash, ...)

// Convert to pointer for model
var secretHashPtr *string
if secretHash.Valid {
    secretHashPtr = &secretHash.String
}
```

### PKCE Enforcement Logic

- **Create**: If `confidential=false`, force `pkce_required=true`
- **Update**: If client is public, reject attempts to disable PKCE

## Files to Modify

- `internal/domain/oauth_client/model.go` - Add Confidential field
- `internal/domain/oauth_client/repo.go` - Handle nullable secret, include confidential
- `internal/domain/oauth_client/service.go` - Enforce PKCE for public, immutable type
- `internal/domain/oauth_client/mapper.go` - Map confidential field
- `internal/domain/oauth_client/errors.go` - Add new errors

## Testing Requirements

- Create public client via API, verify no secret returned
- Create confidential client via API, verify secret returned
- Verify public client has pkce_required=true in database
- Attempt to update public client to disable PKCE, verify rejection
- Query clients, verify confidential field is returned

## Commands to Run

```bash
# Build after changes
make build

# Run the server
./bin/app serve -c config.yaml

# Test create public client (via frontend or curl)
```

## Validation Checklist

- [ ] Model structs compile with new field
- [ ] Repository queries include confidential column
- [ ] Service enforces PKCE for public clients
- [ ] Mapper correctly sets confidential and client_secret_set
- [ ] Error messages are clear and actionable
- [ ] No regressions in existing confidential client flow

## Definition of Done

- [ ] All model structs updated with Confidential field
- [ ] Repository handles nullable secret correctly
- [ ] Service enforces PKCE for public clients on create
- [ ] Service prevents PKCE disable for public clients on update
- [ ] Mapper converts confidential field to proto
- [ ] Build succeeds without errors
- [ ] Existing confidential client tests still pass

## Dependencies

- T39: Database migration must be applied first
- Generated proto code from T39

## Risk Factors

- **Medium Risk**: Nullable column handling requires careful SQL scanning
- **Low Risk**: Changes are additive, default behavior preserved for confidential clients

## Notes

- The `ClientSecretSet` proto field should be `true` only for confidential clients
- Public clients will have `ClientSecretSet=false` and `Confidential=false`
- Confidential clients will have `ClientSecretSet=true` and `Confidential=true`
- Consider logging when public clients are created for audit purposes
