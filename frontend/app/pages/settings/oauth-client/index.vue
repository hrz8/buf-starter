<script setup lang="ts">
import { OAuthClientTable } from '@/components/features/oauth-client';
import { useProjectStore } from '@/stores/project';

definePageMeta({
  layout: 'default',
});

const projectStore = useProjectStore();
const projectId = computed(() => projectStore.activeProjectId);
</script>

<template>
  <div v-if="!projectId" class="container mx-auto text-center py-8">
    <p class="text-muted-foreground">
      Please select a project to view OAuth clients
    </p>
  </div>
  <OAuthClientTable
    v-else
    :project-id="projectId"
  />
</template>
