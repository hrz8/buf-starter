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

import ApiKeyCreateForm from './ApiKeyCreateForm.vue';

const props = defineProps<{
  projectId: string;
  loading?: boolean;
}>();

const emit = defineEmits<{
  success: [result: { apiKey: ApiKey | null; keyValue: string }];
  cancel: [];
}>();

const isSheetOpen = ref(false);

function handleApiKeyCreated(result: { apiKey: ApiKey | null; keyValue: string }) {
  isSheetOpen.value = false;
  emit('success', result);
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
        <SheetTitle>Create New API Key</SheetTitle>
        <SheetDescription>
          Create a new API key for your project. The generated key will only be shown once, so make
          sure to save it securely.
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <ApiKeyCreateForm
          :project-id="props.projectId"
          :loading="props.loading"
          @success="handleApiKeyCreated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
