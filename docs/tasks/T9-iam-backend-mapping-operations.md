# Task T9: IAM Backend Mapping Operations

**Story Reference:** US3-iam-core-entities-and-mappings.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 10-12 hours
**Prerequisites:** T8 (Core entities must be implemented)

## Objective

Implement the IAM Mapper domain to handle all relationship mappings (users↔roles, roles↔permissions, users↔permissions, users↔projects) with batch operations, efficient JOIN queries, and support for inline permission creation.

## Acceptance Criteria

- [ ] Create protobuf schema for IAMMapperService
- [ ] Implement complete IAM Mapper domain (7 files)
- [ ] All mapping operations functional (assign, remove, get)
- [ ] Batch insert/delete with ON CONFLICT handling
- [ ] Efficient JOIN queries for retrieving mappings
- [ ] Project members management (project_members table)
- [ ] Inline permission creation support
- [ ] Error codes defined (608XX range)
- [ ] Container wiring complete
- [ ] Service registered in gRPC and HTTP routes
- [ ] Run `buf generate` successfully
- [ ] All integration tests pass

## Technical Requirements

### 1. Protobuf Schema

#### api/proto/altalune/v1/iam_mapper.proto

**Mapping Request/Response Messages:**

```protobuf
// User-Role Mappings
message AssignUserRolesRequest {
  string user_id = 1 [(buf.validate.field).string.min_len = 14];
  repeated string role_ids = 2;
}

message RemoveUserRolesRequest {
  string user_id = 1 [(buf.validate.field).string.min_len = 14];
  repeated string role_ids = 2;
}

message GetUserRolesRequest {
  string user_id = 1 [(buf.validate.field).string.min_len = 14];
}

message GetUserRolesResponse {
  repeated Role roles = 1;
}

// Role-Permission Mappings
message AssignRolePermissionsRequest {
  string role_id = 1 [(buf.validate.field).string.min_len = 14];
  repeated string permission_ids = 2;
}

message RemoveRolePermissionsRequest {
  string role_id = 1 [(buf.validate.field).string.min_len = 14];
  repeated string permission_ids = 2;
}

message GetRolePermissionsRequest {
  string role_id = 1 [(buf.validate.field).string.min_len = 14];
}

message GetRolePermissionsResponse {
  repeated Permission permissions = 1;
}

// User-Permission Mappings (Direct Assignments)
message AssignUserPermissionsRequest {
  string user_id = 1 [(buf.validate.field).string.min_len = 14];
  repeated string permission_ids = 2;
}

message RemoveUserPermissionsRequest {
  string user_id = 1 [(buf.validate.field).string.min_len = 14];
  repeated string permission_ids = 2;
}

message GetUserPermissionsRequest {
  string user_id = 1 [(buf.validate.field).string.min_len = 14];
}

message GetUserPermissionsResponse {
  repeated Permission permissions = 1;
}

// Project Members
message AssignProjectMembersRequest {
  string project_id = 1 [(buf.validate.field).string.min_len = 14];
  repeated ProjectMember members = 2;
}

message ProjectMember {
  string user_id = 1;
  string role = 2;  // owner, admin, member, viewer
}

message RemoveProjectMembersRequest {
  string project_id = 1 [(buf.validate.field).string.min_len = 14];
  repeated string user_ids = 2;
}

message GetProjectMembersRequest {
  string project_id = 1 [(buf.validate.field).string.min_len = 14];
}

message GetProjectMembersResponse {
  repeated ProjectMemberWithUser members = 1;
}

message ProjectMemberWithUser {
  User user = 1;
  string role = 2;
  google.protobuf.Timestamp created_at = 3;
}
```

