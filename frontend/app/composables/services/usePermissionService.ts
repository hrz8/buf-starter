import type { MessageInitShape } from '@bufbuild/protobuf';
import type { Permission } from '~~/gen/altalune/v1/permission_pb';
import { permissionRepository } from '#shared/repository/permission';
import { create } from '@bufbuild/protobuf';
import {
  CreatePermissionRequestSchema,
  DeletePermissionRequestSchema,
  GetPermissionRequestSchema,
  QueryPermissionsRequestSchema,
  UpdatePermissionRequestSchema,
} from '~~/gen/altalune/v1/permission_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function usePermissionService() {
  const { $permissionClient } = useNuxtApp();
  const permission = permissionRepository($permissionClient);
  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryPermissionsRequestSchema);
  const createValidator = useConnectValidator(CreatePermissionRequestSchema);
  const getValidator = useConnectValidator(GetPermissionRequestSchema);
  const updateValidator = useConnectValidator(UpdatePermissionRequestSchema);
  const deleteValidator = useConnectValidator(DeletePermissionRequestSchema);

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
    req: MessageInitShape<typeof QueryPermissionsRequestSchema>,
  ): Promise<{
    data: Permission[];
    meta: any;
  }> {
    queryValidator.reset();

    if (!queryValidator.validate(req)) {
      console.warn('Validation failed for QueryPermissionsRequest:', queryValidator.errors.value);
      return {
        data: [],
        meta: {
          rowCount: 0,
          pageCount: 0,
        },
      };
    }

    const message = create(QueryPermissionsRequestSchema, req);
    const result = await permission.queryPermissions(message);
    return {
      data: result.data,
      meta: result.meta,
    };
  }

  async function createPermission(
    req: MessageInitShape<typeof CreatePermissionRequestSchema>,
  ): Promise<Permission | null> {
    createState.loading = true;
    createState.error = '';
    createState.success = false;

    createValidator.reset();

    if (!createValidator.validate(req)) {
      createState.loading = false;
      return null;
    }

    try {
      const message = create(CreatePermissionRequestSchema, req);
      const result = await permission.createPermission(message);
      createState.success = true;
      return result.permission || null;
    }
    catch (err) {
      createState.error = parseError(err);
      throw err;
    }
    finally {
      createState.loading = false;
    }
  }

  async function getPermission(
    req: MessageInitShape<typeof GetPermissionRequestSchema>,
  ): Promise<Permission | null> {
    getState.loading = true;
    getState.error = '';
    getState.success = false;

    getValidator.reset();

    if (!getValidator.validate(req)) {
      getState.loading = false;
      return null;
    }

    try {
      const message = create(GetPermissionRequestSchema, req);
      const result = await permission.getPermission(message);
      getState.success = true;
      return result.permission || null;
    }
    catch (err) {
      getState.error = parseError(err);
      throw err;
    }
    finally {
      getState.loading = false;
    }
  }

  async function updatePermission(
    req: MessageInitShape<typeof UpdatePermissionRequestSchema>,
  ): Promise<Permission | null> {
    updateState.loading = true;
    updateState.error = '';
    updateState.success = false;

    updateValidator.reset();

    if (!updateValidator.validate(req)) {
      updateState.loading = false;
      return null;
    }

    try {
      const message = create(UpdatePermissionRequestSchema, req);
      const result = await permission.updatePermission(message);
      updateState.success = true;
      return result.permission || null;
    }
    catch (err) {
      updateState.error = parseError(err);
      throw err;
    }
    finally {
      updateState.loading = false;
    }
  }

  async function deletePermission(
    req: MessageInitShape<typeof DeletePermissionRequestSchema>,
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
      const message = create(DeletePermissionRequestSchema, req);
      await permission.deletePermission(message);
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
  }

  function resetUpdateState() {
    updateState.loading = false;
    updateState.error = '';
    updateState.success = false;
  }

  function resetGetState() {
    getState.loading = false;
    getState.error = '';
    getState.success = false;
  }

  function resetDeleteState() {
    deleteState.loading = false;
    deleteState.error = '';
    deleteState.success = false;
  }

  return {
    query,
    createPermission,
    getPermission,
    updatePermission,
    deletePermission,
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
    resetCreateState,
    resetUpdateState,
    resetGetState,
    resetDeleteState,
  };
}
