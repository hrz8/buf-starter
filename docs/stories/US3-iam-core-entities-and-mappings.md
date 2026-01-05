# User Story US3: Identity & Access Management - Core Entities and Mappings

## Story Overview

**As a** system administrator managing the Altalune platform
**I want** to manage users, roles, permissions, and their relationships through a comprehensive IAM dashboard
**So that** I can control access to system resources, assign appropriate permissions, and manage user memberships across projects with a flexible role-based access control system

## Acceptance Criteria

### Core CRUD Operations

#### Create User

- **Given** I am on the IAM Users management page
- **When** I click "Create User"
- **Then** I should see a form to create a new user
- **And** I can provide email, name, and avatar URL
- **And** upon successful creation, a new user is created with OAuth identity placeholder
- **And** the user is assigned a unique public_id (14-character nanoid)
- **And** the user has default status of "active"
- **And** no password is stored (OAuth-only authentication)

#### List/Query Users

- **Given** I am on the IAM Users management page
- **When** the page loads
- **Then** I should see a table of all users in the system
- **And** I can see user email, name, status, creation date, and last updated date
- **And** I can search/filter users by email or name (keyword search)
- **And** I can filter by status (active, inactive)
- **And** I can sort by email, name, status, or creation date
- **And** I can see pagination controls when there are many users

#### View User Details

- **Given** I have users in the system
- **When** I click on a user row or view action
- **Then** I should see detailed information about that user
- **And** I can see all assigned roles with their names
- **And** I can see all direct permissions assigned to the user
- **And** I can see all projects the user has access to with their roles

#### Update User

- **Given** I am viewing a user
- **When** I click "Edit" action
- **Then** I should see a tabbed form with Profile, Roles, Permissions, and Projects tabs
- **And** in the Profile tab, I can update email, name, and avatar URL
- **And** in the Roles tab, I can assign/remove roles using dual-list transfer component
- **And** in the Permissions tab, I can assign/remove direct permissions using dual-list transfer component
- **And** in the Projects tab, I can add/remove user from projects with specific project roles
- **And** role and permission mappings are saved immediately upon assignment/removal
- **And** profile changes are saved when I click Update button

#### Delete User

- **Given** I am viewing a user
- **When** I click "Delete" action
- **Then** I should see a confirmation dialog
- **And** when I confirm, the user is permanently deleted
- **And** all user's role assignments are removed (cascade delete)
- **And** all user's permission assignments are removed (cascade delete)
- **And** all user's project memberships are removed (cascade delete)

#### Activate/Deactivate User

- **Given** I am viewing a user
- **When** I click "Activate" or "Deactivate" action
- **Then** the user's status changes accordingly
- **And** inactive users cannot access the system (enforced in future auth story)
- **And** the status change is reflected immediately in the UI

#### Create Role

- **Given** I am on the IAM Roles management page
- **When** I click "Create Role"
- **Then** I should see a form to create a new role
- **And** I can provide role name (alphanumeric and underscores only) and description
- **And** upon successful creation, a new role is created with unique public_id
- **And** role name must be unique system-wide (case-insensitive)

#### List/Query Roles

- **Given** I am on the IAM Roles management page
- **When** the page loads
- **Then** I should see a table of all roles in the system
- **And** I can see role name, description, creation date, and last updated date
- **And** I can search/filter roles by name or description (keyword search)
- **And** I can sort by name or creation date
- **And** I can see the count of users assigned to each role
- **And** I can see the count of permissions assigned to each role

#### Update Role

- **Given** I am viewing a role
- **When** I click "Edit" action
- **Then** I should see a tabbed form with Details and Permissions tabs
- **And** in the Details tab, I can update role name and description
- **And** in the Permissions tab, I can assign/remove permissions using dual-list transfer component
- **And** permission mappings are saved immediately upon assignment/removal
- **And** role detail changes are saved when I click Update button

#### Delete Role

- **Given** I am viewing a role
- **When** I click "Delete" action
- **Then** I should see a confirmation dialog
- **And** I should see a warning if the role has users assigned
- **And** when I confirm, the role is permanently deleted
- **And** all role-permission mappings are removed (cascade delete)
- **And** all user-role assignments are removed (cascade delete)

#### Create Permission

- **Given** I am on the IAM Permissions management page
- **When** I click "Create Permission"
- **Then** I should see a form to create a new permission
- **And** I can provide permission name (alphanumeric and underscores only)
- **And** I can select effect (allow or deny)
- **And** I can provide an optional description
- **And** upon successful creation, a new permission is created with unique public_id
- **And** permission name must be unique system-wide (case-insensitive)

#### List/Query Permissions

