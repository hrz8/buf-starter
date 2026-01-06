# Task T10: IAM Frontend Foundation (TransferList + Services)

**Story Reference:** US3-iam-core-entities-and-mappings.md
**Type:** Frontend Foundation
**Priority:** High
**Estimated Effort:** 12-15 hours
**Prerequisites:** T9 (Backend mapping operations must be complete)

## Objective

Build the frontend foundation for IAM features including a reusable TransferList component (composed from shadcn-vue primitives), repositories, service composables, and feature organization files with inline permission creation support.

## Acceptance Criteria

- [x] TransferList component built from shadcn-vue primitives (Command, Button, etc.)
- [x] InlinePermissionCreateDialog component for creating permissions on-the-fly
- [x] Repositories for all 4 services (User, Role, Permission, IAMMapper)
- [x] Service composables with reactive state management
- [x] Feature organization files (schema.ts, error.ts, constants.ts) for each domain
- [x] Description tooltip/info icon in TransferList for permissions
- [x] Immediate save with optimistic updates for mapping operations
- [x] Comprehensive error handling with toast notifications
- [x] Components follow vee-validate patterns (isLoading starts true, no :key)

## Technical Requirements

### 1. TransferList Component

#### frontend/app/components/ui/transfer-list/TransferList.vue

**Component Composition (shadcn-vue primitives):**

