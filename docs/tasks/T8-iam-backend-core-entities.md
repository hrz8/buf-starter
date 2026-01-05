# Task T8: IAM Backend Core Entities (User, Role, Permission)

**Story Reference:** US3-iam-core-entities-and-mappings.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 12-15 hours
**Prerequisites:** T7 (Database schema must be created)

## Objective

Implement three complete backend domains (User, Role, Permission) following the established 7-file pattern with full CRUD operations, protobuf schemas, validation, and service registration.

## Acceptance Criteria

- [ ] Create protobuf schemas for UserService, RoleService, PermissionService
- [ ] Implement 3 complete 7-file domains (user, role, permission)
- [ ] All CRUD operations functional for each entity
- [ ] Email uniqueness validation (case-insensitive)
- [ ] Permission name validation with colon support (^[a-zA-Z0-9_:]+$)
- [ ] User activate/deactivate operations
- [ ] Error codes defined (605XX, 606XX, 607XX ranges)
- [ ] Container wiring complete
- [ ] Services registered in gRPC and HTTP routes
- [ ] Run `buf generate` successfully
- [ ] All unit tests pass

## Technical Requirements

### 1. Protobuf Schemas

#### api/proto/altalune/v1/user.proto

**User Message:**
```protobuf
message User {
  string id = 1;                                    // Public nanoid
  string email = 2;                                 // Unique, case-insensitive
  string first_name = 3;
  string last_name = 4;
  bool is_active = 5;
  google.protobuf.Timestamp created_at = 98;
  google.protobuf.Timestamp updated_at = 99;
}
```

**Service Operations:**
```protobuf
service UserService {
  rpc QueryUsers(QueryUsersRequest) returns (QueryUsersResponse) {}
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {}
  rpc GetUser(GetUserRequest) returns (GetUserResponse) {}
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse) {}
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse) {}
  rpc ActivateUser(ActivateUserRequest) returns (ActivateUserResponse) {}
  rpc DeactivateUser(DeactivateUserRequest) returns (DeactivateUserResponse) {}
}
```

**Validation Rules:**
- email: required, valid email format, max 255 characters, will be lowercased
- first_name: optional, 1-100 characters
- last_name: optional, 1-100 characters
- user_id: required for get/update/delete/activate/deactivate, 14-20 characters

#### api/proto/altalune/v1/role.proto

**Role Message:**
```protobuf
message Role {
  string id = 1;                                    // Public nanoid
  string name = 2;                                  // Unique
  string description = 3;
  google.protobuf.Timestamp created_at = 98;
  google.protobuf.Timestamp updated_at = 99;
}
```

**Service Operations:**
```protobuf
service RoleService {
  rpc QueryRoles(QueryRolesRequest) returns (QueryRolesResponse) {}
  rpc CreateRole(CreateRoleRequest) returns (CreateRoleResponse) {}
  rpc GetRole(GetRoleRequest) returns (GetRoleResponse) {}
  rpc UpdateRole(UpdateRoleRequest) returns (UpdateRoleResponse) {}
  rpc DeleteRole(DeleteRoleRequest) returns (DeleteRoleResponse) {}
}
```

**Validation Rules:**
- name: required, 2-100 characters, alphanumeric + spaces/hyphens/underscores, unique
- description: optional, max 500 characters
- role_id: required for get/update/delete, 14-20 characters

#### api/proto/altalune/v1/permission.proto

**Permission Message:**
```protobuf
message Permission {
  string id = 1;                                    // Public nanoid
  string name = 2;                                  // Machine-readable: "project:read"
  string effect = 3;                                // "allow" or "deny"
  string description = 4;                           // Human-readable (optional)
  google.protobuf.Timestamp created_at = 98;
  google.protobuf.Timestamp updated_at = 99;
}
```

**Service Operations:**
```protobuf
service PermissionService {
  rpc QueryPermissions(QueryPermissionsRequest) returns (QueryPermissionsResponse) {}
  rpc CreatePermission(CreatePermissionRequest) returns (CreatePermissionResponse) {}
  rpc GetPermission(GetPermissionRequest) returns (GetPermissionResponse) {}
  rpc UpdatePermission(UpdatePermissionRequest) returns (UpdatePermissionResponse) {}
  rpc DeletePermission(DeletePermissionRequest) returns (DeletePermissionResponse) {}
}
```

