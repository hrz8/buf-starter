import type { AuthExchangeResponse } from '~~/shared/repository/auth';
import { authRepository } from '~~/shared/repository/auth';
import { useAuthStore } from '@/stores/auth';

export function useAuthService() {
  const { $api } = useNuxtApp();
  const config = useRuntimeConfig();
  const authStore = useAuthStore();

  // Create API client for OAuth BFF endpoints using oauthBackendUrl
  const client = $api.createClient(config.public.oauthBackendUrl);
  const repo = authRepository(client);

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
    const result = await repo.refresh();
    authStore.setUser(result.user, result.expires_in);
    return result;
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

  return {
    handleCallback,
    refreshTokens,
    fetchCurrentUser,
    logout,
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
