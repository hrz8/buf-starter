# Task T64: Frontend Permission System

**Story Reference:** US15-authorization-rbac.md
**Type:** Frontend Implementation
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** T61 (Database + JWT Claims)

## Objective

Implement the frontend permission system including the `usePermissions` composable, auth store extension with permissions/memberships, and permission constants.

## Acceptance Criteria

- [x] Auth store extended with permissions and memberships from JWT
- [x] JWT decoded client-side for claims access
- [x] `usePermissions` composable implemented with helper functions
- [x] `can()`, `canAny()`, `canAll()` helper functions work correctly
- [x] `isMemberOf()`, `getProjectRole()` helper functions work correctly
- [x] `isSuperAdmin` computed property works correctly
- [x] Permission constants file created

## Technical Requirements

### Auth Repository Extension

Update `frontend/shared/repository/auth.ts` to include permissions and memberships:

```typescript
import type { $Fetch } from 'nitropack';

export interface AuthExchangeRequest {
  code: string;
  code_verifier: string;
  redirect_uri: string;
}

export interface AuthUserInfo {
  sub: string;
  email?: string;
  email_verified?: boolean;
  name?: string;
  given_name?: string;
  family_name?: string;
  picture?: string;
  perms?: string[];           // NEW: permissions array
  memberships?: Record<string, string>;  // NEW: project_public_id -> role
}

export interface AuthExchangeResponse {
  user: AuthUserInfo;
  expires_in: number;
}

// ... rest of repository unchanged
```

### Auth Store Extension

Update `frontend/app/stores/auth.ts`:

```typescript
import type { AuthUserInfo } from '~~/shared/repository/auth';

const RETURN_URL_KEY = 'oauth_return_url';
const ROOT_PERMISSION = 'root';

export const useAuthStore = defineStore('auth', () => {
  // User info from JWT (stored in memory, not the token itself)
  const user = ref<AuthUserInfo | null>(null);
  const expiresAt = ref<number | null>(null);

  const isAuthenticated = computed(() => {
    return !!user.value && (expiresAt.value ? Date.now() < expiresAt.value : true);
  });

  // Email verification status computed properties
  const isEmailVerified = computed(() => {
    return user.value?.email_verified ?? false;
  });

  const isEmailVerificationRequired = computed(() => {
    return isAuthenticated.value && !isEmailVerified.value;
  });

  // ============================================
  // NEW: Permission-related computed properties
  // ============================================

  /**
   * User's permissions array from JWT
   */
  const permissions = computed(() => {
    return user.value?.perms ?? [];
  });

  /**
   * User's project memberships from JWT
   * Format: { "proj_abc123": "admin", "proj_xyz789": "member" }
   */
  const memberships = computed(() => {
    return user.value?.memberships ?? {};
  });

  /**
   * Check if user is superadmin (has 'root' permission)
   */
  const isSuperAdmin = computed(() => {
    return permissions.value.includes(ROOT_PERMISSION);
  });

  /**
   * Get list of project IDs user is member of
   */
  const memberProjectIds = computed(() => {
    return Object.keys(memberships.value);
  });

  // ============================================
  // Existing methods
  // ============================================

  function setUser(userData: AuthUserInfo, expiresIn: number) {
    user.value = userData;
    expiresAt.value = Date.now() + expiresIn * 1000;
  }

  function clearAuth() {
    user.value = null;
    expiresAt.value = null;
  }

  // Return URL is stored in sessionStorage to persist across OAuth redirects
  function setReturnUrl(url: string | null) {
    if (import.meta.client) {
      if (url) {
        sessionStorage.setItem(RETURN_URL_KEY, url);
      }
      else {
        sessionStorage.removeItem(RETURN_URL_KEY);
      }
    }
  }

  function getAndClearReturnUrl(): string | null {
    if (import.meta.client) {
      const url = sessionStorage.getItem(RETURN_URL_KEY);
      sessionStorage.removeItem(RETURN_URL_KEY);
      return url;
    }
    return null;
  }

  return {
    user: readonly(user),
    expiresAt: readonly(expiresAt),
    isAuthenticated,
    isEmailVerified,
    isEmailVerificationRequired,
    // NEW exports
    permissions,
    memberships,
    isSuperAdmin,
    memberProjectIds,
    // Existing methods
    setUser,
    clearAuth,
    setReturnUrl,
    getAndClearReturnUrl,
  };
});
```

