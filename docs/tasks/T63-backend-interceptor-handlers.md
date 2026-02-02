# Task T63: Backend Interceptor + Handler Integration

**Story Reference:** US15-authorization-rbac.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T62 (Backend Auth Package)

## Objective

Create Connect-RPC auth interceptor that extracts and validates JWT from requests, and integrate authorization checks into all domain handlers.

## Acceptance Criteria

- [x] Connect-RPC interceptor extracts Bearer token from Authorization header
- [x] Interceptor validates JWT and injects AuthContext into request context
- [x] All domain handlers call authorizer.CheckProjectAccess
- [x] Unauthenticated requests return `UNAUTHENTICATED` error
- [x] Unauthorized requests return `PERMISSION_DENIED` error
- [x] Superadmin users bypass all permission checks

## Technical Requirements

### Auth Interceptor Implementation

Create `internal/auth/interceptor.go`:

```go
package auth

import (
    "context"
    "strings"

    "connectrpc.com/connect"
)

// AuthInterceptor creates a Connect-RPC interceptor for JWT validation
func AuthInterceptor(validator *JWTValidator) connect.UnaryInterceptorFunc {
    return func(next connect.UnaryFunc) connect.UnaryFunc {
        return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
            // Extract Authorization header
            authHeader := req.Header().Get("Authorization")

            // If no auth header, continue with unauthenticated context
            // Individual handlers will decide if auth is required
            if authHeader == "" {
                authCtx := &AuthContext{IsAuthenticated: false}
                ctx = WithAuthContext(ctx, authCtx)
                return next(ctx, req)
            }

            // Validate Bearer token format
            if !strings.HasPrefix(authHeader, "Bearer ") {
                return nil, connect.NewError(connect.CodeUnauthenticated,
                    fmt.Errorf("invalid authorization header format"))
            }

            tokenString := strings.TrimPrefix(authHeader, "Bearer ")

            // Validate JWT
            claims, err := validator.Validate(ctx, tokenString)
            if err != nil {
                return nil, connect.NewError(connect.CodeUnauthenticated, err)
            }

            // Create AuthContext from claims
            authCtx := NewAuthContextFromClaims(claims)
            ctx = WithAuthContext(ctx, authCtx)

            return next(ctx, req)
        }
    }
}

// StreamAuthInterceptor creates a streaming interceptor (if needed)
func StreamAuthInterceptor(validator *JWTValidator) connect.StreamingInterceptorFunc {
    return func(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
        return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
            // Extract Authorization header
            authHeader := conn.RequestHeader().Get("Authorization")

            if authHeader == "" {
                authCtx := &AuthContext{IsAuthenticated: false}
                ctx = WithAuthContext(ctx, authCtx)
                return next(ctx, conn)
            }

            if !strings.HasPrefix(authHeader, "Bearer ") {
                return connect.NewError(connect.CodeUnauthenticated,
                    fmt.Errorf("invalid authorization header format"))
            }

            tokenString := strings.TrimPrefix(authHeader, "Bearer ")

            claims, err := validator.Validate(ctx, tokenString)
            if err != nil {
                return connect.NewError(connect.CodeUnauthenticated, err)
            }

            authCtx := NewAuthContextFromClaims(claims)
            ctx = WithAuthContext(ctx, authCtx)

            return next(ctx, conn)
        }
    }
}
```

### Server Interceptor Registration

Modify `internal/server/server.go`:

```go
import (
    "github.com/your-project/internal/auth"
)

func (s *Server) setupRoutes() {
    // Create JWT validator
    validator := auth.NewJWTValidator(
        s.cfg.AuthValidation.JWKS.URL,
        s.cfg.AuthValidation.Issuer,
        s.cfg.AuthValidation.Audiences,
        s.cfg.AuthValidation.JWKS.CacheTTL,
        s.cfg.AuthValidation.JWKS.RefreshRetryLimit,
    )

    // Create auth interceptor
    authInterceptor := auth.AuthInterceptor(validator)

    // Apply to all Connect handlers
    interceptors := connect.WithInterceptors(authInterceptor)

    // Register handlers with interceptor
    employeeHandler := employee.NewHandler(s.container.EmployeeService, s.container.Authorizer)
    path, handler := employeev1connect.NewEmployeeServiceHandler(employeeHandler, interceptors)
    s.mux.Handle(path, handler)

    // ... repeat for all handlers
}
```

