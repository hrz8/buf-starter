# User Story US15: Authorization - Project Membership & RBAC Validation

## Story Overview

**As a** system administrator and backend developer
**I want** the backend/resource server to validate access tokens and enforce authorization using project membership and RBAC permissions
**So that** users can only access resources they are authorized to, with secure JWT validation and proper permission checks at both backend and frontend levels

## Acceptance Criteria

### Core Functionality

#### JWT Validation in Resource Server

- **Given** a request arrives at a Connect-RPC endpoint
- **When** the Authorization header contains a Bearer token
- **Then** the resource server validates the JWT signature using the public key from JWKS endpoint
- **And** the JWT claims (sub, perms, memberships, email_verified) are extracted into request context
- **And** expired or invalid tokens return `UNAUTHENTICATED` error

- **Given** the resource server starts up
- **When** it initializes the JWT validator
- **Then** it fetches the public key from auth server's `/.well-known/jwks.json`
- **And** it caches the JWKS response in memory
- **And** it refreshes the cache periodically (every 1 hour) and on signature validation failure

#### JWT Claims Extension - Memberships Map

- **Given** a user authenticates via OAuth and receives an access token
- **When** the token is generated
- **Then** the JWT contains a `memberships` claim with format: `{ "project_public_id": "role" }`
- **Example**: `{ "proj_abc123": "admin", "proj_xyz789": "member" }`

- **Given** a user's project membership changes (added/removed/role changed)
- **When** they refresh their access token
- **Then** the new token reflects the updated memberships

#### Authorization Flow

- **Given** a user requests access to a project-scoped resource
- **When** the authorization check runs
- **Then** it follows this order:
  1. Check if user is superadmin (has `root` permission) → Allow all
  2. Check if user has the required global permission (e.g., `employee:read`)
  3. If yes, check if user is a member of the target project
  4. If both checks pass → Allow
  5. If either check fails → Reject with `PERMISSION_DENIED`

- **Given** a superadmin user (has `root` permission)
- **When** they access any resource in any project
- **Then** they bypass project membership checks
- **And** they have access to all resources

#### Project Membership Validation

- **Given** a user accesses the dashboard
- **When** they view the project sidebar
- **Then** only projects they are members of are displayed
- **And** the project list is derived from the `memberships` claim in JWT

- **Given** a user tries to switch to a project
- **When** they are not a member of that project
- **Then** the request is rejected with `PERMISSION_DENIED`
- **And** the error message indicates they don't have access to this project

#### RBAC Permission Validation

- **Given** a user with role "user" (has `dashboard:read` permission)
- **When** they try to access employee list
- **Then** the request is rejected because they lack `employee:read` permission

- **Given** a user with role having `employee:read` permission
- **When** they try to create an employee
- **Then** the request is rejected because they lack `employee:write` permission

- **Given** a user with role having `employee:read` and `employee:write` permissions
- **When** they try to delete an employee
- **Then** the request is rejected because they lack `employee:delete` permission

### Security Requirements

#### JWT Validation Security

- JWT signature MUST be validated using RS256 algorithm
- Public key MUST be fetched from trusted JWKS endpoint
- Token expiration MUST be enforced
- Token audience claim MUST match expected client IDs (optional, configurable)
- Invalid tokens MUST return clear error without exposing internal details

#### JWKS Caching

- JWKS response cached in memory (global variable)
- Cache TTL: 1 hour (configurable)
- Force refresh on signature validation failure (key rotation support)
- Maximum 3 refresh attempts per minute (rate limiting)

#### Permission Enforcement

- Authorization checks MUST happen server-side (frontend hiding is UX only)
- All Connect-RPC handlers MUST call authorization helper
- Missing authorization check MUST default to deny
- Audit log for authorization failures (optional, future)

### Data Model

#### JWT Claims Structure

```go
type AccessTokenClaims struct {
    jwt.RegisteredClaims
    Scope         string            `json:"scope,omitempty"`
    Email         string            `json:"email,omitempty"`
    Name          string            `json:"name,omitempty"`
    Perms         []string          `json:"perms"`
    Memberships   map[string]string `json:"memberships"` // NEW: project_public_id -> role
    EmailVerified bool              `json:"email_verified"`
}
```

#### Predefined Permissions (Migration)

