import {
  QueryEmployeesRequestSchema,
  CreateEmployeeRequestSchema,
  type Employee,
} from '~~/gen/altalune/v1/employee_pb';
import { type MessageInitShape, create } from '@bufbuild/protobuf';

import type { QueryMetaResponseSchema } from '~~/gen/altalune/v1/common_pb';

import { employeeRepository } from '#shared/repository/employee';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useEmployeeService() {
  const { $employeeClient } = useNuxtApp();
  const employee = employeeRepository($employeeClient);
  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryEmployeesRequestSchema);
  const createValidator = useConnectValidator(CreateEmployeeRequestSchema);

  // Create state for form submission
  const createState = reactive({
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
    } catch (err) {
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
    } catch (err) {
      createState.error = parseError(err);
      throw new Error(createState.error);
    } finally {
      createState.loading = false;
    }
  }

  function resetCreateState() {
    createState.loading = false;
    createState.error = '';
    createState.success = false;
    createValidator.reset();
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
  };
}
