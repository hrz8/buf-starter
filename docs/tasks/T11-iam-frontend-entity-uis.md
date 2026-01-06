# Task T11: IAM Frontend Entity UIs (Users, Roles, Permissions)

**Story Reference:** US3-iam-core-entities-and-mappings.md
**Type:** Frontend UI Implementation
**Priority:** High
**Estimated Effort:** 15-18 hours
**Prerequisites:** T10 (Frontend foundation must be complete)

## Objective

Implement complete user interfaces for all three IAM entities (Users, Roles, Permissions) including tables, CRUD forms, TransferList integration for mappings, row actions, and view mappings features following established patterns.

## Acceptance Criteria

- [ ] User feature complete (Table, Create, Update with 4 tabs, Delete, RowActions, ViewMappings)
- [ ] Role feature complete (Table, Create, Update with 2 tabs, Delete, RowActions)
- [ ] Permission feature complete (Table, Create, Update, Delete, RowActions)
- [ ] All tables support pagination, filtering, sorting
- [ ] TransferList integration in User and Role update forms
- [ ] Inline permission creation in TransferList
- [ ] Description tooltip/info icon for permissions
- [ ] View mappings sheet (read-only) for users
- [ ] vee-validate patterns followed (isLoading starts true, no :key)
- [ ] Immediate mapping save with optimistic updates
- [ ] Comprehensive error handling
- [ ] Three pages created (/iam/users, /iam/roles, /iam/permissions)

## Technical Requirements

### 1. Users Feature

#### File Structure
```
frontend/app/components/features/iam/user/
├── UserTable.vue
├── UserCreateSheet.vue
├── UserCreateForm.vue
├── UserUpdateSheet.vue
├── UserUpdateForm.vue        (4 tabs: Profile, Roles, Permissions, Projects)
├── UserDeleteDialog.vue
├── UserRowActions.vue
├── UserViewMappingsSheet.vue
├── schema.ts                  (Already created in T10)
├── error.ts                   (Already created in T10)
├── constants.ts               (Already created in T10)
└── index.ts                   (Export all components)
```

#### UserTable.vue

**Key Features:**
- Data table with columns: Email, Name, Status (Active/Inactive), Created At
- Search by email
- Filter by active status
- Sort by email, created_at
- Row actions (Edit, Delete, View Mappings, Activate/Deactivate)
- Pagination

**Example Column Definition:**
```typescript
const columns: ColumnDef<User>[] = [
  {
    accessorKey: 'email',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Email" />
    ),
  },
  {
    accessorFn: (row) => `${row.firstName} ${row.lastName}`,
    id: 'name',
    header: 'Name',
  },
  {
    accessorKey: 'isActive',
    header: 'Status',
    cell: ({ row }) => (
      <Badge variant={row.original.isActive ? 'success' : 'secondary'}>
        {row.original.isActive ? 'Active' : 'Inactive'}
      </Badge>
    ),
  },
  {
    accessorKey: 'createdAt',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Created" />
    ),
    cell: ({ row }) => formatDate(row.original.createdAt),
  },
  {
    id: 'actions',
    cell: ({ row }) => <UserRowActions user={row.original} />,
  },
]
```

#### UserCreateForm.vue

**Fields:**
- Email (required, email validation)
- First Name (required)
- Last Name (required)

**vee-validate Pattern:**
```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { userSchema } from './schema'
import type { UserFormData } from './schema'

const isLoading = ref(true) // CRITICAL: Start as true

const form = useForm({
  validationSchema: toTypedSchema(userSchema),
})

onMounted(() => {
  // Set to false after any async operations
  isLoading.value = false
})

async function onSubmit(values: UserFormData) {
  try {
    await userService.create(values)
    toast({ title: 'User created successfully' })
    emit('success')
  } catch (error) {
    toast({
      title: 'Failed to create user',
      description: getConnectRPCError(error),
      variant: 'destructive',
    })
  }
}
</script>

<template>
  <form @submit="form.handleSubmit(onSubmit)">
    <template v-if="!isLoading">
      <!-- NO :key on FormField -->
      <FormField v-slot="{ componentField }" name="email">
        <FormItem>
          <FormLabel>Email</FormLabel>
          <FormControl>
            <Input type="email" v-bind="componentField" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="firstName">
        <FormItem>
          <FormLabel>First Name</FormLabel>
          <FormControl>
            <Input v-bind="componentField" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="lastName">
        <FormItem>
          <FormLabel>Last Name</FormLabel>
          <FormControl>
            <Input v-bind="componentField" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>
    </template>

    <Button type="submit" :disabled="isLoading">
      Create User
    </Button>
  </form>
</template>
```

