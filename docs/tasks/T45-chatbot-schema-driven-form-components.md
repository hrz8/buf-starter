# Task T45: Chatbot Schema-Driven Form Components

**Story Reference:** US11-chatbot-configuration-foundation.md
**Type:** Frontend Foundation
**Priority:** High
**Estimated Effort:** 4-6 hours
**Prerequisites:** T44 (Backend Domain)

## Objective

Implement reusable schema-driven form rendering system that auto-generates forms from JSON schema definitions. This creates the foundation for all chatbot module configuration UIs.

## Acceptance Criteria

- [ ] SchemaForm.vue renders forms from JSON schema
- [ ] SchemaField.vue handles string, number, boolean, enum types
- [ ] SchemaArrayField.vue handles array of objects with add/remove
- [ ] JSON Editor component integrated for `additionalTypeInfo: "json"` fields
- [ ] Module schemas defined for: llm, mcpServer, widget, prompt
- [ ] Forms integrate with vee-validate for validation

## Technical Requirements

### Schema Structure Pattern

Each module schema follows this structure:

```typescript
// types/schema.ts
export interface JSONSchemaProperty {
  type: 'string' | 'number' | 'boolean' | 'array' | 'object';
  title: string;
  description?: string;
  default?: any;
  enum?: string[];
  minimum?: number;
  maximum?: number;
  minLength?: number;
  maxLength?: number;
  format?: 'textarea' | 'json';
  additionalTypeInfo?: 'json';
  items?: JSONSchemaProperty | { type: 'object'; properties: Record<string, JSONSchemaProperty> };
  properties?: Record<string, JSONSchemaProperty>;
  required?: string[];
}

export interface ModuleSchema {
  key: string;
  title: string;
  icon: string;
  type: 'object';
  properties: Record<string, JSONSchemaProperty>;
  required?: string[];
}
```

### Schema Registry

File: `frontend/app/components/features/chatbot/schemas/index.ts`

```typescript
import { llmSchema } from './llm.schema';
import { mcpServerSchema } from './mcpServer.schema';
import { widgetSchema } from './widget.schema';
import { promptSchema } from './prompt.schema';

export const MODULE_SCHEMAS = {
  llm: llmSchema,
  mcpServer: mcpServerSchema,
  widget: widgetSchema,
  prompt: promptSchema,
} as const;

export type ModuleName = keyof typeof MODULE_SCHEMAS;

export function getModuleSchema(name: string): ModuleSchema | undefined {
  return MODULE_SCHEMAS[name as ModuleName];
}

export function getModuleList(): Array<{ key: ModuleName; title: string; icon: string }> {
  return Object.entries(MODULE_SCHEMAS).map(([key, schema]) => ({
    key: key as ModuleName,
    title: schema.title,
    icon: schema.icon,
  }));
}
```

### Module Schema Definitions

#### llm.schema.ts

```typescript
import type { ModuleSchema } from '../types/schema';

export const llmSchema: ModuleSchema = {
  key: 'llm',
  title: 'LLM Configuration',
  icon: 'lucide:brain',
  type: 'object',
  properties: {
    model: {
      type: 'string',
      title: 'Model',
      description: 'LLM model identifier (e.g., gpt-4, claude-3)',
      maxLength: 200,
    },
    temperature: {
      type: 'number',
      title: 'Temperature',
      description: 'Controls randomness (0 = deterministic, 2 = creative)',
      minimum: 0,
      maximum: 2,
      default: 0.7,
    },
    maxToolCalls: {
      type: 'number',
      title: 'Max Tool Calls',
      description: 'Maximum number of tool calls per conversation turn',
      minimum: 1,
      maximum: 50,
      default: 5,
    },
  },
  required: ['model'],
};
```

#### mcpServer.schema.ts

