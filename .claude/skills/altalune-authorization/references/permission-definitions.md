# Permission Definitions

## Permission Naming Convention

Format: `{resource}:{action}` or `{resource}:{sub-resource}:{action}`

**Resources:** project, employee, user, role, permission, oauth_client, api_key, chatbot, etc.
**Actions:** read, create, update, delete, manage, admin

## Predefined Permissions

### System-Level

| Permission | Description | Typical Roles |
|------------|-------------|---------------|
| `root` | Full system access (wildcard) | superadmin |
| `admin:access` | Access admin dashboard | owner, admin |
| `api:access` | Access API endpoints | all authenticated |

### Project Permissions

| Permission | Description | Typical Roles |
|------------|-------------|---------------|
| `project:read` | View projects | member+ |
| `project:create` | Create new projects | owner |
| `project:update` | Update project settings | admin+ |
| `project:delete` | Delete projects | owner |
| `project:member:read` | View project members | member+ |
| `project:member:manage` | Add/remove members | admin+ |

### Employee Permissions (Example Entity)

| Permission | Description | Typical Roles |
|------------|-------------|---------------|
| `employee:read` | View employees | member+ |
| `employee:create` | Create employees | member+ |
| `employee:update` | Update employees | member+ |
| `employee:delete` | Delete employees | admin+ |

### User Permissions

| Permission | Description | Typical Roles |
|------------|-------------|---------------|
| `user:read` | View user profiles | member+ |
| `user:update` | Update user profiles | admin+ (or self) |
| `user:delete` | Delete users | owner |
| `user:role:manage` | Assign roles to users | admin+ |

### OAuth Client Permissions

| Permission | Description | Typical Roles |
|------------|-------------|---------------|
| `oauth_client:read` | View OAuth clients | member+ |
| `oauth_client:create` | Create OAuth clients | owner |
| `oauth_client:update` | Update OAuth clients | admin+ |
| `oauth_client:delete` | Delete OAuth clients | owner |
| `oauth_client:secret:reveal` | Reveal client secrets | admin+ |

### API Key Permissions

| Permission | Description | Typical Roles |
|------------|-------------|---------------|
| `api_key:read` | View API keys | admin+ |
| `api_key:create` | Create API keys | admin+ |
| `api_key:delete` | Delete API keys | admin+ |

### IAM Permissions

| Permission | Description | Typical Roles |
|------------|-------------|---------------|
| `role:read` | View roles | admin+ |
| `role:create` | Create roles | owner |
| `role:update` | Update roles | owner |
| `role:delete` | Delete roles | owner |
| `permission:read` | View permissions | admin+ |
| `permission:manage` | Manage permissions | owner |

### Chatbot Permissions

| Permission | Description | Typical Roles |
|------------|-------------|---------------|
| `chatbot:module:read` | View chatbot module config | member+ |
| `chatbot:module:update` | Update chatbot module config | admin+ |
| `chatbot:node:read` | View chatbot nodes | member+ |
| `chatbot:node:manage` | Manage chatbot nodes | admin+ |

## Database Migration

### Generate Public IDs

Use the public ID generator for migrations:

```bash
# Single ID
go run cmd/public_id/main.go

# Batch (e.g., 40 IDs)
go run cmd/public_id/main.go -b -c 40
```

### Seed Predefined Permissions