#### UserUpdateForm.vue (Most Complex - 4 Tabs)

**Tab Structure:**
1. **Profile Tab**: Email, First Name, Last Name
2. **Roles Tab**: TransferList for user-role mappings
3. **Permissions Tab**: TransferList for user-permission mappings (with inline creation)
4. **Projects Tab**: TransferList for project member mappings

**Key Implementation:**
```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import TransferList from '@/components/ui/transfer-list/TransferList.vue'
import { useUserService } from '@/composables/services/useUserService'
import { useRoleService } from '@/composables/services/useRoleService'
import { usePermissionService } from '@/composables/services/usePermissionService'
import { useIAMMapperService } from '@/composables/services/useIAMMapperService'

const props = defineProps<{
  userId: string
}>()

const isLoading = ref(true)
const user = ref<User | null>(null)

// Roles
const allRoles = ref<Role[]>([])
const assignedRoles = ref<Role[]>([])
const availableRoles = computed(() =>
  allRoles.value.filter(r => !assignedRoles.value.some(ar => ar.id === r.id))
)

// Permissions
const allPermissions = ref<Permission[]>([])
const assignedPermissions = ref<Permission[]>([])
const availablePermissions = computed(() =>
  allPermissions.value.filter(p => !assignedPermissions.value.some(ap => ap.id === p.id))
)

onMounted(async () => {
  await Promise.all([
    loadUser(),
    loadRoles(),
    loadPermissions(),
  ])
  isLoading.value = false
})

async function loadUser() {
  user.value = await userService.get(props.userId)
}

async function loadRoles() {
  allRoles.value = await roleService.query()
  assignedRoles.value = await iamMapperService.getUserRoles(props.userId)
}

async function loadPermissions() {
  allPermissions.value = await permissionService.query()
  assignedPermissions.value = await iamMapperService.getUserPermissions(props.userId)
}

async function handleAssignRoles(roleIds: string[]) {
  // Optimistic update
  const newlyAssigned = allRoles.value.filter(r => roleIds.includes(r.id))
  assignedRoles.value.push(...newlyAssigned)

  try {
    await iamMapperService.assignUserRoles(props.userId, roleIds)
    toast({ title: 'Roles assigned' })
  } catch (error) {
    // Rollback
    assignedRoles.value = assignedRoles.value.filter(r => !roleIds.includes(r.id))
    toast({
      title: 'Failed to assign roles',
      description: getConnectRPCError(error),
      variant: 'destructive',
    })
  }
}

async function handleRemoveRoles(roleIds: string[]) {
  // Optimistic update
  const removed = assignedRoles.value.filter(r => roleIds.includes(r.id))
  assignedRoles.value = assignedRoles.value.filter(r => !roleIds.includes(r.id))

  try {
    await iamMapperService.removeUserRoles(props.userId, roleIds)
    toast({ title: 'Roles removed' })
  } catch (error) {
    // Rollback
    assignedRoles.value.push(...removed)
    toast({
      title: 'Failed to remove roles',
      description: getConnectRPCError(error),
      variant: 'destructive',
    })
  }
}

// Similar handlers for permissions and projects...
</script>

<template>
  <div v-if="!isLoading">
    <Tabs default-value="profile">
      <TabsList>
        <TabsTrigger value="profile">Profile</TabsTrigger>
        <TabsTrigger value="roles">Roles</TabsTrigger>
        <TabsTrigger value="permissions">Permissions</TabsTrigger>
        <TabsTrigger value="projects">Projects</TabsTrigger>
      </TabsList>

      <TabsContent value="profile">
        <!-- Profile form fields -->
      </TabsContent>

      <TabsContent value="roles">
        <TransferList
          :available-items="availableRoles"
          :assigned-items="assignedRoles"
          label="Roles"
          singular-label="Role"
          @assign="handleAssignRoles"
          @remove="handleRemoveRoles"
        />
      </TabsContent>

      <TabsContent value="permissions">
        <TransferList
          :available-items="availablePermissions"
          :assigned-items="assignedPermissions"
          label="Permissions"
          singular-label="Permission"
          allow-inline-create
          show-tooltip
          @assign="handleAssignPermissions"
          @remove="handleRemovePermissions"
        />
      </TabsContent>

      <TabsContent value="projects">
        <!-- Project members TransferList -->
      </TabsContent>
    </Tabs>
  </div>
</template>
```