```typescript
import type { ModuleSchema } from '../types/schema';

export const mcpServerSchema: ModuleSchema = {
  key: 'mcpServer',
  title: 'MCP Server',
  icon: 'lucide:server',
  type: 'object',
  properties: {
    urls: {
      type: 'array',
      title: 'Server URLs',
      description: 'MCP server endpoint URLs',
      items: {
        type: 'string',
      },
    },
    structuredOutputs: {
      type: 'array',
      title: 'Structured Outputs',
      description: 'JSON schema definitions for structured outputs',
      items: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            title: 'Output Name',
          },
          schema: {
            type: 'string',
            title: 'JSON Schema',
            additionalTypeInfo: 'json',
          },
        },
      },
    },
  },
};
```

#### widget.schema.ts

```typescript
import type { ModuleSchema } from '../types/schema';

export const widgetSchema: ModuleSchema = {
  key: 'widget',
  title: 'Widget Configuration',
  icon: 'lucide:layout',
  type: 'object',
  properties: {
    cors: {
      type: 'object',
      title: 'CORS Settings',
      properties: {
        allowedOrigins: {
          type: 'array',
          title: 'Allowed Origins',
          description: 'List of allowed origins for CORS',
          items: {
            type: 'string',
          },
        },
        allowedHeaders: {
          type: 'array',
          title: 'Allowed Headers',
          description: 'List of allowed headers',
          items: {
            type: 'string',
          },
        },
        credentials: {
          type: 'boolean',
          title: 'Allow Credentials',
          description: 'Whether to allow credentials in CORS requests',
          default: false,
        },
      },
    },
  },
};
```

#### prompt.schema.ts

```typescript
import type { ModuleSchema } from '../types/schema';

export const promptSchema: ModuleSchema = {
  key: 'prompt',
  title: 'Prompt Configuration',
  icon: 'lucide:message-square',
  type: 'object',
  properties: {
    systemPrompt: {
      type: 'string',
      title: 'System Prompt',
      description: 'The system prompt that defines the chatbot personality and behavior',
      format: 'textarea',
      maxLength: 10000,
      default: 'You are a helpful assistant.',
    },
  },
  required: ['systemPrompt'],
};
```

### Component Implementations

#### SchemaForm.vue

```vue
<script setup lang="ts">
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { computed, watch } from 'vue';
import { z } from 'zod';
import type { ModuleSchema } from './types/schema';
import SchemaField from './SchemaField.vue';
import { buildZodSchema } from './utils/schemaToZod';

const props = defineProps<{
  schema: ModuleSchema;
  modelValue: Record<string, any>;
  disabled?: boolean;
  loading?: boolean;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: Record<string, any>];
  'submit': [value: Record<string, any>];
}>();

// Build Zod schema from JSON schema
const zodSchema = computed(() => buildZodSchema(props.schema));

const form = useForm({
  validationSchema: toTypedSchema(zodSchema.value),
  initialValues: props.modelValue,
});

// Sync external changes
watch(() => props.modelValue, (newVal) => {
  form.resetForm({ values: newVal });
}, { deep: true });

// Emit changes
watch(() => form.values, (newVal) => {
  emit('update:modelValue', newVal);
}, { deep: true });

const onSubmit = form.handleSubmit((values) => {
  emit('submit', values);
});

defineExpose({ form, validate: form.validate, resetForm: form.resetForm });
</script>

<template>
  <form class="space-y-6" @submit="onSubmit">
    <SchemaField
      v-for="(property, key) in schema.properties"
      :key="key"
      :name="String(key)"
      :property="property"
      :disabled="disabled || loading"
    />
    <slot name="actions" :loading="loading" :disabled="disabled" />
  </form>
</template>
```

#### SchemaField.vue

