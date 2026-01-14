import { useAuthService } from '@/composables/useAuthService';
import { useAuthStore } from '@/stores/auth';

export default defineNuxtRouteMiddleware(async (to) => {
  const authStore = useAuthStore();
  const authService = useAuthService();

  if (to.path.startsWith('/auth/')) {
    return;
  }

  if (authStore.isAuthenticated) {
    const expiresAt = authStore.expiresAt;
    if (expiresAt && Date.now() > expiresAt - 60000) {
      try {
        await authService.refreshTokens();
        return;
      }
      catch {
        // Refresh failed, will redirect to login below
      }
    }
    else {
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
