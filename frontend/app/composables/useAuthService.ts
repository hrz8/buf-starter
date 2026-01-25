import type { AuthExchangeResponse } from '~~/shared/repository/auth';
import { authRepository } from '~~/shared/repository/auth';
import { useAuthStore } from '@/stores/auth';

// Module-level state for refresh deduplication (shared across all instances)
let refreshPromise: Promise<AuthExchangeResponse> | null = null;

// Minimum interval between refresh attempts (prevent rapid retries)
const MIN_REFRESH_INTERVAL_MS = 10 * 1000;
let lastRefreshAttempt = 0;

export function useAuthService() {
  const { $api } = useNuxtApp();
  const config = useRuntimeConfig();
  const authStore = useAuthStore();

  // Create API client for OAuth BFF endpoints using oauthBackendUrl
  const bffClient = $api.createClient(config.public.oauthBackendUrl);
  const repo = authRepository(bffClient);

  // Create API client for Auth Server direct endpoints (e.g., resend-verification)
  const authServerClient = $api.createClient(config.public.authServerUrl);
  const authServerRepo = authRepository(authServerClient);

  async function handleCallback(code: string, state: string): Promise<AuthExchangeResponse> {
    // Validate state parameter
    const storedState = sessionStorage.getItem('oauth_state');
    if (state !== storedState) {
      throw new AuthError('invalid_state', 'Invalid state parameter');
    }

    // Get code verifier
    const codeVerifier = sessionStorage.getItem('oauth_code_verifier');
    if (!codeVerifier) {
      throw new AuthError('missing_verifier', 'Missing code verifier');
    }

    // Exchange code for tokens (backend sets httpOnly cookies)
    const result = await repo.exchange({
      code,
      code_verifier: codeVerifier,
      redirect_uri: config.public.oauthRedirectUri,
    });

    // Store user info in Pinia (token is in httpOnly cookie)
    authStore.setUser(result.user, result.expires_in);

    // Clear PKCE data from sessionStorage
    sessionStorage.removeItem('oauth_state');
    sessionStorage.removeItem('oauth_code_verifier');

    return result;
  }

  async function refreshTokens(): Promise<AuthExchangeResponse> {
    // Deduplicate concurrent refresh calls
    if (refreshPromise) {
      return refreshPromise;
    }

    // Prevent rapid refresh attempts
    const now = Date.now();
    if (now - lastRefreshAttempt < MIN_REFRESH_INTERVAL_MS) {
      throw new AuthError('refresh_throttled', 'Refresh attempted too recently');
    }
    lastRefreshAttempt = now;

    refreshPromise = (async () => {
      try {
        const result = await repo.refresh();
        authStore.setUser(result.user, result.expires_in);
        return result;
      }
      catch (error) {
        console.error('[Auth] Token refresh failed:', error);
        throw error;
      }
      finally {
        refreshPromise = null;
      }
    })();

    return refreshPromise;
  }

  /**
   * Check if token is expired or about to expire (within 1 minute)
   */
  function isTokenExpired(): boolean {
    const expiresAt = authStore.expiresAt;
    if (!expiresAt) {
      return true;
    }
    // Consider expired if within 1 minute of expiration
    return Date.now() >= expiresAt - 60 * 1000;
  }

  /**
   * Check token expiration and refresh if needed.
   * Call this before accessing protected resources.
   * Returns true if we have valid auth, false otherwise.
   */
  async function checkAndRefreshIfNeeded(): Promise<boolean> {
    // If not authenticated at all, return false
    if (!authStore.isAuthenticated) {
      return false;
    }

    // Token still valid
    if (!isTokenExpired()) {
      return true;
    }

    // Token expired, attempt refresh
    try {
      await refreshTokens();
      return true;
    }
    catch (error) {
      const err = error as { code?: string; status?: number };
      // If refresh failed due to invalid grant, clear auth
      if (err.code === 'invalid_grant' || err.status === 401) {
        console.warn('[Auth] Refresh token invalid, clearing auth');
        authStore.clearAuth();
      }
      return false;
    }
  }

  async function fetchCurrentUser(): Promise<AuthExchangeResponse | null> {
    try {
      const result = await repo.me();
      authStore.setUser(result.user, result.expires_in);
      return result;
    }
    catch {
      // Not authenticated or token expired
      return null;
    }
  }

  async function logout(): Promise<void> {
    try {
      await repo.logout();
    }
    catch {
      // Ignore errors during logout
    }

    authStore.clearAuth();
    navigateTo('/auth/login');
  }

  async function resendVerificationEmail(): Promise<void> {
    try {
      await authServerRepo.resendVerification();
    }
    catch (error) {
      const err = error as { data?: { error?: string; error_description?: string } };
      throw new AuthError(
        err.data?.error || 'resend_failed',
        err.data?.error_description || 'Failed to resend verification email',
      );
    }
  }

  return {
    handleCallback,
    refreshTokens,
    fetchCurrentUser,
    logout,
    isTokenExpired,
    checkAndRefreshIfNeeded,
    resendVerificationEmail,
  };
}

// Custom error class for auth errors
export class AuthError extends Error {
  code: string;

  constructor(code: string, message: string) {
    super(message);
    this.code = code;
    this.name = 'AuthError';
  }
}
