# Task T22: OAuth Client Frontend UI Components

**Story Reference:** US6-oauth-client-management.md
**Type:** Frontend UI
**Priority:** High (P0)
**Estimated Effort:** 8-10 hours
**Prerequisites:** T21-oauth-client-frontend-foundation (Repository + service exist)

## Objective

Implement all UI components for OAuth client management including data table, forms, dialogs, and the main page.

## Acceptance Criteria

- [ ] OAuthClientTable component created with Tanstack Table integration
- [ ] OAuthClientCreateSheet + CreateForm components created
- [ ] OAuthClientEditSheet + EditForm components created
- [ ] OAuthClientDeleteDialog component created
- [ ] OAuthClientRevealDialog component created with 30-second timer
- [ ] OAuthClientSecretDisplay component created (one-time display)
- [ ] OAuthClientRowActions component created with role-based actions
- [ ] Main page created at pages/iam/oauth-client/index.vue
- [ ] All forms use vee-validate with Zod schemas
- [ ] Dual-layer validation working (vee-validate + ConnectRPC)
- [ ] Loading states for all operations
- [ ] Error handling comprehensive
- [ ] Role-based UI visibility implemented
- [ ] Default client special handling (badges, disabled actions)
- [ ] Responsive design (mobile + desktop)

## Technical Requirements

### Component Structure Overview

**9 Components + 1 Page**:
1. **OAuthClientTable.vue** - Main data table
2. **OAuthClientCreateSheet.vue** - Sheet wrapper for create
3. **OAuthClientCreateForm.vue** - Create form with validation
4. **OAuthClientEditSheet.vue** - Sheet wrapper for edit
5. **OAuthClientEditForm.vue** - Edit form with validation
6. **OAuthClientDeleteDialog.vue** - Delete confirmation
7. **OAuthClientRevealDialog.vue** - Reveal secret with timer
8. **OAuthClientSecretDisplay.vue** - One-time secret display
9. **OAuthClientRowActions.vue** - Action menu
10. **pages/iam/oauth-client/index.vue** - Main page

### 1. OAuthClientTable.vue

**Features**:
- Tanstack Table with pagination, sorting, filtering
- Columns: name, client_id (masked), redirect URIs count, PKCE badge, dates
- Special badge for default client
- Faceted filters (PKCE, default)
- Row actions menu
- Refresh functionality

**Key Implementation**:

