import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { GreeterService } from '../../gen/greeter/v1/greeter_pb';

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig();

  const transport = createConnectTransport({
    baseUrl: config.public.apiUrl,
  });

  const greeterClient = createClient(GreeterService, transport);

  return {
    provide: { greeterClient },
  };
});
