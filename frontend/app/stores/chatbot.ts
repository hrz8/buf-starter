import type { JsonObject, MessageInitShape } from '@bufbuild/protobuf';
import { cleanForProtobuf, getProtobufErrorMessage } from '#shared/helpers/protobuf';
import { chatbotRepository } from '#shared/repository/chatbot';
import { create } from '@bufbuild/protobuf';
import {
  GetChatbotConfigRequestSchema,
  UpdateModuleConfigRequestSchema,
} from '~~/gen/altalune/v1/chatbot_pb';
import {
  getAllModuleDefaults,
  getModuleNames,
  mergeAllWithDefaults,
  mergeWithDefaults,
} from '@/lib/chatbot-modules';
import { useProjectStore } from './project';

export type ModulesConfig = Record<string, Record<string, unknown>>;

export const useChatbotStore = defineStore('chatbot', () => {
  const { $chatbotClient } = useNuxtApp();
  const chatbot = chatbotRepository($chatbotClient);
  const projectStore = useProjectStore();

  // State
  const rawConfig = ref<ModulesConfig>({}); // Config from backend (without defaults)
  const configId = ref<string | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const initialized = ref(false);
  const lastFetchedProjectId = ref<string | null>(null);

  // Computed: merged config (defaults + actual)
  const modulesConfig = computed<ModulesConfig>(() => {
    return mergeAllWithDefaults(rawConfig.value);
  });

  // Getters
  function isModuleEnabled(moduleName: string): boolean {
    return !!modulesConfig.value[moduleName]?.enabled;
  }

  function getModuleConfig(moduleName: string): Record<string, unknown> {
    return modulesConfig.value[moduleName] || mergeWithDefaults(moduleName, {});
  }

  function getRawModuleConfig(moduleName: string): Record<string, unknown> {
    return rawConfig.value[moduleName] || {};
  }

  // Actions
  async function fetchConfig(projectId: string): Promise<void> {
    // Skip if already fetched for this project
    if (initialized.value && lastFetchedProjectId.value === projectId) {
      return;
    }

    loading.value = true;
    error.value = null;

    try {
      const message = create(GetChatbotConfigRequestSchema, { projectId });
      const result = await chatbot.getChatbotConfig(message);

      if (result.chatbotConfig) {
        configId.value = result.chatbotConfig.id;
        // Store raw config from backend (may be partial/empty)
        rawConfig.value = (result.chatbotConfig.modulesConfig as unknown as ModulesConfig) || {};
      }
      else {
        // No config yet - start with empty (will use defaults)
        configId.value = null;
        rawConfig.value = {};
      }

      initialized.value = true;
      lastFetchedProjectId.value = projectId;
    }
    catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch chatbot config';
      throw err;
    }
    finally {
      loading.value = false;
    }
  }

  async function updateModuleConfig(
    projectId: string,
    moduleName: string,
    moduleConfig: Record<string, unknown>,
  ): Promise<void> {
    loading.value = true;
    error.value = null;

    // Clean config for protobuf serialization (removes undefined values)
    const cleanedConfig = cleanForProtobuf(moduleConfig);

    // Optimistically update local state
    const previousRawConfig = { ...rawConfig.value };
    rawConfig.value = {
      ...rawConfig.value,
      [moduleName]: cleanedConfig,
    };

    try {
      const req: MessageInitShape<typeof UpdateModuleConfigRequestSchema> = {
        projectId,
        moduleName,
        config: cleanedConfig as JsonObject,
      };

      const message = create(UpdateModuleConfigRequestSchema, req);
      const result = await chatbot.updateModuleConfig(message);

      if (result.chatbotConfig) {
        configId.value = result.chatbotConfig.id;
        // Update with server response to ensure consistency
        rawConfig.value = (result.chatbotConfig.modulesConfig as unknown as ModulesConfig) || {};
      }
    }
    catch (err) {
      // Revert on error
      rawConfig.value = previousRawConfig;

      // Parse error message for user-friendly display
      const rawMessage = err instanceof Error ? err.message : 'Failed to update module config';
      error.value = getProtobufErrorMessage(rawMessage);

      throw err;
    }
    finally {
      loading.value = false;
    }
  }

  // Reset store when project changes
  function reset(): void {
    rawConfig.value = {};
    configId.value = null;
    initialized.value = false;
    lastFetchedProjectId.value = null;
    error.value = null;
  }

  // Watch for project changes and reset
  watch(
    () => projectStore.activeProjectId,
    (newProjectId, oldProjectId) => {
      if (newProjectId !== oldProjectId) {
        reset();
      }
    },
  );

  // Ensure config is loaded for current project
  async function ensureLoaded(): Promise<void> {
    const projectId = projectStore.activeProjectId;
    if (!projectId) {
      return;
    }

    if (!initialized.value || lastFetchedProjectId.value !== projectId) {
      await fetchConfig(projectId);
    }
  }

  return {
    // State (readonly)
    modulesConfig: readonly(modulesConfig),
    configId: readonly(configId),
    loading: readonly(loading),
    error: readonly(error),
    initialized: readonly(initialized),

    // Getters
    isModuleEnabled,
    getModuleConfig,
    getRawModuleConfig,

    // Actions
    fetchConfig,
    updateModuleConfig,
    ensureLoaded,
    reset,

    // Utilities
    getModuleNames,
    getAllModuleDefaults,
  };
});
