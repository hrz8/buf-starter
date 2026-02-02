# Auth Package

The `internal/auth/` package provides JWT validation and authorization for Connect-RPC endpoints.

## Package Structure

```
internal/auth/
├── interceptor.go  # Connect-RPC interceptor
├── validator.go    # JWKS-based JWT validation
├── jwks.go         # JWKS fetcher with caching
├── context.go      # AuthContext struct and helpers
└── authorizer.go   # Authorization check methods
```

## AuthContext

Holds authenticated user information extracted from JWT:

```go
type AuthContext struct {
    UserID          string            // JWT subject (user public_id)
    Email           string
    Name            string
    Permissions     []string          // e.g., ["employee:read", "project:write"]
    Memberships     map[string]string // project_public_id -> role
    EmailVerified   bool
    IsAuthenticated bool
}
```

### Context Helpers

```go
// Extract from context
auth := auth.FromContext(ctx)

// Check if authenticated
if !auth.IsAuthenticated {
    return nil, connect.NewError(connect.CodeUnauthenticated, ...)
}

// Access user info
userID := auth.UserID
perms := auth.Permissions
role := auth.Memberships["proj_abc123"]
```

## Authorizer

Provides authorization check methods:

```go
type Authorizer struct{}

func NewAuthorizer() *Authorizer

// Check methods return nil on success, connect error on failure
func (a *Authorizer) CheckAuthenticated(ctx context.Context) error
func (a *Authorizer) CheckPermission(ctx context.Context, permission string) error
func (a *Authorizer) CheckProjectAccess(ctx context.Context, permission string, projectID string) error
func (a *Authorizer) CheckProjectMembership(ctx context.Context, projectID string) error

// Query methods
func (a *Authorizer) IsSuperAdmin(ctx context.Context) bool
func (a *Authorizer) HasPermission(ctx context.Context, permission string) bool
func (a *Authorizer) HasAnyPermission(ctx context.Context, permissions []string) bool
func (a *Authorizer) GetUserProjects(ctx context.Context) []string
func (a *Authorizer) GetProjectRole(ctx context.Context, projectID string) string
```

### Superadmin Bypass

Users with `root` permission bypass all authorization checks:

```go
const RootPermission = "root"

func (a *Authorizer) IsSuperAdmin(ctx context.Context) bool {
    auth := FromContext(ctx)
    for _, p := range auth.Permissions {
        if p == RootPermission {
            return true
        }
    }
    return false
}
```

## Connect-RPC Interceptor

The interceptor extracts and validates JWT tokens:

```go
func NewAuthInterceptor(validator *JWTValidator) connect.Interceptor
```

Token extraction priority:
1. `Authorization: Bearer <token>` header
2. `access_token` cookie (for httpOnly cookie auth)

If no token found, continues with unauthenticated context (handlers decide if auth required).

### Registration

```go
// internal/server/http_routes.go
var handlerOptions []connect.HandlerOption
if validator := s.c.GetJWTValidator(); validator != nil {
    authInterceptor := auth.NewAuthInterceptor(validator)
    handlerOptions = append(handlerOptions, connect.WithInterceptors(authInterceptor))
}

// Register handler with interceptor
employeeHandler := employee_domain.NewHandler(s.c.GetEmployeeService(), authorizer)
employeePath, employeeConnectHandler := altalunev1connect.NewEmployeeServiceHandler(employeeHandler, handlerOptions...)
```

## JWT Validator

JWKS-based JWT validation:

```go
func NewJWTValidator(jwksURL, issuer string, audiences []string, cacheTTL int, refreshLimit int) *JWTValidator

func (v *JWTValidator) Validate(ctx context.Context, tokenString string) (*AccessTokenClaims, error)
```

Features:
- JWKS caching with configurable TTL
- Automatic key refresh on signature failure (key rotation support)
- Rate-limited refresh (prevents abuse)
- Issuer and audience validation

### AccessTokenClaims

```go
type AccessTokenClaims struct {
    jwt.RegisteredClaims
    Scope         string            `json:"scope,omitempty"`
    Email         string            `json:"email,omitempty"`
    Name          string            `json:"name,omitempty"`
    Perms         []string          `json:"perms"`
    Memberships   map[string]string `json:"memberships,omitempty"`
    EmailVerified bool              `json:"email_verified"`
}
```

## Container Wiring

```go
// internal/container/container.go
type Container struct {
    jwtValidator *auth.JWTValidator
    authorizer   *auth.Authorizer
}

func (c *Container) initAuthComponents() {
    c.authorizer = auth.NewAuthorizer()

    if c.config.GetAuthValidationJWKSURL() != "" {
        c.jwtValidator = auth.NewJWTValidator(
            c.config.GetAuthValidationJWKSURL(),
            c.config.GetAuthValidationIssuer(),
            c.config.GetAuthValidationAudiences(),
            c.config.GetAuthValidationJWKSCacheTTL(),
            c.config.GetAuthValidationJWKSRefreshLimit(),
        )
    }
}

// internal/container/getter.go
func (c *Container) GetJWTValidator() *auth.JWTValidator { return c.jwtValidator }
func (c *Container) GetAuthorizer() *auth.Authorizer { return c.authorizer }
```

## Error Codes

| Scenario | Connect Code |
|----------|-------------|
| Missing/invalid token | `CodeUnauthenticated` |
| Expired token | `CodeUnauthenticated` |
| Missing permission | `CodePermissionDenied` |
| Not project member | `CodePermissionDenied` |
