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
