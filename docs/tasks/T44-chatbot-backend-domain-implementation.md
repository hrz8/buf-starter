# Task T44: Chatbot Backend Domain Implementation

**Story Reference:** US11-chatbot-configuration-foundation.md
**Type:** Backend
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** T43 (Database Schema)

## Objective

Implement complete chatbot domain following the 7-file pattern with proto definitions and service registration. This provides the backend API for reading and updating chatbot module configurations.

## Acceptance Criteria

- [ ] Proto schema defines GetChatbotConfig and UpdateModuleConfig RPCs
- [ ] Domain implements 7-file pattern (model, interface, repo, service, handler, mapper, errors)
- [ ] Service registered in container and HTTP routes
- [ ] API endpoints functional via Connect-RPC
- [ ] JSONB modules_config handled correctly with google.protobuf.Struct

## Technical Requirements

### Proto Schema

File: `api/proto/altalune/v1/chatbot.proto`

```protobuf
syntax = "proto3";

package altalune.v1;

option go_package = "github.com/hrz8/altalune/gen/altalune/v1;altalunev1";

import "buf/validate/validate.proto";
import "google/protobuf/struct.proto";
import "google/protobuf/timestamp.proto";

// ChatbotService provides operations for managing chatbot configurations
service ChatbotService {
  // GetChatbotConfig retrieves the chatbot configuration for a project
  rpc GetChatbotConfig(GetChatbotConfigRequest) returns (GetChatbotConfigResponse) {}

  // UpdateModuleConfig updates a specific module's configuration
  rpc UpdateModuleConfig(UpdateModuleConfigRequest) returns (UpdateModuleConfigResponse) {}
}

message GetChatbotConfigRequest {
  string project_id = 1 [(buf.validate.field).required = true];
}

message GetChatbotConfigResponse {
  string id = 1;
  google.protobuf.Struct modules_config = 2;
  google.protobuf.Timestamp updated_at = 3;
}

message UpdateModuleConfigRequest {
  string project_id = 1 [(buf.validate.field).required = true];
  string module_name = 2 [(buf.validate.field).required = true];
  google.protobuf.Struct config = 3 [(buf.validate.field).required = true];
}

message UpdateModuleConfigResponse {
  string message = 1;
  google.protobuf.Struct updated_config = 2;
}
```

### Domain Error Codes

- 61001: Chatbot config not found
- 61002: Invalid module configuration
- 61003: Invalid module name

### 7-File Pattern Implementation

#### model.go

```go
package chatbot

import (
    "encoding/json"
    "time"

    "google.golang.org/protobuf/types/known/structpb"
    "google.golang.org/protobuf/types/known/timestamppb"

    altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

// ModulesConfig represents the complete modules configuration
type ModulesConfig map[string]interface{}

// ChatbotConfig represents the domain model for chatbot configuration
type ChatbotConfig struct {
    ID            int64
    PublicID      string
    ProjectID     int64
    ModulesConfig ModulesConfig
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

// ToProto converts domain model to proto response
func (c *ChatbotConfig) ToProto() (*altalunev1.GetChatbotConfigResponse, error) {
    configStruct, err := structpb.NewStruct(c.ModulesConfig)
    if err != nil {
        return nil, err
    }

    return &altalunev1.GetChatbotConfigResponse{
        Id:            c.PublicID,
        ModulesConfig: configStruct,
        UpdatedAt:     timestamppb.New(c.UpdatedAt),
    }, nil
}

// ValidModuleNames contains all valid module names
var ValidModuleNames = map[string]bool{
    "llm":       true,
    "mcpServer": true,
    "widget":    true,
    "prompt":    true,
}
```

#### interface.go

```go
package chatbot

import "context"

// Repositor defines the interface for chatbot repository operations
type Repositor interface {
    // GetByProjectID retrieves the chatbot config for a project
    GetByProjectID(ctx context.Context, projectID int64) (*ChatbotConfig, error)

    // UpdateModuleConfig updates a specific module's configuration
    UpdateModuleConfig(ctx context.Context, projectID int64, moduleName string, config map[string]interface{}) (*ChatbotConfig, error)
}
```

