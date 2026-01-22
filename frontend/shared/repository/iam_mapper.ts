import type { Client } from '@connectrpc/connect';
import type {
  AssignProjectMembersRequest,
  AssignRolePermissionsRequest,
  AssignUserPermissionsRequest,
  AssignUserRolesRequest,
  GetProjectMembersRequest,
  GetProjectMembersResponse,
  GetRolePermissionsRequest,
  GetRolePermissionsResponse,
  GetUserPermissionsRequest,
  GetUserPermissionsResponse,
  GetUserProjectsRequest,
  GetUserProjectsResponse,
  GetUserRolesRequest,
  GetUserRolesResponse,
  IAMMapperService,
  RemoveProjectMembersRequest,
  RemoveRolePermissionsRequest,
  RemoveUserPermissionsRequest,
  RemoveUserRolesRequest,
} from '~~/gen/altalune/v1/iam_mapper_pb';
import { ConnectError } from '@connectrpc/connect';

export function iamMapperRepository(client: Client<typeof IAMMapperService>) {
  return {
    // User-Role Mappings
    async assignUserRoles(req: AssignUserRolesRequest) {
      try {
        const response = await client.assignUserRoles(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async removeUserRoles(req: RemoveUserRolesRequest) {
      try {
        const response = await client.removeUserRoles(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getUserRoles(req: GetUserRolesRequest): Promise<GetUserRolesResponse> {
      try {
        const response = await client.getUserRoles(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    // Role-Permission Mappings
    async assignRolePermissions(req: AssignRolePermissionsRequest) {
      try {
        const response = await client.assignRolePermissions(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async removeRolePermissions(req: RemoveRolePermissionsRequest) {
      try {
        const response = await client.removeRolePermissions(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getRolePermissions(req: GetRolePermissionsRequest): Promise<GetRolePermissionsResponse> {
      try {
        const response = await client.getRolePermissions(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    // User-Permission Mappings
    async assignUserPermissions(req: AssignUserPermissionsRequest) {
      try {
        const response = await client.assignUserPermissions(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async removeUserPermissions(req: RemoveUserPermissionsRequest) {
      try {
        const response = await client.removeUserPermissions(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getUserPermissions(req: GetUserPermissionsRequest): Promise<GetUserPermissionsResponse> {
      try {
        const response = await client.getUserPermissions(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    // Project Members
    async assignProjectMembers(req: AssignProjectMembersRequest) {
      try {
        const response = await client.assignProjectMembers(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async removeProjectMembers(req: RemoveProjectMembersRequest) {
      try {
        const response = await client.removeProjectMembers(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getProjectMembers(req: GetProjectMembersRequest): Promise<GetProjectMembersResponse> {
      try {
        const response = await client.getProjectMembers(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    // User Projects (reverse lookup - projects a user belongs to)
    async getUserProjects(req: GetUserProjectsRequest): Promise<GetUserProjectsResponse> {
      try {
        const response = await client.getUserProjects(req);
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
