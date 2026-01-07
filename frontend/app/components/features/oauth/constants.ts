import { ProviderType as ProtoProviderType } from '~~/gen/altalune/v1/oauth_provider_pb';
import { ProviderTypeEnum } from './schema';

/**
 * Provider metadata for displaying OAuth provider information
 */
export interface ProviderMetadata {
  name: string;
  icon: string; // Icon name or class (e.g., lucide icon name)
  defaultScopes: string;
  docsUrl: string;
  color: string; // Tailwind color class for badge
}

/**
 * OAuth provider metadata mapping
 * Provides display information, default scopes, and documentation links
 */
export const PROVIDER_METADATA: Record<number, ProviderMetadata> = {
  [ProviderTypeEnum.PROVIDER_TYPE_GOOGLE]: {
    name: 'Google',
    icon: 'i-lucide-chrome', // Using lucide Chrome icon for Google
    defaultScopes: 'openid,email,profile',
    docsUrl: 'https://developers.google.com/identity/protocols/oauth2',
    color: 'blue',
  },
  [ProviderTypeEnum.PROVIDER_TYPE_GITHUB]: {
    name: 'GitHub',
    icon: 'i-lucide-github',
    defaultScopes: 'user:email,read:user',
    docsUrl: 'https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps',
    color: 'gray',
  },
  [ProviderTypeEnum.PROVIDER_TYPE_MICROSOFT]: {
    name: 'Microsoft',
    icon: 'i-lucide-boxes', // Using lucide Boxes icon for Microsoft
    defaultScopes: 'openid,email,profile',
    docsUrl: 'https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-auth-code-flow',
    color: 'sky',
  },
  [ProviderTypeEnum.PROVIDER_TYPE_APPLE]: {
    name: 'Apple',
    icon: 'i-lucide-apple',
    defaultScopes: 'name,email',
    docsUrl: 'https://developer.apple.com/documentation/sign_in_with_apple',
    color: 'slate',
  },
};

/**
 * Get provider metadata by provider type
 */
export function getProviderMetadata(providerType: number): ProviderMetadata | undefined {
  return PROVIDER_METADATA[providerType];
}

/**
 * Get provider name by provider type
 */
export function getProviderName(providerType: number): string {
  return PROVIDER_METADATA[providerType]?.name ?? 'Unknown';
}

/**
 * Provider type dropdown options for form select
 */
export const PROVIDER_TYPE_OPTIONS = [
  {
    value: ProviderTypeEnum.PROVIDER_TYPE_GOOGLE,
    label: 'Google',
    icon: 'i-lucide-chrome',
    protoValue: ProtoProviderType.GOOGLE,
  },
  {
    value: ProviderTypeEnum.PROVIDER_TYPE_GITHUB,
    label: 'GitHub',
    icon: 'i-lucide-github',
    protoValue: ProtoProviderType.GITHUB,
  },
  {
    value: ProviderTypeEnum.PROVIDER_TYPE_MICROSOFT,
    label: 'Microsoft',
    icon: 'i-lucide-boxes',
    protoValue: ProtoProviderType.MICROSOFT,
  },
  {
    value: ProviderTypeEnum.PROVIDER_TYPE_APPLE,
    label: 'Apple',
    icon: 'i-lucide-apple',
    protoValue: ProtoProviderType.APPLE,
  },
];
