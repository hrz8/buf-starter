import type { Client } from '@connectrpc/connect';
import type {
  CreateRoleRequest,
  CreateRoleResponse,
  DeleteRoleRequest,
  DeleteRoleResponse,
  GetRoleRequest,
  GetRoleResponse,
  QueryRolesRequest,
  QueryRolesResponse,
  RoleService,
  UpdateRoleRequest,
  UpdateRoleResponse,
} from '~~/gen/altalune/v1/role_pb';
import { ConnectError } from '@connectrpc/connect';

export function roleRepository(client: Client<typeof RoleService>) {
  return {
    async queryRoles(req: QueryRolesRequest): Promise<QueryRolesResponse> {
      try {
        const response = await client.queryRoles(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async createRole(req: CreateRoleRequest): Promise<CreateRoleResponse> {
      try {
        const response = await client.createRole(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getRole(req: GetRoleRequest): Promise<GetRoleResponse> {
      try {
        const response = await client.getRole(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async updateRole(req: UpdateRoleRequest): Promise<UpdateRoleResponse> {
      try {
        const response = await client.updateRole(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async deleteRole(req: DeleteRoleRequest): Promise<DeleteRoleResponse> {
      try {
        const response = await client.deleteRole(req);
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
