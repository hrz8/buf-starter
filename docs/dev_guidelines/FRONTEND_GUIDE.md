## Frontend Development Workflow

This systematic workflow ensures consistency, type safety, and proper architecture adherence when implementing new frontend features.

## Code Formatting

**Manual Format Command:**

If format on save is not enabled, always run before committing:

```bash
cd frontend && pnpm lint:fix
```

**Why This Matters:**

- The project uses ESLint with strict rules for Vue component structure
- Auto-formatting prevents linting errors during development
- Ensures consistent code style across the team
- Required for CI/CD pipeline to pass

### Frontend Architecture Overview

**Layer Structure:**

1. **Repository Layer**: `frontend/shared/repository/` - Connect-RPC client wrappers
2. **Service Layer**: `frontend/app/composables/services/` - Business logic and state management
3. **Component Layer**: `frontend/app/components/features/` - Domain-specific UI components
4. **Page Layer**: `frontend/app/pages/` - Route-level components
5. **Store Layer**: `frontend/app/stores/` - Global state management (Pinia)

### Step 1: Frontend Reusability Analysis

**✅ Infrastructure Components (Reuse These):**

- `useConnectValidator` - Protobuf validation with error mapping
- `useErrorMessage` - Connect error parsing with i18n support
- `useQueryRequest` - Query parameter management for data tables
- DataTable components and utilities (`@/components/custom/datatable`)
- Form components (shadcn-vue based: `@/components/ui/form`)
- Sheet/Modal components for overlays (`@/components/ui/sheet`)
- Toast notifications (`vue-sonner`)
- Loading and empty state components

**✅ Architecture Patterns (Follow These):**

- Repository pattern: Connect client wrappers with error handling
- Service composables: Reactive state + validation + API calls
- Form validation: Protobuf schemas with `useConnectValidator`
- Error handling: `useErrorMessage` + toast notifications
- Data fetching: `useLazyAsyncData` with reactive query parameters
- State management: Reactive objects for form state, computed properties for UI state

**✅ Existing Components (Extend If Needed):**

- Domain feature components in `frontend/app/components/features/`
- Service composables in `frontend/app/composables/services/`
- Repository modules in `frontend/shared/repository/`
- Store modules in `frontend/app/stores/`

### Step 2: Repository Layer Extension

**File: `frontend/shared/repository/{domain}.ts`**

```typescript
import { ConnectError } from "@connectrpc/connect";
import type {
  QueryEntityResponse,
  CreateEntityResponse,
  QueryEntityRequest,
  CreateEntityRequest,
  EntityService,
} from "~~/gen/domain/v1/entity_pb";
import type { Client } from "@connectrpc/connect";

export const entityRepository = (client: Client<typeof EntityService>) => ({
  async queryEntities(req: QueryEntityRequest): Promise<QueryEntityResponse> {
    try {
      const response = await client.queryEntities(req);
      return response;
    } catch (err) {
      if (err instanceof ConnectError) {
        console.error("ConnectError:", err);
      }
      throw err;
    }
  },

  async createEntity(req: CreateEntityRequest): Promise<CreateEntityResponse> {
    try {
      const response = await client.createEntity(req);
      return response;
    } catch (err) {
      if (err instanceof ConnectError) {
        console.error("ConnectError:", err);
      }
      throw err;
    }
  },
});
```

**Key Requirements:**

- Follow exact pattern from existing repositories
- Handle ConnectError properly with logging
- Use generated protobuf types from `~~/gen/`
- Return promises with proper typing

### Step 3: Service Composable Extension

**File: `frontend/app/composables/services/useEntityService.ts`**

