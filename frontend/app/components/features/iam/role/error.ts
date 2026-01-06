import { ConnectError } from '@connectrpc/connect';

/**
 * Type for validation errors - accepts ref-like objects with .value or direct records
 * Supports both mutable and readonly string arrays
 */
type ValidationErrors
  = | { readonly value: Record<string, readonly string[]> }
    | Record<string, readonly string[]>;

/**
 * Get the error record from validation errors (handles both ref and plain object)
 */
function getErrorRecord(validationErrors: ValidationErrors): Record<string, readonly string[]> {
  if ('value' in validationErrors && typeof validationErrors.value === 'object') {
    return validationErrors.value as Record<string, readonly string[]>;
  }
  return validationErrors as Record<string, readonly string[]>;
}

/**
 * Get ConnectRPC validation error for a specific field
 * Checks both direct field name and nested 'value.fieldName' format
 *
 * @param validationErrors - Validation errors from service
 * @param fieldName - Field name to get error for
 * @returns Error message string or empty string
 */
export function getConnectRPCError(
  validationErrors: ValidationErrors,
  fieldName: string,
): string {
  const errorObj = getErrorRecord(validationErrors);
  const errors = errorObj[fieldName] || errorObj[`value.${fieldName}`];
  return errors?.[0] || '';
}

/**
 * Check if a field has ConnectRPC validation errors
 *
 * @param validationErrors - Validation errors from service
 * @param fieldName - Field name to check
 * @returns True if field has errors
 */
export function hasConnectRPCError(
  validationErrors: ValidationErrors,
  fieldName: string,
): boolean {
  const errorObj = getErrorRecord(validationErrors);
  return !!(errorObj[fieldName] || errorObj[`value.${fieldName}`]);
}

/**
 * Extract error code from ConnectError details
 * @param error - ConnectError instance
 * @returns Error code string or null
 */
function extractErrorCode(error: ConnectError): string | null {
  // Try to extract from error details
  if (error.details && error.details.length > 0) {
    for (const detail of error.details) {
      const value = detail.value as any;
      if (value?.code) {
        return value.code;
      }
    }
  }

  // Fallback: Try to extract from message - find any 5-digit code followed by colon
  // Handles formats: "60601: Message", "[code] 60601: Message", "ConnectError: [code] 60601: Message"
  const match = error.message.match(/(\d{5}):/);
  if (match && match[1]) {
    return match[1];
  }

  return null;
}

/**
 * Get translated Connect error message using i18n
 * @param error - Error from Connect RPC call
 * @param t - i18n translate function
 * @returns Translated error message string
 */
export function getTranslatedConnectError(error: unknown, t: (key: string) => string): string {
  if (error instanceof ConnectError) {
    const errorCode = extractErrorCode(error);
    if (errorCode) {
      const translationKey = `errorCodes.${errorCode}`;
      const translated = t(translationKey);
      // If translation exists (not the same as the key), return it
      if (translated !== translationKey) {
        return translated;
      }
    }
    // Fallback to raw message
    return error.message;
  }
  return t('errorCodes.69901'); // Server Error
}

/**
 * Get generic Connect error message
 * @param error - Error from Connect RPC call
 * @returns Error message string
 * @deprecated Use getTranslatedConnectError instead for i18n support
 */
export function getGenericConnectError(error: unknown): string {
  if (error instanceof ConnectError) {
    return error.message;
  }
  return 'An unexpected error occurred';
}
