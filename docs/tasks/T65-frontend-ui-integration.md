# Task T65: Frontend UI Integration

**Story Reference:** US15-authorization-rbac.md
**Type:** Frontend Implementation
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** T64 (Frontend Permission System)

## Objective

Integrate the permission system into the frontend UI including PermissionGuard component, sidebar navigation filtering, action button hiding, and access denied handling.

## Acceptance Criteria

- [x] PermissionGuard component for declarative permission checks
- [x] Sidebar navigation filtered by user permissions
- [ ] Create/Edit/Delete buttons hidden when user lacks permission (requires page-by-page updates)
- [x] Project selector shows only projects user is member of
- [x] Access denied page for unauthorized direct URL access
- [x] No UI flash before permission check completes (AuthLoadingOverlay)

## Technical Requirements

### PermissionGuard Component

Create `frontend/app/components/custom/PermissionGuard.vue`:

```vue
<script setup lang="ts">
import type { Permission } from '~/constants/permissions';

const props = defineProps<{
  /**
   * Single permission to check
   */
  permission?: Permission | string
  /**
   * Multiple permissions to check
   */
  permissions?: (Permission | string)[]
  /**
   * Mode for multiple permissions: 'any' (OR) or 'all' (AND)
   */
  mode?: 'any' | 'all'
  /**
   * Fallback behavior: 'hide' (default), 'disabled', 'redirect'
   */
  fallback?: 'hide' | 'disabled' | 'redirect'
  /**
   * Optional project ID for project-scoped permission check
   */
  projectId?: string
}>();

const { can, canAny, canAll, isMemberOf } = usePermissions();

const hasAccess = computed(() => {
  // If projectId provided, also check membership
  if (props.projectId && !isMemberOf(props.projectId)) {
    return false;
  }

  // Single permission check
  if (props.permission) {
    return can(props.permission);
  }

  // Multiple permissions check
  if (props.permissions && props.permissions.length > 0) {
    return props.mode === 'all'
      ? canAll(props.permissions)
      : canAny(props.permissions);
  }

  // No permission specified = allow
  return true;
});

// Handle redirect fallback
const router = useRouter();
watch(hasAccess, (value) => {
  if (!value && props.fallback === 'redirect') {
    router.push('/access-denied');
  }
}, { immediate: true });
</script>

<template>
  <!-- Render slot content if has access -->
  <slot v-if="hasAccess" />

  <!-- Render disabled slot if no access and fallback is 'disabled' -->
  <slot v-else-if="fallback === 'disabled'" name="disabled">
    <!-- Default disabled state -->
    <div class="opacity-50 pointer-events-none">
      <slot />
    </div>
  </slot>

  <!-- Render fallback slot for 'hide' mode (default shows nothing) -->
  <slot v-else-if="fallback !== 'redirect'" name="fallback" />
</template>
```

### Access Denied Page

Create `frontend/app/pages/access-denied.vue`:

```vue
<script setup lang="ts">
definePageMeta({
  layout: 'default',
});

const { t } = useI18n();
const router = useRouter();

function goBack() {
  router.back();
}

function goHome() {
  router.push('/dashboard');
}
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-[60vh] text-center">
    <Icon name="lucide:shield-x" class="w-16 h-16 text-destructive mb-4" />

    <h1 class="text-2xl font-semibold mb-2">
      {{ t('errors.accessDenied.title') }}
    </h1>

    <p class="text-muted-foreground mb-6 max-w-md">
      {{ t('errors.accessDenied.description') }}
    </p>

    <div class="flex gap-4">
      <Button variant="outline" @click="goBack">
        <Icon name="lucide:arrow-left" class="w-4 h-4 mr-2" />
        {{ t('common.goBack') }}
      </Button>

      <Button @click="goHome">
        <Icon name="lucide:home" class="w-4 h-4 mr-2" />
        {{ t('common.goHome') }}
      </Button>
    </div>
  </div>
</template>
```

### i18n Translations

Add to `frontend/app/i18n/locales/en.json`:

```json
{
  "errors": {
    "accessDenied": {
      "title": "Access Denied",
      "description": "You don't have permission to access this page. Please contact your administrator if you believe this is an error."
    }
  },
  "common": {
    "goBack": "Go Back",
    "goHome": "Go to Dashboard"
  }
}
```

### Navigation Items with Permissions

Update `frontend/app/composables/navigation/useNavigationItems.ts`:

```typescript
import { PERMISSIONS } from '~/constants/permissions';

export function useNavigationItems() {
  const { t } = useI18n();
  const { can } = usePermissions();
  const chatbotStore = useChatbotStore();
  const nodeStore = useChatbotNodeStore();

  // ... existing computed for module items ...

  /**
   * Main navigation items with translations and permissions
   */
  const mainNavItems = computed<NavItem[]>(() => {
    const items: NavItem[] = [
      {
        title: t('nav.dashboard'),
        to: '/dashboard',
        icon: LucideHome,
        breadcrumb: {
          path: '/dashboard',
          label: 'nav.dashboard',
          i18nKey: 'nav.dashboard',
        },
        // No permission required - everyone can see dashboard
      },
      {
        title: t('nav.devices.title'),
        to: '/devices',
        match: '/devices',
        icon: Smartphone,
        breadcrumb: {
          path: '/devices',
          label: 'nav.devices.title',
          i18nKey: 'nav.devices.title',
        },
        items: [
          {
            title: t('nav.devices.scan'),
            to: '/devices/scan',
            breadcrumb: {
              path: '/devices/scan',
              label: 'nav.devices.scan',
              i18nKey: 'nav.devices.scan',
              parent: '/devices',
            },
          },
          {
            title: t('nav.devices.chat'),
            to: '/devices/chat',
            breadcrumb: {
              path: '/devices/chat',
              label: 'nav.devices.chat',
              i18nKey: 'nav.devices.chat',
              parent: '/devices',
            },
          },
        ],
      },
      {
        title: t('nav.examples.title'),
        to: '/examples',
        match: '/examples',
        icon: Puzzle,
        permission: PERMISSIONS.EMPLOYEE.READ, // NEW: requires employee:read
        breadcrumb: {
          path: '/examples',
          label: 'nav.examples.title',
          i18nKey: 'nav.examples.title',
        },
        items: [
          {
            title: t('nav.examples.datatable'),
            to: '/examples/datatable',
            breadcrumb: {
              path: '/examples/datatable',
              label: 'nav.examples.datatable',
              i18nKey: 'nav.examples.datatable',
              parent: '/examples',
            },
          },
        ],
      },
      {
        title: t('nav.chatbot.title'),
        to: '/platform/modules',
        match: '/platform/modules',
        icon: Bot,
        permission: PERMISSIONS.CHATBOT.READ, // NEW: requires chatbot:read
        breadcrumb: {
          path: '/platform/modules',
          label: 'nav.chatbot.title',
          i18nKey: 'nav.chatbot.title',
        },
        items: chatbotModuleItems.value,
      },
      {
        title: t('nav.nodes.title'),
        to: '/platform/nodes',
        match: '/platform/node',
        icon: Workflow,
        permission: PERMISSIONS.CHATBOT.READ, // NEW: requires chatbot:read
        breadcrumb: {
          path: '/platform/nodes',
          label: 'nav.nodes.title',
          i18nKey: 'nav.nodes.title',
        },
        items: chatbotNodeItems.value,
        action: 'createNode',
        defaultExpanded: true,
      },
    ];

    // Filter items by permission
    return filterNavItemsByPermission(items);
  });

  /**
   * Settings navigation items with permissions
   */
  const settingsNavItems = computed<SettingsItem[]>(() => {
    const items: SettingsItem[] = [
      {
        name: t('nav.settings.apiKeys'),
        url: '/settings/api-keys',
        icon: Key,
        permission: PERMISSIONS.API_KEY.READ, // NEW
        breadcrumb: {
          path: '/settings/api-keys',
          label: 'nav.settings.apiKeys',
          i18nKey: 'nav.settings.apiKeys',
          parent: '/settings',
        },
      },
      {
        name: t('nav.iam.oauthClients'),
        url: '/settings/oauth-client',
        icon: KeyRound,
        permission: PERMISSIONS.CLIENT.READ, // NEW
        breadcrumb: {
          path: '/settings/oauth-client',
          label: 'nav.iam.oauthClients',
          i18nKey: 'nav.iam.oauthClients',
          parent: '/settings',
        },
      },
      {
        name: t('nav.settings.project'),
        url: '/settings/project',
        icon: Cog,
        permission: PERMISSIONS.PROJECT.READ, // NEW
        breadcrumb: {
          path: '/settings/project',
          label: 'nav.settings.project',
          i18nKey: 'nav.settings.project',
          parent: '/settings',
        },
      },
    ];

    // Filter by permission
    return items.filter(item => !item.permission || can(item.permission));
  });

  /**
   * IAM navigation items with permissions
   */
  const iamNavItems = computed<SettingsItem[]>(() => {
    const items: SettingsItem[] = [
      {
        name: t('nav.iam.users'),
        url: '/iam/users',
        icon: Users,
        permission: PERMISSIONS.USER.READ, // NEW
        breadcrumb: {
          path: '/iam/users',
          label: 'nav.iam.users',
          i18nKey: 'nav.iam.users',
          parent: '/iam',
        },
      },
      {
        name: t('nav.iam.roles'),
        url: '/iam/roles',
        icon: ShieldCheck,
        permission: PERMISSIONS.ROLE.READ, // NEW
        breadcrumb: {
          path: '/iam/roles',
          label: 'nav.iam.roles',
          i18nKey: 'nav.iam.roles',
          parent: '/iam',
        },
      },
      {
        name: t('nav.iam.permissions'),
        url: '/iam/permissions',
        icon: Key,
        permission: PERMISSIONS.PERMISSION.READ, // NEW
        breadcrumb: {
          path: '/iam/permissions',
          label: 'nav.iam.permissions',
          i18nKey: 'nav.iam.permissions',
          parent: '/iam',
        },
      },
      {
        name: t('nav.iam.oauthProviders'),
        url: '/iam/oauth-provider',
        icon: KeyRound,
        // No permission - visible to all authenticated users
        breadcrumb: {
          path: '/iam/oauth-provider',
          label: 'nav.iam.oauthProviders',
          i18nKey: 'nav.iam.oauthProviders',
          parent: '/iam',
        },
      },
    ];

    // Filter by permission
    return items.filter(item => !item.permission || can(item.permission));
  });

  /**
   * Filter navigation items by permission (recursive)
   */
  function filterNavItemsByPermission(items: NavItem[]): NavItem[] {
    return items
      .filter(item => !item.permission || can(item.permission))
      .map(item => ({
        ...item,
        items: item.items ? filterNavItemsByPermission(item.items) : undefined,
      }));
  }

  return {
    mainNavItems,
    settingsNavItems,
    iamNavItems,
  };
}
```