**Service Operations:**
```protobuf
service IAMMapperService {
  // User-Role Mappings
  rpc AssignUserRoles(AssignUserRolesRequest) returns (google.protobuf.Empty) {}
  rpc RemoveUserRoles(RemoveUserRolesRequest) returns (google.protobuf.Empty) {}
  rpc GetUserRoles(GetUserRolesRequest) returns (GetUserRolesResponse) {}

  // Role-Permission Mappings
  rpc AssignRolePermissions(AssignRolePermissionsRequest) returns (google.protobuf.Empty) {}
  rpc RemoveRolePermissions(RemoveRolePermissionsRequest) returns (google.protobuf.Empty) {}
  rpc GetRolePermissions(GetRolePermissionsRequest) returns (GetRolePermissionsResponse) {}

  // User-Permission Mappings
  rpc AssignUserPermissions(AssignUserPermissionsRequest) returns (google.protobuf.Empty) {}
  rpc RemoveUserPermissions(RemoveUserPermissionsRequest) returns (google.protobuf.Empty) {}
  rpc GetUserPermissions(GetUserPermissionsRequest) returns (GetUserPermissionsResponse) {}

  // Project Members
  rpc AssignProjectMembers(AssignProjectMembersRequest) returns (google.protobuf.Empty) {}
  rpc RemoveProjectMembers(RemoveProjectMembersRequest) returns (google.protobuf.Empty) {}
  rpc GetProjectMembers(GetProjectMembersRequest) returns (GetProjectMembersResponse) {}
}
```

### 2. Domain Implementation (7-File Pattern)

#### internal/domain/iam_mapper/model.go

Define models for junction table records:

```go
type UserRole struct {
    ID        int64     `db:"id"`
    UserID    int64     `db:"user_id"`
    RoleID    int64     `db:"role_id"`
    CreatedAt time.Time `db:"created_at"`
}

type RolePermission struct {
    ID           int64     `db:"id"`
    RoleID       int64     `db:"role_id"`
    PermissionID int64     `db:"permission_id"`
    CreatedAt    time.Time `db:"created_at"`
}

type UserPermission struct {
    ID           int64     `db:"id"`
    UserID       int64     `db:"user_id"`
    PermissionID int64     `db:"permission_id"`
    CreatedAt    time.Time `db:"created_at"`
}

type ProjectMember struct {
    ID        int64     `db:"id"`
    ProjectID int64     `db:"project_id"`
    UserID    int64     `db:"user_id"`
    Role      string    `db:"role"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}
