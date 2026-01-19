/**
 * Utility functions for working with Protocol Buffers, particularly
 * google.protobuf.Struct and google.protobuf.Value types.
 */

/**
 * Clean an object for protobuf serialization.
 *
 * Protobuf's google.protobuf.Value cannot represent `undefined` values.
 * This function recursively removes undefined values and ensures
 * the object can be safely serialized.
 *
 * @param obj - The object to clean
 * @returns A new object safe for protobuf serialization
 */
export function cleanForProtobuf<T extends Record<string, unknown>>(
  obj: T,
): Record<string, unknown> {
  const result: Record<string, unknown> = {};

  for (const [key, value] of Object.entries(obj)) {
    const cleaned = cleanValue(value);
    // Only include if the cleaned value is not undefined
    if (cleaned !== undefined) {
      result[key] = cleaned;
    }
  }

  return result;
}

/**
 * Clean a single value for protobuf serialization.
 * Returns undefined for values that should be omitted.
 */
function cleanValue(value: unknown): unknown {
  // Remove undefined values
  if (value === undefined) {
    return undefined;
  }

  // Keep null as-is (protobuf can handle null)
  if (value === null) {
    return null;
  }

  // Handle arrays
  if (Array.isArray(value)) {
    return value
      .map(item => cleanValue(item))
      .filter(item => item !== undefined);
  }

  // Handle nested objects
  if (typeof value === 'object') {
    return cleanForProtobuf(value as Record<string, unknown>);
  }

  // Primitives (string, number, boolean) pass through
  return value;
}

/**
 * Common protobuf error messages mapped to user-friendly messages.
 */
const PROTOBUF_ERROR_MESSAGES: Record<string, string> = {
  'google.protobuf.Value must have a value': 'Please fill in all required fields before saving.',
  'invalid wire type': 'Invalid data format. Please refresh and try again.',
};

/**
 * Check if an error message is a known protobuf error and return a user-friendly message.
 *
 * @param errorMessage - The raw error message
 * @returns User-friendly message if known, otherwise the original message
 */
export function getProtobufErrorMessage(errorMessage: string): string {
  for (const [pattern, friendlyMessage] of Object.entries(PROTOBUF_ERROR_MESSAGES)) {
    if (errorMessage.includes(pattern)) {
      return friendlyMessage;
    }
  }
  return errorMessage;
}
