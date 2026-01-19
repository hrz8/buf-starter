/**
 * Schema definition for a single property/field.
 * Used by form components to render appropriate inputs.
 *
 * This is the merged result of JSON Schema (from proto) and UI metadata.
 */
export interface PropertySchema {
  /**
   * The data type of the property.
   */
  type: 'string' | 'number' | 'boolean' | 'array' | 'object';

  /**
   * Display title for the field label.
   */
  title: string;

  /**
   * Help text or description for the field.
   */
  description?: string;

  /**
   * Default value for the field.
   */
  default?: unknown;

  /**
   * Placeholder text for input fields.
   */
  placeholder?: string;

  /**
   * Maximum string length (from proto validation).
   */
  maxLength?: number;

  /**
   * Minimum string length (from proto validation).
   */
  minLength?: number;

  /**
   * Format hint for string fields.
   * - 'textarea': Use multi-line textarea instead of single-line input
   * - 'json': Use JSON editor with syntax highlighting
   */
  format?: 'textarea' | 'json';

  /**
   * Additional type information for special handling.
   * - 'json': Field contains JSON string that should be pretty-printed
   */
  additionalTypeInfo?: 'json';

  /**
   * Minimum value for number fields (from proto validation).
   */
  minimum?: number;

  /**
   * Maximum value for number fields (from proto validation).
   */
  maximum?: number;

  /**
   * Step value for number inputs.
   * Controls the increment/decrement amount.
   */
  step?: number;

  /**
   * Enum options for select fields.
   */
  enum?: string[];

  /**
   * Human-readable labels for enum values.
   * Maps enum value to display label.
   */
  enumLabels?: Record<string, string>;

  /**
   * Schema for array items.
   */
  items?: PropertySchema;

  /**
   * Nested property schemas for object types.
   */
  properties?: Record<string, PropertySchema>;

  /**
   * List of required field names for object types.
   */
  required?: string[];

  /**
   * For arrays of objects: which field to use as the item title in collapsible headers.
   */
  titleKey?: string;
}

/**
 * Complete schema definition for a module.
 * Contains all information needed to render the module's configuration form.
 *
 * Built by merging JSON Schema (from proto) with UI metadata.
 */
export interface ModuleSchema {
  /**
   * Module key (must match the key used in MODULE_SCHEMAS).
   */
  key: string;

  /**
   * Display title for the module.
   */
  title: string;

  /**
   * Short description of what the module does.
   */
  description: string;

  /**
   * Icon name (from Lucide icons).
   */
  icon: string;

  /**
   * Always 'object' for module schemas.
   */
  type: 'object';

  /**
   * Property schemas for each field in the module.
   */
  properties: Record<string, PropertySchema>;

  /**
   * List of required field names.
   */
  required?: string[];
}

/**
 * UI-only metadata for a field.
 * These properties cannot be expressed in protobuf/JSON Schema.
 *
 * All properties are optional - the schema merger will:
 * - Auto-generate titles from field keys using toStartCase()
 * - Use description from JSON Schema (proto comments)
 * - Use validation constraints from JSON Schema (proto validation)
 */
export interface FieldMetadata {
  /**
   * Override the auto-generated title.
   * By default, titles are generated from field keys using toStartCase().
   * Only specify this for edge cases where auto-generation doesn't work well.
   */
  title?: string;

  /**
   * Placeholder text for input fields.
   * Not expressible in proto.
   */
  placeholder?: string;

  /**
   * Format hint for string fields.
   * - 'textarea': Use multi-line textarea instead of single-line input
   * - 'json': Use JSON editor with syntax highlighting
   */
  format?: 'textarea' | 'json';

  /**
   * Additional type information for special handling.
   * - 'json': Field contains JSON string that should be pretty-printed
   */
  additionalTypeInfo?: 'json';

  /**
   * Step value for number inputs.
   * Controls the increment/decrement amount.
   */
  step?: number;

  /**
   * Enum options for select fields.
   * Use when the proto doesn't define an enum but the field has fixed options.
   */
  enum?: string[];

  /**
   * Human-readable labels for enum values.
   * Maps enum value to display label.
   */
  enumLabels?: Record<string, string>;

  /**
   * For arrays of objects: which field to use as the item title in collapsible headers.
   */
  titleKey?: string;

  /**
   * Field ordering for nested objects.
   * List field keys in the order they should appear.
   */
  fieldOrder?: string[];

  /**
   * Nested field metadata for object types.
   */
  properties?: Record<string, FieldMetadata>;

  /**
   * Metadata for array items.
   */
  items?: FieldMetadata;
}

/**
 * Module-level UI metadata.
 * Provides display information for the module card/header.
 */
export interface ModuleMetadata {
  /**
   * Module key (must match the key used in MODULE_SCHEMAS).
   */
  key: string;

  /**
   * Display title for the module.
   */
  title: string;

  /**
   * Short description of what the module does.
   */
  description: string;

  /**
   * Icon name (from Lucide icons).
   */
  icon: string;

  /**
   * Field ordering for the form.
   * List field keys in the order they should appear.
   * Fields not listed will appear after ordered fields.
   */
  fieldOrder?: string[];

  /**
   * Field-level metadata.
   * Keys must match the field names in the proto message.
   */
  fields: Record<string, FieldMetadata>;
}
