import type { Component } from 'vue';
import type { PropertySchema } from '@/lib/chatbot-modules';

import ArrayField from './ArrayField.vue';
import JsonField from './JsonField.vue';
import NumberField from './NumberField.vue';
import ObjectArrayField from './ObjectArrayField.vue';
import SelectField from './SelectField.vue';
import SwitchField from './SwitchField.vue';
import TextareaField from './TextareaField.vue';
import TextField from './TextField.vue';

export type FieldType
  = | 'text'
    | 'textarea'
    | 'number'
    | 'switch'
    | 'select'
    | 'json'
    | 'array'
    | 'objectArray';

// Field registry - maps type to component
const FIELD_REGISTRY: Record<FieldType, Component> = {
  text: TextField,
  textarea: TextareaField,
  number: NumberField,
  switch: SwitchField,
  select: SelectField,
  json: JsonField,
  array: ArrayField,
  objectArray: ObjectArrayField,
};

/**
 * Resolve a PropertySchema to its corresponding field type
 */
export function resolveFieldType(schema: PropertySchema): FieldType {
  // JSON editor takes priority
  if (schema.additionalTypeInfo === 'json') {
    return 'json';
  }

  // Array handling
  if (schema.type === 'array') {
    return schema.items?.type === 'object' ? 'objectArray' : 'array';
  }

  // Boolean
  if (schema.type === 'boolean') {
    return 'switch';
  }

  // Number
  if (schema.type === 'number') {
    return 'number';
  }

  // String variants
  if (schema.type === 'string') {
    if (schema.format === 'textarea') {
      return 'textarea';
    }
    if (schema.enum && schema.enum.length > 0) {
      return 'select';
    }
    return 'text';
  }

  // Default
  return 'text';
}

/**
 * Get the Vue component for a given PropertySchema
 */
export function resolveFieldComponent(schema: PropertySchema): Component {
  const fieldType = resolveFieldType(schema);
  return FIELD_REGISTRY[fieldType];
}

/**
 * Check if a schema represents a nested object (not array, not primitive)
 */
export function isNestedObject(schema: PropertySchema): boolean {
  return schema.type === 'object' && !!schema.properties;
}

// Export all field components for direct use if needed
export {
  ArrayField,
  JsonField,
  NumberField,
  ObjectArrayField,
  SelectField,
  SwitchField,
  TextareaField,
  TextField,
};
