/**
 * OAuth Utilities
 * Multi-provider OAuth 2.0 client library with PKCE support
 */

// Types
export type {
  CookieOptions,
  OAuthCallbackParams,
  OAuthProviderConfig,
  OAuthTokens,
  OAuthUser,
  PendingOAuthRequest,
} from './types';

// PKCE utilities
export { generateCodeChallenge, generateCodeVerifier, generateState } from './pkce';

// Cookie utilities
export {
  clearAuthCookies,
  getAuthCookies,
  hasValidAuthCookies,
  setAuthCookies,
} from './cookies';

// OAuth Client
export { OAuthClient, parseCallbackParams } from './client';

// Providers
export { altaluneProvider } from './providers/altalune';
