# Task T12: IAM Navigation, Integration & Polish

**Story Reference:** US3-iam-core-entities-and-mappings.md
**Type:** Integration & Polish
**Priority:** High
**Estimated Effort:** 6-8 hours
**Prerequisites:** T11 (Frontend entity UIs must be complete)

## Objective

Integrate the IAM feature into the application navigation, add internationalization support, implement project members management, configure breadcrumbs, and perform final testing and polish.

## Acceptance Criteria

- [ ] Top-level IAM menu with 3 subpages in navigation
- [ ] Breadcrumb configuration for all IAM pages
- [ ] i18n translations for all IAM features (en-US)
- [ ] Project Settings → Members tab implemented
- [ ] Bidirectional member management (from user form and project settings)
- [ ] All cascade deletes verified working
- [ ] All CRUD + mapping operations tested end-to-end
- [ ] Responsive design verified on mobile/tablet/desktop
- [ ] No console errors or warnings
- [ ] All features working as specified in user story

## Technical Requirements

### 1. Navigation Menu Integration

#### frontend/app/composables/navigation/useNavigationItems.ts

Add IAM menu item with subpages:

```typescript
export function useNavigationItems() {
  const { t } = useI18n()

  return computed(() => [
    // ... existing menu items ...

    {
      title: t('navigation.iam'),
      icon: Shield,
      to: '/iam',
      children: [
        {
          title: t('navigation.iam.users'),
          to: '/iam/users',
          icon: Users,
        },
        {
          title: t('navigation.iam.roles'),
          to: '/iam/roles',
          icon: ShieldCheck,
        },
        {
          title: t('navigation.iam.permissions'),
          to: '/iam/permissions',
          icon: Key,
        },
      ],
    },
  ])
}
```

**Menu Structure:**
```
📁 IAM
  ├─ 👥 Users (/iam/users)
  ├─ ✓ Roles (/iam/roles)
  └─ 🔑 Permissions (/iam/permissions)
```

### 2. Breadcrumb Configuration

Add breadcrumb metadata to all IAM pages:

#### Users Page
```vue
definePageMeta({
  breadcrumb: {
    title: 'Users',
    items: [
      { title: 'IAM', to: '/iam' },
      { title: 'Users' },
    ],
  },
})
```

#### Roles Page
```vue
definePageMeta({
  breadcrumb: {
    title: 'Roles',
    items: [
      { title: 'IAM', to: '/iam' },
      { title: 'Roles' },
    ],
  },
})
```

#### Permissions Page
```vue
definePageMeta({
  breadcrumb: {
    title: 'Permissions',
    items: [
      { title: 'IAM', to: '/iam' },
      { title: 'Permissions' },
    ],
  },
})
```

### 3. Internationalization (i18n)

#### frontend/i18n/locales/en-US.json

Add comprehensive translations:

