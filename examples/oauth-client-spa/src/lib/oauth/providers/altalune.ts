/**
 * Altalune OAuth Provider Configuration
 */

import type { OAuthProviderConfig } from '../types';

const authServerUrl = import.meta.env.VITE_ALTALUNE_AUTH_SERVER || 'http://localhost:3300';
const clientId = import.meta.env.VITE_ALTALUNE_CLIENT_ID || '';

export const altaluneProvider: OAuthProviderConfig = {
  id: 'altalune',
  name: 'Altalune',
  authorizationUrl: `${authServerUrl}/oauth/authorize`,
  tokenUrl: `${authServerUrl}/oauth/token`,
  clientId,
  redirectUri: 'http://localhost:5173/auth/altalune/callback',
  scopes: ['openid', 'profile', 'email'],
  pkceRequired: true,
};
