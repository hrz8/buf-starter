/**
 * Cookie Management Utilities
 */

import type { CookieOptions, OAuthTokens, OAuthUser } from './types';

const AUTH_COOKIES = {
  ACCESS_TOKEN: 'auth_access_token',
  REFRESH_TOKEN: 'auth_refresh_token',
  USER: 'auth_user',
  EXPIRES_AT: 'auth_expires_at',
} as const;

const DEFAULT_OPTIONS: CookieOptions = {
  path: '/',
  sameSite: 'Lax',
};

/**
 * Set authentication cookies
 */
export function setAuthCookies(
  tokens: OAuthTokens,
  user: OAuthUser,
  options: CookieOptions = {}
): void {
  const opts = { ...DEFAULT_OPTIONS, ...options };
  const expiresAt = Date.now() + tokens.expiresIn * 1000;

  setCookie(AUTH_COOKIES.ACCESS_TOKEN, tokens.accessToken, {
    ...opts,
    maxAge: tokens.expiresIn,
  });

  if (tokens.refreshToken) {
    setCookie(AUTH_COOKIES.REFRESH_TOKEN, tokens.refreshToken, {
      ...opts,
      maxAge: 60 * 60 * 24 * 30, // 30 days
    });
  }

  setCookie(AUTH_COOKIES.USER, encodeURIComponent(JSON.stringify(user)), {
    ...opts,
    maxAge: tokens.expiresIn,
  });

  setCookie(AUTH_COOKIES.EXPIRES_AT, String(expiresAt), {
    ...opts,
    maxAge: tokens.expiresIn,
  });
}

/**
 * Get authentication data from cookies
 */
export function getAuthCookies(): {
  accessToken: string | null;
  refreshToken: string | null;
  user: OAuthUser | null;
  expiresAt: number | null;
} {
  const accessToken = getCookie(AUTH_COOKIES.ACCESS_TOKEN);
  const refreshToken = getCookie(AUTH_COOKIES.REFRESH_TOKEN);
  const userStr = getCookie(AUTH_COOKIES.USER);
  const expiresAtStr = getCookie(AUTH_COOKIES.EXPIRES_AT);

  let user: OAuthUser | null = null;
  if (userStr) {
    try {
      user = JSON.parse(decodeURIComponent(userStr));
    } catch {
      // Invalid user cookie
    }
  }

  return {
    accessToken,
    refreshToken,
    user,
    expiresAt: expiresAtStr ? parseInt(expiresAtStr, 10) : null,
  };
}

/**
 * Clear all authentication cookies
 */
export function clearAuthCookies(): void {
  deleteCookie(AUTH_COOKIES.ACCESS_TOKEN);
  deleteCookie(AUTH_COOKIES.REFRESH_TOKEN);
  deleteCookie(AUTH_COOKIES.USER);
  deleteCookie(AUTH_COOKIES.EXPIRES_AT);
}

/**
 * Check if auth cookies exist and are not expired
 */
export function hasValidAuthCookies(): boolean {
  const { accessToken, expiresAt } = getAuthCookies();
  if (!accessToken || !expiresAt) return false;
  return Date.now() < expiresAt;
}

// Cookie helpers
function setCookie(name: string, value: string, options: CookieOptions): void {
  let cookie = `${name}=${value}`;

  if (options.maxAge !== undefined) {
    cookie += `; max-age=${options.maxAge}`;
  }
  if (options.path) {
    cookie += `; path=${options.path}`;
  }
  if (options.sameSite) {
    cookie += `; samesite=${options.sameSite}`;
  }
  if (options.secure) {
    cookie += '; secure';
  }

  document.cookie = cookie;
}

function getCookie(name: string): string | null {
  const cookies = document.cookie.split(';');
  for (const cookie of cookies) {
    const [key, value] = cookie.trim().split('=');
    if (key === name) {
      return value || null;
    }
  }
  return null;
}

function deleteCookie(name: string): void {
  document.cookie = `${name}=; path=/; max-age=0`;
}
