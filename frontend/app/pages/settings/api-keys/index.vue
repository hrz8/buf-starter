<script setup lang="ts">
import { ApiKeyTable } from '@/components/features/api_key';
import { useProjectStore } from '@/stores/project';

definePageMeta({
  layout: 'default',
});

const { t } = useI18n();
const projectStore = useProjectStore();
const projectId = computed(() => projectStore.activeProjectId);
</script>

<template>
  <div v-if="!projectId" class="container mx-auto text-center py-8">
    <p class="text-muted-foreground">
      {{ t('features.api_keys.page.selectProjectMessage') }}
    </p>
  </div>
  <ApiKeyTable
    v-else
    :project-id="projectId"
  />
</template>
