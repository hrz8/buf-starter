import type { MessageInitShape } from '@bufbuild/protobuf';
import type { Role } from '~~/gen/altalune/v1/role_pb';
import { roleRepository } from '#shared/repository/role';
import { create } from '@bufbuild/protobuf';
import {
  CreateRoleRequestSchema,
  DeleteRoleRequestSchema,
  GetRoleRequestSchema,
  QueryRolesRequestSchema,
  UpdateRoleRequestSchema,
} from '~~/gen/altalune/v1/role_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useRoleService() {
  const { $roleClient } = useNuxtApp();
  const role = roleRepository($roleClient);
  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryRolesRequestSchema);
  const createValidator = useConnectValidator(CreateRoleRequestSchema);
  const getValidator = useConnectValidator(GetRoleRequestSchema);
  const updateValidator = useConnectValidator(UpdateRoleRequestSchema);
  const deleteValidator = useConnectValidator(DeleteRoleRequestSchema);

  const createState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  const updateState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  const getState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  const deleteState = reactive({
    loading: false,
    error: '',
    success: false,
  });

  async function query(
    req: MessageInitShape<typeof QueryRolesRequestSchema>,
  ): Promise<{
    data: Role[];
    meta: any;
  }> {
    queryValidator.reset();

    if (!queryValidator.validate(req)) {
      console.warn('Validation failed for QueryRolesRequest:', queryValidator.errors.value);
      return {
        data: [],
        meta: {
          rowCount: 0,
          pageCount: 0,
        },
      };
    }

    const message = create(QueryRolesRequestSchema, req);
    const result = await role.queryRoles(message);
    return {
      data: result.data,
      meta: result.meta,
    };
  }

  async function createRole(
    req: MessageInitShape<typeof CreateRoleRequestSchema>,
  ): Promise<Role | null> {
    createState.loading = true;
    createState.error = '';
    createState.success = false;

    createValidator.reset();

    if (!createValidator.validate(req)) {
      createState.loading = false;
      return null;
    }

    try {
      const message = create(CreateRoleRequestSchema, req);
      const result = await role.createRole(message);
      createState.success = true;
      return result.role || null;
    }
    catch (err) {
      createState.error = parseError(err);
      throw err;
    }
    finally {
      createState.loading = false;
    }
  }

  async function getRole(
    req: MessageInitShape<typeof GetRoleRequestSchema>,
  ): Promise<Role | null> {
    getState.loading = true;
    getState.error = '';
    getState.success = false;

    getValidator.reset();

    if (!getValidator.validate(req)) {
      getState.loading = false;
      return null;
    }

    try {
      const message = create(GetRoleRequestSchema, req);
      const result = await role.getRole(message);
      getState.success = true;
      return result.role || null;
    }
    catch (err) {
      getState.error = parseError(err);
      throw err;
    }
    finally {
      getState.loading = false;
    }
  }

  async function updateRole(
    req: MessageInitShape<typeof UpdateRoleRequestSchema>,
  ): Promise<Role | null> {
    updateState.loading = true;
    updateState.error = '';
    updateState.success = false;

    updateValidator.reset();

    if (!updateValidator.validate(req)) {
      updateState.loading = false;
      return null;
    }

    try {
      const message = create(UpdateRoleRequestSchema, req);
      const result = await role.updateRole(message);
      updateState.success = true;
      return result.role || null;
    }
    catch (err) {
      updateState.error = parseError(err);
      throw err;
    }
    finally {
      updateState.loading = false;
    }
  }

  async function deleteRole(
    req: MessageInitShape<typeof DeleteRoleRequestSchema>,
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
      const message = create(DeleteRoleRequestSchema, req);
      await role.deleteRole(message);
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

  function resetCreateState() {
    createState.loading = false;
    createState.error = '';
    createState.success = false;
    createValidator.reset();
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

  return {
    query,
    createRole,
    getRole,
    updateRole,
    deleteRole,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    updateLoading: computed(() => updateState.loading),
    updateError: computed(() => updateState.error),
    updateSuccess: computed(() => updateState.success),
    getLoading: computed(() => getState.loading),
    getError: computed(() => getState.error),
    getSuccess: computed(() => getState.success),
    deleteLoading: computed(() => deleteState.loading),
    deleteError: computed(() => deleteState.error),
    deleteSuccess: computed(() => deleteState.success),
    queryValidationErrors: queryValidator.errors,
    createValidationErrors: createValidator.errors,
    updateValidationErrors: updateValidator.errors,
    getValidationErrors: getValidator.errors,
    deleteValidationErrors: deleteValidator.errors,
    resetCreateState,
    resetUpdateState,
    resetGetState,
    resetDeleteState,
  };
}