```vue
<template>
  <div class="grid grid-cols-[1fr_auto_1fr] gap-4 w-full">
    <!-- Available List -->
    <div class="space-y-2">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-medium">Available {{ label }}</h3>
        <Badge variant="secondary">{{ availableItems.length }}</Badge>
      </div>

      <Command class="border rounded-lg">
        <CommandInput :placeholder="`Search ${label.toLowerCase()}...`" />
        <CommandList>
          <CommandEmpty>No {{ label.toLowerCase() }} found.</CommandEmpty>
          <CommandGroup>
            <CommandItem
              v-for="item in availableItems"
              :key="item.id"
              :value="item.id"
              @select="toggleAvailableSelection(item.id)"
            >
              <Checkbox
                :checked="availableSelected.includes(item.id)"
                class="mr-2"
              />
              <div class="flex-1">
                <span>{{ getItemLabel(item) }}</span>
                <!-- Description Tooltip for Permissions -->
                <TooltipProvider v-if="showTooltip && item.description">
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <Button variant="ghost" size="icon" class="h-4 w-4 ml-1">
                        <Info class="h-3 w-3" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p class="text-xs">{{ item.description }}</p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </div>
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>

      <!-- Inline Create Button (for permissions only) -->
      <Button
        v-if="allowInlineCreate"
        variant="outline"
        size="sm"
        class="w-full"
        @click="openInlineCreateDialog"
      >
        <Plus class="h-4 w-4 mr-2" />
        Create New {{ singularLabel }}
      </Button>
    </div>

    <!-- Arrow Buttons -->
    <div class="flex flex-col justify-center gap-2">
      <Button
        variant="outline"
        size="icon"
        :disabled="availableSelected.length === 0 || isLoading"
        @click="assignSelected"
      >
        <ChevronRight class="h-4 w-4" />
      </Button>
      <Button
        variant="outline"
        size="icon"
        :disabled="assignedSelected.length === 0 || isLoading"
        @click="removeSelected"
      >
        <ChevronLeft class="h-4 w-4" />
      </Button>
    </div>

    <!-- Assigned List -->
    <div class="space-y-2">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-medium">Assigned {{ label }}</h3>
        <Badge variant="secondary">{{ assignedItems.length }}</Badge>
      </div>

      <Command class="border rounded-lg">
        <CommandInput :placeholder="`Search ${label.toLowerCase()}...`" />
        <CommandList>
          <CommandEmpty>No {{ label.toLowerCase() }} assigned.</CommandEmpty>
          <CommandGroup>
            <CommandItem
              v-for="item in assignedItems"
              :key="item.id"
              :value="item.id"
              @select="toggleAssignedSelection(item.id)"
            >
              <Checkbox
                :checked="assignedSelected.includes(item.id)"
                class="mr-2"
              />
              <div class="flex-1">
                <span>{{ getItemLabel(item) }}</span>
                <!-- Description Tooltip for Permissions -->
                <TooltipProvider v-if="showTooltip && item.description">
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <Button variant="ghost" size="icon" class="h-4 w-4 ml-1">
                        <Info class="h-3 w-3" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p class="text-xs">{{ item.description }}</p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </div>
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>
    </div>
  </div>

  <!-- Inline Create Dialog -->
  <InlinePermissionCreateDialog
    v-if="allowInlineCreate"
    v-model:open="inlineCreateDialogOpen"
    @created="handleInlineCreated"
  />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Command,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
} from '@/components/ui/command'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { ChevronRight, ChevronLeft, Plus, Info } from 'lucide-vue-next'
import InlinePermissionCreateDialog from './InlinePermissionCreateDialog.vue'

interface TransferListItem {
  id: string
  [key: string]: any
}

interface Props {
  availableItems: TransferListItem[]
  assignedItems: TransferListItem[]
  label: string // e.g., "Roles", "Permissions"
  singularLabel: string // e.g., "Role", "Permission"
  labelKey?: string // Field to use for display (default: "name")
  isLoading?: boolean
  allowInlineCreate?: boolean // Only true for permissions
  showTooltip?: boolean // Show description tooltip
}

const props = withDefaults(defineProps<Props>(), {
  labelKey: 'name',
  isLoading: false,
  allowInlineCreate: false,
  showTooltip: false,
})

const emit = defineEmits<{
  assign: [ids: string[]]
  remove: [ids: string[]]
}>()

const availableSelected = ref<string[]>([])
const assignedSelected = ref<string[]>([])
const inlineCreateDialogOpen = ref(false)

function toggleAvailableSelection(id: string) {
  const index = availableSelected.value.indexOf(id)
  if (index > -1) {
    availableSelected.value.splice(index, 1)
  } else {
    availableSelected.value.push(id)
  }
}

function toggleAssignedSelection(id: string) {
  const index = assignedSelected.value.indexOf(id)
  if (index > -1) {
    assignedSelected.value.splice(index, 1)
  } else {
    assignedSelected.value.push(id)
  }
}

function assignSelected() {
  if (availableSelected.value.length > 0) {
    emit('assign', [...availableSelected.value])
    availableSelected.value = []
  }
}

function removeSelected() {
  if (assignedSelected.value.length > 0) {
    emit('remove', [...assignedSelected.value])
    assignedSelected.value = []
  }
}

function getItemLabel(item: TransferListItem): string {
  return item[props.labelKey] || item.id
}

function openInlineCreateDialog() {
  inlineCreateDialogOpen.value = true
}

function handleInlineCreated(permission: TransferListItem) {
  // Auto-assign newly created permission
  emit('assign', [permission.id])
}
</script>
```

**Key Features:**
- Dual Command components for searchable, filterable lists
- Multi-select with checkboxes
- Arrow buttons for assign/remove
- Badge counters for item counts
- Empty states when lists are empty
- Inline create button (permissions only)
- Description tooltip with info icon (permissions only)
- Immediate emit of assign/remove events

#### frontend/app/components/ui/transfer-list/InlinePermissionCreateDialog.vue

