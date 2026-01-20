<script setup lang="ts">
import type { ChatbotNode } from '~~/gen/chatbot/nodes/v1/node_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';

import NodeCreateForm from './NodeCreateForm.vue';

interface Props {
  open?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  open: false,
});

const emit = defineEmits<{
  'success': [node: ChatbotNode];
  'cancel': [];
  'update:open': [open: boolean];
}>();

const { t } = useI18n();

const isSheetOpen = computed({
  get: () => props.open,
  set: (value: boolean) => emit('update:open', value),
});

function handleNodeCreated(node: ChatbotNode) {
  isSheetOpen.value = false;
  emit('success', node);
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
        <SheetTitle>{{ t('features.chatbotNode.sheet.createTitle') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.chatbotNode.sheet.createDescription') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <NodeCreateForm
          @success="handleNodeCreated"
          @cancel="handleSheetClose"
        />
      </div>
      <SheetFooter />
    </SheetContent>
  </Sheet>
</template>
