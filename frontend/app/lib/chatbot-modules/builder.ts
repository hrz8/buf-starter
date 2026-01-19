import type { FieldMetadata, ModuleMetadata, ModuleSchema, PropertySchema } from './types';
import { toStartCase } from '~~/shared/helpers/string';

/**
 * JSON Schema types (subset used by protoschema-plugins)
 */
interface JsonSchemaProperty {
  type?: string;
  description?: string;
  minimum?: number;
  maximum?: number;
  maxLength?: number;
  minLength?: number;
  pattern?: string;
  items?: JsonSchemaProperty | { $ref: string };
  properties?: Record<string, JsonSchemaProperty>;
  required?: string[];
  $ref?: string;
}

interface JsonSchemaBundle {
  $schema: string;
  $id: string;
  $ref: string;
  $defs: Record<string, JsonSchemaDefinition>;
}

interface JsonSchemaDefinition {
  $schema: string;
  title?: string;
  description?: string;
  type: string;
  properties?: Record<string, JsonSchemaProperty>;
  required?: string[];
  additionalProperties?: boolean;
}

/**
 * Resolve a $ref to its definition in the bundle
 */
function resolveRef(
  ref: string,
  defs: Record<string, JsonSchemaDefinition>,
): JsonSchemaDefinition | undefined {
  // $ref format: "#/$defs/chatbot.modules.v1.McpServerUrl.jsonschema.strict.json"
  const match = ref.match(/^#\/\$defs\/(.+)$/);
  if (!match || !match[1])
    return undefined;
  return defs[match[1]];
}

/**
 * Map JSON Schema type to PropertySchema type.
 * Also considers $ref as indicating an object type.
 */
function mapJsonSchemaType(
  jsonProp: JsonSchemaProperty,
): PropertySchema['type'] {
  // If property has $ref, it's referencing another definition (object)
  if (jsonProp.$ref) {
    return 'object';
  }

  switch (jsonProp.type) {
    case 'boolean':
      return 'boolean';
    case 'integer':
    case 'number':
      return 'number';
    case 'array':
      return 'array';
    case 'object':
      return 'object';
    case 'string':
    default:
      return 'string';
  }
}

/**
 * Transform a JSON Schema property to PropertySchema
 */
function transformProperty(
  key: string,
  jsonProp: JsonSchemaProperty,
  defs: Record<string, JsonSchemaDefinition>,
  fieldMetadata?: FieldMetadata,
): PropertySchema {
  const type = mapJsonSchemaType(jsonProp);

  const prop: PropertySchema = {
    type,
    // Auto-generate title from key, allow metadata override
    title: fieldMetadata?.title ?? toStartCase(key),
    // Use description from JSON Schema (proto comments)
    description: jsonProp.description,
  };

  // Add metadata properties (defaults are handled separately via defaults.ts)
  if (fieldMetadata?.placeholder) {
    prop.placeholder = fieldMetadata.placeholder;
  }
  if (fieldMetadata?.format) {
    prop.format = fieldMetadata.format;
  }
  if (fieldMetadata?.additionalTypeInfo) {
    prop.additionalTypeInfo = fieldMetadata.additionalTypeInfo;
  }
  if (fieldMetadata?.step) {
    prop.step = fieldMetadata.step;
  }
  if (fieldMetadata?.enum) {
    prop.enum = fieldMetadata.enum;
  }
  if (fieldMetadata?.enumLabels) {
    prop.enumLabels = fieldMetadata.enumLabels;
  }
  if (fieldMetadata?.titleKey) {
    prop.titleKey = fieldMetadata.titleKey;
  }

  // Add JSON Schema validation constraints
  if (jsonProp.minimum !== undefined) {
    prop.minimum = jsonProp.minimum;
  }
  if (jsonProp.maximum !== undefined) {
    prop.maximum = jsonProp.maximum;
  }
  if (jsonProp.maxLength !== undefined) {
    prop.maxLength = jsonProp.maxLength;
  }
  if (jsonProp.minLength !== undefined) {
    prop.minLength = jsonProp.minLength;
  }

  // Handle array items
  if (type === 'array' && jsonProp.items) {
    const itemsMetadata = fieldMetadata?.items;

    if ('$ref' in jsonProp.items && jsonProp.items.$ref) {
      // Items reference another definition
      const refDef = resolveRef(jsonProp.items.$ref, defs);
      if (refDef && refDef.properties) {
        prop.items = {
          type: 'object',
          title: fieldMetadata?.items?.title ?? toStartCase(key.replace(/s$/, '')),
          titleKey: itemsMetadata?.titleKey ?? fieldMetadata?.titleKey,
          properties: transformProperties(
            refDef.properties,
            defs,
            itemsMetadata?.properties,
            itemsMetadata?.fieldOrder,
          ),
          required: refDef.required,
        };
      }
    }
    else {
      // Inline items definition
      const itemsProp = jsonProp.items as JsonSchemaProperty;
      if (itemsProp.type === 'object' && itemsProp.properties) {
        prop.items = {
          type: 'object',
          title: itemsMetadata?.title ?? toStartCase(key.replace(/s$/, '')),
          titleKey: itemsMetadata?.titleKey ?? fieldMetadata?.titleKey,
          properties: transformProperties(
            itemsProp.properties,
            defs,
            itemsMetadata?.properties,
            itemsMetadata?.fieldOrder,
          ),
        };
      }
      else {
        // Simple array items (e.g., string[])
        prop.items = {
          type: mapJsonSchemaType(itemsProp),
          title: itemsMetadata?.title ?? toStartCase(key.replace(/s$/, '')),
          placeholder: itemsMetadata?.placeholder,
        };
      }
    }
  }

  // Handle nested objects
  if (type === 'object') {
    if (jsonProp.$ref) {
      // Object references another definition
      const refDef = resolveRef(jsonProp.$ref, defs);
      if (refDef && refDef.properties) {
        prop.properties = transformProperties(
          refDef.properties,
          defs,
          fieldMetadata?.properties,
          fieldMetadata?.fieldOrder,
        );
        prop.required = refDef.required;
      }
    }
    else if (jsonProp.properties) {
      // Inline object definition
      prop.properties = transformProperties(
        jsonProp.properties,
        defs,
        fieldMetadata?.properties,
        fieldMetadata?.fieldOrder,
      );
      prop.required = jsonProp.required;
    }
  }

  return prop;
}

/**
 * Transform JSON Schema properties to PropertySchema record.
 * Properties are ordered based on fieldOrder if provided.
 */
function transformProperties(
  jsonProps: Record<string, JsonSchemaProperty>,
  defs: Record<string, JsonSchemaDefinition>,
  fieldsMetadata?: Record<string, FieldMetadata>,
  fieldOrder?: string[],
): Record<string, PropertySchema> {
  const result: Record<string, PropertySchema> = {};
  const jsonKeys = Object.keys(jsonProps);

  // Determine field order: use explicit order if provided, otherwise use JSON Schema order
  let orderedKeys: string[];
  if (fieldOrder && fieldOrder.length > 0) {
    // Start with explicitly ordered fields
    const orderedSet = new Set(fieldOrder);
    orderedKeys = [...fieldOrder.filter(key => jsonKeys.includes(key))];
    // Add remaining fields not in explicit order
    for (const key of jsonKeys) {
      if (!orderedSet.has(key)) {
        orderedKeys.push(key);
      }
    }
  }
  else {
    orderedKeys = jsonKeys;
  }

  for (const key of orderedKeys) {
    const jsonProp = jsonProps[key];
    if (jsonProp) {
      result[key] = transformProperty(
        key,
        jsonProp,
        defs,
        fieldsMetadata?.[key],
      );
    }
  }

  return result;
}

/**
 * Build a ModuleSchema from JSON Schema bundle and UI metadata
 */
export function buildModuleSchema(
  jsonSchemaBundle: JsonSchemaBundle,
  metadata: ModuleMetadata,
): ModuleSchema {
  // Get the main definition from $ref
  const mainRefMatch = jsonSchemaBundle.$ref.match(/^#\/\$defs\/(.+)$/);
  if (!mainRefMatch || !mainRefMatch[1]) {
    throw new Error(`Invalid $ref in JSON Schema bundle: ${jsonSchemaBundle.$ref}`);
  }

  const mainDef = jsonSchemaBundle.$defs[mainRefMatch[1]];
  if (!mainDef || !mainDef.properties) {
    throw new Error(`Main definition not found or has no properties: ${mainRefMatch[1]}`);
  }

  return {
    key: metadata.key,
    title: metadata.title,
    description: metadata.description,
    icon: metadata.icon,
    type: 'object',
    properties: transformProperties(
      mainDef.properties,
      jsonSchemaBundle.$defs,
      metadata.fields,
      metadata.fieldOrder,
    ),
    required: mainDef.required,
  };
}