**Validation Rules:**
- name: required, 2-100 characters, pattern ^[a-zA-Z0-9_:]+$ (allows colons), unique
- effect: required, enum ["allow", "deny"], defaults to "allow"
- description: optional, max 500 characters
- permission_id: required for get/update/delete, 14-20 characters

### 2. Domain Implementation (7-File Pattern)

For each domain (user, role, permission), create these files:

#### internal/domain/{domain}/model.go
- Go struct matching database schema
- Both internal ID (int64) and public ID (nanoid string)
- created_at and updated_at timestamps
- Validation tags

#### internal/domain/{domain}/interface.go
- Repository interface with all CRUD methods
- Service interface with all business operations
- Clear method signatures with context.Context

#### internal/domain/{domain}/repo.go
- PostgreSQL repository implementation
- **CRITICAL**: Non-partitioned queries (NO project_id filtering)
- Email uniqueness check (case-insensitive) for users
- Name uniqueness checks for roles and permissions
- Efficient SELECT queries with proper indexing
- INSERT with RETURNING for created resources
- UPDATE with optimistic locking via updated_at
- DELETE with proper error handling

**Example User Query (NO project_id):**
```go
const queryUsersSQL = `
  SELECT id, public_id, email, first_name, last_name, is_active, created_at, updated_at
  FROM altalune_users
  WHERE ($1 = '' OR email ILIKE '%' || $1 || '%')
    AND ($2::boolean IS NULL OR is_active = $2)
  ORDER BY created_at DESC
  LIMIT $3 OFFSET $4
`
```

#### internal/domain/{domain}/service.go
- Business logic layer
- Validation before repository calls
- Nanoid generation for public IDs
- Email lowercase enforcement for users
- Duplicate checking (email, role name, permission name)
- Error wrapping with domain errors

#### internal/domain/{domain}/handler.go
- Connect-RPC handler implementation
- Request validation via protovalidate
- Call service methods
- Map domain errors to gRPC codes
- Response mapping

#### internal/domain/{domain}/mapper.go
- Convert domain models to protobuf messages
- Convert protobuf messages to domain models
- Handle optional fields properly
- Timestamp conversions

#### internal/domain/{domain}/errors.go
- Domain-specific error types
- Error constructors
- Error checking helpers
- Connect error code mapping

### 3. Error Code Ranges

Add to `errors.go`:

**User Domain (605XX):**
- 60500: UserNotFound
- 60501: UserAlreadyExists (email conflict)
- 60502: UserInvalidEmail
- 60503: UserAlreadyActive
- 60504: UserAlreadyInactive
- 60505: UserCannotDeleteSelf

**Role Domain (606XX):**
- 60600: RoleNotFound
- 60601: RoleAlreadyExists (name conflict)
- 60602: RoleInvalidName
- 60603: RoleInUse (cannot delete, has users)

**Permission Domain (607XX):**
- 60700: PermissionNotFound
- 60701: PermissionAlreadyExists (name conflict)
- 60702: PermissionInvalidName
- 60703: PermissionInvalidEffect
- 60704: PermissionInUse (cannot delete, assigned to roles/users)

### 4. Container Wiring

Update `internal/container/container.go`:

```go
// Add repository fields
UserRepo       user.Repository
RoleRepo       role.Repository
PermissionRepo permission.Repository

// Add service fields
UserService       user.Service
RoleService       role.Service
PermissionService permission.Service

// Add handler fields
UserHandler       *user.Handler
RoleHandler       *role.Handler
PermissionHandler *permission.Handler
```

Wire up in `NewContainer()`:
1. Initialize repositories with db connection
2. Initialize services with repositories
3. Initialize handlers with services

### 5. Service Registration

**Update internal/server/grpc_services.go:**
```go
userv1connect.RegisterUserServiceHandler(mux, container.UserHandler)
rolev1connect.RegisterRoleServiceHandler(mux, container.RoleHandler)
permissionv1connect.RegisterPermissionServiceHandler(mux, container.PermissionHandler)
```

