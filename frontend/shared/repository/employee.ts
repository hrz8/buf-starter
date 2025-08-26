import { ConnectError } from '@connectrpc/connect';

import type {
  QueryEmployeesResponse,
  CreateEmployeeResponse,
  QueryEmployeesRequest,
  CreateEmployeeRequest,
  EmployeeService,
} from '~~/gen/altalune/v1/employee_pb';
import type { Client } from '@connectrpc/connect';

export const employeeRepository = (client: Client<typeof EmployeeService>) => ({
  async queryEmployees(req: QueryEmployeesRequest): Promise<QueryEmployeesResponse> {
    try {
      const response = await client.queryEmployees(req);
      return response;
    } catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err);
      }
      throw err;
    }
  },

  async createEmployee(req: CreateEmployeeRequest): Promise<CreateEmployeeResponse> {
    try {
      const response = await client.createEmployee(req);
      return response;
    } catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err);
      }
      throw err;
    }
  },
});