```vue
<template>
  <Dialog v-model:open="isOpen">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Create New Permission</DialogTitle>
        <DialogDescription>
          Create a permission to assign immediately. Name is required (e.g., "project:read"), description is optional.
        </DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div class="space-y-2">
          <Label for="permission-name">Permission Name *</Label>
          <Input
            id="permission-name"
            v-model="name"
            placeholder="project:read"
            :disabled="isCreating"
          />
          <p class="text-xs text-muted-foreground">
            Format: alphanumeric, underscores, and colons (e.g., "project:read:metadata")
          </p>
        </div>

        <div class="space-y-2">
          <Label for="permission-description">Description (Optional)</Label>
          <Textarea
            id="permission-description"
            v-model="description"
            placeholder="Access read project"
            :disabled="isCreating"
          />
        </div>

        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            @click="handleCancel"
            :disabled="isCreating"
          >
            Cancel
          </Button>
          <Button type="submit" :disabled="!name || isCreating">
            <Loader2 v-if="isCreating" class="mr-2 h-4 w-4 animate-spin" />
            Create & Assign
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Loader2 } from 'lucide-vue-next'
import { useToast } from '@/components/ui/toast/use-toast'
import { usePermissionService } from '@/composables/services/usePermissionService'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  created: [permission: any]
}>()

const { toast } = useToast()
const permissionService = usePermissionService()

const name = ref('')
const description = ref('')
const isCreating = ref(false)

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value),
})

// Reset form when dialog closes
watch(isOpen, (newValue) => {
  if (!newValue) {
    name.value = ''
    description.value = ''
  }
})

async function handleSubmit() {
  if (!name.value) return

  isCreating.value = true

  try {
    const created = await permissionService.create({
      name: name.value,
      description: description.value || undefined,
      effect: 'allow',
    })

    toast({
      title: 'Permission created',
      description: `Permission "${name.value}" has been created and assigned.`,
    })

    emit('created', created)
    isOpen.value = false
  } catch (error: any) {
    toast({
      title: 'Failed to create permission',
      description: error.message || 'An error occurred',
      variant: 'destructive',
    })
    // Keep dialog open for retry
  } finally {
    isCreating.value = false
  }
}

function handleCancel() {
  isOpen.value = false
}
</script>
```

**Key Features:**
- Simple form with name (required) and description (optional)
- Client-side validation
- Toast notification on success/error
- Keeps dialog open on error for retry
- Resets form on close
- Emits created permission for auto-assignment

### 2. Repositories

Create Connect-RPC client wrappers for all 4 services.

#### frontend/shared/repository/user.ts

```typescript
import { createPromiseClient } from '@connectrpc/connect'
import { UserService } from '@/gen/altalune/v1/user_connect'
import { useConnectTransport } from '@/composables/useConnectTransport'
import type {
  QueryUsersRequest,
  CreateUserRequest,
  GetUserRequest,
  UpdateUserRequest,
  DeleteUserRequest,
  ActivateUserRequest,
  DeactivateUserRequest,
} from '@/gen/altalune/v1/user_pb'

export function useUserRepository() {
  const transport = useConnectTransport()
  const client = createPromiseClient(UserService, transport)

  return {
    query: (request: QueryUsersRequest) => client.queryUsers(request),
    create: (request: CreateUserRequest) => client.createUser(request),
    get: (request: GetUserRequest) => client.getUser(request),
    update: (request: UpdateUserRequest) => client.updateUser(request),
    delete: (request: DeleteUserRequest) => client.deleteUser(request),
    activate: (request: ActivateUserRequest) => client.activateUser(request),
    deactivate: (request: DeactivateUserRequest) => client.deactivateUser(request),
  }
}
```

#### frontend/shared/repository/role.ts

Similar structure with RoleService operations.

#### frontend/shared/repository/permission.ts

Similar structure with PermissionService operations.

#### frontend/shared/repository/iam_mapper.ts