```sql
-- Entity permissions (entity:action format)
INSERT INTO altalune_permissions (public_id, name, description) VALUES
-- Employee entity
('perm_emp_r', 'employee:read', 'View employees'),
('perm_emp_w', 'employee:write', 'Create and update employees'),
('perm_emp_d', 'employee:delete', 'Delete employees'),
-- User entity (IAM)
('perm_usr_r', 'user:read', 'View users'),
('perm_usr_w', 'user:write', 'Create and update users'),
('perm_usr_d', 'user:delete', 'Delete users'),
-- Role entity (IAM)
('perm_rol_r', 'role:read', 'View roles'),
('perm_rol_w', 'role:write', 'Create and update roles'),
('perm_rol_d', 'role:delete', 'Delete roles'),
-- Permission entity (IAM)
('perm_prm_r', 'permission:read', 'View permissions'),
('perm_prm_w', 'permission:write', 'Create and update permissions'),
('perm_prm_d', 'permission:delete', 'Delete permissions'),
-- Project entity
('perm_prj_r', 'project:read', 'View projects'),
('perm_prj_w', 'project:write', 'Create and update projects'),
('perm_prj_d', 'project:delete', 'Delete projects'),
-- API Key entity
('perm_api_r', 'apikey:read', 'View API keys'),
('perm_api_w', 'apikey:write', 'Create and update API keys'),
('perm_api_d', 'apikey:delete', 'Delete API keys'),
-- Chatbot Config entity
('perm_bot_r', 'chatbot:read', 'View chatbot configurations'),
('perm_bot_w', 'chatbot:write', 'Create and update chatbot configurations'),
('perm_bot_d', 'chatbot:delete', 'Delete chatbot configurations'),
-- OAuth Client entity
('perm_cli_r', 'client:read', 'View OAuth clients'),
('perm_cli_w', 'client:write', 'Create and update OAuth clients'),
('perm_cli_d', 'client:delete', 'Delete OAuth clients'),
-- Project Members management
('perm_mem_r', 'member:read', 'View project members'),
('perm_mem_w', 'member:write', 'Add and update project members'),
('perm_mem_d', 'member:delete', 'Remove project members'),
-- IAM Mapper (role-permission assignments)
('perm_iam_r', 'iam:read', 'View IAM mappings'),
('perm_iam_w', 'iam:write', 'Manage IAM mappings');

-- Also verify and drop 'effect' column if still exists
-- ALTER TABLE altalune_permissions DROP COLUMN IF EXISTS effect;
```

### User Experience

#### Frontend Permission-Based UI

- **Sidebar Navigation**: Show/hide menu items based on user permissions
- **Action Buttons**: Show/hide create/edit/delete buttons based on permissions
- **Page Access**: Redirect to "Access Denied" page if user lacks permission
- **Loading States**: Don't flash unauthorized UI before hiding

#### Permission Store Design (Frontend)

```typescript
// composables/usePermissions.ts
interface PermissionState {
  permissions: string[]       // ['employee:read', 'employee:write', ...]
  memberships: Record<string, string>  // { 'proj_abc': 'admin', ... }
  isLoaded: boolean
}

// Helper functions
function hasPermission(permission: string): boolean
function hasAnyPermission(permissions: string[]): boolean
function hasAllPermissions(permissions: string[]): boolean
function isMemberOf(projectId: string): boolean
function getProjectRole(projectId: string): string | null
function isSuperAdmin(): boolean
```

#### Component Examples

```vue
<!-- Permission-based rendering -->
<template>
  <Button v-if="can('employee:write')" @click="createEmployee">
    Create Employee
  </Button>

  <MenuItem v-if="can('user:read')" :to="'/users'">
    User Management
  </MenuItem>
</template>

<script setup>
const { can, canAny, canAll } = usePermissions()
</script>
```

### Technical Requirements

### Backend Architecture

#### Directory Structure

```
internal/
├── auth/
│   ├── context.go        # Request context with claims
│   ├── middleware.go     # Connect-RPC interceptor
│   ├── validator.go      # JWT validation logic
│   ├── jwks.go           # JWKS fetching and caching
│   └── authorizer.go     # Authorization helper functions
```

#### JWT Validator

```go
// internal/auth/validator.go
type JWTValidator struct {
    jwksURL     string
    jwksCache   *JWKSCache
    issuer      string
    audiences   []string // optional audience validation
}

func NewJWTValidator(jwksURL, issuer string, audiences []string) *JWTValidator
func (v *JWTValidator) Validate(tokenString string) (*AccessTokenClaims, error)
func (v *JWTValidator) RefreshJWKS(ctx context.Context) error
```

#### JWKS Cache

