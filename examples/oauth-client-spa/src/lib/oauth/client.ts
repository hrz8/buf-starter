/**
 * OAuth Client
 * Provider-agnostic OAuth 2.0 client with PKCE support
 */

import type {
  OAuthCallbackParams,
  OAuthProviderConfig,
  OAuthTokens,
  OAuthUser,
  PendingOAuthRequest,
} from './types';
import { generateCodeChallenge, generateCodeVerifier, generateState } from './pkce';

const PENDING_AUTH_KEY = 'oauth_pending_auth';

export class OAuthClient {
  private provider: OAuthProviderConfig;

  constructor(provider: OAuthProviderConfig) {
    this.provider = provider;
  }

  /**
   * Start OAuth authorization flow
   */
  async authorize(returnTo?: string): Promise<void> {
    const state = generateState();
    const codeVerifier = generateCodeVerifier();
    const codeChallenge = await generateCodeChallenge(codeVerifier);

    // Store pending auth request
    const pendingAuth: PendingOAuthRequest = {
      state,
      codeVerifier,
      providerId: this.provider.id,
      returnTo,
      createdAt: Date.now(),
    };
    sessionStorage.setItem(PENDING_AUTH_KEY, JSON.stringify(pendingAuth));

    // Build authorization URL
    const params = new URLSearchParams({
      response_type: 'code',
      client_id: this.provider.clientId,
      redirect_uri: this.provider.redirectUri,
      scope: this.provider.scopes.join(' '),
      state,
      ...(this.provider.pkceRequired && {
        code_challenge: codeChallenge,
        code_challenge_method: 'S256',
      }),
      ...this.provider.authParams,
    });

    const authUrl = `${this.provider.authorizationUrl}?${params.toString()}`;
    window.location.href = authUrl;
  }

  /**
   * Handle OAuth callback
   */
  async handleCallback(params: OAuthCallbackParams): Promise<{
    tokens: OAuthTokens;
    user: OAuthUser;
    returnTo?: string;
  }> {
    // Check for errors
    if (params.error) {
      throw new Error(params.errorDescription || params.error);
    }

    // Get pending auth request
    const pendingAuthStr = sessionStorage.getItem(PENDING_AUTH_KEY);
    if (!pendingAuthStr) {
      throw new Error('No pending OAuth request found');
    }

    const pendingAuth: PendingOAuthRequest = JSON.parse(pendingAuthStr);

    // Validate state
    if (pendingAuth.state !== params.state) {
      throw new Error('Invalid state parameter');
    }

    // Validate provider
    if (pendingAuth.providerId !== this.provider.id) {
      throw new Error('Provider mismatch');
    }

    // Check expiration (5 minutes)
    if (Date.now() - pendingAuth.createdAt > 5 * 60 * 1000) {
      throw new Error('OAuth request expired');
    }

    // Exchange code for tokens
    const tokens = await this.exchangeCode(params.code, pendingAuth.codeVerifier);

    // Extract user from access token
    const user = this.extractUserFromToken(tokens.accessToken);

    // Clear pending auth
    sessionStorage.removeItem(PENDING_AUTH_KEY);

    return {
      tokens,
      user,
      returnTo: pendingAuth.returnTo,
    };
  }

  /**
   * Exchange authorization code for tokens
   */
  private async exchangeCode(code: string, codeVerifier: string): Promise<OAuthTokens> {
    const body = new URLSearchParams({
      grant_type: 'authorization_code',
      code,
      redirect_uri: this.provider.redirectUri,
      client_id: this.provider.clientId,
      ...(this.provider.pkceRequired && { code_verifier: codeVerifier }),
    });

    const response = await fetch(this.provider.tokenUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.error_description || error.error || 'Token exchange failed');
    }

    const data = await response.json();

    return {
      accessToken: data.access_token,
      refreshToken: data.refresh_token,
      tokenType: data.token_type,
      expiresIn: data.expires_in,
      scope: data.scope,
    };
  }

  /**
   * Extract user info from JWT access token
   */
  private extractUserFromToken(accessToken: string): OAuthUser {
    try {
      const parts = accessToken.split('.');
      if (parts.length !== 3) {
        throw new Error('Invalid JWT');
      }

      const payload = JSON.parse(atob(parts[1]));

      return {
        sub: payload.sub,
        name: payload.name,
        email: payload.email,
        picture: payload.picture,
        ...payload,
      };
    } catch {
      return { sub: 'unknown' };
    }
  }

  /**
   * Refresh access token using refresh token
   */
  async refreshToken(refreshToken: string): Promise<OAuthTokens> {
    const body = new URLSearchParams({
      grant_type: 'refresh_token',
      refresh_token: refreshToken,
      client_id: this.provider.clientId,
    });

    const response = await fetch(this.provider.tokenUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      const errorCode = error.error || 'unknown_error';
      const errorMessage = error.error_description || 'Token refresh failed';

      // Create error with code for handling
      const refreshError = new Error(errorMessage) as Error & { code: string };
      refreshError.code = errorCode;
      throw refreshError;
    }

    const data = await response.json();

    return {
      accessToken: data.access_token,
      refreshToken: data.refresh_token,
      tokenType: data.token_type,
      expiresIn: data.expires_in,
      scope: data.scope,
    };
  }

  /**
   * Extract user info from JWT access token
   */
  extractUser(accessToken: string): OAuthUser {
    return this.extractUserFromToken(accessToken);
  }

  /**
   * Get provider info
   */
  getProvider(): OAuthProviderConfig {
    return this.provider;
  }
}

/**
 * Parse callback URL parameters
 */
export function parseCallbackParams(search: string): OAuthCallbackParams {
  const params = new URLSearchParams(search);
  return {
    code: params.get('code') || '',
    state: params.get('state') || '',
    error: params.get('error') || undefined,
    errorDescription: params.get('error_description') || undefined,
  };
}
