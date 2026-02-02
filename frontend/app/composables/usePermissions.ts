import { useAuthStore } from '~/stores/auth';

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
  function can(permission: string): boolean {
    // Superadmin has all permissions
    if (authStore.isSuperAdmin)
      return true;
    return permissions.value.includes(permission);
  }

  /**
   * Check if user has ANY of the specified permissions
   * @param perms - Array of permissions to check
   */
  function canAny(perms: string[]): boolean {
    if (authStore.isSuperAdmin)
      return true;
    return perms.some(p => permissions.value.includes(p));
  }

  /**
   * Check if user has ALL of the specified permissions
   * @param perms - Array of permissions to check
   */
  function canAll(perms: string[]): boolean {
    if (authStore.isSuperAdmin)
      return true;
    return perms.every(p => permissions.value.includes(p));
  }

  /**
   * Check if user is a member of a specific project
   * @param projectId - Project public ID
   */
  function isMemberOf(projectId: string): boolean {
    if (authStore.isSuperAdmin)
      return true;
    return projectId in memberships.value;
  }

  /**
   * Get user's role in a specific project
   * @param projectId - Project public ID
   * @returns Role string or null if not a member
   */
  function getProjectRole(projectId: string): string | null {
    return memberships.value[projectId] ?? null;
  }

  /**
   * Check if user has permission AND is member of project
   * This mirrors the backend CheckProjectAccess logic
   * @param permission - Permission to check
   * @param projectId - Project public ID
   */
  function canAccessProject(permission: string, projectId: string): boolean {
    if (authStore.isSuperAdmin)
      return true;
    return can(permission) && isMemberOf(projectId);
  }

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
