import { useAuthService } from '@/composables/useAuthService';
import { useAuthStore } from '@/stores/auth';

export default defineNuxtRouteMiddleware(async (to) => {
  const authStore = useAuthStore();
  const authService = useAuthService();

  // Skip auth check for auth-related pages
  if (to.path.startsWith('/auth/')) {
    return;
  }

  // Skip auth check for access-denied page (prevent redirect loop)
  if (to.path === '/access-denied') {
    authStore.setInitialized();
    return;
  }

  // If already authenticated and initialized, just check token refresh
  if (authStore.isAuthenticated) {
    const hasValidAuth = await authService.checkAndRefreshIfNeeded();
    if (hasValidAuth) {
      return;
    }
  }

  // Try to fetch current user (validates session cookie)
  const result = await authService.fetchCurrentUser();
  if (result) {
    // Auth successful - isInitialized is set in setUser()
    return;
  }

  // fetchCurrentUser failed (access_token expired or missing)
  // Try to refresh tokens using refresh_token cookie
  try {
    await authService.refreshTokens();
    // Refresh successful - tokens renewed, user set in store
    return;
  }
  catch {
    // Refresh also failed - no valid session
    authStore.setInitialized();
  }

  // Not authenticated - redirect to login
  const nextUrl = to.fullPath !== '/' ? to.fullPath : undefined;
  authStore.setReturnUrl(nextUrl || null);

  const loginPath = nextUrl ? `/auth/login?next=${encodeURIComponent(nextUrl)}` : '/auth/login';
  return navigateTo(loginPath);
});
