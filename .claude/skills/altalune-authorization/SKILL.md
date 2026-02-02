---
name: altalune-authorization
description: |
  Authorization implementation for Altalune including project membership and RBAC. Use when: (1) Implementing authorization middleware for Connect-RPC endpoints, (2) Adding permission checks to services, (3) Working with project membership validation, (4) Defining predefined permissions, (5) Implementing role-based access control. Covers: Project membership (owner/admin/member/user roles), RBAC (roles, permissions, assignments), JWT token permissions, middleware patterns.
---

# Altalune Authorization

## Overview

Altalune has TWO authorization layers:

1. **Project Membership** - Coarse-grained access (owner/admin/member/user)
2. **RBAC (Role-Based Access Control)** - Fine-grained permissions

```
User authenticates → JWT contains permissions
                   ↓
Request hits endpoint → Middleware checks:
  1. Is user member of project? (project_id from request)
  2. Does user have required permission? (from JWT perms)
```

## Current State

| Component | Status |
|-----------|--------|
| Permission/Role CRUD | ✅ Implemented |
| IAM Mapper (assignments) | ✅ Implemented |
| JWT includes `perms` | ✅ Implemented |
| Project membership table | ✅ Exists |
| Auth middleware | ❌ Not enforced yet |
| Predefined permissions | ❌ Only `root` exists |

## Quick Reference

### Project Membership Roles

| Role | Scope | Dashboard | Manage Members |
|------|-------|-----------|----------------|
| `owner` | Superadmin only (user_id=1) | ✅ | ✅ All |
| `admin` | Project-level admin | ✅ | ✅ (member/user only) |
| `member` | Regular project access | ✅ | ❌ |
| `user` | OAuth user (default) | ❌ | ❌ |

### JWT Token Structure

```json
{
  "sub": "user_public_id",
  "aud": ["client_id"],
  "scope": "openid profile email",
  "email": "user@example.com",
  "perms": ["project:read", "employee:write"],
  "email_verified": true
}
```

## Reference Files

- **[middleware-patterns.md](references/middleware-patterns.md)** - Authorization middleware implementation
- **[permission-definitions.md](references/permission-definitions.md)** - Predefined permissions list
- **[project-membership.md](references/project-membership.md)** - Project membership rules and patterns

## Key Files

**Backend:**
- `internal/shared/jwt/claims.go` - JWT claims with `Perms []string`
- `internal/domain/permission/` - Permission CRUD domain
- `internal/domain/role/` - Role CRUD domain
- `internal/domain/iam_mapper/` - Assignment operations
- `internal/domain/oauth_auth/permission_service.go` - Permission fetcher
- `internal/server/middleware.go` - Current middleware (needs auth)

**Database:**
- `altalune_permissions` - Permission definitions
- `altalune_roles` - Role definitions
- `altalune_users_roles` - User-role assignments
- `altalune_roles_permissions` - Role-permission assignments
- `altalune_users_permissions` - Direct user-permission assignments
- `altalune_project_members` - Project membership (project_id, user_id, role)
