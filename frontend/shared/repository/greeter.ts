import type { Client } from '@connectrpc/connect';

import type { SayHelloRequest, SayHelloResponse } from '~~/gen/greeter/v1/hello_pb';
import type { GreeterService } from '~~/gen/greeter/v1/greeter_pb';

export const greeterRepository = (client: Client<typeof GreeterService>) => ({
  async sayHello(req: SayHelloRequest): Promise<SayHelloResponse> {
    try {
      const response = await client.sayHello(req);
      return response;
    } catch (error) {
      console.error('Error greeting:', error);
      throw error;
    }
  },
});
