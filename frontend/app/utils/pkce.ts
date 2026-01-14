/**
 * PKCE (Proof Key for Code Exchange) utilities for OAuth 2.0 authorization code flow
 * @see https://datatracker.ietf.org/doc/html/rfc7636
 */

/**
 * Base64 URL encode a Uint8Array
 * Converts to base64 and replaces characters for URL safety
 */
function base64UrlEncode(buffer: Uint8Array): string {
  return btoa(String.fromCharCode(...buffer))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '');
}

/**
 * Generate a cryptographically random code verifier
 * The code verifier is a high-entropy cryptographic random string
 * between 43 and 128 characters in length
 */
export function generateCodeVerifier(): string {
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return base64UrlEncode(array);
}

/**
 * Generate a code challenge from a code verifier using SHA-256
 * The code challenge is the base64url-encoded SHA-256 hash of the code verifier
 */
export async function generateCodeChallenge(verifier: string): Promise<string> {
  const encoder = new TextEncoder();
  const data = encoder.encode(verifier);
  const hash = await crypto.subtle.digest('SHA-256', data);
  return base64UrlEncode(new Uint8Array(hash));
}

/**
 * Generate a cryptographically random state parameter for CSRF protection
 */
export function generateState(): string {
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return base64UrlEncode(array);
}
