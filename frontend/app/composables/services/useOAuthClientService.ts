import type { MessageInitShape } from '@bufbuild/protobuf';
import type { QueryMetaResponse } from '~~/gen/altalune/v1/common_pb';

import type { OAuthClient } from '~~/gen/altalune/v1/oauth_client_pb';
import { oauthClientRepository } from '#shared/repository/oauth_client';

import { create } from '@bufbuild/protobuf';
import {
  CreateOAuthClientRequestSchema,
  DeleteOAuthClientRequestSchema,
  GetOAuthClientRequestSchema,
  QueryOAuthClientsRequestSchema,
  RevealOAuthClientSecretRequestSchema,
  UpdateOAuthClientRequestSchema,
} from '~~/gen/altalune/v1/oauth_client_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useOAuthClientService() {
  const { $oauthClientClient } = useNuxtApp();
  const oauthClient = oauthClientRepository($oauthClientClient);
  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryOAuthClientsRequestSchema);
  const createValidator = useConnectValidator(CreateOAuthClientRequestSchema);
  const getValidator = useConnectValidator(GetOAuthClientRequestSchema);
  const updateValidator = useConnectValidator(UpdateOAuthClientRequestSchema);
  const deleteValidator = useConnectValidator(DeleteOAuthClientRequestSchema);
  const revealValidator = useConnectValidator(RevealOAuthClientSecretRequestSchema);

  // Create state for form submission
  const createState = reactive({
    loading: false,
    error: '',
    success: false,
    clientSecret: '', // Store plaintext secret after creation
  });

  // Query state for list
  const queryState = reactive({
    loading: false,
    error: '',
  });

  // Get state for fetching single OAuth client
  const getState = reactive({
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

  // Delete state for confirmation and deletion
  const deleteState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  // Reveal state for revealing client secret
  const revealState = reactive({
    loading: false,
    error: '',
    clientSecret: '',
  });

  async function query(
    req: MessageInitShape<typeof QueryOAuthClientsRequestSchema>,
  ): Promise<{
    data: OAuthClient[];
    meta: QueryMetaResponse | undefined;
  }> {
    queryValidator.reset();
    queryState.loading = true;
    queryState.error = '';

    if (!queryValidator.validate(req)) {
      console.warn('Validation failed for QueryOAuthClientsRequest:', queryValidator.errors.value);
      queryState.loading = false;
      return {
        data: [],
        meta: undefined,
      };
    }

    try {
      const message = create(QueryOAuthClientsRequestSchema, req);
      const result = await oauthClient.queryOAuthClients(message);
      return {
        data: result.clients,
        meta: result.meta,
      };
    }
    catch (err) {
      const errorMessage = parseError(err);
      queryState.error = errorMessage;
      throw new Error(errorMessage);
    }
    finally {
      queryState.loading = false;
    }
  }

  async function createOAuthClient(
    req: MessageInitShape<typeof CreateOAuthClientRequestSchema>,
  ): Promise<{ client: OAuthClient | null; clientSecret: string }> {
    createState.loading = true;
    createState.error = '';
    createState.success = false;
    createState.clientSecret = '';

    createValidator.reset();

    if (!createValidator.validate(req)) {
      createState.loading = false;
      return { client: null, clientSecret: '' };
    }

    try {
      const message = create(CreateOAuthClientRequestSchema, req);
      const result = await oauthClient.createOAuthClient(message);
      createState.success = true;
      createState.clientSecret = result.clientSecret;
      return {
        client: result.client || null,
        clientSecret: result.clientSecret || '',
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
    createState.clientSecret = '';
    createValidator.reset();
  }

  async function getOAuthClient(
    req: MessageInitShape<typeof GetOAuthClientRequestSchema>,
  ): Promise<OAuthClient | null> {
    getState.loading = true;
    getState.error = '';
    getState.success = false;

    getValidator.reset();

    if (!getValidator.validate(req)) {
      getState.loading = false;
      return null;
    }

    try {
      const message = create(GetOAuthClientRequestSchema, req);
      const result = await oauthClient.getOAuthClient(message);
      getState.success = true;
      return result.client || null;
    }
    catch (err) {
      getState.error = parseError(err);
      throw new Error(getState.error);
    }
    finally {
      getState.loading = false;
    }
  }

  function resetGetState() {
    getState.loading = false;
    getState.error = '';
    getState.success = false;
    getValidator.reset();
  }

  async function updateOAuthClient(
    req: MessageInitShape<typeof UpdateOAuthClientRequestSchema>,
  ): Promise<OAuthClient | null> {
    updateState.loading = true;
    updateState.error = '';
    updateState.success = false;

    updateValidator.reset();

    if (!updateValidator.validate(req)) {
      updateState.loading = false;
      return null;
    }

    try {
      const message = create(UpdateOAuthClientRequestSchema, req);
      const result = await oauthClient.updateOAuthClient(message);
      updateState.success = true;
      return result.client || null;
    }
    catch (err) {
      updateState.error = parseError(err);
      throw new Error(updateState.error);
    }
    finally {
      updateState.loading = false;
    }
  }

  function resetUpdateState() {
    updateState.loading = false;
    updateState.error = '';
    updateState.success = false;
    updateValidator.reset();
  }

  async function deleteOAuthClient(
    req: MessageInitShape<typeof DeleteOAuthClientRequestSchema>,
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
      const message = create(DeleteOAuthClientRequestSchema, req);
      await oauthClient.deleteOAuthClient(message);
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

  function resetDeleteState() {
    deleteState.loading = false;
    deleteState.error = '';
    deleteState.success = false;
    deleteValidator.reset();
  }

  async function revealOAuthClientSecret(
    req: MessageInitShape<typeof RevealOAuthClientSecretRequestSchema>,
  ): Promise<string> {
    revealState.loading = true;
    revealState.error = '';
    revealState.clientSecret = '';

    revealValidator.reset();

    if (!revealValidator.validate(req)) {
      revealState.loading = false;
      return '';
    }

    try {
      const message = create(RevealOAuthClientSecretRequestSchema, req);
      const result = await oauthClient.revealOAuthClientSecret(message);
      revealState.clientSecret = result.clientSecret;
      return result.clientSecret;
    }
    catch (err) {
      revealState.error = parseError(err);
      throw new Error(revealState.error);
    }
    finally {
      revealState.loading = false;
    }
  }

  function resetRevealState() {
    revealState.loading = false;
    revealState.error = '';
    revealState.clientSecret = '';
    revealValidator.reset();
  }

  return {
    // Query
    query,
    queryLoading: computed(() => queryState.loading),
    queryError: computed(() => queryState.error),
    queryValidationErrors: queryValidator.errors,

    // Create
    createOAuthClient,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    clientSecret: computed(() => createState.clientSecret),
    createValidationErrors: createValidator.errors,
    resetCreateState,

    // Get
    getOAuthClient,
    getLoading: computed(() => getState.loading),
    getError: computed(() => getState.error),
    getSuccess: computed(() => getState.success),
    getValidationErrors: getValidator.errors,
    resetGetState,

    // Update
    updateOAuthClient,
    updateLoading: computed(() => updateState.loading),
    updateError: computed(() => updateState.error),
    updateSuccess: computed(() => updateState.success),
    updateValidationErrors: updateValidator.errors,
    resetUpdateState,

    // Delete
    deleteOAuthClient,
    deleteLoading: computed(() => deleteState.loading),
    deleteError: computed(() => deleteState.error),
    deleteSuccess: computed(() => deleteState.success),
    deleteValidationErrors: deleteValidator.errors,
    resetDeleteState,

    // Reveal
    revealOAuthClientSecret,
    revealLoading: computed(() => revealState.loading),
    revealError: computed(() => revealState.error),
    revealedSecret: computed(() => revealState.clientSecret),
    revealValidationErrors: revealValidator.errors,
    resetRevealState,
  };
}