```go
// internal/auth/jwks.go
type JWKSCache struct {
    mu          sync.RWMutex
    keys        map[string]*rsa.PublicKey  // kid -> public key
    lastFetch   time.Time
    ttl         time.Duration
}

func NewJWKSCache(ttl time.Duration) *JWKSCache
func (c *JWKSCache) GetKey(kid string) (*rsa.PublicKey, bool)
func (c *JWKSCache) SetKeys(keys map[string]*rsa.PublicKey)
func (c *JWKSCache) IsExpired() bool
```

#### Request Context

```go
// internal/auth/context.go
type AuthContext struct {
    UserID        string            // JWT subject
    Email         string
    Permissions   []string
    Memberships   map[string]string // project_public_id -> role
    EmailVerified bool
    IsAuthenticated bool
}

func FromContext(ctx context.Context) *AuthContext
func WithAuthContext(ctx context.Context, auth *AuthContext) context.Context
```

#### Authorization Helper

```go
// internal/auth/authorizer.go
type Authorizer struct{}

// CheckPermission validates user has required permission
// Returns error if unauthorized
func (a *Authorizer) CheckPermission(ctx context.Context, permission string) error

// CheckProjectAccess validates user has permission AND is project member
// projectID is the project's public_id
func (a *Authorizer) CheckProjectAccess(ctx context.Context, permission string, projectID string) error

// CheckProjectMembership validates user is member of project
func (a *Authorizer) CheckProjectMembership(ctx context.Context, projectID string) error

// IsSuperAdmin checks if user has root permission
func (a *Authorizer) IsSuperAdmin(ctx context.Context) bool

// GetUserProjects returns list of project IDs user has access to
func (a *Authorizer) GetUserProjects(ctx context.Context) []string
```

#### Connect-RPC Interceptor

```go
// internal/auth/middleware.go
func AuthInterceptor(validator *JWTValidator) connect.UnaryInterceptorFunc {
    return func(next connect.UnaryFunc) connect.UnaryFunc {
        return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
            // Extract Bearer token from Authorization header
            // Validate JWT
            // Inject AuthContext into request context
            // Call next handler (handler does its own authorization check)
            return next(ctx, req)
        }
    }
}
```

#### Handler Authorization Pattern

```go
// Example: internal/domain/employee/handler.go
func (h *Handler) ListEmployees(ctx context.Context, req *connect.Request[v1.ListEmployeesRequest]) (*connect.Response[v1.ListEmployeesResponse], error) {
    // Authorization check
    if err := h.authorizer.CheckProjectAccess(ctx, "employee:read", req.Msg.ProjectId); err != nil {
        return nil, err
    }

    // Proceed with business logic
    employees, err := h.service.List(ctx, req.Msg.ProjectId, req.Msg.Page, req.Msg.PageSize)
    // ...
}
```

### Frontend Architecture

#### Directory Structure

```
frontend/
├── shared/
│   ├── composables/
│   │   └── usePermissions.ts    # Permission state and helpers
│   └── stores/
│       └── auth.ts              # Auth store with JWT decoding
├── app/
│   ├── plugins/
│   │   └── auth.client.ts       # Initialize auth on app load
│   ├── middleware/
│   │   └── auth.ts              # Route guard
│   └── components/
│       └── PermissionGuard.vue  # Declarative permission check
```

#### Permission Composable

```typescript
// shared/composables/usePermissions.ts
export function usePermissions() {
  const authStore = useAuthStore()

  const permissions = computed(() => authStore.permissions)
  const memberships = computed(() => authStore.memberships)

  const can = (permission: string): boolean => {
    if (authStore.isSuperAdmin) return true
    return permissions.value.includes(permission)
  }

  const canAny = (perms: string[]): boolean => {
    if (authStore.isSuperAdmin) return true
    return perms.some(p => permissions.value.includes(p))
  }

  const canAll = (perms: string[]): boolean => {
    if (authStore.isSuperAdmin) return true
    return perms.every(p => permissions.value.includes(p))
  }

  const isMemberOf = (projectId: string): boolean => {
    if (authStore.isSuperAdmin) return true
    return projectId in memberships.value
  }

  const getProjectRole = (projectId: string): string | null => {
    return memberships.value[projectId] ?? null
  }

  return { can, canAny, canAll, isMemberOf, getProjectRole, permissions, memberships }
}
```

#### Auth Store Extension

```typescript
// shared/stores/auth.ts
export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(null)
  const decodedToken = ref<DecodedToken | null>(null)

  const permissions = computed(() => decodedToken.value?.perms ?? [])
  const memberships = computed(() => decodedToken.value?.memberships ?? {})
  const isSuperAdmin = computed(() => permissions.value.includes('root'))

  function decodeToken(token: string) {
    // Decode JWT payload (no verification - that's server's job)
    const payload = JSON.parse(atob(token.split('.')[1]))
    decodedToken.value = payload
  }

  return { accessToken, permissions, memberships, isSuperAdmin, decodeToken }
})
```

