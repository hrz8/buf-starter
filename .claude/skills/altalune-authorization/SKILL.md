---
name: altalune-authorization
description: |
  Authorization implementation for Altalune including project membership and RBAC. Use when: (1) Adding authorization checks to new handlers, (2) Working with project membership validation, (3) Understanding JWT token structure, (4) Adding new permissions to the system. Covers: Connect-RPC auth interceptor, Authorizer methods (CheckPermission, CheckProjectAccess), JWT claims with permissions and memberships, handler authorization patterns.
---

# Altalune Authorization

## Overview

Altalune uses TWO authorization layers validated from JWT claims:

1. **RBAC Permissions** - Fine-grained permissions (`employee:read`, `project:write`)
2. **Project Membership** - Project-scoped access (`memberships` map in JWT)

```
Request → Auth Interceptor → JWT Validation → AuthContext injected
                                  ↓
                          Handler checks:
                          - CheckPermission (global resources)
                          - CheckProjectAccess (project-scoped resources)
```

## Architecture

```
internal/auth/
├── interceptor.go  # Connect-RPC interceptor (JWT extraction + validation)
├── validator.go    # JWKS-based JWT validation
├── jwks.go         # JWKS fetcher with caching
├── context.go      # AuthContext struct and context helpers
└── authorizer.go   # Authorization check methods
```

## Quick Reference

### Authorizer Methods

| Method | Use Case | Checks |
|--------|----------|--------|
| `CheckAuthenticated(ctx)` | Public but logged-in routes | Is authenticated? |
| `CheckPermission(ctx, perm)` | Global resources (users, roles) | Auth + Permission |
| `CheckProjectAccess(ctx, perm, projectID)` | Project-scoped resources | Auth + Permission + Membership |

### Handler Patterns

**Global resource (no project_id):**
```go
func (h *Handler) QueryUsers(ctx context.Context, req *connect.Request[...]) (..., error) {
    if err := h.auth.CheckPermission(ctx, "user:read"); err != nil {
        return nil, err
    }
    // ... business logic
}
```

**Project-scoped resource:**
```go
func (h *Handler) QueryEmployees(ctx context.Context, req *connect.Request[...]) (..., error) {
    if err := h.auth.CheckProjectAccess(ctx, "employee:read", req.Msg.ProjectId); err != nil {
        return nil, err
    }
    // ... business logic
}
```

### Permission Matrix

| Domain | Read | Write | Delete |
|--------|------|-------|--------|
| Employee | `employee:read` | `employee:write` | `employee:delete` |
| Project | `project:read` | `project:write` | `project:delete` |
| User | `user:read` | `user:write` | `user:delete` |
| Role | `role:read` | `role:write` | `role:delete` |
| Permission | `permission:read` | `permission:write` | `permission:delete` |
| OAuth Client | `client:read` | `client:write` | `client:delete` |
| API Key | `apikey:read` | `apikey:write` | `apikey:delete` |
| Chatbot | `chatbot:read` | `chatbot:write` | `chatbot:delete` |
| IAM Mapper | `iam:read` | `iam:write` | - |
| Project Member | `member:read` | `member:write` | - |

### JWT Token Structure

```json
{
  "sub": "user_public_id",
  "aud": ["client_id"],
  "iss": "http://localhost:3000",
  "email": "user@example.com",
  "name": "User Name",
  "perms": ["employee:read", "employee:write", "project:read"],
  "memberships": {
    "proj_abc123": "admin",
    "proj_def456": "member"
  },
  "email_verified": true
}
```

## Reference Files

- **[auth-package.md](references/auth-package.md)** - Full auth package implementation
- **[handler-patterns.md](references/handler-patterns.md)** - Handler authorization patterns
- **[permission-list.md](references/permission-list.md)** - All predefined permissions

## Key Files

**Auth Package:**
- `internal/auth/interceptor.go` - Connect-RPC interceptor
- `internal/auth/validator.go` - JWKS JWT validation
- `internal/auth/context.go` - AuthContext struct
- `internal/auth/authorizer.go` - Authorization methods

**Handler Integration:**
- `internal/server/http_routes.go` - Handler registration with authorizer
- `internal/domain/*/handler.go` - Authorization checks in each handler

**Container Wiring:**
- `internal/container/container.go` - JWTValidator and Authorizer initialization
- `internal/container/getter.go` - GetJWTValidator() and GetAuthorizer()

**Configuration:**
- `config.yaml` - `auth_validation` section for JWKS settings

## Configuration

```yaml
# config.yaml
auth_validation:
  jwks:
    url: "http://localhost:3000/oauth/.well-known/jwks.json"
    cache_ttl: 3600        # seconds
    refresh_retry_limit: 3 # per minute
  issuer: "http://localhost:3000"
  audiences: []            # optional, validates aud claim
```
