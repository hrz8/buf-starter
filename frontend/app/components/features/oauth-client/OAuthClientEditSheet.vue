<script setup lang="ts">
import type { OAuthClient } from '~~/gen/altalune/v1/oauth_client_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';

import OAuthClientEditForm from './OAuthClientEditForm.vue';

const props = defineProps<{
  client: OAuthClient;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [client: OAuthClient];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t } = useI18n();
const isSheetOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

function handleClientUpdated(client: OAuthClient) {
  isSheetOpen.value = false;
  emit('success', client);
}

function handleSheetClose() {
  isSheetOpen.value = false;
  emit('cancel');
}
</script>

<template>
  <Sheet v-model:open="isSheetOpen">
    <!-- Only show trigger when not controlled externally -->
    <SheetTrigger
      v-if="!props.open && $slots.default"
      as-child
    >
      <slot />
    </SheetTrigger>
    <SheetContent class="w-full sm:max-w-[540px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>{{ t('features.oauth_clients.sheets.edit.title') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.oauth_clients.sheets.edit.description') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <OAuthClientEditForm
          :client-id="props.client.id"
          @success="handleClientUpdated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
