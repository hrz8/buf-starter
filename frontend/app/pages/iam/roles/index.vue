<script setup lang="ts">
import { RoleTable } from '@/components/features/iam/role';
import { usePageTitle } from '@/composables/usePageTitle';

definePageMeta({
  layout: 'default',
});

const { t } = useI18n();
usePageTitle(computed(() => t('features.roles.page.title')));
</script>

<template>
  <RoleTable />
</template>