```json
{
  "navigation": {
    "iam": "Identity & Access",
    "iam.users": "Users",
    "iam.roles": "Roles",
    "iam.permissions": "Permissions"
  },

  "iam": {
    "users": {
      "title": "Users",
      "description": "Manage user accounts and access",
      "create": "Create User",
      "edit": "Edit User",
      "delete": "Delete User",
      "activate": "Activate User",
      "deactivate": "Deactivate User",
      "viewMappings": "View Mappings",

      "fields": {
        "email": "Email",
        "firstName": "First Name",
        "lastName": "Last Name",
        "isActive": "Active Status",
        "createdAt": "Created At"
      },

      "tabs": {
        "profile": "Profile",
        "roles": "Roles",
        "permissions": "Permissions",
        "projects": "Projects"
      },

      "messages": {
        "created": "User created successfully",
        "updated": "User updated successfully",
        "deleted": "User deleted successfully",
        "activated": "User activated successfully",
        "deactivated": "User deactivated successfully",
        "rolesAssigned": "Roles assigned successfully",
        "rolesRemoved": "Roles removed successfully",
        "permissionsAssigned": "Permissions assigned successfully",
        "permissionsRemoved": "Permissions removed successfully",
        "projectsAssigned": "Projects assigned successfully",
        "projectsRemoved": "Projects removed successfully"
      },

      "errors": {
        "createFailed": "Failed to create user",
        "updateFailed": "Failed to update user",
        "deleteFailed": "Failed to delete user",
        "emailExists": "Email already exists",
        "invalidEmail": "Invalid email format"
      }
    },

    "roles": {
      "title": "Roles",
      "description": "Manage roles and permissions",
      "create": "Create Role",
      "edit": "Edit Role",
      "delete": "Delete Role",

      "fields": {
        "name": "Name",
        "description": "Description",
        "createdAt": "Created At"
      },

      "tabs": {
        "details": "Details",
        "permissions": "Permissions"
      },

      "messages": {
        "created": "Role created successfully",
        "updated": "Role updated successfully",
        "deleted": "Role deleted successfully",
        "permissionsAssigned": "Permissions assigned successfully",
        "permissionsRemoved": "Permissions removed successfully"
      },

      "errors": {
        "createFailed": "Failed to create role",
        "updateFailed": "Failed to update role",
        "deleteFailed": "Failed to delete role",
        "nameExists": "Role name already exists",
        "roleInUse": "Cannot delete role that is assigned to users"
      }
    },

    "permissions": {
      "title": "Permissions",
      "description": "Manage system permissions",
      "create": "Create Permission",
      "createInline": "Create New Permission",
      "edit": "Edit Permission",
      "delete": "Delete Permission",

      "fields": {
        "name": "Name",
        "description": "Description",
        "effect": "Effect",
        "createdAt": "Created At"
      },

      "effects": {
        "allow": "Allow",
        "deny": "Deny"
      },

      "messages": {
        "created": "Permission created successfully",
        "updated": "Permission updated successfully",
        "deleted": "Permission deleted successfully",
        "createdAndAssigned": "Permission created and assigned successfully"
      },

      "errors": {
        "createFailed": "Failed to create permission",
        "updateFailed": "Failed to update permission",
        "deleteFailed": "Failed to delete permission",
        "nameExists": "Permission name already exists",
        "invalidName": "Permission name can only contain letters, numbers, underscores, and colons",
        "permissionInUse": "Cannot delete permission that is assigned to roles or users"
      }
    },

    "transferList": {
      "available": "Available {label}",
      "assigned": "Assigned {label}",
      "searchPlaceholder": "Search {label}...",
      "empty": "No {label} found",
      "noAssigned": "No {label} assigned"
    },

    "projectMembers": {
      "title": "Project Members",
      "description": "Manage project members and their roles",
      "addMember": "Add Member",
      "removeMember": "Remove Member",

      "roles": {
        "owner": "Owner",
        "admin": "Admin",
        "member": "Member",
        "viewer": "Viewer"
      },

      "messages": {
        "added": "Member added successfully",
        "removed": "Member removed successfully",
        "roleUpdated": "Member role updated successfully"
      },

      "errors": {
        "addFailed": "Failed to add member",
        "removeFailed": "Failed to remove member",
        "cannotRemoveLastOwner": "Cannot remove the last owner from the project"
      }
    }
  }
}
```

**Usage in Components:**
```vue
<script setup lang="ts">
const { t } = useI18n()
</script>

<template>
  <h1>{{ t('iam.users.title') }}</h1>
  <Button>{{ t('iam.users.create') }}</Button>
</template>
```

### 4. Project Members Integration

#### frontend/app/components/features/project/ProjectMembersTab.vue