### Permission Composable

Create `frontend/app/composables/usePermissions.ts`:

```typescript
/**
 * Permission checking composable
 * Provides reactive permission helpers that respect superadmin status
 */
export function usePermissions() {
  const authStore = useAuthStore();

  /**
   * User's permissions array
   */
  const permissions = computed(() => authStore.permissions);

  /**
   * User's project memberships
   */
  const memberships = computed(() => authStore.memberships);

  /**
   * Check if user has a specific permission
   * @param permission - Permission to check (e.g., 'employee:read')
   */
  const can = (permission: string): boolean => {
    // Superadmin has all permissions
    if (authStore.isSuperAdmin) return true;
    return permissions.value.includes(permission);
  };

  /**
   * Check if user has ANY of the specified permissions
   * @param perms - Array of permissions to check
   */
  const canAny = (perms: string[]): boolean => {
    if (authStore.isSuperAdmin) return true;
    return perms.some(p => permissions.value.includes(p));
  };

  /**
   * Check if user has ALL of the specified permissions
   * @param perms - Array of permissions to check
   */
  const canAll = (perms: string[]): boolean => {
    if (authStore.isSuperAdmin) return true;
    return perms.every(p => permissions.value.includes(p));
  };

  /**
   * Check if user is a member of a specific project
   * @param projectId - Project public ID
   */
  const isMemberOf = (projectId: string): boolean => {
    if (authStore.isSuperAdmin) return true;
    return projectId in memberships.value;
  };

  /**
   * Get user's role in a specific project
   * @param projectId - Project public ID
   * @returns Role string or null if not a member
   */
  const getProjectRole = (projectId: string): string | null => {
    return memberships.value[projectId] ?? null;
  };

  /**
   * Check if user has permission AND is member of project
   * This mirrors the backend CheckProjectAccess logic
   * @param permission - Permission to check
   * @param projectId - Project public ID
   */
  const canAccessProject = (permission: string, projectId: string): boolean => {
    if (authStore.isSuperAdmin) return true;
    return can(permission) && isMemberOf(projectId);
  };

  /**
   * Check if user is superadmin
   */
  const isSuperAdmin = computed(() => authStore.isSuperAdmin);

  /**
   * Get list of project IDs user is member of
   */
  const memberProjects = computed(() => authStore.memberProjectIds);

  return {
    // Reactive state
    permissions,
    memberships,
    isSuperAdmin,
    memberProjects,
    // Helper functions
    can,
    canAny,
    canAll,
    isMemberOf,
    getProjectRole,
    canAccessProject,
  };
}
```

### Permission Constants

Create `frontend/app/constants/permissions.ts`:

