<script setup lang="ts">
import type { Project } from '~~/gen/altalune/v1/project_pb';
import { toast } from 'vue-sonner';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useProjectService } from '~/composables/services/useProjectService';
import { useProjectStore } from '~/stores/project';

const props = defineProps<{
  project: Project;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t } = useI18n();
const projectStore = useProjectStore();
const { deleteProject, deleteLoading, deleteError, resetDeleteState } = useProjectService();

// Dialog state
const isDialogOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

// Confirmation text
const confirmText = ref('');
const isLastProject = computed(() => projectStore.projects.length <= 1);

// Validation
const confirmationValid = computed(() => {
  if (!confirmText.value || !props.project)
    return false;

  const text = confirmText.value.trim();
  const projectName = props.project.name;

  // Case-insensitive for command, case-sensitive for project name
  const pattern = new RegExp(`^confirm delete ${projectName}$`, 'i');

  return pattern.test(text);
});

// Handle delete
async function handleDelete() {
  if (!confirmationValid.value) {
    toast.error(t('features.projects.settings.messages.deleteConfirmMismatch'));
    return;
  }

  try {
    const success = await deleteProject({ id: props.project.id });
    if (success) {
      toast.success(t('features.projects.settings.messages.deleteSuccess'), {
        description: t('features.projects.settings.messages.deleteSuccessDesc', { name: props.project.name }),
      });

      projectStore.removeProject(props.project.id);
      emit('success');

      // Hard refresh per user story
      window.location.reload();
    }
  }
  catch (error) {
    console.error('Failed to delete project:', error);
    toast.error(t('features.projects.settings.messages.deleteError'), {
      description: deleteError.value || t('features.projects.settings.messages.deleteErrorDesc'),
    });
  }
}

function handleCancel() {
  confirmText.value = '';
  resetDeleteState();
  emit('cancel');
  isDialogOpen.value = false;
}

// Reset on open
watch(isDialogOpen, (newValue) => {
  if (!newValue) {
    confirmText.value = '';
    resetDeleteState();
  }
});
</script>

<template>
  <AlertDialog v-model:open="isDialogOpen">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>
          {{ t('features.projects.settings.deleteDialog.title') }}
        </AlertDialogTitle>
        <AlertDialogDescription>
          {{ t('features.projects.settings.deleteDialog.description') }}
        </AlertDialogDescription>
      </AlertDialogHeader>

      <div v-if="isLastProject" class="rounded-md bg-yellow-50 p-4 dark:bg-yellow-950">
        <div class="flex">
          <Icon
            name="lucide:alert-triangle"
            size="1em"
            mode="svg"
            class="text-yellow-600 dark:text-yellow-500"
          />
          <div class="ml-3">
            <p class="text-sm text-yellow-800 dark:text-yellow-200">
              {{ t('features.projects.settings.deleteDialog.lastProjectWarning') }}
            </p>
          </div>
        </div>
      </div>

      <div v-else class="space-y-2">
        <Label>{{ t('features.projects.settings.deleteDialog.confirmLabel') }}</Label>
        <Input
          v-model="confirmText"
          :placeholder="t('features.projects.settings.deleteDialog.confirmPlaceholder')"
          :disabled="deleteLoading"
        />
        <p class="text-sm text-muted-foreground">
          {{
            t(
              'features.projects.settings.deleteDialog.confirmDescription',
              {
                name: project.name,
              },
            )
          }}
        </p>
      </div>

      <AlertDialogFooter>
        <AlertDialogCancel :disabled="deleteLoading" @click="handleCancel">
          {{ t('features.projects.settings.actions.cancel') }}
        </AlertDialogCancel>
        <AlertDialogAction
          v-if="!isLastProject"
          :disabled="deleteLoading || !confirmationValid"
          class="bg-destructive text-white hover:bg-destructive/90"
          @click="handleDelete"
        >
          <Icon
            v-if="deleteLoading"
            name="lucide:loader-2"
            size="1em"
            mode="svg"
            class="mr-2 animate-spin"
          />
          {{ deleteLoading
            ? t('features.projects.settings.actions.deleting')
            : t('features.projects.settings.actions.delete') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
