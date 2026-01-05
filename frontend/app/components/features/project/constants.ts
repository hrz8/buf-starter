/**
 * Timezone options for project configuration
 * Common timezones sorted by geographical region
 */
export const TIMEZONE_OPTIONS = [
  // UTC
  'UTC',
  // Americas
  'America/New_York',
  'America/Chicago',
  'America/Denver',
  'America/Los_Angeles',
  'America/Toronto',
  'America/Mexico_City',
  'America/Sao_Paulo',
  // Europe
  'Europe/London',
  'Europe/Paris',
  'Europe/Berlin',
  'Europe/Moscow',
  // Asia
  'Asia/Dubai',
  'Asia/Kolkata',
  'Asia/Singapore',
  'Asia/Jakarta',
  'Asia/Tokyo',
  'Asia/Shanghai',
  'Asia/Hong_Kong',
  // Pacific
  'Pacific/Auckland',
  'Australia/Sydney',
] as const;

export type Timezone = typeof TIMEZONE_OPTIONS[number];

/**
 * Project environment types
 */
export const PROJECT_ENVIRONMENTS = {
  SANDBOX: 'sandbox',
  LIVE: 'live',
} as const;

export type ProjectEnvironment = typeof PROJECT_ENVIRONMENTS[keyof typeof PROJECT_ENVIRONMENTS];
