<script setup lang="ts">
import type { Permission } from '~~/gen/altalune/v1/permission_pb';

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
import { usePermissionService } from '@/composables/services/usePermissionService';
import { useI18nSafe } from '@/composables/useI18nSafe';

const props = defineProps<{
  permission: Permission;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t, tFormatted } = useI18nSafe();
const { deletePermission, deleteLoading, deleteError, resetDeleteState } = usePermissionService();
const isDialogOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

async function handleDelete() {
  try {
    const success = await deletePermission({
      id: props.permission.id,
    });

    if (success) {
      toast.success(t('features.permissions.messages.deleteSuccess'), {
        description: t('features.permissions.messages.deleteSuccessDesc', { name: props.permission.name }),
      });

      isDialogOpen.value = false;
      emit('success');
    }
  }
  catch {
    toast.error(t('features.permissions.messages.deleteError'), {
      description: deleteError.value || t('features.permissions.messages.deleteErrorDesc'),
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
        <AlertDialogTitle>{{ t('features.permissions.deleteDialog.title') }}</AlertDialogTitle>
        <AlertDialogDescription>
          <component
            :is="tFormatted(
              'features.permissions.deleteDialog.confirmMessage',
              { name: permission.name },
            )"
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
          {{ deleteLoading
            ? t('common.status.deleting')
            : t('features.permissions.actions.delete') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