### Handler Authorization Pattern

Update each domain handler to check authorization. Example for employee:

**Modify `internal/domain/employee/handler.go`:**

```go
package employee

import (
    "context"

    "connectrpc.com/connect"
    "github.com/your-project/internal/auth"
    employeev1 "github.com/your-project/gen/altalune/employee/v1"
)

type Handler struct {
    service    *Service
    authorizer *auth.Authorizer
}

func NewHandler(service *Service, authorizer *auth.Authorizer) *Handler {
    return &Handler{
        service:    service,
        authorizer: authorizer,
    }
}

// ListEmployees lists employees with authorization check
func (h *Handler) ListEmployees(ctx context.Context, req *connect.Request[employeev1.ListEmployeesRequest]) (*connect.Response[employeev1.ListEmployeesResponse], error) {
    // Authorization check
    if err := h.authorizer.CheckProjectAccess(ctx, "employee:read", req.Msg.ProjectId); err != nil {
        return nil, err
    }

    // Proceed with business logic
    result, err := h.service.List(ctx, req.Msg)
    if err != nil {
        return nil, err
    }

    return connect.NewResponse(result), nil
}

// GetEmployee gets a single employee with authorization check
func (h *Handler) GetEmployee(ctx context.Context, req *connect.Request[employeev1.GetEmployeeRequest]) (*connect.Response[employeev1.GetEmployeeResponse], error) {
    // Authorization check
    if err := h.authorizer.CheckProjectAccess(ctx, "employee:read", req.Msg.ProjectId); err != nil {
        return nil, err
    }

    result, err := h.service.Get(ctx, req.Msg)
    if err != nil {
        return nil, err
    }

    return connect.NewResponse(result), nil
}

// CreateEmployee creates an employee with authorization check
func (h *Handler) CreateEmployee(ctx context.Context, req *connect.Request[employeev1.CreateEmployeeRequest]) (*connect.Response[employeev1.CreateEmployeeResponse], error) {
    // Authorization check - requires write permission
    if err := h.authorizer.CheckProjectAccess(ctx, "employee:write", req.Msg.ProjectId); err != nil {
        return nil, err
    }

    result, err := h.service.Create(ctx, req.Msg)
    if err != nil {
        return nil, err
    }

    return connect.NewResponse(result), nil
}

// UpdateEmployee updates an employee with authorization check
func (h *Handler) UpdateEmployee(ctx context.Context, req *connect.Request[employeev1.UpdateEmployeeRequest]) (*connect.Response[employeev1.UpdateEmployeeResponse], error) {
    // Authorization check - requires write permission
    if err := h.authorizer.CheckProjectAccess(ctx, "employee:write", req.Msg.ProjectId); err != nil {
        return nil, err
    }

    result, err := h.service.Update(ctx, req.Msg)
    if err != nil {
        return nil, err
    }

    return connect.NewResponse(result), nil
}

// DeleteEmployee deletes an employee with authorization check
func (h *Handler) DeleteEmployee(ctx context.Context, req *connect.Request[employeev1.DeleteEmployeeRequest]) (*connect.Response[employeev1.DeleteEmployeeResponse], error) {
    // Authorization check - requires delete permission
    if err := h.authorizer.CheckProjectAccess(ctx, "employee:delete", req.Msg.ProjectId); err != nil {
        return nil, err
    }

    result, err := h.service.Delete(ctx, req.Msg)
    if err != nil {
        return nil, err
    }

    return connect.NewResponse(result), nil
}
```

### Handler Authorization Matrix

Apply authorization checks to all handlers following this pattern:

| Domain | List | Get | Create | Update | Delete |
|--------|------|-----|--------|--------|--------|
| Employee | employee:read | employee:read | employee:write | employee:write | employee:delete |
| User | user:read | user:read | user:write | user:write | user:delete |
| Role | role:read | role:read | role:write | role:write | role:delete |
| Permission | permission:read | permission:read | permission:write | permission:write | permission:delete |
| Project | project:read | project:read | project:write | project:write | project:delete |
| API Key | apikey:read | apikey:read | apikey:write | apikey:write | apikey:delete |
| Chatbot | chatbot:read | chatbot:read | chatbot:write | chatbot:write | chatbot:delete |
| OAuth Client | client:read | client:read | client:write | client:write | client:delete |
| Member | member:read | member:read | member:write | member:write | member:delete |
| IAM Mapper | iam:read | - | iam:write | iam:write | iam:write |

### Container Wiring

Update `internal/container/container.go`:

```go
import "github.com/your-project/internal/auth"

type Container struct {
    // ... existing fields
    Authorizer *auth.Authorizer
    JWTValidator *auth.JWTValidator
}

func NewContainer(cfg *config.Config, db *pgxpool.Pool) *Container {
    // ... existing setup

    // Create authorizer
    authorizer := auth.NewAuthorizer()

    // Create JWT validator
    jwtValidator := auth.NewJWTValidator(
        cfg.AuthValidation.JWKS.URL,
        cfg.AuthValidation.Issuer,
        cfg.AuthValidation.Audiences,
        cfg.AuthValidation.JWKS.CacheTTL,
        cfg.AuthValidation.JWKS.RefreshRetryLimit,
    )

    return &Container{
        // ... existing fields
        Authorizer:   authorizer,
        JWTValidator: jwtValidator,
    }
}
```

### Project-Scoped vs Global Handlers

Some handlers may not be project-scoped (e.g., listing projects):

```go
// ListProjects - user can see all projects they're members of
func (h *Handler) ListProjects(ctx context.Context, req *connect.Request[projectv1.ListProjectsRequest]) (*connect.Response[projectv1.ListProjectsResponse], error) {
    // Just check if authenticated
    if err := h.authorizer.CheckAuthenticated(ctx); err != nil {
        return nil, err
    }

    // Get user's project memberships from context
    authCtx := auth.FromContext(ctx)

    // If superadmin, return all projects
    if h.authorizer.IsSuperAdmin(ctx) {
        result, err := h.service.ListAll(ctx, req.Msg)
        if err != nil {
            return nil, err
        }
        return connect.NewResponse(result), nil
    }

    // Otherwise, filter to user's projects
    projectIDs := h.authorizer.GetUserProjects(ctx)
    result, err := h.service.ListByIDs(ctx, projectIDs, req.Msg)
    if err != nil {
        return nil, err
    }

    return connect.NewResponse(result), nil
}
```

## Files to Create

- `internal/auth/interceptor.go` - Connect-RPC auth interceptor

## Files to Modify

- `internal/server/server.go` - Register interceptor
- `internal/container/container.go` - Wire Authorizer and JWTValidator
- `internal/domain/employee/handler.go` - Add authorization checks
- `internal/domain/user/handler.go` - Add authorization checks
- `internal/domain/role/handler.go` - Add authorization checks
- `internal/domain/permission/handler.go` - Add authorization checks
- `internal/domain/project/handler.go` - Add authorization checks
- `internal/domain/api_key/handler.go` - Add authorization checks
- `internal/domain/chatbot_config/handler.go` - Add authorization checks
- `internal/domain/oauth_client/handler.go` - Add authorization checks
- `internal/domain/project_member/handler.go` - Add authorization checks
- `internal/domain/iam_mapper/handler.go` - Add authorization checks

## Testing Requirements

