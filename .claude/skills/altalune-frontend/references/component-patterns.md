# Component Implementation Patterns

## Table of Contents
- [Repository Layer](#repository-layer)
- [Service Composable](#service-composable)
- [Feature Components](#feature-components)
- [Data Table](#data-table)
- [Centralized Feature Files](#centralized-feature-files)

## Repository Layer

**File: `frontend/shared/repository/{domain}.ts`**

```typescript
import { ConnectError } from '@connectrpc/connect'
import type { Client } from '@connectrpc/connect'
import type {
  QueryEntityRequest,
  QueryEntityResponse,
  CreateEntityRequest,
  CreateEntityResponse,
  EntityService,
} from '~~/gen/altalune/v1/entity_pb'

export const entityRepository = (client: Client<typeof EntityService>) => ({
  async queryEntities(req: QueryEntityRequest): Promise<QueryEntityResponse> {
    try {
      const response = await client.queryEntities(req)
      return response
    }
    catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err)
      }
      throw err
    }
  },

  async createEntity(req: CreateEntityRequest): Promise<CreateEntityResponse> {
    try {
      const response = await client.createEntity(req)
      return response
    }
    catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err)
      }
      throw err
    }
  },

  async getEntity(req: GetEntityRequest): Promise<GetEntityResponse> {
    try {
      return await client.getEntity(req)
    }
    catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err)
      }
      throw err
    }
  },

  async updateEntity(req: UpdateEntityRequest): Promise<UpdateEntityResponse> {
    try {
      return await client.updateEntity(req)
    }
    catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err)
      }
      throw err
    }
  },

  async deleteEntity(req: DeleteEntityRequest): Promise<DeleteEntityResponse> {
    try {
      return await client.deleteEntity(req)
    }
    catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err)
      }
      throw err
    }
  },
})
```

## Service Composable

**File: `frontend/app/composables/services/use{Domain}Service.ts`**

```typescript
import {
  QueryEntityRequestSchema,
  CreateEntityRequestSchema,
  type Entity,
} from '~~/gen/altalune/v1/entity_pb'
import { type MessageInitShape, create } from '@bufbuild/protobuf'

import { entityRepository } from '#shared/repository/entity'
import { useConnectValidator } from '../useConnectValidator'
import { useErrorMessage } from '../useErrorMessage'

export function useEntityService() {
  const { $entityClient } = useNuxtApp()
  const entity = entityRepository($entityClient)
  const { parseError } = useErrorMessage()

  // Query functionality
  const queryValidator = useConnectValidator(QueryEntityRequestSchema)

  // Create functionality
  const createValidator = useConnectValidator(CreateEntityRequestSchema)
  const createState = reactive({
    loading: false,
    error: '',
    success: false,
  })

  async function query(req: MessageInitShape<typeof QueryEntityRequestSchema>) {
    queryValidator.reset()
    if (!queryValidator.validate(req)) {
      console.warn('Validation failed:', queryValidator.errors.value)
      return { data: [], meta: { rowCount: 0, pageCount: 0, filters: {} } }
    }

    try {
      const message = create(QueryEntityRequestSchema, req)
      const result = await entity.queryEntities(message)
      return { data: result.data, meta: result.meta }
    }
    catch (err) {
      throw new Error(parseError(err))
    }
  }

  async function createEntity(
    req: MessageInitShape<typeof CreateEntityRequestSchema>,
  ): Promise<Entity | null> {
    createState.loading = true
    createState.error = ''
    createState.success = false

    createValidator.reset()

    if (!createValidator.validate(req)) {
      createState.loading = false
      return null
    }

    try {
      const message = create(CreateEntityRequestSchema, req)
      const result = await entity.createEntity(message)
      createState.success = true
      return result.entity || null
    }
    catch (err) {
      createState.error = parseError(err)
      throw new Error(createState.error)
    }
    finally {
      createState.loading = false
    }
  }

  function resetCreateState() {
    createState.loading = false
    createState.error = ''
    createState.success = false
    createValidator.reset()
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
  }
}
```

## Feature Components

### Sheet Component (Modal Wrapper)

```vue
<script setup lang="ts">
import type { Entity } from '~~/gen/altalune/v1/entity_pb'
import EntityCreateForm from './EntityCreateForm.vue'
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from '@/components/ui/sheet'

const props = defineProps<{
  projectId?: string
  open?: boolean  // For v-model:open support
}>()

const emit = defineEmits<{
  success: [entity: Entity]
  cancel: []
  'update:open': [value: boolean]
}>()

// Support both internal state and v-model:open
const isSheetOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
})

function handleEntityCreated(entity: Entity) {
  isSheetOpen.value = false
  emit('success', entity)
}

function handleSheetClose() {
  isSheetOpen.value = false
  emit('cancel')
}
</script>

<template>
  <Sheet v-model:open="isSheetOpen">
    <!-- Only show trigger when not controlled externally -->
    <SheetTrigger v-if="!props.open && $slots.default" as-child>
      <slot />
    </SheetTrigger>
    <SheetContent class="w-full sm:max-w-[540px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>{{ t('features.entity.sheet.createTitle') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.entity.sheet.createDescription') }}
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

### Form Component

```vue
<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { toast } from 'vue-sonner'
import type { Entity } from '~~/gen/altalune/v1/entity_pb'

import { useEntityService } from '@/composables/services/useEntityService'
import { entityCreateSchema } from './schema'
import { getConnectRPCError, hasConnectRPCError } from './error'
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  FormDescription,
} from '@/components/ui/form'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'

const props = defineProps<{
  projectId?: string
}>()

const emit = defineEmits<{
  success: [entity: Entity]
  cancel: []
}>()

const { t } = useI18n()
const {
  createEntity,
  createLoading,
  createError,
  createValidationErrors,
  resetCreateState,
} = useEntityService()

// Zod schema with vee-validate
const formSchema = toTypedSchema(entityCreateSchema)

const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    projectId: props.projectId || '',
    name: '',
    email: '',
  },
})

// Handle form submission
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const entity = await createEntity(values)

    if (entity) {
      toast.success(t('features.entity.messages.createSuccess'), {
        description: t('features.entity.messages.createSuccessDesc', { name: values.name }),
      })

      emit('success', entity)
      resetForm()
    }
  }
  catch (error) {
    console.error('Failed to create entity:', error)
    toast.error(t('features.entity.messages.createError'), {
      description: createError.value || t('common.errors.unexpected'),
    })
  }
})

