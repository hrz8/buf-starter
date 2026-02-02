# Task T62: Backend Auth Package

**Story Reference:** US15-authorization-rbac.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T61 (Database + JWT Claims)

## Objective

Create the `internal/auth/` package with JWKS fetcher, JWT validator, AuthContext, and Authorizer helper functions for resource server authorization.

## Acceptance Criteria

- [x] JWKS fetcher implemented with in-memory caching
- [x] JWT validator uses public key from JWKS
- [x] Cache refreshes on TTL expiry (1 hour default)
- [x] Cache refreshes on signature validation failure
- [x] Rate limiting on JWKS refresh (max 3/minute)
- [x] AuthContext struct holds parsed JWT claims
- [x] Authorizer helper provides permission checking functions
- [x] Configuration extended with JWKS settings

## Technical Requirements

### Directory Structure

```
internal/
├── auth/
│   ├── context.go        # Request context with claims
│   ├── jwks.go           # JWKS fetching and caching
│   ├── validator.go      # JWT validation logic
│   └── authorizer.go     # Authorization helper functions
```

### JWKS Cache Implementation

Create `internal/auth/jwks.go`:

```go
package auth

import (
    "context"
    "crypto/rsa"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "math/big"
    "net/http"
    "sync"
    "time"
)

// JWKSResponse represents the JWKS endpoint response
type JWKSResponse struct {
    Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
    Kty string `json:"kty"`
    Kid string `json:"kid"`
    Use string `json:"use"`
    Alg string `json:"alg"`
    N   string `json:"n"`
    E   string `json:"e"`
}

// JWKSCache manages caching of JWKS keys
type JWKSCache struct {
    mu              sync.RWMutex
    keys            map[string]*rsa.PublicKey  // kid -> public key
    lastFetch       time.Time
    ttl             time.Duration
    refreshCount    int
    lastRefreshMin  time.Time
    refreshLimit    int
}

// NewJWKSCache creates a new JWKS cache
func NewJWKSCache(ttl time.Duration, refreshLimit int) *JWKSCache {
    return &JWKSCache{
        keys:         make(map[string]*rsa.PublicKey),
        ttl:          ttl,
        refreshLimit: refreshLimit,
    }
}

// GetKey returns the public key for the given kid
func (c *JWKSCache) GetKey(kid string) (*rsa.PublicKey, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    key, ok := c.keys[kid]
    return key, ok
}

// SetKeys stores the fetched keys
func (c *JWKSCache) SetKeys(keys map[string]*rsa.PublicKey) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.keys = keys
    c.lastFetch = time.Now()
}

// IsExpired checks if the cache has expired
func (c *JWKSCache) IsExpired() bool {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return time.Since(c.lastFetch) > c.ttl
}

// CanRefresh checks if refresh is allowed (rate limiting)
func (c *JWKSCache) CanRefresh() bool {
    c.mu.Lock()
    defer c.mu.Unlock()

    now := time.Now()
    if now.Sub(c.lastRefreshMin) > time.Minute {
        c.refreshCount = 0
        c.lastRefreshMin = now
    }

    if c.refreshCount >= c.refreshLimit {
        return false
    }

    c.refreshCount++
    return true
}

// JWKSFetcher fetches and parses JWKS from the auth server
type JWKSFetcher struct {
    url        string
    httpClient *http.Client
    cache      *JWKSCache
}

// NewJWKSFetcher creates a new JWKS fetcher
func NewJWKSFetcher(url string, cacheTTL time.Duration, refreshLimit int) *JWKSFetcher {
    return &JWKSFetcher{
        url:        url,
        httpClient: &http.Client{Timeout: 10 * time.Second},
        cache:      NewJWKSCache(cacheTTL, refreshLimit),
    }
}

// GetPublicKey returns the public key for the given kid
func (f *JWKSFetcher) GetPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
    // Try cache first
    if key, ok := f.cache.GetKey(kid); ok && !f.cache.IsExpired() {
        return key, nil
    }

    // Fetch fresh keys
    if err := f.Refresh(ctx); err != nil {
        return nil, err
    }

    key, ok := f.cache.GetKey(kid)
    if !ok {
        return nil, fmt.Errorf("key not found: %s", kid)
    }

    return key, nil
}

// Refresh fetches keys from JWKS endpoint
func (f *JWKSFetcher) Refresh(ctx context.Context) error {
    if !f.cache.CanRefresh() {
        return fmt.Errorf("JWKS refresh rate limit exceeded")
    }

    req, err := http.NewRequestWithContext(ctx, "GET", f.url, nil)
    if err != nil {
        return err
    }

    resp, err := f.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("JWKS fetch failed: status %d", resp.StatusCode)
    }

    var jwks JWKSResponse
    if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
        return err
    }

    keys := make(map[string]*rsa.PublicKey)
    for _, jwk := range jwks.Keys {
        if jwk.Kty != "RSA" {
            continue
        }

        pubKey, err := parseRSAPublicKey(jwk)
        if err != nil {
            continue
        }

        keys[jwk.Kid] = pubKey
    }

    f.cache.SetKeys(keys)
    return nil
}

// ForceRefresh forces a refresh on validation failure (key rotation)
func (f *JWKSFetcher) ForceRefresh(ctx context.Context) error {
    return f.Refresh(ctx)
}

func parseRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
    nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
    if err != nil {
        return nil, err
    }

    eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
    if err != nil {
        return nil, err
    }

    n := new(big.Int).SetBytes(nBytes)
    e := int(new(big.Int).SetBytes(eBytes).Int64())

    return &rsa.PublicKey{N: n, E: e}, nil
}
```

