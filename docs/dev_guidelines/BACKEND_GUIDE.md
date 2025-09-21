## Backend Development Workflow

This systematic workflow ensures consistency, reusability, and proper architecture adherence when implementing new backend endpoints/services.

### Step 1: Reusability Analysis

**Before creating anything new, check for reusable components:**

**✅ Infrastructure Components (Reuse These):**

- `nanoid.GeneratePublicID()` - Public ID generation
- `query.QueryParams/QueryResult[T]` - Query architecture with filtering/pagination
- `postgres.DB` interface and utilities (`IsUniqueViolation`, etc.)
- `altalune.NewXXXError()` helpers - Standardized error creation
- `protovalidate.Validator` - Request validation
- Container/DI patterns in `container.go`

**✅ Architecture Patterns (Follow These):**

- Dual ID system: public nanoid (string) + internal database ID (int64)
- Enum mapping: domain constants ↔ database strings ↔ protobuf enums
- Service structure: validator + logger + repository dependency injection
- Repository query pattern: build query → count → sort → paginate → get filters
- Error flow: domain errors → `altalune.NewXXXError()` → `altalune.ToConnectError()`

**⚠️ Database Partitioning Requirements:**

- Tables that store project-specific data **MUST** be partitioned by `project_id`
- When creating new project-related tables, add to `partitionedTables` in `internal/domain/project/repo.go`
- Partition pattern: `PARTITION BY LIST (project_id)` with `PRIMARY KEY (project_id, id)`
- Partitions are auto-created during project creation via `createPartitionsForProject()`

**✅ Existing Domain Components (Extend If Needed):**

- Repository interfaces and implementations
- Service methods and business logic
- Handler methods and Connect-RPC wrappers
- Domain models and conversion methods

### Step 2: Define Protobuf Schema

**File: `api/proto/altalune/v1/{domain}.proto`**

```protobuf
// Required imports
import "google/protobuf/timestamp.proto";
import "buf/validate/validate.proto";
import "altalune/v1/common.proto";

// Define main entity message
message EntityName {
  string id = 1;                                    // Public ID (nanoid)
  // ... domain-specific fields with validation
  google.protobuf.Timestamp created_at = 98;
  google.protobuf.Timestamp updated_at = 99;
}

// Define request/response messages
message CreateEntityRequest {
  // Use comprehensive buf.validate constraints
  string field = 1 [(buf.validate.field).required = true];
}

message CreateEntityResponse {
  EntityName entity = 1;
  string message = 2;
}

// Extend or create service
service EntityService {
  rpc CreateEntity(CreateEntityRequest) returns (CreateEntityResponse) {}
}
```

**Key Requirements:**

- Use `buf.validate` for all input validation
- Follow established field naming and numbering conventions
- Include comprehensive validation rules (required, patterns, length limits)
- Use consistent timestamp field positions (98, 99)

### Step 3: Create Domain Models

**File: `internal/domain/{domain}/model.go`**

```go
// Define domain enums as string constants
type EntityStatus string
const (
    EntityStatusActive   EntityStatus = "active"
    EntityStatusInactive EntityStatus = "inactive"
)

// Query result struct (with internal int64 ID)
type EntityQueryResult struct {
    ID       int64          // Database primary key
    PublicID string         // External nanoid
    // ... other fields
    Status   EntityStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Domain model (with public string ID)
type Entity struct {
    ID       string         // Public nanoid for external APIs
    // ... other fields
    Status   EntityStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Conversion methods
func (r *EntityQueryResult) ToEntity() *Entity { ... }
func (m *Entity) ToEntityProto() *protov1.Entity { ... }

// Input/Result structs for operations
type CreateEntityInput struct { ... }
type CreateEntityResult struct { ... }
```

### Step 4: Create/Extend Repository

**File: `internal/domain/{domain}/repo.go`**

**Only create new repository if domain doesn't exist. Otherwise extend existing interface:**