```typescript
import { createPromiseClient } from '@connectrpc/connect'
import { IAMMapperService } from '@/gen/altalune/v1/iam_mapper_connect'
import { useConnectTransport } from '@/composables/useConnectTransport'
import type {
  AssignUserRolesRequest,
  RemoveUserRolesRequest,
  GetUserRolesRequest,
  AssignRolePermissionsRequest,
  RemoveRolePermissionsRequest,
  GetRolePermissionsRequest,
  AssignUserPermissionsRequest,
  RemoveUserPermissionsRequest,
  GetUserPermissionsRequest,
  AssignProjectMembersRequest,
  RemoveProjectMembersRequest,
  GetProjectMembersRequest,
} from '@/gen/altalune/v1/iam_mapper_pb'

export function useIAMMapperRepository() {
  const transport = useConnectTransport()
  const client = createPromiseClient(IAMMapperService, transport)

  return {
    // User-Role Mappings
    assignUserRoles: (request: AssignUserRolesRequest) => client.assignUserRoles(request),
    removeUserRoles: (request: RemoveUserRolesRequest) => client.removeUserRoles(request),
    getUserRoles: (request: GetUserRolesRequest) => client.getUserRoles(request),

    // Role-Permission Mappings
    assignRolePermissions: (request: AssignRolePermissionsRequest) =>
      client.assignRolePermissions(request),
    removeRolePermissions: (request: RemoveRolePermissionsRequest) =>
      client.removeRolePermissions(request),
    getRolePermissions: (request: GetRolePermissionsRequest) =>
      client.getRolePermissions(request),

    // User-Permission Mappings
    assignUserPermissions: (request: AssignUserPermissionsRequest) =>
      client.assignUserPermissions(request),
    removeUserPermissions: (request: RemoveUserPermissionsRequest) =>
      client.removeUserPermissions(request),
    getUserPermissions: (request: GetUserPermissionsRequest) =>
      client.getUserPermissions(request),

    // Project Members
    assignProjectMembers: (request: AssignProjectMembersRequest) =>
      client.assignProjectMembers(request),
    removeProjectMembers: (request: RemoveProjectMembersRequest) =>
      client.removeProjectMembers(request),
    getProjectMembers: (request: GetProjectMembersRequest) =>
      client.getProjectMembers(request),
  }
}
```

### 3. Service Composables

#### frontend/app/composables/services/useUserService.ts

```typescript
import { ref } from 'vue'
import { useUserRepository } from '@/shared/repository/user'
import type { User } from '@/gen/altalune/v1/user_pb'

export function useUserService() {
  const repository = useUserRepository()
  const users = ref<User[]>([])
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  async function query(searchQuery: string = '') {
    isLoading.value = true
    error.value = null

    try {
      const response = await repository.query({ query: searchQuery })
      users.value = response.users
      return response.users
    } catch (err: any) {
      error.value = err
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function create(data: { email: string; firstName?: string; lastName?: string }) {
    isLoading.value = true
    error.value = null

    try {
      const response = await repository.create(data)
      return response.user
    } catch (err: any) {
      error.value = err
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function get(userId: string) {
    isLoading.value = true
    error.value = null

    try {
      const response = await repository.get({ userId })
      return response.user
    } catch (err: any) {
      error.value = err
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function update(userId: string, data: Partial<User>) {
    isLoading.value = true
    error.value = null

    try {
      const response = await repository.update({ userId, ...data })
      return response.user
    } catch (err: any) {
      error.value = err
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function remove(userId: string) {
    isLoading.value = true
    error.value = null

    try {
      await repository.delete({ userId })
    } catch (err: any) {
      error.value = err
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function activate(userId: string) {
    isLoading.value = true
    error.value = null

    try {
      const response = await repository.activate({ userId })
      return response.user
    } catch (err: any) {
      error.value = err
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function deactivate(userId: string) {
    isLoading.value = true
    error.value = null

    try {
      const response = await repository.deactivate({ userId })
      return response.user
    } catch (err: any) {
      error.value = err
      throw err
    } finally {
      isLoading.value = false
    }
  }

  return {
    users,
    isLoading,
    error,
    query,
    create,
    get,
    update,
    remove,
    activate,
    deactivate,
  }
}
```

#### frontend/app/composables/services/useRoleService.ts

Similar structure with role-specific operations.

#### frontend/app/composables/services/usePermissionService.ts

Similar structure with permission-specific operations.

#### frontend/app/composables/services/useIAMMapperService.ts

