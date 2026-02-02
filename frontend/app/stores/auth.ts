import type { AuthUserInfo } from '~~/shared/repository/auth';

const RETURN_URL_KEY = 'oauth_return_url';
const ROOT_PERMISSION = 'root';

export const useAuthStore = defineStore('auth', () => {
  // User info from JWT (stored in memory, not the token itself)
  const user = ref<AuthUserInfo | null>(null);
  const expiresAt = ref<number | null>(null);

  /**
   * Tracks if auth initialization has completed (middleware checked auth)
   * This is set to true after the first auth check, regardless of whether
   * the user is authenticated or not.
   */
  const isInitialized = ref(false);

  const isAuthenticated = computed(() => {
    return !!user.value && (expiresAt.value ? Date.now() < expiresAt.value : true);
  });

  /**
   * Returns true when auth is fully ready (initialized AND authenticated)
   * Use this to gate API calls and content rendering
   */
  const isReady = computed(() => {
    return isInitialized.value && isAuthenticated.value;
  });

  // Email verification status computed properties
  const isEmailVerified = computed(() => {
    return user.value?.email_verified ?? false;
  });

  const isEmailVerificationRequired = computed(() => {
    return isAuthenticated.value && !isEmailVerified.value;
  });

  // ============================================
  // Permission-related computed properties
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
  // Methods
  // ============================================

  function setUser(userData: AuthUserInfo, expiresIn: number) {
    user.value = userData;
    expiresAt.value = Date.now() + expiresIn * 1000;
    isInitialized.value = true;
  }

  function clearAuth() {
    user.value = null;
    expiresAt.value = null;
  }

  /**
   * Mark auth as initialized (called by auth middleware after checking auth)
   */
  function setInitialized() {
    isInitialized.value = true;
  }

  // Return URL is stored in sessionStorage to persist across OAuth redirects
  // (Pinia store is in-memory and lost when page redirects to auth server)
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
    isInitialized: readonly(isInitialized),
    isReady,
    isEmailVerified,
    isEmailVerificationRequired,
    // Permission exports
    permissions,
    memberships,
    isSuperAdmin,
    memberProjectIds,
    // Methods
    setUser,
    clearAuth,
    setInitialized,
    setReturnUrl,
    getAndClearReturnUrl,
  };
});
