import type { Client } from '@connectrpc/connect';

import type {
  ChatbotNodeService,
  CreateNodeRequest,
  CreateNodeResponse,
  DeleteNodeRequest,
  DeleteNodeResponse,
  GetNodeRequest,
  GetNodeResponse,
  ListNodesRequest,
  ListNodesResponse,
  UpdateNodeRequest,
  UpdateNodeResponse,
} from '~~/gen/altalune/v1/chatbot_node_pb';
import { ConnectError } from '@connectrpc/connect';

export function chatbotNodeRepository(client: Client<typeof ChatbotNodeService>) {
  return {
    async listNodes(req: ListNodesRequest): Promise<ListNodesResponse> {
      try {
        const response = await client.listNodes(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async createNode(req: CreateNodeRequest): Promise<CreateNodeResponse> {
      try {
        const response = await client.createNode(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getNode(req: GetNodeRequest): Promise<GetNodeResponse> {
      try {
        const response = await client.getNode(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async updateNode(req: UpdateNodeRequest): Promise<UpdateNodeResponse> {
      try {
        const response = await client.updateNode(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async deleteNode(req: DeleteNodeRequest): Promise<DeleteNodeResponse> {
      try {
        const response = await client.deleteNode(req);
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
