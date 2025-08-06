<template>
  <div class="max-w-md mx-auto mt-10 p-6 bg-white shadow-md rounded-md space-y-4">
    <h2 class="text-xl font-semibold text-gray-800">
      Say Hello
    </h2>

    <input
      v-model="name"
      type="text"
      placeholder="Enter your name"
      class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
    >

    <button
      :disabled="loading || !name"
      class="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 disabled:opacity-50"
      @click="submit"
    >
      {{ loading ? 'Sending...' : 'Say Hello' }}
    </button>

    <p
      v-if="message"
      class="text-green-600 font-medium"
    >
      Response: {{ message }}
    </p>
    <p
      v-if="error"
      class="text-red-600 font-medium"
    >
      Error: {{ error }}
    </p>
  </div>
</template>

<script lang="ts" setup>
import { greeterRepository } from '#shared/repository/greeter';

const name = ref('');
const message = ref('');
const error = ref('');
const loading = ref(false);

const { $greeterClient } = useNuxtApp();
const greeter = greeterRepository($greeterClient);

const submit = async () => {
  loading.value = true;
  message.value = '';
  error.value = '';

  try {
    const result = await greeter.sayHello({ name: name.value });
    message.value = result;
  } catch (err: any) {
    error.value = err?.message || 'Something went wrong';
  } finally {
    loading.value = false;
  }
};
</script>