### JWT Validator Implementation

Create `internal/auth/validator.go`:

```go
package auth

import (
    "context"
    "fmt"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

// AccessTokenClaims mirrors the JWT claims structure
type AccessTokenClaims struct {
    jwt.RegisteredClaims
    Scope         string            `json:"scope,omitempty"`
    Email         string            `json:"email,omitempty"`
    Name          string            `json:"name,omitempty"`
    Perms         []string          `json:"perms"`
    Memberships   map[string]string `json:"memberships,omitempty"`
    EmailVerified bool              `json:"email_verified"`
}

// JWTValidator validates JWT tokens using JWKS
type JWTValidator struct {
    fetcher   *JWKSFetcher
    issuer    string
    audiences []string // Optional audience validation
}

// NewJWTValidator creates a new JWT validator
func NewJWTValidator(jwksURL, issuer string, audiences []string, cacheTTL int, refreshLimit int) *JWTValidator {
    ttl := time.Duration(cacheTTL) * time.Second
    if cacheTTL == 0 {
        ttl = time.Hour // Default 1 hour
    }
    if refreshLimit == 0 {
        refreshLimit = 3 // Default 3 per minute
    }

    return &JWTValidator{
        fetcher:   NewJWKSFetcher(jwksURL, ttl, refreshLimit),
        issuer:    issuer,
        audiences: audiences,
    }
}

// Validate validates a JWT token and returns the claims
func (v *JWTValidator) Validate(ctx context.Context, tokenString string) (*AccessTokenClaims, error) {
    // Parse token without validation to get kid
    parser := jwt.NewParser()
    token, _, err := parser.ParseUnverified(tokenString, &AccessTokenClaims{})
    if err != nil {
        return nil, fmt.Errorf("invalid token format: %w", err)
    }

    kid, ok := token.Header["kid"].(string)
    if !ok {
        return nil, fmt.Errorf("missing kid in token header")
    }

    // Get public key from JWKS
    publicKey, err := v.fetcher.GetPublicKey(ctx, kid)
    if err != nil {
        return nil, fmt.Errorf("failed to get public key: %w", err)
    }

    // Parse and validate token with signature verification
    claims := &AccessTokenClaims{}
    token, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return publicKey, nil
    })

    if err != nil {
        // Try refreshing JWKS on validation failure (key rotation)
        if strings.Contains(err.Error(), "signature") {
            if refreshErr := v.fetcher.ForceRefresh(ctx); refreshErr == nil {
                publicKey, _ = v.fetcher.GetPublicKey(ctx, kid)
                token, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
                    return publicKey, nil
                })
            }
        }

        if err != nil {
            if strings.Contains(err.Error(), "expired") {
                return nil, fmt.Errorf("token has expired")
            }
            return nil, fmt.Errorf("invalid token signature")
        }
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    // Validate issuer
    if v.issuer != "" && claims.Issuer != v.issuer {
        return nil, fmt.Errorf("invalid token issuer")
    }

    // Validate audience (optional)
    if len(v.audiences) > 0 {
        valid := false
        for _, aud := range v.audiences {
            for _, tokenAud := range claims.Audience {
                if aud == tokenAud {
                    valid = true
                    break
                }
            }
            if valid {
                break
            }
        }
        if !valid {
            return nil, fmt.Errorf("invalid token audience")
        }
    }

    return claims, nil
}

// RefreshJWKS forces a refresh of the JWKS cache
func (v *JWTValidator) RefreshJWKS(ctx context.Context) error {
    return v.fetcher.ForceRefresh(ctx)
}
```