```go
func TestAuthInterceptor(t *testing.T) {
    // Setup mock validator
    validator := NewMockJWTValidator()
    interceptor := AuthInterceptor(validator)

    // Test missing auth header
    req := connect.NewRequest(&testpb.Request{})
    _, err := interceptor(mockNext)(context.Background(), req)
    // Should succeed with unauthenticated context
    assert.NoError(t, err)

    // Test valid token
    req.Header().Set("Authorization", "Bearer valid-token")
    validator.SetValidToken("valid-token", &AccessTokenClaims{
        RegisteredClaims: jwt.RegisteredClaims{Subject: "user1"},
        Perms:            []string{"employee:read"},
    })

    _, err = interceptor(mockNext)(context.Background(), req)
    assert.NoError(t, err)

    // Test invalid token
    req.Header().Set("Authorization", "Bearer invalid-token")
    _, err = interceptor(mockNext)(context.Background(), req)
    assert.Error(t, err)
    assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
}

func TestHandlerAuthorization(t *testing.T) {
    handler := NewHandler(mockService, auth.NewAuthorizer())

    // Test unauthorized access
    ctx := auth.WithAuthContext(context.Background(), &auth.AuthContext{
        UserID:          "user1",
        Permissions:     []string{}, // No permissions
        IsAuthenticated: true,
    })

    req := connect.NewRequest(&employeev1.ListEmployeesRequest{ProjectId: "proj_123"})
    _, err := handler.ListEmployees(ctx, req)

    assert.Error(t, err)
    assert.Equal(t, connect.CodePermissionDenied, connect.CodeOf(err))

    // Test authorized access
    ctx = auth.WithAuthContext(context.Background(), &auth.AuthContext{
        UserID:          "user2",
        Permissions:     []string{"employee:read"},
        Memberships:     map[string]string{"proj_123": "member"},
        IsAuthenticated: true,
    })

    _, err = handler.ListEmployees(ctx, req)
    assert.NoError(t, err)
}
```

## Commands to Run

```bash
# Build to verify compilation
make build

# Run tests
go test ./internal/auth/...
go test ./internal/domain/employee/...
go test ./internal/domain/user/...
# ... etc for all domains

# Start server and test manually
air
# Then test with curl:
curl -H "Authorization: Bearer <token>" http://localhost:3000/api/employees
```

## Validation Checklist

- [ ] Interceptor extracts Bearer token correctly
- [ ] Invalid token format returns UNAUTHENTICATED
- [ ] Expired token returns UNAUTHENTICATED
- [ ] Valid token injects AuthContext
- [ ] Employee handlers check employee:read/write/delete
- [ ] User handlers check user:read/write/delete
- [ ] Role handlers check role:read/write/delete
- [ ] Permission handlers check permission:read/write/delete
- [ ] Project handlers check project:read/write/delete
- [ ] API Key handlers check apikey:read/write/delete
- [ ] Chatbot handlers check chatbot:read/write/delete
- [ ] OAuth Client handlers check client:read/write/delete
- [ ] Member handlers check member:read/write/delete
- [ ] IAM Mapper handlers check iam:read/write
- [ ] Superadmin bypasses all checks
- [ ] Project membership validated for project-scoped resources

## Definition of Done

- [x] Auth interceptor created and registered
- [x] All handlers updated with authorization checks
- [x] Container wiring complete
- [ ] All tests pass
- [ ] Manual testing confirms authorization works
- [x] Build succeeds

## Dependencies

- T62: Auth package must be complete

## Risk Factors

- **High Risk**: Breaking existing functionality if auth checks added incorrectly
- **Medium Risk**: Performance impact from additional context checks
- **Low Risk**: Pattern is consistent across all handlers

## Notes

- Interceptor runs on EVERY request but only validates if token present
- Handlers decide if authentication is required
- Always use `CheckProjectAccess` for project-scoped resources
- Use `CheckPermission` for global resources (like listing all projects)
- Use `CheckAuthenticated` for endpoints that just need a logged-in user
- Superadmin check is first in authorization flow for efficiency
- Error messages should not leak internal details

### Manual Testing Scenarios

1. **No token**: `curl http://localhost:3000/api/employees` → depends on handler
2. **Valid token, has permission**: `curl -H "Authorization: Bearer <token>" ...` → 200 OK
3. **Valid token, no permission**: → 403 PERMISSION_DENIED
4. **Valid token, not project member**: → 403 PERMISSION_DENIED
5. **Superadmin token**: → 200 OK (always)
6. **Expired token**: → 401 UNAUTHENTICATED
7. **Invalid signature**: → 401 UNAUTHENTICATED
