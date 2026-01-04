import { z } from 'zod';

/**
 * Base API Key field validations
 * Reused across create and update schemas
 */
const apiKeyNameValidation = z
  .string()
  .min(2, 'Name must be at least 2 characters')
  .max(50, 'Name must not exceed 50 characters')
  .regex(/^[\w\s\-]+$/, 'Name can only contain letters, numbers, spaces, hyphens, and underscores');

const apiKeyExpirationValidation = z.string().min(1, 'Expiration date is required');

/**
 * API Key Creation Form Schema
 * Matches CreateAPIKeyRequest protobuf validation
 */
export const apiKeyCreateSchema = z.object({
  projectId: z.string().length(14),
  name: apiKeyNameValidation,
  expiration: apiKeyExpirationValidation,
});

export type ApiKeyCreateFormData = z.infer<typeof apiKeyCreateSchema>;

/**
 * API Key Update Form Schema
 * Matches UpdateAPIKeyRequest protobuf validation
 */
export const apiKeyUpdateSchema = z.object({
  projectId: z.string().length(14),
  apiKeyId: z.string().min(1),
  name: apiKeyNameValidation,
  expiration: apiKeyExpirationValidation,
});

export type ApiKeyUpdateFormData = z.infer<typeof apiKeyUpdateSchema>;