- **Given** I am on the IAM Permissions management page
- **When** the page loads
- **Then** I should see a table of all permissions in the system
- **And** I can see permission name, effect, description, creation date
- **And** I can search/filter permissions by name or description (keyword search)
- **And** I can filter by effect (allow, deny)
- **And** I can sort by name, effect, or creation date
- **And** I can see the count of roles using each permission
- **And** I can see the count of users with direct assignment of each permission

#### Update Permission

- **Given** I am viewing a permission
- **When** I click "Edit" action
- **Then** I should see a form to update the permission
- **And** I can update permission name, effect, and description
- **And** changes are saved when I click Update button

#### Delete Permission

- **Given** I am viewing a permission
- **When** I click "Delete" action
- **Then** I should see a confirmation dialog
- **And** I should see a warning if the permission is assigned to roles or users
- **And** when I confirm, the permission is permanently deleted
- **And** all role-permission mappings are removed (cascade delete)
- **And** all user-permission mappings are removed (cascade delete)

### Relationship Mapping Operations

#### Assign Roles to User (Dual-List Transfer)

- **Given** I am editing a user in the Roles tab
- **When** I view the dual-list transfer component
- **Then** I should see "Available Roles" list on the left showing all unassigned roles
- **And** I should see "Assigned Roles" list on the right showing currently assigned roles
- **And** I can search/filter roles in both lists by name
- **And** I can select multiple roles in the Available list using checkboxes
- **And** when I click the right arrow (→) button, selected roles are assigned immediately via API
- **And** assigned roles move from Available to Assigned list (optimistic update)
- **And** I see a success toast notification
- **And** if assignment fails, I see an error toast and roles remain in Available list

#### Remove Roles from User

- **Given** I am editing a user in the Roles tab
- **When** I select roles in the "Assigned Roles" list
- **And** I click the left arrow (←) button
- **Then** selected roles are removed immediately via API
- **And** removed roles move from Assigned to Available list (optimistic update)
- **And** I see a success toast notification
- **And** if removal fails, I see an error toast and roles remain in Assigned list

#### Assign Permissions to Role (Dual-List Transfer)

- **Given** I am editing a role in the Permissions tab
- **When** I view the dual-list transfer component
- **Then** I should see "Available Permissions" and "Assigned Permissions" lists
- **And** I can search/filter permissions in both lists
- **And** I can see permission effect (allow/deny) as visual indicator in the list
- **And** I can select and assign multiple permissions at once
- **And** assignment is immediate with optimistic UI updates

#### Assign Direct Permissions to User

- **Given** I am editing a user in the Permissions tab
- **When** I view the dual-list transfer component
- **Then** I should see permissions grouped or filtered by effect (allow/deny)
- **And** I can assign direct permissions that override or supplement role permissions
- **And** assignment is immediate with optimistic UI updates
- **And** I can see which permissions come from roles vs direct assignment (in view mode)

#### Add User to Project (Project Members Management)

- **Given** I am editing a user in the Projects tab
- **When** I click "Add to Project"
- **Then** I should see a dialog or inline form
- **And** I can select a project from dropdown
- **And** I can select a project role (owner, admin, member, viewer)
- **And** upon confirmation, user is added to project with selected role
- **And** the project appears in the user's Projects list
- **And** one user can only have one role per project (enforced by unique constraint)

#### Remove User from Project

- **Given** I am viewing user's projects in the Projects tab
- **When** I click remove/delete icon on a project
- **Then** I should see a confirmation dialog
- **And** when confirmed, user is removed from that project
- **And** the project disappears from the user's Projects list

#### Update User's Project Role

- **Given** I am viewing user's projects in the Projects tab
- **When** I change the role dropdown for a project
- **Then** the user's role for that project is updated immediately
- **And** I see a success notification

#### View User Mappings (Read-Only)

- **Given** I am on the Users table
- **When** I click "View Mappings" action on a user row
- **Then** a sheet opens showing read-only view of user's mappings
- **And** I see three tabs: Roles, Permissions, Projects
- **And** in Roles tab, I see list of all assigned roles with role names and descriptions
- **And** in Permissions tab, I see both role-inherited and direct permissions with effect indicators
- **And** in Projects tab, I see all projects with user's role in each project
- **And** I see creation/join dates as metadata
- **And** I can close the sheet without making any changes

### Project Members Management (From Project Settings)

#### View Project Members

- **Given** I am in Project Settings → Members tab
- **When** the page loads
- **Then** I should see a table of all members with access to this project
- **And** I can see member email, name, role, and join date
- **And** I can filter by role (owner, admin, member, viewer)
- **And** I can search by member name or email

#### Add Member to Project (From Project Side)

- **Given** I am in Project Settings → Members tab
- **When** I click "Add Member"
- **Then** I should see a dialog with user selector and role dropdown
- **And** I can search and select an existing user from the system
- **And** I can select the project role for this user
- **And** upon confirmation, user is added as project member
- **And** the member appears in the project members table

