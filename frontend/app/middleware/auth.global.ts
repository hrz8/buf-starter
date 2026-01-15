import { useAuthService } from '@/composables/useAuthService';
import { useAuthStore } from '@/stores/auth';

export default defineNuxtRouteMiddleware(async (to) => {
  const authStore = useAuthStore();
  const authService = useAuthService();

  if (to.path.startsWith('/auth/')) {
    return;
  }

  // If authenticated, check if token needs refresh
  if (authStore.isAuthenticated) {
    const hasValidAuth = await authService.checkAndRefreshIfNeeded();
    if (hasValidAuth) {
      return;
    }
  }

  try {
    const result = await authService.fetchCurrentUser();
    if (result) {
      return;
    }
  }
  catch {
    try {
      await authService.refreshTokens();
      return;
    }
    catch {
      // All attempts failed, redirect to login
    }
  }

  const nextUrl = to.fullPath !== '/' ? to.fullPath : undefined;
  authStore.setReturnUrl(nextUrl || null);

  const loginPath = nextUrl ? `/auth/login?next=${encodeURIComponent(nextUrl)}` : '/auth/login';
  return navigateTo(loginPath);
});
