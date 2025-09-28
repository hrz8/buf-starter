import type { MessageInitShape } from '@bufbuild/protobuf';
import type { QueryMetaResponseSchema } from '~~/gen/altalune/v1/common_pb';

import type { Employee } from '~~/gen/altalune/v1/employee_pb';
import { employeeRepository } from '#shared/repository/employee';

import { create } from '@bufbuild/protobuf';
import {
  CreateEmployeeRequestSchema,
  DeleteEmployeeRequestSchema,

  GetEmployeeRequestSchema,
  QueryEmployeesRequestSchema,
  UpdateEmployeeRequestSchema,
} from '~~/gen/altalune/v1/employee_pb';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useEmployeeService() {
  const { $employeeClient } = useNuxtApp();
  const employee = employeeRepository($employeeClient);
  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryEmployeesRequestSchema);
  const createValidator = useConnectValidator(CreateEmployeeRequestSchema);
  const getValidator = useConnectValidator(GetEmployeeRequestSchema);
  const updateValidator = useConnectValidator(UpdateEmployeeRequestSchema);
  const deleteValidator = useConnectValidator(DeleteEmployeeRequestSchema);

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

  // Get state for fetching single employee
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

  async function query(
    req: MessageInitShape<typeof QueryEmployeesRequestSchema>,
  ): Promise<{
    data: Employee[];
    meta: MessageInitShape<typeof QueryMetaResponseSchema> | undefined;
  }> {
    queryValidator.reset();
    if (!queryValidator.validate(req)) {
      console.warn('Validation failed for QueryEmployeesRequest:', queryValidator.errors.value);
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
      const message = create(QueryEmployeesRequestSchema, req);
      const result = await employee.queryEmployees(message);
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

  async function createEmployee(
    req: MessageInitShape<typeof CreateEmployeeRequestSchema>,
  ): Promise<Employee | null> {
    createState.loading = true;
    createState.error = '';
    createState.success = false;

    createValidator.reset();

    if (!createValidator.validate(req)) {
      createState.loading = false;
      return null;
    }

    try {
      const message = create(CreateEmployeeRequestSchema, req);
      const result = await employee.createEmployee(message);
      createState.success = true;
      return result.employee || null;
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

  async function getEmployee(
    req: MessageInitShape<typeof GetEmployeeRequestSchema>,
  ): Promise<Employee | null> {
    getState.loading = true;
    getState.error = '';
    getState.success = false;

    getValidator.reset();

    if (!getValidator.validate(req)) {
      getState.loading = false;
      return null;
    }

    try {
      const message = create(GetEmployeeRequestSchema, req);
      const result = await employee.getEmployee(message);
      getState.success = true;
      return result.employee || null;
    }
    catch (err) {
      getState.error = parseError(err);
      throw new Error(getState.error);
    }
    finally {
      getState.loading = false;
    }
  }

  async function updateEmployee(
    req: MessageInitShape<typeof UpdateEmployeeRequestSchema>,
  ): Promise<Employee | null> {
    updateState.loading = true;
    updateState.error = '';
    updateState.success = false;

    updateValidator.reset();

    if (!updateValidator.validate(req)) {
      updateState.loading = false;
      return null;
    }

    try {
      const message = create(UpdateEmployeeRequestSchema, req);
      const result = await employee.updateEmployee(message);
      updateState.success = true;
      return result.employee || null;
    }
    catch (err) {
      updateState.error = parseError(err);
      throw new Error(updateState.error);
    }
    finally {
      updateState.loading = false;
    }
  }

  async function deleteEmployee(
    req: MessageInitShape<typeof DeleteEmployeeRequestSchema>,
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
      const message = create(DeleteEmployeeRequestSchema, req);
      await employee.deleteEmployee(message);
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

  return {
    // Query
    query,
    queryValidationErrors: queryValidator.errors,

    // Create
    createEmployee,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    createValidationErrors: createValidator.errors,
    resetCreateState,

    // Get
    getEmployee,
    getLoading: computed(() => getState.loading),
    getError: computed(() => getState.error),
    getSuccess: computed(() => getState.success),
    getValidationErrors: getValidator.errors,
    resetGetState,

    // Update
    updateEmployee,
    updateLoading: computed(() => updateState.loading),
    updateError: computed(() => updateState.error),
    updateSuccess: computed(() => updateState.success),
    updateValidationErrors: updateValidator.errors,
    resetUpdateState,

    // Delete
    deleteEmployee,
    deleteLoading: computed(() => deleteState.loading),
    deleteError: computed(() => deleteState.error),
    deleteSuccess: computed(() => deleteState.success),
    deleteValidationErrors: deleteValidator.errors,
    resetDeleteState,
  };
}