#### errors.go

```go
package chatbot

import "errors"

var (
    ErrChatbotConfigNotFound     = errors.New("chatbot config not found")
    ErrInvalidModuleConfig       = errors.New("invalid module configuration")
    ErrInvalidModuleName         = errors.New("invalid module name")
)
```

#### repo.go

Key methods:
- `GetByProjectID(ctx, projectID)` - Retrieves config, parses JSONB
- `UpdateModuleConfig(ctx, projectID, moduleName, config)` - Uses JSONB path update

JSONB Update Query:
```sql
UPDATE altalune_chatbot_configs
SET modules_config = jsonb_set(modules_config, $2::text[], $3::jsonb),
    updated_at = NOW()
WHERE project_id = $1
RETURNING id, public_id, project_id, modules_config, created_at, updated_at
```

#### service.go

```go
type Service struct {
    validator   protovalidate.Validator
    log         altalune.Logger
    projectRepo project.Repositor
    chatbotRepo Repositor
}

func NewService(
    v protovalidate.Validator,
    log altalune.Logger,
    projectRepo project.Repositor,
    chatbotRepo Repositor,
) *Service {
    return &Service{
        validator:   v,
        log:         log,
        projectRepo: projectRepo,
        chatbotRepo: chatbotRepo,
    }
}
```

Service methods:
1. Validate request with protovalidate
2. Get internal project ID from public ID
3. Call repository method
4. Convert to proto response
5. Handle errors with domain error codes

#### handler.go

Thin Connect-RPC wrapper around service.

#### mapper.go

Helper functions for proto ↔ domain conversions.

## Files to Create

```
api/proto/altalune/v1/chatbot.proto

internal/domain/chatbot/
├── model.go
├── interface.go
├── errors.go
├── mapper.go
├── repo.go
├── service.go
└── handler.go
```

## Files to Modify

- `internal/container/container.go`
  - Add `chatbotRepo` and `chatbotService` fields to Container struct
  - Initialize in `initRepositories()` and `initServices()`

- `internal/container/getter.go`
  - Add `GetChatbotService()` method

- `internal/server/http_routes.go`
  - Register ChatbotService handler with Connect-RPC mux

## Commands to Run

```bash
# Regenerate proto code
buf generate

# Lint proto files
buf lint

# Build Go binary
make build

# Test endpoint (after server running)
curl -X POST http://localhost:8080/altalune.v1.ChatbotService/GetChatbotConfig \
  -H "Content-Type: application/json" \
  -d '{"projectId": "PROJECT_PUBLIC_ID"}'
```

## Validation Checklist

- [ ] `buf lint` passes with no errors
- [ ] `buf generate` creates Go and TypeScript types
- [ ] Service compiles without errors
- [ ] API responds at `/altalune.v1.ChatbotService/GetChatbotConfig`
- [ ] API responds at `/altalune.v1.ChatbotService/UpdateModuleConfig`
- [ ] GetChatbotConfig returns correct JSONB structure
- [ ] UpdateModuleConfig updates specific module only
- [ ] Error codes returned correctly for invalid requests

## Definition of Done

- [ ] Proto schema created and generates correctly
- [ ] All 7 domain files implemented
- [ ] Service registered in container
- [ ] Routes registered in HTTP server
- [ ] GetChatbotConfig endpoint working
- [ ] UpdateModuleConfig endpoint working
- [ ] JSONB handling works correctly
- [ ] Error handling with proper codes
- [ ] Build and lint pass

## Dependencies

- T43: Database schema must exist
- Project domain: Need projectRepo for ID lookup

## Risk Factors

- **Medium Risk**: JSONB handling with google.protobuf.Struct requires careful conversion
- **Low Risk**: Standard 7-file pattern implementation

## Notes

- Use `google.protobuf.Struct` for flexible JSONB handling in proto
- The `modules_config` field contains all module configs in one JSONB column
- UpdateModuleConfig uses PostgreSQL `jsonb_set` for atomic path updates
- Module names are validated against `ValidModuleNames` map
- Project ID lookup is required since frontend uses public IDs