### Auth Context Implementation

Create `internal/auth/context.go`:

```go
package auth

import (
    "context"
)

type contextKey string

const authContextKey contextKey = "auth_context"

// AuthContext holds the authenticated user's information from JWT
type AuthContext struct {
    UserID          string            // JWT subject (user public_id)
    Email           string
    Name            string
    Permissions     []string
    Memberships     map[string]string // project_public_id -> role
    EmailVerified   bool
    IsAuthenticated bool
}

// FromContext extracts AuthContext from request context
func FromContext(ctx context.Context) *AuthContext {
    auth, ok := ctx.Value(authContextKey).(*AuthContext)
    if !ok {
        return &AuthContext{IsAuthenticated: false}
    }
    return auth
}

// WithAuthContext adds AuthContext to request context
func WithAuthContext(ctx context.Context, auth *AuthContext) context.Context {
    return context.WithValue(ctx, authContextKey, auth)
}

// NewAuthContextFromClaims creates AuthContext from JWT claims
func NewAuthContextFromClaims(claims *AccessTokenClaims) *AuthContext {
    return &AuthContext{
        UserID:          claims.Subject,
        Email:           claims.Email,
        Name:            claims.Name,
        Permissions:     claims.Perms,
        Memberships:     claims.Memberships,
        EmailVerified:   claims.EmailVerified,
        IsAuthenticated: true,
    }
}
```

### Authorizer Implementation

Create `internal/auth/authorizer.go`:

