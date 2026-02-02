# Authorization Middleware Patterns

## Overview

Connect-RPC uses interceptors (similar to middleware) for authorization. The pattern:

```
Request → JWT Validation → Extract Claims → Check Permissions → Handler
```

## Implementation Approach

### 1. Create Auth Context

**File:** `internal/shared/authctx/context.go`

```go
package authctx

import (
    "context"

    "github.com/hrz8/altalune/internal/shared/jwt"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

// WithClaims adds JWT claims to context
func WithClaims(ctx context.Context, claims *jwt.AccessTokenClaims) context.Context {
    return context.WithValue(ctx, claimsKey, claims)
}

// GetClaims retrieves JWT claims from context
func GetClaims(ctx context.Context) (*jwt.AccessTokenClaims, bool) {
    claims, ok := ctx.Value(claimsKey).(*jwt.AccessTokenClaims)
    return claims, ok
}

// GetUserID returns user public_id from context
func GetUserID(ctx context.Context) (string, bool) {
    claims, ok := GetClaims(ctx)
    if !ok {
        return "", false
    }
    return claims.Subject, true
}

// GetPermissions returns user permissions from context
func GetPermissions(ctx context.Context) []string {
    claims, ok := GetClaims(ctx)
    if !ok {
        return nil
    }
    return claims.Perms
}

// HasPermission checks if user has specific permission
func HasPermission(ctx context.Context, permission string) bool {
    perms := GetPermissions(ctx)
    for _, p := range perms {
        // Check exact match or wildcard "root"
        if p == permission || p == "root" {
            return true
        }
    }
    return false
}

// HasAnyPermission checks if user has any of the given permissions
func HasAnyPermission(ctx context.Context, permissions ...string) bool {
    for _, p := range permissions {
        if HasPermission(ctx, p) {
            return true
        }
    }
    return false
}
```

### 2. Create Connect-RPC Auth Interceptor

**File:** `internal/server/interceptor.go`

```go
package server

import (
    "context"
    "strings"

    "connectrpc.com/connect"
    "github.com/hrz8/altalune/internal/shared/authctx"
    "github.com/hrz8/altalune/internal/shared/jwt"
)

// AuthInterceptor validates JWT and adds claims to context
func AuthInterceptor(signer *jwt.Signer) connect.UnaryInterceptorFunc {
    return func(next connect.UnaryFunc) connect.UnaryFunc {
        return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
            // Extract token from Authorization header
            authHeader := req.Header().Get("Authorization")
            if authHeader == "" {
                return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
            }

            // Parse "Bearer <token>"
            parts := strings.SplitN(authHeader, " ", 2)
            if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
                return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid authorization header format"))
            }

            token := parts[1]

            // Validate token
            claims, err := signer.ValidateAccessToken(token)
            if err != nil {
                return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid token: %w", err))
            }

            // Add claims to context
            ctx = authctx.WithClaims(ctx, claims)

            return next(ctx, req)
        }
    }
}

// PermissionInterceptor checks if user has required permission
func PermissionInterceptor(requiredPerm string) connect.UnaryInterceptorFunc {
    return func(next connect.UnaryFunc) connect.UnaryFunc {
        return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
            if !authctx.HasPermission(ctx, requiredPerm) {
                return nil, connect.NewError(connect.CodePermissionDenied,
                    fmt.Errorf("permission denied: requires %s", requiredPerm))
            }
            return next(ctx, req)
        }
    }
}
```

### 3. Apply Interceptors to Handlers

**File:** `internal/server/http_routes.go`

```go
import "connectrpc.com/connect"

func (s *Server) setupRoutes() *http.ServeMux {
    // Create interceptor chain
    authInterceptor := AuthInterceptor(s.jwtSigner)

    // Interceptor options for authenticated routes
    authOpts := connect.WithInterceptors(authInterceptor)

    // Register handlers with interceptors
    employeePath, employeeHandler := altalunev1connect.NewEmployeeServiceHandler(
        employeeHandler,
        authOpts,  // Require authentication
    )

    // For permission-specific routes, chain interceptors
    adminOpts := connect.WithInterceptors(
        authInterceptor,
        PermissionInterceptor("admin:access"),
    )

    // Public routes (no interceptors)
    configPath, configHandler := altalunev1connect.NewConfigServiceHandler(configHandler)

    // ...
}
```

## Project Membership Validation