#### Change Member Role in Project

- **Given** I am viewing project members
- **When** I change the role dropdown for a member
- **Then** the member's project role is updated immediately
- **And** I see a success notification

#### Remove Member from Project (From Project Side)

- **Given** I am viewing project members
- **When** I click remove/delete action on a member
- **Then** I should see a confirmation dialog
- **And** when confirmed, the member is removed from the project
- **And** the member disappears from the project members table

### Security Requirements

#### User Email Uniqueness

- Email addresses must be unique across the system (case-insensitive)
- Attempting to create a user with duplicate email shows validation error
- Email validation follows RFC 5322 format

#### Role and Permission Name Validation

- Role names can only contain alphanumeric characters and underscores (^[a-zA-Z0-9_]+$)
- Permission names follow same pattern as role names
- Names are unique system-wide (case-insensitive comparison)
- Reserved role name: "super_admin" (seeded by migration)
- Reserved permission name: "root" (seeded by migration)

#### Data Isolation and Access Control

- Users, roles, and permissions are global (system-wide), not project-scoped
- Project access is controlled via separate project_members table
- Super admin role users can access all projects without explicit membership (enforced at application level)
- Cascade deletes ensure referential integrity when entities are removed

#### OAuth-Only Authentication Model

- No password fields exist in users table
- Users must be linked to OAuth provider via user_identities (in future OAuth flow story)
- Mock super admin user created with placeholder Google OAuth identity
- Actual OAuth flow implementation is out of scope for this story

### Data Validation

#### User Fields

- **email**: Required, valid email format, max 255 characters, unique (case-insensitive)
- **name**: Required, 1-100 characters
- **avatar_url**: Optional, valid URL format, max 500 characters
- **status**: Required, enum ('active', 'inactive'), default 'active'
- **public_id**: System-generated, 14-character nanoid, unique

#### Role Fields

- **name**: Required, 2-50 characters, alphanumeric + underscore, unique (case-insensitive)
- **description**: Optional, max 500 characters
- **public_id**: System-generated, 14-character nanoid, unique

#### Permission Fields

- **name**: Required, 2-100 characters, alphanumeric + underscore, unique (case-insensitive)
- **effect**: Required, enum ('allow', 'deny'), default 'allow'
- **description**: Optional, max 500 characters
- **public_id**: System-generated, 14-character nanoid, unique

#### Project Member Fields

- **project_id**: Required, foreign key to altalune_projects
- **user_id**: Required, foreign key to altalune_users
- **role**: Required, enum ('owner', 'admin', 'member', 'viewer'), default 'viewer'
- **Unique constraint**: One user can only have one role per project

#### Mapping Table Constraints

- **users_roles**: Composite unique constraint on (user_id, role_id)
- **roles_permissions**: Composite unique constraint on (role_id, permission_id)
- **users_permissions**: Composite unique constraint on (user_id, permission_id)
- **project_members**: Composite unique constraint on (project_id, user_id)
- All foreign keys cascade on delete to maintain referential integrity

### User Experience

#### Responsive Design

- Works on desktop and mobile devices
- Tables should be scrollable/responsive on small screens
- Dual-list transfer component adapts to screen size
- Forms should be touch-friendly on mobile
- Tabs in edit forms stack vertically on mobile

#### Feedback and Notifications

- Success toast messages for create/update/delete/activate/deactivate operations
- Success toast messages for immediate mapping operations (assign/remove)
- Clear error messages for validation failures
- Loading states during API operations
- Loading skeletons while fetching user/role/permission data
- Confirmation dialogs for destructive actions (delete user, delete role with assignments)
- Warning dialogs when deleting entities with relationships

#### Integration with Existing UI

- Follows existing design patterns and components (shadcn-vue)
- Uses established DataTable component with pagination/filtering/sorting
- Uses Sheet pattern for edit forms (UserUpdateSheet + UserUpdateForm)
- Uses Dialog for delete confirmations
- Follows established form patterns with vee-validate
- Integrates with existing navigation structure (top-level IAM menu)
- Uses Breadcrumb pattern for navigation context
- Consistent with API Keys feature UI/UX

#### Dual-List Transfer Component Features

- Two side-by-side lists: Available (left) and Assigned (right)
- Search/filter input at top of each list
- Checkbox multi-select in both lists
- Arrow buttons between lists: → (assign), ← (remove)
- Badge showing item count in each list
- Visual disabled state when loading or no items selected
- Scrollable list areas with fixed height
- Empty state messages when lists are empty
- Loading skeletons during data fetch

#### Tab-Based Edit Form UX