```typescript
import {
  QueryEntityRequestSchema,
  CreateEntityRequestSchema,
  type Entity,
} from "~~/gen/domain/v1/entity_pb";
import { type MessageInitShape, create } from "@bufbuild/protobuf";
import type { QueryMetaResponseSchema } from "~~/gen/altalune/v1/common_pb";

import { entityRepository } from "#shared/repository/entity";
import { useConnectValidator } from "../useConnectValidator";
import { useErrorMessage } from "../useErrorMessage";

export function useEntityService() {
  const { $entityClient } = useNuxtApp();
  const entity = entityRepository($entityClient);
  const { parseError } = useErrorMessage();

  // Query functionality
  const queryValidator = useConnectValidator(QueryEntityRequestSchema);

  // Create functionality
  const createValidator = useConnectValidator(CreateEntityRequestSchema);
  const createState = reactive({
    loading: false,
    error: "",
    success: false,
  });

  async function query(req: MessageInitShape<typeof QueryEntityRequestSchema>) {
    queryValidator.reset();
    if (!queryValidator.validate(req)) {
      console.warn("Validation failed:", queryValidator.errors.value);
      return { data: [], meta: { rowCount: 0, pageCount: 0, filters: {} } };
    }

    try {
      const message = create(QueryEntityRequestSchema, req);
      const result = await entity.queryEntities(message);
      return { data: result.data, meta: result.meta };
    } catch (err) {
      throw new Error(parseError(err));
    }
  }

  async function createEntity(
    req: MessageInitShape<typeof CreateEntityRequestSchema>
  ): Promise<Entity | null> {
    createState.loading = true;
    createState.error = "";
    createState.success = false;

    createValidator.reset();

    if (!createValidator.validate(req)) {
      createState.loading = false;
      return null;
    }

    try {
      const message = create(CreateEntityRequestSchema, req);
      const result = await entity.createEntity(message);
      createState.success = true;
      return result.entity || null;
    } catch (err) {
      createState.error = parseError(err);
      throw new Error(createState.error);
    } finally {
      createState.loading = false;
    }
  }

  function resetCreateState() {
    createState.loading = false;
    createState.error = "";
    createState.success = false;
    createValidator.reset();
  }

  return {
    // Query
    query,
    queryValidationErrors: queryValidator.errors,

    // Create
    createEntity,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    createValidationErrors: createValidator.errors,
    resetCreateState,
  };
}
```

**Key Requirements:**

- Use reactive state management for form states
- Implement dual-layer validation: vee-validate (primary) + ConnectRPC (fallback)
- Handle errors with `useErrorMessage` and i18n
- Return computed properties for reactive UI binding
- Follow exact pattern from existing services

### Step 4: UI Components Creation

**File Structure: `frontend/app/components/features/{domain}/`**

- `EntityCreateSheet.vue` - Modal wrapper component
- `EntityCreateForm.vue` - Form component with validation
- `EntityTableLoading.vue` - Loading state component (optional)
- `index.ts` - Component exports

**Sheet Component Pattern:**

```vue
<script setup lang="ts">
import type { Entity } from "~~/gen/domain/v1/entity_pb";
import EntityCreateForm from "./EntityCreateForm.vue";
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from "@/components/ui/sheet";

const props = defineProps<{
  projectId?: string; // If entity belongs to a project
}>();

const emit = defineEmits<{
  success: [entity: Entity];
  cancel: [];
}>();

const isSheetOpen = ref(false);

function handleEntityCreated(entity: Entity) {
  isSheetOpen.value = false;
  emit("success", entity);
}

function handleSheetClose() {
  isSheetOpen.value = false;
  emit("cancel");
}
</script>

<template>
  <Sheet v-model:open="isSheetOpen">
    <SheetTrigger as-child>
      <slot />
    </SheetTrigger>
    <SheetContent class="w-full sm:max-w-[540px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>Add New Entity</SheetTitle>
        <SheetDescription>
          Fill in the details below. All fields marked with * are required.
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <EntityCreateForm
          :project-id="props.projectId"
          @success="handleEntityCreated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
```

**Form Component Pattern (vee-validate + ConnectRPC Dual Validation):**

**Important: Form Validation Strategy**

The project uses a dual-layer validation approach:

1. **Primary**: vee-validate with Zod schemas for immediate client-side feedback
2. **Fallback**: ConnectRPC validation for development/edge cases

This ensures robust validation while maintaining good UX with immediate feedback.

**Critical: Avoid Dual Field Binding**

Never bind multiple inputs to the same vee-validate field. Each field name should have only ONE input controlling it:

```vue
<!-- WRONG - Two inputs bound to same field -->
<FormField v-slot="{ componentField }" name="role">
  <Select v-bind="componentField">...</Select>      <!-- Field binding #1 -->
  <Input v-bind="componentField">...</Input>        <!-- Field binding #2 - CONFLICT! -->
</FormField>

<!-- CORRECT - Single input with enhanced UX -->
<FormField v-slot="{ componentField }" name="role">
  <Select v-bind="componentField">...</Select>      <!-- Only field binding -->
  <div class="text-xs text-muted-foreground">
    You can also type directly in the field above for custom values
  </div>
</FormField>
```

```vue
<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { useForm } from "vee-validate";
import * as z from "zod";
import { toast } from "vue-sonner";
import type { Entity } from "~~/gen/domain/v1/entity_pb";

import { useEntityService } from "@/composables/services/useEntityService";
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const props = defineProps<{
  projectId?: string;
}>();

const emit = defineEmits<{
  success: [entity: Entity];
  cancel: [];
}>();

const {
  createEntity,
  createLoading,
  createError,
  createValidationErrors,
  resetCreateState,
} = useEntityService();

// Create Zod schema matching protobuf validation rules
const formSchema = toTypedSchema(
  z.object({
    projectId: z.string().length(14), // Nanoid project ID
    name: z.string().min(2).max(50),
    email: z.string().email("Must be a valid email address"),
    // ... other fields with appropriate validation
  })
);

// Initialize vee-validate form
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    projectId: props.projectId || "",
    name: "",
    email: "",
    // ... other fields
  },
});

// ConnectRPC validation helpers (fallback layer)
const getConnectRPCError = (fieldName: string): string => {
  const errors =
    createValidationErrors.value[fieldName] ||
    createValidationErrors.value[`value.${fieldName}`];
  return errors?.[0] || "";
};

const hasConnectRPCError = (fieldName: string): boolean => {
  return !!(
    createValidationErrors.value[fieldName] ||
    createValidationErrors.value[`value.${fieldName}`]
  );
};

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const entity = await createEntity(values);

    if (entity) {
      toast.success("Entity created successfully", {
        description: `${values.name} has been created.`,
      });

      emit("success", entity);
      resetForm();
    }
  } catch (error) {
    console.error("Failed to create entity:", error);
    toast.error("Failed to create entity", {
      description:
        createError.value || "An unexpected error occurred. Please try again.",
    });
  }
});

function handleCancel() {
  resetForm();
  emit("cancel");
}

function resetForm() {
  form.resetForm({
    values: {
      projectId: props.projectId || "",
      name: "",
      email: "",
      // Reset all fields to initial values
    },
  });
  resetCreateState();
}

// Watch for prop changes
watch(
  () => props.projectId,
  (newProjectId) => {
    if (newProjectId) {
      form.setFieldValue("projectId", newProjectId);
    }
  }
);

onUnmounted(() => {
  resetCreateState();
});
</script>

<template>
  <form class="space-y-6" @submit="onSubmit">
    <!-- Error Alert -->
    <Alert v-if="createError" variant="destructive">
      <Icon name="lucide:alert-circle" size="1em" mode="svg" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ createError }}</AlertDescription>
    </Alert>

    <!-- Form Fields with vee-validate + ConnectRPC fallback -->
    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>Name *</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            placeholder="Entity name"
            :class="{ 'border-destructive': hasConnectRPCError('name') }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription> Entity name (required) </FormDescription>
        <!-- Primary validation message from vee-validate -->
        <FormMessage />
        <!-- Fallback validation message from ConnectRPC -->
        <div v-if="hasConnectRPCError('name')" class="text-sm text-destructive">
          {{ getConnectRPCError("name") }}
        </div>
      </FormItem>
    </FormField>

    <!-- Action Buttons -->
    <div class="flex justify-end space-x-2 pt-4">
      <Button
        type="button"
        variant="outline"
        :disabled="createLoading"
        @click="handleCancel"
      >
        Cancel
      </Button>
      <Button type="submit" :disabled="createLoading">
        <Icon
          v-if="createLoading"
          name="lucide:loader-2"
          class="mr-2 animate-spin"
          size="1em"
          mode="svg"
        />
        {{ createLoading ? "Creating..." : "Create Entity" }}
      </Button>
    </div>
  </form>
</template>
```

### Step 5: Data Table Integration

**Row Actions Architecture:**

The project provides two approaches for data table row actions:

