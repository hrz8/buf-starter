# Project Membership & OAuth Architecture Guide

## Overview

This guide documents the project membership system, role hierarchy, and OAuth client-project relationship in Altalune. Understanding these concepts is critical for implementing authentication, authorization, and multi-tenant features.

## Table of Contents

- [Core Concepts](#core-concepts)
- [Role Hierarchy](#role-hierarchy)
- [Project Membership Rules](#project-membership-rules)
- [OAuth Client-Project Relationship](#oauth-client-project-relationship)
- [Auto-Registration Behavior](#auto-registration-behavior)
- [Permission Matrix](#permission-matrix)
- [Database Schema](#database-schema)
- [Implementation Guidelines](#implementation-guidelines)

---

## Core Concepts

### Project-Centric Multi-Tenancy

Altalune uses a **project-centric multi-tenant architecture** where:

1. **Projects** are the primary isolation boundary
2. **OAuth clients** are attached to specific projects
3. **Users** gain access to projects through **project membership**
4. **Roles** define what users can do within projects

### Key Principle: OAuth Client → Project → User

```
OAuth Client ──belongs to──> Project ──has many──> Project Members (Users)
```

**Critical Rule**: Every OAuth client is attached to a specific project. When users authenticate via an OAuth client, they are automatically registered as members of that client's project.

---

## Role Hierarchy

Altalune defines **4 distinct roles** with hierarchical permissions:

### 1. Owner (Superadmin Only)

**Scope**: System-wide superadmin role

**Characteristics**:
- Only assigned to the superadmin user (user_id=1)
- Automatically assigned to ALL projects via auto-registration
- Cannot be manually assigned to other users
- Highest level of access across the entire system

**Permissions**:
- Full access to all projects
- Create new projects
- Delete projects
- Full OAuth client management (create, edit, delete, reveal secrets)
- Manage all project members
- Access all settings and configuration

**Implementation Note**: The `internal/domain/project/repo.go` automatically registers superadmin (user_id=1) as owner when any new project is created. See `registerSuperadminAsOwner()` method.

### 2. Admin

**Scope**: Project-level administrator role

**Characteristics**:
- Can be assigned to multiple projects
- Has elevated permissions but cannot delete projects
- Can manage other members (except owners)

**Permissions**:
- Manage project settings (name, description, configuration)
- Manage project members (add, remove, change roles for member/user)
- Edit OAuth clients (cannot delete)
- Reveal OAuth client secrets (with audit logging)
- Manage API keys
- Access project data and analytics
- **Cannot**: Delete projects, assign/remove owner role, create new projects

### 3. Member

**Scope**: Project-level regular user role

**Characteristics**:
- Standard project member with read/write access to project data
- Can be assigned to multiple projects
- No administrative capabilities

**Permissions**:
- Access project data (employees, resources, etc.)
- Read project settings (view only)
- View OAuth clients (read-only, cannot see secrets)
- Use project features and APIs
- **Cannot**: Modify settings, manage members, manage OAuth clients, manage API keys

### 4. User (OAuth User)

**Scope**: OAuth-authenticated user with minimal permissions

**Characteristics**:
- **Default role** assigned on first OAuth login/registration
- Attached to the specific OAuth client they signed up through
- No dashboard access by default
- Lowest privilege level

**Permissions**:
- Authenticate via OAuth
- Access APIs using access tokens (if scopes permit)
- View own profile information
- **Cannot**: Access dashboard, view project data, modify anything

**Upgrade Path**: Admins or owners can upgrade users to member/admin roles to grant dashboard access.

---

## Project Membership Rules

### Rule 1: Auto-Registration on OAuth Login

When a user authenticates via an OAuth client for the **first time**:

1. User record is created in `altalune_users` table
2. User identity is created in `altalune_user_identities` table, linked to the OAuth client
3. **Project membership is automatically created** in `altalune_project_members` with:
   - `project_id` = OAuth client's project_id
   - `user_id` = newly created user's ID
   - `role` = **'user'** (default)

```sql
-- Pseudo-code for auto-registration
1. Create user (email, first_name, last_name from OAuth provider)
2. Create user_identity (user_id, provider='google', oauth_client_id)
3. Create project_members (project_id, user_id, role='user')
```

**Important**: Users are scoped to the OAuth client's project. If an app uses OAuth client A (project 1), users can only access project 1 unless manually granted access to other projects.

### Rule 2: Superadmin Owner Auto-Registration

When a **new project is created**, the superadmin is automatically registered as owner:

```go
// internal/domain/project/repo.go:622-678
func (r *Repo) registerSuperadminAsOwner(ctx context.Context, projectID int64) error {
    const superadminID = int64(1)
    // Check superadmin exists, create owner membership
}
```

This ensures every project has at least one owner (the superadmin).

### Rule 3: Multi-Project Assignment

Users with roles **admin** or **member** can be assigned to multiple projects:

```
User A (email: user@example.com)
├─ Project 1: admin
├─ Project 2: member
└─ Project 3: member
```

**Exception**: The 'user' role typically stays with the original project they signed up through, unless manually upgraded and assigned elsewhere.

### Rule 4: Role Restrictions

- **Owner role**: Cannot be assigned to anyone except superadmin (user_id=1)
- **Role hierarchy enforcement**: Admins cannot modify owner memberships
- **Self-demotion prevention**: Users cannot remove their own owner/admin role if they're the last admin

### Rule 5: Project Deletion & Cascade

When a project is deleted:
- All project members are removed (CASCADE)
- All OAuth clients for that project are removed (CASCADE)
- All partitioned tables (e.g., `altalune_example_employees_p{project_id}`) are dropped

---

## OAuth Client-Project Relationship

### Client-Project Binding

Every OAuth client **must be attached to exactly one project**:

```sql
CREATE TABLE oauth_clients (
  id BIGSERIAL,
  project_id BIGINT NOT NULL,  -- Required, defines project ownership
  client_id UUID UNIQUE NOT NULL,
  client_secret_hash VARCHAR(255) NOT NULL,
  redirect_uris TEXT[] NOT NULL,
  ...
  FOREIGN KEY (project_id) REFERENCES altalune_projects(id) ON DELETE CASCADE
) PARTITION BY LIST (project_id);
```

**Why This Matters**:
1. **Data isolation**: OAuth clients can only access their project's data
2. **Permission scoping**: Users authenticated via a client inherit project context
3. **Multi-tenancy**: Different apps/clients can have separate data spaces

### Default Dashboard Client

The dashboard itself is an OAuth client:

```yaml
# config.yaml
seeder:
  defaultOAuthClient:
    name: "Altalune Dashboard"
    clientId: "e730207a-0fce-495d-bac3-6211963ac423"
    clientSecret: "0cMw4XzRZcRI4YDEqoY9AYWui3y4eZTQ"
    redirectUris:
      - "http://localhost:3000/auth/callback"
    pkceRequired: true
```

**Special Rules for Dashboard Client**:
- Marked with `is_default = true` in database
- Cannot be deleted (UI blocks deletion)
- Used for dashboard authentication
- Users who sign up via dashboard get project membership in dashboard's project

### External App Clients

External applications register their own OAuth clients:

```
Example App "TaskManager"
├─ OAuth Client ID: xxx-yyy-zzz
├─ Project ID: 5
└─ Users who sign up via TaskManager
   └─ Automatically become members of Project 5 with 'user' role
```

This means:
- TaskManager app can only access Project 5's data
- Users via TaskManager cannot access dashboard by default (role='user')
- Admins can upgrade TaskManager users to 'member' to grant dashboard access

---

## Auto-Registration Behavior

### Scenario 1: First-Time User via Dashboard OAuth

```
1. User clicks "Login with Google" on dashboard
2. Dashboard redirects to /oauth/authorize
3. User authenticates with Google
4. System checks: User exists? No
5. System creates:
   - altalune_users (email, first_name, last_name)
   - altalune_user_identities (provider='google', oauth_client_id=dashboard_client_id)
   - altalune_project_members (project_id=dashboard_project_id, role='user')
6. User is now authenticated but has no dashboard access (role='user')
7. Admin must upgrade user to 'member' or 'admin' for dashboard access
```

### Scenario 2: First-Time User via External App OAuth

```
1. User clicks "Login with Google" on ExternalApp
2. ExternalApp redirects to Altalune /oauth/authorize with client_id=external_client_id
3. User authenticates with Google
4. System checks: User exists? No
5. System creates:
   - altalune_users (email, first_name, last_name)
   - altalune_user_identities (provider='google', oauth_client_id=external_client_id)
   - altalune_project_members (project_id=external_app_project_id, role='user')
6. User can now use ExternalApp APIs (if scopes permit)
7. User has NO access to dashboard or other projects
```

### Scenario 3: Existing User Logs Into Different App

```
Given:
- User already exists (signed up via Dashboard)
- User is member of Project 1 (dashboard)

When:
- User logs into ExternalApp (Project 5) for the first time

Then:
1. System checks: User exists? Yes
2. System checks: User identity for this client exists? No
3. System creates:
   - altalune_user_identities (provider='google', oauth_client_id=external_client_id)
   - altalune_project_members (project_id=5, role='user')
4. User now has access to BOTH projects:
   - Project 1: member (dashboard access)
   - Project 5: user (API access only)
```

### Scenario 4: New Project Creation

```
When: Admin/Owner creates a new project

Then:
1. Project record is created in altalune_projects
2. Partitions are auto-created for partitioned tables
3. Superadmin (user_id=1) is AUTOMATICALLY added as owner:
   - altalune_project_members (project_id=new_project_id, user_id=1, role='owner')
4. This happens in internal/domain/project/repo.go:380-384
```

**Implementation**:
```go
// Auto-register superadmin as owner to the new project
if err := r.registerSuperadminAsOwner(ctx, result.ID); err != nil {
    // Log error but don't fail the project creation
    fmt.Printf("Warning: failed to register superadmin to project %d: %v\n", result.ID, err)
}
```

---

## Permission Matrix

### Dashboard Access

| Role | Dashboard Login | View Projects | Settings | OAuth Clients | API Keys | Manage Members |
|------|----------------|---------------|----------|---------------|----------|----------------|
| Owner | ✅ | All | ✅ Edit All | ✅ Full CRUD + Secrets | ✅ | ✅ All roles |
| Admin | ✅ | Assigned | ✅ Edit (no delete) | ✅ Edit (no delete) + Secrets | ✅ | ✅ (except owner) |
| Member | ✅ | Assigned | 👁️ View Only | 👁️ View (no secrets) | ❌ | ❌ |
| User | ❌ | None | ❌ | ❌ | ❌ | ❌ |

### Project Operations

| Role | Create Project | Delete Project | Edit Project | View Project Data | Modify Project Data |
|------|---------------|----------------|--------------|-------------------|---------------------|
| Owner | ✅ | ✅ | ✅ | ✅ All | ✅ All |
| Admin | ❌ | ❌ | ✅ | ✅ Assigned | ✅ Assigned |
| Member | ❌ | ❌ | ❌ | ✅ Assigned | ✅ Assigned |
| User | ❌ | ❌ | ❌ | ❌ | ❌ |

### OAuth Client Operations

| Role | Create Client | Delete Client | Edit Client | View Client Secret | Reveal Secret |
|------|--------------|---------------|-------------|-------------------|---------------|
| Owner | ✅ | ✅ | ✅ | ✅ (on creation) | ✅ (anytime) |
| Admin | ❌ | ❌ | ✅ | ✅ (on creation) | ✅ (anytime) |
| Member | ❌ | ❌ | ❌ | ❌ | ❌ |
| User | ❌ | ❌ | ❌ | ❌ | ❌ |

**Note**: Default dashboard client cannot be deleted by anyone (enforced in UI and backend).

### Member Management

| Role | Add Member | Remove Member | Change Role to Owner | Change Role to Admin | Change Role to Member/User |
|------|-----------|---------------|---------------------|---------------------|---------------------------|
| Owner | ✅ | ✅ (except self) | ❌ (reserved) | ✅ | ✅ |
| Admin | ✅ | ✅ (member/user only) | ❌ | ❌ | ✅ |
| Member | ❌ | ❌ | ❌ | ❌ | ❌ |
| User | ❌ | ❌ | ❌ | ❌ | ❌ |

---

## Database Schema

### altalune_project_members (Global Table)

```sql
CREATE TABLE altalune_project_members (
  id BIGSERIAL PRIMARY KEY,
  public_id VARCHAR(14) NOT NULL UNIQUE,
  project_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  role VARCHAR(20) NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'user')),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (project_id) REFERENCES altalune_projects(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES altalune_users(id) ON DELETE CASCADE,
  UNIQUE(project_id, user_id)
);

CREATE INDEX project_members_project_id_idx ON altalune_project_members(project_id);
CREATE INDEX project_members_user_id_idx ON altalune_project_members(user_id);
CREATE INDEX project_members_role_idx ON altalune_project_members(role);
```

**Key Points**:
- **Global table** (not partitioned) for cross-project member queries
- **Unique constraint** on (project_id, user_id) prevents duplicate memberships
- Users can have multiple rows (one per project) for multi-project access
- CASCADE deletes when project or user is deleted

### altalune_user_identities

```sql
CREATE TABLE altalune_user_identities (
  id BIGSERIAL PRIMARY KEY,
  public_id VARCHAR(14) NOT NULL UNIQUE,
  user_id BIGINT NOT NULL,
  provider VARCHAR(50) NOT NULL,  -- 'google', 'github', 'system'
  provider_user_id VARCHAR(255) NOT NULL,
  oauth_client_id UUID,  -- Links to specific OAuth client
  email VARCHAR(255) NOT NULL,
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  avatar_url TEXT,
  last_login_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (user_id) REFERENCES altalune_users(id) ON DELETE CASCADE,
  UNIQUE(provider, provider_user_id)
);
```

**Key Points**:
- **oauth_client_id**: Links identity to the OAuth client used for signup
- Users can have multiple identities (Google, GitHub, etc.)
- **last_login_at**: Tracks authentication activity

### oauth_clients (Partitioned by project_id)

```sql
CREATE TABLE oauth_clients (
  id BIGSERIAL,
  project_id BIGINT NOT NULL,
  public_id VARCHAR(14) NOT NULL,
  name VARCHAR(100) NOT NULL,
  client_id UUID UNIQUE NOT NULL,
  client_secret_hash VARCHAR(255) NOT NULL,  -- bcrypt hash
  redirect_uris TEXT[] NOT NULL,
  pkce_required BOOLEAN DEFAULT false,
  is_default BOOLEAN DEFAULT false,  -- Dashboard client flag
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  PRIMARY KEY (project_id, id),
  FOREIGN KEY (project_id) REFERENCES altalune_projects(id) ON DELETE CASCADE
) PARTITION BY LIST (project_id);
```

**Key Points**:
- **Partitioned by project_id** for data isolation and performance
- **is_default**: Marks the dashboard OAuth client (cannot be deleted)
- **pkce_required**: Optional PKCE enforcement per client
- **redirect_uris**: Array allows multiple callback URLs

---

## Implementation Guidelines

### Backend: Creating Project Membership

```go
// internal/domain/project_member/repo.go (example)
func (r *Repo) Create(ctx context.Context, input CreateProjectMemberInput) (*ProjectMember, error) {
    // Validate role
    if input.Role == "owner" {
        return nil, ErrOwnerRoleReserved
    }

    // Check if membership already exists
    exists, err := r.membershipExists(ctx, input.ProjectID, input.UserID)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, ErrMembershipAlreadyExists
    }

    // Generate public_id
    publicID, err := nanoid.GeneratePublicID()
    if err != nil {
        return nil, err
    }

    // Insert membership
    query := `
        INSERT INTO altalune_project_members
        (public_id, project_id, user_id, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    // ... execute query
}
```

### Backend: Enforcing Role Permissions

```go
// Middleware example
func RequireRole(minRole string) middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userID := getUserIDFromToken(r)
            projectID := getProjectIDFromRequest(r)

            member, err := projectMemberRepo.GetByProjectAndUser(r.Context(), projectID, userID)
            if err != nil {
                http.Error(w, "Unauthorized", 403)
                return
            }

            if !hasMinimumRole(member.Role, minRole) {
                http.Error(w, "Insufficient permissions", 403)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
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

### Frontend: Role-Based UI

```vue
<template>
  <div>
    <!-- Only show to admin+ -->
    <Button v-if="canManageMembers" @click="openAddMemberDialog">
      Add Member
    </Button>

    <!-- Only show to owner -->
    <Button v-if="isOwner" @click="deleteProject">
      Delete Project
    </Button>

    <!-- Show to all roles with different actions -->
    <DataTable :columns="columns" :data="members">
      <template #actions="{ row }">
        <Button v-if="canEditRole(row)" @click="editRole(row)">
          Edit Role
        </Button>
        <Button v-if="canRemoveMember(row)" @click="removeMember(row)">
          Remove
        </Button>
      </template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
const currentUser = useCurrentUser()
const currentProject = useCurrentProject()

const isOwner = computed(() => currentUser.value?.role === 'owner')
const isAdmin = computed(() => ['owner', 'admin'].includes(currentUser.value?.role))
const canManageMembers = computed(() => isAdmin.value)

const canEditRole = (member: ProjectMember) => {
  if (member.role === 'owner') return false // Cannot edit owner
  if (isOwner.value) return true
  if (isAdmin.value && ['member', 'user'].includes(member.role)) return true
  return false
}

const canRemoveMember = (member: ProjectMember) => {
  if (member.role === 'owner') return false
  if (member.id === currentUser.value?.id) return false // Cannot remove self
  if (isOwner.value) return true
  if (isAdmin.value && ['member', 'user'].includes(member.role)) return true
  return false
}
</script>
```

### Querying User's Project Access

```go
// Get all projects a user has access to
func (r *Repo) GetProjectsByUser(ctx context.Context, userID int64) ([]*UserProject, error) {
    query := `
        SELECT
            p.id,
            p.public_id,
            p.name,
            pm.role,
            pm.created_at as joined_at
        FROM altalune_projects p
        INNER JOIN altalune_project_members pm ON p.id = pm.project_id
        WHERE pm.user_id = $1
        ORDER BY pm.created_at DESC
    `
    // ... execute and return
}

// Check if user has specific role in project
func (r *Repo) HasRole(ctx context.Context, userID, projectID int64, minRole string) (bool, error) {
    var role string
    err := r.db.QueryRowContext(ctx, `
        SELECT role FROM altalune_project_members
        WHERE user_id = $1 AND project_id = $2
    `, userID, projectID).Scan(&role)

    if err == sql.ErrNoRows {
        return false, nil
    }
    if err != nil {
        return false, err
    }

    return hasMinimumRole(role, minRole), nil
}
```

---

## Common Patterns & Best Practices

### Pattern 1: Check Project Access Before Data Operations

```go
func (s *EmployeeService) GetEmployee(ctx context.Context, projectID, employeeID int64) (*Employee, error) {
    // 1. Check user has access to project
    userID := getUserIDFromContext(ctx)
    member, err := s.projectMemberRepo.GetByProjectAndUser(ctx, projectID, userID)
    if err != nil {
        return nil, ErrProjectAccessDenied
    }

    // 2. Check minimum role (member+ can view)
    if !hasMinimumRole(member.Role, "member") {
        return nil, ErrInsufficientPermissions
    }

    // 3. Fetch data (already scoped to project due to partitioning)
    return s.employeeRepo.GetByID(ctx, projectID, employeeID)
}
```

### Pattern 2: Auto-Create Project Membership on OAuth Signup

```go
// internal/domain/oauth_auth/service.go (conceptual)
func (s *AuthService) HandleOAuthCallback(ctx context.Context, code string, clientID uuid.UUID) error {
    // 1. Exchange code for user info from provider
    userInfo, err := s.exchangeCode(ctx, code)

    // 2. Get OAuth client to find project_id
    client, err := s.oauthClientRepo.GetByClientID(ctx, clientID)
    if err != nil {
        return err
    }

    // 3. Check if user exists
    user, err := s.userRepo.GetByEmail(ctx, userInfo.Email)
    if err == ErrUserNotFound {
        // 3a. Create new user
        user, err = s.userRepo.Create(ctx, CreateUserInput{
            Email:     userInfo.Email,
            FirstName: userInfo.FirstName,
            LastName:  userInfo.LastName,
        })

        // 3b. Create user identity
        _, err = s.userIdentityRepo.Create(ctx, CreateUserIdentityInput{
            UserID:         user.ID,
            Provider:       userInfo.Provider,
            ProviderUserID: userInfo.ProviderUserID,
            OAuthClientID:  &clientID,
            Email:          userInfo.Email,
        })

        // 3c. AUTO-CREATE PROJECT MEMBERSHIP with 'user' role
        _, err = s.projectMemberRepo.Create(ctx, CreateProjectMemberInput{
            ProjectID: client.ProjectID,
            UserID:    user.ID,
            Role:      "user",
        })
    }

    // 4. Continue with session/token creation...
}
```

### Pattern 3: Superadmin Auto-Register on Project Creation

Already implemented in `internal/domain/project/repo.go`:

```go
// Called automatically after project creation
func (r *Repo) Create(ctx context.Context, input CreateProjectInput) (*Project, error) {
    // ... create project logic

    // Auto-register superadmin as owner to the new project
    if err := r.registerSuperadminAsOwner(ctx, result.ID); err != nil {
        // Log error but don't fail the project creation
        fmt.Printf("Warning: failed to register superadmin to project %d: %v\n", result.ID, err)
    }

    return result, nil
}

func (r *Repo) registerSuperadminAsOwner(ctx context.Context, projectID int64) error {
    const superadminID = int64(1)  // Fixed user_id from SQL migration

    // Check superadmin exists
    var exists bool
    err := r.db.QueryRowContext(ctx, `
        SELECT EXISTS(SELECT 1 FROM altalune_users WHERE id = $1)
    `, superadminID).Scan(&exists)

    if !exists {
        return nil  // Skip if migrations haven't run yet
    }

    // Check if membership already exists (idempotent)
    err = r.db.QueryRowContext(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM altalune_project_members
            WHERE project_id = $1 AND user_id = $2
        )
    `, projectID, superadminID).Scan(&exists)

    if exists {
        return nil  // Already registered
    }

    // Create owner membership
    publicID, _ := nanoid.GeneratePublicID()
    _, err = r.db.ExecContext(ctx, `
        INSERT INTO altalune_project_members (
            public_id, project_id, user_id, role, created_at, updated_at
        ) VALUES ($1, $2, $3, 'owner', NOW(), NOW())
    `, publicID, projectID, superadminID)

    return err
}
```

---

## Testing Scenarios

### Scenario: User Role Upgrade

```
Given:
- User "john@example.com" (user_id=100)
- Signed up via External App OAuth client
- project_members: (project_id=5, user_id=100, role='user')

When:
- Admin upgrades user to 'member' role

Then:
- User can now access dashboard
- User can view project 5's data
- User still cannot modify settings or manage members

SQL:
UPDATE altalune_project_members
SET role = 'member', updated_at = NOW()
WHERE project_id = 5 AND user_id = 100;
```

### Scenario: Multi-Project Assignment

```
Given:
- User "jane@example.com" (user_id=200)
- Currently member of project 1

When:
- Admin assigns user to project 2 as 'admin'

Then:
- New project_members row created: (project_id=2, user_id=200, role='admin')
- User can now switch between projects in dashboard
- User has 'member' permissions in project 1
- User has 'admin' permissions in project 2

SQL:
INSERT INTO altalune_project_members (public_id, project_id, user_id, role)
VALUES ('abc123def45678', 2, 200, 'admin');
```

### Scenario: Owner Role Restriction

```
Given:
- Admin tries to assign 'owner' role to user

Then:
- Backend returns error: ErrOwnerRoleReserved
- Only superadmin (user_id=1) can have owner role
- UI should not show 'owner' option in role dropdown

Validation:
if input.Role == "owner" && input.UserID != 1 {
    return ErrOwnerRoleReserved
}
```

---

## FAQ

### Q: Can a user belong to multiple projects?

**A**: Yes, for roles **admin** and **member**. Users with 'user' role typically belong to one project (the OAuth client's project), but can be assigned to others if upgraded.

### Q: Can I manually assign the 'owner' role?

**A**: No. The 'owner' role is reserved exclusively for the superadmin user (user_id=1). This is enforced in code and should be validated in the UI.

### Q: What happens when a user logs in via a different OAuth client?

**A**: If the user already exists (same email):
1. New `user_identity` is created linking to the new OAuth client
2. New `project_members` entry is created for the new client's project with 'user' role
3. User now has access to multiple projects (their original project + new client's project)

### Q: How do I give dashboard access to an OAuth user?

**A**: Upgrade their role from 'user' to 'member' or higher. Only 'member', 'admin', and 'owner' roles can access the dashboard.

### Q: Can I delete the default dashboard OAuth client?

**A**: No. The client marked with `is_default = true` cannot be deleted. This is enforced in both the backend service logic and frontend UI.

### Q: What's the difference between 'admin' and 'owner'?

**A**:
- **Owner**: System-wide superadmin, can delete projects, automatically assigned to all projects, only user_id=1
- **Admin**: Project-level administrator, can manage settings and members but cannot delete projects, can be assigned to any user

### Q: How are permissions checked in the API?

**A**: Typically via middleware that:
1. Extracts user_id from JWT token
2. Extracts project_id from request (path param, body, or query)
3. Queries `altalune_project_members` to get user's role in that project
4. Validates role meets minimum requirement for the endpoint

### Q: Can superadmin access all projects without being a member?

**A**: No. While superadmin has the 'owner' role, they must still have a `project_members` entry for each project. However, the auto-registration feature ensures they're automatically added as owner to all new projects.

---

## Related Documentation

- **[BACKEND_GUIDE.md](./BACKEND_GUIDE.md)** - Backend domain implementation patterns
- **[FRONTEND_GUIDE.md](./FRONTEND_GUIDE.md)** - Frontend role-based UI patterns
- **[CLAUDE.md](../../CLAUDE.md)** - Project overview and partitioned tables
- **[oauth_server_prepare/plan.md](../../oauth_server_prepare/plan.md)** - Complete OAuth implementation plan

---

## Version History

- **2026-01-09**: Initial documentation created based on T18 OAuth foundation implementation
- Covers: Role hierarchy, project membership rules, OAuth client-project relationship, auto-registration behavior
