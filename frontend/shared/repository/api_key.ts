import type { Client } from '@connectrpc/connect';

import type {
  ActivateApiKeyRequest,
  ActivateApiKeyResponse,
  ApiKeyService,
  CreateApiKeyRequest,
  CreateApiKeyResponse,
  DeactivateApiKeyRequest,
  DeactivateApiKeyResponse,
  DeleteApiKeyRequest,
  DeleteApiKeyResponse,
  GetApiKeyRequest,
  GetApiKeyResponse,
  QueryApiKeysRequest,
  QueryApiKeysResponse,
  UpdateApiKeyRequest,
  UpdateApiKeyResponse,
} from '~~/gen/altalune/v1/api_key_pb';
import { ConnectError } from '@connectrpc/connect';

export function apiKeyRepository(client: Client<typeof ApiKeyService>) {
  return {
    async queryApiKeys(req: QueryApiKeysRequest): Promise<QueryApiKeysResponse> {
      try {
        const response = await client.queryApiKeys(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async createApiKey(req: CreateApiKeyRequest): Promise<CreateApiKeyResponse> {
      try {
        const response = await client.createApiKey(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getApiKey(req: GetApiKeyRequest): Promise<GetApiKeyResponse> {
      try {
        const response = await client.getApiKey(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async updateApiKey(req: UpdateApiKeyRequest): Promise<UpdateApiKeyResponse> {
      try {
        const response = await client.updateApiKey(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async deleteApiKey(req: DeleteApiKeyRequest): Promise<DeleteApiKeyResponse> {
      try {
        const response = await client.deleteApiKey(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async activateApiKey(req: ActivateApiKeyRequest): Promise<ActivateApiKeyResponse> {
      try {
        const response = await client.activateApiKey(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async deactivateApiKey(req: DeactivateApiKeyRequest): Promise<DeactivateApiKeyResponse> {
      try {
        const response = await client.deactivateApiKey(req);
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