**Key Features:**
- 4 tabs with separate concerns
- Profile tab uses vee-validate form
- Roles/Permissions/Projects tabs use TransferList
- Permissions tab has inline creation and description tooltips
- Immediate save with optimistic updates
- Separate from profile form submission

#### UserViewMappingsSheet.vue

**Read-only sheet showing:**
- Assigned roles (table with name, description)
- Assigned permissions (table with name, description)
- Project memberships (table with project name, role)

**Purpose**: Quick inspection without editing

#### UserRowActions.vue

**Actions:**
- Edit (opens UserUpdateSheet)
- Delete (opens UserDeleteDialog)
- View Mappings (opens UserViewMappingsSheet)
- Activate/Deactivate (direct API call with confirmation)

### 2. Roles Feature

#### File Structure
```
frontend/app/components/features/iam/role/
├── RoleTable.vue
├── RoleCreateSheet.vue
├── RoleCreateForm.vue
├── RoleUpdateSheet.vue
├── RoleUpdateForm.vue        (2 tabs: Details, Permissions)
├── RoleDeleteDialog.vue
├── RoleRowActions.vue
├── schema.ts                  (Already created in T10)
├── error.ts                   (Already created in T10)
└── index.ts
```

#### RoleTable.vue

**Columns:**
- Name
- Description
- Created At
- Actions

**Features:**
- Search by name
- Sort by name, created_at
- Pagination

#### RoleCreateForm.vue

**Fields:**
- Name (required, 2-100 characters, unique)
- Description (optional, max 500 characters)

#### RoleUpdateForm.vue (2 Tabs)

**Tab Structure:**
1. **Details Tab**: Name, Description (vee-validate form)
2. **Permissions Tab**: TransferList for role-permission mappings (with inline creation)

**Key Features:**
- Details tab separate from Permissions tab
- Permissions tab uses TransferList with inline creation
- Immediate save for permission assignments
- Description tooltips on permissions

#### RoleRowActions.vue

**Actions:**
- Edit
- Delete (check if role is in use before deleting)

### 3. Permissions Feature

#### File Structure
```
frontend/app/components/features/iam/permission/
├── PermissionTable.vue
├── PermissionCreateSheet.vue
├── PermissionCreateForm.vue
├── PermissionUpdateSheet.vue
├── PermissionUpdateForm.vue
├── PermissionDeleteDialog.vue
├── PermissionRowActions.vue
├── schema.ts                  (Already created in T10)
├── error.ts                   (Already created in T10)
├── constants.ts               (Already created in T10)
└── index.ts
```

#### PermissionTable.vue

**Columns:**
- Name (machine-readable: "project:read")
- Description (human-readable: "Access read project")
- Effect (Allow/Deny with badge)
- Created At
- Actions

**Features:**
- Search by name
- Filter by effect (allow/deny)
- Sort by name, created_at
- Pagination

**Description Column Display:**
```vue
{
  accessorKey: 'description',
  header: 'Description',
  cell: ({ row }) => (
    <span class="text-muted-foreground">
      {row.original.description || '-'}
    </span>
  ),
}
```

#### PermissionCreateForm.vue

**Fields:**
- Name (required, regex ^[a-zA-Z0-9_:]+$)
- Effect (required, default "allow")
- Description (optional, max 500 characters)

**Name Field Validation:**
```vue
<FormField v-slot="{ componentField }" name="name">
  <FormItem>
    <FormLabel>Permission Name</FormLabel>
    <FormControl>
      <Input
        v-bind="componentField"
        placeholder="project:read"
      />
    </FormControl>
    <FormDescription>
      Format: letters, numbers, underscores, and colons (e.g., "project:read:metadata")
    </FormDescription>
    <FormMessage />
  </FormItem>
</FormField>
```

