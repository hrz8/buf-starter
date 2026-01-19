# T48: Schema-Driven Form Refactoring

## Problem Analysis

### Current Issues

1. **Schemas in wrong location**
   - Currently: `components/features/chatbot/schemas/`
   - Schemas are **data/configuration**, not UI components
   - Violates separation of concerns
   - Makes reuse across features harder

2. **SchemaField.vue has multiple if-else**
   - Single component handles 5+ field types via `v-if`/`v-else-if`
   - Not scalable for new field types
   - Mixed concerns (each type has different props/behavior)
   - Hard to test individual field types

### Current SchemaField.vue Flow
```
fieldType computed → switch on type
  ├── 'switch' → Switch component
  ├── 'number' → Input type="number"
  ├── 'textarea' → Textarea component
  ├── 'select' → Select component
  └── default → Input type="text"
```

---

## Proposed Architecture

### 1. Move Schemas to `lib/modules/`

**New Location:**
```
frontend/lib/
└── modules/
    ├── index.ts              # Barrel exports
    ├── types.ts              # PropertySchema, ModuleSchema
    └── schemas/
        ├── index.ts          # Registry + helpers
        ├── llm.schema.ts
        ├── mcpServer.schema.ts
        ├── widget.schema.ts
        └── prompt.schema.ts
```

**Why `lib/` instead of `shared/`:**
- `shared/` contains repository/API layer
- `lib/` is for pure utilities, types, and configuration
- Follows common Nuxt conventions

### 2. Field Registry Pattern

**New Structure:**
```
frontend/app/components/features/chatbot/
├── fields/
│   ├── index.ts              # Field registry + resolveField()
│   ├── FieldWrapper.vue      # Common FormField wrapper
│   ├── TextField.vue         # type: 'string' (no format/enum)
│   ├── TextareaField.vue     # type: 'string', format: 'textarea'
│   ├── NumberField.vue       # type: 'number'
│   ├── SwitchField.vue       # type: 'boolean'
│   ├── SelectField.vue       # type: 'string', enum: [...]
│   ├── JsonField.vue         # additionalTypeInfo: 'json'
│   ├── ArrayField.vue        # type: 'array' (primitives)
│   └── ObjectArrayField.vue  # type: 'array', items.type: 'object'
├── SchemaField.vue           # Thin resolver (uses registry)
├── SchemaForm.vue            # Form container
├── ModuleConfigForm.vue      # Module-specific form
├── ModuleToggle.vue          # Enable/disable toggle
└── index.ts                  # Barrel exports
```

### 3. Field Registry Implementation

**`fields/index.ts`:**
```typescript
import type { Component } from 'vue';
import type { PropertySchema } from '@/lib/modules';

import TextField from './TextField.vue';
import TextareaField from './TextareaField.vue';
import NumberField from './NumberField.vue';
import SwitchField from './SwitchField.vue';
import SelectField from './SelectField.vue';
import JsonField from './JsonField.vue';
import ArrayField from './ArrayField.vue';
import ObjectArrayField from './ObjectArrayField.vue';

export type FieldType =
  | 'text'
  | 'textarea'
  | 'number'
  | 'switch'
  | 'select'
  | 'json'
  | 'array'
  | 'objectArray';

// Field registry - maps type to component
const FIELD_REGISTRY: Record<FieldType, Component> = {
  text: TextField,
  textarea: TextareaField,
  number: NumberField,
  switch: SwitchField,
  select: SelectField,
  json: JsonField,
  array: ArrayField,
  objectArray: ObjectArrayField,
};

// Resolve schema to field type
export function resolveFieldType(schema: PropertySchema): FieldType {
  // JSON editor takes priority
  if (schema.additionalTypeInfo === 'json') {
    return 'json';
  }

  // Array handling
  if (schema.type === 'array') {
    return schema.items?.type === 'object' ? 'objectArray' : 'array';
  }

  // Boolean
  if (schema.type === 'boolean') {
    return 'switch';
  }

  // Number
  if (schema.type === 'number') {
    return 'number';
  }

  // String variants
  if (schema.type === 'string') {
    if (schema.format === 'textarea') {
      return 'textarea';
    }
    if (schema.enum && schema.enum.length > 0) {
      return 'select';
    }
    return 'text';
  }

  // Default
  return 'text';
}

// Get component for schema
export function resolveFieldComponent(schema: PropertySchema): Component {
  const fieldType = resolveFieldType(schema);
  return FIELD_REGISTRY[fieldType];
}

// Export all field components for direct use if needed
export {
  TextField,
  TextareaField,
  NumberField,
  SwitchField,
  SelectField,
  JsonField,
  ArrayField,
  ObjectArrayField,
};
```

### 4. Simplified SchemaField.vue

```vue
<script setup lang="ts">
import type { PropertySchema } from '@/lib/modules';
import { resolveFieldComponent } from './fields';

const props = defineProps<{
  name: string;
  schema: PropertySchema;
  disabled?: boolean;
}>();

const FieldComponent = computed(() => resolveFieldComponent(props.schema));
</script>

<template>
  <component
    :is="FieldComponent"
    :name="name"
    :schema="schema"
    :disabled="disabled"
  />
</template>
```

