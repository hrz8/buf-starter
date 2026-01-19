import type { Client } from '@connectrpc/connect';

import type {
  ChatbotService,
  GetChatbotConfigRequest,
  GetChatbotConfigResponse,
  UpdateModuleConfigRequest,
  UpdateModuleConfigResponse,
} from '~~/gen/altalune/v1/chatbot_pb';
import { ConnectError } from '@connectrpc/connect';

export function chatbotRepository(client: Client<typeof ChatbotService>) {
  return {
    async getChatbotConfig(req: GetChatbotConfigRequest): Promise<GetChatbotConfigResponse> {
      try {
        const response = await client.getChatbotConfig(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async updateModuleConfig(req: UpdateModuleConfigRequest): Promise<UpdateModuleConfigResponse> {
      try {
        const response = await client.updateModuleConfig(req);
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
