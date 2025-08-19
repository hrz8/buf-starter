import { ConnectError } from '@connectrpc/connect';

import type {
  QueryProjectsResponse,
  QueryProjectsRequest, ProjectService,
} from '~~/gen/altalune/v1/project_pb';
import type { Client } from '@connectrpc/connect';

export const projectRepository = (client: Client<typeof ProjectService>) => ({
  async queryProjects(req: QueryProjectsRequest): Promise<QueryProjectsResponse> {
    try {
      const response = await client.queryProjects(req);
      return response;
    } catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err);
      }
      throw err;
    }
  },
});
