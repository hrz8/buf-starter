# User Story US11: Chatbot Configuration Foundation

## Story Overview

**As a** project administrator using Altalune
**I want** to configure chatbot modules (LLM, MCP Server, Widget, Prompt) for my project
**So that** I can customize the AI chatbot behavior and appearance for my specific use case

## Key Architecture Decisions

### 1. Database Design: Single Row with JSONB
- **1:1 relationship**: One `altalune_chatbot_configs` row per project
- **No separate `enabled` column at table level** - each module has its own `enabled` flag inside JSONB
- **Structure matches module-editor.json pattern** for consistency

### 2. Schema-Driven Form Rendering (Core Feature)
- JSON Schema defined in frontend determines form structure
- Forms auto-rendered from schema using `SchemaForm.vue` component
- Adding new modules = adding new JSON schema file (no code changes to form components)
- Supports: string, number, boolean, enum/select, arrays, nested objects

### 3. Module-Specific Pages with Dynamic Routing
- Route: `/platform/modules/[name].vue` (e.g., `/platform/modules/llm`)
- Sidebar shows all modules with enabled/disabled indicators
- Each module page has enabled toggle at top + schema-rendered form below

### 4. Simplified Module Configs (Per User Requirements)
- **LLM Module**: `enabled`, `model` (string), `temperature`, `maxToolCalls`
- **MCP Server Module**: `enabled`, `urls[]`, `structuredOutputs[]`
- **Widget Module**: `enabled`, `cors` config
- **Prompt Module**: `enabled`, `systemPrompt` (text only)

## Acceptance Criteria

### Core Functionality

#### View Module List in Sidebar
- **Given** I am logged in and have selected a project
- **When** I view the sidebar
- **Then** I should see "Chatbot" as a parent menu item under Platform group
- **And** under "Chatbot" I see child items: Prompt, MCP Server, Widget
- **And** each item shows enabled/disabled indicator (icon or badge)

#### Navigate to Module Configuration
- **Given** I see the Chatbot menu in sidebar
- **When** I click on a module (e.g., "Prompt")
- **Then** I navigate to `/platform/modules/prompt`
- **And** I see the module configuration page

#### Module Configuration Page Structure
- **Given** I am on a module configuration page
- **When** the page loads
- **Then** I should see:
  - Module name as page title
  - Enabled/Disabled toggle at the top (special component)
  - Schema-driven form rendered below based on module's JSON schema
  - Save button at the bottom

#### Toggle Module Enabled/Disabled
- **Given** I am on a module configuration page
- **When** I toggle the enabled switch
- **Then** the module's `enabled` status is updated in the database
- **And** sidebar indicator updates to reflect new status
- **And** form fields may be disabled when module is disabled (optional UX)

#### Update Module Configuration
- **Given** I am on a module configuration page with enabled module
- **When** I modify form fields and click Save
- **Then** changes are validated against JSON schema
- **And** if valid, changes are persisted to `modules_config` JSONB
- **And** I receive success feedback via toast notification

### Schema-Driven Form Rendering

#### SchemaForm Component Requirements
- **Given** a JSON schema is provided to SchemaForm.vue
- **When** the component renders
- **Then** it generates appropriate form fields based on schema types:
  - `string` -> Input or Textarea (based on `format` or `additionalTypeInfo`)
  - `number` -> Number input with min/max validation
  - `boolean` -> Switch/Toggle component
  - `string` with `enum` -> Select dropdown
  - `array` of objects -> Repeatable form section with add/remove
  - `object` -> Nested fieldset with child fields
  - `string` with `additionalTypeInfo: "json"` -> JSON Editor component

#### Schema Definition Pattern
Each module has a schema file in `frontend/app/components/features/chatbot/schemas/`:
```typescript
// llm.schema.ts
export const llmConfigSchema = {
  type: "object",
  properties: {
    enabled: { type: "boolean", title: "Enabled", default: true },
    model: { type: "string", title: "Model", default: "gpt-4" },
    temperature: { type: "number", title: "Temperature", minimum: 0, maximum: 2, default: 0.7 },
    maxToolCalls: { type: "number", title: "Max Tool Calls", minimum: 1, default: 5 }
  }
};
```

### Auto-Creation

#### Default Configuration on Project Creation
- **Given** a new project is being created
- **When** the project creation completes
- **Then** a default chatbot configuration is automatically created
- **And** all modules have sensible defaults with `enabled: false` (except Prompt: true)
- **And** the configuration is associated with the new project

### Data Validation

#### LLM Module Validation
- Model: Required string when enabled, 1-200 characters (e.g., `us.anthropic.claude-sonnet-4-5-20250929-v1:0`)
- Temperature: Number between 0 and 2, default 0.7
- Max Tool Calls: Number >= 1, default 5

#### MCP Server Module Validation
- URLs array: Each item has `name` (lowercase_snake_case), `url` (valid URL), `apiKey` (optional)
- Structured Outputs array: Each has `target` (format: `{name}__{tool}`), `model`, `prompt`, `inputType`, `outputSchema` (JSON)

#### Widget Module Validation
- CORS allowedOrigins: Array of valid URLs or "*"
- CORS allowedHeaders: Array of strings
- CORS credentials: Boolean

#### Prompt Module Validation
- System Prompt: Required string when enabled, max 50000 characters

### User Experience

#### Responsive Design
- Works on desktop and mobile devices
- Form adapts to screen width
- JSON Editor usable on mobile (may use fullscreen mode)

#### Feedback and Notifications
- Success messages when saving configuration
- Clear error messages for validation failures (per-field)
- Loading states during save operations

