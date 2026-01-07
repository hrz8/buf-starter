import type { Client } from '@connectrpc/connect';
import type {
  CreateOAuthProviderRequest,
  CreateOAuthProviderResponse,
  DeleteOAuthProviderRequest,
  DeleteOAuthProviderResponse,
  GetOAuthProviderRequest,
  GetOAuthProviderResponse,
  OAuthProviderService,
  QueryOAuthProvidersRequest,
  QueryOAuthProvidersResponse,
  RevealClientSecretRequest,
  RevealClientSecretResponse,
  UpdateOAuthProviderRequest,
  UpdateOAuthProviderResponse,
} from '~~/gen/altalune/v1/oauth_provider_pb';
import { ConnectError } from '@connectrpc/connect';

export function oauthProviderRepository(client: Client<typeof OAuthProviderService>) {
  return {
    async queryOAuthProviders(
      req: QueryOAuthProvidersRequest,
    ): Promise<QueryOAuthProvidersResponse> {
      try {
        const response = await client.queryOAuthProviders(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async createOAuthProvider(
      req: CreateOAuthProviderRequest,
    ): Promise<CreateOAuthProviderResponse> {
      try {
        const response = await client.createOAuthProvider(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getOAuthProvider(req: GetOAuthProviderRequest): Promise<GetOAuthProviderResponse> {
      try {
        const response = await client.getOAuthProvider(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async updateOAuthProvider(
      req: UpdateOAuthProviderRequest,
    ): Promise<UpdateOAuthProviderResponse> {
      try {
        const response = await client.updateOAuthProvider(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async deleteOAuthProvider(
      req: DeleteOAuthProviderRequest,
    ): Promise<DeleteOAuthProviderResponse> {
      try {
        const response = await client.deleteOAuthProvider(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async revealClientSecret(req: RevealClientSecretRequest): Promise<RevealClientSecretResponse> {
      try {
        const response = await client.revealClientSecret(req);
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
