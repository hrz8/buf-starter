/**
 * Auth Store (Zustand)
 * Manages authentication state with cookie persistence and on-demand token refresh
 */

import { create } from 'zustand';

import type { OAuthTokens, OAuthUser } from './oauth';
import {
  altaluneProvider,
  clearAuthCookies,
  getAuthCookies,
  hasValidAuthCookies,
  OAuthClient,
  setAuthCookies,
} from './oauth';

// Refresh buffer - consider token expired if within this time of expiration
const REFRESH_BUFFER_MS = 60 * 1000; // 1 minute

// Minimum interval between refresh attempts (prevent rapid retries)
const MIN_REFRESH_INTERVAL_MS = 10 * 1000;

interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  isRefreshing: boolean;
  accessToken: string | null;
  refreshToken: string | null;
  user: OAuthUser | null;
  expiresAt: number | null;
  provider: string | null;
}

interface AuthActions {
  initFromCookies: () => void;
  setAuth: (tokens: OAuthTokens, user: OAuthUser, provider: string) => void;
  clearAuth: () => void;
  login: (returnTo?: string) => Promise<void>;
  logout: () => void;
  refreshTokens: () => Promise<boolean>;
  checkAndRefreshIfNeeded: () => Promise<boolean>;
  isTokenExpired: () => boolean;
}

type AuthStore = AuthState & AuthActions;

// OAuth client instance
const altaluneClient = new OAuthClient(altaluneProvider);

let lastRefreshAttempt = 0;
let refreshPromise: Promise<boolean> | null = null;

export const useAuthStore = create<AuthStore>((set, get) => ({
  // Initial state
  isAuthenticated: false,
  isLoading: true,
  isRefreshing: false,
  accessToken: null,
  refreshToken: null,
  user: null,
  expiresAt: null,
  provider: null,

  // Initialize from cookies (call on app mount)
  initFromCookies: () => {
    if (hasValidAuthCookies()) {
      const { accessToken, refreshToken, user, expiresAt } = getAuthCookies();
      set({
        isAuthenticated: true,
        isLoading: false,
        accessToken,
        refreshToken,
        user,
        expiresAt,
        provider: 'altalune',
      });
    } else {
      clearAuthCookies();
      set({
        isAuthenticated: false,
        isLoading: false,
        accessToken: null,
        refreshToken: null,
        user: null,
        expiresAt: null,
        provider: null,
      });
    }
  },

  // Set auth after successful login
  setAuth: (tokens: OAuthTokens, user: OAuthUser, provider: string) => {
    setAuthCookies(tokens, user);
    const expiresAt = Date.now() + tokens.expiresIn * 1000;
    set({
      isAuthenticated: true,
      isLoading: false,
      isRefreshing: false,
      accessToken: tokens.accessToken,
      refreshToken: tokens.refreshToken || null,
      user,
      expiresAt,
      provider,
    });
  },

  // Clear auth state and cookies
  clearAuth: () => {
    clearAuthCookies();
    set({
      isAuthenticated: false,
      isLoading: false,
      isRefreshing: false,
      accessToken: null,
      refreshToken: null,
      user: null,
      expiresAt: null,
      provider: null,
    });
  },

  // Start OAuth login flow
  login: async (returnTo?: string) => {
    await altaluneClient.authorize(returnTo);
  },

  // Logout
  logout: () => {
    clearAuthCookies();
    set({
      isAuthenticated: false,
      isLoading: false,
      isRefreshing: false,
      accessToken: null,
      refreshToken: null,
      user: null,
      expiresAt: null,
      provider: null,
    });
  },

  // Check if access token is expired or about to expire
  isTokenExpired: () => {
    const state = get();
    if (!state.expiresAt) return true;
    return Date.now() >= state.expiresAt - REFRESH_BUFFER_MS;
  },

  // Check token expiration and refresh if needed
  // Call this before accessing protected resources
  checkAndRefreshIfNeeded: async () => {
    const state = get();

    if (!state.isAuthenticated || !state.refreshToken) {
      return false;
    }

    if (!get().isTokenExpired()) {
      return true;
    }

    console.log('[Auth] Token expired, attempting refresh...');
    return get().refreshTokens();
  },

  refreshTokens: async () => {
    const state = get();

    if (refreshPromise) {
      return refreshPromise;
    }

    const now = Date.now();
    if (now - lastRefreshAttempt < MIN_REFRESH_INTERVAL_MS) {
      console.log('[Auth] Skipping refresh, too soon since last attempt');
      return state.isAuthenticated && !get().isTokenExpired();
    }
    lastRefreshAttempt = now;

    if (!state.refreshToken) {
      console.warn('[Auth] No refresh token available');
      get().clearAuth();
      return false;
    }

    set({ isRefreshing: true });

    refreshPromise = (async () => {
      try {
        const tokens = await altaluneClient.refreshToken(state.refreshToken!);
        const user = altaluneClient.extractUser(tokens.accessToken);

        setAuthCookies(tokens, user);
        const expiresAt = Date.now() + tokens.expiresIn * 1000;

        set({
          isRefreshing: false,
          accessToken: tokens.accessToken,
          refreshToken: tokens.refreshToken || state.refreshToken,
          user,
          expiresAt,
        });

        console.log('[Auth] Token refreshed successfully');
        return true;
      } catch (error) {
        const err = error as Error & { code?: string };
        console.error('[Auth] Token refresh failed:', err.message);

        if (err.code === 'invalid_grant') {
          console.warn('[Auth] Refresh token invalid, clearing auth');
          get().clearAuth();
        }

        set({ isRefreshing: false });
        return false;
      } finally {
        refreshPromise = null;
      }
    })();

    return refreshPromise;
  },
}));

// Export OAuth client for callback handling
export { altaluneClient };

export function isTokenExpiredError(status: number, errorCode?: string): boolean {
  return status === 401 && errorCode === 'invalid_token';
}

export function isRefreshTokenInvalidError(errorCode?: string): boolean {
  return errorCode === 'invalid_grant';
}
