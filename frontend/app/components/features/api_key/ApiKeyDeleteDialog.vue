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
      toast.success('API key deleted successfully', {
        description: `${props.apiKey.name} has been removed.`,
      });

      isDialogOpen.value = false;
      emit('success');
    }
  }
  catch {
    toast.error('Failed to delete API key', {
      description: deleteError.value || 'An unexpected error occurred. Please try again.',
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
        <AlertDialogTitle>Delete API Key</AlertDialogTitle>
        <AlertDialogDescription>
          Are you sure you want to delete <strong>{{ apiKey.name }}</strong>?
          This action cannot be undone and will permanently revoke this API key.
          Any applications using this key will lose access immediately.
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel
          :disabled="deleteLoading"
          @click="handleCancel"
        >
          Cancel
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
          {{ deleteLoading ? 'Deleting...' : 'Delete API Key' }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