```sql
-- +goose Up
-- +goose StatementBegin

-- System permissions
INSERT INTO altalune_permissions (public_id, name, effect, description) VALUES
  ('h2lbbe8ptky8cf', 'root', 'allow', 'Full system access'),
  ('vbdw8s5z8d3aht', 'admin:access', 'allow', 'Access admin dashboard'),
  ('fgv3s8n2xj3buv', 'api:access', 'allow', 'Access API endpoints');

-- Project permissions
INSERT INTO altalune_permissions (public_id, name, effect, description) VALUES
  ('av273atq8ckqbr', 'project:read', 'allow', 'View projects'),
  ('he48qxa8yjyqvs', 'project:create', 'allow', 'Create projects'),
  ('r3mf7jw28xkhdk', 'project:update', 'allow', 'Update projects'),
  ('rcy4b8m9jzq2sa', 'project:delete', 'allow', 'Delete projects'),
  ('7jm98fvssmb8n9', 'project:member:read', 'allow', 'View project members'),
  ('stt6a7jz8l88d7', 'project:member:manage', 'allow', 'Manage project members');

-- Employee permissions
INSERT INTO altalune_permissions (public_id, name, effect, description) VALUES
  ('4tdshx5eung286', 'employee:read', 'allow', 'View employees'),
  ('99bgnwkjfusvnt', 'employee:create', 'allow', 'Create employees'),
  ('4r3nzx6aza897t', 'employee:update', 'allow', 'Update employees'),
  ('wdf83azqs3mbrs', 'employee:delete', 'allow', 'Delete employees');

-- User permissions
INSERT INTO altalune_permissions (public_id, name, effect, description) VALUES
  ('4j263hzxcbv5a4', 'user:read', 'allow', 'View users'),
  ('qzuquegjyvmskc', 'user:update', 'allow', 'Update users'),
  ('5e23c2yena2r3h', 'user:delete', 'allow', 'Delete users'),
  ('5vsry8xkf49cpr', 'user:role:manage', 'allow', 'Manage user roles');

-- OAuth Client permissions
INSERT INTO altalune_permissions (public_id, name, effect, description) VALUES
  ('csechpgvckj2p6', 'oauth_client:read', 'allow', 'View OAuth clients'),
  ('delpywt72wgwkw', 'oauth_client:create', 'allow', 'Create OAuth clients'),
  ('xyqycuvnsf8yqs', 'oauth_client:update', 'allow', 'Update OAuth clients'),
  ('59bqh4p2se9y8g', 'oauth_client:delete', 'allow', 'Delete OAuth clients'),
  ('fhbu38552c7dfm', 'oauth_client:secret:reveal', 'allow', 'Reveal client secrets');

-- API Key permissions
INSERT INTO altalune_permissions (public_id, name, effect, description) VALUES
  ('z4bxq55829mjka', 'api_key:read', 'allow', 'View API keys'),
  ('mq3r3kp2n23kwj', 'api_key:create', 'allow', 'Create API keys'),
  ('mr9zmarxtgc3z7', 'api_key:delete', 'allow', 'Delete API keys');

-- IAM permissions
INSERT INTO altalune_permissions (public_id, name, effect, description) VALUES
  ('vb3frg3f3ua2sr', 'role:read', 'allow', 'View roles'),
  ('8dkn6gxw464rzy', 'role:create', 'allow', 'Create roles'),
  ('ncq8m5nceasdbt', 'role:update', 'allow', 'Update roles'),
  ('9x92rzn528f6da', 'role:delete', 'allow', 'Delete roles'),
  ('khfkrrp3xkr9ds', 'permission:read', 'allow', 'View permissions'),
  ('9u8mt9apz5nub5', 'permission:manage', 'allow', 'Manage permissions');

-- Chatbot permissions
INSERT INTO altalune_permissions (public_id, name, effect, description) VALUES
  ('4dn9cuwterml83', 'chatbot:module:read', 'allow', 'View chatbot module config'),
  ('7azdgsj89a6m6k', 'chatbot:module:update', 'allow', 'Update chatbot module config'),
  ('jt49ynzk45uflk', 'chatbot:node:read', 'allow', 'View chatbot nodes'),
  ('lmhszxstzc6x8h', 'chatbot:node:manage', 'allow', 'Manage chatbot nodes');

-- +goose StatementEnd
```

### Create Default Roles with Permissions

