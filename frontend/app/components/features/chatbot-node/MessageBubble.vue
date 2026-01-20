<script setup lang="ts">
import { ChevronDown, ChevronUp, GripVertical, Trash2 } from 'lucide-vue-next';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { VALIDATION_LIMITS } from '.';

interface Props {
  content: string;
  index: number;
  totalCount: number;
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
});

const emit = defineEmits<{
  'update:content': [content: string];
  'remove': [];
  'moveUp': [];
  'moveDown': [];
}>();

const { t } = useI18n();

const localContent = computed({
  get: () => props.content,
  set: (value: string) => emit('update:content', value),
});

const charCount = computed(() => props.content.length);
const isOverLimit = computed(() => charCount.value > VALIDATION_LIMITS.messageContentMaxLength);

function handleRemove() {
  emit('remove');
}

function handleMoveUp() {
  emit('moveUp');
}

function handleMoveDown() {
  emit('moveDown');
}
</script>

<template>
  <div class="group relative flex gap-2">
    <!-- Drag handle and reorder buttons -->
    <div
      class="flex flex-col items-center justify-center gap-1
        opacity-0 group-hover:opacity-100 transition-opacity"
    >
      <GripVertical class="h-4 w-4 text-muted-foreground cursor-grab" />
      <Button
        variant="ghost"
        size="icon"
        class="h-6 w-6"
        :disabled="disabled || index === 0"
        @click="handleMoveUp"
      >
        <ChevronUp class="h-3 w-3" />
      </Button>
      <Button
        variant="ghost"
        size="icon"
        class="h-6 w-6"
        :disabled="disabled || index === totalCount - 1"
        @click="handleMoveDown"
      >
        <ChevronDown class="h-3 w-3" />
      </Button>
    </div>

    <!-- Chat bubble -->
    <div class="flex-1 max-w-[85%]">
      <div class="bg-primary/10 rounded-2xl rounded-tl-sm p-4 shadow-sm">
        <Textarea
          v-model="localContent"
          :placeholder="t('features.chatbotNode.message.placeholder')"
          :disabled="disabled"
          class="min-h-[80px] bg-transparent border-0 p-0 resize-none
            focus-visible:ring-0 focus-visible:ring-offset-0"
          rows="3"
        />
        <div class="flex items-center justify-between mt-2 text-xs text-muted-foreground">
          <span :class="{ 'text-destructive': isOverLimit }">
            {{ charCount }} / {{ VALIDATION_LIMITS.messageContentMaxLength }}
          </span>
          <Button
            v-if="totalCount > 1"
            variant="ghost"
            size="icon"
            class="h-6 w-6 text-destructive hover:text-destructive"
            :disabled="disabled"
            @click="handleRemove"
          >
            <Trash2 class="h-3 w-3" />
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>
