# Task T25: OAuth Auth Domain - Authorization Code & Token Management

**Story Reference:** US7-oauth-authorization-server.md
**Type:** Backend Domain
**Priority:** High
**Estimated Effort:** 5-6 hours
**Prerequisites:** None (uses existing database tables from US5)

## Objective

Implement the oauth_auth domain following the 7-file pattern for managing authorization codes, refresh tokens, and user consents. This domain provides the core business logic for OAuth 2.0 token operations.

## Acceptance Criteria

- [ ] Repository interface for authorization codes, refresh tokens, and consents
- [ ] Domain models for all OAuth entities
- [ ] Repository implementation with soft-delete pattern
- [ ] Service layer with business logic for code/token generation and validation
- [ ] Client authentication utilities (Basic Auth parsing, secret verification)
- [ ] Proper error handling with domain-specific errors

## Technical Requirements

### Domain Models (`model.go`)

#### AuthorizationCode
```go
type AuthorizationCode struct {
    ID                  int64
    Code                uuid.UUID
    ClientID            uuid.UUID
    UserID              int64
    RedirectURI         string
    Scope               string
    Nonce               *string
    CodeChallenge       *string
    CodeChallengeMethod *string
    ExpiresAt           time.Time
    ExchangeAt          *time.Time // Soft-delete timestamp
    CreatedAt           time.Time
}
```

#### RefreshToken
```go
type RefreshToken struct {
    ID         int64
    Token      uuid.UUID
    ClientID   uuid.UUID
    UserID     int64
    Scope      string
    Nonce      *string
    ExpiresAt  time.Time
    ExchangeAt *time.Time // Soft-delete timestamp
    CreatedAt  time.Time
}
```

#### UserConsent
```go
type UserConsent struct {
    ID        int64
    UserID    int64
    ClientID  uuid.UUID
    Scope     string
    GrantedAt time.Time
    RevokedAt *time.Time
    CreatedAt time.Time
}
```

### Repository Interface (`interface.go`)

```go
type Repositor interface {
    // Authorization Codes
    CreateAuthorizationCode(ctx context.Context, input CreateAuthCodeInput) (*AuthorizationCode, error)
    GetAuthorizationCodeByCode(ctx context.Context, code uuid.UUID) (*AuthorizationCode, error)
    MarkCodeExchanged(ctx context.Context, code uuid.UUID) error

    // Refresh Tokens
    CreateRefreshToken(ctx context.Context, input CreateRefreshTokenInput) (*RefreshToken, error)
    GetRefreshTokenByToken(ctx context.Context, token uuid.UUID) (*RefreshToken, error)
    MarkRefreshTokenExchanged(ctx context.Context, token uuid.UUID) error

    // User Consents
    GetUserConsent(ctx context.Context, userID int64, clientID uuid.UUID) (*UserConsent, error)
    CreateOrUpdateUserConsent(ctx context.Context, input UserConsentInput) (*UserConsent, error)
}
```

### Service Interface

```go
type Service interface {
    // Authorization Code Operations
    GenerateAuthorizationCode(ctx context.Context, input GenerateAuthCodeInput) (*AuthorizationCode, error)
    ValidateAndExchangeCode(ctx context.Context, code string, clientID uuid.UUID, redirectURI string, codeVerifier *string) (*CodeExchangeResult, error)

    // Token Operations
    GenerateTokenPair(ctx context.Context, userID int64, clientID uuid.UUID, scope string) (*TokenPair, error)
    ValidateAndRefreshToken(ctx context.Context, refreshToken string, clientID uuid.UUID) (*TokenPair, error)

    // Consent Operations
    CheckUserConsent(ctx context.Context, userID int64, clientID uuid.UUID, requestedScope string) (bool, error)
    SaveUserConsent(ctx context.Context, userID int64, clientID uuid.UUID, scope string) error

    // Client Authentication
    AuthenticateClient(ctx context.Context, clientID, clientSecret string) (*OAuthClient, error)
}
```

