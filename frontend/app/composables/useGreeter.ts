import type { ValidationResult } from '@bufbuild/protovalidate';
import { create, type MessageInitShape } from '@bufbuild/protobuf';

import { greeterRepository } from '#shared/repository/greeter';
import { SayHelloRequestSchema, type SayHelloRequest } from '~~/gen/greeter/v1/hello_pb';

export function useGreeter() {
  const { $validator, $greeterClient } = useNuxtApp();
  const greeter = greeterRepository($greeterClient);

  const response = ref('');
  const error = ref('');
  const loading = ref(false);

  function validate(req: SayHelloRequest): ValidationResult {
    return $validator.validate(SayHelloRequestSchema, req);
  }

  async function submit(req: MessageInitShape<typeof SayHelloRequestSchema>): Promise<void> {
    loading.value = true;
    response.value = '';
    error.value = '';

    const request = create(SayHelloRequestSchema, req);
    const validated = validate(request);

    if (validated.kind === 'invalid') {
      error.value = validated.violations.map((e) => e.message).join(', ');
      loading.value = false;
      return;
    }

    try {
      const result = await greeter.sayHello(request);
      response.value = result.message;
    } catch (err: any) {
      error.value = err?.message || 'Something went wrong';
    } finally {
      loading.value = false;
    }
  };

  return {
    response: readonly(response),
    error: readonly(error),
    loading: readonly(loading),
    submit,
  };
}