```vue
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useOAuthClientService } from '~/composables/services/useOAuthClientService'
import { useQueryRequest } from '~/composables/useQueryRequest'
import { useDataTableState } from '~/composables/useDataTableState'
import OAuthClientRowActions from './OAuthClientRowActions.vue'

const { queryClients, queryLoading, clients, clientsMeta } = useOAuthClientService()
const projectId = ref('') // From route or context

// Data table state
const dataTableRef = ref()
const { columnFilters, sorting, pagination } = useDataTableState(dataTableRef)

// Query request builder
const queryRequest = useQueryRequest({
  page: computed(() => pagination.value.pageIndex + 1),
  pageSize: computed(() => pagination.value.pageSize),
  keyword: ref(''),
  columnFilters,
  sorting,
})

// Columns definition
const columns = [
  {
    accessorKey: 'name',
    header: 'Client Name',
    cell: ({ row }) => {
      const client = row.original
      return h('div', { class: 'flex items-center gap-2' }, [
        h('span', client.name),
        client.isDefault && h(Badge, { variant: 'secondary' }, 'Default'),
      ])
    },
  },
  {
    accessorKey: 'clientId',
    header: 'Client ID',
    cell: ({ row }) => {
      const clientId = row.original.clientId
      // Mask: show only last 4 chars
      const masked = `****-****-****-${clientId.slice(-4)}`
      return masked
    },
  },
  {
    accessorKey: 'redirectUris',
    header: 'Redirect URIs',
    cell: ({ row }) => row.original.redirectUris.length,
  },
  {
    accessorKey: 'pkceRequired',
    header: 'PKCE',
    cell: ({ row }) => {
      return row.original.pkceRequired
        ? h(Badge, { variant: 'outline' }, 'Required')
        : h('span', { class: 'text-muted-foreground' }, 'Optional')
    },
  },
  {
    accessorKey: 'createdAt',
    header: 'Created',
    cell: ({ row }) => formatDate(row.original.createdAt),
  },
  {
    id: 'actions',
    cell: ({ row }) => {
      return h(OAuthClientRowActions, {
        client: row.original,
        onRefresh: fetchClients,
      })
    },
  },
]

// Fetch clients
async function fetchClients() {
  await queryClients({
    projectId: projectId.value,
    params: queryRequest.value,
  })
}

onMounted(() => {
  fetchClients()
})

// Watch for filter changes
watch([columnFilters, sorting, pagination], () => {
  fetchClients()
})
</script>

<template>
  <div class="space-y-4">
    <!-- Filters -->
    <div class="flex items-center gap-4">
      <Input
        v-model="keyword"
        placeholder="Search clients..."
        class="max-w-sm"
      />
      <Button @click="fetchClients" variant="outline">
        <Icon name="lucide:refresh-cw" class="mr-2 h-4 w-4" />
        Refresh
      </Button>
    </div>

    <!-- Table -->
    <DataTable
      ref="dataTableRef"
      :columns="columns"
      :data="clients"
      :loading="queryLoading"
      :total-rows="clientsMeta?.totalRows"
      @update:pagination="pagination = $event"
      @update:sorting="sorting = $event"
      @update:column-filters="columnFilters = $event"
    />
  </div>
</template>
```

### 2. OAuthClientCreateSheet.vue

**Pattern**: Sheet wrapper that opens/closes form

```vue
<script setup lang="ts">
import { ref } from 'vue'
import OAuthClientCreateForm from './OAuthClientCreateForm.vue'

const isOpen = ref(false)

function handleSuccess() {
  isOpen.value = false
  // Emit refresh event
  emit('refresh')
}

function handleCancel() {
  isOpen.value = false
}
</script>

<template>
  <Sheet v-model:open="isOpen">
    <SheetTrigger as-child>
      <slot />
    </SheetTrigger>
    <SheetContent class="sm:max-w-[600px]">
      <SheetHeader>
        <SheetTitle>{{ t('oauthClient.create') }}</SheetTitle>
        <SheetDescription>
          {{ t('oauthClient.createDescription') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <OAuthClientCreateForm
          @success="handleSuccess"
          @cancel="handleCancel"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
```

### 3. OAuthClientCreateForm.vue

**Critical vee-validate Patterns**:
- ✅ Start `isLoading = ref(true)`
- ✅ NO `:key` on FormField
- ✅ Simple `v-if` rendering
- ✅ Dual validation (vee-validate + ConnectRPC)

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { useOAuthClientService } from '~/composables/services/useOAuthClientService'
import { oauthClientCreateSchema } from './schema'
import { getConnectRPCError, hasConnectRPCError } from './error'
import OAuthClientSecretDisplay from './OAuthClientSecretDisplay.vue'

const emit = defineEmits(['success', 'cancel'])

// Service
const {
  createClient,
  createLoading,
  createValidationErrors,
  clientSecret,
  resetCreateState,
} = useOAuthClientService()

// Loading state - MUST start as true
const isLoading = ref(true)

// Form setup
const formSchema = toTypedSchema(oauthClientCreateSchema)
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    projectId: '', // From context
    name: '',
    redirectUris: [''],
    pkceRequired: false,
    allowedScopes: [],
  },
})

// Redirect URI array management
const redirectUris = ref([''])

function addRedirectUri() {
  redirectUris.value.push('')
}

function removeRedirectUri(index: number) {
  if (redirectUris.value.length > 1) {
    redirectUris.value.splice(index, 1)
  }
}