#### Permission Guard Component

```vue
<!-- components/PermissionGuard.vue -->
<script setup lang="ts">
const props = defineProps<{
  permission?: string
  permissions?: string[]
  mode?: 'any' | 'all'
  fallback?: 'hide' | 'disabled' | 'redirect'
}>()

const { can, canAny, canAll } = usePermissions()

const hasAccess = computed(() => {
  if (props.permission) return can(props.permission)
  if (props.permissions) {
    return props.mode === 'all'
      ? canAll(props.permissions)
      : canAny(props.permissions)
  }
  return true
})
</script>

<template>
  <slot v-if="hasAccess" />
  <slot v-else-if="fallback === 'disabled'" name="disabled" />
  <slot v-else name="fallback" />
</template>
```

#### Sidebar Navigation with Permissions

```typescript
// Navigation items with required permissions
const navigationItems = [
  {
    label: 'Employees',
    to: '/employees',
    icon: 'lucide:users',
    permission: 'employee:read'
  },
  {
    label: 'Users',
    to: '/users',
    icon: 'lucide:user-cog',
    permission: 'user:read'
  },
  {
    label: 'Roles',
    to: '/roles',
    icon: 'lucide:shield',
    permission: 'role:read'
  },
  // ... etc
]

// Filter by permission
const visibleNavItems = computed(() =>
  navigationItems.filter(item => !item.permission || can(item.permission))
)
```

### Configuration

```yaml
# config.yaml additions
auth:
  jwks:
    url: "http://localhost:3300/.well-known/jwks.json"  # Auth server JWKS endpoint
    cacheTTL: 3600                                       # Cache TTL in seconds (1 hour)
    refreshRetryLimit: 3                                 # Max refresh attempts per minute
  issuer: "http://localhost:3300"                        # Expected JWT issuer
  audiences: []                                          # Optional: expected audiences (empty = skip check)
```

## Out of Scope

- Token introspection endpoint (using local validation with JWKS)
- Permission audit logging (future enhancement)
- Rate limiting per user/permission (future enhancement)
- Dynamic permission loading from server (using JWT claims)
- Project-scoped RBAC (permissions are global, membership is per-project)
- Role hierarchy/inheritance (explicit permissions per role)
- UI for role-permission management (already exists, this story adds enforcement)

## Dependencies

- US5: OAuth Server Foundation (IAM tables)
- US7: OAuth Authorization Server (JWT generation)
- US8: Dashboard OAuth Integration (token refresh)
- US14: Standalone IDP Application (email_verified claim)
- Existing JWKS endpoint at `/.well-known/jwks.json`
- Existing project membership table (`altalune_project_members`)
- Existing RBAC tables (`altalune_roles`, `altalune_permissions`, etc.)

## Definition of Done

### Database

- [ ] Migration creates all predefined permissions (entity:action format)
- [ ] Migration drops `effect` column from permissions table if exists
- [ ] All existing entities have read/write/delete permissions

### Backend - JWT Validation

- [ ] JWKS fetcher implemented with in-memory caching
- [ ] JWT validator uses public key from JWKS
- [ ] Cache refreshes on TTL expiry (1 hour)
- [ ] Cache refreshes on signature validation failure
- [ ] Rate limiting on JWKS refresh (max 3/minute)

### Backend - Authorization

- [ ] `AuthContext` struct holds parsed JWT claims
- [ ] `AuthInterceptor` extracts and validates JWT from requests
- [ ] `AuthInterceptor` injects `AuthContext` into request context
- [ ] `Authorizer` helper provides permission checking functions
- [ ] `CheckPermission` validates global permission
- [ ] `CheckProjectAccess` validates permission + membership
- [ ] `IsSuperAdmin` checks for `root` permission
- [ ] Unauthenticated requests return `UNAUTHENTICATED` error
- [ ] Unauthorized requests return `PERMISSION_DENIED` error

### Backend - Handler Integration

- [ ] Employee handlers check `employee:read/write/delete`
- [ ] User handlers check `user:read/write/delete`
- [ ] Role handlers check `role:read/write/delete`
- [ ] Permission handlers check `permission:read/write/delete`
- [ ] Project handlers check `project:read/write/delete`
- [ ] API Key handlers check `apikey:read/write/delete`
- [ ] Chatbot handlers check `chatbot:read/write/delete`
- [ ] OAuth Client handlers check `client:read/write/delete`
- [ ] Member handlers check `member:read/write/delete`
- [ ] IAM Mapper handlers check `iam:read/write`
- [ ] All project-scoped handlers validate membership

