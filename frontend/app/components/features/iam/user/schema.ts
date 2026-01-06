import { z } from 'zod';

/**
 * Base User field validations
 * Reused across create and update schemas
 */
const userEmailValidation = z
  .string()
  .email('Invalid email format')
  .min(1, 'Email is required');

const userFirstNameValidation = z
  .string()
  .min(1, 'First name is required')
  .max(100);

const userLastNameValidation = z
  .string()
  .min(1, 'Last name is required')
  .max(100);

/**
 * User Creation Form Schema
 * Matches CreateUserRequest protobuf validation
 */
export const userSchema = z.object({
  email: userEmailValidation,
  firstName: userFirstNameValidation,
  lastName: userLastNameValidation,
});

export type UserFormData = z.infer<typeof userSchema>;

/**
 * User Update Form Schema
 * Matches UpdateUserRequest protobuf validation
 */
export const userUpdateSchema = z.object({
  id: z.string().min(1),
  email: userEmailValidation,
  firstName: userFirstNameValidation,
  lastName: userLastNameValidation,
});

export type UserUpdateFormData = z.infer<typeof userUpdateSchema>;
