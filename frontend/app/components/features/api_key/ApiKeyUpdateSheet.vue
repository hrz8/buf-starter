<script setup lang="ts">
import type { ApiKey } from '~~/gen/altalune/v1/api_key_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';

import ApiKeyUpdateForm from './ApiKeyUpdateForm.vue';

const props = defineProps<{
  projectId: string;
  apiKey: ApiKey;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [apiKey: ApiKey];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const isSheetOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

function handleApiKeyUpdated(apiKey: ApiKey) {
  isSheetOpen.value = false;
  emit('success', apiKey);
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
        <SheetTitle>Edit API Key</SheetTitle>
        <SheetDescription>
          Update API key details below. All fields marked with * are required.
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <ApiKeyUpdateForm
          :project-id="props.projectId"
          :api-key-id="props.apiKey.id"
          @success="handleApiKeyUpdated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
