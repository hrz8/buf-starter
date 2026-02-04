import type { MessageInitShape } from '@bufbuild/protobuf';
import type {
  ChatbotNode,
  ChatbotNodeMessage,
  ChatbotNodeTrigger,
  NodeCondition,
  NodeEffect,
  NodeNextAction,
} from '~~/gen/chatbot/nodes/v1/node_pb';

import { chatbotNodeRepository } from '#shared/repository/chatbot-node';

import { create } from '@bufbuild/protobuf';
import {
  CreateNodeRequestSchema,
  DeleteNodeRequestSchema,
  GetNodeRequestSchema,
  UpdateNodeRequestSchema,
} from '~~/gen/altalune/v1/chatbot_node_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useChatbotNodeService() {
  const { $chatbotNodeClient } = useNuxtApp();
  const nodeRepo = chatbotNodeRepository($chatbotNodeClient);
  const { parseError } = useErrorMessage();

  // Validators
  const createValidator = useConnectValidator(CreateNodeRequestSchema);
  const getValidator = useConnectValidator(GetNodeRequestSchema);
  const updateValidator = useConnectValidator(UpdateNodeRequestSchema);
  const deleteValidator = useConnectValidator(DeleteNodeRequestSchema);

  // Create node state
  const createState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // Get node state
  const getState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // Update node state
  const updateState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // Delete node state
  const deleteState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  async function createNode(
    projectId: string,
    name: string,
    lang: string,
    tags: string[] = [],
    version?: string,
  ): Promise<ChatbotNode | null> {
    createState.loading = true;
    createState.error = '';
    createState.success = false;
    createValidator.reset();

    const req: MessageInitShape<typeof CreateNodeRequestSchema> = {
      projectId,
      name,
      lang,
      tags,
      version,
    };

    if (!createValidator.validate(req)) {
      createState.loading = false;
      return null;
    }

    try {
      const message = create(CreateNodeRequestSchema, req);
      const result = await nodeRepo.createNode(message);
      createState.success = true;
      return result.node || null;
    }
    catch (err) {
      createState.error = parseError(err);
      throw new Error(createState.error);
    }
    finally {
      createState.loading = false;
    }
  }

  async function getNode(
    projectId: string,
    nodeId: string,
  ): Promise<ChatbotNode | null> {
    getState.loading = true;
    getState.error = '';
    getState.success = false;
    getValidator.reset();

    const req: MessageInitShape<typeof GetNodeRequestSchema> = {
      projectId,
      nodeId,
    };

    if (!getValidator.validate(req)) {
      getState.loading = false;
      return null;
    }

    try {
      const message = create(GetNodeRequestSchema, req);
      const result = await nodeRepo.getNode(message);
      getState.success = true;
      return result.node || null;
    }
    catch (err) {
      getState.error = parseError(err);
      throw new Error(getState.error);
    }
    finally {
      getState.loading = false;
    }
  }

  async function updateNode(
    projectId: string,
    nodeId: string,
    data: {
      name?: string;
      tags?: string[];
      enabled?: boolean;
      triggers?: ChatbotNodeTrigger[];
      messages?: ChatbotNodeMessage[];
      priority?: number;
      condition?: NodeCondition;
      effect?: NodeEffect;
      nextAction?: NodeNextAction;
    },
  ): Promise<ChatbotNode | null> {
    updateState.loading = true;
    updateState.error = '';
    updateState.success = false;
    updateValidator.reset();

    const req: MessageInitShape<typeof UpdateNodeRequestSchema> = {
      projectId,
      nodeId,
      name: data.name,
      tags: data.tags || [],
      enabled: data.enabled,
      triggers: data.triggers || [],
      messages: data.messages || [],
      priority: data.priority,
      condition: data.condition,
      effect: data.effect,
      nextAction: data.nextAction,
      // Set clear flags when the field is explicitly undefined/removed
      clearCondition: data.condition === undefined,
      clearEffect: data.effect === undefined,
      clearNextAction: data.nextAction === undefined,
    };

    if (!updateValidator.validate(req)) {
      updateState.loading = false;
      return null;
    }

    try {
      const message = create(UpdateNodeRequestSchema, req);
      const result = await nodeRepo.updateNode(message);
      updateState.success = true;
      return result.node || null;
    }
    catch (err) {
      updateState.error = parseError(err);
      throw new Error(updateState.error);
    }
    finally {
      updateState.loading = false;
    }
  }

  async function deleteNode(
    projectId: string,
    nodeId: string,
  ): Promise<boolean> {
    deleteState.loading = true;
    deleteState.error = '';
    deleteState.success = false;
    deleteValidator.reset();

    const req: MessageInitShape<typeof DeleteNodeRequestSchema> = {
      projectId,
      nodeId,
    };

    if (!deleteValidator.validate(req)) {
      deleteState.loading = false;
      return false;
    }

    try {
      const message = create(DeleteNodeRequestSchema, req);
      await nodeRepo.deleteNode(message);
      deleteState.success = true;
      return true;
    }
    catch (err) {
      deleteState.error = parseError(err);
      throw new Error(deleteState.error);
    }
    finally {
      deleteState.loading = false;
    }
  }

  // Reset functions
  function resetCreateState() {
    createState.loading = false;
    createState.error = '';
    createState.success = false;
    createValidator.reset();
  }

  function resetGetState() {
    getState.loading = false;
    getState.error = '';
    getState.success = false;
    getValidator.reset();
  }

  function resetUpdateState() {
    updateState.loading = false;
    updateState.error = '';
    updateState.success = false;
    updateValidator.reset();
  }

  function resetDeleteState() {
    deleteState.loading = false;
    deleteState.error = '';
    deleteState.success = false;
    deleteValidator.reset();
  }

  return {
    // Create
    createNode,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    createValidationErrors: createValidator.errors,
    resetCreateState,

    // Get
    getNode,
    getLoading: computed(() => getState.loading),
    getError: computed(() => getState.error),
    getSuccess: computed(() => getState.success),
    getValidationErrors: getValidator.errors,
    resetGetState,

    // Update
    updateNode,
    updateLoading: computed(() => updateState.loading),
    updateError: computed(() => updateState.error),
    updateSuccess: computed(() => updateState.success),
    updateValidationErrors: updateValidator.errors,
    resetUpdateState,

    // Delete
    deleteNode,
    deleteLoading: computed(() => deleteState.loading),
    deleteError: computed(() => deleteState.error),
    deleteSuccess: computed(() => deleteState.success),
    deleteValidationErrors: deleteValidator.errors,
    resetDeleteState,
  };
}
