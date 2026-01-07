<script setup lang="ts">
import type { OAuthProvider } from '~~/gen/altalune/v1/oauth_provider_pb';

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
import { useOAuthProviderService } from '@/composables/services/useOAuthProviderService';
import { useI18nSafe } from '@/composables/useI18nSafe';
import { getProviderName } from './constants';

const props = defineProps<{
  provider: OAuthProvider;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t, tFormatted } = useI18nSafe();
const {
  deleteOAuthProvider,
  deleteLoading,
  deleteError,
  resetDeleteState,
} = useOAuthProviderService();

const isDialogOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

const providerName = computed(() => getProviderName(props.provider.providerType));

async function handleDelete() {
  try {
    const success = await deleteOAuthProvider({
      id: props.provider.id,
    });

    if (success) {
      toast.success(t('features.oauth.messages.deleteSuccess'), {
        description: t('features.oauth.messages.deleteSuccessDesc', { name: providerName.value }),
      });

      isDialogOpen.value = false;
      emit('success');
    }
  }
  catch {
    toast.error(t('features.oauth.messages.deleteError'), {
      description: deleteError.value || t('features.oauth.messages.deleteErrorDesc'),
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
        <AlertDialogTitle>{{ t('features.oauth.deleteDialog.title') }}</AlertDialogTitle>
        <AlertDialogDescription>
          <component
            :is="tFormatted('features.oauth.deleteDialog.confirmMessage', { name: providerName })"
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
          {{ deleteLoading ? t('common.status.deleting') : t('features.oauth.actions.delete') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
