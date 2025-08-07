<script lang="ts" setup>
import { useGreeter } from '@/composables/useGreeter';
import ExampleLanguageSelector from './ExampleLanguageSelector.vue';

const name = ref('');
const { t } = useI18n();

const {
  response,
  error: serverError,
  validationErrors,
  loading,
  submit,
} = useGreeter();

function onSubmit() {
  submit({ name: name.value.trim() });
}
</script>

<template>
  <div class="max-w-md mx-auto mt-10 p-6 bg-white shadow-md rounded-md space-y-4">
    <ExampleLanguageSelector />

    <h2 class="text-xl font-semibold text-gray-800">
      {{ t('example.header') }}
    </h2>

    <input
      v-model="name"
      type="text"
      :placeholder="t('example.inputPlaceholder')"
      class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
    >
    <p
      v-for="(errorMsg, index) in validationErrors.name"
      :key="index"
      class="text-red-600 text-sm"
    >
      {{ errorMsg }}
    </p>

    <button
      :disabled="loading || !name"
      class="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 disabled:opacity-50"
      @click="onSubmit"
    >
      {{ loading ? t('example.submitBtnProgress') : t('example.submitBtn') }}
    </button>

    <p
      v-if="response"
      class="text-green-600 font-medium"
    >
      {{ t('example.response') }}: {{ response }}
    </p>

    <p
      v-if="serverError"
      class="text-red-600 font-medium"
    >
      {{ t('example.error') }}: {{ serverError }}
    </p>
  </div>
</template>
