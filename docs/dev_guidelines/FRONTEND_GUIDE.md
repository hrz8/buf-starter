## Frontend Development Workflow

This systematic workflow ensures consistency, type safety, and proper architecture adherence when implementing new frontend features.

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
- Implement proper validation with `useConnectValidator`
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

**Form Component Pattern:**

```vue
<script setup lang="ts">
import { toast } from "vue-sonner";
import type {
  CreateEntityRequestSchema,
  Entity,
} from "~~/gen/domain/v1/entity_pb";
import type { MessageInitShape } from "@bufbuild/protobuf";

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

const formData = reactive<MessageInitShape<typeof CreateEntityRequestSchema>>({
  // Initialize form fields based on schema
  name: "",
  description: "",
  // ... other fields
});

const getFieldError = (fieldName: string): string => {
  const errors =
    createValidationErrors.value[fieldName] ||
    createValidationErrors.value[`value.${fieldName}`];
  return errors?.[0] || "";
};

const hasFieldError = (fieldName: string): boolean => {
  return !!(
    createValidationErrors.value[fieldName] ||
    createValidationErrors.value[`value.${fieldName}`]
  );
};

async function handleSubmit() {
  try {
    const entity = await createEntity(formData);

    if (entity) {
      toast.success("Entity created successfully", {
        description: `${formData.name} has been created.`,
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
}

function handleCancel() {
  resetForm();
  emit("cancel");
}

function resetForm() {
  Object.assign(formData, {
    name: "",
    description: "",
    // Reset all fields
  });
  resetCreateState();
}

onUnmounted(() => {
  resetCreateState();
});
</script>

<template>
  <form class="space-y-6" @submit.prevent="handleSubmit">
    <!-- Error Alert -->
    <Alert v-if="createError" variant="destructive">
      <AlertCircle class="w-4 h-4" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ createError }}</AlertDescription>
    </Alert>

    <!-- Form Fields -->
    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>Name *</FormLabel>
        <FormControl>
          <Input
            v-model="formData.name"
            v-bind="componentField"
            placeholder="Entity name"
            :class="{ 'border-destructive': hasFieldError('name') }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription> Entity name (required) </FormDescription>
        <FormMessage v-if="hasFieldError('name')" class="text-destructive">
          {{ getFieldError("name") }}
        </FormMessage>
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
          class="mr-2 h-4 w-4 animate-spin"
        />
        {{ createLoading ? "Creating..." : "Create Entity" }}
      </Button>
    </div>
  </form>
</template>
```

### Step 5: Data Table Integration

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
- [ ] Checked if repository/service already exists for the domain
- [ ] Reviewed similar domain implementations for patterns
- [ ] Generated latest protobuf types (`buf generate`)

**During Implementation:**

- [ ] Used established repository pattern with Connect error handling
- [ ] Followed service composable pattern with reactive state
- [ ] Implemented proper form validation with `useConnectValidator`
- [ ] Used shadcn-vue components for consistent UI
- [ ] Added comprehensive error handling with toast notifications
- [ ] Followed established naming conventions and file structure

**After Implementation:**

- [ ] Tested form validation (required fields, patterns, edge cases)
- [ ] Tested successful creation flow with proper feedback
- [ ] Tested error scenarios and error message display
- [ ] Verified data table refresh after creation
- [ ] Tested responsive design on mobile/desktop
- [ ] Verified accessibility (keyboard navigation, screen readers)

This workflow ensures type-safe, consistent, and maintainable frontend development while leveraging the power of protobuf validation and Connect-RPC.