function handleCancel() {
  resetForm()
  emit('cancel')
}

function resetForm() {
  form.resetForm({
    values: {
      projectId: props.projectId || '',
      name: '',
      email: '',
    },
  })
  resetCreateState()
}

// Watch for prop changes
watch(
  () => props.projectId,
  (newProjectId) => {
    if (newProjectId) {
      form.setFieldValue('projectId', newProjectId)
    }
  },
)

onUnmounted(() => {
  resetCreateState()
})
</script>

<template>
  <form class="space-y-6" @submit="onSubmit">
    <!-- Error Alert -->
    <Alert v-if="createError" variant="destructive">
      <Icon name="lucide:alert-circle" size="1em" mode="svg" />
      <AlertTitle>{{ t('common.error') }}</AlertTitle>
      <AlertDescription>{{ createError }}</AlertDescription>
    </Alert>

    <!-- Form Fields -->
    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>{{ t('features.entity.form.nameLabel') }} *</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            :placeholder="t('features.entity.form.namePlaceholder')"
            :class="{ 'border-destructive': hasConnectRPCError(createValidationErrors, 'name') }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>{{ t('features.entity.form.nameDescription') }}</FormDescription>
        <FormMessage />
        <!-- ConnectRPC fallback error -->
        <div v-if="hasConnectRPCError(createValidationErrors, 'name')" class="text-sm text-destructive">
          {{ getConnectRPCError(createValidationErrors, 'name') }}
        </div>
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="email">
      <FormItem>
        <FormLabel>{{ t('features.entity.form.emailLabel') }} *</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="email"
            :placeholder="t('features.entity.form.emailPlaceholder')"
            :disabled="createLoading"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- Action Buttons -->
    <div class="flex justify-end space-x-2 pt-4">
      <Button type="button" variant="outline" :disabled="createLoading" @click="handleCancel">
        {{ t('common.btn.cancel') }}
      </Button>
      <Button type="submit" :disabled="createLoading">
        <Icon
          v-if="createLoading"
          name="lucide:loader-2"
          class="mr-2 animate-spin"
          size="1em"
          mode="svg"
        />
        {{ createLoading ? t('common.btn.creating') : t('common.btn.create') }}
      </Button>
    </div>
  </form>
</template>
```

## Data Table

### Table Page Component

```vue
<script setup lang="ts">
import { useEntityService } from '@/composables/services/useEntityService'
import { useQueryRequest } from '@/composables/useQueryRequest'
import {
  useDataTableState,
  useDataTableFilter,
} from '@/components/custom/datatable/utils'
import { DataTable } from '@/components/custom/datatable'
import EntityCreateSheet from '@/components/features/entity/EntityCreateSheet.vue'
import { columns } from './columns'

const { t } = useI18n()
const { query, resetCreateState } = useEntityService()

// Data table state
const page = ref(1)
const pageSize = ref(10)
const keyword = ref('')
const dataTableRef = ref<InstanceType<typeof DataTable> | null>(null)
const table = computed(() => dataTableRef.value?.table)
const { columnFilters, sorting } = useDataTableState(dataTableRef)

// Query request management
const { queryRequest } = useQueryRequest({
  page,
  pageSize,
  keyword,
  columnFilters,
  sorting,
})

