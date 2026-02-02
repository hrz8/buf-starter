import { usePermissions } from '~/composables/usePermissions';
import { PERMISSIONS } from '~/constants/permissions';
import { useAuthStore } from '~/stores/auth';

/**
 * Route permission requirements
 * Maps route patterns to required permissions
 */
const ROUTE_PERMISSIONS: Record<string, string | string[]> = {
  '/examples/datatable': PERMISSIONS.EMPLOYEE.READ,
  '/iam/users': PERMISSIONS.USER.READ,
  '/iam/roles': PERMISSIONS.ROLE.READ,
  '/iam/permissions': PERMISSIONS.PERMISSION.READ,
  '/iam/oauth-provider': PERMISSIONS.PROVIDER.READ,
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
  if (!authStore.isAuthenticated)
    return;

  // Check route against permission requirements
  for (const [pattern, permission] of Object.entries(ROUTE_PERMISSIONS)) {
    if (to.path.startsWith(pattern)) {
      const hasPermission = Array.isArray(permission)
        ? canAny(permission)
        : can(permission);

      if (!hasPermission)
        return navigateTo('/access-denied');

      break;
    }
  }
});