```

#### internal/domain/iam_mapper/interface.go

```go
type Repository interface {
    // User-Role Mappings
    AssignUserRoles(ctx context.Context, userID int64, roleIDs []int64) error
    RemoveUserRoles(ctx context.Context, userID int64, roleIDs []int64) error
    GetUserRoles(ctx context.Context, userID int64) ([]role.Role, error)

    // Role-Permission Mappings
    AssignRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
    RemoveRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
    GetRolePermissions(ctx context.Context, roleID int64) ([]permission.Permission, error)

    // User-Permission Mappings
    AssignUserPermissions(ctx context.Context, userID int64, permissionIDs []int64) error
    RemoveUserPermissions(ctx context.Context, userID int64, permissionIDs []int64) error
    GetUserPermissions(ctx context.Context, userID int64) ([]permission.Permission, error)

    // Project Members
    AssignProjectMembers(ctx context.Context, projectID int64, members []ProjectMemberInput) error
    RemoveProjectMembers(ctx context.Context, projectID int64, userIDs []int64) error
    GetProjectMembers(ctx context.Context, projectID int64) ([]ProjectMemberWithUser, error)
}
```

#### internal/domain/iam_mapper/repo.go

**Key Implementation Patterns:**

**Batch Insert with ON CONFLICT:**
```go
func (r *repository) AssignUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
    if len(roleIDs) == 0 {
        return nil
    }

    // Build multi-row INSERT with ON CONFLICT DO NOTHING
    query := `
        INSERT INTO altalune_users_roles (user_id, role_id)
        VALUES `

    args := []interface{}{userID}
    placeholders := []string{}

    for i, roleID := range roleIDs {
        placeholders = append(placeholders, fmt.Sprintf("($1, $%d)", i+2))
        args = append(args, roleID)
    }

    query += strings.Join(placeholders, ", ") + " ON CONFLICT (user_id, role_id) DO NOTHING"

    _, err := r.db.ExecContext(ctx, query, args...)
    return err
}
```

**Batch Delete:**
```go
func (r *repository) RemoveUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
    if len(roleIDs) == 0 {
        return nil
    }

    query := `
        DELETE FROM altalune_users_roles
        WHERE user_id = $1 AND role_id = ANY($2)
    `

    _, err := r.db.ExecContext(ctx, query, userID, pq.Array(roleIDs))
    return err
}
```

**Efficient JOIN Query:**
```go
func (r *repository) GetUserRoles(ctx context.Context, userID int64) ([]role.Role, error) {
    query := `
        SELECT r.id, r.public_id, r.name, r.description, r.created_at, r.updated_at
        FROM altalune_roles r
        INNER JOIN altalune_users_roles ur ON ur.role_id = r.id
        WHERE ur.user_id = $1
        ORDER BY r.name ASC
    `

    var roles []role.Role
    err := r.db.SelectContext(ctx, &roles, query, userID)
    return roles, err
}
```

**Project Members with User Data:**
```go
func (r *repository) GetProjectMembers(ctx context.Context, projectID int64) ([]ProjectMemberWithUser, error) {
    query := `
        SELECT
            u.id, u.public_id, u.email, u.first_name, u.last_name, u.is_active,
            pm.role, pm.created_at
        FROM altalune_project_members pm
        INNER JOIN altalune_users u ON u.id = pm.user_id
        WHERE pm.project_id = $1
        ORDER BY pm.created_at DESC
    `

    var members []ProjectMemberWithUser
    err := r.db.SelectContext(ctx, &members, query, projectID)
    return members, err
}
```

#### internal/domain/iam_mapper/service.go

- Convert public IDs to internal IDs before repository calls
- Validate that entities exist before creating mappings
- Validate project role enum (owner, admin, member, viewer)
- Handle empty arrays gracefully
- Proper error wrapping

#### internal/domain/iam_mapper/handler.go

- Implement all 12 RPC methods
- Request validation via protovalidate
- Call service methods
- Map domain errors to gRPC codes

#### internal/domain/iam_mapper/mapper.go

- Convert domain models to protobuf messages
- Handle nested objects (ProjectMemberWithUser contains User)
- Timestamp conversions

#### internal/domain/iam_mapper/errors.go

- Domain-specific error types for mapping operations
- Error constructors

### 3. Error Code Range

Add to `errors.go`:

**IAM Mapper Domain (608XX):**
- 60800: MappingNotFound
- 60801: MappingAlreadyExists (idempotent, can be warning instead of error)
- 60802: InvalidProjectRole
- 60803: CannotRemoveLastOwner (at least one owner required per project)
- 60804: UserNotFound (when assigning)
- 60805: RoleNotFound (when assigning)
- 60806: PermissionNotFound (when assigning)
- 60807: ProjectNotFound (when assigning members)

### 4. Container Wiring

Update `internal/container/container.go`:

```go
// Add repository field
IAMMapperRepo iam_mapper.Repository

// Add service field
IAMMapperService iam_mapper.Service