### 5. Individual Field Components

Each field component is self-contained with its own FormField wrapper:

**Example: `TextField.vue`**
```vue
<script setup lang="ts">
import type { PropertySchema } from '@/lib/modules';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';

defineProps<{
  name: string;
  schema: PropertySchema;
  disabled?: boolean;
}>();
</script>

<template>
  <FormField v-slot="{ componentField }" :name="name">
    <FormItem>
      <FormLabel>{{ schema.title }}</FormLabel>
      <FormControl>
        <Input
          v-bind="componentField"
          type="text"
          :placeholder="schema.placeholder"
          :maxlength="schema.maxLength"
          :disabled="disabled"
        />
      </FormControl>
      <FormDescription v-if="schema.description">
        {{ schema.description }}
      </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
</template>
```

**Example: `SwitchField.vue`** (different layout)
```vue
<script setup lang="ts">
import type { PropertySchema } from '@/lib/modules';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { Switch } from '@/components/ui/switch';

defineProps<{
  name: string;
  schema: PropertySchema;
  disabled?: boolean;
}>();
</script>

<template>
  <FormField v-slot="{ componentField, value }" :name="name">
    <FormItem class="flex items-center justify-between gap-4">
      <div class="space-y-0.5">
        <FormLabel>{{ schema.title }}</FormLabel>
        <FormDescription v-if="schema.description">
          {{ schema.description }}
        </FormDescription>
      </div>
      <FormControl>
        <Switch
          :model-value="value"
          :disabled="disabled"
          @update:model-value="componentField['onUpdate:modelValue']"
        />
      </FormControl>
    </FormItem>
  </FormField>
</template>
```

---

## Benefits

### Developer Experience (DX)
1. **Single Responsibility** - Each field component handles one type
2. **Easy to Add New Types** - Create component + add to registry
3. **Easy to Test** - Test each field type in isolation
4. **Easy to Customize** - Override specific field without affecting others

### Scalability
1. **New field types** - Just add component + registry entry
2. **Custom schemas** - Can add domain-specific fields (e.g., ColorField, DateField)
3. **Validation** - Can add field-specific validation logic per component

### Maintainability
1. **Clear separation** - Schemas (data) vs Components (UI)
2. **Smaller files** - Each field ~30-50 lines instead of 120+ line monolith
3. **Import clarity** - `@/lib/modules` for schemas, `@/components/...` for UI

---

## Implementation Steps

### Phase 1: Create lib/modules structure
1. Create `frontend/lib/modules/` directory
2. Move types.ts to `lib/modules/types.ts`
3. Move schemas to `lib/modules/schemas/`
4. Update barrel exports

### Phase 2: Create field components
1. Create `fields/` directory
2. Extract TextField.vue from SchemaField.vue
3. Extract TextareaField.vue
4. Extract NumberField.vue
5. Extract SwitchField.vue
6. Extract SelectField.vue
7. Rename JsonEditorField.vue → fields/JsonField.vue
8. Rename SchemaArrayField.vue → fields/ArrayField.vue
9. Rename SchemaObjectArrayField.vue → fields/ObjectArrayField.vue
10. Create field registry (index.ts)

### Phase 3: Simplify existing components
1. Simplify SchemaField.vue to use registry
2. Update SchemaForm.vue to use new field resolution
3. Update imports across codebase

### Phase 4: Cleanup
1. Remove old files
2. Update barrel exports
3. Verify all imports work
4. Run lint + build

---

## Files to Create

```
frontend/lib/modules/
├── index.ts
├── types.ts
└── schemas/
    ├── index.ts
    ├── llm.schema.ts
    ├── mcpServer.schema.ts
    ├── widget.schema.ts
    └── prompt.schema.ts

frontend/app/components/features/chatbot/fields/
├── index.ts
├── TextField.vue
├── TextareaField.vue
├── NumberField.vue
├── SwitchField.vue
├── SelectField.vue
├── JsonField.vue
├── ArrayField.vue
└── ObjectArrayField.vue
```

## Files to Modify

```
frontend/app/components/features/chatbot/SchemaField.vue     # Simplify
frontend/app/components/features/chatbot/SchemaForm.vue      # Update imports
frontend/app/components/features/chatbot/ModuleConfigForm.vue # Update imports
frontend/app/components/features/chatbot/ModuleToggle.vue    # Update imports
frontend/app/components/features/chatbot/index.ts            # Update exports
frontend/app/pages/platform/modules/[name].vue               # Update imports
frontend/app/pages/platform/modules/index.vue                # Update imports
frontend/app/composables/navigation/useNavigationItems.ts    # Update imports
```

## Files to Delete

```
frontend/app/components/features/chatbot/schemas/            # Entire folder
frontend/app/components/features/chatbot/JsonEditorField.vue # Moved to fields/
frontend/app/components/features/chatbot/SchemaArrayField.vue # Moved to fields/
frontend/app/components/features/chatbot/SchemaObjectArrayField.vue # Moved to fields/
```

---

## Estimated Impact

- **New files:** 13
- **Modified files:** 8
- **Deleted files:** 7
- **Net change:** +6 files (but much smaller, focused files)
