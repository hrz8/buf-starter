// Trigger types for node activation
export const TRIGGER_TYPES = ['keyword', 'contains', 'regex', 'equals'] as const;
export type TriggerType = typeof TRIGGER_TYPES[number];

// Supported languages
export const LANGUAGES = [
  { value: 'en-US', label: 'English (US)' },
  { value: 'en-GB', label: 'English (UK)' },
  { value: 'id-ID', label: 'Bahasa Indonesia' },
  { value: 'ms-MY', label: 'Bahasa Melayu' },
] as const;

export type LanguageCode = typeof LANGUAGES[number]['value'];

// Default trigger for new nodes
export const DEFAULT_TRIGGER = {
  type: 'keyword' as TriggerType,
  value: '',
};

// Default message for new nodes
export const DEFAULT_MESSAGE = {
  role: 'assistant' as const,
  content: '',
};

// Trigger type labels for UI
export const TRIGGER_TYPE_LABELS: Record<TriggerType, string> = {
  keyword: 'Keyword',
  contains: 'Contains',
  regex: 'Regex',
  equals: 'Equals',
};

// Trigger type descriptions for UI
export const TRIGGER_TYPE_DESCRIPTIONS: Record<TriggerType, string> = {
  keyword: 'Match exact word or phrase',
  contains: 'Match if message contains text (case-insensitive)',
  regex: 'Match using regular expression pattern',
  equals: 'Exact match (case-sensitive)',
};

// Name pattern regex (lowercase_snake_case)
export const NODE_NAME_PATTERN = /^[a-z][a-z0-9_]*$/;

// Validation limits
export const VALIDATION_LIMITS = {
  nameMinLength: 2,
  nameMaxLength: 100,
  triggerValueMaxLength: 500,
  messageContentMaxLength: 5000,
  tagMaxLength: 50,
} as const;