```vue
<script setup lang="ts">
import { computed } from 'vue';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import SchemaArrayField from './SchemaArrayField.vue';
import JsonEditor from './JsonEditor.vue';
import type { JSONSchemaProperty } from './types/schema';

const props = defineProps<{
  name: string;
  property: JSONSchemaProperty;
  disabled?: boolean;
}>();

const fieldType = computed(() => {
  if (props.property.additionalTypeInfo === 'json') return 'json';
  if (props.property.enum) return 'select';
  if (props.property.type === 'array') return 'array';
  if (props.property.type === 'object') return 'object';
  if (props.property.type === 'boolean') return 'switch';
  if (props.property.type === 'number') return 'number';
  if (props.property.format === 'textarea') return 'textarea';
  return 'text';
});
</script>

<template>
  <!-- Array Field -->
  <SchemaArrayField
    v-if="fieldType === 'array'"
    :name="name"
    :property="property"
    :disabled="disabled"
  />

  <!-- Nested Object -->
  <fieldset
    v-else-if="fieldType === 'object' && property.properties"
    class="space-y-4 rounded-lg border p-4"
  >
    <legend class="px-2 text-sm font-medium">{{ property.title }}</legend>
    <SchemaField
      v-for="(nestedProp, nestedKey) in property.properties"
      :key="nestedKey"
      :name="`${name}.${nestedKey}`"
      :property="nestedProp"
      :disabled="disabled"
    />
  </fieldset>

  <!-- Standard Fields -->
  <FormField v-else v-slot="{ componentField }" :name="name">
    <FormItem>
      <FormLabel>{{ property.title }}</FormLabel>
      <FormControl>
        <!-- Switch for boolean -->
        <Switch
          v-if="fieldType === 'switch'"
          :checked="componentField.modelValue"
          :disabled="disabled"
          @update:checked="componentField['onUpdate:modelValue']"
        />

        <!-- Select for enum -->
        <Select
          v-else-if="fieldType === 'select'"
          :model-value="componentField.modelValue"
          :disabled="disabled"
          @update:model-value="componentField['onUpdate:modelValue']"
        >
          <SelectTrigger>
            <SelectValue :placeholder="`Select ${property.title}`" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="opt in property.enum" :key="opt" :value="opt">
              {{ opt }}
            </SelectItem>
          </SelectContent>
        </Select>

        <!-- JSON Editor -->
        <JsonEditor
          v-else-if="fieldType === 'json'"
          :model-value="componentField.modelValue"
          :disabled="disabled"
          @update:model-value="componentField['onUpdate:modelValue']"
        />

        <!-- Textarea -->
        <Textarea
          v-else-if="fieldType === 'textarea'"
          v-bind="componentField"
          :disabled="disabled"
          :placeholder="property.description"
          rows="6"
        />

        <!-- Number Input -->
        <Input
          v-else-if="fieldType === 'number'"
          v-bind="componentField"
          type="number"
          :disabled="disabled"
          :placeholder="property.description"
          :min="property.minimum"
          :max="property.maximum"
        />

        <!-- Text Input (default) -->
        <Input
          v-else
          v-bind="componentField"
          type="text"
          :disabled="disabled"
          :placeholder="property.description"
          :maxlength="property.maxLength"
        />
      </FormControl>
      <FormDescription v-if="property.description">
        {{ property.description }}
      </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
</template>
```

#### SchemaArrayField.vue

Handles arrays with add/remove functionality.

#### JsonEditor.vue

Monaco editor or simple textarea for JSON editing with validation.

#### ModuleToggle.vue

Header component with enabled/disabled toggle.

```vue
<script setup lang="ts">
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';
import { Icon } from '@/components/ui/icon';

defineProps<{
  title: string;
  icon: string;
  enabled: boolean;
  loading?: boolean;
}>();

const emit = defineEmits<{
  'update:enabled': [value: boolean];
}>();
</script>

<template>
  <div class="flex items-center justify-between rounded-lg border p-4">
    <div class="flex items-center gap-3">
      <Icon :name="icon" class="h-5 w-5" />
      <div>
        <h3 class="font-medium">{{ title }}</h3>
        <Badge :variant="enabled ? 'default' : 'secondary'" class="mt-1">
          {{ enabled ? 'Enabled' : 'Disabled' }}
        </Badge>
      </div>
    </div>
    <Switch
      :checked="enabled"
      :disabled="loading"
      @update:checked="$emit('update:enabled', $event)"
    />
  </div>
</template>
```

### Chatbot Repository

File: `frontend/shared/repository/chatbot.ts`

```typescript
import type { GetChatbotConfigRequest, UpdateModuleConfigRequest } from '@/gen/altalune/v1/chatbot_pb';
import { createClient } from '../../utils/connect';

const client = createClient('ChatbotService');

export async function getChatbotConfig(projectId: string) {
  return client.getChatbotConfig({ projectId });
}

export async function updateModuleConfig(
  projectId: string,
  moduleName: string,
  config: Record<string, any>
) {
  return client.updateModuleConfig({
    projectId,
    moduleName,
    config: config as any,
  });
}
```

