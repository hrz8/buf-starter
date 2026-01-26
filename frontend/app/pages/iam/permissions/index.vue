<script setup lang="ts">
import { PermissionTable } from '@/components/features/iam/permission';
import { usePageTitle } from '@/composables/usePageTitle';

definePageMeta({
  layout: 'default',
});

const { t } = useI18n();
usePageTitle(computed(() => t('features.permissions.page.title')));
</script>

<template>
  <PermissionTable />
</template>
