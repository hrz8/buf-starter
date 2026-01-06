import { z } from 'zod';

/**
 * Base Permission field validations
 * Reused across create and update schemas
 */
const permissionNameValidation = z
  .string()
  .min(2, 'Name must be at least 2 characters')
  .max(100)
  .regex(/^[\w:]+$/, 'Name can only contain letters, numbers, underscores, and colons');

const permissionEffectValidation = z.enum(['allow', 'deny']);

const permissionDescriptionValidation = z.string().max(500).optional();

/**
 * Permission Creation Form Schema
 * Matches CreatePermissionRequest protobuf validation
 */
export const permissionSchema = z.object({
  name: permissionNameValidation,
  effect: permissionEffectValidation,
  description: permissionDescriptionValidation,
});

export type PermissionFormData = z.infer<typeof permissionSchema>;

/**
 * Permission Update Form Schema
 * Matches UpdatePermissionRequest protobuf validation
 * Note: name is read-only, only effect and description can be updated
 */
export const permissionUpdateSchema = z.object({
  id: z.string().min(1),
  effect: permissionEffectValidation,
  description: permissionDescriptionValidation,
});

export type PermissionUpdateFormData = z.infer<typeof permissionUpdateSchema>;