#### Navigation
- Sidebar: Platform > Chatbot > [Prompt, MCP Server, Widget]
- Breadcrumb: Platform > Chatbot > {Module Name}

## Technical Requirements

### Backend Architecture

#### Database Schema
```sql
CREATE TABLE altalune_chatbot_configs (
  id BIGINT GENERATED BY DEFAULT AS IDENTITY,
  public_id VARCHAR(20) NOT NULL,
  project_id BIGINT NOT NULL,
  modules_config JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (project_id, id),
  FOREIGN KEY (project_id) REFERENCES altalune_projects (id) ON DELETE CASCADE
) PARTITION BY LIST (project_id);

-- One config per project constraint
CREATE UNIQUE INDEX ON altalune_chatbot_configs (project_id);
CREATE UNIQUE INDEX ON altalune_chatbot_configs (project_id, public_id);
```

**JSONB Structure (modules_config):**
```json
{
  "llm": {
    "enabled": true,
    "model": "us.anthropic.claude-sonnet-4-5-20250929-v1:0",
    "temperature": 0.7,
    "maxToolCalls": 5
  },
  "mcpServer": {
    "enabled": true,
    "urls": [
      { "name": "booking", "url": "https://...", "apiKey": "" }
    ],
    "structuredOutputs": [
      { "target": "booking__search_flights", "model": "...", "prompt": "...", "inputType": "string", "outputSchema": "{...}" }
    ]
  },
  "widget": {
    "enabled": true,
    "cors": {
      "allowedOrigins": ["http://localhost:3000"],
      "allowedHeaders": ["Content-Type"],
      "credentials": true
    }
  },
  "prompt": {
    "enabled": true,
    "systemPrompt": "You are a helpful assistant..."
  }
}
```

- Follow established 7-file domain pattern
- Partitioned by project_id (add to `partitionedTables`)
- Use `google.protobuf.Struct` for JSONB handling in proto

#### gRPC Service Design
- `GetChatbotConfig(project_id)` - Get full configuration
- `UpdateModuleConfig(project_id, module_name, config)` - Update single module

### Frontend Architecture

#### Directory Structure
```
frontend/app/
├── components/features/chatbot/
│   ├── schemas/                    # JSON Schema definitions
│   │   ├── llm.schema.ts
│   │   ├── mcpServer.schema.ts
│   │   ├── widget.schema.ts
│   │   └── prompt.schema.ts
│   ├── SchemaForm.vue              # Generic schema-driven form renderer
│   ├── SchemaField.vue             # Individual field renderer
│   ├── SchemaArrayField.vue        # Array field with add/remove
│   ├── ModuleToggle.vue            # Enabled/disabled toggle header
│   └── index.ts
├── composables/services/
│   └── useChatbotService.ts
├── pages/platform/modules/
│   └── [name].vue                  # Dynamic module page
└── shared/repository/
    └── chatbot.ts
```

#### SchemaForm.vue Responsibilities
1. Receive JSON schema + current values
2. Recursively render SchemaField for each property
3. Handle nested objects and arrays
4. Integrate with vee-validate for validation
5. Emit form values on change

#### Module Page ([name].vue) Flow
1. Get `name` from route params (llm, mcpServer, widget, prompt)
2. Load corresponding schema from `schemas/` directory
3. Fetch current config from API
4. Render ModuleToggle + SchemaForm
5. On save, call `updateModuleConfig(projectId, name, values)`

### API Design

#### Endpoints
- `POST /altalune.v1.ChatbotService/GetChatbotConfig`
- `POST /altalune.v1.ChatbotService/UpdateModuleConfig`

#### Error Codes
- 61001: Chatbot config not found
- 61002: Invalid module configuration
- 61003: Invalid module name

## Out of Scope

- Node Editor functionality (covered in US12)
- Export/Import functionality (covered in US13)
- Real-time configuration preview
- Configuration versioning/history
- Module-level permissions

## Dependencies

- Existing project management functionality
- Database partitioning system
- Frontend UI component library (shadcn-vue)
- JSON Editor component (need to add: vue-json-pretty-editor or similar)

## Definition of Done

- [ ] Database migration created for `altalune_chatbot_configs` table
- [ ] Table added to `partitionedTables` for auto-partition creation
- [ ] Default configuration auto-created when project is created
- [ ] Backend domain implemented (7-file pattern)
- [ ] Protobuf schema defined with Struct for JSONB
- [ ] gRPC service registered in container and routes
- [ ] SchemaForm.vue component implemented (renders forms from JSON schema)
- [ ] SchemaField.vue supports all required field types
- [ ] SchemaArrayField.vue handles array of objects
- [ ] JSON Editor integrated for `additionalTypeInfo: "json"` fields
- [ ] Module schemas defined for: llm, mcpServer, widget, prompt
- [ ] Dynamic page `/platform/modules/[name].vue` implemented
- [ ] ModuleToggle.vue component for enabled/disabled
- [ ] Sidebar navigation shows Chatbot parent with module children
- [ ] Sidebar shows enabled/disabled indicator per module
- [ ] Form validation from schema constraints
- [ ] Error handling with toast notifications
- [ ] Loading states implemented
- [ ] Responsive design verified
- [ ] i18n translations added (en-US, id-ID)

## Notes

- **Schema-first approach**: JSON schemas are the source of truth for both validation and form rendering
- **No provider field**: Use full model ID string (e.g., `us.anthropic.claude-sonnet-4-5-20250929-v1:0`) for flexibility
- **Enabled per module**: No table-level enabled flag; each module manages its own enabled state in JSONB
- **MCP Server UI rendering**: The `ui` config in structuredOutputs (carousel, card) is for webchat app, not dashboard