// Form submission
const onSubmit = form.handleSubmit(async (values) => {
  try {
    await createClient({
      projectId: values.projectId,
      name: values.name,
      redirectUris: values.redirectUris.filter((uri) => uri.trim()),
      pkceRequired: values.pkceRequired,
      allowedScopes: values.allowedScopes,
    })

    // Success - show secret display (don't close immediately)
    // User must acknowledge and copy secret
  } catch (error) {
    console.error('Failed to create client:', error)
  }
})

// Close form
function handleCancel() {
  resetCreateState()
  emit('cancel')
}

// After user acknowledges secret
function handleSecretAcknowledged() {
  resetCreateState()
  emit('success')
}

// Set loading to false after mount
onMounted(() => {
  isLoading.value = false
})
</script>

<template>
  <div v-if="!isLoading" class="space-y-6">
    <!-- Show secret display after creation -->
    <OAuthClientSecretDisplay
      v-if="clientSecret"
      :client-secret="clientSecret"
      @acknowledged="handleSecretAcknowledged"
    />

    <!-- Show form if no secret yet -->
    <form v-else @submit="onSubmit" class="space-y-4">
      <!-- Client Name -->
      <FormField v-slot="{ componentField }" name="name">
        <FormItem>
          <FormLabel>Client Name *</FormLabel>
          <FormControl>
            <Input
              v-bind="componentField"
              placeholder="My Application"
              :disabled="createLoading"
            />
          </FormControl>
          <FormDescription>
            A friendly name for this OAuth client
          </FormDescription>
          <FormMessage />
          <div
            v-if="hasConnectRPCError(createValidationErrors, 'name')"
            class="text-sm text-destructive"
          >
            {{ getConnectRPCError(createValidationErrors, 'name') }}
          </div>
        </FormItem>
      </FormField>

      <!-- Redirect URIs (Array) -->
      <div class="space-y-2">
        <Label>Redirect URIs *</Label>
        <div
          v-for="(uri, index) in redirectUris"
          :key="index"
          class="flex items-center gap-2"
        >
          <FormField v-slot="{ componentField }" :name="`redirectUris.${index}`">
            <FormItem class="flex-1">
              <FormControl>
                <Input
                  v-bind="componentField"
                  placeholder="https://example.com/callback"
                  :disabled="createLoading"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
          <Button
            v-if="redirectUris.length > 1"
            type="button"
            variant="ghost"
            size="icon"
            @click="removeRedirectUri(index)"
            :disabled="createLoading"
          >
            <Icon name="lucide:x" class="h-4 w-4" />
          </Button>
        </div>
        <Button
          type="button"
          variant="outline"
          size="sm"
          @click="addRedirectUri"
          :disabled="createLoading"
        >
          <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
          Add Redirect URI
        </Button>
        <p class="text-sm text-muted-foreground">
          Valid HTTP/HTTPS URLs. Localhost allowed for development.
        </p>
      </div>

      <!-- PKCE Required -->
      <FormField v-slot="{ componentField }" name="pkceRequired">
        <FormItem class="flex items-center justify-between rounded-lg border p-4">
          <div class="space-y-0.5">
            <FormLabel>PKCE Required</FormLabel>
            <FormDescription>
              Enable for public clients (SPAs, mobile apps)
            </FormDescription>
          </div>
          <FormControl>
            <Switch
              v-bind="componentField"
              :disabled="createLoading"
            />
          </FormControl>
        </FormItem>
      </FormField>

      <!-- Actions -->
      <div class="flex justify-end gap-2">
        <Button
          type="button"
          variant="outline"
          @click="handleCancel"
          :disabled="createLoading"
        >
          Cancel
        </Button>
        <Button type="submit" :disabled="createLoading">
          <Icon
            v-if="createLoading"
            name="lucide:loader-2"
            class="mr-2 h-4 w-4 animate-spin"
          />
          Create Client
        </Button>
      </div>
    </form>
  </div>
