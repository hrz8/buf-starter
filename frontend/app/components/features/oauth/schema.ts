import { z } from 'zod';

/**
 * ProviderType enum values matching proto definition
 */
export const ProviderTypeEnum = {
  PROVIDER_TYPE_UNSPECIFIED: 0,
  PROVIDER_TYPE_GOOGLE: 1,
  PROVIDER_TYPE_GITHUB: 2,
  PROVIDER_TYPE_MICROSOFT: 3,
  PROVIDER_TYPE_APPLE: 4,
} as const;

export type ProviderType = typeof ProviderTypeEnum[keyof typeof ProviderTypeEnum];

/**
 * Zod schema for creating a new OAuth provider
 * Matches CreateOAuthProviderRequest proto validation rules
 * NO .default() - defaults handled in initialValues
 */
export const createOAuthProviderSchema = z.object({
  providerType: z
    .number()
    .int()
    .min(1, 'Provider type is required')
    .max(4, 'Invalid provider type'),
  clientId: z
    .string()
    .min(1, 'Client ID is required')
    .max(500, 'Client ID must be at most 500 characters'),
  clientSecret: z
    .string()
    .min(1, 'Client secret is required')
    .max(500, 'Client secret must be at most 500 characters'),
  redirectUrl: z
    .string()
    .url('Redirect URL must be a valid URL')
    .max(500, 'Redirect URL must be at most 500 characters'),
  scopes: z
    .string()
    .max(1000, 'Scopes must be at most 1000 characters'),
  enabled: z.boolean(),
});

export type CreateOAuthProviderInput = z.infer<typeof createOAuthProviderSchema>;

/**
 * Zod schema for updating an existing OAuth provider
 * Matches UpdateOAuthProviderRequest proto validation rules
 * NOTE: provider_type is NOT included (immutable after creation)
 * NO .default() - defaults handled in initialValues
 */
export const updateOAuthProviderSchema = z.object({
  id: z
    .string()
    .min(14, 'Provider ID must be at least 14 characters')
    .max(20, 'Provider ID must be at most 20 characters'),
  clientId: z
    .string()
    .min(1, 'Client ID is required')
    .max(500, 'Client ID must be at most 500 characters'),
  clientSecret: z
    .string()
    .max(500, 'Client secret must be at most 500 characters'),
  redirectUrl: z
    .string()
    .url('Redirect URL must be a valid URL')
    .max(500, 'Redirect URL must be at most 500 characters'),
  scopes: z
    .string()
    .max(1000, 'Scopes must be at most 1000 characters'),
  enabled: z.boolean(),
});

export type UpdateOAuthProviderInput = z.infer<typeof updateOAuthProviderSchema>;
