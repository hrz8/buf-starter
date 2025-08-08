import type { Client } from '@connectrpc/connect';
import { ConnectError } from '@connectrpc/connect';

import type { SayHelloRequest, SayHelloResponse } from '~~/gen/greeter/v1/hello_pb';
import type { GetAllowedNamesRequest, GetAllowedNamesResponse } from '~~/gen/greeter/v1/name_pb';
import type { GreeterService } from '~~/gen/greeter/v1/greeter_pb';

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