</template>
```

### 4. OAuthClientSecretDisplay.vue

**One-Time Secret Display** - Shows after creation

```vue
<script setup lang="ts">
import { ref } from 'vue'
import { useClipboard } from '@vueuse/core'

const props = defineProps<{
  clientSecret: string
}>()

const emit = defineEmits(['acknowledged'])

const { copy, copied } = useClipboard()

function copySecret() {
  copy(props.clientSecret)
}

function acknowledge() {
  emit('acknowledged')
}
</script>

<template>
  <Alert variant="default" class="border-yellow-500 bg-yellow-50">
    <Icon name="lucide:alert-triangle" class="h-5 w-5 text-yellow-600" />
    <AlertTitle class="text-lg font-semibold">
      Client Secret Created
    </AlertTitle>
    <AlertDescription class="space-y-4">
      <p class="text-sm text-yellow-800">
        <strong>Important:</strong> Save this secret now - it won't be shown again!
      </p>

      <!-- Secret Display -->
      <div class="rounded-md bg-white p-4 border">
        <div class="flex items-center justify-between gap-4">
          <code class="font-mono text-sm break-all">
            {{ clientSecret }}
          </code>
          <Button
            size="sm"
            variant="outline"
            @click="copySecret"
          >
            <Icon
              :name="copied ? 'lucide:check' : 'lucide:copy'"
              class="h-4 w-4 mr-2"
            />
            {{ copied ? 'Copied!' : 'Copy' }}
          </Button>
        </div>
      </div>

      <!-- Acknowledge Button -->
      <div class="flex justify-end">
        <Button @click="acknowledge" variant="default">
          I've Saved the Secret
        </Button>
      </div>
    </AlertDescription>
  </Alert>
</template>
```

### 5. OAuthClientEditForm.vue

**Similar to CreateForm but**:
- Loads existing client data
- Cannot edit client_id (read-only, copyable)
- Cannot disable PKCE for default client
- All fields optional (only update what changed)

```vue
<script setup lang="ts">
// Similar structure to CreateForm
// Key differences:
// 1. Load client data on mount
// 2. Show client_id (read-only, copyable)
// 3. Check isDefault and disable PKCE toggle if true
// 4. Use updateClient instead of createClient

const props = defineProps<{
  clientId: string
}>()

// Disable PKCE toggle if default client
const isPKCEDisabled = computed(() => {
  return client.value?.isDefault && client.value?.pkceRequired
})
</script>

<template>
  <!-- ... -->

  <!-- PKCE Toggle - Disabled for default client -->
  <FormField v-slot="{ componentField }" name="pkceRequired">
    <FormItem>
      <div class="flex items-center justify-between">
        <FormLabel>PKCE Required</FormLabel>
        <FormControl>
          <Switch
            v-bind="componentField"
            :disabled="updateLoading || isPKCEDisabled"
          />
        </FormControl>
      </div>
      <FormDescription v-if="isPKCEDisabled">
        PKCE cannot be disabled for the default client
      </FormDescription>
    </FormItem>
  </FormField>

  <!-- ... -->
</template>
```

### 6. OAuthClientDeleteDialog.vue

**Confirmation with warnings**

```vue
<script setup lang="ts">
import { ref } from 'vue'
import { useOAuthClientService } from '~/composables/services/useOAuthClientService'

const props = defineProps<{
  client: any
  open: boolean
}>()

const emit = defineEmits(['update:open', 'success'])

const { deleteClient, deleteLoading } = useOAuthClientService()

async function handleDelete() {
  try {
    await deleteClient(props.client.projectId, props.client.id)
    emit('success')
    emit('update:open', false)
  } catch (error) {
    console.error('Delete failed:', error)
  }
}
</script>

