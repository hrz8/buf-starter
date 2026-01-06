import type { Client } from '@connectrpc/connect';
import type {
  ActivateUserRequest,
  ActivateUserResponse,
  CreateUserRequest,
  CreateUserResponse,
  DeactivateUserRequest,
  DeactivateUserResponse,
  DeleteUserRequest,
  DeleteUserResponse,
  GetUserRequest,
  GetUserResponse,
  QueryUsersRequest,
  QueryUsersResponse,
  UpdateUserRequest,
  UpdateUserResponse,
  UserService,
} from '~~/gen/altalune/v1/user_pb';
import { ConnectError } from '@connectrpc/connect';

export function userRepository(client: Client<typeof UserService>) {
  return {
    async queryUsers(req: QueryUsersRequest): Promise<QueryUsersResponse> {
      try {
        const response = await client.queryUsers(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async createUser(req: CreateUserRequest): Promise<CreateUserResponse> {
      try {
        const response = await client.createUser(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async getUser(req: GetUserRequest): Promise<GetUserResponse> {
      try {
        const response = await client.getUser(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async updateUser(req: UpdateUserRequest): Promise<UpdateUserResponse> {
      try {
        const response = await client.updateUser(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async deleteUser(req: DeleteUserRequest): Promise<DeleteUserResponse> {
      try {
        const response = await client.deleteUser(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async activateUser(req: ActivateUserRequest): Promise<ActivateUserResponse> {
      try {
        const response = await client.activateUser(req);
        return response;
      }
      catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },

    async deactivateUser(req: DeactivateUserRequest): Promise<DeactivateUserResponse> {
      try {
        const response = await client.deactivateUser(req);
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
