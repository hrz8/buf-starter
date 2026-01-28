export const ERROR_CODE = {
  NOT_FOUND: 'INF_3001',
  INTERNAL_SERVER_ERROR: 'INF_3002',
} as const;

export type ErrorCode = typeof ERROR_CODE[keyof typeof ERROR_CODE];
