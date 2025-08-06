import type { MessageInitShape } from '@bufbuild/protobuf';
import { create } from '@bufbuild/protobuf';
import { SayHelloRequestSchema } from '../../gen/greeter/v1/hello_pb';
import type { Client } from '@connectrpc/connect';
import type { GreeterService } from '../../gen/greeter/v1/greeter_pb';

export const greeterRepository = (client: Client<typeof GreeterService>) => ({
  async sayHello(req: MessageInitShape<typeof SayHelloRequestSchema>) {
    const request = create(SayHelloRequestSchema, req);

    try {
      const response = await client.sayHello(request);
      return response.message;
    } catch (error) {
      console.error('Error greeting:', error);
      throw error;
    }
  },
});
