import type { Client } from '@connectrpc/connect';

import type {
  CreateProjectRequest,
  CreateProjectResponse,
  DeleteProjectRequest,
  DeleteProjectResponse,
  GetProjectRequest,
  GetProjectResponse,
  ProjectService,
  QueryProjectsRequest,
  QueryProjectsResponse,
  UpdateProjectRequest,
  UpdateProjectResponse,
} from '~~/gen/altalune/v1/project_pb';
import { ConnectError } from '@connectrpc/connect';

export function projectRepository(client: Client<typeof ProjectService>) {
  return {
    async queryProjects(req: QueryProjectsRequest): Promise<QueryProjectsResponse> {
      try {
        const response = await client.queryProjects(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async createProject(req: CreateProjectRequest): Promise<CreateProjectResponse> {
      try {
        const response = await client.createProject(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getProject(req: GetProjectRequest): Promise<GetProjectResponse> {
      try {
        const response = await client.getProject(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async updateProject(req: UpdateProjectRequest): Promise<UpdateProjectResponse> {
      try {
        const response = await client.updateProject(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async deleteProject(req: DeleteProjectRequest): Promise<DeleteProjectResponse> {
      try {
        const response = await client.deleteProject(req);
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