#### PermissionUpdateForm.vue

**Fields:**
- Name (read-only, cannot be changed to prevent breaking assignments)
- Effect (can be changed)
- Description (can be changed)

**Note**: Name is read-only in update form to prevent breaking existing mappings.

#### PermissionRowActions.vue

**Actions:**
- Edit
- Delete (check if permission is in use before deleting)

### 4. Pages

#### frontend/app/pages/iam/users/index.vue

```vue
<script setup lang="ts">
import { UserTable } from '@/components/features/iam/user'

definePageMeta({
  breadcrumb: {
    title: 'Users',
    items: [
      { title: 'IAM', to: '/iam' },
      { title: 'Users' },
    ],
  },
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold">Users</h1>
    </div>

    <UserTable />
  </div>
</template>
```

#### frontend/app/pages/iam/roles/index.vue

Similar structure with RoleTable

#### frontend/app/pages/iam/permissions/index.vue

Similar structure with PermissionTable

## Implementation Notes

### vee-validate FormField Best Practices

**CRITICAL PATTERN:**

```typescript
// 1. isLoading MUST start as true
const isLoading = ref(true)

// 2. Load data in onMounted, then set to false
onMounted(async () => {
  await loadData()
  isLoading.value = false
})

// 3. NO :key on FormField
<FormField v-slot="{ componentField }" name="email">
  <!-- Never add :key here -->
</FormField>

// 4. Simple conditional rendering
<template v-if="!isLoading">
  <FormField>...</FormField>
</template>
```

**Why**: This prevents "useFormField should be used within <FormField>" errors.

### Multi-Tab Forms with TransferList

**Key Pattern**: Profile form submission is SEPARATE from mapping operations.

- **Profile Tab**: Traditional form with submit button
- **Mapping Tabs**: Immediate save when assign/remove clicked
- **No Mixed State**: Don't combine profile data with mapping data in same form state

### Optimistic Updates

For all mapping operations, implement optimistic updates:

```typescript
async function handleAssign(ids: string[]) {
  // 1. Optimistic update (immediate UI change)
  const newItems = available.value.filter(item => ids.includes(item.id))
  assigned.value.push(...newItems)
  available.value = available.value.filter(item => !ids.includes(item.id))

  try {
    // 2. API call
    await mapperService.assign(entityId, ids)
    toast({ title: 'Assigned successfully' })
  } catch (error) {
    // 3. Rollback on error
    assigned.value = assigned.value.filter(item => !ids.includes(item.id))
    available.value.push(...newItems)
    toast({
      title: 'Failed to assign',
      description: getConnectRPCError(error),
      variant: 'destructive',
    })
  }
}
```

### Description Tooltip for Permissions

In TransferList, permissions should show description on hover or via info icon:

```vue
<div class="flex items-center">
  <span>{{ permission.name }}</span>
  <TooltipProvider v-if="permission.description">
    <Tooltip>
      <TooltipTrigger as-child>
        <Button variant="ghost" size="icon" class="h-4 w-4 ml-1">
          <Info class="h-3 w-3 text-muted-foreground" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>
        <p class="text-xs max-w-xs">{{ permission.description }}</p>
      </TooltipContent>
    </Tooltip>
  </TooltipProvider>
</div>
```

### Table Search and Filter Patterns

Follow existing DataTable patterns:

```vue
<script setup lang="ts">
const searchQuery = ref('')
const statusFilter = ref<'all' | 'active' | 'inactive'>('all')

const filteredData = computed(() => {
  let result = data.value

  // Search
  if (searchQuery.value) {
    result = result.filter(item =>
      item.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  // Filter
  if (statusFilter.value !== 'all') {
    const isActive = statusFilter.value === 'active'
    result = result.filter(item => item.isActive === isActive)
  }

  return result
})
</script>
```

### Delete Confirmation

For delete operations, check if entity is in use:

