<script setup lang="ts">
import { OAuthProviderTable } from '@/components/features/oauth';
import { usePageTitle } from '@/composables/usePageTitle';

definePageMeta({
  layout: 'default',
});

const { t } = useI18n();
usePageTitle(computed(() => t('features.oauth.page.title')));
</script>

<template>
  <OAuthProviderTable />
</template>
