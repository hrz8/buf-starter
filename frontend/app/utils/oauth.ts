/**
 * OAuth 2.0 utilities for initiating the authorization code flow with PKCE
 */

import { generateCodeChallenge, generateCodeVerifier, generateState } from './pkce';

export interface OAuthConfig {
  authServerUrl: string;
  clientId: string;
  redirectUri: string;
  scopes: string[];
}

/**
 * Initiate OAuth authorization code flow with PKCE
 * Generates PKCE parameters, stores them in sessionStorage, and redirects to auth server
 *
 * @param config OAuth configuration
 * @param forceConsent Force the consent screen to be shown (useful for testing)
 */
export async function initiateOAuthFlow(config: OAuthConfig, forceConsent = true) {
  // 1. Generate PKCE code verifier and challenge
  const codeVerifier = generateCodeVerifier();
  const codeChallenge = await generateCodeChallenge(codeVerifier);

  // 2. Generate state for CSRF protection
  const state = generateState();

  // 3. Store in sessionStorage (cleared on browser close for security)
  sessionStorage.setItem('oauth_code_verifier', codeVerifier);
  sessionStorage.setItem('oauth_state', state);

  // 4. Build authorization URL
  const params = new URLSearchParams({
    response_type: 'code',
    client_id: config.clientId,
    redirect_uri: config.redirectUri,
    scope: config.scopes.join(' '),
    state,
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
  });

  // Always force consent for now (user requirement: single button that forces consent)
  if (forceConsent) {
    params.set('prompt', 'consent');
  }

  // 5. Redirect to authorization server
  window.location.href = `${config.authServerUrl}/oauth/authorize?${params.toString()}`;
}