### Input/Output Types

```go
type CreateAuthCodeInput struct {
    ClientID            uuid.UUID
    UserID              int64
    RedirectURI         string
    Scope               string
    Nonce               *string
    CodeChallenge       *string
    CodeChallengeMethod *string
    ExpiresAt           time.Time
}

type GenerateAuthCodeInput struct {
    ClientID            uuid.UUID
    UserID              int64
    RedirectURI         string
    Scope               string
    Nonce               *string
    CodeChallenge       *string
    CodeChallengeMethod *string
}

type CodeExchangeResult struct {
    UserID int64
    Scope  string
    Nonce  *string
}

type TokenPair struct {
    AccessToken  string
    RefreshToken string
    TokenType    string
    ExpiresIn    int
    Scope        string
}
```

## Implementation Details

### Repository Implementation (`repo.go`)

Key considerations:
- Use soft-delete pattern (set `exchange_at` instead of DELETE)
- Check `exchange_at IS NULL` for active codes/tokens
- Check `expires_at > NOW()` for non-expired codes/tokens

```go
func (r *Repo) GetAuthorizationCodeByCode(ctx context.Context, code uuid.UUID) (*AuthorizationCode, error) {
    query := `
        SELECT id, code, client_id, user_id, redirect_uri, scope, nonce,
               code_challenge, code_challenge_method, expires_at, exchange_at, created_at
        FROM altalune_oauth_authorization_codes
        WHERE code = $1
          AND exchange_at IS NULL
          AND expires_at > NOW()
    `
    // Execute query and scan result
}

func (r *Repo) MarkCodeExchanged(ctx context.Context, code uuid.UUID) error {
    query := `
        UPDATE altalune_oauth_authorization_codes
        SET exchange_at = NOW(), updated_at = NOW()
        WHERE code = $1 AND exchange_at IS NULL
    `
    // Execute update
}
```

### Service Implementation (`service.go`)

The service handles business logic and integrates with:
- oauth_client domain for client verification
- JWT utilities (from T24) for access token generation
- PKCE utilities (from T24) for code challenge verification
- Config for expiry settings

```go
type Service struct {
    repo          Repositor
    oauthClientSvc *oauth_client.Service
    jwtSigner     *jwt.Signer
    config        altalune.Config
    logger        *slog.Logger
}

func (s *Service) ValidateAndExchangeCode(ctx context.Context, codeStr string, clientID uuid.UUID, redirectURI string, codeVerifier *string) (*CodeExchangeResult, error) {
    // 1. Parse code UUID
    code, err := uuid.Parse(codeStr)

    // 2. Get authorization code from database
    authCode, err := s.repo.GetAuthorizationCodeByCode(ctx, code)
    if err != nil {
        return nil, ErrInvalidAuthorizationCode
    }

    // 3. Validate client_id matches
    if authCode.ClientID != clientID {
        return nil, ErrClientMismatch
    }

    // 4. Validate redirect_uri matches
    if authCode.RedirectURI != redirectURI {
        return nil, ErrRedirectURIMismatch
    }

    // 5. Validate PKCE if code_challenge was stored
    if authCode.CodeChallenge != nil {
        if codeVerifier == nil {
            return nil, ErrMissingCodeVerifier
        }
        if !pkce.VerifyCodeChallenge(*codeVerifier, *authCode.CodeChallenge, *authCode.CodeChallengeMethod) {
            return nil, ErrInvalidCodeVerifier
        }
    }

    // 6. Mark code as exchanged (soft-delete)
    if err := s.repo.MarkCodeExchanged(ctx, code); err != nil {
        return nil, err
    }

    return &CodeExchangeResult{
        UserID: authCode.UserID,
        Scope:  authCode.Scope,
        Nonce:  authCode.Nonce,
    }, nil
}
```

### Client Authentication

