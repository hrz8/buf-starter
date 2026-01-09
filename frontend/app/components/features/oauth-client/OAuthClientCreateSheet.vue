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

import OAuthClientCreateForm from './OAuthClientCreateForm.vue';

const props = defineProps<{
  projectId: string;
  loading?: boolean;
}>();

const emit = defineEmits<{
  success: [result: { client: OAuthClient | null; clientSecret: string }];
  cancel: [];
}>();

const { t } = useI18n();
const isSheetOpen = ref(false);

function handleClientCreated(result: { client: OAuthClient | null; clientSecret: string }) {
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
    <SheetContent class="w-full sm:max-w-[600px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>{{ t('features.oauth_clients.sheets.create.title') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.oauth_clients.sheets.create.description') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <OAuthClientCreateForm
          :project-id="props.projectId"
          :loading="props.loading"
          @success="handleClientCreated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
