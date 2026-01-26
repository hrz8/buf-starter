<script setup lang="ts">
import ProjectSettingsForm from '@/components/features/project/ProjectSettingsForm.vue';
import { usePageTitle } from '@/composables/usePageTitle';
import { useProjectStore } from '~/stores/project';

const { t } = useI18n();
const projectStore = useProjectStore();
const projectId = computed(() => projectStore.activeProjectId);

usePageTitle(computed(() => t('features.projects.settings.title')));
</script>

<template>
  <div class="container mx-auto px-2 py-3">
    <div v-if="!projectId" class="text-center py-8">
      <p class="text-muted-foreground">
        {{ t('features.projects.settings.messages.noProjectSelected') }}
      </p>
    </div>
    <ProjectSettingsForm v-else :project-id="projectId" />
  </div>
</template>