```typescript
import { ref } from 'vue'
import { useIAMMapperRepository } from '@/shared/repository/iam_mapper'
import type { Role, Permission } from '@/gen/altalune/v1/user_pb'

export function useIAMMapperService() {
  const repository = useIAMMapperRepository()
  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  // User-Role Mappings
  async function assignUserRoles(userId: string, roleIds: string[]) {
    isLoading.value = true
    error.value = null

    try {
      await repository.assignUserRoles({ userId, roleIds })
    } catch (err: any) {
      error.value = err
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function removeUserRoles(userId: string, roleIds: string[]) {
    isLoading.value = true
    error.value = null

    try {
      await repository.removeUserRoles({ userId, roleIds })
    } catch (err: any) {
      error.value = err
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function getUserRoles(userId: string): Promise<Role[]> {
    isLoading.value = true
    error.value = null

    try {
      const response = await repository.getUserRoles({ userId })
      return response.roles
    } catch (err: any) {
      error.value = err
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Similar methods for role-permission and user-permission mappings...

  return {
    isLoading,
    error,
    assignUserRoles,
    removeUserRoles,
    getUserRoles,
    // ... other mapping operations
  }
}
```

### 4. Feature Organization Files

#### frontend/app/components/features/iam/user/schema.ts

```typescript
import { z } from 'zod'

export const userSchema = z.object({
  email: z.string().email('Invalid email format').min(1, 'Email is required'),
  firstName: z.string().min(1, 'First name is required').max(100),
  lastName: z.string().min(1, 'Last name is required').max(100),
})

export type UserFormData = z.infer<typeof userSchema>
```

#### frontend/app/components/features/iam/user/error.ts

```typescript
import { ConnectError } from '@connectrpc/connect'

export function getConnectRPCError(error: unknown): string {
  if (error instanceof ConnectError) {
    return error.message
  }
  return 'An unexpected error occurred'
}

export function hasConnectRPCError(error: unknown, code: number): boolean {
  return error instanceof ConnectError && error.code === code
}
```

#### frontend/app/components/features/iam/user/constants.ts

```typescript
export const USER_ACTIVE_OPTIONS = [
  { label: 'Active', value: 'true' },
  { label: 'Inactive', value: 'false' },
] as const
```

Similar files for role and permission features.

#### frontend/app/components/features/iam/permission/schema.ts

```typescript
import { z } from 'zod'

export const permissionSchema = z.object({
  name: z
    .string()
    .min(2, 'Name must be at least 2 characters')
    .max(100)
    .regex(/^[a-zA-Z0-9_:]+$/, 'Name can only contain letters, numbers, underscores, and colons'),
  effect: z.enum(['allow', 'deny']).default('allow'),
  description: z.string().max(500).optional(),
})

export type PermissionFormData = z.infer<typeof permissionSchema>
```

#### frontend/app/components/features/iam/permission/constants.ts

```typescript
export const PERMISSION_EFFECT_OPTIONS = [
  { label: 'Allow', value: 'allow' },
  { label: 'Deny', value: 'deny' },
] as const
```

## Implementation Notes

### TransferList Component Patterns

1. **Command Component as Base**: Use shadcn-vue Command component which has built-in:
   - Search/filter functionality via CommandInput
   - Selection state management
   - Keyboard navigation
   - Empty states via CommandEmpty

2. **Multi-Select Pattern**: Track selected IDs in separate arrays for available/assigned lists

3. **Immediate Save**: Emit assign/remove events immediately when arrow buttons clicked - parent handles API calls with optimistic updates

4. **Description Tooltip**: Use shadcn-vue Tooltip component with info icon button for permissions

### Inline Permission Creation Flow

1. User clicks "+ Create New Permission" in TransferList
2. Dialog opens with name (required) and description (optional) fields
3. On submit, call permission service create()
4. On success:
   - Show toast notification
   - Close dialog
   - Emit 'created' event with new permission
   - Parent component auto-assigns to user/role
5. On error:
   - Show toast with error message
   - Keep dialog open with input preserved
   - User can retry

### vee-validate Patterns

**CRITICAL**: Follow these patterns to avoid FormField errors:

