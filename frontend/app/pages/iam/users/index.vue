<script setup lang="ts">
import { UserTable } from '@/components/features/iam/user';
import { usePageTitle } from '@/composables/usePageTitle';

definePageMeta({
  layout: 'default',
});

const { t } = useI18n();
usePageTitle(computed(() => t('features.users.page.title')));
</script>

<template>
  <UserTable />
</template>
