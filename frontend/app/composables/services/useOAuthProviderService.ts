import type { MessageInitShape } from '@bufbuild/protobuf';
import type { QueryMetaResponseSchema } from '~~/gen/altalune/v1/common_pb';
import type { OAuthProvider } from '~~/gen/altalune/v1/oauth_provider_pb';

import { oauthProviderRepository } from '#shared/repository/oauth_provider';
import { create } from '@bufbuild/protobuf';
import {
  CreateOAuthProviderRequestSchema,
  DeleteOAuthProviderRequestSchema,
  GetOAuthProviderRequestSchema,
  QueryOAuthProvidersRequestSchema,
  RevealClientSecretRequestSchema,
  UpdateOAuthProviderRequestSchema,
} from '~~/gen/altalune/v1/oauth_provider_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

/**
 * Reveal state for managing client secret visibility with auto-hide timer
 */
interface RevealState {
  providerId: string | null; // ID of provider with revealed secret
  clientSecret: string; // Plaintext client secret
  secondsRemaining: number; // Countdown timer (30 -> 0)
  loading: boolean;
  error: string;
}

export function useOAuthProviderService() {
  const { $oauthProviderClient } = useNuxtApp();
  const oauthProvider = oauthProviderRepository($oauthProviderClient);
  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryOAuthProvidersRequestSchema);
  const createValidator = useConnectValidator(CreateOAuthProviderRequestSchema);
  const getValidator = useConnectValidator(GetOAuthProviderRequestSchema);
  const updateValidator = useConnectValidator(UpdateOAuthProviderRequestSchema);
  const deleteValidator = useConnectValidator(DeleteOAuthProviderRequestSchema);
  const revealValidator = useConnectValidator(RevealClientSecretRequestSchema);

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

  // Get state for fetching single OAuth provider
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

  // Reveal state with auto-hide timer management
  const revealState = reactive<RevealState>({
    providerId: null,
    clientSecret: '',
    secondsRemaining: 0,
    loading: false,
    error: '',
  });

  // Timer references for cleanup
  let countdownInterval: ReturnType<typeof setInterval> | null = null;
  let autoHideTimeout: ReturnType<typeof setTimeout> | null = null;

  /**
   * Clear all reveal timers and reset reveal state
   */
  function clearRevealTimers() {
    if (countdownInterval) {
      clearInterval(countdownInterval);
      countdownInterval = null;
    }
    if (autoHideTimeout) {
      clearTimeout(autoHideTimeout);
      autoHideTimeout = null;
    }
  }

  /**
   * Hide revealed client secret and cleanup timers
   */
  function hideRevealedSecret() {
    clearRevealTimers();
    revealState.providerId = null;
    revealState.clientSecret = '';
    revealState.secondsRemaining = 0;
    revealState.error = '';
  }

  /**
   * Start 30-second auto-hide timer with countdown
   */
  function startRevealTimer() {
    clearRevealTimers(); // Clear any existing timers
    revealState.secondsRemaining = 30;

    // Countdown interval: decrements every second
    countdownInterval = setInterval(() => {
      if (revealState.secondsRemaining > 0) {
        revealState.secondsRemaining--;
      }
    }, 1000);

    // Auto-hide timeout: hides secret after 30 seconds
    autoHideTimeout = setTimeout(() => {
      hideRevealedSecret();
    }, 30000);
  }

  /**
   * Query OAuth providers with pagination
   */
  async function query(
    req: MessageInitShape<typeof QueryOAuthProvidersRequestSchema>,
  ): Promise<{
    data: OAuthProvider[];
    meta: MessageInitShape<typeof QueryMetaResponseSchema> | undefined;
  }> {
    queryValidator.reset();
    if (!queryValidator.validate(req)) {
      console.warn('Validation failed for QueryOAuthProvidersRequest:', queryValidator.errors.value);
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
      const message = create(QueryOAuthProvidersRequestSchema, req);
      const result = await oauthProvider.queryOAuthProviders(message);
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

  /**
   * Create a new OAuth provider (encrypts client secret)
   */
  async function createOAuthProvider(
    req: MessageInitShape<typeof CreateOAuthProviderRequestSchema>,
  ): Promise<OAuthProvider | null> {
    createState.loading = true;
    createState.error = '';
    createState.success = false;

    createValidator.reset();

    if (!createValidator.validate(req)) {
      createState.loading = false;
      return null;
    }

    try {
      const message = create(CreateOAuthProviderRequestSchema, req);
      const result = await oauthProvider.createOAuthProvider(message);
      createState.success = true;
      return result.provider || null;
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

  /**
   * Get a single OAuth provider by ID
   */
  async function getOAuthProvider(
    req: MessageInitShape<typeof GetOAuthProviderRequestSchema>,
  ): Promise<OAuthProvider | null> {
    getState.loading = true;
    getState.error = '';
    getState.success = false;

    getValidator.reset();

    if (!getValidator.validate(req)) {
      getState.loading = false;
      return null;
    }

    try {
      const message = create(GetOAuthProviderRequestSchema, req);
      const result = await oauthProvider.getOAuthProvider(message);
      getState.success = true;
      return result.provider || null;
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

  /**
   * Update an existing OAuth provider
   * If client_secret is empty, existing secret is retained
   * If client_secret is provided, it will be re-encrypted
   */
  async function updateOAuthProvider(
    req: MessageInitShape<typeof UpdateOAuthProviderRequestSchema>,
  ): Promise<OAuthProvider | null> {
    updateState.loading = true;
    updateState.error = '';
    updateState.success = false;

    updateValidator.reset();

    if (!updateValidator.validate(req)) {
      updateState.loading = false;
      return null;
    }

    try {
      const message = create(UpdateOAuthProviderRequestSchema, req);
      const result = await oauthProvider.updateOAuthProvider(message);
      updateState.success = true;
      return result.provider || null;
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

  /**
   * Delete an OAuth provider
   */
  async function deleteOAuthProvider(
    req: MessageInitShape<typeof DeleteOAuthProviderRequestSchema>,
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
      const message = create(DeleteOAuthProviderRequestSchema, req);
      await oauthProvider.deleteOAuthProvider(message);
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

  /**
   * Reveal client secret (decrypts and displays plaintext)
   * SECURITY: Starts 30-second auto-hide timer
   * After 30 seconds, secret is automatically hidden
   */
  async function revealClientSecret(
    req: MessageInitShape<typeof RevealClientSecretRequestSchema>,
  ): Promise<boolean> {
    revealState.loading = true;
    revealState.error = '';

    revealValidator.reset();

    if (!revealValidator.validate(req)) {
      revealState.loading = false;
      return false;
    }

    try {
      const message = create(RevealClientSecretRequestSchema, req);
      const result = await oauthProvider.revealClientSecret(message);

      // Set revealed secret and start 30-second timer
      revealState.providerId = req.id || null;
      revealState.clientSecret = result.clientSecret;
      startRevealTimer();

      return true;
    }
    catch (err) {
      revealState.error = parseError(err);
      throw new Error(revealState.error);
    }
    finally {
      revealState.loading = false;
    }
  }

  /**
   * Check if a provider's secret is currently revealed
   */
  function isSecretRevealed(providerId: string): boolean {
    return revealState.providerId === providerId && revealState.clientSecret !== '';
  }

  /**
   * Get revealed secret for a specific provider
   */
  function getRevealedSecret(providerId: string): string {
    if (isSecretRevealed(providerId)) {
      return revealState.clientSecret;
    }
    return '';
  }

  // Cleanup timers on component unmount
  onUnmounted(() => {
    clearRevealTimers();
  });

  return {
    // Query
    query,
    queryValidationErrors: queryValidator.errors,

    // Create
    createOAuthProvider,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    createValidationErrors: createValidator.errors,
    resetCreateState,

    // Get
    getOAuthProvider,
    getLoading: computed(() => getState.loading),
    getError: computed(() => getState.error),
    getSuccess: computed(() => getState.success),
    getValidationErrors: getValidator.errors,
    resetGetState,

    // Update
    updateOAuthProvider,
    updateLoading: computed(() => updateState.loading),
    updateError: computed(() => updateState.error),
    updateSuccess: computed(() => updateState.success),
    updateValidationErrors: updateValidator.errors,
    resetUpdateState,

    // Delete
    deleteOAuthProvider,
    deleteLoading: computed(() => deleteState.loading),
    deleteError: computed(() => deleteState.error),
    deleteSuccess: computed(() => deleteState.success),
    deleteValidationErrors: deleteValidator.errors,
    resetDeleteState,

    // Reveal (with 30-second auto-hide timer)
    revealClientSecret,
    revealLoading: computed(() => revealState.loading),
    revealError: computed(() => revealState.error),
    revealSecondsRemaining: computed(() => revealState.secondsRemaining),
    isSecretRevealed,
    getRevealedSecret,
    hideRevealedSecret,
    revealValidationErrors: revealValidator.errors,
  };
}