- Profile tab for basic user information (email, name, avatar)
- Roles tab with dual-list transfer for role assignment
- Permissions tab with dual-list transfer for direct permission assignment
- Projects tab with project membership management
- Badge on each tab showing count of assigned items
- Clear visual separation between tabs
- Tab state persists within edit session
- Form actions (Cancel, Update) always visible at bottom

## Technical Requirements

### Backend Architecture

#### Domain Structure (7-File Pattern)

**User Domain** (`internal/domain/user/`):
- `model.go` - User structs, status enum, input/result types
- `interface.go` - Repository interface with CRUD and query methods
- `repo.go` - PostgreSQL implementation with pgx, no project_id filtering (global table)
- `service.go` - Business logic with protovalidate, email uniqueness checks
- `handler.go` - Connect-RPC handler (thin wrapper around service)
- `mapper.go` - Proto ↔ Domain conversions, status enum mappings
- `errors.go` - User-specific errors (UserNotFound, UserAlreadyExists, etc.)

**Role Domain** (`internal/domain/role/`):
- Same 7-file structure as User domain
- Name validation with regex pattern in repo and service
- No status field (roles don't have active/inactive)

**Permission Domain** (`internal/domain/permission/`):
- Same 7-file structure as User and Role domains
- Effect enum mappings (allow, deny)
- Name validation with regex pattern

**IAM Mapper Domain** (`internal/domain/iam_mapper/`):
- `model.go` - Mapping result structs, project role enum
- `interface.go` - Mapping operations interface (assign, remove, get mappings)
- `repo.go` - Batch insert/delete for mappings, efficient JOIN queries
- `service.go` - Mapping business logic, validation
- `handler.go` - Connect-RPC handler for mapping operations
- `mapper.go` - Proto ↔ Domain conversions for mapping results
- `errors.go` - Mapping-specific errors

#### Database Design

**Global Tables (Non-Partitioned):**
- `altalune_users` - No project_id column, simple primary key
- `altalune_roles` - No project_id column
- `altalune_permissions` - No project_id column
- `altalune_users_roles` - Junction table with surrogate key
- `altalune_roles_permissions` - Junction table with surrogate key
- `altalune_users_permissions` - Junction table with surrogate key
- `altalune_project_members` - Project access mapping (global, not partitioned)
- `altalune_user_identities` - OAuth identity links (for future use)

**Key Differences from Partitioned Tables:**
- No composite primary keys (just `id`, not `(project_id, id)`)
- Simpler unique indexes (no partition key in index)
- No partition routing in queries
- Standard foreign key constraints without partition key

**Indexes Strategy:**
- Unique index on public_id for all entity tables
- Unique index on email (users), name (roles, permissions)
- Composite unique indexes on junction tables for duplicate prevention
- Individual indexes on foreign keys for efficient joins
- Status/effect indexes for filtering

#### Error Handling and Codes

- **605XX**: User domain errors
  - 60501: UserNotFound
  - 60502: UserAlreadyExists (email conflict)
  - 60503: InvalidUserStatus
  - 60504: UserCannotBeDeleted
- **606XX**: Role domain errors
  - 60601: RoleNotFound
  - 60602: RoleAlreadyExists (name conflict)
  - 60603: RoleCannotBeDeleted (has users)
- **607XX**: Permission domain errors
  - 60701: PermissionNotFound
  - 60702: PermissionAlreadyExists (name conflict)
  - 60703: PermissionCannotBeDeleted (in use)
- **608XX**: IAM Mapper domain errors
  - 60801: MappingAlreadyExists
  - 60802: MappingNotFound
  - 60803: InvalidProjectRole

#### Repository Patterns

**Non-Partitioned Query Pattern:**
```go
query := `
    SELECT id, public_id, email, name, status, created_at, updated_at
    FROM altalune_users
    WHERE 1=1  -- NO project_id filtering!
`
```

**Email Uniqueness Check (Case-Insensitive):**
```go
query := `SELECT COUNT(*) FROM altalune_users WHERE LOWER(email) = LOWER($1)`
```

**Batch Assignment (ON CONFLICT DO NOTHING):**
```go
INSERT INTO altalune_users_roles (user_id, role_id, created_by)
VALUES ($1, $2, $3), ($4, $5, $6)
ON CONFLICT (user_id, role_id) DO NOTHING
```

**Efficient Mapping Retrieval with JOINs:**
```go
SELECT ur.user_id, ur.role_id, r.name as role_name, ur.created_at
FROM altalune_users_roles ur
INNER JOIN altalune_roles r ON ur.role_id = r.id
WHERE ur.user_id = $1
ORDER BY r.name ASC
```

#### Dual ID System

- **Internal ID**: BIGINT auto-increment for database operations
- **Public ID**: 14-character nanoid for external API exposure
- All proto messages use public_id (string)
- Repository methods resolve public_id → internal id when needed
- mapper.go handles conversion between domain models and proto messages

### Frontend Architecture

#### Repository Layer

**Pattern**: Simple pass-through with ConnectError logging
- `frontend/shared/repository/user.ts`
- `frontend/shared/repository/role.ts`
- `frontend/shared/repository/permission.ts`
- `frontend/shared/repository/iam_mapper.ts`

Each repository wraps Connect-RPC client calls with try-catch error handling.

#### Service Composables

**Pattern**: Reactive state per operation with dual validation
- `frontend/app/composables/services/useUserService.ts`
- `frontend/app/composables/services/useRoleService.ts`
- `frontend/app/composables/services/usePermissionService.ts`
- `frontend/app/composables/services/useIAMMapperService.ts`

Each service provides:
- Reactive loading/error/success states per operation
- Validator composables (useConnectValidator)
- Reset functions for state cleanup
- Computed refs for UI binding

#### Feature Organization

**Pattern**: Centralized schema/error/constants files
```
frontend/app/components/features/iam/
├── user/
│   ├── index.ts (exports)
│   ├── schema.ts (Zod validation - single source of truth)
│   ├── error.ts (ConnectRPC error utilities)
│   ├── constants.ts (status options, etc.)
│   ├── UserTable.vue
│   ├── UserCreateSheet.vue + UserCreateForm.vue
│   ├── UserUpdateSheet.vue + UserUpdateForm.vue
│   ├── UserDeleteDialog.vue
│   ├── UserRowActions.vue
│   └── UserViewMappingsSheet.vue
├── role/ (similar structure)
├── permission/ (similar structure)
└── shared/
    └── TransferList.vue (reusable component)
```

#### TransferList Component Specification

**Component**: `frontend/app/components/ui/transfer-list/TransferList.vue`

**Props:**
```typescript
interface TransferListItem {
  id: string;
  label: string;
  description?: string;
  disabled?: boolean;
}

props: {
  availableItems: TransferListItem[];
  selectedItems: TransferListItem[];
  availableTitle?: string;
  selectedTitle?: string;
  loading?: boolean;
  disabled?: boolean;
  height?: string;
}
```

**Events:**
```typescript
emit('assign', itemIds: string[]);
emit('remove', itemIds: string[]);
```

**Features:**
- Dual-list layout with grid cols
- Search/filter in both lists
- Multi-select with checkboxes
- Arrow buttons with disabled states
- Badge counts for each list
- ScrollArea for scrollable lists
- Loading skeletons
- Empty state messages

#### Form Validation Strategy

**Dual-Layer Validation:**
1. **Primary**: vee-validate with Zod schemas (client-side, immediate feedback)
2. **Fallback**: ConnectRPC protovalidate (server-side, edge cases)

**Critical vee-validate Pattern:**
```vue
<script setup>
// ✅ CRITICAL: Loading starts as TRUE
const isLoading = ref(true);

const formSchema = toTypedSchema(userCreateSchema);
const form = useForm({ validationSchema: formSchema });

onMounted(async () => {
  await fetchData();
  isLoading.value = false; // Set false after loading
});
</script>

<template>
  <div v-if="isLoading">
    <Skeleton class="h-10 w-full" />
  </div>
  <form v-else @submit="onSubmit">
    <!-- ✅ NO :key attribute on FormField -->
    <FormField v-slot="{ componentField }" name="email">
      <FormItem>
        <FormLabel>Email</FormLabel>
        <FormControl>
          <Input v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>
  </form>
</template>
```

#### State Management for Mappings

**Pattern**: Immediate save with optimistic updates

```typescript
async function handleAssignRoles(roleIds: string[]) {
  try {
    // Immediate API call
    await assignRolesToUser({ userId, roleIds });

    // Optimistic update
    const assignedItems = availableRoles.value.filter(r => roleIds.includes(r.id));
    selectedRoles.value.push(...assignedItems);
    availableRoles.value = availableRoles.value.filter(r => !roleIds.includes(r.id));

    toast.success('Roles assigned successfully');
  } catch (error) {
    // Rollback on error
    toast.error('Failed to assign roles');
    throw error;
  }
}
```

#### Navigation Integration

**Top-Level IAM Menu:**
```typescript
{
  title: 'IAM',
  to: '/iam',
  icon: Shield,
  items: [
    { title: 'Users', to: '/iam/users' },
    { title: 'Roles', to: '/iam/roles' },
    { title: 'Permissions', to: '/iam/permissions' },
  ]
}
```

**Breadcrumb Support:**
Each page defines breadcrumb metadata for navigation context.

### API Design

#### Protocol Buffer Services

**UserService:**
- `QueryUsers(QueryUsersRequest) → QueryUsersResponse`
- `CreateUser(CreateUserRequest) → CreateUserResponse`
- `GetUser(GetUserRequest) → GetUserResponse`
- `UpdateUser(UpdateUserRequest) → UpdateUserResponse`
- `DeleteUser(DeleteUserRequest) → DeleteUserResponse`
- `ActivateUser(ActivateUserRequest) → ActivateUserResponse`
- `DeactivateUser(DeactivateUserRequest) → DeactivateUserResponse`

**RoleService, PermissionService:** Similar structure to UserService (no activate/deactivate for roles/permissions)

**IAMMapperService:**
- `AssignRolesToUser(AssignRolesToUserRequest) → AssignRolesToUserResponse`
- `RemoveRolesFromUser(RemoveRolesFromUserRequest) → RemoveRolesFromUserResponse`
- `GetUserRoles(GetUserRolesRequest) → GetUserRolesResponse`
- Similar RPCs for role-permission and user-permission mappings
- `AddProjectMember(AddProjectMemberRequest) → AddProjectMemberResponse`
- `RemoveProjectMember(RemoveProjectMemberRequest) → RemoveProjectMemberResponse`
- `GetProjectMembers(GetProjectMembersRequest) → GetProjectMembersResponse`
- `GetUserProjects(GetUserProjectsRequest) → GetUserProjectsResponse`

#### Request/Response Validation

**buf.validate Rules:**
- All IDs validated as exactly 14 characters (nanoid length)
- Email fields validated with email format
- Name fields with min/max length and regex patterns
- Enum fields with defined_only constraint
- Required fields marked explicitly

**Example:**
```protobuf
message CreateUserRequest {
  string email = 1 [
    (buf.validate.field).required = true,
    (buf.validate.field).string = {
      email: true,
      max_len: 255
    }
  ];
  string name = 2 [
    (buf.validate.field).required = true,
    (buf.validate.field).string = {
      min_len: 1,
      max_len: 100
    }
  ];
}
```

#### Status Codes and Error Handling

- **200 OK**: Successful CRUD operations
- **400 Bad Request**: Validation errors (protovalidate failures)
- **404 Not Found**: Entity not found (user, role, permission)
- **409 Conflict**: Unique constraint violations (duplicate email, duplicate name)
- **500 Internal Server Error**: Database errors, unexpected failures

**Consistent Error Response:**
```protobuf
message ErrorDetail {
  string code = 1;    // Error code (60501, etc.)
  string message = 2; // Human-readable message
  map<string, string> meta = 3; // Additional context
}
```

#### Pagination and Filtering

**Using shared query.QueryParams:**
- `page`: Page number (1-indexed)
- `page_size`: Items per page (default: 10, max: 100)
- `keyword`: Global text search
- `filters`: Column-specific filters (e.g., statuses=['active'])
- `sort_by`: Column name to sort by
- `sort_order`: 'asc' or 'desc'

**Response includes:**
- `data`: Array of entities
- `meta.row_count`: Total count of matching rows
- `meta.page_count`: Total pages
- `meta.filters`: Available filter values (for dropdowns)

## Out of Scope

### Permission Evaluation Logic

- Middleware for checking user permissions is NOT included in this story
- Logic for combining role permissions + direct permissions is deferred
- "Allow" vs "Deny" precedence rules are deferred
- `validateAccess()` or `hasPermission()` functions are NOT implemented
- Integration with actual route/API protection is deferred to separate auth middleware story

### OAuth Flow Implementation

- Actual OAuth2 login flow (redirect, callback, token exchange) is NOT included
- OAuth provider integration (Google, Github SDKs) is NOT implemented
- User creation from OAuth callback is NOT implemented
- Token storage and session management is NOT implemented
- This story only creates the database schema and UI for managing OAuth providers

### Advanced IAM Features

- Permission wildcards (e.g., "project:*") are supported in data model but evaluation logic is deferred
- Hierarchical permissions (parent-child relationships) are NOT implemented
- Time-based permissions (temporary access) are NOT implemented
- IP-based access restrictions are NOT implemented
- Audit logging of permission changes (who changed what when) is deferred to separate audit story
- Permission templates or role templates are NOT implemented
- Bulk user import/export is NOT implemented

### User Management Features

- User profile pages (public-facing) are NOT implemented
- User preferences/settings are NOT implemented
- User notifications are NOT implemented
- User activity tracking/history is NOT implemented
- User search with advanced filters (by role, by permission) is deferred
- User onboarding workflow is NOT implemented

### Project-Level Features

- Project-specific roles (beyond owner/admin/member/viewer) are NOT implemented
- Per-project permission customization is NOT implemented
- Project invitation system is NOT implemented
- Project access requests/approvals are NOT implemented

### API Features

- Rate limiting per user or role is NOT implemented
- API key scoping to specific permissions is NOT implemented (different from project API keys)
- Webhook notifications for IAM changes are NOT implemented

## Dependencies

### Existing Infrastructure

- **Database**: PostgreSQL with pgx driver (existing)
- **Migration system**: Goose migrations (existing)
- **Protobuf tooling**: Buf for code generation (existing)
- **Backend framework**: Connect-RPC for dual HTTP/gRPC (existing)
- **Frontend framework**: Nuxt.js with Vue 3 and TypeScript (existing)
- **UI components**: shadcn-vue component library (existing)
- **Form validation**: vee-validate with Zod (existing)
- **Routing**: Nuxt file-based routing (existing)
- **State management**: Nuxt composables pattern (existing)

### Data Dependencies

- **Projects table**: `altalune_projects` must exist for project_members foreign key
- **NanoID generator**: `internal/shared/nanoid/` utility must exist
- **Query utilities**: `internal/shared/query/` for pagination/filtering must exist
- **Error handling**: Base error system in `errors.go` must exist

### Frontend Dependencies

- **DataTable component**: Existing reusable data table with pagination
- **Sheet component**: Modal sheet for edit forms
- **Dialog component**: Confirmation dialogs
- **Form components**: FormField, FormItem, FormControl, FormLabel, FormMessage
- **Input components**: Input, Textarea, Select, Checkbox
- **Feedback components**: Toast notifications (vue-sonner), Alert
- **Layout components**: Badge, Button, ScrollArea, Separator, Tabs

### Development Dependencies

- **Air**: Hot reload for development (existing)
- **Buf CLI**: Protobuf linting and generation (existing)
- **pnpm**: Frontend package manager (existing)
- **Go toolchain**: 1.21+ for backend development (existing)

## Definition of Done

### Backend Completion

- [ ] Database migrations created and tested (up and down)
- [ ] All 8 tables created with proper constraints and indexes
- [ ] Seed data migration creates super_admin role, root permission, and mock user
- [ ] User domain fully implemented (7 files, all CRUD + activate/deactivate)
- [ ] Role domain fully implemented (7 files, all CRUD)
- [ ] Permission domain fully implemented (7 files, all CRUD)
- [ ] IAM Mapper domain fully implemented (7 files, all mapping operations)
- [ ] All protocol buffer schemas defined with comprehensive validation
- [ ] `buf generate` runs successfully and generates Go and TypeScript code
- [ ] Error codes added to errors.go with constructor functions
- [ ] Container wiring complete (repositories, services, handlers registered)
- [ ] All services registered in grpc_services.go and http_routes.go
- [ ] Repository queries work correctly without project_id (non-partitioned pattern)
- [ ] Email uniqueness enforced (case-insensitive checks)
- [ ] Name validation enforced with regex patterns
- [ ] Cascade deletes work correctly for all relationships

### Frontend Completion

- [ ] Repository layer implemented for all four services
- [ ] Service composables implemented with reactive state management
- [ ] Zod schemas defined for all form validations (schema.ts files)
- [ ] Error utilities created for ConnectRPC error handling (error.ts files)
- [ ] Constants defined for dropdowns and options (constants.ts files)
- [ ] TransferList component fully implemented and tested in isolation
- [ ] User feature complete (Table, Create, Update, Delete, ViewMappings, RowActions)
- [ ] Role feature complete (Table, Create, Update, Delete, RowActions)
- [ ] Permission feature complete (Table, Create, Update, Delete, RowActions)
- [ ] User edit form has 4 tabs (Profile, Roles, Permissions, Projects) working correctly
- [ ] Role edit form has 2 tabs (Details, Permissions) working correctly
- [ ] Dual-list transfer works for all mapping operations
- [ ] Immediate save works with optimistic updates and error rollback
- [ ] View mappings sheet shows read-only data correctly
- [ ] Project members can be managed from both user form and project settings
- [ ] Navigation menu includes top-level IAM with three submenus
- [ ] All three IAM pages created with proper breadcrumb metadata
- [ ] i18n translations added for all IAM features

### Quality Assurance

- [ ] All CRUD operations tested manually for users, roles, permissions
- [ ] All mapping operations tested (assign, remove, view)
- [ ] Search and filtering work correctly in all data tables
- [ ] Sorting works correctly in all data tables
- [ ] Pagination works correctly with accurate counts
- [ ] Form validation shows appropriate error messages
- [ ] Toast notifications appear for all operations
- [ ] Confirmation dialogs appear for destructive actions
- [ ] Loading states display correctly during operations
- [ ] Error states display correctly when operations fail
- [ ] Responsive design tested on mobile and desktop
- [ ] vee-validate forms work without "useFormField" errors
- [ ] TransferList component handles edge cases (empty lists, search no results)

### Code Quality

- [ ] Code follows established patterns (7-file domain, feature organization)
- [ ] No console errors in browser
- [ ] No Go compilation errors
- [ ] Protobuf schemas pass `buf lint`
- [ ] Database migrations are reversible (down migrations work)
- [ ] Repository methods have error handling
- [ ] Service methods validate input and handle edge cases
- [ ] Frontend components handle loading and error states
- [ ] Code is properly formatted (Go: gofmt, Frontend: prettier)

### Documentation

- [ ] CLAUDE.md updated with IAM patterns (if needed)
- [ ] Migration files have clear comments
- [ ] Complex repository queries have explanatory comments
- [ ] Frontend component props and events are documented
- [ ] Schema.ts files have JSDoc comments for validation rules
- [ ] README updated with IAM feature overview (if applicable)

### Deployment Readiness

- [ ] All migrations run successfully on clean database
- [ ] Seed data creates expected super_admin, root permission, mock user
- [ ] Air hot reload works during development
- [ ] Frontend build succeeds (`pnpm build`)
- [ ] Backend build succeeds (`make build`)
- [ ] No breaking changes to existing API endpoints
- [ ] Feature can be deployed without downtime (migrations are additive)

## Notes

### Implementation Sequencing

**Recommended order:**
1. Database migrations + seed data
2. Backend domains (user, role, permission, iam_mapper)
3. Protobuf schemas + buf generate
4. Container wiring and service registration
5. Frontend repositories + service composables
6. TransferList component (test in isolation first)
7. User feature (complete before moving to next)
8. Role feature
9. Permission feature
10. Project members management integration
11. Navigation and i18n
12. Testing and polish

### Potential Challenges

**Challenge 1: vee-validate FormField Errors**
- **Solution**: Follow FRONTEND_GUIDE.md patterns strictly (loading starts as true, no :key attributes)

**Challenge 2: Non-Partitioned Query Patterns**
- **Solution**: Add clear comments in repo.go files, verify queries in testing

**Challenge 3: TransferList Component Complexity**
- **Solution**: Build as standalone component first, test in Storybook or demo page before integration

**Challenge 4: Immediate Save State Management**
- **Solution**: Use optimistic updates with error rollback pattern, clear toast notifications

**Challenge 5: Mapping Data Consistency**
- **Solution**: Always refetch after save to ensure UI matches database, handle concurrent modifications gracefully

### Future Enhancements

- **Audit logging**: Track who made what changes to IAM entities
- **Permission evaluation**: Middleware to check permissions on routes and API calls
- **OAuth flow**: Complete OAuth login flow for Google and Github
- **Advanced permissions**: Wildcards, hierarchical permissions, time-based access
- **Bulk operations**: Import/export users, batch role assignments
- **User profiles**: Public-facing user profile pages
- **Project invitations**: Invite users to projects via email
- **API key permissions**: Scope API keys to specific permissions
- **Role templates**: Pre-defined role configurations for common use cases

### Architecture Benefits

**Global IAM Design:**
- Users, roles, and permissions are system-wide (not per-project)
- Enables cross-project user management
- Single source of truth for identity
- Scalable for enterprise multi-tenant scenarios

**Separate Project Access:**
- Project membership is independent of global permissions
- Clear separation of concerns (system permissions vs project access)
- Simple project role hierarchy (owner/admin/member/viewer)
- Super admin can access all projects without explicit membership

**RBAC + Direct Permissions:**
- Flexible permission model
- Users inherit permissions from roles
- Direct permissions allow exceptions
- Effect field (allow/deny) supports fine-grained control

**Reusable Components:**
- TransferList component works for all mapping scenarios
- Consistent UI patterns across all three entities
- Easy to extend to new mapping types in future

### Cross-Feature Integration Points

**With Project Settings (US2):**
- Project Settings → Members tab shows project_members
- User Edit → Projects tab shows same data from user perspective
- Bidirectional management ensures consistency

**With API Keys (US1):**
- API keys are project-scoped (existing feature)
- IAM users are global (this feature)
- Future: Could scope API keys to user permissions

**With Future OAuth Story (US4):**
- user_identities table prepared in this story
- OAuth provider configuration in US4
- OAuth flow implementation in separate future story

### Relationship with Permission Evaluation

This story creates the **data model and UI** for permissions. The actual **evaluation logic** (checking if user has permission to do X) is intentionally deferred to a separate middleware/authorization story.

**What this story provides:**
- Database tables for users, roles, permissions, and mappings
- UI to assign roles and permissions to users
- Data structure with effect field (allow/deny)

**What is deferred:**
- `validateAccess(user, permission)` function
- Middleware to protect routes/API endpoints
- Logic for combining role permissions + direct permissions
- Precedence rules (deny overrides allow)
- Wildcard permission matching

This separation allows the IAM dashboard to be built and tested independently while authorization logic is developed in parallel or later.