1. **DataTableBasicRowActions** - Generic component for basic CRUD operations:

   ```vue
   import { DataTableBasicRowActions } from '@/components/custom/datatable'; //
   Usage in column definition columnHelper.display({ id: 'actions', cell: ({ row
   }) => h(DataTableBasicRowActions, { row, actions: { edit: true, delete: true,
   duplicate: false, favorite: false }, onEdit: handleEdit, onDelete:
   handleDelete, onDuplicate: handleDuplicate, onFavorite: handleFavorite, }),
   })
   ```

2. **Domain-Specific Row Actions** - For complex domain logic:

   ```vue
   // Create in: components/features/{domain}/{Domain}RowActions.vue // Example:
   components/features/employee/EmployeeRowActions.vue

   <script setup lang="ts">
   // Domain-specific logic, validation, permissions, etc.
   </script>
   ```

**When to use which approach:**

- **Use DataTableBasicRowActions** for simple CRUD operations with standard edit/delete patterns
- **Create Domain Row Actions** when you need:
  - Complex business logic or validation
  - Domain-specific permissions/role checks
  - Custom actions beyond basic CRUD
  - Multi-step workflows or confirmations

**Query Implementation with Data Tables:**

```vue
<script setup lang="ts">
import { useEntityService } from "@/composables/services/useEntityService";
import { useQueryRequest } from "@/composables/useQueryRequest";
import {
  useDataTableState,
  useDataTableFilter,
} from "@/components/custom/datatable/utils";
import { DataTable } from "@/components/custom/datatable";

const { query, resetCreateState } = useEntityService();

// Data table state
const page = ref(1);
const pageSize = ref(10);
const keyword = ref("");
const dataTableRef = ref<InstanceType<typeof DataTable> | null>(null);
const table = computed(() => dataTableRef.value?.table);
const { columnFilters, sorting } = useDataTableState(dataTableRef);

// Query request management
const { queryRequest } = useQueryRequest({
  page,
  pageSize,
  keyword,
  columnFilters,
  sorting,
});

// Async data fetching
const asyncDataKey = computed(() => {
  const { pagination, keyword, filters, sorting } = queryRequest.value;
  const { page, pageSize } = pagination!;
  return [
    "entity-table",
    page,
    pageSize,
    keyword,
    filters ? JSON.stringify(filters) : null,
    sorting ? `${sorting.field}:${sorting.order}` : null,
  ]
    .filter(Boolean)
    .join("-");
});

const {
  data: response,
  pending,
  refresh,
} = useLazyAsyncData(asyncDataKey, () => query(queryRequest.value), {
  server: false,
  watch: [queryRequest],
  immediate: false,
});

const data = computed(() => response.value?.data ?? []);
const rowCount = computed(() => response.value?.meta?.rowCount ?? 0);
const filters = computed(() => response.value?.meta?.filters);

function handleEntityCreated() {
  resetCreateState();
  refresh();
}
</script>
```

### Step 6: Store Integration (If Needed)

**For domain-specific global state:**

```typescript
// frontend/app/stores/entity.ts
export const useEntityStore = defineStore("entity", () => {
  const entities = ref<Entity[]>([]);
  const pending = ref(false);
  const error = ref<Error | null>(null);

  function setEntities(newEntities: Entity[]) {
    entities.value = newEntities;
  }

  function addEntity(entity: Entity) {
    entities.value.unshift(entity);
  }

  return {
    entities: readonly(entities),
    pending: readonly(pending),
    error: readonly(error),
    setEntities,
    addEntity,
  };
});
```

### Step 7: Code Generation & Testing

```bash
# Generate protobuf types
buf generate

# Start frontend development server
cd frontend && pnpm dev

# Test the implementation
# 1. Navigate to the page with the create functionality
# 2. Test form validation (required fields, patterns)
# 3. Test successful creation flow
# 4. Test error handling scenarios
# 5. Verify data table refresh after creation
```

### Frontend Development Checklist

**Before Implementation:**

- [ ] Analyzed existing frontend components for reusability
- [ ] Add new shadcn-vue component if not exist yet using `pnpm dlx shadcn-vue@latest add <component>`
- [ ] Checked if repository/service already exists for the domain
- [ ] Reviewed similar domain implementations for patterns
- [ ] Generated latest protobuf types (`buf generate`)