<template>
  <AlertDialog :open="open" @update:open="$emit('update:open', $event)">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>Delete OAuth Client?</AlertDialogTitle>
        <AlertDialogDescription class="space-y-2">
          <p>You are about to delete <strong>{{ client.name }}</strong>.</p>
          <p class="text-destructive font-semibold">
            ⚠️ All applications using this client will stop working.
          </p>
          <p class="text-muted-foreground">
            This action cannot be undone.
          </p>
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel :disabled="deleteLoading">
          Cancel
        </AlertDialogCancel>
        <AlertDialogAction
          @click="handleDelete"
          :disabled="deleteLoading"
          class="bg-destructive hover:bg-destructive/90"
        >
          <Icon
            v-if="deleteLoading"
            name="lucide:loader-2"
            class="mr-2 h-4 w-4 animate-spin"
          />
          Delete Client
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
```

### 7. OAuthClientRevealDialog.vue

**Reveal with 30-second auto-hide timer**

```vue
<script setup lang="ts">
import { ref, watch } from 'vue'
import { useOAuthClientService } from '~/composables/services/useOAuthClientService'
import { useClipboard } from '@vueuse/core'

const props = defineProps<{
  client: any
  open: boolean
}>()

const emit = defineEmits(['update:open'])

const {
  revealClientSecret,
  revealLoading,
  revealedSecret,
  resetRevealState,
} = useOAuthClientService()

const { copy, copied } = useClipboard()
const countdown = ref(30)
let intervalId: NodeJS.Timeout | null = null

async function handleReveal() {
  await revealClientSecret(props.client.projectId, props.client.id)
  startCountdown()
}

function startCountdown() {
  countdown.value = 30
  intervalId = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      handleClose()
    }
  }, 1000)
}

function handleClose() {
  if (intervalId) {
    clearInterval(intervalId)
  }
  resetRevealState()
  emit('update:open', false)
}

function copySecret() {
  if (revealedSecret.value) {
    copy(revealedSecret.value)
  }
}

watch(() => props.open, (isOpen) => {
  if (!isOpen && intervalId) {
    clearInterval(intervalId)
  }
})
</script>

<template>
  <Dialog :open="open" @update:open="$emit('update:open', $event)">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Reveal Client Secret</DialogTitle>
        <DialogDescription>
          This will expose the client secret. This action is logged for security.
        </DialogDescription>
      </DialogHeader>

      <!-- Before reveal: Confirmation -->
      <div v-if="!revealedSecret" class="space-y-4">
        <Alert variant="destructive">
          <Icon name="lucide:shield-alert" class="h-5 w-5" />
          <AlertTitle>Security Warning</AlertTitle>
          <AlertDescription>
            The client secret will be displayed. Make sure you're in a secure environment.
          </AlertDescription>
        </Alert>
      </div>

      <!-- After reveal: Display secret -->
      <div v-else class="space-y-4">
        <div class="rounded-md bg-muted p-4">
          <div class="flex items-center justify-between gap-4">
            <code class="font-mono text-sm break-all">
              {{ revealedSecret }}
            </code>
            <Button size="sm" variant="outline" @click="copySecret">
              <Icon
                :name="copied ? 'lucide:check' : 'lucide:copy'"
                class="h-4 w-4 mr-2"
              />
              {{ copied ? 'Copied!' : 'Copy' }}
            </Button>
          </div>
        </div>

        <!-- Countdown timer -->
        <p class="text-sm text-center text-muted-foreground">
          Auto-hiding in {{ countdown }} seconds
        </p>
      </div>

      <DialogFooter>
        <Button
          v-if="!revealedSecret"
          @click="handleReveal"
          :disabled="revealLoading"
        >
          <Icon
            v-if="revealLoading"
            name="lucide:loader-2"
            class="mr-2 h-4 w-4 animate-spin"
          />
          Reveal Secret
        </Button>
        <Button v-else variant="outline" @click="handleClose">
          Close
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
```

### 8. OAuthClientRowActions.vue

**Dropdown menu with role-based actions**

```vue
<script setup lang="ts">
import { computed } from 'vue'
import { OAUTH_CLIENT_PERMISSIONS } from './constants'

const props = defineProps<{
  client: any
}>()

const emit = defineEmits(['refresh'])

const userRole = ref('owner') // From context

