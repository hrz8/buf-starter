import { create, type MessageInitShape } from '@bufbuild/protobuf';
import { SayHelloRequestSchema } from '~~/gen/greeter/v1/hello_pb';
import { GetAllowedNamesRequestSchema } from '~~/gen/greeter/v1/name_pb';
import { greeterRepository } from '#shared/repository/greeter';
import { useErrorMessage } from '../useErrorMessage';
import { useConnectValidator } from '../useConnectValidator';
import type { PaginationMetaSchema } from '~~/gen/greeter/v1/common_pb';

export function useGreeter() {
  const { $greeterClient } = useNuxtApp();
  const greeter = greeterRepository($greeterClient);

  const helloValidator = useConnectValidator(SayHelloRequestSchema);
  const nameValidator = useConnectValidator(GetAllowedNamesRequestSchema);

  const { parseError } = useErrorMessage();

  const submitState = reactive({
    loading: false,
    error: '',
    response: '',
  });

  async function list(
    req: MessageInitShape<typeof GetAllowedNamesRequestSchema>,
  ): Promise<{
    names: string[];
    meta: MessageInitShape<typeof PaginationMetaSchema> | undefined;
  }> {
    nameValidator.reset();

    if (!nameValidator.validate(req)) {
      console.warn('Validation failed for GetAllowedNamesRequest:', nameValidator.errors.value);
      return {
        names: [],
        meta: { total: 0, page: 1, limit: 10, totalPages: 1, hasNext: false, hasPrev: false },
      };
    }

    try {
      const message = create(GetAllowedNamesRequestSchema, req);
      const result = await greeter.getAllowedNames(message);
      return {
        names: result.names,
        meta: result.meta,
      };
    } catch (err) {
      const errorMessage = parseError(err);
      throw new Error(errorMessage);
    }
  }

  async function submit(req: MessageInitShape<typeof SayHelloRequestSchema>): Promise<void> {
    submitState.loading = true;
    submitState.response = '';
    submitState.error = '';

    if (!helloValidator.validate(req)) {
      submitState.loading = false;
      return;
    }

    try {
      const message = create(SayHelloRequestSchema, req);
      const result = await greeter.sayHello(message);
      submitState.response = result.message;
    } catch (err) {
      submitState.error = parseError(err);
    } finally {
      submitState.loading = false;
    }
  }

  return {
    // List
    list,

    // Submit
    submit,
    submitLoading: computed(() => submitState.loading),
    submitError: computed(() => submitState.error),
    submitResponse: computed(() => submitState.response),

    // Validation
    helloValidationErrors: helloValidator.errors,
    nameValidationErrors: nameValidator.errors,
  };
}
