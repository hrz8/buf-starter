/**
 * OAuth Client role permissions
 */
export const OAUTH_CLIENT_PERMISSIONS = {
  CREATE: ['owner', 'admin'],
  UPDATE: ['owner', 'admin'],
  DELETE: ['owner'], // Only owner can delete
  REVEAL_SECRET: ['owner', 'admin'],
  VIEW: ['owner', 'admin', 'member'], // Member read-only
} as const;

/**
 * PKCE options
 */
export const PKCE_OPTIONS = [
  { value: true, label: 'Required (Recommended for public clients)' },
  { value: false, label: 'Not Required (Confidential clients only)' },
] as const;

/**
 * Default scopes
 */
export const DEFAULT_SCOPES = [
  'openid',
  'profile',
  'email',
  'offline_access',
] as const;