**During Implementation:**

- [ ] Used established repository pattern with Connect error handling
- [ ] Followed service composable pattern with reactive state
- [ ] Implemented dual-layer validation: vee-validate (primary) + ConnectRPC (fallback)
- [ ] Used shadcn-vue components for consistent UI
- [ ] Added comprehensive error handling with toast notifications
- [ ] Followed established naming conventions and file structure
- [ ] Chose appropriate row actions approach (DataTableBasicRowActions vs domain-specific)

**After Implementation:**

- [ ] Tested form validation (required fields, patterns, edge cases)
- [ ] Tested successful creation flow with proper feedback
- [ ] Tested error scenarios and error message display
- [ ] Verified data table refresh after creation
- [ ] Tested responsive design on mobile/desktop
- [ ] Verified accessibility (keyboard navigation, screen readers)

## Sheet/Dialog Best Practices

**Common Issue: Sheets/Dialogs Inside Dropdown Menus**

When placing Sheet or AlertDialog components inside DropdownMenu items, you may encounter an issue where the sheet/dialog opens but immediately closes. This happens due to event propagation conflicts.

**Problematic Pattern (Causes Immediate Closing):**

```vue
<DropdownMenu>
  <DropdownMenuContent>
    <!-- DON'T DO THIS - Sheet will close immediately -->
    <MyCustomSheet>
      <DropdownMenuItem>Edit</DropdownMenuItem>
    </MyCustomSheet>
  </DropdownMenuContent>
</DropdownMenu>
```

**Correct Pattern (Manual Control):**

```vue
<template>
  <DropdownMenu>
    <DropdownMenuContent>
      <!-- Use direct click handler -->
      <DropdownMenuItem @click="openMySheet">Edit</DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>

  <!-- Place Sheet/Dialog outside dropdown -->
  <MyCustomSheet
    v-model:open="isSheetOpen"
    @success="handleSuccess"
    @cancel="closeMySheet"
  />
</template>

<script setup>
const isSheetOpen = ref(false);

function openMySheet() {
  // Use nextTick to ensure dropdown closes first
  nextTick(() => {
    isSheetOpen.value = true;
  });
}

function closeMySheet() {
  isSheetOpen.value = false;
}

function handleSuccess() {
  closeMySheet();
  // Handle success logic
}
</script>
```

**Key Requirements for Sheet/Dialog Components:**

1. **Support v-model:open**: Components should accept an optional `open` prop and emit `update:open`
2. **Conditional Trigger**: Only show trigger slot when not controlled externally
3. **nextTick Usage**: Use `nextTick()` when opening from dropdown items

**Component Implementation Pattern:**

```vue
<!-- MyCustomSheet.vue -->
<script setup>
const props = defineProps<{
  // Your component props
  data: SomeType;
  open?: boolean; // For v-model:open support
}>();

const emit = defineEmits<{
  success: [result: SomeType];
  cancel: [];
  'update:open': [value: boolean]; // For v-model:open support
}>();

// Support both internal state and v-model:open
const isSheetOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});
</script>

<template>
  <Sheet v-model:open="isSheetOpen">
    <!-- Only show trigger when not controlled externally -->
    <SheetTrigger v-if="!props.open && $slots.default" as-child>
      <slot />
    </SheetTrigger>
    <SheetContent>
      <!-- Your content -->
    </SheetContent>
  </Sheet>
</template>
```

**Why This Pattern Works:**

1. **Event Isolation**: Sheet/Dialog is outside the dropdown, preventing conflicts
2. **Proper Timing**: `nextTick()` ensures dropdown closes before sheet opens
3. **Flexible Usage**: Component works both with triggers and manual control
4. **Clean Separation**: UI state management is explicit and predictable

## Icon Usage Guidelines

### Overview

This project uses **Nuxt Icon component** as the standard way to display icons throughout the application. All icon usage should follow this pattern to ensure consistency and maintainability.

### Correct Icon Usage Pattern

**Always use the Nuxt Icon component:**

- `name` (required): icon name or global component name
- `size`: icon size (default: `1em`)
- `mode`: icon rendering mode (`svg` or `css`, default: `css`)

