import type { AuthUserInfo } from '~~/shared/repository/auth';

const RETURN_URL_KEY = 'oauth_return_url';

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

  function setUser(userData: AuthUserInfo, expiresIn: number) {
    user.value = userData;
    expiresAt.value = Date.now() + expiresIn * 1000;
  }

  function clearAuth() {
    user.value = null;
    expiresAt.value = null;
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
    isEmailVerified,
    isEmailVerificationRequired,
    setUser,
    clearAuth,
    setReturnUrl,
    getAndClearReturnUrl,
  };
});