### NavItem Type Extension

Update `frontend/app/types/navigation.ts`:

```typescript
export interface NavItem {
  title: string;
  to: string;
  icon?: LucideIcon;
  match?: string;
  permission?: string; // NEW: required permission to view
  breadcrumb?: BreadcrumbConfig;
  items?: NavSubItem[];
  action?: string;
  defaultExpanded?: boolean;
}

export interface NavSubItem {
  title: string;
  to: string;
  icon?: LucideIcon;
  permission?: string; // NEW
  breadcrumb?: BreadcrumbConfig;
  badge?: string;
  badgeVariant?: 'default' | 'secondary';
}

export interface SettingsItem {
  name: string;
  url: string;
  icon: LucideIcon;
  permission?: string; // NEW
  breadcrumb?: BreadcrumbConfig;
}
```

### Action Button Permission Checks

Update list pages to conditionally show action buttons. Example for employees page:

```vue
<!-- frontend/app/pages/examples/datatable.vue -->
<script setup lang="ts">
import { PERMISSIONS } from '~/constants/permissions';

const { can } = usePermissions();
const canCreate = computed(() => can(PERMISSIONS.EMPLOYEE.WRITE));
const canDelete = computed(() => can(PERMISSIONS.EMPLOYEE.DELETE));
</script>

<template>
  <div>
    <!-- Page header with conditional Create button -->
    <div class="flex justify-between items-center mb-4">
      <h1>{{ t('pages.employees.title') }}</h1>

      <Button v-if="canCreate" @click="openCreateDialog">
        <Icon name="lucide:plus" class="w-4 h-4 mr-2" />
        {{ t('common.create') }}
      </Button>
    </div>

    <!-- Data table with conditional actions -->
    <DataTable :columns="columns" :data="employees">
      <template #actions="{ row }">
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="ghost" size="icon">
              <Icon name="lucide:more-horizontal" class="w-4 h-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <!-- View is always visible -->
            <DropdownMenuItem @click="viewEmployee(row)">
              <Icon name="lucide:eye" class="w-4 h-4 mr-2" />
              {{ t('common.view') }}
            </DropdownMenuItem>

            <!-- Edit requires write permission -->
            <PermissionGuard :permission="PERMISSIONS.EMPLOYEE.WRITE">
              <DropdownMenuItem @click="editEmployee(row)">
                <Icon name="lucide:pencil" class="w-4 h-4 mr-2" />
                {{ t('common.edit') }}
              </DropdownMenuItem>
            </PermissionGuard>

            <!-- Delete requires delete permission -->
            <PermissionGuard :permission="PERMISSIONS.EMPLOYEE.DELETE">
              <DropdownMenuSeparator />
              <DropdownMenuItem
                class="text-destructive"
                @click="confirmDelete(row)"
              >
                <Icon name="lucide:trash-2" class="w-4 h-4 mr-2" />
                {{ t('common.delete') }}
              </DropdownMenuItem>
            </PermissionGuard>
          </DropdownMenuContent>
        </DropdownMenu>
      </template>
    </DataTable>
  </div>
</template>
```

