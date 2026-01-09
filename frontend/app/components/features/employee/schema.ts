import { z } from 'zod';

/**
 * Create Employee Schema
 */
export const employeeCreateSchema = z.object({
  projectId: z.string().length(14, 'Project ID must be 14 characters'),
  name: z
    .string()
    .min(2, 'Name must be at least 2 characters')
    .max(50, 'Name must be at most 50 characters')
    .trim(),
  email: z
    .string()
    .email('Must be a valid email address')
    .trim(),
  role: z
    .string()
    .min(1, 'Role is required')
    .trim(),
  department: z
    .string()
    .min(1, 'Department is required')
    .trim(),
  status: z
    .number()
    .int()
    .min(0, 'Status must be a valid value'),
});

export type EmployeeCreateFormData = z.infer<typeof employeeCreateSchema>;

/**
 * Update Employee Schema
 */
export const employeeUpdateSchema = z.object({
  projectId: z.string().length(14, 'Project ID must be 14 characters'),
  employeeId: z.string().min(1, 'Employee ID is required'),
  name: z
    .string()
    .min(2, 'Name must be at least 2 characters')
    .max(50, 'Name must be at most 50 characters')
    .trim(),
  email: z
    .string()
    .email('Must be a valid email address')
    .trim(),
  role: z
    .string()
    .min(1, 'Role is required')
    .trim(),
  department: z
    .string()
    .min(1, 'Department is required')
    .trim(),
  status: z
    .number()
    .int()
    .min(0, 'Status must be a valid value'),
});

export type EmployeeUpdateFormData = z.infer<typeof employeeUpdateSchema>;
