import { z } from 'zod';

/**
 * Create OAuth Client Schema
 * OAuth clients are GLOBAL entities (not project-scoped)
 */
export const oauthClientCreateSchema = z.object({
  name: z
    .string()
    .min(1, 'Client name is required')
    .max(100, 'Client name must be at most 100 characters')
    .trim(),
  redirectUris: z
    .array(
      z.string().url('Must be a valid URL').trim(),
    )
    .min(1, 'At least one redirect URI is required'),
  pkceRequired: z.boolean(),
  allowedScopes: z.array(z.string()),
});

export type OAuthClientCreateFormData = z.infer<typeof oauthClientCreateSchema>;

/**
 * Update OAuth Client Schema
 * OAuth clients are GLOBAL entities (not project-scoped)
 */
export const oauthClientUpdateSchema = z.object({
  id: z.string().length(14, 'Client ID must be 14 characters'),
  name: z
    .string()
    .min(1, 'Client name is required')
    .max(100, 'Client name must be at most 100 characters')
    .trim()
    .optional(),
  redirectUris: z
    .array(
      z.string().url('Must be a valid URL').trim(),
    )
    .min(1, 'At least one redirect URI is required')
    .optional(),
  pkceRequired: z.boolean().optional(),
  allowedScopes: z.array(z.string()).optional(),
});

export type OAuthClientUpdateFormData = z.infer<typeof oauthClientUpdateSchema>;