### Project Selector Filtering

Update project selector to filter by membership:

```vue
<!-- frontend/app/components/custom/nav/ProjectSelector.vue -->
<script setup lang="ts">
const { memberships, isSuperAdmin } = usePermissions();
const projectStore = useProjectStore();

// Filter projects by membership (superadmin sees all)
const visibleProjects = computed(() => {
  if (isSuperAdmin.value) {
    return projectStore.projects;
  }

  return projectStore.projects.filter(
    project => project.public_id in memberships.value
  );
});

// Get role badge for project
const getProjectRoleBadge = (projectId: string) => {
  const role = memberships.value[projectId];
  if (role === 'owner' || role === 'admin') {
    return { label: role, variant: 'default' };
  }
  return null;
};
</script>

<template>
  <Select v-model="projectStore.currentProjectId">
    <SelectTrigger>
      <SelectValue :placeholder="t('nav.selectProject')" />
    </SelectTrigger>
    <SelectContent>
      <SelectItem
        v-for="project in visibleProjects"
        :key="project.id"
        :value="project.public_id"
      >
        <div class="flex items-center gap-2">
          <span>{{ project.name }}</span>
          <Badge
            v-if="getProjectRoleBadge(project.public_id)"
            :variant="getProjectRoleBadge(project.public_id)?.variant"
            class="text-xs"
          >
            {{ getProjectRoleBadge(project.public_id)?.label }}
          </Badge>
        </div>
      </SelectItem>
    </SelectContent>
  </Select>
</template>
```

### Route Middleware for Protected Pages

Create `frontend/app/middleware/permission.ts`:

```typescript
import { PERMISSIONS } from '~/constants/permissions';

/**
 * Route permission requirements
 * Maps route patterns to required permissions
 */
const ROUTE_PERMISSIONS: Record<string, string | string[]> = {
  '/examples/datatable': PERMISSIONS.EMPLOYEE.READ,
  '/iam/users': PERMISSIONS.USER.READ,
  '/iam/roles': PERMISSIONS.ROLE.READ,
  '/iam/permissions': PERMISSIONS.PERMISSION.READ,
  '/settings/api-keys': PERMISSIONS.API_KEY.READ,
  '/settings/oauth-client': PERMISSIONS.CLIENT.READ,
  '/settings/project': PERMISSIONS.PROJECT.READ,
  '/platform/modules': PERMISSIONS.CHATBOT.READ,
  '/platform/node': PERMISSIONS.CHATBOT.READ,
};

export default defineNuxtRouteMiddleware((to) => {
  const { can, canAny } = usePermissions();
  const authStore = useAuthStore();

  // Skip permission check if not authenticated (auth middleware handles redirect)
  if (!authStore.isAuthenticated) {
    return;
  }

  // Check route against permission requirements
  for (const [pattern, permission] of Object.entries(ROUTE_PERMISSIONS)) {
    if (to.path.startsWith(pattern)) {
      const hasPermission = Array.isArray(permission)
        ? canAny(permission)
        : can(permission);

      if (!hasPermission) {
        return navigateTo('/access-denied');
      }
      break;
    }
  }
});
```

