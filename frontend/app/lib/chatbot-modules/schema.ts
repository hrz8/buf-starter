import type { ModuleMetadata, ModuleSchema } from './types';
import { buildModuleSchema } from './builder';

/**
 * Auto-discover all module metadata from module directories.
 * Each module must export `{moduleName}Metadata` from its index.ts
 */
const moduleExports = import.meta.glob<Record<string, ModuleMetadata>>(
  './*/index.ts',
  { eager: true },
);

/**
 * Auto-discover all JSON schemas generated from proto files.
 * Naming convention: chatbot.modules.v1.{PascalCaseModuleName}Config.jsonschema.strict.bundle.json
 */
const jsonSchemas = import.meta.glob<{ default: unknown }>(
  '../../../gen/jsonschema/chatbot.modules.v1.*.jsonschema.strict.bundle.json',
  { eager: true },
);

/**
 * Convert module key to PascalCase for proto message name lookup.
 * e.g., 'llm' -> 'Llm', 'mcpServer' -> 'McpServer'
 */
function toPascalCase(str: string): string {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

/**
 * Build MODULE_SCHEMAS dynamically by discovering modules and matching JSON schemas.
 */
function buildModuleSchemas(): Record<string, ModuleSchema> {
  const schemas: Record<string, ModuleSchema> = {};

  for (const [path, exports] of Object.entries(moduleExports)) {
    // Extract module name from path: './llm/index.ts' -> 'llm'
    const match = path.match(/^\.\/([^/]+)\/index\.ts$/);
    if (!match || !match[1])
      continue;

    const moduleName = match[1];
    const metadataKey = `${moduleName}Metadata`;
    const metadata = exports[metadataKey];

    if (!metadata)
      continue;

    // Find corresponding JSON schema by convention
    const protoName = `${toPascalCase(moduleName)}Config`;
    const jsonSchemaKey = Object.keys(jsonSchemas).find(k =>
      k.includes(`chatbot.modules.v1.${protoName}`),
    );

    if (jsonSchemaKey && jsonSchemas[jsonSchemaKey]) {
      const jsonSchema = jsonSchemas[jsonSchemaKey].default;
      if (jsonSchema) {
        schemas[moduleName] = buildModuleSchema(jsonSchema as never, metadata);
      }
    }
  }

  return schemas;
}

/**
 * Registry of all module schemas.
 * Auto-populated by discovering modules in subdirectories.
 */
export const MODULE_SCHEMAS: Record<string, ModuleSchema> = buildModuleSchemas();

/**
 * Type for valid module names (derived from MODULE_SCHEMAS keys).
 */
export type ModuleName = keyof typeof MODULE_SCHEMAS;

/**
 * Get a module schema by name.
 */
export function getModuleSchema(name: string): ModuleSchema | undefined {
  return MODULE_SCHEMAS[name];
}

/**
 * Get all module names.
 */
export function getModuleNames(): string[] {
  return Object.keys(MODULE_SCHEMAS);
}

/**
 * Check if a string is a valid module name.
 */
export function isValidModuleName(name: string): name is ModuleName {
  return name in MODULE_SCHEMAS;
}
