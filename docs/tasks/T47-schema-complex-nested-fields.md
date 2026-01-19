# T47: Schema-Driven Complex Nested Fields

## Overview

Enhance the schema-driven form system to support complex nested structures including:
- Arrays of objects (not just strings)
- Nested objects within arrays
- JSON editor fields
- Deep nesting (3-4 levels)

## Current State

The current implementation handles:
- ✅ Simple fields (string, number, boolean, enum)
- ✅ Arrays of strings
- ✅ Single-level nested objects
- ❌ Arrays of objects
- ❌ JSON editor fields
- ❌ Deep nested structures (object → array → object → field)

## Target Structure (from module-editor.json)

```
mcpServer.structuredOutputs[] (array of objects)
├── target: string
├── model: string (enum)
├── prompt: string (textarea)
├── inputType: string (enum)
├── outputSchema: string (JSON editor)
└── ui: object (nested)
    ├── enabled: boolean
    ├── component: string (enum)
    ├── isArray: boolean
    ├── dataPath: string
    └── itemConfig: string (JSON editor)
```

## Implementation Plan

### Phase 1: Update Schema Types & Definitions

**1.1 Update `types.ts`**
- Add support for `items` being an object schema (for array of objects)
- Ensure recursive property definitions work

**1.2 Update `mcpServer.schema.ts`**
- Define full `urls` structure (array of objects with name, url, apiKey)
- Define full `structuredOutputs` structure with nested `ui` object

### Phase 2: New Components

**2.1 Create `JsonEditorField.vue`**
- Textarea-based JSON editor with validation
- Real-time JSON syntax checking
- Error message display
- Optional: Pretty-print/format button

**2.2 Create `SchemaObjectArrayField.vue`**
- Handles arrays where `items.type === 'object'`
- Collapsible items for better UX
- Each item shows a summary/title
- Add/Remove functionality
- Renders nested properties for each item

### Phase 3: Update Existing Components

**3.1 Update `SchemaField.vue`**
- Add case for `additionalTypeInfo === 'json'` → render `JsonEditorField`

**3.2 Update `SchemaArrayField.vue`**
- Detect if `items.type === 'object'`
- If object: delegate to `SchemaObjectArrayField`
- If primitive: current behavior

**3.3 Update `SchemaForm.vue`**
- Ensure recursive rendering works for nested objects within arrays

### Phase 4: UX Enhancements

**4.1 Collapsible Array Items**
- Use Collapsible component for array of objects
- Show item summary when collapsed (e.g., "booking__search_flights")

**4.2 Visual Hierarchy**
- Clear indentation for nesting levels
- Borders/cards to separate array items
- Section headers for nested objects

## Component Architecture

```
SchemaForm.vue
├── SchemaField.vue
│   └── Handles: string, number, boolean, enum
│   └── NEW: additionalTypeInfo === 'json' → JsonEditorField
├── SchemaArrayField.vue
│   ├── items.type === 'string' → Current behavior (Input per item)
│   └── items.type === 'object' → SchemaObjectArrayField
├── SchemaObjectArrayField.vue (NEW)
│   └── Renders array of objects with:
│       ├── Collapsible item headers
│       ├── Add/Remove controls
│       └── Nested property rendering (recursive)
└── JsonEditorField.vue (NEW)
    └── JSON editing with validation
```

## Schema Definition Example

```typescript
// mcpServer.schema.ts
export const mcpServerSchema: ModuleSchema = {
  key: 'mcpServer',
  title: 'MCP Server',
  // ...
  properties: {
    urls: {
      type: 'array',
      title: 'MCP Server URLs',
      items: {
        type: 'object',
        title: 'Server',
        properties: {
          name: { type: 'string', title: 'Name' },
          url: { type: 'string', title: 'URL' },
          apiKey: { type: 'string', title: 'API Key' },
        },
      },
    },
    structuredOutputs: {
      type: 'array',
      title: 'Structured Outputs',
      items: {
        type: 'object',
        title: 'Output Configuration',
        // Used for collapsible header display
        titleKey: 'target',
        properties: {
          target: { type: 'string', title: 'Target' },
          model: {
            type: 'string',
            title: 'Model',
            enum: ['openai.gpt-oss-120b-1:0'],
          },
          prompt: {
            type: 'string',
            title: 'Prompt',
            format: 'textarea',
          },
          inputType: {
            type: 'string',
            title: 'Input Type',
            enum: ['string', 'json'],
          },
          outputSchema: {
            type: 'string',
            title: 'Output Schema',
            additionalTypeInfo: 'json',
          },
          ui: {
            type: 'object',
            title: 'UI Configuration',
            properties: {
              enabled: { type: 'boolean', title: 'Enable UI', default: false },
              component: {
                type: 'string',
                title: 'Component',
                enum: ['carousel', 'card'],
                default: 'card',
              },
              isArray: { type: 'boolean', title: 'Is Array', default: false },
              dataPath: { type: 'string', title: 'Data Path' },
              itemConfig: {
                type: 'string',
                title: 'Item Config',
                additionalTypeInfo: 'json',
              },
            },
          },
        },
      },
    },
  },
};
```

## Task Breakdown

### T47.1: Schema Types & Definitions (1-2 hours)
- [ ] Update `types.ts` for recursive/nested support
- [ ] Update `mcpServer.schema.ts` with full structure
- [ ] Verify TypeScript compiles

### T47.2: JsonEditorField Component (2-3 hours)
- [ ] Create `JsonEditorField.vue`
- [ ] JSON validation with error display
- [ ] Integrate with vee-validate (useField)
- [ ] Basic formatting/pretty-print

### T47.3: SchemaObjectArrayField Component (3-4 hours)
- [ ] Create `SchemaObjectArrayField.vue`
- [ ] Collapsible items with Collapsible component
- [ ] Dynamic title from `titleKey` or first string field
- [ ] Add/Remove item functionality
- [ ] Recursive property rendering

### T47.4: Update Existing Components (1-2 hours)
- [ ] Update `SchemaField.vue` for JSON fields
- [ ] Update `SchemaArrayField.vue` to detect and delegate
- [ ] Update `SchemaForm.vue` if needed

### T47.5: Testing & Polish (1-2 hours)
- [ ] Test with MCP Server module
- [ ] Verify data persistence (save/load)
- [ ] Visual polish and responsive design
- [ ] Handle edge cases (empty arrays, validation)

## Acceptance Criteria

1. **MCP Server URLs**: Can add/edit/remove server objects with name, url, apiKey
2. **Structured Outputs**: Can add/edit/remove output configurations with all fields
3. **Nested UI Object**: Can configure the `ui` sub-object within each structured output
4. **JSON Editor**: outputSchema and itemConfig render as JSON editors with validation
5. **Data Persistence**: All nested data saves and loads correctly via API
6. **Collapsible UX**: Array items are collapsible with meaningful headers
7. **Validation**: JSON fields show validation errors, required fields enforced

## Dependencies

- Existing: vee-validate, shadcn-vue components (Collapsible, Card, Input, etc.)
- No new external dependencies required (JSON editor is textarea-based)

## Future Enhancements (Out of Scope)

- Monaco/CodeMirror for advanced JSON editing
- Drag-and-drop reordering of array items
- Import/Export JSON for entire configurations
- JSON Schema validation against outputSchema definitions