```go
func (s *Service) AuthenticateClient(ctx context.Context, clientID, clientSecret string) (*oauth_client.OAuthClient, error) {
    // Parse client_id as UUID
    clientUUID, err := uuid.Parse(clientID)
    if err != nil {
        return nil, ErrInvalidClientID
    }

    // Get client from oauth_client domain
    client, err := s.oauthClientSvc.GetByClientID(ctx, clientUUID)
    if err != nil {
        return nil, ErrClientNotFound
    }

    // Verify secret using argon2id (reuse from oauth_client domain)
    if !oauth_client.VerifyClientSecret(clientSecret, client.ClientSecretHash) {
        return nil, ErrInvalidClientSecret
    }

    return client, nil
}
```

### Consent Management

```go
func (s *Service) CheckUserConsent(ctx context.Context, userID int64, clientID uuid.UUID, requestedScope string) (bool, error) {
    consent, err := s.repo.GetUserConsent(ctx, userID, clientID)
    if err != nil {
        return false, nil // No consent exists
    }

    if consent.RevokedAt != nil {
        return false, nil // Consent was revoked
    }

    // Check if all requested scopes are in granted scopes
    requestedScopes := strings.Split(requestedScope, " ")
    grantedScopes := strings.Split(consent.Scope, " ")
    grantedSet := make(map[string]bool)
    for _, s := range grantedScopes {
        grantedSet[s] = true
    }

    for _, s := range requestedScopes {
        if !grantedSet[s] {
            return false, nil // Missing scope, need new consent
        }
    }

    return true, nil
}
```

## Files to Create

- `internal/domain/oauth_auth/interface.go` - Repository and service interfaces
- `internal/domain/oauth_auth/model.go` - Domain models and DTOs
- `internal/domain/oauth_auth/repo.go` - Repository implementation
- `internal/domain/oauth_auth/service.go` - Service implementation
- `internal/domain/oauth_auth/mapper.go` - DTO to model conversions (if needed)

## Files to Modify

- `internal/domain/oauth_auth/errors.go` - Extend with additional error types
- `internal/container/container.go` - Register oauth_auth service

## Testing Requirements

- Unit tests for repository methods
- Unit tests for service business logic
- Test soft-delete pattern (exchange_at timestamp)
- Test expiration validation
- Test PKCE verification in code exchange
- Test consent scope matching logic

## Commands to Run

```bash
# Build and verify no compile errors
make build

# Run tests (when implemented)
go test ./internal/domain/oauth_auth/...
```

## Validation Checklist

- [ ] Authorization codes are UUID v4 format
- [ ] Codes expire after configured duration (default 10 min)
- [ ] Codes can only be exchanged once (soft-delete)
- [ ] Refresh tokens are UUID v4 format
- [ ] Refresh tokens are single-use
- [ ] PKCE validation works for S256 method
- [ ] Consent tracking prevents re-prompting

## Definition of Done

- [ ] Repository interface defined with all required methods
- [ ] Domain models match database schema
- [ ] Repository implementation uses soft-delete pattern
- [ ] Service implements code generation with PKCE support
- [ ] Service implements code exchange with all validations
- [ ] Service implements token refresh flow
- [ ] Consent checking compares scope sets correctly
- [ ] Client authentication verifies argon2id hashed secrets
- [ ] Registered in container for dependency injection
- [ ] Code follows established domain patterns

## Dependencies

- `internal/domain/oauth_client` - For client lookup and secret verification
- `internal/shared/jwt` - For access token generation (from T24)
- `internal/shared/pkce` - For PKCE verification (from T24)
- Database tables from US5 migration

## Risk Factors

- **Low Risk**: Standard CRUD operations following existing patterns
- **Medium Risk**: PKCE verification must match RFC 7636 exactly

## Notes

- The service will depend on T24 utilities for JWT and PKCE, but repository can be implemented independently
- Soft-delete pattern preserves audit trail - DO NOT hard delete codes/tokens
- Consider adding cleanup job for expired codes/tokens in future
- Scope comparison should be case-sensitive per OAuth 2.0 spec
