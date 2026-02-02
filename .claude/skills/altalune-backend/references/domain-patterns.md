# Domain Implementation Patterns

## Table of Contents
- [Protobuf Schema](#protobuf-schema)
- [Domain Models](#domain-models)
- [Repository Layer](#repository-layer)
- [Service Layer](#service-layer)
- [Handler Layer](#handler-layer)
- [Error Handling](#error-handling)
- [Container Registration](#container-registration)
- [Route Registration](#route-registration)

## Protobuf Schema

**File: `api/proto/altalune/v1/{domain}.proto`**

```protobuf
syntax = "proto3";
package altalune.v1;

import "google/protobuf/timestamp.proto";
import "buf/validate/validate.proto";
import "altalune/v1/common.proto";

// Main entity message
message Entity {
  string id = 1;                                    // Public ID (nanoid)
  string name = 2;
  EntityStatus status = 3;
  google.protobuf.Timestamp created_at = 98;
  google.protobuf.Timestamp updated_at = 99;
}

enum EntityStatus {
  ENTITY_STATUS_UNSPECIFIED = 0;
  ENTITY_STATUS_ACTIVE = 1;
  ENTITY_STATUS_INACTIVE = 2;
}

// Query request with standard patterns
message QueryEntityRequest {
  altalune.v1.QueryRequestParams query = 1;
  string project_id = 2 [(buf.validate.field).string.len = 14];
}

message QueryEntityResponse {
  repeated Entity data = 1;
  altalune.v1.QueryMeta meta = 2;
}

// Create request with validation
message CreateEntityRequest {
  string project_id = 1 [(buf.validate.field).string.len = 14];
  string name = 2 [(buf.validate.field) = {
    required: true,
    string: { min_len: 1, max_len: 50 }
  }];
}

message CreateEntityResponse {
  Entity entity = 1;
  string message = 2;
}

// Service definition
service EntityService {
  rpc QueryEntities(QueryEntityRequest) returns (QueryEntityResponse) {}
  rpc CreateEntity(CreateEntityRequest) returns (CreateEntityResponse) {}
  rpc GetEntity(GetEntityRequest) returns (GetEntityResponse) {}
  rpc UpdateEntity(UpdateEntityRequest) returns (UpdateEntityResponse) {}
  rpc DeleteEntity(DeleteEntityRequest) returns (DeleteEntityResponse) {}
}
```

**Key Requirements:**
- Use `buf.validate` for all input validation
- Timestamp fields at positions 98, 99
- Include `project_id` for project-scoped entities
- Follow enum naming: `{ENTITY}_{STATUS}_{VALUE}`

## Domain Models

**File: `internal/domain/{domain}/model.go`**

```go
package entity

import (
    "time"
    protov1 "your-module/gen/altalune/v1"
    "google.golang.org/protobuf/types/known/timestamppb"
)

// Domain enums as string constants
type EntityStatus string

const (
    EntityStatusActive   EntityStatus = "active"
    EntityStatusInactive EntityStatus = "inactive"
)

// Query result (includes internal int64 ID)
type EntityQueryResult struct {
    ID        int64
    PublicID  string
    Name      string
    Status    EntityStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Domain model (uses public string ID)
type Entity struct {
    ID        string
    Name      string
    Status    EntityStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Convert query result to domain model
func (r *EntityQueryResult) ToEntity() *Entity {
    return &Entity{
        ID:        r.PublicID,
        Name:      r.Name,
        Status:    r.Status,
        CreatedAt: r.CreatedAt,
        UpdatedAt: r.UpdatedAt,
    }
}

// Convert to protobuf
func (m *Entity) ToProto() *protov1.Entity {
    return &protov1.Entity{
        Id:        m.ID,
        Name:      m.Name,
        Status:    mapStatusToProto(m.Status),
        CreatedAt: timestamppb.New(m.CreatedAt),
        UpdatedAt: timestamppb.New(m.UpdatedAt),
    }
}

// Input/Result structs for operations
type CreateEntityInput struct {
    ProjectID int64
    Name      string
}

type CreateEntityResult struct {
    ID        int64
    PublicID  string
    Name      string
    Status    EntityStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Repository Layer

**File: `internal/domain/{domain}/interface.go`**

```go
package entity

import (
    "context"
    "your-module/internal/shared/query"
)

type Repositor interface {
    Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[EntityQueryResult], error)
    Create(ctx context.Context, input *CreateEntityInput) (*CreateEntityResult, error)
    GetByPublicID(ctx context.Context, projectID int64, publicID string) (*EntityQueryResult, error)
    Update(ctx context.Context, input *UpdateEntityInput) (*UpdateEntityResult, error)
    Delete(ctx context.Context, projectID int64, publicID string) error
}
```

**File: `internal/domain/{domain}/repo.go`**

```go
package entity

import (
    "context"
    "fmt"
    "your-module/internal/shared/nanoid"
    "your-module/internal/shared/postgres"
    "your-module/internal/shared/query"
)

type Repo struct {
    db postgres.DB
}

func NewRepo(db postgres.DB) *Repo {
    return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, input *CreateEntityInput) (*CreateEntityResult, error) {
    // Generate public ID
    publicID, err := nanoid.GeneratePublicID()
    if err != nil {
        return nil, fmt.Errorf("generate public id: %w", err)
    }

    // Insert with RETURNING clause
    query := `
        INSERT INTO altalune_entities (public_id, project_id, name, status, created_at, updated_at)
        VALUES ($1, $2, $3, 'active', NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

    result := &CreateEntityResult{
        PublicID: publicID,
        Name:     input.Name,
        Status:   EntityStatusActive,
    }

    err = r.db.QueryRowContext(ctx, query, publicID, input.ProjectID, input.Name).
        Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)
    if err != nil {
        if postgres.IsUniqueViolation(err) {
            return nil, ErrEntityAlreadyExists
        }
        return nil, fmt.Errorf("insert entity: %w", err)
    }

    return result, nil
}

func (r *Repo) Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[EntityQueryResult], error) {
    // Base query
    baseQuery := `
        SELECT id, public_id, name, status, created_at, updated_at
        FROM altalune_entities
        WHERE project_id = $1
    `

    // Build WHERE clause with params
    whereClause, args := query.BuildWhereClause(params, 2) // Start from $2
    args = append([]interface{}{projectID}, args...)

    // Count total
    countQuery := fmt.Sprintf("SELECT COUNT(*) FROM altalune_entities WHERE project_id = $1 %s", whereClause)
    var totalCount int
    if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
        return nil, fmt.Errorf("count entities: %w", err)
    }

    // Add sorting and pagination
    fullQuery := fmt.Sprintf("%s %s %s %s",
        baseQuery, whereClause,
        query.BuildSortClause(params),
        query.BuildPaginationClause(params),
    )

    rows, err := r.db.QueryContext(ctx, fullQuery, args...)
    if err != nil {
        return nil, fmt.Errorf("query entities: %w", err)
    }
    defer rows.Close()

    var results []EntityQueryResult
    for rows.Next() {
        var e EntityQueryResult
        var statusStr string
        if err := rows.Scan(&e.ID, &e.PublicID, &e.Name, &statusStr, &e.CreatedAt, &e.UpdatedAt); err != nil {
            return nil, fmt.Errorf("scan entity: %w", err)
        }
        e.Status = EntityStatus(statusStr)
        results = append(results, e)
    }

    return &query.QueryResult[EntityQueryResult]{
        Data:       results,
        TotalCount: totalCount,
        PageCount:  query.CalculatePageCount(totalCount, params.PageSize),
    }, nil
}
```

## Service Layer

**File: `internal/domain/{domain}/service.go`**

```go
package entity

import (
    "context"
    "log/slog"

    "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
    "github.com/bufbuild/protovalidate-go"
    protov1 "your-module/gen/altalune/v1"
    "your-module/internal/domain/altalune"
    "your-module/internal/shared/query"
)

type Service struct {
    protov1.UnimplementedEntityServiceServer
    validator *protovalidate.Validator
    logger    *slog.Logger
    repo      Repositor
    projectRepo ProjectRepositor // For resolving project public ID to internal ID
}

func NewService(
    validator *protovalidate.Validator,
    logger *slog.Logger,
    repo Repositor,
    projectRepo ProjectRepositor,
) *Service {
    return &Service{
        validator:   validator,
        logger:      logger,
        repo:        repo,
        projectRepo: projectRepo,
    }
}

func (s *Service) CreateEntity(ctx context.Context, req *protov1.CreateEntityRequest) (*protov1.CreateEntityResponse, error) {
    // 1. Validate request
    if err := s.validator.Validate(req); err != nil {
        return nil, altalune.NewInvalidInputError(err)
    }

    // 2. Resolve project ID (public to internal)
    project, err := s.projectRepo.GetByPublicID(ctx, req.ProjectId)
    if err != nil {
        s.logger.Error("failed to get project", "error", err, "project_id", req.ProjectId)
        return nil, altalune.NewNotFoundError(err)
    }

    // 3. Map protobuf to domain input
    input := &CreateEntityInput{
        ProjectID: project.ID,
        Name:      req.Name,
    }

    // 4. Call repository
    result, err := s.repo.Create(ctx, input)
    if err != nil {
        if err == ErrEntityAlreadyExists {
            return nil, altalune.NewAlreadyExistsError(err)
        }
        s.logger.Error("failed to create entity", "error", err)
        return nil, altalune.NewInternalError(err)
    }

    // 5. Map result to protobuf response
    entity := &Entity{
        ID:        result.PublicID,
        Name:      result.Name,
        Status:    result.Status,
        CreatedAt: result.CreatedAt,
        UpdatedAt: result.UpdatedAt,
    }

    return &protov1.CreateEntityResponse{
        Entity:  entity.ToProto(),
        Message: "Entity created successfully",
    }, nil
}

func (s *Service) QueryEntities(ctx context.Context, req *protov1.QueryEntityRequest) (*protov1.QueryEntityResponse, error) {
    // Validate
    if err := s.validator.Validate(req); err != nil {
        return nil, altalune.NewInvalidInputError(err)
    }

    // Resolve project
    project, err := s.projectRepo.GetByPublicID(ctx, req.ProjectId)
    if err != nil {
        return nil, altalune.NewNotFoundError(err)
    }

    // Map query params
    params := query.FromProto(req.Query)

    // Execute query
    result, err := s.repo.Query(ctx, project.ID, params)
    if err != nil {
        s.logger.Error("failed to query entities", "error", err)
        return nil, altalune.NewInternalError(err)
    }

    // Map to proto
    var data []*protov1.Entity
    for _, e := range result.Data {
        data = append(data, e.ToEntity().ToProto())
    }

    return &protov1.QueryEntityResponse{
        Data: data,
        Meta: &protov1.QueryMeta{
            RowCount:  int32(result.TotalCount),
            PageCount: int32(result.PageCount),
        },
    }, nil
}
```

## Handler Layer

**File: `internal/domain/{domain}/handler.go`**

```go
package entity

import (
    "context"
    "connectrpc.com/connect"
    protov1 "your-module/gen/altalune/v1"
    "your-module/internal/domain/altalune"
)

type Handler struct {
    svc *Service
}

func NewHandler(svc *Service) *Handler {
    return &Handler{svc: svc}
}

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

func (h *Handler) QueryEntities(
    ctx context.Context,
    req *connect.Request[protov1.QueryEntityRequest],
) (*connect.Response[protov1.QueryEntityResponse], error) {
    response, err := h.svc.QueryEntities(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(response), nil
}
```

## Error Handling

**File: `internal/domain/{domain}/errors.go`**

```go
package entity

import "errors"

var (
    ErrEntityNotFound      = errors.New("entity not found")
    ErrEntityAlreadyExists = errors.New("entity with this name already exists")
)
```

**Using shared error helpers (`internal/domain/altalune/`):**

```go
// Create errors with proper codes
altalune.NewInvalidInputError(err)    // Invalid input
altalune.NewNotFoundError(err)        // Resource not found
altalune.NewAlreadyExistsError(err)   // Duplicate/conflict
altalune.NewInternalError(err)        // Internal server error
altalune.NewUnauthorizedError(err)    // Authentication required
altalune.NewPermissionDeniedError(err) // Authorization denied

// Convert to Connect error for handlers
altalune.ToConnectError(err)
```

## Container Registration

**File: `internal/container/container.go`**

```go
// In initRepositories()
func (c *Container) initRepositories() {
    // ... existing repos ...
    c.entityRepo = entity.NewRepo(c.db)
}

// In initServices()
func (c *Container) initServices() {
    // ... existing services ...
    c.entitySvc = entity.NewService(
        c.validator,
        c.logger,
        c.entityRepo,
        c.projectRepo,
    )
}

// Add getter methods
func (c *Container) GetEntityService() *entity.Service {
    return c.entitySvc
}
```

## Route Registration

**File: `internal/server/http_routes.go`**

```go
// In setupConnectRPCRoutes()
func (s *Server) setupConnectRPCRoutes(mux *http.ServeMux) {
    // ... existing routes ...

    entityHandler := entity.NewHandler(s.c.GetEntityService())
    entityPath, entityConnectHandler := protov1connect.NewEntityServiceHandler(entityHandler)
    connectrpcMux.Handle(entityPath, entityConnectHandler)
}
```

**For gRPC registration:**

**File: `internal/server/grpc_services.go`**

```go
func (s *Server) registerGRPCServices(grpcServer *grpc.Server) {
    // ... existing services ...
    protov1.RegisterEntityServiceServer(grpcServer, s.c.GetEntityService())
}
```
