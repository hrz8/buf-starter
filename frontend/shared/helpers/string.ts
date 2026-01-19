/**
 * Common acronyms that should be fully uppercased in titles.
 * Add new acronyms here as needed.
 */
const ACRONYMS = new Set([
  'api',
  'cors',
  'cpu',
  'css',
  'db',
  'dns',
  'html',
  'http',
  'https',
  'id',
  'io',
  'ip',
  'json',
  'jwt',
  'llm',
  'mcp',
  'oauth',
  'ram',
  'rest',
  'rpc',
  'sdk',
  'sql',
  'ssh',
  'ssl',
  'tls',
  'ui',
  'uri',
  'url',
  'urls',
  'xml',
]);

/**
 * Converts a camelCase or PascalCase string to Start Case with proper acronym handling.
 *
 * @example
 * toStartCase('apiKey')        // 'API Key'
 * toStartCase('maxToolCalls')  // 'Max Tool Calls'
 * toStartCase('cors')          // 'CORS'
 * toStartCase('systemPrompt')  // 'System Prompt'
 * toStartCase('isEnabled')     // 'Is Enabled'
 * toStartCase('httpUrl')       // 'HTTP URL'
 *
 * @param str - The camelCase or PascalCase string to convert
 * @returns The string converted to Start Case with acronyms uppercased
 */
export function toStartCase(str: string): string {
  if (!str)
    return '';

  // Split on uppercase letters, keeping the delimiter
  // 'apiKey' -> ['api', 'Key'] -> ['api', 'key']
  // 'maxToolCalls' -> ['max', 'Tool', 'Calls'] -> ['max', 'tool', 'calls']
  const words = str
    .replace(/([A-Z])/g, ' $1')
    .trim()
    .toLowerCase()
    .split(/\s+/);

  return words
    .map((word) => {
      if (ACRONYMS.has(word)) {
        return word.toUpperCase();
      }
      return word.charAt(0).toUpperCase() + word.slice(1);
    })
    .join(' ');
}

/**
 * Converts a snake_case string to Start Case with proper acronym handling.
 *
 * @example
 * snakeToStartCase('api_key')        // 'API Key'
 * snakeToStartCase('max_tool_calls') // 'Max Tool Calls'
 * snakeToStartCase('system_prompt')  // 'System Prompt'
 *
 * @param str - The snake_case string to convert
 * @returns The string converted to Start Case with acronyms uppercased
 */
export function snakeToStartCase(str: string): string {
  if (!str)
    return '';

  const words = str.toLowerCase().split('_');

  return words
    .map((word) => {
      if (ACRONYMS.has(word)) {
        return word.toUpperCase();
      }
      return word.charAt(0).toUpperCase() + word.slice(1);
    })
    .join(' ');
}

/**
 * Converts a kebab-case string to Start Case with proper acronym handling.
 *
 * @example
 * kebabToStartCase('api-key')        // 'API Key'
 * kebabToStartCase('max-tool-calls') // 'Max Tool Calls'
 *
 * @param str - The kebab-case string to convert
 * @returns The string converted to Start Case with acronyms uppercased
 */
export function kebabToStartCase(str: string): string {
  if (!str)
    return '';

  const words = str.toLowerCase().split('-');

  return words
    .map((word) => {
      if (ACRONYMS.has(word)) {
        return word.toUpperCase();
      }
      return word.charAt(0).toUpperCase() + word.slice(1);
    })
    .join(' ');
}