Implement project members management UI (if doesn't exist):

```vue
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useProjectService } from '@/composables/services/useProjectService'
import { useIAMMapperService } from '@/composables/services/useIAMMapperService'
import { useUserService } from '@/composables/services/useUserService'
import TransferList from '@/components/ui/transfer-list/TransferList.vue'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const props = defineProps<{
  projectId: string
}>()

const { t } = useI18n()
const mapperService = useIAMMapperService()
const userService = useUserService()

const allUsers = ref<User[]>([])
const members = ref<ProjectMemberWithUser[]>([])
const isLoading = ref(true)

const assignedUsers = computed(() =>
  members.value.map(m => ({ ...m.user, role: m.role }))
)

const availableUsers = computed(() =>
  allUsers.value.filter(u => !members.value.some(m => m.user.id === u.id))
)

onMounted(async () => {
  await Promise.all([loadUsers(), loadMembers()])
  isLoading.value = false
})

async function loadUsers() {
  allUsers.value = await userService.query()
}

async function loadMembers() {
  const response = await mapperService.getProjectMembers(props.projectId)
  members.value = response.members
}

async function handleAddMembers(userIds: string[]) {
  const newMembers = userIds.map(userId => ({
    userId,
    role: 'member' as const, // Default role
  }))

  try {
    await mapperService.assignProjectMembers(props.projectId, newMembers)
    await loadMembers()
    toast({ title: t('iam.projectMembers.messages.added') })
  } catch (error) {
    toast({
      title: t('iam.projectMembers.errors.addFailed'),
      description: getConnectRPCError(error),
      variant: 'destructive',
    })
  }
}

async function handleRemoveMembers(userIds: string[]) {
  try {
    await mapperService.removeProjectMembers(props.projectId, userIds)
    await loadMembers()
    toast({ title: t('iam.projectMembers.messages.removed') })
  } catch (error) {
    if (hasConnectRPCError(error, 60803)) {
      // CannotRemoveLastOwner error
      toast({
        title: t('iam.projectMembers.errors.cannotRemoveLastOwner'),
        variant: 'destructive',
      })
    } else {
      toast({
        title: t('iam.projectMembers.errors.removeFailed'),
        description: getConnectRPCError(error),
        variant: 'destructive',
      })
    }
  }
}

async function handleRoleChange(userId: string, newRole: string) {
  // Remove then re-add with new role
  try {
    await mapperService.removeProjectMembers(props.projectId, [userId])
    await mapperService.assignProjectMembers(props.projectId, [
      { userId, role: newRole },
    ])
    await loadMembers()
    toast({ title: t('iam.projectMembers.messages.roleUpdated') })
  } catch (error) {
    toast({
      title: 'Failed to update role',
      description: getConnectRPCError(error),
      variant: 'destructive',
    })
  }
}
</script>

<template>
  <div v-if="!isLoading" class="space-y-6">
    <div>
      <h3 class="text-lg font-medium">{{ t('iam.projectMembers.title') }}</h3>
      <p class="text-sm text-muted-foreground">
        {{ t('iam.projectMembers.description') }}
      </p>
    </div>

    <!-- Members List with Role Dropdown -->
    <div class="space-y-2">
      <h4 class="text-sm font-medium">Current Members</h4>
      <div class="border rounded-lg">
        <div
          v-for="member in members"
          :key="member.user.id"
          class="flex items-center justify-between p-4 border-b last:border-b-0"
        >
          <div>
            <p class="font-medium">{{ member.user.email }}</p>
            <p class="text-sm text-muted-foreground">
              {{ member.user.firstName }} {{ member.user.lastName }}
            </p>
          </div>

          <div class="flex items-center gap-2">
            <Select
              :model-value="member.role"
              @update:model-value="handleRoleChange(member.user.id, $event)"
            >
              <SelectTrigger class="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="owner">
                  {{ t('iam.projectMembers.roles.owner') }}
                </SelectItem>
                <SelectItem value="admin">
                  {{ t('iam.projectMembers.roles.admin') }}
                </SelectItem>
                <SelectItem value="member">
                  {{ t('iam.projectMembers.roles.member') }}
                </SelectItem>
                <SelectItem value="viewer">
                  {{ t('iam.projectMembers.roles.viewer') }}
                </SelectItem>
              </SelectContent>
            </Select>

            <Button
              variant="ghost"
              size="icon"
              @click="handleRemoveMembers([member.user.id])"
            >
              <Trash2 class="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Members -->
    <div class="space-y-2">
      <h4 class="text-sm font-medium">Add Members</h4>
      <TransferList
        :available-items="availableUsers"
        :assigned-items="assignedUsers"
        label="Users"
        singular-label="User"
        label-key="email"
        @assign="handleAddMembers"
        @remove="handleRemoveMembers"
      />
    </div>
  </div>
</template>
```

**Integration in Project Settings:**
If `frontend/app/pages/projects/[id]/settings.vue` doesn't exist, create it with tabs for Settings and Members.

### 5. Final Testing Checklist

#### Database Tests
- [ ] Create user, role, permission via API
- [ ] Delete user → verify all mappings deleted (cascade)
- [ ] Delete role → verify all role mappings deleted (cascade)
- [ ] Delete permission → verify all permission mappings deleted (cascade)
- [ ] Verify unique constraints (email, role name, permission name)
- [ ] Verify CHECK constraints (permission name regex, effect enum, project role enum)

#### CRUD Operations
- [ ] Users: Create, Read, Update, Delete, Activate, Deactivate
- [ ] Roles: Create, Read, Update, Delete
- [ ] Permissions: Create, Read, Update, Delete

#### Mapping Operations
- [ ] Assign/remove user roles
- [ ] Assign/remove role permissions
- [ ] Assign/remove user permissions
- [ ] Assign/remove project members

#### UI Features
- [ ] TransferList search/filter working
- [ ] TransferList multi-select working
- [ ] Inline permission creation working
- [ ] Description tooltips showing for permissions
- [ ] Optimistic updates with rollback on error
- [ ] Toast notifications for all operations
- [ ] View mappings sheet working for users
- [ ] Tables pagination/sorting/filtering working

#### Error Handling
- [ ] Email uniqueness error shows correct message
- [ ] Role/permission name uniqueness error shows correct message
- [ ] Cannot delete role in use
- [ ] Cannot delete permission in use
- [ ] Cannot remove last project owner
- [ ] Invalid permission name shows validation error
- [ ] All ConnectRPC errors handled gracefully

#### Responsive Design
- [ ] Desktop (1920x1080) - all features accessible
- [ ] Tablet (768x1024) - tables scrollable, forms usable
- [ ] Mobile (375x667) - navigation working, forms accessible

#### Accessibility
- [ ] Keyboard navigation working in TransferList
- [ ] Focus states visible
- [ ] Form labels associated with inputs
- [ ] Error messages announced

### 6. Performance Optimization

If needed, add these optimizations:

**Lazy Load TransferList Data:**
```typescript
// Only load permissions when Permissions tab is activated
async function onTabChange(tabValue: string) {
  if (tabValue === 'permissions' && !permissionsLoaded.value) {
    await loadPermissions()
    permissionsLoaded.value = true
  }
}
```

**Debounce Search:**
```typescript
import { useDebounceFn } from '@vueuse/core'

const debouncedSearch = useDebounceFn((query: string) => {
  performSearch(query)
}, 300)
```

## Implementation Notes

### Bidirectional Member Management

**From User Form (Projects Tab):**
- User can see all projects they're a member of
- Can assign user to new projects
- Can remove user from projects

**From Project Settings (Members Tab):**
- Project admin can see all members
- Can add new members to project
- Can remove members from project
- Can change member roles

**Consistency**: Both UIs use the same backend API (IAMMapperService.assignProjectMembers/removeProjectMembers)

### Navigation Icons

Use lucide-vue-next icons:
- IAM: `Shield`
- Users: `Users`
- Roles: `ShieldCheck`
- Permissions: `Key`

### Cascade Delete Verification

Test cascade deletes work correctly:

```bash
# In psql:
-- Delete a user
DELETE FROM altalune_users WHERE id = 1;

-- Verify all related records deleted
SELECT * FROM altalune_users_roles WHERE user_id = 1;  -- Should be empty
SELECT * FROM altalune_users_permissions WHERE user_id = 1;  -- Should be empty
SELECT * FROM altalune_project_members WHERE user_id = 1;  -- Should be empty
```

## Files to Modify

- `frontend/app/composables/navigation/useNavigationItems.ts` - Add IAM menu
- `frontend/i18n/locales/en-US.json` - Add all i18n translations

## Files to Create

- `frontend/app/components/features/project/ProjectMembersTab.vue` (if doesn't exist)
- `frontend/app/pages/projects/[id]/settings.vue` (if doesn't exist)

## Commands to Run

```bash
# Frontend development
cd frontend && pnpm dev

# Lint and fix
cd frontend && pnpm lint:fix

# Type check
cd frontend && pnpm type-check

# Build for production (to verify)
cd frontend && pnpm build

# Backend testing
./bin/app serve -c config.yaml
```

## Definition of Done

- [ ] IAM menu appears in navigation with 3 subpages
- [ ] All breadcrumbs configured correctly
- [ ] All i18n translations added
- [ ] Project members tab working
- [ ] Bidirectional member management working
- [ ] All cascade deletes verified
- [ ] All CRUD operations tested
- [ ] All mapping operations tested
- [ ] Inline permission creation tested
- [ ] Description tooltips tested
- [ ] Responsive design verified
- [ ] No console errors
- [ ] No TypeScript errors
- [ ] All features match acceptance criteria in US3
- [ ] Production build successful

## Dependencies

- T11 (frontend entity UIs) must be complete
- All previous tasks (T7-T10) must be complete
- Navigation infrastructure
- i18n infrastructure
- Project settings page structure

## Risk Factors

- **Low risk**: Mostly integration and polish work
- **Watch out for**: Missing i18n keys causing blank labels
- **Test carefully**: Cascade deletes to ensure data integrity
- **Verify**: Project members bidirectional management doesn't cause race conditions
