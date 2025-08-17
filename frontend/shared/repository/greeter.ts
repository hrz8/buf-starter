import { ConnectError } from '@connectrpc/connect';

import type { GetAllowedNamesResponse, GetAllowedNamesRequest } from '~~/gen/greeter/v1/name_pb';
import type { SayHelloResponse, SayHelloRequest } from '~~/gen/greeter/v1/hello_pb';
import type { GreeterService } from '~~/gen/greeter/v1/greeter_pb';
import type { Client } from '@connectrpc/connect';

export const greeterRepository = (client: Client<typeof GreeterService>) => ({
  async sayHello(req: SayHelloRequest): Promise<SayHelloResponse> {
    try {
      const response = await client.sayHello(req);
      return response;
    } catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err);
      }
      throw err;
    }
  },
  async getAllowedNames(req: GetAllowedNamesRequest): Promise<GetAllowedNamesResponse> {
    try {
      const response = await client.getAllowedNames(req);
      return response;
    } catch (err) {
      if (err instanceof ConnectError) {
        console.error('ConnectError:', err);
      }
      throw err;
    }
  },
});
