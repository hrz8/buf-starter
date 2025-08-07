import { create, type MessageInitShape } from '@bufbuild/protobuf';
import { SayHelloRequestSchema } from '~~/gen/greeter/v1/hello_pb';

import { greeterRepository } from '#shared/repository/greeter';
import { useErrorMessage } from './useErrorMessage';
import { useConnectValidator } from './useConnectValidator';

export function useGreeter() {
  const { $greeterClient } = useNuxtApp();
  const greeter = greeterRepository($greeterClient);

  const { validate, errors: validationErrors } = useConnectValidator(SayHelloRequestSchema);
  const { parseError } = useErrorMessage();

  const response = ref('');
  const error = ref('');
  const loading = ref(false);

  async function submit(req: MessageInitShape<typeof SayHelloRequestSchema>): Promise<void> {
    loading.value = true;
    response.value = '';
    error.value = '';

    if (!validate(req)) {
      loading.value = false;
      return;
    }

    try {
      const message = create(SayHelloRequestSchema, req);
      const result = await greeter.sayHello(message);
      response.value = result.message;
    } catch (err) {
      error.value = parseError(err);
    } finally {
      loading.value = false;
    }
  }

  return {
    response: readonly(response),
    error: readonly(error),
    validationErrors,
    loading: readonly(loading),
    submit,
  };
}
