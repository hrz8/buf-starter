import { QueryEmployeesRequestSchema, type Employee } from '~~/gen/altalune/v1/employee_pb';
import { type MessageInitShape, create } from '@bufbuild/protobuf';

import type { QueryMetaResponseSchema } from '~~/gen/altalune/v1/common_pb';

import { employeeRepository } from '#shared/repository/employee';
import { useConnectValidator } from '../useConnectValidator';
import { useErrorMessage } from '../useErrorMessage';

export function useEmployee() {
  const { $employeeClient } = useNuxtApp();
  const employee = employeeRepository($employeeClient);

  const { parseError } = useErrorMessage();

  const queryValidator = useConnectValidator(QueryEmployeesRequestSchema);

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

  return {
    // Query
    query,

    // Validation
    queryValidationErrors: queryValidator.errors,
  };
}
