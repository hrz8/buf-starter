import type { Client } from '@connectrpc/connect';

import type {
  CreateEmployeeRequest,
  CreateEmployeeResponse,
  DeleteEmployeeRequest,
  DeleteEmployeeResponse,
  EmployeeService,
  GetEmployeeRequest,
  GetEmployeeResponse,
  QueryEmployeesRequest,
  QueryEmployeesResponse,
  UpdateEmployeeRequest,
  UpdateEmployeeResponse,
} from '~~/gen/altalune/v1/employee_pb';
import { ConnectError } from '@connectrpc/connect';

export function employeeRepository(client: Client<typeof EmployeeService>) {
  return {
    async queryEmployees(req: QueryEmployeesRequest): Promise<QueryEmployeesResponse> {
      try {
        const response = await client.queryEmployees(req);
        return response;
      }
      catch (err) {
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
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getEmployee(req: GetEmployeeRequest): Promise<GetEmployeeResponse> {
      try {
        const response = await client.getEmployee(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async updateEmployee(req: UpdateEmployeeRequest): Promise<UpdateEmployeeResponse> {
      try {
        const response = await client.updateEmployee(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async deleteEmployee(req: DeleteEmployeeRequest): Promise<DeleteEmployeeResponse> {
      try {
        const response = await client.deleteEmployee(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },
  };
}
