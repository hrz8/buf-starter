import type { Client } from '@connectrpc/connect';

import type {
  CreateProjectRequest,
  CreateProjectResponse,
  ProjectService,
  QueryProjectsRequest,
  QueryProjectsResponse,
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
  };
}
