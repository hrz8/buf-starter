import type { MessageInitShape } from '@bufbuild/protobuf';
import type { User } from '~~/gen/altalune/v1/user_pb';
import { userRepository } from '#shared/repository/user';
import { create } from '@bufbuild/protobuf';
import {
  ActivateUserRequestSchema,
  CreateUserRequestSchema,
  DeactivateUserRequestSchema,
  DeleteUserRequestSchema,
  GetUserRequestSchema,
  QueryUsersRequestSchema,
  UpdateUserRequestSchema,
} from '~~/gen/altalune/v1/user_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useUserService() {
  const { $userClient } = useNuxtApp();
  const user = userRepository($userClient);
  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryUsersRequestSchema);
  const createValidator = useConnectValidator(CreateUserRequestSchema);
  const getValidator = useConnectValidator(GetUserRequestSchema);
  const updateValidator = useConnectValidator(UpdateUserRequestSchema);
  const deleteValidator = useConnectValidator(DeleteUserRequestSchema);
  const activateValidator = useConnectValidator(ActivateUserRequestSchema);
  const deactivateValidator = useConnectValidator(DeactivateUserRequestSchema);

  const createState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  const deleteState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  const getState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  const updateState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  function resetCreateState() {
    createState.loading = false;
    createState.error = '';
    createState.success = false;
    createValidator.reset();
  }

  function resetDeleteState() {
    deleteState.loading = false;
    deleteState.error = '';
    deleteState.success = false;
    deleteValidator.reset();
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

  async function query(
    req: MessageInitShape<typeof QueryUsersRequestSchema>,
  ): Promise<{
    data: User[];
    meta: any;
  }> {
    queryValidator.reset();

    if (!queryValidator.validate(req)) {
      console.warn('Validation failed for QueryUsersRequest:', queryValidator.errors.value);
      return {
        data: [],
        meta: {
          rowCount: 0,
          pageCount: 0,
        },
      };
    }

    const message = create(QueryUsersRequestSchema, req);
    const result = await user.queryUsers(message);
    return {
      data: result.data,
      meta: result.meta,
    };
  }

  async function createUser(
    req: MessageInitShape<typeof CreateUserRequestSchema>,
  ): Promise<User | null> {
    createState.loading = true;
    createState.error = '';
    createState.success = false;

    createValidator.reset();

    if (!createValidator.validate(req)) {
      createState.loading = false;
      return null;
    }

    try {
      const message = create(CreateUserRequestSchema, req);
      const result = await user.createUser(message);
      createState.success = true;
      return result.user || null;
    }
    catch (err) {
      createState.error = parseError(err);
      throw err;
    }
    finally {
      createState.loading = false;
    }
  }

  async function getUser(
    req: MessageInitShape<typeof GetUserRequestSchema>,
  ): Promise<User | null> {
    getState.loading = true;
    getState.error = '';
    getState.success = false;

    getValidator.reset();

    if (!getValidator.validate(req)) {
      getState.loading = false;
      return null;
    }

    try {
      const message = create(GetUserRequestSchema, req);
      const result = await user.getUser(message);
      getState.success = true;
      return result.user || null;
    }
    catch (err) {
      getState.error = parseError(err);
      throw err;
    }
    finally {
      getState.loading = false;
    }
  }

  async function updateUser(
    req: MessageInitShape<typeof UpdateUserRequestSchema>,
  ): Promise<User | null> {
    updateState.loading = true;
    updateState.error = '';
    updateState.success = false;

    updateValidator.reset();

    if (!updateValidator.validate(req)) {
      updateState.loading = false;
      return null;
    }

    try {
      const message = create(UpdateUserRequestSchema, req);
      const result = await user.updateUser(message);
      updateState.success = true;
      return result.user || null;
    }
    catch (err) {
      updateState.error = parseError(err);
      throw err;
    }
    finally {
      updateState.loading = false;
    }
  }

  async function deleteUser(
    req: MessageInitShape<typeof DeleteUserRequestSchema>,
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
      const message = create(DeleteUserRequestSchema, req);
      await user.deleteUser(message);
      deleteState.success = true;
      return true;
    }
    catch (err) {
      deleteState.error = parseError(err);
      throw err;
    }
    finally {
      deleteState.loading = false;
    }
  }

  async function activateUser(
    req: MessageInitShape<typeof ActivateUserRequestSchema>,
  ): Promise<User | null> {
    activateValidator.reset();

    if (!activateValidator.validate(req)) {
      return null;
    }

    const message = create(ActivateUserRequestSchema, req);
    const result = await user.activateUser(message);
    return result.user || null;
  }

  async function deactivateUser(
    req: MessageInitShape<typeof DeactivateUserRequestSchema>,
  ): Promise<User | null> {
    deactivateValidator.reset();

    if (!deactivateValidator.validate(req)) {
      return null;
    }

    const message = create(DeactivateUserRequestSchema, req);
    const result = await user.deactivateUser(message);
    return result.user || null;
  }

  return {
    query,
    createUser,
    getUser,
    updateUser,
    deleteUser,
    activateUser,
    deactivateUser,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    deleteLoading: computed(() => deleteState.loading),
    deleteError: computed(() => deleteState.error),
    deleteSuccess: computed(() => deleteState.success),
    getLoading: computed(() => getState.loading),
    getError: computed(() => getState.error),
    getSuccess: computed(() => getState.success),
    updateLoading: computed(() => updateState.loading),
    updateError: computed(() => updateState.error),
    updateSuccess: computed(() => updateState.success),
    queryValidationErrors: queryValidator.errors,
    createValidationErrors: createValidator.errors,
    updateValidationErrors: updateValidator.errors,
    resetCreateState,
    resetDeleteState,
    resetGetState,
    resetUpdateState,
  };
}