```typescript
/**
 * Permission constants for type-safe permission checks
 * Matches backend predefined permissions
 */
export const PERMISSIONS = {
  // Employee entity
  EMPLOYEE: {
    READ: 'employee:read',
    WRITE: 'employee:write',
    DELETE: 'employee:delete',
  },
  // User entity (IAM)
  USER: {
    READ: 'user:read',
    WRITE: 'user:write',
    DELETE: 'user:delete',
  },
  // Role entity (IAM)
  ROLE: {
    READ: 'role:read',
    WRITE: 'role:write',
    DELETE: 'role:delete',
  },
  // Permission entity (IAM)
  PERMISSION: {
    READ: 'permission:read',
    WRITE: 'permission:write',
    DELETE: 'permission:delete',
  },
  // Project entity
  PROJECT: {
    READ: 'project:read',
    WRITE: 'project:write',
    DELETE: 'project:delete',
  },
  // API Key entity
  API_KEY: {
    READ: 'apikey:read',
    WRITE: 'apikey:write',
    DELETE: 'apikey:delete',
  },
  // Chatbot Config entity
  CHATBOT: {
    READ: 'chatbot:read',
    WRITE: 'chatbot:write',
    DELETE: 'chatbot:delete',
  },
  // OAuth Client entity
  CLIENT: {
    READ: 'client:read',
    WRITE: 'client:write',
    DELETE: 'client:delete',
  },
  // Project Members management
  MEMBER: {
    READ: 'member:read',
    WRITE: 'member:write',
    DELETE: 'member:delete',
  },
  // IAM Mapper (role-permission assignments)
  IAM: {
    READ: 'iam:read',
    WRITE: 'iam:write',
  },
  // Special permissions
  ROOT: 'root',
} as const;

/**
 * Type for all permission values
 */
export type Permission =
  | typeof PERMISSIONS.EMPLOYEE[keyof typeof PERMISSIONS.EMPLOYEE]
  | typeof PERMISSIONS.USER[keyof typeof PERMISSIONS.USER]
  | typeof PERMISSIONS.ROLE[keyof typeof PERMISSIONS.ROLE]
  | typeof PERMISSIONS.PERMISSION[keyof typeof PERMISSIONS.PERMISSION]
  | typeof PERMISSIONS.PROJECT[keyof typeof PERMISSIONS.PROJECT]
  | typeof PERMISSIONS.API_KEY[keyof typeof PERMISSIONS.API_KEY]
  | typeof PERMISSIONS.CHATBOT[keyof typeof PERMISSIONS.CHATBOT]
  | typeof PERMISSIONS.CLIENT[keyof typeof PERMISSIONS.CLIENT]
  | typeof PERMISSIONS.MEMBER[keyof typeof PERMISSIONS.MEMBER]
  | typeof PERMISSIONS.IAM[keyof typeof PERMISSIONS.IAM]
  | typeof PERMISSIONS.ROOT;

/**
 * Project membership roles
 */
export const PROJECT_ROLES = {
  OWNER: 'owner',
  ADMIN: 'admin',
  MEMBER: 'member',
  USER: 'user',
} as const;

export type ProjectRole = typeof PROJECT_ROLES[keyof typeof PROJECT_ROLES];

/**
 * Permission descriptions for UI display
 * Can be used in permission management interfaces
 */
export const PERMISSION_DESCRIPTIONS: Record<string, string> = {
  'employee:read': 'View employee records',
  'employee:write': 'Create and update employee records',
  'employee:delete': 'Delete employee records',
  'user:read': 'View user accounts',
  'user:write': 'Create and update user accounts',
  'user:delete': 'Delete user accounts',
  'role:read': 'View roles',
  'role:write': 'Create and update roles',
  'role:delete': 'Delete roles',
  'permission:read': 'View permissions',
  'permission:write': 'Create and update permissions',
  'permission:delete': 'Delete permissions',
  'project:read': 'View projects',
  'project:write': 'Create and update projects',
  'project:delete': 'Delete projects',
  'apikey:read': 'View API keys',
  'apikey:write': 'Create and update API keys',
  'apikey:delete': 'Delete API keys',
  'chatbot:read': 'View chatbot configurations',
  'chatbot:write': 'Create and update chatbot configurations',
  'chatbot:delete': 'Delete chatbot configurations',
  'client:read': 'View OAuth clients',
  'client:write': 'Create and update OAuth clients',
  'client:delete': 'Delete OAuth clients',
  'member:read': 'View project members',
  'member:write': 'Add and update project members',
  'member:delete': 'Remove project members',
  'iam:read': 'View IAM mappings',
  'iam:write': 'Manage IAM mappings',
  'root': 'Full superadmin access',
};
```

### Backend BFF Update (if needed)

Ensure the BFF `/me` endpoint returns permissions and memberships. Update `internal/domain/bff/handler.go` if needed:

```go
// MeResponse in the handler should include perms and memberships
type MeResponse struct {
    User struct {
        Sub           string            `json:"sub"`
        Email         string            `json:"email,omitempty"`
        EmailVerified bool              `json:"email_verified"`
        Name          string            `json:"name,omitempty"`
        Perms         []string          `json:"perms,omitempty"`
        Memberships   map[string]string `json:"memberships,omitempty"`
    } `json:"user"`
    ExpiresIn int64 `json:"expires_in"`
}
```