### Service Composable

File: `frontend/app/composables/services/useChatbotService.ts`

```typescript
import { ref } from 'vue';
import * as chatbotRepo from '@shared/repository/chatbot';
import type { ModuleName } from '@/components/features/chatbot/schemas';

export function useChatbotService(projectId: string) {
  const config = ref<Record<string, any>>({});
  const loading = ref(true);
  const error = ref<Error | null>(null);

  async function fetchConfig() {
    loading.value = true;
    error.value = null;
    try {
      const response = await chatbotRepo.getChatbotConfig(projectId);
      config.value = response.modulesConfig?.toJson() as Record<string, any> || {};
    } catch (e) {
      error.value = e as Error;
    } finally {
      loading.value = false;
    }
  }

  async function updateModule(moduleName: ModuleName, moduleConfig: Record<string, any>) {
    loading.value = true;
    error.value = null;
    try {
      const response = await chatbotRepo.updateModuleConfig(
        projectId,
        moduleName,
        moduleConfig
      );
      // Update local config with response
      config.value = response.updatedConfig?.toJson() as Record<string, any> || config.value;
      return response;
    } catch (e) {
      error.value = e as Error;
      throw e;
    } finally {
      loading.value = false;
    }
  }

  function getModuleConfig(moduleName: ModuleName) {
    return config.value[moduleName] || {};
  }

  return {
    config,
    loading,
    error,
    fetchConfig,
    updateModule,
    getModuleConfig,
  };
}
```

## Files to Create

```
frontend/app/components/features/chatbot/
├── schemas/
│   ├── index.ts
│   ├── llm.schema.ts
│   ├── mcpServer.schema.ts
│   ├── widget.schema.ts
│   └── prompt.schema.ts
├── types/
│   └── schema.ts
├── utils/
│   └── schemaToZod.ts
├── SchemaForm.vue
├── SchemaField.vue
├── SchemaArrayField.vue
├── JsonEditor.vue
├── ModuleToggle.vue
├── schema.ts
├── error.ts
├── constants.ts
└── index.ts

frontend/shared/repository/
└── chatbot.ts

frontend/app/composables/services/
└── useChatbotService.ts
```

## Commands to Run

```bash
cd frontend

# Install dependencies if needed
pnpm install

# Run lint to check for errors
pnpm lint

# Run dev server to test components
pnpm dev
```

## Validation Checklist

- [ ] SchemaForm renders all field types correctly
- [ ] String fields render as text input
- [ ] Number fields render with min/max constraints
- [ ] Boolean fields render as switches
- [ ] Textarea fields render with proper sizing
- [ ] Array fields support add/remove
- [ ] Nested objects render as fieldsets
- [ ] JSON editor validates JSON syntax
- [ ] Form validation errors display per-field
- [ ] Forms emit update events correctly
- [ ] Repository methods work with backend API

## Definition of Done

- [ ] All 4 module schemas defined
- [ ] SchemaForm component working
- [ ] SchemaField handles all property types
- [ ] SchemaArrayField handles arrays
- [ ] JsonEditor component working
- [ ] ModuleToggle component working
- [ ] Repository layer created
- [ ] Service composable created
- [ ] TypeScript types correct
- [ ] No lint errors

## Dependencies

- T44: Backend API must be available
- Generated TypeScript types from buf generate
- vee-validate, zod for validation
- shadcn-vue components

## Risk Factors

- **Medium Risk**: Schema-to-Zod conversion complexity
- **Medium Risk**: vee-validate integration with dynamic schemas
- **Low Risk**: Standard component patterns

## Notes

- Schema registry is single source of truth for module definitions
- No hardcoded module list in navigation - derive from MODULE_SCHEMAS
- i18n translations are for UI chrome only, not module titles (per user requirement)
- Follow vee-validate FormField best practices from FRONTEND_GUIDE.md
- Loading state should start as `true` to avoid provide/inject errors