### Auth Server - JWT Enhancement

- [ ] JWT includes `memberships` claim
- [ ] Memberships fetched from `altalune_project_members` table
- [ ] Memberships format: `{ "project_public_id": "role" }`

### Frontend - Permission System

- [ ] `usePermissions` composable implemented
- [ ] `useAuthStore` extended with permissions/memberships
- [ ] JWT decoded client-side for claims access
- [ ] `can()`, `canAny()`, `canAll()` helper functions
- [ ] `isMemberOf()`, `getProjectRole()` helper functions
- [ ] `isSuperAdmin` computed property

### Frontend - UI Integration

- [ ] Sidebar navigation filtered by permissions
- [ ] Create/Edit/Delete buttons hidden when unauthorized
- [ ] `PermissionGuard` component for declarative checks
- [ ] Project selector shows only member projects
- [ ] Access denied page for unauthorized routes
- [ ] No UI flash before permission check completes

### Configuration

- [ ] `config.yaml` updated with JWKS settings
- [ ] `config.example.yaml` updated with examples

### Testing

- [ ] JWT validation tested (valid, expired, invalid signature)
- [ ] JWKS caching tested (TTL, refresh on failure)
- [ ] Permission checks tested (has/doesn't have)
- [ ] Project access checks tested (member/not member)
- [ ] Superadmin bypass tested
- [ ] Frontend permission helpers tested

## Notes

### Authorization Decision Flow

```
Request arrives at Connect-RPC endpoint
    │
    ▼
AuthInterceptor extracts Bearer token
    │
    ▼
JWT validated using JWKS public key
    │
    ├─ Invalid/Expired → Return UNAUTHENTICATED
    │
    ▼
AuthContext injected into request context
    │
    ▼
Handler calls authorizer.CheckProjectAccess(ctx, "entity:action", projectID)
    │
    ├─ User has 'root' permission? → ALLOW (superadmin bypass)
    │
    ├─ User has required permission?
    │   │
    │   ├─ No → Return PERMISSION_DENIED
    │   │
    │   ▼
    │   User is member of project?
    │       │
    │       ├─ No → Return PERMISSION_DENIED
    │       │
    │       ▼
    │       ALLOW
```

### JWT Claims Example

```json
{
  "iss": "http://localhost:3300",
  "sub": "usr_abc123xyz",
  "aud": ["client_dashboard"],
  "exp": 1704067200,
  "iat": 1704063600,
  "scope": "openid profile email",
  "email": "user@example.com",
  "name": "John Doe",
  "email_verified": true,
  "perms": ["employee:read", "employee:write", "dashboard:read"],
  "memberships": {
    "proj_abc123": "admin",
    "proj_xyz789": "member"
  }
}
```

### Error Messages

| Scenario | Error Code | Message |
|----------|------------|---------|
| No Authorization header | UNAUTHENTICATED | "missing authorization header" |
| Invalid token format | UNAUTHENTICATED | "invalid token format" |
| Expired token | UNAUTHENTICATED | "token has expired" |
| Invalid signature | UNAUTHENTICATED | "invalid token signature" |
| Missing permission | PERMISSION_DENIED | "permission denied: requires {permission}" |
| Not project member | PERMISSION_DENIED | "permission denied: not a member of this project" |

### Frontend Permission Constants

```typescript
// constants/permissions.ts
export const PERMISSIONS = {
  EMPLOYEE: {
    READ: 'employee:read',
    WRITE: 'employee:write',
    DELETE: 'employee:delete',
  },
  USER: {
    READ: 'user:read',
    WRITE: 'user:write',
    DELETE: 'user:delete',
  },
  // ... etc
} as const
```

### JWKS Response Format

```json
{
  "keys": [
    {
      "kty": "RSA",
      "kid": "key-id-1",
      "use": "sig",
      "alg": "RS256",
      "n": "base64url-encoded-modulus",
      "e": "AQAB"
    }
  ]
}
```

### Related Stories

- US5: OAuth Server Foundation (IAM tables)
- US7: OAuth Authorization Server (JWT generation)
- US8: Dashboard OAuth Integration (token refresh)
- US14: Standalone IDP (email_verified claim)

### Future Enhancements

- Permission audit logging
- Token introspection for real-time revocation
- Per-endpoint rate limiting
- Permission caching with Redis
- GraphQL-style field-level permissions
- Attribute-based access control (ABAC)
