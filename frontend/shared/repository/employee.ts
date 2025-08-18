import { ConnectError } from '@connectrpc/connect';

import type {
  QueryEmployeesResponse,
  QueryEmployeesRequest, EmployeeService,
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
});
