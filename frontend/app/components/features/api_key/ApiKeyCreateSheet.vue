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

const { t } = useI18n();
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
        <SheetTitle>{{ t('features.api_keys.sheet.createTitle') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.api_keys.sheet.createDescription') }}
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
