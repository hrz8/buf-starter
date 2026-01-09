<script setup lang="ts">
import type { OAuthClient } from '~~/gen/altalune/v1/oauth_client_pb';

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
import { useOAuthClientService } from '@/composables/services/useOAuthClientService';

const props = defineProps<{
  projectId: string;
  client: OAuthClient;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t } = useI18n();
const { deleteOAuthClient, deleteLoading, deleteError, resetDeleteState } = useOAuthClientService();
const isDialogOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

async function handleDelete() {
  try {
    const success = await deleteOAuthClient({
      id: props.client.id,
      projectId: props.projectId,
    });

    if (success) {
      toast.success(t('features.oauth_clients.toasts.deleted'), {
        description: t('features.oauth_clients.toasts.deletedDesc', { name: props.client.name }),
      });

      isDialogOpen.value = false;
      emit('success');
    }
  }
  catch {
    toast.error(t('features.oauth_clients.toasts.deleteFailed'), {
      description: deleteError.value || t('features.oauth_clients.toasts.deleteFailedDesc'),
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
        <AlertDialogTitle>
          {{ t('features.oauth_clients.dialogs.delete.title') }}
        </AlertDialogTitle>
        <AlertDialogDescription>
          <i18n-t keypath="features.oauth_clients.dialogs.delete.description" tag="span">
            <template #name>
              <strong>{{ client.name }}</strong>
            </template>
          </i18n-t>
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel
          :disabled="deleteLoading"
          @click="handleCancel"
        >
          {{ t('features.oauth_clients.actions.cancel') }}
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
          {{
            deleteLoading
              ? t('features.oauth_clients.actions.deleting')
              : t('features.oauth_clients.actions.delete')
          }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