**Update internal/server/http_routes.go:**
```go
// User routes
mux.Handle(userv1connect.NewUserServiceHandler(container.UserHandler))

// Role routes
mux.Handle(rolev1connect.NewRoleServiceHandler(container.RoleHandler))

// Permission routes
mux.Handle(permissionv1connect.NewPermissionServiceHandler(container.PermissionHandler))
```

## Implementation Notes

### Non-Partitioned Query Patterns

**CRITICAL DIFFERENCE**: These IAM tables are NOT partitioned by project_id.

**Standard partitioned query:**
```go
// DON'T DO THIS for IAM tables
WHERE project_id = $1 AND id = $2
```

**IAM query pattern:**
```go
// DO THIS instead
WHERE id = $1  // No project_id!
```

### User Email Handling

- Always convert email to lowercase before database operations
- Use case-insensitive comparison in queries (ILIKE or LOWER())
- Database has CHECK constraint for lowercase enforcement

### Permission Name Validation

- Regex: ^[a-zA-Z0-9_:]+$
- Allows: letters, numbers, underscores, colons
- Examples: "root", "project:read", "user:manage:delete", "api_key:create"

### Activate/Deactivate User

- Separate RPCs from Update to make intent explicit
- Check current state before toggling (prevent redundant operations)
- Return appropriate errors if already in target state

## Files to Create

**Protobuf:**
- `api/proto/altalune/v1/user.proto`
- `api/proto/altalune/v1/role.proto`
- `api/proto/altalune/v1/permission.proto`

**User Domain:**
- `internal/domain/user/model.go`
- `internal/domain/user/interface.go`
- `internal/domain/user/repo.go`
- `internal/domain/user/service.go`
- `internal/domain/user/handler.go`
- `internal/domain/user/mapper.go`
- `internal/domain/user/errors.go`

**Role Domain:**
- `internal/domain/role/model.go`
- `internal/domain/role/interface.go`
- `internal/domain/role/repo.go`
- `internal/domain/role/service.go`
- `internal/domain/role/handler.go`
- `internal/domain/role/mapper.go`
- `internal/domain/role/errors.go`

**Permission Domain:**
- `internal/domain/permission/model.go`
- `internal/domain/permission/interface.go`
- `internal/domain/permission/repo.go`
- `internal/domain/permission/service.go`
- `internal/domain/permission/handler.go`
- `internal/domain/permission/mapper.go`
- `internal/domain/permission/errors.go`

## Files to Modify

- `errors.go` - Add error codes and constructors
- `internal/container/container.go` - Wire up all 3 domains
- `internal/server/grpc_services.go` - Register gRPC services
- `internal/server/http_routes.go` - Register HTTP handlers

## Commands to Run

```bash
# Generate protobuf code
buf generate

# Build the application
make build

# Run migrations (if not already done)
./bin/app migrate -c config.yaml

# Start server to test
./bin/app serve -c config.yaml
```

## Definition of Done

- [ ] All 3 protobuf schemas compile without errors
- [ ] Generated Go code in `gen/altalune/v1/`
- [ ] Generated TypeScript code in `frontend/gen/`
- [ ] All 21 domain files created (7 per entity)
- [ ] Container wiring complete and compiles
- [ ] Services registered in gRPC and HTTP
- [ ] Application builds successfully
- [ ] Can create users, roles, permissions via API
- [ ] Can query, get, update, delete each entity
- [ ] User activate/deactivate works correctly
- [ ] Email uniqueness enforced (case-insensitive)
- [ ] Role/permission name uniqueness enforced
- [ ] Permission name validation allows colons
- [ ] Proper error responses for all failure cases
- [ ] Non-partitioned queries working (no project_id)

## Dependencies

- T7 (database schema) must be complete
- Existing protobuf infrastructure
- Existing container and server setup
- buf.validate for validation rules

## Risk Factors

- **Medium risk**: Non-partitioned pattern different from existing code
- **Watch out for**: Accidentally adding project_id to queries
- **Test carefully**: Email case-insensitivity and uniqueness
- **Verify**: Permission name regex allows colons but rejects spaces/special chars