// Async data fetching
const asyncDataKey = computed(() => {
  const { pagination, keyword, filters, sorting } = queryRequest.value
  const { page, pageSize } = pagination!
  return [
    'entity-table',
    page,
    pageSize,
    keyword,
    filters ? JSON.stringify(filters) : null,
    sorting ? `${sorting.field}:${sorting.order}` : null,
  ]
    .filter(Boolean)
    .join('-')
})

const {
  data: response,
  pending,
  refresh,
} = useLazyAsyncData(asyncDataKey, () => query(queryRequest.value), {
  server: false,
  watch: [queryRequest],
  immediate: false,
})

const data = computed(() => response.value?.data ?? [])
const rowCount = computed(() => response.value?.meta?.rowCount ?? 0)
const filters = computed(() => response.value?.meta?.filters)

function handleEntityCreated() {
  resetCreateState()
  refresh()
}
</script>

<template>
  <div class="space-y-5 px-4 py-3 sm:px-6 lg:px-8">
    <!-- Page Title & Description -->
    <div class="container mx-auto">
      <h2 class="text-2xl font-bold">{{ t('features.entity.page.title') }}</h2>
      <p class="text-muted-foreground">{{ t('features.entity.page.description') }}</p>
    </div>

    <!-- Action Buttons -->
    <div class="container mx-auto flex justify-end">
      <EntityCreateSheet @success="handleEntityCreated">
        <Button>
          <Icon name="lucide:plus" class="mr-2" size="1em" mode="svg" />
          {{ t('features.entity.actions.create') }}
        </Button>
      </EntityCreateSheet>
    </div>

    <!-- Data Table -->
    <div class="container mx-auto">
      <DataTable
        ref="dataTableRef"
        v-model:page="page"
        v-model:page-size="pageSize"
        v-model:keyword="keyword"
        :columns="columns"
        :data="data"
        :row-count="rowCount"
        :filters="filters"
        :loading="pending"
      />
    </div>
  </div>
</template>
```

### Column Definitions

```typescript
// columns.ts
import { h } from 'vue'
import { createColumnHelper } from '@tanstack/vue-table'
import type { Entity } from '~~/gen/altalune/v1/entity_pb'
import { DataTableColumnHeader } from '@/components/custom/datatable'
import EntityRowActions from './EntityRowActions.vue'

const columnHelper = createColumnHelper<Entity>()

export const columns = computed(() => [
  columnHelper.accessor('name', {
    header: ({ column }) =>
      h(DataTableColumnHeader, {
        column,
        title: t('features.entity.columns.name'),
      }),
    cell: ({ row }) => h('span', { class: 'font-medium' }, row.original.name),
  }),

  columnHelper.accessor('email', {
    header: ({ column }) =>
      h(DataTableColumnHeader, {
        column,
        title: t('features.entity.columns.email'),
      }),
    cell: ({ row }) => row.original.email,
  }),

  columnHelper.accessor('status', {
    header: ({ column }) =>
      h(DataTableColumnHeader, {
        column,
        title: t('features.entity.columns.status'),
      }),
    cell: ({ row }) =>
      h(Badge, { variant: row.original.status === 'active' ? 'default' : 'secondary' }, () =>
        row.original.status,
      ),
    filterFn: (row, id, value) => value.includes(row.getValue(id)),
  }),

  columnHelper.display({
    id: 'actions',
    cell: ({ row }) =>
      h(EntityRowActions, {
        entity: row.original,
        onEdit: () => handleEdit(row.original),
        onDelete: () => handleDelete(row.original),
      }),
  }),
])
```

## Centralized Feature Files

### schema.ts

```typescript
import { z } from 'zod'

export const entityCreateSchema = z.object({
  projectId: z.string().length(14),
  name: z.string().min(1, 'Name is required').max(50, 'Name must be 50 characters or less'),
  email: z.string().email('Must be a valid email address'),
})

export type EntityCreateFormData = z.infer<typeof entityCreateSchema>

export const entityUpdateSchema = entityCreateSchema.extend({
  id: z.string().length(14),
})

export type EntityUpdateFormData = z.infer<typeof entityUpdateSchema>
```

### error.ts

```typescript
import type { ComputedRef } from 'vue'

export function getConnectRPCError(
  validationErrors: ComputedRef<Record<string, string[]>>,
  fieldName: string,
): string {
  const errors =
    validationErrors.value[fieldName] ||
    validationErrors.value[`value.${fieldName}`]
  return errors?.[0] || ''
}

export function hasConnectRPCError(
  validationErrors: ComputedRef<Record<string, string[]>>,
  fieldName: string,
): boolean {
  return !!(
    validationErrors.value[fieldName] ||
    validationErrors.value[`value.${fieldName}`]
  )
}
```

### constants.ts

```typescript
export const ENTITY_STATUS_OPTIONS = [
  { value: 'active', label: 'Active' },
  { value: 'inactive', label: 'Inactive' },
] as const

export type EntityStatus = (typeof ENTITY_STATUS_OPTIONS)[number]['value']
```
