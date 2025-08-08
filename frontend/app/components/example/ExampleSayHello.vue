<script setup lang="ts">
import { useGreeter } from '~/composables/services/useGreeter';

const name = ref('');
const { t } = useI18n();

const {
  submit,
  submitLoading,
  submitError,
  submitResponse,
  helloValidationErrors,
} = useGreeter();

function onSubmit() {
  submit({ name: name.value.trim() });
}
</script>

<template>
  <div class="space-y-4">
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
      v-for="(errorMsg, index) in helloValidationErrors"
      :key="index"
      class="text-red-600 text-sm"
    >
      {{ errorMsg }}
    </p>

    <button
      :disabled="submitLoading || !name"
      class="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 disabled:opacity-50"
      @click="onSubmit"
    >
      {{ submitLoading ? t('example.submitBtnProgress') : t('example.submitBtn') }}
    </button>

    <p
      v-if="submitResponse"
      class="text-green-600 font-medium"
    >
      {{ t('example.response') }}: {{ submitResponse }}
    </p>

    <p
      v-if="submitError"
      class="text-red-600 font-medium"
    >
      {{ t('example.error') }}: {{ submitError }}
    </p>
  </div>
</template>