```sql
-- Generate IDs: go run cmd/public_id/main.go -b -c 3
INSERT INTO altalune_roles (public_id, name, description) VALUES
  ('b6bswz6wunyll6', 'admin', 'Project administrator'),
  ('trfg8acrqvrq62', 'member', 'Project member'),
  ('mlp9gfhax3m228', 'user', 'OAuth user (minimal)');

-- Assign permissions to admin role
INSERT INTO altalune_roles_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM altalune_roles r, altalune_permissions p
WHERE r.name = 'admin'
  AND p.name IN (
    'admin:access', 'api:access',
    'project:read', 'project:update', 'project:member:read', 'project:member:manage',
    'employee:read', 'employee:create', 'employee:update', 'employee:delete',
    'user:read', 'user:update', 'user:role:manage',
    'oauth_client:read', 'oauth_client:update', 'oauth_client:secret:reveal',
    'api_key:read', 'api_key:create', 'api_key:delete',
    'role:read', 'permission:read',
    'chatbot:module:read', 'chatbot:module:update', 'chatbot:node:read', 'chatbot:node:manage'
  );

-- Assign permissions to member role
INSERT INTO altalune_roles_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM altalune_roles r, altalune_permissions p
WHERE r.name = 'member'
  AND p.name IN (
    'api:access',
    'project:read', 'project:member:read',
    'employee:read', 'employee:create', 'employee:update',
    'user:read',
    'oauth_client:read',
    'chatbot:module:read', 'chatbot:node:read'
  );

-- Assign permissions to user role
INSERT INTO altalune_roles_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM altalune_roles r, altalune_permissions p
WHERE r.name = 'user'
  AND p.name IN ('api:access');
```

## Go Constants

**File:** `internal/domain/permission/constants.go`

```go
package permission

// Predefined permission names
const (
    // System
    PermRoot        = "root"
    PermAdminAccess = "admin:access"
    PermAPIAccess   = "api:access"

    // Project
    PermProjectRead         = "project:read"
    PermProjectCreate       = "project:create"
    PermProjectUpdate       = "project:update"
    PermProjectDelete       = "project:delete"
    PermProjectMemberRead   = "project:member:read"
    PermProjectMemberManage = "project:member:manage"

    // Employee
    PermEmployeeRead   = "employee:read"
    PermEmployeeCreate = "employee:create"
    PermEmployeeUpdate = "employee:update"
    PermEmployeeDelete = "employee:delete"

    // User
    PermUserRead       = "user:read"
    PermUserUpdate     = "user:update"
    PermUserDelete     = "user:delete"
    PermUserRoleManage = "user:role:manage"

    // OAuth Client
    PermOAuthClientRead         = "oauth_client:read"
    PermOAuthClientCreate       = "oauth_client:create"
    PermOAuthClientUpdate       = "oauth_client:update"
    PermOAuthClientDelete       = "oauth_client:delete"
    PermOAuthClientSecretReveal = "oauth_client:secret:reveal"

    // API Key
    PermAPIKeyRead   = "api_key:read"
    PermAPIKeyCreate = "api_key:create"
    PermAPIKeyDelete = "api_key:delete"

    // IAM
    PermRoleRead         = "role:read"
    PermRoleCreate       = "role:create"
    PermRoleUpdate       = "role:update"
    PermRoleDelete       = "role:delete"
    PermPermissionRead   = "permission:read"
    PermPermissionManage = "permission:manage"

    // Chatbot
    PermChatbotModuleRead   = "chatbot:module:read"
    PermChatbotModuleUpdate = "chatbot:module:update"
    PermChatbotNodeRead     = "chatbot:node:read"
    PermChatbotNodeManage   = "chatbot:node:manage"
)
```

## Frontend Permission Check

```typescript
// composables/usePermissions.ts
export function usePermissions() {
  const authStore = useAuthStore()

  const permissions = computed(() => authStore.user?.perms ?? [])

  const hasPermission = (perm: string): boolean => {
    return permissions.value.includes('root') || permissions.value.includes(perm)
  }

  const hasAnyPermission = (...perms: string[]): boolean => {
    return perms.some(p => hasPermission(p))
  }

  return { permissions, hasPermission, hasAnyPermission }
}
```