**Note:** Use database query approach for membership validation (not JWT claims). See [project-membership.md](project-membership.md#validation-approach-database-query-vs-jwt-claims) for rationale.

### Check Project Access in Service Layer

```go
// internal/domain/employee/service.go

func (s *Service) GetEmployee(ctx context.Context, req *connect.Request[altalunev1.GetEmployeeRequest]) (*connect.Response[altalunev1.GetEmployeeResponse], error) {
    // 1. Get user from context
    userID, ok := authctx.GetUserID(ctx)
    if !ok {
        return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user not authenticated"))
    }

    // 2. Get project ID from request
    projectID := req.Msg.GetProjectId()

    // 3. Check project membership (DB query, optionally cached)
    member, err := s.projectMemberRepo.GetByProjectAndUser(ctx, projectID, userID)
    if err != nil {
        return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("not a member of this project"))
    }

    // 4. Optional: Check minimum role
    if !hasMinimumRole(member.Role, "member") {
        return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("insufficient role"))
    }

    // 5. Proceed with operation
    // ...
}

func hasMinimumRole(userRole, minRole string) bool {
    roleHierarchy := map[string]int{
        "user":   1,
        "member": 2,
        "admin":  3,
        "owner":  4,
    }
    return roleHierarchy[userRole] >= roleHierarchy[minRole]
}
```

### Project Membership Interceptor (Alternative)

```go
// ProjectMemberInterceptor checks project membership
// projectIDExtractor is a function that extracts project_id from the request
func ProjectMemberInterceptor(
    memberRepo ProjectMemberRepository,
    minRole string,
) connect.UnaryInterceptorFunc {
    return func(next connect.UnaryFunc) connect.UnaryFunc {
        return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
            userID, ok := authctx.GetUserID(ctx)
            if !ok {
                return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("not authenticated"))
            }

            // Extract project_id from request (implementation depends on proto)
            // This is tricky because request types vary
            projectID := extractProjectID(req)
            if projectID == "" {
                return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("project_id required"))
            }

            member, err := memberRepo.GetByProjectAndUser(ctx, projectID, userID)
            if err != nil {
                return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("not a project member"))
            }

            if !hasMinimumRole(member.Role, minRole) {
                return nil, connect.NewError(connect.CodePermissionDenied,
                    fmt.Errorf("requires %s role, have %s", minRole, member.Role))
            }

            return next(ctx, req)
        }
    }
}
```

## Permission Checking Patterns

### Pattern 1: Interceptor-Based (Route Level)

```go
// Apply to entire service
projectOpts := connect.WithInterceptors(
    AuthInterceptor(signer),
    PermissionInterceptor("project:read"),
)
projectPath, projectHandler := altalunev1connect.NewProjectServiceHandler(handler, projectOpts)
```

### Pattern 2: Service-Level (Method Level)

```go
func (s *Service) DeleteProject(ctx context.Context, req *connect.Request[...]) (*connect.Response[...], error) {
    // Check specific permission in service method
    if !authctx.HasPermission(ctx, "project:delete") {
        return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("requires project:delete permission"))
    }
    // ...
}
```

### Pattern 3: Permission Matrix

```go
// Define permission requirements per operation
var permissionMatrix = map[string]string{
    "ProjectService/GetProject":    "project:read",
    "ProjectService/CreateProject": "project:create",
    "ProjectService/UpdateProject": "project:update",
    "ProjectService/DeleteProject": "project:delete",
    "EmployeeService/GetEmployee":  "employee:read",
    // ...
}

// Dynamic permission interceptor
func DynamicPermissionInterceptor() connect.UnaryInterceptorFunc {
    return func(next connect.UnaryFunc) connect.UnaryFunc {
        return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
            procedure := req.Spec().Procedure  // e.g., "/altalune.v1.ProjectService/GetProject"

            // Extract method name
            parts := strings.Split(procedure, "/")
            methodKey := parts[len(parts)-2] + "/" + parts[len(parts)-1]

            requiredPerm, exists := permissionMatrix[methodKey]
            if exists && !authctx.HasPermission(ctx, requiredPerm) {
                return nil, connect.NewError(connect.CodePermissionDenied,
                    fmt.Errorf("requires %s permission", requiredPerm))
            }

            return next(ctx, req)
        }
    }
}
```

## Best Practices

### 1. Defense in Depth

```go
// Layer 1: Auth interceptor (validates JWT)
// Layer 2: Permission interceptor (checks permissions)
// Layer 3: Service-level checks (project membership, ownership)

opts := connect.WithInterceptors(
    AuthInterceptor(signer),           // Must be authenticated
    PermissionInterceptor("api:access"), // Must have API access
)
```

### 2. Public vs Protected Routes

```go
// Public routes - no interceptors
configPath, configHandler := altalunev1connect.NewConfigServiceHandler(handler)

// Protected routes - with auth
projectPath, projectHandler := altalunev1connect.NewProjectServiceHandler(handler, authOpts)

// Admin routes - with auth + admin permission
adminPath, adminHandler := altalunev1connect.NewAdminServiceHandler(handler, adminOpts)
```

### 3. Error Messages

```go
// Don't leak information
// BAD: "user admin@example.com not found"
// GOOD: "authentication failed"

// Be specific about permission requirements
// BAD: "access denied"
// GOOD: "requires project:delete permission"
```

### 4. Audit Logging

```go
func AuditInterceptor(logger Logger) connect.UnaryInterceptorFunc {
    return func(next connect.UnaryFunc) connect.UnaryFunc {
        return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
            userID, _ := authctx.GetUserID(ctx)
            procedure := req.Spec().Procedure

            logger.Info("API access",
                "user_id", userID,
                "procedure", procedure,
                "timestamp", time.Now(),
            )

            return next(ctx, req)
        }
    }
}
```

## Testing

### Unit Test Authorization

```go
func TestAuthInterceptor(t *testing.T) {
    // Create test signer
    signer := createTestSigner(t)

    // Generate valid token
    token, _ := signer.GenerateAccessToken(jwt.GenerateTokenParams{
        UserPublicID: "test-user-id",
        Perms:        []string{"project:read"},
        Expiry:       time.Hour,
    })

    // Create request with token
    req := connect.NewRequest(&pb.GetProjectRequest{})
    req.Header().Set("Authorization", "Bearer "+token)

    // Test interceptor
    interceptor := AuthInterceptor(signer)
    // ...
}
```

### Integration Test

```bash
# Test with valid token
curl -X POST http://localhost:8080/api/altalune.v1.ProjectService/GetProject \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id": "project-id"}'

# Test without token (should fail)
curl -X POST http://localhost:8080/api/altalune.v1.ProjectService/GetProject \
  -H "Content-Type: application/json" \
  -d '{"id": "project-id"}'
# Expected: 401 Unauthenticated
```
