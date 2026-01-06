import { z } from 'zod';

/**
 * Base Role field validations
 * Reused across create and update schemas
 */
const roleNameValidation = z
  .string()
  .min(2, 'Name must be at least 2 characters')
  .max(100)
  .regex(/^\w+$/, 'Name can only contain letters, numbers, and underscores');

const roleDescriptionValidation = z.string().max(500).optional();

/**
 * Role Creation Form Schema
 * Matches CreateRoleRequest protobuf validation
 */
export const roleSchema = z.object({
  name: roleNameValidation,
  description: roleDescriptionValidation,
});

export type RoleFormData = z.infer<typeof roleSchema>;

/**
 * Role Update Form Schema
 * Matches UpdateRoleRequest protobuf validation
 */
export const roleUpdateSchema = z.object({
  id: z.string().min(1),
  name: roleNameValidation,
  description: roleDescriptionValidation,
});

export type RoleUpdateFormData = z.infer<typeof roleUpdateSchema>;
