/**
 * OAuth Provider Configuration
 * Provider-agnostic interface for configuring OAuth providers
 */
export interface OAuthProviderConfig {
  id: string;
  name: string;
  authorizationUrl: string;
  tokenUrl: string;
  clientId: string;
  redirectUri: string;
  scopes: string[];
  pkceRequired: boolean;
  authParams?: Record<string, string>;
}

/**
 * OAuth Tokens Response
 */
export interface OAuthTokens {
  accessToken: string;
  refreshToken?: string;
  tokenType: string;
  expiresIn: number;
  scope?: string;
}

/**
 * OAuth User from JWT claims
 */
export interface OAuthUser {
  sub: string;
  name?: string;
  email?: string;
  picture?: string;
  [key: string]: unknown;
}

/**
 * Pending OAuth Request (stored during auth flow)
 */
export interface PendingOAuthRequest {
  state: string;
  codeVerifier: string;
  providerId: string;
  returnTo?: string;
  createdAt: number;
}

/**
 * OAuth Callback Parameters
 */
export interface OAuthCallbackParams {
  code: string;
  state: string;
  error?: string;
  errorDescription?: string;
}

/**
 * Cookie Options
 */
export interface CookieOptions {
  maxAge?: number;
  path?: string;
  sameSite?: 'Strict' | 'Lax' | 'None';
  secure?: boolean;
}
