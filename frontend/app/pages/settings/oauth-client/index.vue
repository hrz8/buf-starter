<script setup lang="ts">
import { OAuthClientTable } from '@/components/features/oauth-client';
import { usePageTitle } from '@/composables/usePageTitle';

definePageMeta({
  layout: 'default',
});

// OAuth clients are GLOBAL entities (not project-scoped)
// Following Keycloak/Auth0 architecture patterns

const { t } = useI18n();
usePageTitle(computed(() => t('features.oauth_clients.page.title')));
</script>

<template>
  <OAuthClientTable />
</template>