```typescript
// CORRECT: isLoading starts as true
const isLoading = ref(true)

// Load data then set to false
onMounted(async () => {
  await loadData()
  isLoading.value = false
})

// NO :key on FormField
<FormField v-slot="{ componentField }" name="email">
  <!-- No :key here! -->
</FormField>

// Simple conditional rendering
<template v-if="!isLoading">
  <FormField>...</FormField>
</template>
```

### Optimistic Updates

For mapping operations, implement optimistic updates:

```typescript
async function handleAssign(permissionIds: string[]) {
  // Optimistic update
  const newlyAssigned = availablePermissions.value.filter(p => permissionIds.includes(p.id))
  assignedPermissions.value.push(...newlyAssigned)
  availablePermissions.value = availablePermissions.value.filter(p => !permissionIds.includes(p.id))

  try {
    await iamMapperService.assignUserPermissions(userId, permissionIds)
    toast({ title: 'Permissions assigned' })
  } catch (error) {
    // Rollback on error
    assignedPermissions.value = assignedPermissions.value.filter(p => !permissionIds.includes(p.id))
    availablePermissions.value.push(...newlyAssigned)
    toast({ title: 'Failed to assign', variant: 'destructive' })
  }
}
```

## Files to Create

**TransferList Component:**
- `frontend/app/components/ui/transfer-list/TransferList.vue`
- `frontend/app/components/ui/transfer-list/InlinePermissionCreateDialog.vue`

**Repositories:**
- `frontend/shared/repository/user.ts`
- `frontend/shared/repository/role.ts`
- `frontend/shared/repository/permission.ts`
- `frontend/shared/repository/iam_mapper.ts`

**Service Composables:**
- `frontend/app/composables/services/useUserService.ts`
- `frontend/app/composables/services/useRoleService.ts`
- `frontend/app/composables/services/usePermissionService.ts`
- `frontend/app/composables/services/useIAMMapperService.ts`

**Feature Organization (User):**
- `frontend/app/components/features/iam/user/schema.ts`
- `frontend/app/components/features/iam/user/error.ts`
- `frontend/app/components/features/iam/user/constants.ts`

**Feature Organization (Role):**
- `frontend/app/components/features/iam/role/schema.ts`
- `frontend/app/components/features/iam/role/error.ts`

**Feature Organization (Permission):**
- `frontend/app/components/features/iam/permission/schema.ts`
- `frontend/app/components/features/iam/permission/error.ts`
- `frontend/app/components/features/iam/permission/constants.ts`

## Commands to Run

```bash
# Install any missing dependencies (if needed)
cd frontend && pnpm install

# Start development server
cd frontend && pnpm dev

# Lint and fix
cd frontend && pnpm lint:fix
```

## Definition of Done

- [x] TransferList component renders correctly
- [x] Command component provides search/filter functionality
- [x] Multi-select with checkboxes works
- [x] Arrow buttons assign/remove items correctly
- [x] Badge counters display accurate counts
- [x] Empty states show when lists are empty
- [x] "+ Create New Permission" button appears for permissions only
- [x] InlinePermissionCreateDialog opens and functions correctly
- [x] Description tooltip shows on hover/info icon for permissions
- [x] All 4 repositories created and working
- [x] All 4 service composables created with reactive state
- [x] Feature organization files (schema, error, constants) created
- [x] No TypeScript errors (only minor lint warnings)
- [x] Components follow vee-validate patterns
- [x] Inline creation creates permission and emits event
- [x] Error handling with toast notifications working

## Dependencies

- T9 (backend mapping operations) must be complete
- shadcn-vue components (Command, Button, Checkbox, Dialog, Tooltip, etc.)
- Connect-RPC client libraries
- Zod for schema validation
- lucide-vue-next for icons

## Risk Factors

- **Low risk**: Composing from existing shadcn-vue primitives
- **Watch out for**: vee-validate FormField patterns (isLoading, no :key)
- **Test carefully**: Inline permission creation and auto-assignment flow
- **Verify**: Description tooltip works correctly and doesn't clutter UI
