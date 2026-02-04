import { z } from 'zod';
import { NODE_NAME_PATTERN, TRIGGER_TYPES, VALIDATION_LIMITS } from './constants';

// Trigger schema
export const triggerSchema = z.object({
  type: z.enum(TRIGGER_TYPES),
  value: z.string()
    .min(1, 'Trigger value is required')
    .max(VALIDATION_LIMITS.triggerValueMaxLength, `Trigger value must be ${VALIDATION_LIMITS.triggerValueMaxLength} characters or less`),
});

export type TriggerFormData = z.infer<typeof triggerSchema>;

// Message schema
export const messageSchema = z.object({
  role: z.literal('assistant'),
  content: z.string()
    .min(1, 'Message content is required')
    .max(VALIDATION_LIMITS.messageContentMaxLength, `Message must be ${VALIDATION_LIMITS.messageContentMaxLength} characters or less`),
});

export type MessageFormData = z.infer<typeof messageSchema>;

// Node create schema (for creating new nodes)
// Note: version is NOT included here - users create default nodes first,
// then use "Add Version" to create variants
export const nodeCreateSchema = z.object({
  name: z.string()
    .min(VALIDATION_LIMITS.nameMinLength, `Name must be at least ${VALIDATION_LIMITS.nameMinLength} characters`)
    .max(VALIDATION_LIMITS.nameMaxLength, `Name must be ${VALIDATION_LIMITS.nameMaxLength} characters or less`)
    .regex(NODE_NAME_PATTERN, 'Name must be lowercase letters, numbers, and underscores (e.g., greeting, faq_pricing)'),
  lang: z.string().min(1, 'Language is required'),
  tags: z.array(z.string().max(VALIDATION_LIMITS.tagMaxLength)).optional(),
});

export type NodeCreateFormData = z.infer<typeof nodeCreateSchema>;

// Node edit schema (for editing existing nodes)
// Note: triggers are optional - a node can match based on conditions only
export const nodeEditSchema = z.object({
  name: z.string()
    .min(VALIDATION_LIMITS.nameMinLength, `Name must be at least ${VALIDATION_LIMITS.nameMinLength} characters`)
    .max(VALIDATION_LIMITS.nameMaxLength, `Name must be ${VALIDATION_LIMITS.nameMaxLength} characters or less`)
    .regex(NODE_NAME_PATTERN, 'Name must be lowercase letters, numbers, and underscores'),
  tags: z.array(z.string().max(VALIDATION_LIMITS.tagMaxLength)).optional(),
  enabled: z.boolean(),
  triggers: z.array(triggerSchema), // Optional - can be empty if node uses conditions
  messages: z.array(messageSchema).min(1, 'At least one message is required'),
});

export type NodeEditFormData = z.infer<typeof nodeEditSchema>;

// Validate regex pattern
export function isValidRegex(pattern: string): boolean {
  try {
    const _regex = new RegExp(pattern);
    return Boolean(_regex);
  }
  catch {
    return false;
  }
}

// Validate trigger (including regex check)
export function validateTrigger(trigger: TriggerFormData): string | null {
  if (trigger.type === 'regex' && !isValidRegex(trigger.value)) {
    return 'Invalid regex pattern';
  }
  return null;
}
