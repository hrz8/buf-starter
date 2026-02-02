# Chatbot Module Development Guide

Complete step-by-step guide for creating a new chatbot configuration module.

## Architecture Overview

```
Proto Definition
    ↓ buf generate
┌───────────────────────────────────────────────────┐
│  Generated Code                                    │
│  ├─ Go: gen/chatbot/modules/v1/{module}.pb.go    │
│  ├─ TS: frontend/gen/chatbot/modules/v1/*_pb.ts  │
│  └─ JSON Schema: frontend/gen/jsonschema/*.json  │
└───────────────────────────────────────────────────┘
    ↓
┌─────────────────┐    ┌─────────────────────────────┐
│  Backend        │    │  Frontend                   │
│  ├─ model.go    │    │  ├─ {module}/default.ts    │
│  │  (ValidNames)│    │  ├─ {module}/metadata.ts   │
│  └─ repo.go     │    │  └─ {module}/index.ts      │
│    (defaults)   │    │                             │
└─────────────────┘    └─────────────────────────────┘
```

## Step 1: Define Proto Schema

**File:** `api/proto/chatbot/modules/v1/{module_name}.proto`

```protobuf
syntax = "proto3";

package chatbot.modules.v1;

option go_package = "github.com/hrz8/altalune/gen/chatbot/modules/v1;chatbotmodulesv1";

import "buf/validate/validate.proto";

// {ModuleName}Config defines the configuration for the {module} module.
// Add detailed description here - it appears in JSON Schema.
message {ModuleName}Config {
  // Whether the module is enabled
  bool enabled = 1;

  // Field with string validation
  string some_field = 2 [(buf.validate.field).string.max_len = 500];

  // Field with number validation
  double temperature = 3 [(buf.validate.field).double = {gte: 0, lte: 2}];

  // Field with enum-like validation (string in)
  string provider = 4 [(buf.validate.field).string = {
    in: ["option1", "option2", ""]
  }];

  // Integer with range
  int32 max_value = 5 [(buf.validate.field).int32 = {gte: 1, lte: 100}];
}
```

**Naming Convention:**
- Proto file: `snake_case.proto` (e.g., `mcp_server.proto`)
- Message name: `{PascalCase}Config` (e.g., `McpServerConfig`)
- Module key: `camelCase` (e.g., `mcpServer`)

## Step 2: Generate Code

```bash
buf generate
```

**Generated files:**

| Generated File | Purpose |
|----------------|---------|
| `gen/chatbot/modules/v1/{module}.pb.go` | Go protobuf types |
| `frontend/gen/chatbot/modules/v1/{module}_pb.ts` | TypeScript types |
| `frontend/gen/jsonschema/chatbot.modules.v1.{PascalCase}Config.jsonschema.strict.bundle.json` | JSON Schema with validation rules |

**JSON Schema is critical** - it contains:
- Field types from proto
- Validation constraints from `buf.validate`
- Field descriptions from proto comments
- Used by frontend to auto-generate forms

## Step 3: Backend - Register Module

### 3.1 Add to ValidModuleNames

**File:** `internal/domain/chatbot/model.go`

```go
// ValidModuleNames defines the allowed module names
var ValidModuleNames = map[string]bool{
    "llm":       true,
    "mcpServer": true,
    "widget":    true,
    "prompt":    true,
    "newModule": true,  // ADD YOUR MODULE
}
```

### 3.2 Add Default Config

**File:** `internal/domain/chatbot/repo.go`

Update `defaultChatbotModulesConfig`:

```go
const defaultChatbotModulesConfig = `{
    "llm": { ... },
    "mcpServer": { ... },
    "widget": { ... },
    "prompt": { ... },
    "newModule": {
        "enabled": false,
        "someField": "",
        "temperature": 0.7,
        "provider": "option1",
        "maxValue": 10
    }
}`
```

