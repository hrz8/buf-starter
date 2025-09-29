import type { MessageInitShape } from '@bufbuild/protobuf';
import type { ApiKey } from '~~/gen/altalune/v1/api_key_pb';

import type { QueryMetaResponseSchema } from '~~/gen/altalune/v1/common_pb';
import { apiKeyRepository } from '#shared/repository/api_key';

import { create } from '@bufbuild/protobuf';
import {
  ActivateApiKeyRequestSchema,
  CreateApiKeyRequestSchema,
  DeactivateApiKeyRequestSchema,
  DeleteApiKeyRequestSchema,
  GetApiKeyRequestSchema,
  QueryApiKeysRequestSchema,
  UpdateApiKeyRequestSchema,
} from '~~/gen/altalune/v1/api_key_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useApiKeyService() {
  const { $apiKeyClient } = useNuxtApp();
  const apiKey = apiKeyRepository($apiKeyClient);
  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryApiKeysRequestSchema);
  const createValidator = useConnectValidator(CreateApiKeyRequestSchema);
  const getValidator = useConnectValidator(GetApiKeyRequestSchema);
  const updateValidator = useConnectValidator(UpdateApiKeyRequestSchema);
  const deleteValidator = useConnectValidator(DeleteApiKeyRequestSchema);
  const activateValidator = useConnectValidator(ActivateApiKeyRequestSchema);
  const deactivateValidator = useConnectValidator(DeactivateApiKeyRequestSchema);

  // Create state for form submission
  const createState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // Update state for form submission
  const updateState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // Get state for fetching single API key
  const getState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // Delete state for confirmation and deletion
  const deleteState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // Activate state for toggling
  const activateState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // Deactivate state for toggling
  const deactivateState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  async function query(
    req: MessageInitShape<typeof QueryApiKeysRequestSchema>,
  ): Promise<{
    data: ApiKey[];
    meta: MessageInitShape<typeof QueryMetaResponseSchema> | undefined;
  }> {
    queryValidator.reset();
    if (!queryValidator.validate(req)) {
      console.warn('Validation failed for QueryApiKeysRequest:', queryValidator.errors.value);
      return {
        data: [],
        meta: {
          rowCount: 0,
          pageCount: 0,
          filters: {},
        },
      };
    }

    try {
      const message = create(QueryApiKeysRequestSchema, req);
      const result = await apiKey.queryApiKeys(message);
      return {
        data: result.data,
        meta: result.meta,
      };
    }
    catch (err) {
      const errorMessage = parseError(err);
      throw new Error(errorMessage);
    }
  }

  async function createApiKey(
    req: MessageInitShape<typeof CreateApiKeyRequestSchema>,
  ): Promise<{ apiKey: ApiKey | null; keyValue: string }> {
    createState.loading = true;
    createState.error = '';
    createState.success = false;

    createValidator.reset();

    if (!createValidator.validate(req)) {
      createState.loading = false;
      return { apiKey: null, keyValue: '' };
    }

    try {
      const message = create(CreateApiKeyRequestSchema, req);
      const result = await apiKey.createApiKey(message);
      createState.success = true;
      return {
        apiKey: result.apiKey || null,
        keyValue: result.keyValue || '',
      };
    }
    catch (err) {
      createState.error = parseError(err);
      throw new Error(createState.error);
    }
    finally {
      createState.loading = false;
    }
  }

  function resetCreateState() {
    createState.loading = false;
    createState.error = '';
    createState.success = false;
    createValidator.reset();
  }

  async function getApiKey(
    req: MessageInitShape<typeof GetApiKeyRequestSchema>,
  ): Promise<ApiKey | null> {
    getState.loading = true;
    getState.error = '';
    getState.success = false;

    getValidator.reset();

    if (!getValidator.validate(req)) {
      getState.loading = false;
      return null;
    }

    try {
      const message = create(GetApiKeyRequestSchema, req);
      const result = await apiKey.getApiKey(message);
      getState.success = true;
      return result.apiKey || null;
    }
    catch (err) {
      getState.error = parseError(err);
      throw new Error(getState.error);
    }
    finally {
      getState.loading = false;
    }
  }

  async function updateApiKey(
    req: MessageInitShape<typeof UpdateApiKeyRequestSchema>,
  ): Promise<ApiKey | null> {
    updateState.loading = true;
    updateState.error = '';
    updateState.success = false;

    updateValidator.reset();

    if (!updateValidator.validate(req)) {
      updateState.loading = false;
      return null;
    }

    try {
      const message = create(UpdateApiKeyRequestSchema, req);
      const result = await apiKey.updateApiKey(message);
      updateState.success = true;
      return result.apiKey || null;
    }
    catch (err) {
      updateState.error = parseError(err);
      throw new Error(updateState.error);
    }
    finally {
      updateState.loading = false;
    }
  }

  async function deleteApiKey(
    req: MessageInitShape<typeof DeleteApiKeyRequestSchema>,
  ): Promise<boolean> {
    deleteState.loading = true;
    deleteState.error = '';
    deleteState.success = false;

    deleteValidator.reset();

    if (!deleteValidator.validate(req)) {
      deleteState.loading = false;
      return false;
    }

    try {
      const message = create(DeleteApiKeyRequestSchema, req);
      await apiKey.deleteApiKey(message);
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

  function resetUpdateState() {
    updateState.loading = false;
    updateState.error = '';
    updateState.success = false;
    updateValidator.reset();
  }

  function resetGetState() {
    getState.loading = false;
    getState.error = '';
    getState.success = false;
    getValidator.reset();
  }

  function resetDeleteState() {
    deleteState.loading = false;
    deleteState.error = '';
    deleteState.success = false;
    deleteValidator.reset();
  }

  async function activateApiKey(
    req: MessageInitShape<typeof ActivateApiKeyRequestSchema>,
  ): Promise<ApiKey | null> {
    activateState.loading = true;
    activateState.error = '';
    activateState.success = false;

    activateValidator.reset();

    if (!activateValidator.validate(req)) {
      activateState.loading = false;
      return null;
    }

    try {
      const message = create(ActivateApiKeyRequestSchema, req);
      const result = await apiKey.activateApiKey(message);
      activateState.success = true;
      return result.apiKey || null;
    }
    catch (err) {
      activateState.error = parseError(err);
      throw new Error(activateState.error);
    }
    finally {
      activateState.loading = false;
    }
  }

  async function deactivateApiKey(
    req: MessageInitShape<typeof DeactivateApiKeyRequestSchema>,
  ): Promise<ApiKey | null> {
    deactivateState.loading = true;
    deactivateState.error = '';
    deactivateState.success = false;

    deactivateValidator.reset();

    if (!deactivateValidator.validate(req)) {
      deactivateState.loading = false;
      return null;
    }

    try {
      const message = create(DeactivateApiKeyRequestSchema, req);
      const result = await apiKey.deactivateApiKey(message);
      deactivateState.success = true;
      return result.apiKey || null;
    }
    catch (err) {
      deactivateState.error = parseError(err);
      throw new Error(deactivateState.error);
    }
    finally {
      deactivateState.loading = false;
    }
  }

  function resetActivateState() {
    activateState.loading = false;
    activateState.error = '';
    activateState.success = false;
    activateValidator.reset();
  }

  function resetDeactivateState() {
    deactivateState.loading = false;
    deactivateState.error = '';
    deactivateState.success = false;
    deactivateValidator.reset();
  }

  return {
    // Query
    query,
    queryValidationErrors: queryValidator.errors,

    // Create
    createApiKey,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    createValidationErrors: createValidator.errors,
    resetCreateState,

    // Get
    getApiKey,
    getLoading: computed(() => getState.loading),
    getError: computed(() => getState.error),
    getSuccess: computed(() => getState.success),
    getValidationErrors: getValidator.errors,
    resetGetState,

    // Update
    updateApiKey,
    updateLoading: computed(() => updateState.loading),
    updateError: computed(() => updateState.error),
    updateSuccess: computed(() => updateState.success),
    updateValidationErrors: updateValidator.errors,
    resetUpdateState,

    // Delete
    deleteApiKey,
    deleteLoading: computed(() => deleteState.loading),
    deleteError: computed(() => deleteState.error),
    deleteSuccess: computed(() => deleteState.success),
    deleteValidationErrors: deleteValidator.errors,
    resetDeleteState,

    // Activate
    activateApiKey,
    activateLoading: computed(() => activateState.loading),
    activateError: computed(() => activateState.error),
    activateSuccess: computed(() => activateState.success),
    activateValidationErrors: activateValidator.errors,
    resetActivateState,

    // Deactivate
    deactivateApiKey,
    deactivateLoading: computed(() => deactivateState.loading),
    deactivateError: computed(() => deactivateState.error),
    deactivateSuccess: computed(() => deactivateState.success),
    deactivateValidationErrors: deactivateValidator.errors,
    resetDeactivateState,
  };
}
