<script setup lang="ts">
import type { OAuthProvider } from '~~/gen/altalune/v1/oauth_provider_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';

import OAuthProviderUpdateForm from './OAuthProviderUpdateForm.vue';

const props = defineProps<{
  provider: OAuthProvider;
  open: boolean;
  loading?: boolean;
}>();

const emit = defineEmits<{
  'success': [provider: OAuthProvider | null];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t } = useI18n();

const isSheetOpen = computed({
  get: () => props.open,
  set: (value: boolean) => emit('update:open', value),
});

function handleProviderUpdated(provider: OAuthProvider | null) {
  isSheetOpen.value = false;
  emit('success', provider);
}

function handleSheetClose() {
  isSheetOpen.value = false;
  emit('cancel');
}
</script>

<template>
  <Sheet v-model:open="isSheetOpen">
    <SheetContent class="w-full sm:max-w-[540px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>{{ t('features.oauth.sheet.editTitle') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.oauth.sheet.editDescription') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <OAuthProviderUpdateForm
          :provider="props.provider"
          :loading="props.loading"
          @success="handleProviderUpdated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
