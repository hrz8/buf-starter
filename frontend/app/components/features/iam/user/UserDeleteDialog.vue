<script setup lang="ts">
import type { User } from '~~/gen/altalune/v1/user_pb';
import { AlertCircle, Loader2 } from 'lucide-vue-next';
import { toast } from 'vue-sonner';

import {
  Alert,
  AlertDescription,
  AlertTitle,
} from '@/components/ui/alert';
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
import { useUserService } from '@/composables/services/useUserService';
import { getTranslatedConnectError } from './error';

const props = defineProps<{
  user: User;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t } = useI18n();

const {
  deleteUser,
  deleteLoading,
  deleteError,
  resetDeleteState,
} = useUserService();

const isDialogOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

async function handleDelete() {
  try {
    const success = await deleteUser({
      id: props.user.id,
    });

    if (success) {
      toast.success(t('features.users.messages.deleteSuccess'), {
        description: t('features.users.messages.deleteSuccessDesc', {
          name: `${props.user.firstName} ${props.user.lastName}`,
        }),
      });

      emit('success');
      isDialogOpen.value = false;
    }
  }
  catch (error) {
    console.error('Failed to delete user:', error);
    toast.error(t('features.users.messages.deleteError'), {
      description: deleteError.value || getTranslatedConnectError(error, t),
    });
  }
}

function handleCancel() {
  resetDeleteState();
  emit('cancel');
}

onUnmounted(() => {
  resetDeleteState();
});
</script>

<template>
  <AlertDialog v-model:open="isDialogOpen">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ t('features.users.dialog.deleteTitle') }}</AlertDialogTitle>
        <AlertDialogDescription>
          {{
            t('features.users.dialog.deleteDescription', {
              name: `${user.firstName} ${user.lastName}`,
            })
          }}
        </AlertDialogDescription>
      </AlertDialogHeader>

      <Alert
        v-if="deleteError"
        variant="destructive"
      >
        <AlertCircle class="w-4 h-4" />
        <AlertTitle>{{ t('common.status.error') }}</AlertTitle>
        <AlertDescription>{{ deleteError }}</AlertDescription>
      </Alert>

      <AlertDialogFooter>
        <AlertDialogCancel
          :disabled="deleteLoading"
          @click="handleCancel"
        >
          {{ t('common.btn.cancel') }}
        </AlertDialogCancel>
        <AlertDialogAction
          :disabled="deleteLoading"
          @click="handleDelete"
        >
          <Loader2
            v-if="deleteLoading"
            class="mr-2 h-4 w-4 animate-spin"
          />
          {{ deleteLoading ? t('common.status.deleting') : t('common.btn.delete') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