```vue
<!-- CORRECT: Use Nuxt Icon component -->
<Icon name="lucide:user-plus" class="mr-2" size="1em" mode="svg" />
<Icon name="lucide:edit" class="mr-2" size="1em" mode="svg" />
<Icon name="lucide:trash-2" class="mr-2" size="1em" mode="svg" />
<Icon name="lucide:loader-2" class="mr-2 animate-spin" size="1em" mode="svg" />
```

### Incorrect Patterns to Avoid

**Never import icons directly from lucide packages:**

```vue
<!-- WRONG: Direct lucide imports -->
<script setup lang="ts">
import { UserPlus, Edit, Trash2 } from "lucide-vue-next";
</script>

<template>
  <UserPlus class="mr-2 h-4 w-4" />
  <Edit class="mr-2 h-4 w-4" />
  <Trash2 class="mr-2 h-4 w-4" />
</template>
```

### Icon Naming Convention

When using Nuxt Icon component, follow the `lucide:` prefix with kebab-case naming:

```vue
<!-- Lucide icon names in kebab-case -->
<Icon name="lucide:user-plus" />
<!-- UserPlus becomes user-plus -->
<Icon name="lucide:more-horizontal" />
<!-- MoreHorizontal becomes more-horizontal -->
<Icon name="lucide:chevron-down" />
<!-- ChevronDown becomes chevron-down -->
<Icon name="lucide:alert-circle" />
<!-- AlertCircle becomes alert-circle -->
<Icon name="lucide:shield-alert" />
<!-- ShieldAlert becomes shield-alert -->
<Icon name="lucide:eye-off" />
<!-- EyeOff becomes eye-off -->
```

### Shadcn Components Icon Migration

**Important:** When using shadcn-vue components that have direct lucide imports, you must modify them to use the Nuxt Icon component instead.

**Before (shadcn default):**

```vue
<script setup lang="ts">
import { ChevronDown } from "lucide-vue-next";
</script>

<template>
  <SelectTrigger>
    <slot />
    <SelectIcon as-child>
      <ChevronDown class="size-4 opacity-50" />
    </SelectIcon>
  </SelectTrigger>
</template>
```

**After (Nuxt Icon component):**

```vue
<script setup lang="ts">
// Remove direct lucide import
</script>

<template>
  <SelectTrigger>
    <slot />
    <SelectIcon as-child>
      <Icon name="lucide:chevron-down" class="size-4 opacity-50" />
    </SelectIcon>
  </SelectTrigger>
</template>
```

### Common Icon Migration Examples

| Direct Import                                   | Nuxt Icon Component                                                    |
| ----------------------------------------------- | ---------------------------------------------------------------------- |
| `<AlertCircle class="w-4 h-4" />`               | `<Icon name="lucide:alert-circle" />`                                  |
| `<UserPlus class="mr-2 h-4 w-4" />`             | `<Icon name="lucide:user-plus" />`                                     |
| `<MoreHorizontal class="h-4 w-4" />`            | `<Icon name="lucide:more-horizontal" />`                               |
| `<X class="size-4" />`                          | `<Icon name="lucide:x" size="2em" />`                                  |
| `<Search class="size-4 shrink-0 opacity-50" />` | `<Icon name="lucide:search" size="1em" class="shrink-0 opacity-50" />` |

### Icon Usage in Components

**Data Table Row Actions:**

```vue
<DropdownMenuItem @click="handleEdit">
  <Icon name="lucide:edit" size="1em" class="mr-2" />
  Edit
</DropdownMenuItem>

<DropdownMenuItem @click="handleDelete">
  <Icon name="lucide:trash-2" size="1em" class="mr-2" />
  Delete
</DropdownMenuItem>
```

**Loading States:**

```vue
<Button :disabled="loading">
  <Icon v-if="loading" name="lucide:loader-2" size="1em" class="mr-2 animate-spin" />
  {{ loading ? 'Saving...' : 'Save' }}
</Button>
```

**Form Validation:**

```vue
<Alert variant="destructive">
  <Icon name="lucide:alert-circle" size="1em" />
  <AlertTitle>Error</AlertTitle>
  <AlertDescription>{{ error }}</AlertDescription>
</Alert>
```

### Benefits of This Approach

