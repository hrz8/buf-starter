<script setup lang="ts">
import type { NuxtError } from '#app';

const props = defineProps({
  error: {
    type: Object as () => NuxtError,
    default: () => ({
      statusCode: 500,
      message: 'Unknown error',
    }),
  },
});

const { t } = useI18n();
</script>

<template>
  <div>
    <h1>{{ props.error.statusCode }}</h1>
    <NuxtLink to="/">
      {{ t('errorPage.goHome') }}
    </NuxtLink>
  </div>
</template>