**Important:**
- Use camelCase for JSON keys (matches proto JSON serialization)
- Provide sensible defaults for all fields
- This is the source of truth for new project configs

## Step 4: Frontend - Create Module Plugin

### 4.1 Create Directory Structure

```
frontend/app/lib/chatbot-modules/{moduleName}/
├── default.ts      # Default values (typed)
├── metadata.ts     # UI hints
└── index.ts        # Re-exports
```

### 4.2 Create default.ts

**File:** `frontend/app/lib/chatbot-modules/{moduleName}/default.ts`

```typescript
import type { MessageInitShape } from '@bufbuild/protobuf';
import type { GenMessage } from '@bufbuild/protobuf/codegenv2';
import type { NewModuleConfig } from '~~/gen/chatbot/modules/v1/{module}_pb';

// Type-safe default using generated proto types
export type NewModuleConfigInit = MessageInitShape<GenMessage<NewModuleConfig>>;

export const newModuleDefaults: NewModuleConfigInit = {
  enabled: false,
  someField: '',
  temperature: 0.7,
  provider: 'option1',
  maxValue: 10,
};
```

**Key points:**
- Use generated types from `~~/gen/chatbot/modules/v1/{module}_pb`
- `MessageInitShape` provides type-safe initialization
- Export name MUST be `{moduleName}Defaults` (auto-discovered)

### 4.3 Create metadata.ts

**File:** `frontend/app/lib/chatbot-modules/{moduleName}/metadata.ts`

```typescript
import type { ModuleMetadata } from '../types';

export const newModuleMetadata: ModuleMetadata = {
  key: 'newModule',           // Must match module key
  title: 'New Module',        // Display title
  description: 'Configure the new module settings.',
  icon: 'Settings',           // Lucide icon name
  fieldOrder: ['enabled', 'provider', 'someField', 'temperature', 'maxValue'],
  fields: {
    enabled: {},  // No override needed - auto-generates title
    provider: {
      title: 'Provider',  // Override auto-generated title
      enum: ['option1', 'option2'],
      enumLabels: {
        option1: 'Option One',
        option2: 'Option Two',
      },
    },
    someField: {
      placeholder: 'Enter value here...',
    },
    temperature: {
      step: 0.1,  // Input step for number field
    },
    maxValue: {
      title: 'Max Value',
      step: 1,
    },
  },
};
```

**Export name MUST be `{moduleName}Metadata`** (auto-discovered)

**FieldMetadata options:**
| Property | Purpose |
|----------|---------|
| `title` | Override auto-generated title |
| `placeholder` | Input placeholder text |
| `format` | `'textarea'` or `'json'` for string fields |
| `step` | Increment for number inputs |
| `enum` | Options for select dropdown |
| `enumLabels` | Human-readable labels for enum values |
| `fieldOrder` | Order nested object fields |
| `titleKey` | For arrays: which field shows in header |

### 4.4 Create index.ts

**File:** `frontend/app/lib/chatbot-modules/{moduleName}/index.ts`

```typescript
export { newModuleDefaults } from './default';
export type { NewModuleConfigInit } from './default';
export { newModuleMetadata } from './metadata';
```

## Step 5: Auto-Discovery (No Manual Registration)

The frontend module system auto-discovers modules:

**`schema.ts`** - Finds all `{module}/index.ts` files and matches with JSON schemas:
```typescript
// Auto-discovers: ./*/index.ts
// Matches with: ../../../gen/jsonschema/chatbot.modules.v1.{PascalCase}Config.jsonschema.strict.bundle.json
```

**`defaults.ts`** - Finds all `{module}Defaults` exports:
```typescript
// Auto-discovers: ./*/index.ts
// Looks for: {moduleName}Defaults export
```

**Convention is critical:**
- Directory name = module key (camelCase)
- Export `{moduleName}Defaults` from default.ts
- Export `{moduleName}Metadata` from metadata.ts
- Proto message = `{PascalCase}Config`

