import { ConnectError } from '@connectrpc/connect';

import type {
  QueryEmployeesResponse,
  CreateEmployeeResponse,
  UpdateEmployeeResponse,
  DeleteEmployeeResponse,
  QueryEmployeesRequest,
  CreateEmployeeRequest,
  UpdateEmployeeRequest,
  DeleteEmployeeRequest,
  GetEmployeeResponse,
  GetEmployeeRequest,
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

  async getEmployee(req: GetEmployeeRequest): Promise<GetEmployeeResponse> {
    try {
      const response = await client.getEmployee(req);
      return response;
    } catch (err) {
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
    } catch (err) {
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
    } catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err);
      }
      throw err;
    }
  },
});
