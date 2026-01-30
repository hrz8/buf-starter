import type { Bot } from '../bot/bot.js';
import type { ModuleConfigMap, ModuleName } from '../bot/types.js';
import { setup as llmSetup } from './llm/index.js';
import { setup as mcpServerSetup } from './mcpServer/index.js';
import { setup as promptSetup } from './prompt/index.js';
import { setup as widgetSetup } from './widget/index.js';

export type ModuleSetupFn<K extends ModuleName = ModuleName> = (
  bot: Bot,
  config: ModuleConfigMap[K],
) => void | Promise<void>;

export interface ModuleRegistry {
  llm?: ModuleSetupFn<'llm'>;
  mcpServer?: ModuleSetupFn<'mcpServer'>;
  prompt?: ModuleSetupFn<'prompt'>;
  widget?: ModuleSetupFn<'widget'>;
}

export const moduleRegistry: ModuleRegistry = {
  llm: llmSetup,
  mcpServer: mcpServerSetup,
  prompt: promptSetup,
  widget: widgetSetup,
};

export function getModuleSetup<K extends ModuleName>(moduleName: K): ModuleSetupFn<K> | undefined {
  return moduleRegistry[moduleName] as ModuleSetupFn<K> | undefined;
}

export function hasModule(moduleName: ModuleName): boolean {
  return moduleName in moduleRegistry && moduleRegistry[moduleName] !== undefined;
}
