# Project Membership

## Overview

Project membership is the coarse-grained access control layer. Users access projects through membership, with one of 4 roles defining their capabilities.

## Role Hierarchy

```
owner (4) > admin (3) > member (2) > user (1)
```

| Role | Scope | Dashboard | Manage Settings | Manage Members |
|------|-------|-----------|-----------------|----------------|
| `owner` | Superadmin only | ✅ | ✅ All | ✅ All roles |
| `admin` | Project-level | ✅ | ✅ (no delete) | ✅ (member/user) |
| `member` | Project-level | ✅ | ❌ View only | ❌ |
| `user` | OAuth default | ❌ | ❌ | ❌ |

## Key Rules

### Rule 1: Owner is Reserved

- Only user_id=1 (superadmin) can have `owner` role
- Automatically registered to ALL projects on project creation
- Cannot be assigned to other users

### Rule 2: Auto-Registration on OAuth

When user authenticates via OAuth client for first time:
1. User created in `altalune_users`
2. Identity created in `altalune_user_identities`
3. Membership created with `role='user'` in OAuth client's project

### Rule 3: Multi-Project Access

Users (admin/member) can belong to multiple projects with different roles:

```
User A
├─ Project 1: admin
├─ Project 2: member
└─ Project 3: member
```

### Rule 4: Cascade on Delete

- Project deleted → All memberships deleted
- User deleted → All memberships deleted

## Database Schema

```sql
CREATE TABLE altalune_project_members (
  id BIGSERIAL PRIMARY KEY,
  public_id VARCHAR(14) NOT NULL UNIQUE,
  project_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  role VARCHAR(20) NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'user')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(project_id, user_id)  -- One membership per user per project
);
```

## Implementation Patterns

### Check Project Membership

```go
// Get user's role in project
func (r *Repo) GetMemberRole(ctx context.Context, projectID, userID int64) (string, error) {
    var role string
    err := r.db.QueryRowContext(ctx, `
        SELECT role FROM altalune_project_members
        WHERE project_id = $1 AND user_id = $2
    `, projectID, userID).Scan(&role)

    if err == sql.ErrNoRows {
        return "", ErrNotAMember
    }
    return role, err
}

// Check minimum role
func hasMinimumRole(userRole, minRole string) bool {
    hierarchy := map[string]int{"user": 1, "member": 2, "admin": 3, "owner": 4}
    return hierarchy[userRole] >= hierarchy[minRole]
}
```

### Service-Level Authorization

```go
func (s *Service) UpdateProjectSettings(ctx context.Context, req *pb.UpdateProjectRequest) error {
    userID, _ := authctx.GetUserID(ctx)

    // Check membership with minimum admin role
    role, err := s.memberRepo.GetMemberRole(ctx, req.ProjectId, userID)
    if err != nil {
        return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("not a project member"))
    }

    if !hasMinimumRole(role, "admin") {
        return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("requires admin role"))
    }

    // Proceed with update...
}
```

### Auto-Register Superadmin

```go
// Called when new project is created
func (r *Repo) registerSuperadminAsOwner(ctx context.Context, projectID int64) error {
    const superadminID = int64(1)

    _, err := r.db.ExecContext(ctx, `
        INSERT INTO altalune_project_members (public_id, project_id, user_id, role)
        VALUES ($1, $2, $3, 'owner')
        ON CONFLICT (project_id, user_id) DO NOTHING
    `, nanoid.GeneratePublicID(), projectID, superadminID)

    return err
}
```

### Auto-Register OAuth User

```go
// Called during OAuth callback when user doesn't exist in project
func (s *AuthService) autoRegisterProjectMember(ctx context.Context, projectID, userID int64) error {
    return s.memberRepo.Create(ctx, CreateMemberInput{
        ProjectID: projectID,
        UserID:    userID,
        Role:      "user",  // Default role for OAuth users
    })
}
```

## Frontend Usage

### Get Current User's Role

```typescript
// stores/auth.ts
export const useAuthStore = defineStore('auth', () => {
  const currentProject = ref<string | null>(null)
  const projectMemberships = ref<ProjectMembership[]>([])

  const currentRole = computed(() => {
    const membership = projectMemberships.value.find(
      m => m.projectId === currentProject.value
    )
    return membership?.role ?? null
  })

  const isOwner = computed(() => currentRole.value === 'owner')
  const isAdmin = computed(() => ['owner', 'admin'].includes(currentRole.value ?? ''))
  const isMember = computed(() => ['owner', 'admin', 'member'].includes(currentRole.value ?? ''))

  return { currentRole, isOwner, isAdmin, isMember }
})
```