// Add handler field
IAMMapperHandler *iam_mapper.Handler
```

Wire up in `NewContainer()`:
1. Initialize repository with db connection
2. Initialize service with repository + user/role/permission services (for validation)
3. Initialize handler with service

### 5. Service Registration

**Update internal/server/grpc_services.go:**
```go
iammapperv1connect.RegisterIAMMapperServiceHandler(mux, container.IAMMapperHandler)
```

**Update internal/server/http_routes.go:**
```go
mux.Handle(iammapperv1connect.NewIAMMapperServiceHandler(container.IAMMapperHandler))
```

## Implementation Notes

### Batch Operations Efficiency

For TransferList UI, users may assign/remove multiple items at once. Use batch operations to minimize database round-trips:

**Bad (N+1 queries):**
```go
for _, roleID := range roleIDs {
    INSERT INTO altalune_users_roles (user_id, role_id) VALUES ($1, $2)
}
```

**Good (single query):**
```go
INSERT INTO altalune_users_roles (user_id, role_id)
VALUES ($1, $2), ($1, $3), ($1, $4), ...
ON CONFLICT DO NOTHING
```

### ON CONFLICT DO NOTHING

Use `ON CONFLICT DO NOTHING` for idempotent assign operations. This allows:
- Retry-safe operations
- Simplifies frontend logic (no need to check if already assigned)
- Graceful handling of concurrent requests

### Project Role Validation

Validate project role enum in service layer before inserting:
- Valid values: "owner", "admin", "member", "viewer"
- Case-sensitive matching
- Return `InvalidProjectRole` error for invalid values

### Last Owner Protection

Implement business rule: **A project must have at least one owner.**

When removing project members:
```go
if member.Role == "owner" {
    // Check if this is the last owner
    remainingOwners := countProjectOwners(ctx, projectID)
    if remainingOwners <= 1 {
        return ErrCannotRemoveLastOwner
    }
}
```

### Public ID to Internal ID Resolution

All protobuf APIs use public IDs (nanoid), but database uses internal IDs (bigint):

```go
func (s *service) AssignUserRoles(ctx context.Context, userPublicID string, rolePublicIDs []string) error {
    // Resolve user public ID
    user, err := s.userService.GetByPublicID(ctx, userPublicID)
    if err != nil {
        return ErrUserNotFound
    }

    // Resolve role public IDs
    roleIDs := make([]int64, len(rolePublicIDs))
    for i, publicID := range rolePublicIDs {
        role, err := s.roleService.GetByPublicID(ctx, publicID)
        if err != nil {
            return ErrRoleNotFound
        }
        roleIDs[i] = role.ID
    }

    return s.repo.AssignUserRoles(ctx, user.ID, roleIDs)
}
```

## Files to Create

**Protobuf:**
- `api/proto/altalune/v1/iam_mapper.proto`

**IAM Mapper Domain:**
- `internal/domain/iam_mapper/model.go`
- `internal/domain/iam_mapper/interface.go`
- `internal/domain/iam_mapper/repo.go`
- `internal/domain/iam_mapper/service.go`
- `internal/domain/iam_mapper/handler.go`
- `internal/domain/iam_mapper/mapper.go`
- `internal/domain/iam_mapper/errors.go`

## Files to Modify

- `errors.go` - Add 608XX error codes
- `internal/container/container.go` - Wire up iam_mapper domain
- `internal/server/grpc_services.go` - Register IAMMapperService
- `internal/server/http_routes.go` - Register HTTP handler

## Commands to Run

```bash
# Generate protobuf code
buf generate

# Build the application
make build

# Start server to test
./bin/app serve -c config.yaml
```

## Definition of Done

- [ ] Protobuf schema compiles without errors
- [ ] Generated Go code in `gen/altalune/v1/`
- [ ] Generated TypeScript code in `frontend/gen/`
- [ ] All 7 domain files created
- [ ] Container wiring complete and compiles
- [ ] Service registered in gRPC and HTTP
- [ ] Application builds successfully
- [ ] Can assign user roles via API
- [ ] Can remove user roles via API
- [ ] Can get user roles via API (returns role objects with JOIN)
- [ ] Can assign role permissions via API
- [ ] Can remove role permissions via API
- [ ] Can get role permissions via API
- [ ] Can assign user permissions via API
- [ ] Can remove user permissions via API
- [ ] Can get user permissions via API
- [ ] Can assign project members via API
- [ ] Can remove project members via API
- [ ] Can get project members via API (with user data)
- [ ] Batch operations work correctly (multiple assignments in single call)
- [ ] ON CONFLICT DO NOTHING prevents duplicate errors
- [ ] Cannot remove last project owner
- [ ] Invalid project role returns proper error
- [ ] Proper error responses when assigning non-existent entities

## Dependencies

- T7 (database schema) must be complete
- T8 (core entities) must be complete
- User, Role, Permission services for ID resolution

## Risk Factors

- **Medium risk**: Batch SQL operations with dynamic placeholders
- **Watch out for**: SQL injection when building dynamic queries (use parameterized queries!)
- **Test carefully**: Last owner protection logic
- **Verify**: ON CONFLICT syntax works with PostgreSQL version
