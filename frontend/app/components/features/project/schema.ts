import { z } from 'zod';

/**
 * Project Settings Form Schema
 * Matches UpdateProjectRequest protobuf validation
 */
export const projectSettingsSchema = z.object({
  name: z.string().min(1, 'Name is required').max(50, 'Name must be 50 characters or less'),
  description: z.string().max(100, 'Description must be 100 characters or less').optional(),
  timezone: z.string().min(1, 'Timezone is required'),
});

export type ProjectSettingsFormData = z.infer<typeof projectSettingsSchema>;

/**
 * Project Creation Form Schema
 * Matches CreateProjectRequest protobuf validation
 */
export const projectCreateSchema = z.object({
  name: z.string().min(1, 'Name is required').max(50, 'Name must be 50 characters or less'),
  description: z.string().max(100, 'Description must be 100 characters or less').optional(),
  timezone: z.string().min(1, 'Timezone is required'),
  environment: z.enum(['sandbox', 'live']).optional(),
});

export type ProjectCreateFormData = z.infer<typeof projectCreateSchema>;
