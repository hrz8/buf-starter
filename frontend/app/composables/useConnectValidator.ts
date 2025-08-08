import type { ValidationResult } from '@bufbuild/protovalidate';
import type { DescMessage, MessageInitShape } from '@bufbuild/protobuf';
import { create as createSchema } from '@bufbuild/protobuf';

export function useConnectValidator<T extends DescMessage>(schema: T) {
  const { $validator } = useNuxtApp();

  const errors = ref<Record<string, string[]>>({});

  function validate(input: MessageInitShape<T>): boolean {
    const message = createSchema(schema, input);
    const result: ValidationResult = $validator.validate(schema, message);

    errors.value = {};

    if (result.kind === 'invalid') {
      for (const violation of result.violations) {
        const fieldPath = violation.field
          .map((f) => {
            if ('name' in f && typeof f.name === 'string' && f.kind === 'field') {
              return f.name;
            }
            return '';
          })
          .filter(Boolean)
          .join('.') || 'form';
        if (!errors.value[fieldPath]) {
          errors.value[fieldPath] = [];
        }
        errors.value[fieldPath].push(violation.message);
      }
      return false;
    }

    return true;
  }

  function reset() {
    errors.value = {};
  }

  return {
    validate,
    reset,
    errors: readonly(errors),
  };
}
