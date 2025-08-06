import { greeterRepository } from '#shared/repository/greeter';

export function useGreeter() {
  const { $greeterClient } = useNuxtApp();
  const greeter = greeterRepository($greeterClient);

  const response = ref('');
  const error = ref('');
  const loading = ref(false);

  async function submit(name: string) {
    loading.value = true;
    response.value = '';
    error.value = '';

    try {
      const result = await greeter.sayHello({ name: name });
      response.value = result;
    } catch (err: any) {
      error.value = err?.message || 'Something went wrong';
    } finally {
      loading.value = false;
    }
  };

  return {
    response,
    error,
    loading,
    submit,
  };
}
