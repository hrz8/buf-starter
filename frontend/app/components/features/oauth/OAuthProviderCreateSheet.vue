<script setup lang="ts">
import type { OAuthProvider } from '~~/gen/altalune/v1/oauth_provider_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';

import OAuthProviderCreateForm from './OAuthProviderCreateForm.vue';

const props = defineProps<{
  loading?: boolean;
}>();

const emit = defineEmits<{
  success: [provider: OAuthProvider | null];
  cancel: [];
}>();

const { t } = useI18n();
const isSheetOpen = ref(false);

function handleProviderCreated(provider: OAuthProvider | null) {
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
    <SheetTrigger as-child>
      <slot />
    </SheetTrigger>
    <SheetContent class="w-full sm:max-w-[540px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>{{ t('features.oauth.sheet.createTitle') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.oauth.sheet.createDescription') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <OAuthProviderCreateForm
          :loading="props.loading"
          @success="handleProviderCreated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