// Role-based visibility
const canEdit = computed(() =>
  OAUTH_CLIENT_PERMISSIONS.UPDATE.includes(userRole.value),
)

const canDelete = computed(() =>
  OAUTH_CLIENT_PERMISSIONS.DELETE.includes(userRole.value),
)

const canReveal = computed(() =>
  OAUTH_CLIENT_PERMISSIONS.REVEAL_SECRET.includes(userRole.value),
)

// Disable delete for default client
const isDeleteDisabled = computed(() => props.client.isDefault)

// Dialog states
const isEditOpen = ref(false)
const isDeleteOpen = ref(false)
const isRevealOpen = ref(false)
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" size="icon">
        <Icon name="lucide:more-horizontal" class="h-4 w-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end">
      <!-- Edit -->
      <DropdownMenuItem v-if="canEdit" @click="isEditOpen = true">
        <Icon name="lucide:pencil" class="mr-2 h-4 w-4" />
        Edit
      </DropdownMenuItem>

      <!-- Reveal Secret -->
      <DropdownMenuItem v-if="canReveal" @click="isRevealOpen = true">
        <Icon name="lucide:eye" class="mr-2 h-4 w-4" />
        Reveal Secret
      </DropdownMenuItem>

      <DropdownMenuSeparator />

      <!-- Delete -->
      <DropdownMenuItem
        v-if="canDelete"
        @click="isDeleteOpen = true"
        :disabled="isDeleteDisabled"
        class="text-destructive"
      >
        <Icon name="lucide:trash-2" class="mr-2 h-4 w-4" />
        Delete
        <span v-if="isDeleteDisabled" class="ml-auto text-xs text-muted-foreground">
          (Default)
        </span>
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>

  <!-- Dialogs -->
  <OAuthClientEditSheet
    v-if="canEdit"
    :client-id="client.id"
    v-model:open="isEditOpen"
    @success="$emit('refresh')"
  />

  <OAuthClientRevealDialog
    v-if="canReveal"
    :client="client"
    v-model:open="isRevealOpen"
  />

  <OAuthClientDeleteDialog
    v-if="canDelete"
    :client="client"
    v-model:open="isDeleteOpen"
    @success="$emit('refresh')"
  />
</template>
```

### 9. Main Page (pages/iam/oauth-client/index.vue)

```vue
<script setup lang="ts">
import { ref } from 'vue'
import OAuthClientTable from '~/components/features/oauth-client/OAuthClientTable.vue'
import OAuthClientCreateSheet from '~/components/features/oauth-client/OAuthClientCreateSheet.vue'

const { t } = useI18n()
const tableRef = ref()

