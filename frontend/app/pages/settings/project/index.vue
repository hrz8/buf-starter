<script setup lang="ts">
import ProjectSettingsForm from '@/components/features/project/ProjectSettingsForm.vue';
import { useProjectStore } from '~/stores/project';

const { t } = useI18n();
const projectStore = useProjectStore();
const projectId = computed(() => projectStore.activeProjectId);
</script>

<template>
  <div class="container mx-auto p-6">
    <div v-if="!projectId" class="text-center py-8">
      <p class="text-muted-foreground">
        {{ t('features.projects.settings.messages.noProjectSelected') }}
      </p>
    </div>

    <ProjectSettingsForm v-else :project-id="projectId" />
  </div>
</template>