### Role-Based UI

```vue
<template>
  <div>
    <!-- Admin+ only -->
    <Button v-if="isAdmin" @click="openSettings">Settings</Button>

    <!-- Owner only -->
    <Button v-if="isOwner" variant="destructive" @click="deleteProject">
      Delete Project
    </Button>

    <!-- Member+ can view -->
    <ProjectData v-if="isMember" :project-id="projectId" />
  </div>
</template>

<script setup>
const { isOwner, isAdmin, isMember } = useAuthStore()
</script>
```

## Permission vs Membership

| Aspect | Project Membership | RBAC Permission |
|--------|-------------------|-----------------|
| Granularity | Coarse (4 roles) | Fine (many permissions) |
| Scope | Per-project | System-wide |
| Check | Project-specific | Any endpoint |
| Use Case | "Can user access this project?" | "Can user do this action?" |

**Recommended Pattern:** Use BOTH

```go
func (s *Service) DeleteEmployee(ctx context.Context, req *pb.DeleteRequest) error {
    // 1. Check project membership (coarse)
    role, err := s.memberRepo.GetMemberRole(ctx, req.ProjectId, userID)
    if err != nil || !hasMinimumRole(role, "member") {
        return connect.NewError(connect.CodePermissionDenied, ...)
    }

    // 2. Check permission (fine)
    if !authctx.HasPermission(ctx, "employee:delete") {
        return connect.NewError(connect.CodePermissionDenied, ...)
    }

    // 3. Proceed
    return s.repo.Delete(ctx, req.ProjectId, req.Id)
}
```

## Validation Approach: Database Query vs JWT Claims

### Option A: Database Query (Recommended)

Query membership on each request that requires project access.

```go
// In service or interceptor
role, err := s.memberRepo.GetMemberRole(ctx, projectID, userID)
if err != nil {
    return connect.NewError(connect.CodePermissionDenied, ...)
}
```

**Pros:**
- Always fresh - role changes take effect immediately
- Simple JWT structure (no membership data)
- No token refresh required after role changes
- Works with any number of projects

**Cons:**
- Extra DB query per request
- Slightly higher latency

**Mitigation:** Use Redis cache with short TTL (30-60s) for frequently accessed memberships.

### Option B: JWT Claims with Membership Map

Embed project memberships in JWT access token.

```json
{
  "sub": "user_public_id",
  "perms": ["employee:read", "employee:create"],
  "memberships": {
    "proj_abc123": "admin",
    "proj_def456": "member"
  }
}
```

**Pros:**
- No DB query needed for membership check
- Stateless validation
- Lower latency

**Cons:**
- Stale data - role changes don't apply until token refresh
- JWT size grows with project count
- Requires token refresh mechanism on role change
- Complex invalidation logic

### Recommendation: Database Query + Caching

Use **Option A (Database Query)** for these reasons:

1. **Multi-project users**: Admin users may belong to many projects; embedding all in JWT bloats token size
2. **Real-time accuracy**: Role changes by owner should take effect immediately
3. **Simplicity**: No need for complex token refresh/invalidation flows
4. **Separation of concerns**: RBAC permissions in JWT, membership in DB

**Implementation Pattern:**

```go
// Cache membership for performance (optional)
type MembershipCache struct {
    cache *redis.Client
    repo  *MemberRepo
}

func (c *MembershipCache) GetRole(ctx context.Context, projectID, userID string) (string, error) {
    key := fmt.Sprintf("membership:%s:%s", projectID, userID)

    // Try cache first
    role, err := c.cache.Get(ctx, key).Result()
    if err == nil {
        return role, nil
    }

    // Fallback to DB
    role, err = c.repo.GetMemberRole(ctx, projectID, userID)
    if err != nil {
        return "", err
    }

    // Cache for 60 seconds
    c.cache.Set(ctx, key, role, 60*time.Second)
    return role, nil
}

// Invalidate on role change
func (c *MembershipCache) InvalidateMembership(ctx context.Context, projectID, userID string) {
    key := fmt.Sprintf("membership:%s:%s", projectID, userID)
    c.cache.Del(ctx, key)
}
```

**When to use Option B (JWT Claims):**
- Single-project applications
- Users typically belong to 1-2 projects max
- Role changes are rare and can wait for token expiry

## Related

- **[middleware-patterns.md](middleware-patterns.md)** - Authorization middleware
- **[permission-definitions.md](permission-definitions.md)** - Predefined permissions