```go
package auth

import (
    "context"

    "connectrpc.com/connect"
)

const RootPermission = "root"

// Authorizer provides authorization helper functions
type Authorizer struct{}

// NewAuthorizer creates a new Authorizer
func NewAuthorizer() *Authorizer {
    return &Authorizer{}
}

// CheckAuthenticated verifies user is authenticated
func (a *Authorizer) CheckAuthenticated(ctx context.Context) error {
    auth := FromContext(ctx)
    if !auth.IsAuthenticated {
        return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
    }
    return nil
}

// CheckPermission validates user has required permission
// Returns error if unauthorized
func (a *Authorizer) CheckPermission(ctx context.Context, permission string) error {
    auth := FromContext(ctx)

    if !auth.IsAuthenticated {
        return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
    }

    // Superadmin bypass
    if a.IsSuperAdmin(ctx) {
        return nil
    }

    // Check permission
    for _, p := range auth.Permissions {
        if p == permission {
            return nil
        }
    }

    return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("permission denied: requires %s", permission))
}

// CheckProjectAccess validates user has permission AND is project member
// projectID is the project's public_id
func (a *Authorizer) CheckProjectAccess(ctx context.Context, permission string, projectID string) error {
    auth := FromContext(ctx)

    if !auth.IsAuthenticated {
        return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
    }

    // Superadmin bypass
    if a.IsSuperAdmin(ctx) {
        return nil
    }

    // Check permission
    hasPermission := false
    for _, p := range auth.Permissions {
        if p == permission {
            hasPermission = true
            break
        }
    }

    if !hasPermission {
        return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("permission denied: requires %s", permission))
    }

    // Check project membership
    if _, ok := auth.Memberships[projectID]; !ok {
        return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("permission denied: not a member of this project"))
    }

    return nil
}

// CheckProjectMembership validates user is member of project
func (a *Authorizer) CheckProjectMembership(ctx context.Context, projectID string) error {
    auth := FromContext(ctx)

    if !auth.IsAuthenticated {
        return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
    }

    // Superadmin bypass
    if a.IsSuperAdmin(ctx) {
        return nil
    }

    if _, ok := auth.Memberships[projectID]; !ok {
        return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("permission denied: not a member of this project"))
    }

    return nil
}

// IsSuperAdmin checks if user has root permission
func (a *Authorizer) IsSuperAdmin(ctx context.Context) bool {
    auth := FromContext(ctx)
    if !auth.IsAuthenticated {
        return false
    }

    for _, p := range auth.Permissions {
        if p == RootPermission {
            return true
        }
    }
    return false
}

// GetUserProjects returns list of project IDs user has access to
func (a *Authorizer) GetUserProjects(ctx context.Context) []string {
    auth := FromContext(ctx)
    if !auth.IsAuthenticated {
        return nil
    }

    projects := make([]string, 0, len(auth.Memberships))
    for projectID := range auth.Memberships {
        projects = append(projects, projectID)
    }
    return projects
}

// GetProjectRole returns user's role in a specific project
func (a *Authorizer) GetProjectRole(ctx context.Context, projectID string) string {
    auth := FromContext(ctx)
    if !auth.IsAuthenticated {
        return ""
    }
    return auth.Memberships[projectID]
}

// HasPermission checks if user has a specific permission
func (a *Authorizer) HasPermission(ctx context.Context, permission string) bool {
    auth := FromContext(ctx)
    if !auth.IsAuthenticated {
        return false
    }

    // Superadmin has all permissions
    if a.IsSuperAdmin(ctx) {
        return true
    }

    for _, p := range auth.Permissions {
        if p == permission {
            return true
        }
    }
    return false
}

// HasAnyPermission checks if user has any of the specified permissions
func (a *Authorizer) HasAnyPermission(ctx context.Context, permissions []string) bool {
    auth := FromContext(ctx)
    if !auth.IsAuthenticated {
        return false
    }

    if a.IsSuperAdmin(ctx) {
        return true
    }

    for _, required := range permissions {
        for _, p := range auth.Permissions {
            if p == required {
                return true
            }
        }
    }
    return false
}
```

### Configuration Extension

Add to `internal/config/app.go`:

```go
type AuthValidationConfig struct {
    JWKS JWKSConfig `yaml:"jwks"`
    Issuer string `yaml:"issuer"`
    Audiences []string `yaml:"audiences,omitempty"`
}

type JWKSConfig struct {
    URL string `yaml:"url" validate:"required,url"`
    CacheTTL int `yaml:"cacheTTL"` // seconds, default 3600
    RefreshRetryLimit int `yaml:"refreshRetryLimit"` // default 3
}

func (c *JWKSConfig) setDefaults() {
    if c.CacheTTL == 0 {
        c.CacheTTL = 3600 // 1 hour
    }
    if c.RefreshRetryLimit == 0 {
        c.RefreshRetryLimit = 3
    }
}
```

### Config YAML Updates

Add to `config.yaml`:

```yaml
authValidation:
  jwks:
    url: "http://localhost:3300/.well-known/jwks.json"
    cacheTTL: 3600          # Cache TTL in seconds (1 hour)
    refreshRetryLimit: 3    # Max refresh attempts per minute
  issuer: "http://localhost:3300"
  audiences: []             # Empty = skip audience validation
```

## Files to Create

- `internal/auth/context.go` - Request context with claims
- `internal/auth/jwks.go` - JWKS fetching and caching
- `internal/auth/validator.go` - JWT validation logic
- `internal/auth/authorizer.go` - Authorization helper functions

## Files to Modify

- `internal/config/app.go` - Add AuthValidationConfig struct
- `config.yaml` - Add authValidation settings
- `config.example.yaml` - Add authValidation settings with examples

