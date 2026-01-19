// Schema builder (for advanced usage)
export { buildModuleSchema } from './builder';

// Defaults and merge utilities
export {
  deepMerge,
  getAllModuleDefaults,
  getModuleDefaults,
  mergeAllWithDefaults,
  mergeWithDefaults,
  MODULE_DEFAULTS,
} from './defaults';

// Schema registry and helpers
export {
  getModuleNames,
  getModuleSchema,
  isValidModuleName,
  MODULE_SCHEMAS,
} from './schema';

export type { ModuleName } from './schema';

// Types
export type {
  FieldMetadata,
  ModuleMetadata,
  ModuleSchema,
  PropertySchema,
} from './types';
