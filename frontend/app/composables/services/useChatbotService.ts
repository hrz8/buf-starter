import type { JsonObject, MessageInitShape } from '@bufbuild/protobuf';
import type { ChatbotConfig } from '~~/gen/altalune/v1/chatbot_pb';

import { chatbotRepository } from '#shared/repository/chatbot';

import { create } from '@bufbuild/protobuf';
import {
  GetChatbotConfigRequestSchema,
  UpdateModuleConfigRequestSchema,
} from '~~/gen/altalune/v1/chatbot_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useChatbotService() {
  const { $chatbotClient } = useNuxtApp();
  const chatbot = chatbotRepository($chatbotClient);
  const { parseError } = useErrorMessage();

  const getConfigValidator = useConnectValidator(GetChatbotConfigRequestSchema);
  const updateModuleValidator = useConnectValidator(UpdateModuleConfigRequestSchema);

  // Get config state
  const getState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // Update module state
  const updateState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  async function getChatbotConfig(
    req: MessageInitShape<typeof GetChatbotConfigRequestSchema>,
  ): Promise<ChatbotConfig | null> {
    getState.loading = true;
    getState.error = '';
    getState.success = false;

    getConfigValidator.reset();

    if (!getConfigValidator.validate(req)) {
      getState.loading = false;
      return null;
    }

    try {
      const message = create(GetChatbotConfigRequestSchema, req);
      const result = await chatbot.getChatbotConfig(message);
      getState.success = true;
      return result.chatbotConfig || null;
    }
    catch (err) {
      getState.error = parseError(err);
      throw new Error(getState.error);
    }
    finally {
      getState.loading = false;
    }
  }

  async function updateModuleConfig(
    projectId: string,
    moduleName: string,
    config: Record<string, unknown>,
  ): Promise<ChatbotConfig | null> {
    updateState.loading = true;
    updateState.error = '';
    updateState.success = false;

    updateModuleValidator.reset();

    const req: MessageInitShape<typeof UpdateModuleConfigRequestSchema> = {
      projectId,
      moduleName,
      config: config as JsonObject,
    };

    if (!updateModuleValidator.validate(req)) {
      updateState.loading = false;
      return null;
    }

    try {
      const message = create(UpdateModuleConfigRequestSchema, req);
      const result = await chatbot.updateModuleConfig(message);
      updateState.success = true;
      return result.chatbotConfig || null;
    }
    catch (err) {
      updateState.error = parseError(err);
      throw new Error(updateState.error);
    }
    finally {
      updateState.loading = false;
    }
  }

  function resetGetState() {
    getState.loading = false;
    getState.error = '';
    getState.success = false;
    getConfigValidator.reset();
  }

  function resetUpdateState() {
    updateState.loading = false;
    updateState.error = '';
    updateState.success = false;
    updateModuleValidator.reset();
  }

  return {
    // Get config
    getChatbotConfig,
    getLoading: computed(() => getState.loading),
    getError: computed(() => getState.error),
    getSuccess: computed(() => getState.success),
    getValidationErrors: getConfigValidator.errors,
    resetGetState,

    // Update module
    updateModuleConfig,
    updateLoading: computed(() => updateState.loading),
    updateError: computed(() => updateState.error),
    updateSuccess: computed(() => updateState.success),
    updateValidationErrors: updateModuleValidator.errors,
    resetUpdateState,
  };
}