## Files to Create

- `frontend/app/composables/usePermissions.ts` - Permission checking composable
- `frontend/app/constants/permissions.ts` - Permission constants

## Files to Modify

- `frontend/shared/repository/auth.ts` - Add perms and memberships to AuthUserInfo
- `frontend/app/stores/auth.ts` - Add permissions, memberships, isSuperAdmin computed

## Testing Requirements

```typescript
// Test usePermissions composable
describe('usePermissions', () => {
  it('should allow superadmin all permissions', () => {
    // Mock auth store with root permission
    const { can, canAny, canAll, isMemberOf } = usePermissions();

    expect(can('employee:read')).toBe(true);
    expect(can('nonexistent:permission')).toBe(true);
    expect(isMemberOf('any-project')).toBe(true);
  });

  it('should check specific permission', () => {
    // Mock auth store with limited permissions
    const { can } = usePermissions();

    expect(can('employee:read')).toBe(true);
    expect(can('employee:write')).toBe(false);
  });

  it('should check project membership', () => {
    // Mock auth store with memberships
    const { isMemberOf, getProjectRole } = usePermissions();

    expect(isMemberOf('proj_abc')).toBe(true);
    expect(isMemberOf('proj_xyz')).toBe(false);
    expect(getProjectRole('proj_abc')).toBe('admin');
    expect(getProjectRole('proj_xyz')).toBeNull();
  });
});
```

## Commands to Run

```bash
# Run frontend dev server
cd frontend && pnpm dev

# Run linting
cd frontend && pnpm lint

# Run type checking
cd frontend && pnpm typecheck
```

## Validation Checklist

- [ ] AuthUserInfo type includes perms and memberships
- [ ] Auth store exposes permissions computed
- [ ] Auth store exposes memberships computed
- [ ] Auth store exposes isSuperAdmin computed
- [ ] usePermissions composable created
- [ ] can() checks single permission
- [ ] canAny() checks any of permissions
- [ ] canAll() checks all permissions
- [ ] isMemberOf() checks project membership
- [ ] getProjectRole() returns role or null
- [ ] canAccessProject() combines permission + membership check
- [ ] Superadmin bypasses all checks
- [ ] Permission constants match backend
- [ ] TypeScript types are correct

## Definition of Done

- [x] Auth store extended with permissions
- [x] usePermissions composable implemented
- [x] Permission constants file created
- [x] All helper functions work correctly
- [x] Superadmin bypass works
- [x] Linting passes
- [x] Type checking passes

## Dependencies

- T61: JWT must include perms and memberships claims

## Risk Factors

- **Low Risk**: Simple reactive state management
- **Low Risk**: Client-side JWT decoding is standard pattern
- **Low Risk**: Matches existing auth store patterns

## Notes

- Frontend permission checks are for UX only (hiding buttons, nav items)
- Backend always re-validates permissions
- JWT decoding doesn't verify signature (server's job)
- Permissions refresh on token refresh
- Consider memoization if performance becomes an issue

### Usage Examples

```vue
<script setup>
import { PERMISSIONS } from '~/constants/permissions';

const { can, canAny, isMemberOf, isSuperAdmin } = usePermissions();

// Check single permission
const canViewEmployees = computed(() => can(PERMISSIONS.EMPLOYEE.READ));

// Check multiple permissions
const canManageIAM = computed(() =>
  canAny([PERMISSIONS.USER.READ, PERMISSIONS.ROLE.READ, PERMISSIONS.PERMISSION.READ])
);

// Check project access
const projectStore = useProjectStore();
const canEditInProject = computed(() =>
  can(PERMISSIONS.EMPLOYEE.WRITE) && isMemberOf(projectStore.currentProjectId)
);
</script>

<template>
  <Button v-if="canViewEmployees" @click="viewEmployees">
    View Employees
  </Button>

  <span v-if="isSuperAdmin" class="badge">Superadmin</span>
</template>
```
