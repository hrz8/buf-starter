<script setup lang="ts">
import type { ChatbotNodeMessage } from '~~/gen/chatbot/nodes/v1/node_pb';
import { create } from '@bufbuild/protobuf';
import { Plus } from 'lucide-vue-next';
import { ChatbotNodeMessageSchema } from '~~/gen/chatbot/nodes/v1/node_pb';
import { Button } from '@/components/ui/button';
import { DEFAULT_MESSAGE } from '.';
import MessageBubble from './MessageBubble.vue';

interface Props {
  messages: ChatbotNodeMessage[];
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
});

const emit = defineEmits<{
  'update:messages': [messages: ChatbotNodeMessage[]];
}>();

const { t } = useI18n();

function updateMessage(index: number, content: string) {
  const newMessages = [...props.messages];
  const existingMessage = newMessages[index];
  if (existingMessage) {
    newMessages[index] = create(ChatbotNodeMessageSchema, {
      role: existingMessage.role,
      content,
    });
  }
  emit('update:messages', newMessages);
}

function removeMessage(index: number) {
  if (props.messages.length <= 1) {
    return;
  }
  const newMessages = props.messages.filter((_, i) => i !== index);
  emit('update:messages', newMessages);
}

function moveMessage(index: number, direction: 'up' | 'down') {
  const newIndex = direction === 'up' ? index - 1 : index + 1;
  if (newIndex < 0 || newIndex >= props.messages.length) {
    return;
  }

  const newMessages = [...props.messages];
  const current = newMessages[index];
  const target = newMessages[newIndex];
  if (current && target) {
    newMessages[index] = target;
    newMessages[newIndex] = current;
  }
  emit('update:messages', newMessages);
}

function addMessage() {
  const newMessage = create(ChatbotNodeMessageSchema, {
    role: DEFAULT_MESSAGE.role,
    content: DEFAULT_MESSAGE.content,
  });
  emit('update:messages', [...props.messages, newMessage]);
}
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-medium text-foreground">
        {{ t('features.chatbotNode.message.title') }}
      </h3>
      <span class="text-xs text-muted-foreground">
        {{ messages.length }} {{ t('features.chatbotNode.message.count', messages.length) }}
      </span>
    </div>

    <p class="text-sm text-muted-foreground">
      {{ t('features.chatbotNode.message.description') }}
    </p>

    <!-- Messages as chat bubbles -->
    <div class="space-y-3 pl-4">
      <MessageBubble
        v-for="(message, index) in messages"
        :key="index"
        :content="message.content"
        :index="index"
        :total-count="messages.length"
        :disabled="disabled"
        @update:content="updateMessage(index, $event)"
        @remove="removeMessage(index)"
        @move-up="moveMessage(index, 'up')"
        @move-down="moveMessage(index, 'down')"
      />
    </div>

    <!-- Add message button -->
    <Button
      variant="outline"
      size="sm"
      class="w-full border-dashed"
      :disabled="disabled"
      @click="addMessage"
    >
      <Plus class="h-4 w-4 mr-2" />
      {{ t('features.chatbotNode.message.add') }}
    </Button>
  </div>
</template>