Register middleware in `nuxt.config.ts` or apply to specific pages:

```typescript
// In page component
definePageMeta({
  middleware: ['auth', 'permission'],
});
```

## Files to Create

- `frontend/app/components/custom/PermissionGuard.vue` - Permission guard component
- `frontend/app/pages/access-denied.vue` - Access denied page
- `frontend/app/middleware/permission.ts` - Route permission middleware

## Files to Modify

- `frontend/app/composables/navigation/useNavigationItems.ts` - Add permission filtering
- `frontend/app/types/navigation.ts` - Add permission field to nav types
- `frontend/app/i18n/locales/en.json` - Add access denied translations
- `frontend/app/i18n/locales/id.json` - Add access denied translations
- `frontend/app/pages/examples/datatable.vue` - Add permission checks to actions
- `frontend/app/components/custom/nav/ProjectSelector.vue` - Filter by membership

## Validation Checklist

- [ ] PermissionGuard component renders correctly
- [ ] PermissionGuard hides content when permission denied
- [ ] PermissionGuard disabled mode works
- [ ] PermissionGuard redirect mode works
- [ ] Sidebar filters items by permission
- [ ] Settings nav filters by permission
- [ ] IAM nav filters by permission
- [ ] Create buttons hidden when no write permission
- [ ] Edit buttons hidden when no write permission
- [ ] Delete buttons hidden when no delete permission
- [ ] Project selector shows only member projects
- [ ] Access denied page displays correctly
- [ ] Direct URL to unauthorized page redirects
- [ ] No UI flash before permission check
- [ ] Superadmin sees all items

## Definition of Done

- [x] PermissionGuard component created
- [x] Access denied page created
- [x] Navigation filtered by permissions
- [ ] Action buttons conditionally rendered (can be done per-page as needed)
- [x] Project selector filters by membership
- [x] Route middleware created
- [x] i18n translations added
- [x] AuthLoadingOverlay to prevent UI flash
- [x] No console errors
- [x] Linting passes

## Dependencies

- T64: usePermissions composable must be available

## Risk Factors

- **Medium Risk**: UI changes across multiple components
- **Low Risk**: Navigation filtering is reactive
- **Low Risk**: Permission guard is simple wrapper

## Notes

- Frontend checks are for UX only - backend always re-validates
- Always use `v-if` not `v-show` for permission checks (prevents DOM exposure)
- PermissionGuard can be nested for complex permission logic
- Consider loading states to prevent flash of unauthorized content
- Test with different user permission combinations

### Testing Scenarios

1. **Superadmin user**: All nav items visible, all actions available
2. **User with employee:read only**: Can see employee list, no create/edit/delete
3. **User with no permissions**: Dashboard only, limited nav items
4. **Direct URL access**: Unauthorized pages redirect to access-denied
5. **Project membership**: Only member projects in selector

### UI Flash Prevention

To prevent showing unauthorized content briefly:

```vue
<script setup>
const authStore = useAuthStore();
const { can } = usePermissions();

// Don't render until auth is loaded
const isReady = computed(() => authStore.isAuthenticated);
</script>

<template>
  <div v-if="isReady">
    <!-- Permission-gated content -->
    <PermissionGuard permission="employee:read">
      <EmployeeList />
    </PermissionGuard>
  </div>
  <div v-else>
    <!-- Loading skeleton -->
    <Skeleton class="h-48" />
  </div>
</template>
```