function refreshTable() {
  tableRef.value?.refresh()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">
          {{ t('oauthClient.title') }}
        </h1>
        <p class="text-muted-foreground">
          {{ t('oauthClient.description') }}
        </p>
      </div>

      <!-- Create Button -->
      <OAuthClientCreateSheet @refresh="refreshTable">
        <Button>
          <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
          {{ t('oauthClient.create') }}
        </Button>
      </OAuthClientCreateSheet>
    </div>

    <!-- Table -->
    <OAuthClientTable ref="tableRef" />
  </div>
</template>
```

## Files to Create

- `frontend/app/components/features/oauth-client/OAuthClientTable.vue`
- `frontend/app/components/features/oauth-client/OAuthClientCreateSheet.vue`
- `frontend/app/components/features/oauth-client/OAuthClientCreateForm.vue`
- `frontend/app/components/features/oauth-client/OAuthClientEditSheet.vue`
- `frontend/app/components/features/oauth-client/OAuthClientEditForm.vue`
- `frontend/app/components/features/oauth-client/OAuthClientDeleteDialog.vue`
- `frontend/app/components/features/oauth-client/OAuthClientRevealDialog.vue`
- `frontend/app/components/features/oauth-client/OAuthClientSecretDisplay.vue`
- `frontend/app/components/features/oauth-client/OAuthClientRowActions.vue`
- `frontend/app/pages/iam/oauth-client/index.vue`

## Files to Modify

None (all new files)

## Commands to Run

```bash
# Frontend dev server
cd frontend && pnpm dev

# Open browser
open http://localhost:3000/iam/oauth-client

# Type check
pnpm typecheck

# Lint
pnpm lint
```

## Validation Checklist

- [ ] All 9 components compile without errors
- [ ] Main page renders correctly
- [ ] Table displays client data
- [ ] Create form opens and validates
- [ ] Create form submits successfully
- [ ] Secret display shows after creation
- [ ] Edit form loads client data
- [ ] Edit form cannot disable PKCE for default
- [ ] Delete dialog blocks default client
- [ ] Reveal dialog shows secret with timer
- [ ] Row actions menu role-based visibility
- [ ] Responsive on mobile (sheets, dialogs)
- [ ] No vee-validate FormField errors

## Definition of Done

- [ ] All 9 components created and functional
- [ ] Main page created with header and table
- [ ] Forms use vee-validate with Zod schemas
- [ ] Dual-layer validation working
- [ ] Loading states for all operations
- [ ] Error messages displayed correctly
- [ ] Role-based UI visibility enforced
- [ ] Default client special handling (badges, disabled delete)
- [ ] Secret display one-time with copy button
- [ ] Reveal dialog with 30-second auto-hide
- [ ] Responsive design tested (mobile + desktop)
- [ ] No TypeScript errors
- [ ] No console errors
- [ ] All CRUD operations tested manually

## Dependencies

**Internal**:
- T21: Repository, service, schemas, error utilities
- Existing: shadcn-vue components
- Existing: DataTable component
- Existing: useI18n, useClipboard

## Risk Factors

- **High Risk**: vee-validate FormField errors
  - **Mitigation**: Follow critical best practices (isLoading=true, no :key, simple v-if)
- **Medium Risk**: Sheet/Dialog inside Dropdown
  - **Mitigation**: Use manual v-model:open control, nextTick if needed
- **Low Risk**: Responsive design issues
  - **Mitigation**: Test on mobile viewport, use responsive classes

## Notes

### Critical vee-validate FormField Best Practices

**MUST FOLLOW**:
1. ✅ Start `isLoading = ref(true)` - Stable provide/inject
2. ✅ NO `:key` attributes on FormField components
3. ✅ Simple `v-if`/`v-else-if` conditional rendering
4. ✅ No Teleport around FormFields
5. ✅ Set `isLoading.value = false` in onMounted()

**If you get "useFormField should be used within FormField" error**:
- Check isLoading starts as true
- Remove any :key attributes
- Simplify conditional rendering

### Secret Handling Patterns

**One-Time Display** (After Creation):
1. Backend returns plaintext in CreateClientResponse
2. Store in service state: `clientSecret`
3. Display with OAuthClientSecretDisplay
4. User must copy before closing
5. Clear state when acknowledged

**Reveal Flow** (Existing Client):
1. Confirmation dialog with warning
2. Call revealClientSecret()
3. Backend returns hashed secret (audit logged)
4. Display in modal with 30-second timer
5. Auto-hide after countdown
6. Clear state on close

### Role-Based Visibility

**Permissions** (from constants.ts):
- CREATE: owner, admin
- UPDATE: owner, admin
- DELETE: owner only
- REVEAL_SECRET: owner, admin
- VIEW: owner, admin, member (read-only)

**Implementation**:
```typescript
const canDelete = computed(() =>
  OAUTH_CLIENT_PERMISSIONS.DELETE.includes(userRole.value)
)
```

### Default Client Special Handling

**UI Indicators**:
- Badge: "Default" shown next to name
- Delete button: Disabled with tooltip
- PKCE toggle: Disabled if default and PKCE=true
- Visual distinction in table (badge color)

**Enforcement**:
- Backend blocks deletion (returns error)
- UI disables button (better UX)
- Tooltip explains why disabled
