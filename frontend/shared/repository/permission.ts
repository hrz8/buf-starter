import type { Client } from '@connectrpc/connect';
import type {
  CreatePermissionRequest,
  CreatePermissionResponse,
  DeletePermissionRequest,
  DeletePermissionResponse,
  GetPermissionRequest,
  GetPermissionResponse,
  PermissionService,
  QueryPermissionsRequest,
  QueryPermissionsResponse,
  UpdatePermissionRequest,
  UpdatePermissionResponse,
} from '~~/gen/altalune/v1/permission_pb';
import { ConnectError } from '@connectrpc/connect';

export function permissionRepository(client: Client<typeof PermissionService>) {
  return {
    async queryPermissions(req: QueryPermissionsRequest): Promise<QueryPermissionsResponse> {
      try {
        const response = await client.queryPermissions(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async createPermission(req: CreatePermissionRequest): Promise<CreatePermissionResponse> {
      try {
        const response = await client.createPermission(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getPermission(req: GetPermissionRequest): Promise<GetPermissionResponse> {
      try {
        const response = await client.getPermission(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async updatePermission(req: UpdatePermissionRequest): Promise<UpdatePermissionResponse> {
      try {
        const response = await client.updatePermission(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async deletePermission(req: DeletePermissionRequest): Promise<DeletePermissionResponse> {
      try {
        const response = await client.deletePermission(req);
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
