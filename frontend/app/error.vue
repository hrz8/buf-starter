<script setup lang="ts">
const error = useError();

const { t } = useI18n();
const isDevelopment = import.meta.dev;

</script>

<template>
  <div class="pl-5">
    <h1>Error</h1>
    <details v-if="isDevelopment && error?.message">
      <summary>
        Technical Details
      </summary>
      <div>
        <p>
          <strong>Error:</strong> {{ error.message }}
        </p>
        <p v-if="error.statusCode">
          <strong>Status:</strong> {{ error.statusCode }}
        </p>
        <p v-if="error.statusMessage">
          <strong>Status Message:</strong> {{ error.statusMessage }}
        </p>
        <textarea
          v-if="error.stack"
          cols="40"
          rows="10"
          class="border"
          readonly
          :value="error.stack"
        />
      </div>
    </details>
    <NuxtLink to="/">
      {{ t('errorPage.goHome') }}
    </NuxtLink>
  </div>
</template>
