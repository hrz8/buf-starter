# Task T50: Node Backend Domain Implementation

**Story Reference:** US12-node-editor.md
**Type:** Backend
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T49 (Node Database Schema)

## Objective

Implement complete chatbot_node domain following the 7-file pattern with proto definitions, all CRUD operations, and service registration.

## Acceptance Criteria

- [ ] Proto schema defines ListNodes, CreateNode, GetNode, UpdateNode, DeleteNode RPCs
- [ ] Domain implements 7-file pattern (model, interface, repo, service, handler, mapper, errors)
- [ ] ListNodes returns nodes sorted alphabetically for sidebar
- [ ] CreateNode validates name format (lowercase_snake_case), language, triggers, messages
- [ ] UpdateNode allows editing name, tags, enabled, triggers, messages (not lang)
- [ ] DeleteNode removes node permanently
- [ ] Service registered in container and HTTP routes
- [ ] All API endpoints functional via Connect-RPC

## Technical Requirements

### Proto Schema (chatbot_node.proto)

```protobuf
syntax = "proto3";

package altalune.v1;

option go_package = "github.com/hrz8/altalune/gen/altalune/v1;altalunev1";

import "google/protobuf/timestamp.proto";
import "buf/validate/validate.proto";

service ChatbotNodeService {
  rpc ListNodes(ListNodesRequest) returns (ListNodesResponse) {}
  rpc CreateNode(CreateNodeRequest) returns (CreateNodeResponse) {}
  rpc GetNode(GetNodeRequest) returns (GetNodeResponse) {}
  rpc UpdateNode(UpdateNodeRequest) returns (UpdateNodeResponse) {}
  rpc DeleteNode(DeleteNodeRequest) returns (DeleteNodeResponse) {}
}

message ChatbotNode {
  string id = 1;                              // public_id
  string name = 2;
  string lang = 3;
  repeated string tags = 4;
  bool enabled = 5;
  repeated Trigger triggers = 6;
  repeated Message messages = 7;
  google.protobuf.Timestamp created_at = 98;
  google.protobuf.Timestamp updated_at = 99;
}

message Trigger {
  string type = 1;   // keyword, contains, regex, equals
  string value = 2;
}

message Message {
  string role = 1;   // assistant
  string content = 2;
}

message ListNodesRequest {
  string project_id = 1 [(buf.validate.field).required = true];
}

message ListNodesResponse {
  repeated ChatbotNode nodes = 1;
}

message CreateNodeRequest {
  string project_id = 1 [(buf.validate.field).required = true];
  string name = 2 [(buf.validate.field).required = true, (buf.validate.field).string = {
    min_len: 2,
    max_len: 100,
    pattern: "^[a-z][a-z0-9_]*$"
  }];
  string lang = 3 [(buf.validate.field).required = true, (buf.validate.field).string = {
    in: ["en-US", "id-ID"]
  }];
  repeated string tags = 4;
}

message CreateNodeResponse {
  ChatbotNode node = 1;
}

message GetNodeRequest {
  string project_id = 1 [(buf.validate.field).required = true];
  string node_id = 2 [(buf.validate.field).required = true];
}

message GetNodeResponse {
  ChatbotNode node = 1;
}

message UpdateNodeRequest {
  string project_id = 1 [(buf.validate.field).required = true];
  string node_id = 2 [(buf.validate.field).required = true];
  string name = 3 [(buf.validate.field).string = {
    min_len: 2,
    max_len: 100,
    pattern: "^[a-z][a-z0-9_]*$"
  }];
  repeated string tags = 4;
  bool enabled = 5;
  repeated Trigger triggers = 6;
  repeated Message messages = 7;
}

message UpdateNodeResponse {
  ChatbotNode node = 1;
  string message = 2;
}

message DeleteNodeRequest {
  string project_id = 1 [(buf.validate.field).required = true];
  string node_id = 2 [(buf.validate.field).required = true];
}

message DeleteNodeResponse {
  string message = 1;
}
```

### Validation Rules

- **name**: required, 2-100 chars, pattern `^[a-z][a-z0-9_]*$`
- **lang**: required, enum `["en-US", "id-ID"]`
- **triggers**: at least 1 required, each value max 500 chars
- **messages**: at least 1 required, each content max 5000 chars

