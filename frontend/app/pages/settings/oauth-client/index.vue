<script setup lang="ts">
import { OAuthClientTable } from '@/components/features/oauth-client';

definePageMeta({
  layout: 'default',
});

// OAuth clients are GLOBAL entities (not project-scoped)
// Following Keycloak/Auth0 architecture patterns
</script>

<template>
  <OAuthClientTable />
</template>
