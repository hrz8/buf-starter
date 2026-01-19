/**
 * Auto-discover all module defaults from module directories.
 * Each module must export `{moduleName}Defaults` from its index.ts
 */
const moduleExports = import.meta.glob<Record<string, Record<string, unknown>>>(
  './*/index.ts',
  { eager: true },
);

/**
 * Build MODULE_DEFAULTS dynamically by discovering modules.
 */
function buildModuleDefaults(): Record<string, Record<string, unknown>> {
  const defaults: Record<string, Record<string, unknown>> = {};

  for (const [path, exports] of Object.entries(moduleExports)) {
    // Extract module name from path: './llm/index.ts' -> 'llm'
    const match = path.match(/^\.\/([^/]+)\/index\.ts$/);
    if (!match || !match[1])
      continue;

    const moduleName = match[1];
    const defaultsKey = `${moduleName}Defaults`;
    const moduleDefaults = exports[defaultsKey];

    if (moduleDefaults) {
      defaults[moduleName] = moduleDefaults;
    }
  }

  return defaults;
}

/**
 * Registry of default values for all modules.
 * Auto-populated by discovering modules in subdirectories.
 */
export const MODULE_DEFAULTS: Record<string, Record<string, unknown>>
  = buildModuleDefaults();

/**
 * Get default values for a specific module.
 * Returns empty object if module is not found.
 */
export function getModuleDefaults(moduleName: string): Record<string, unknown> {
  return MODULE_DEFAULTS[moduleName] ?? {};
}

/**
 * Get default values for all modules.
 */
export function getAllModuleDefaults(): Record<string, Record<string, unknown>> {
  return { ...MODULE_DEFAULTS };
}

/**
 * Deep merge two objects. Source values override target values.
 * Used to merge actual config over defaults.
 */
export function deepMerge<T extends Record<string, unknown>>(
  target: T,
  source: Partial<T>,
): T {
  const result = { ...target } as T;

  for (const [key, sourceValue] of Object.entries(source)) {
    if (sourceValue === undefined || sourceValue === null) {
      continue;
    }

    const targetValue = result[key as keyof T];

    // Deep merge objects (but not arrays)
    if (
      typeof sourceValue === 'object'
      && !Array.isArray(sourceValue)
      && typeof targetValue === 'object'
      && !Array.isArray(targetValue)
      && targetValue !== null
    ) {
      result[key as keyof T] = deepMerge(
        targetValue as Record<string, unknown>,
        sourceValue as Record<string, unknown>,
      ) as T[keyof T];
    }
    else {
      result[key as keyof T] = sourceValue as T[keyof T];
    }
  }

  return result;
}

/**
 * Merge module config with defaults.
 * Actual values take precedence over defaults.
 */
export function mergeWithDefaults(
  moduleName: string,
  actualConfig: Record<string, unknown>,
): Record<string, unknown> {
  const defaults = getModuleDefaults(moduleName);
  return deepMerge(defaults, actualConfig);
}

/**
 * Merge all modules config with their defaults.
 */
export function mergeAllWithDefaults(
  actualConfig: Record<string, Record<string, unknown>>,
): Record<string, Record<string, unknown>> {
  const allDefaults = getAllModuleDefaults();
  const result: Record<string, Record<string, unknown>> = {};

  // Start with all defaults
  for (const [moduleName, defaults] of Object.entries(allDefaults)) {
    const actual = actualConfig[moduleName] || {};
    result[moduleName] = deepMerge(defaults, actual);
  }

  return result;
}