### Error Codes

- 62001: Node not found
- 62002: Invalid node name format
- 62003: Invalid language
- 62004: Duplicate node name+lang
- 62005: At least one trigger required
- 62006: At least one message required
- 62007: Invalid trigger type
- 62008: Invalid regex pattern

### 7-File Pattern Implementation

#### model.go

```go
package chatbot_node

type ChatbotNode struct {
    ID        string
    ProjectID int64
    Name      string
    Lang      string
    Tags      []string
    Enabled   bool
    Triggers  []Trigger
    Messages  []Message
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Trigger struct {
    Type  string // keyword, contains, regex, equals
    Value string
}

type Message struct {
    Role    string // assistant
    Content string
}

// ValidTriggerTypes contains all valid trigger types
var ValidTriggerTypes = map[string]bool{
    "keyword":  true,
    "contains": true,
    "regex":    true,
    "equals":   true,
}

// ValidLanguages contains all valid language codes
var ValidLanguages = map[string]bool{
    "en-US": true,
    "id-ID": true,
}
```

#### interface.go

```go
package chatbot_node

type Repositor interface {
    List(ctx context.Context, projectID int64) ([]*ChatbotNode, error)
    GetByID(ctx context.Context, projectID int64, publicID string) (*ChatbotNode, error)
    Create(ctx context.Context, input *CreateNodeInput) (*ChatbotNode, error)
    Update(ctx context.Context, input *UpdateNodeInput) (*ChatbotNode, error)
    Delete(ctx context.Context, projectID int64, publicID string) error
}
```

#### errors.go

```go
package chatbot_node

var (
    ErrNodeNotFound       = errors.New("chatbot node not found")
    ErrInvalidNodeName    = errors.New("invalid node name format")
    ErrInvalidLanguage    = errors.New("invalid language")
    ErrDuplicateNodeName  = errors.New("duplicate node name and language")
    ErrNoTriggers         = errors.New("at least one trigger required")
    ErrNoMessages         = errors.New("at least one message required")
    ErrInvalidTriggerType = errors.New("invalid trigger type")
    ErrInvalidRegex       = errors.New("invalid regex pattern")
)
```

## Files to Create

```
api/proto/altalune/v1/chatbot_node.proto

internal/domain/chatbot_node/
├── model.go       # ChatbotNode, Trigger, Message structs
├── interface.go   # Repositor interface
├── errors.go      # Domain-specific errors
├── mapper.go      # Proto <-> Domain conversions
├── repo.go        # Database operations
├── service.go     # Business logic, validation
└── handler.go     # Connect-RPC handler
```

## Files to Modify

- `internal/container/container.go` - Add chatbot_nodeRepo and chatbot_nodeService fields
- `internal/container/getter.go` - Add `GetChatbotNodeService()` method
- `internal/server/http_routes.go` - Register ChatbotNodeService handler

## Commands to Run

```bash
# Regenerate proto code
buf generate

# Lint proto files
buf lint

# Build Go binary
make build
```

## Validation Checklist

- [ ] `buf lint` passes
- [ ] Service compiles without errors
- [ ] ListNodes returns nodes sorted alphabetically
- [ ] CreateNode validates all fields and creates node
- [ ] GetNode returns full node with triggers/messages
- [ ] UpdateNode persists changes correctly
- [ ] DeleteNode removes node from database

## Definition of Done

- [ ] Proto schema created and generates correctly
- [ ] All 7 domain files implemented
- [ ] Service registered in container
- [ ] Routes registered in HTTP server
- [ ] All CRUD endpoints working
- [ ] Proper error codes for all failure cases
- [ ] JSONB triggers/messages handled correctly
- [ ] Build and lint pass

## Dependencies

- T49: Database schema must exist
- Project domain: Need projectRepo for ID lookup

## Risk Factors

- **Medium Risk**: JSONB array handling for triggers/messages requires careful parsing
- **Low Risk**: Standard 7-file pattern implementation

## Notes

- ListNodes should return only essential fields (id, name, lang, enabled) for sidebar efficiency
- GetNode returns full data including triggers/messages
- Triggers array stored as JSONB, parsed to []Trigger on read
- Messages follow OpenAI format with role="assistant"
- Regex patterns validated server-side before save
