<template>
  <div class="max-w-md mx-auto mt-10 p-6 bg-white shadow-md rounded-md space-y-4">
    <h2 class="text-xl font-semibold text-gray-800">
      Say Change
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
      @click="onSubmit"
    >
      {{ loading ? 'Sending...' : 'Say Hello' }}
    </button>

    <p
      v-if="response"
      class="text-green-600 font-medium"
    >
      Response: {{ response }}
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
import { useGreeter } from '@/composables/useGreeter';

const name = ref('');

const { response, error, loading, submit } = useGreeter();

function onSubmit() {
  if (name.value.trim()) {
    submit({
      name: name.value.trim(),
    });
  }
}
</script>
