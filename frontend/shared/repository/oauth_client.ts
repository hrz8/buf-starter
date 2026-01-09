import type { Client } from '@connectrpc/connect';

import type {
  CreateOAuthClientRequest,
  CreateOAuthClientResponse,
  DeleteOAuthClientRequest,
  DeleteOAuthClientResponse,
  GetOAuthClientRequest,
  GetOAuthClientResponse,
  OAuthClientService,
  QueryOAuthClientsRequest,
  QueryOAuthClientsResponse,
  RevealOAuthClientSecretRequest,
  RevealOAuthClientSecretResponse,
  UpdateOAuthClientRequest,
  UpdateOAuthClientResponse,
} from '~~/gen/altalune/v1/oauth_client_pb';
import { ConnectError } from '@connectrpc/connect';

export function oauthClientRepository(client: Client<typeof OAuthClientService>) {
  return {
    /**
     * Create OAuth client
     * Returns plaintext client_secret (shown once)
     */
    async createOAuthClient(req: CreateOAuthClientRequest): Promise<CreateOAuthClientResponse> {
      try {
        const response = await client.createOAuthClient(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    /**
     * Query OAuth clients with pagination/filtering/sorting
     * Never returns client_secret_hash
     */
    async queryOAuthClients(req: QueryOAuthClientsRequest): Promise<QueryOAuthClientsResponse> {
      try {
        const response = await client.queryOAuthClients(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    /**
     * Get single OAuth client by ID
     * Never returns client_secret_hash
     */
    async getOAuthClient(req: GetOAuthClientRequest): Promise<GetOAuthClientResponse> {
      try {
        const response = await client.getOAuthClient(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    /**
     * Update OAuth client
     * Optional secret re-hashing
     */
    async updateOAuthClient(req: UpdateOAuthClientRequest): Promise<UpdateOAuthClientResponse> {
      try {
        const response = await client.updateOAuthClient(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    /**
     * Delete OAuth client
     * Blocked for default client
     */
    async deleteOAuthClient(req: DeleteOAuthClientRequest): Promise<DeleteOAuthClientResponse> {
      try {
        const response = await client.deleteOAuthClient(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    /**
     * Reveal client secret (hashed)
     * Audit logged on backend
     */
    async revealOAuthClientSecret(
      req: RevealOAuthClientSecretRequest,
    ): Promise<RevealOAuthClientSecretResponse> {
      try {
        const response = await client.revealOAuthClientSecret(req);
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