## Step 6: Test

### Backend

```bash
make build
./bin/app serve -c config.yaml
```

Use `mcp__postgres__query` to verify:

```sql
-- Check default config structure
SELECT modules_config FROM altalune_chatbot_configs LIMIT 1;

-- Verify new module exists in config
SELECT modules_config->'newModule' FROM altalune_chatbot_configs LIMIT 1;
```

### Frontend

```bash
cd frontend && pnpm dev
```

Use `mcp__playwright__*` to verify:

1. Navigate to chatbot settings page
2. Check new module appears in list
3. Verify form renders correctly
4. Test save functionality

## Complete Example: Adding "analytics" Module

### 1. Proto

**`api/proto/chatbot/modules/v1/analytics.proto`:**

```protobuf
syntax = "proto3";
package chatbot.modules.v1;
option go_package = "github.com/hrz8/altalune/gen/chatbot/modules/v1;chatbotmodulesv1";

import "buf/validate/validate.proto";

// AnalyticsConfig defines conversation analytics settings.
message AnalyticsConfig {
  bool enabled = 1;
  string tracking_id = 2 [(buf.validate.field).string.max_len = 100];
  bool track_messages = 3;
  int32 retention_days = 4 [(buf.validate.field).int32 = {gte: 1, lte: 365}];
}
```

### 2. Generate

```bash
buf generate
```

### 3. Backend

**`internal/domain/chatbot/model.go`:**
```go
var ValidModuleNames = map[string]bool{
    // ... existing
    "analytics": true,
}
```

**`internal/domain/chatbot/repo.go`:**
```go
const defaultChatbotModulesConfig = `{
    // ... existing
    "analytics": {
        "enabled": false,
        "trackingId": "",
        "trackMessages": true,
        "retentionDays": 30
    }
}`
```

### 4. Frontend

**`frontend/app/lib/chatbot-modules/analytics/default.ts`:**
```typescript
import type { MessageInitShape } from '@bufbuild/protobuf';
import type { GenMessage } from '@bufbuild/protobuf/codegenv2';
import type { AnalyticsConfig } from '~~/gen/chatbot/modules/v1/analytics_pb';

export type AnalyticsConfigInit = MessageInitShape<GenMessage<AnalyticsConfig>>;

export const analyticsDefaults: AnalyticsConfigInit = {
  enabled: false,
  trackingId: '',
  trackMessages: true,
  retentionDays: 30,
};
```

**`frontend/app/lib/chatbot-modules/analytics/metadata.ts`:**
```typescript
import type { ModuleMetadata } from '../types';

export const analyticsMetadata: ModuleMetadata = {
  key: 'analytics',
  title: 'Analytics',
  description: 'Track conversation analytics and metrics.',
  icon: 'BarChart3',
  fieldOrder: ['enabled', 'trackingId', 'trackMessages', 'retentionDays'],
  fields: {
    enabled: {},
    trackingId: {
      title: 'Tracking ID',
      placeholder: 'UA-XXXXX-X',
    },
    trackMessages: {
      title: 'Track Messages',
    },
    retentionDays: {
      title: 'Retention Days',
      step: 1,
    },
  },
};
```

**`frontend/app/lib/chatbot-modules/analytics/index.ts`:**
```typescript
export { analyticsDefaults } from './default';
export type { AnalyticsConfigInit } from './default';
export { analyticsMetadata } from './metadata';
```

## Troubleshooting

| Issue | Cause | Fix |
|-------|-------|-----|
| Module not appearing | Export naming wrong | Ensure `{moduleName}Defaults` and `{moduleName}Metadata` |
| Form not rendering | JSON Schema not found | Check proto message is `{PascalCase}Config` |
| Type errors | Missing proto generation | Run `buf generate` |
| Validation not working | Missing buf.validate | Add validation rules to proto |
| Update fails | Module not in ValidModuleNames | Add to model.go |