```go
// Extend interface in interface.go
type Repositor interface {
    // Existing methods...
    Create(ctx context.Context, input *CreateEntityInput) (*CreateEntityResult, error)
    GetByUniqueField(ctx context.Context, field string) (*Entity, error)
}

// Implement in repo.go following established patterns:
func (r *Repo) Create(ctx context.Context, input *CreateEntityInput) (*CreateEntityResult, error) {
    // 1. Generate public ID
    publicID, _ := nanoid.GeneratePublicID()

    // 2. Map domain enums to database strings

    // 3. Execute INSERT with RETURNING clause

    // 4. Handle unique constraint violations with postgres.IsUniqueViolation()

    // 5. Map database strings back to domain enums

    // 6. Return result
}
```

### Step 5: Create/Extend Service Layer

**File: `internal/domain/{domain}/service.go`**

**Extend existing service or create new one following established pattern:**

```go
func (s *Service) CreateEntity(ctx context.Context, req *protov1.CreateEntityRequest) (*protov1.CreateEntityResponse, error) {
    // 1. Validate request using s.validator.Validate(req)

    // 2. Check business rules (duplicates, relationships, etc.)

    // 3. Map protobuf enums to domain enums

    // 4. Call repository with domain input

    // 5. Handle domain errors (convert to altalune.NewXXXError())

    // 6. Map result back to protobuf response with timestamppb.New()

    // 7. Return response with success message
}
```

**Key Requirements:**

- Always validate requests first using `s.validator.Validate(req)`
- Check business rules before repository calls
- Use proper error handling with domain-specific errors
- Include comprehensive logging for failures
- Map between protobuf ↔ domain ↔ database representations

### Step 6: Create/Extend Handler Layer

**File: `internal/domain/{domain}/handler.go`**

**Extend existing handler or create new one:**

```go
func (h *Handler) CreateEntity(
    ctx context.Context,
    req *connect.Request[protov1.CreateEntityRequest],
) (*connect.Response[protov1.CreateEntityResponse], error) {
    response, err := h.svc.CreateEntity(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(response), nil
}
```

### Step 7: Create Custom Domain Errors

**File: `internal/domain/{domain}/errors.go`**

**Only create new errors as needed:**

```go
var (
    ErrEntityNotFound      = errors.New("entity not found")
    ErrEntityAlreadyExists = errors.New("entity with this field already exists")
    // ... other domain-specific errors
)
```

### Step 8: Integration & Registration

**1. Update Container (`internal/container/container.go`):**

```go
// Add to initRepositories() if new repository
// Add to initServices() if new service
```

**2. Register gRPC Service (`internal/server/grpc_services.go`):**

```go
protov1.RegisterEntityServiceServer(grpcServer, s.c.GetEntityService())
```

**3. Register Connect-RPC Handler (`internal/server/http_routes.go`):**

```go
entityHandler := entity_domain.NewHandler(s.c.GetEntityService())
entityPath, entityConnectHandler := protov1connect.NewEntityServiceHandler(entityHandler)
connectrpcMux.Handle(entityPath, entityConnectHandler)
```

### Step 9: Code Generation & Testing

```bash
# Generate protobuf code
buf generate

# Build and test
go build -o ./tmp/test-app cmd/altalune/*.go

# Test via HTTP/Connect-RPC
curl -X POST http://localhost:8080/api/domain.v1.EntityService/CreateEntity \
  -H "Content-Type: application/json" \
  -d '{"field": "value"}'

# Test via gRPC
grpcurl -plaintext -d '{"field": "value"}' \
  localhost:8080 domain.v1.EntityService/CreateEntity
```

### Development Checklist

**Before Implementation:**

- [ ] Analyzed existing components for reusability
- [ ] Checked if domain/repository already exists
- [ ] Reviewed similar domain implementations for patterns

**During Implementation:**

- [ ] Used `buf.validate` for all input validation
- [ ] Followed dual ID system (public nanoid + internal int64)
- [ ] Implemented proper enum mapping patterns
- [ ] Added comprehensive error handling with domain errors
- [ ] Used established query patterns for database operations
- [ ] Included proper logging for debugging

**After Implementation:**

- [ ] Generated protobuf code with `buf generate`
- [ ] Registered services in container, gRPC, and Connect-RPC handlers
- [ ] Tested both HTTP/Connect-RPC and gRPC endpoints
- [ ] Verified error handling and validation works correctly

This workflow ensures consistency with established patterns while maximizing code reuse and maintaining clean architecture boundaries.
