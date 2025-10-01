<script setup lang="ts">
import type { ApiKey } from '~~/gen/altalune/v1/api_key_pb';

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
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { useApiKeyService } from '@/composables/services/useApiKeyService';
import { useI18nSafe } from '@/composables/useI18nSafe';

const props = defineProps<{
  projectId: string;
  apiKey: ApiKey;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t, tFormatted } = useI18nSafe();
const { deleteApiKey, deleteLoading, deleteError, resetDeleteState } = useApiKeyService();
const isDialogOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

async function handleDelete() {
  try {
    const success = await deleteApiKey({
      projectId: props.projectId,
      apiKeyId: props.apiKey.id,
    });

    if (success) {
      toast.success(t('features.api_keys.messages.deleteSuccess'), {
        description: t('features.api_keys.messages.deleteSuccessDesc', { name: props.apiKey.name }),
      });

      isDialogOpen.value = false;
      emit('success');
    }
  }
  catch {
    toast.error(t('features.api_keys.messages.deleteError'), {
      description: deleteError.value || t('features.api_keys.messages.deleteErrorDesc'),
    });
  }
}

function handleCancel() {
  isDialogOpen.value = false;
  resetDeleteState();
  emit('cancel');
}

onUnmounted(() => {
  resetDeleteState();
});
</script>

<template>
  <AlertDialog v-model:open="isDialogOpen">
    <!-- Only show trigger when not controlled externally -->
    <AlertDialogTrigger
      v-if="!props.open && $slots.default"
      as-child
    >
      <slot />
    </AlertDialogTrigger>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ t('features.api_keys.deleteDialog.title') }}</AlertDialogTitle>
        <AlertDialogDescription>
          <component
            :is="tFormatted('features.api_keys.deleteDialog.confirmMessage', { name: apiKey.name })"
          />
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel
          :disabled="deleteLoading"
          @click="handleCancel"
        >
          {{ t('common.btn.cancel') }}
        </AlertDialogCancel>
        <AlertDialogAction
          :disabled="deleteLoading"
          class="bg-destructive text-white hover:bg-destructive/90 focus:ring-destructive"
          @click="handleDelete"
        >
          <Icon
            v-if="deleteLoading"
            name="lucide:loader-2"
            class="mr-2 h-4 w-4 animate-spin"
          />
          {{ deleteLoading ? t('common.status.deleting') : t('features.api_keys.actions.delete') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
