import { createClient } from '@connectrpc/connect';
import { createValidator } from '@bufbuild/protovalidate';
import { createConnectTransport } from '@connectrpc/connect-web';

import { GreeterService } from '~~/gen/greeter/v1/greeter_pb';

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig();

  const transport = createConnectTransport({
    baseUrl: config.public.apiUrl,
  });

  const validator = createValidator();
  const greeterClient = createClient(GreeterService, transport);

  return {
    provide: { validator, greeterClient },
  };
});