## Testing Requirements

```go
func TestJWKSCache(t *testing.T) {
    cache := NewJWKSCache(time.Hour, 3)

    // Test initial state
    assert.True(t, cache.IsExpired())

    // Set keys
    keys := map[string]*rsa.PublicKey{"kid1": mockKey}
    cache.SetKeys(keys)

    assert.False(t, cache.IsExpired())
    key, ok := cache.GetKey("kid1")
    assert.True(t, ok)
    assert.NotNil(t, key)
}

func TestJWKSCacheRateLimiting(t *testing.T) {
    cache := NewJWKSCache(time.Hour, 3)

    // First 3 refreshes allowed
    assert.True(t, cache.CanRefresh())
    assert.True(t, cache.CanRefresh())
    assert.True(t, cache.CanRefresh())

    // 4th blocked
    assert.False(t, cache.CanRefresh())
}

func TestAuthorizer(t *testing.T) {
    authorizer := NewAuthorizer()

    // Test superadmin bypass
    ctx := WithAuthContext(context.Background(), &AuthContext{
        UserID:          "user1",
        Permissions:     []string{"root"},
        IsAuthenticated: true,
    })

    err := authorizer.CheckPermission(ctx, "employee:read")
    assert.NoError(t, err) // Superadmin bypasses

    // Test permission check
    ctx = WithAuthContext(context.Background(), &AuthContext{
        UserID:          "user2",
        Permissions:     []string{"employee:read"},
        Memberships:     map[string]string{"proj_abc": "member"},
        IsAuthenticated: true,
    })

    err = authorizer.CheckProjectAccess(ctx, "employee:read", "proj_abc")
    assert.NoError(t, err)

    err = authorizer.CheckProjectAccess(ctx, "employee:write", "proj_abc")
    assert.Error(t, err) // Missing permission

    err = authorizer.CheckProjectAccess(ctx, "employee:read", "proj_xyz")
    assert.Error(t, err) // Not a member
}
```

## Commands to Run

```bash
# Build to verify compilation
make build

# Run tests
go test ./internal/auth/...
```

## Validation Checklist

- [ ] JWKS fetcher retrieves and parses keys correctly
- [ ] JWKS cache TTL works (1 hour default)
- [ ] JWKS refresh rate limiting works (max 3/minute)
- [ ] JWT validator validates signatures correctly
- [ ] JWT validator handles expired tokens
- [ ] JWT validator handles key rotation (force refresh)
- [ ] AuthContext extracts all claims correctly
- [ ] Authorizer.CheckPermission works
- [ ] Authorizer.CheckProjectAccess works
- [ ] Authorizer.IsSuperAdmin bypasses all checks
- [ ] Configuration loaded correctly

## Definition of Done

- [ ] All auth package files created
- [ ] JWKS caching implemented with rate limiting
- [ ] JWT validation implemented
- [ ] AuthContext implemented
- [ ] Authorizer implemented with all helper methods
- [ ] Configuration extended
- [ ] All tests pass
- [ ] Build succeeds

## Dependencies

- T61: JWT claims must include memberships

## Risk Factors

- **Medium Risk**: JWKS fetching requires network access
- **Low Risk**: JWT validation is well-established pattern
- **Low Risk**: Rate limiting is simple counter-based

## Notes

- JWKS endpoint is `/.well-known/jwks.json` on auth server
- Cache TTL of 1 hour balances freshness and performance
- Force refresh on signature failure handles key rotation
- Superadmin check is first in authorization flow for efficiency
- AuthContext is immutable after creation from JWT claims

### Error Messages Reference

| Scenario | Error Code | Message |
|----------|------------|---------|
| No Authorization header | UNAUTHENTICATED | "missing authorization header" |
| Invalid token format | UNAUTHENTICATED | "invalid token format" |
| Expired token | UNAUTHENTICATED | "token has expired" |
| Invalid signature | UNAUTHENTICATED | "invalid token signature" |
| Missing permission | PERMISSION_DENIED | "permission denied: requires {permission}" |
| Not project member | PERMISSION_DENIED | "permission denied: not a member of this project" |