1. **Consistency**: Single pattern across the entire application
2. **Bundle Optimization**: Nuxt Icon handles icon loading and optimization
3. **Maintainability**: Easy to update or replace icon libraries
4. **Developer Experience**: No need to manage individual icon imports
5. **Type Safety**: TypeScript support for icon names

### Migration Checklist

When migrating components from direct lucide imports:

- [ ] Remove `import { IconName } from 'lucide-vue-next'` statements
- [ ] Replace `<IconName />` with `<Icon name="lucide:icon-name" />`
- [ ] Convert PascalCase icon names to kebab-case
- [ ] Preserve all existing classes and attributes
- [ ] Test the component to ensure icons display correctly

### Finding Direct Icon Imports

To find components that still use direct lucide imports:

```bash
# Search for direct lucide imports
grep -r "from ['\"]lucide" frontend/app/components/

# Search for specific icon usage patterns
grep -r "import.*lucide-vue" frontend/app/
```

## Internationalization (i18n)

### Translation File Structure

**Location:** `frontend/i18n/locales/{locale}.json`

**Supported Locales:**

- `en-US.json` - English (US)
- `id-ID.json` - Bahasa Indonesia

### Key Principles

1. **Fully Nested Structure** - All translations must use proper object nesting
2. **Alphabetical Ordering** - Keys should be organized alphabetically for maintainability
3. **No Mixed Nesting** - Avoid mixing flat keys (with dots) inside nested objects

### Format Rules

✅ **CORRECT - Fully nested structure:**

```json
{
  "common": {
    "btn": {
      "cancel": "Cancel",
      "create": "Create",
      "delete": "Delete"
    },
    "status": {
      "loading": "Loading...",
      "success": "Success"
    }
  },
  "features": {
    "api_keys": {
      "actions": {
        "create": "Create API Key",
        "edit": "Edit API Key"
      },
      "columns": {
        "name": "Name",
        "status": "Status"
      }
    }
  },
  "nav": {
    "dashboard": "Dashboard",
    "devices": {
      "chat": "Chat",
      "scan": "Scan"
    }
  }
}
```

❌ **INCORRECT - Mixed nesting (flat keys inside nested objects):**

```json
{
  "features": {
    "api_keys": {
      "actions.create": "Create API Key", // DON'T DO THIS
      "columns.name": "Name"
    }
  }
}
```

### Why Fully Nested?

1. **Framework Support** - `@nuxtjs/i18n` requires consistent nesting
2. **IDE Autocomplete** - Better TypeScript support for nested structures
3. **Clear Organization** - Visual hierarchy matches logical structure
4. **Standard Practice** - Follows common i18n patterns

### Usage in Components

**Basic Text Translation:**

```vue
<script setup lang="ts">
const { t } = useI18n();
</script>

<template>
  <h1>{{ t("features.api_keys.sheet.createTitle") }}</h1>
  <Button>{{ t("common.btn.create") }}</Button>
</template>
```

**Translation with Variables:**

```json
{
  "features": {
    "employees": {
      "messages": {
        "createSuccessDesc": "{name} has been added to the team."
      }
    }
  }
}
```

```typescript
toast.success(
  t("features.employees.messages.createSuccessDesc", { name: "John Doe" })
);
// Output: "John Doe has been added to the team."
```

**Formatted Text with Markdown:**

For translations that need basic styling (bold, italic), use the `useI18nSafe` composable:

```json
{
  "features": {
    "api_keys": {
      "deleteDialog": {
        "confirmMessage": "Are you sure you want to delete **{name}**? This action cannot be undone."
      }
    }
  }
}
```

```vue
<script setup lang="ts">
const { t, tFormatted } = useI18nSafe();
</script>

<template>
  <AlertDialogDescription>
    <component
      :is="
        tFormatted('features.api_keys.deleteDialog.confirmMessage', {
          name: apiKey.name,
        })
      "
    />
  </AlertDialogDescription>
</template>
```

**Supported Markdown Syntax:**

- `**text**` → `<strong>text</strong>` (bold)
- `*text*` → `<em>text</em>` (italic)
- `__text__` → `<u>text</u>` (underline)

### Special Character Escaping

For special characters that conflict with i18n syntax:

```json
{
  "features": {
    "employees": {
      "form": {
        "emailPlaceholder": "john.doe{'@'}company.com"
      }
    }
  }
}
```

