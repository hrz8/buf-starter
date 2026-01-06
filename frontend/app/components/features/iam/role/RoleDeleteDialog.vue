<script setup lang="ts">
import type { Role } from '~~/gen/altalune/v1/role_pb';

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
import { useRoleService } from '@/composables/services/useRoleService';
import { useI18nSafe } from '@/composables/useI18nSafe';

const props = defineProps<{
  role: Role;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t, tFormatted } = useI18nSafe();
const { deleteRole, deleteLoading, deleteError, resetDeleteState } = useRoleService();
const isDialogOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

async function handleDelete() {
  try {
    const success = await deleteRole({
      id: props.role.id,
    });

    if (success) {
      toast.success(t('features.roles.messages.deleteSuccess'), {
        description: t('features.roles.messages.deleteSuccessDesc', { name: props.role.name }),
      });

      isDialogOpen.value = false;
      emit('success');
    }
  }
  catch {
    toast.error(t('features.roles.messages.deleteError'), {
      description: deleteError.value || t('features.roles.messages.deleteErrorDesc'),
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
        <AlertDialogTitle>{{ t('features.roles.deleteDialog.title') }}</AlertDialogTitle>
        <AlertDialogDescription>
          <component
            :is="tFormatted(
              'features.roles.deleteDialog.confirmMessage',
              { name: role.name },
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
            : t('features.roles.actions.delete') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