```typescript
async function handleDelete() {
  try {
    await roleService.delete(roleId)
    toast({ title: 'Role deleted successfully' })
  } catch (error) {
    if (hasConnectRPCError(error, 60603)) {
      // RoleInUse error
      toast({
        title: 'Cannot delete role',
        description: 'This role is currently assigned to users',
        variant: 'destructive',
      })
    } else {
      toast({
        title: 'Failed to delete role',
        description: getConnectRPCError(error),
        variant: 'destructive',
      })
    }
  }
}
```

## Files to Create

**User Feature (9 files):**
- `frontend/app/components/features/iam/user/UserTable.vue`
- `frontend/app/components/features/iam/user/UserCreateSheet.vue`
- `frontend/app/components/features/iam/user/UserCreateForm.vue`
- `frontend/app/components/features/iam/user/UserUpdateSheet.vue`
- `frontend/app/components/features/iam/user/UserUpdateForm.vue`
- `frontend/app/components/features/iam/user/UserDeleteDialog.vue`
- `frontend/app/components/features/iam/user/UserRowActions.vue`
- `frontend/app/components/features/iam/user/UserViewMappingsSheet.vue`
- `frontend/app/components/features/iam/user/index.ts`

**Role Feature (8 files):**
- `frontend/app/components/features/iam/role/RoleTable.vue`
- `frontend/app/components/features/iam/role/RoleCreateSheet.vue`
- `frontend/app/components/features/iam/role/RoleCreateForm.vue`
- `frontend/app/components/features/iam/role/RoleUpdateSheet.vue`
- `frontend/app/components/features/iam/role/RoleUpdateForm.vue`
- `frontend/app/components/features/iam/role/RoleDeleteDialog.vue`
- `frontend/app/components/features/iam/role/RoleRowActions.vue`
- `frontend/app/components/features/iam/role/index.ts`

**Permission Feature (8 files):**
- `frontend/app/components/features/iam/permission/PermissionTable.vue`
- `frontend/app/components/features/iam/permission/PermissionCreateSheet.vue`
- `frontend/app/components/features/iam/permission/PermissionCreateForm.vue`
- `frontend/app/components/features/iam/permission/PermissionUpdateSheet.vue`
- `frontend/app/components/features/iam/permission/PermissionUpdateForm.vue`
- `frontend/app/components/features/iam/permission/PermissionDeleteDialog.vue`
- `frontend/app/components/features/iam/permission/PermissionRowActions.vue`
- `frontend/app/components/features/iam/permission/index.ts`

**Pages (3 files):**
- `frontend/app/pages/iam/users/index.vue`
- `frontend/app/pages/iam/roles/index.vue`
- `frontend/app/pages/iam/permissions/index.vue`

## Commands to Run

```bash
# Start development server
cd frontend && pnpm dev

# Lint and fix
cd frontend && pnpm lint:fix

# Type check
cd frontend && pnpm type-check
```

## Definition of Done

- [ ] All 28 files created
- [ ] User feature fully functional (CRUD + 4-tab form)
- [ ] Role feature fully functional (CRUD + 2-tab form)
- [ ] Permission feature fully functional (CRUD)
- [ ] All tables render with correct columns
- [ ] Search and filter working on all tables
- [ ] Pagination working on all tables
- [ ] TransferList working in User and Role update forms
- [ ] Inline permission creation working
- [ ] Description tooltips showing for permissions
- [ ] Optimistic updates working for all mappings
- [ ] Error handling comprehensive
- [ ] Toast notifications for all operations
- [ ] View mappings sheet working for users
- [ ] Activate/Deactivate working for users
- [ ] Delete validation (check if in use)
- [ ] vee-validate patterns followed (no FormField errors)
- [ ] All 3 pages render correctly
- [ ] No TypeScript errors
- [ ] No console errors
- [ ] Responsive design working

## Dependencies

- T10 (frontend foundation) must be complete
- TransferList component available
- All service composables and repositories available
- shadcn-vue components (Table, Form, Dialog, Sheet, Tabs, etc.)

## Risk Factors

- **Medium risk**: Complex UserUpdateForm with 4 tabs and multiple TransferLists
- **Watch out for**: vee-validate FormField errors (follow patterns strictly)
- **Test carefully**: Optimistic updates with rollback logic
- **Verify**: Permission name regex validation allows colons
- **Performance**: Multiple TransferLists in tabs may need lazy loading if data sets are large