The `{'@'}` syntax escapes the `@` symbol, which would otherwise be interpreted as a linked message reference.

### Reactive Translations in Computed Properties

When using translations in programmatic rendering (e.g., TanStack Table columns), wrap in `computed()` for reactivity:

```typescript
// ❌ WRONG - Not reactive
const columns = [
  columnHelper.accessor("name", {
    header: ({ column }) =>
      h(DataTableColumnHeader, {
        column,
        title: t("features.api_keys.columns.name"), // Called once only
      }),
  }),
];

// ✅ CORRECT - Reactive
const columns = computed(() => [
  columnHelper.accessor("name", {
    header: ({ column }) =>
      h(DataTableColumnHeader, {
        column,
        title: t("features.api_keys.columns.name"), // Re-evaluated on locale change
      }),
  }),
]);
```

### Configuration

**File:** `frontend/nuxt.config.ts`

```typescript
export default defineNuxtConfig({
  i18n: {
    strategy: "no_prefix",
    defaultLocale: "en-US",
    lazy: true, // Enable lazy loading
    langDir: "locales", // Translation files directory
    locales: [
      {
        code: "en-US",
        name: "English",
        file: "en-US.json",
        dir: "ltr",
      },
      {
        code: "id-ID",
        name: "Bahasa Indonesia",
        file: "id-ID.json",
        dir: "ltr",
      },
    ],
  },
});
```

### Adding New Translations

**Step 1: Add to all locale files**

Update **both** `en-US.json` and `id-ID.json`:

```json
// en-US.json
{
  "features": {
    "new_feature": {
      "title": "New Feature",
      "description": "Feature description"
    }
  }
}

// id-ID.json
{
  "features": {
    "new_feature": {
      "title": "Fitur Baru",
      "description": "Deskripsi fitur"
    }
  }
}
```

**Step 2: Use in components**

```vue
<script setup lang="ts">
const { t } = useI18n();
</script>

<template>
  <h1>{{ t("features.new_feature.title") }}</h1>
  <p>{{ t("features.new_feature.description") }}</p>
</template>
```

### Important Notes

**Database-Sourced Values:**

Do NOT add translations for values that come from the database (roles, departments, categories, etc.). These should remain dynamic:

```typescript
// ❌ WRONG - Don't translate database values
const roles = [{ value: "engineer", label: t("roles.engineer") }];

// ✅ CORRECT - Use database values directly
const roles = [{ value: "engineer", label: "Engineer" }];
```

**Error Codes:**

Backend error codes should be translated for user-facing messages:

```json
{
  "errorCodes": {
    "60001": "Invalid input",
    "60201": "Employee not found",
    "69001": "Server Error"
  }
}
```

Usage:

```typescript
const errorMessage = t(`errorCodes.${errorCode}`);
```

### Translation Checklist

When adding new translations:

- [ ] Translation exists in **all** locale files
- [ ] Keys follow fully nested structure (no mixed nesting)
- [ ] Keys are organized alphabetically where appropriate
- [ ] Variables use `{variable}` syntax
- [ ] Special characters are properly escaped
- [ ] Markdown formatting uses `useI18nSafe` composable
- [ ] Reactive contexts use `computed()` wrapper
- [ ] Database-sourced values are NOT translated

### Common Patterns

**Feature Actions:**

```json
{
  "features": {
    "{feature_name}": {
      "actions": {
        "create": "Create {Entity}",
        "edit": "Edit {Entity}",
        "delete": "Delete {Entity}"
      }
    }
  }
}
```

**Form Labels:**

```json
{
  "features": {
    "{feature_name}": {
      "form": {
        "nameLabel": "Name *",
        "namePlaceholder": "Enter name",
        "nameDescription": "Name description"
      }
    }
  }
}
```

**Messages:**

```json
{
  "features": {
    "{feature_name}": {
      "messages": {
        "createSuccess": "Entity created successfully",
        "createSuccessDesc": "{name} has been created.",
        "createError": "Failed to create entity",
        "createErrorDesc": "An unexpected error occurred."
      }
    }
  }
}
```

This workflow ensures type-safe, consistent, and maintainable frontend development while leveraging the power of protobuf validation and Connect-RPC.
