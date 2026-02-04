<script setup lang="ts">
import type { ChatbotNodeTrigger } from '~~/gen/chatbot/nodes/v1/node_pb';
import { create } from '@bufbuild/protobuf';
import { Plus, Trash2 } from 'lucide-vue-next';
import { ChatbotNodeTriggerSchema } from '~~/gen/chatbot/nodes/v1/node_pb';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  DEFAULT_TRIGGER,
  isValidRegex,
  TRIGGER_TYPE_DESCRIPTIONS,
  TRIGGER_TYPE_LABELS,
  TRIGGER_TYPES,
  VALIDATION_LIMITS,
} from '.';

interface Props {
  triggers: ChatbotNodeTrigger[];
  disabled?: boolean;
  showHeader?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  showHeader: true,
});

const emit = defineEmits<{
  'update:triggers': [triggers: ChatbotNodeTrigger[]];
}>();

const { t } = useI18n();

// Track regex validation errors for each trigger
const regexErrors = ref<Record<number, string>>({});

function validateRegex(index: number, value: string, type: string) {
  if (type === 'regex' && value && !isValidRegex(value)) {
    regexErrors.value[index] = t('features.chatbotNode.trigger.invalidRegex');
  }
  else {
    delete regexErrors.value[index];
  }
}

function updateTriggerType(index: number, type: string) {
  const newTriggers = [...props.triggers];
  const existingTrigger = newTriggers[index];
  if (existingTrigger) {
    newTriggers[index] = create(ChatbotNodeTriggerSchema, {
      type,
      value: existingTrigger.value,
    });
    validateRegex(index, existingTrigger.value, type);
  }
  emit('update:triggers', newTriggers);
}

function updateTriggerValue(index: number, value: string) {
  const newTriggers = [...props.triggers];
  const existingTrigger = newTriggers[index];
  if (existingTrigger) {
    newTriggers[index] = create(ChatbotNodeTriggerSchema, {
      type: existingTrigger.type,
      value,
    });
    validateRegex(index, value, existingTrigger.type);
  }
  emit('update:triggers', newTriggers);
}

function removeTrigger(index: number) {
  // Allow removing all triggers - nodes can match based on conditions only
  const newTriggers = props.triggers.filter((_, i) => i !== index);
  // Re-index regex errors
  const newErrors: Record<number, string> = {};
  Object.entries(regexErrors.value).forEach(([key, val]) => {
    const oldIndex = Number.parseInt(key, 10);
    if (oldIndex < index) {
      newErrors[oldIndex] = val;
    }
    else if (oldIndex > index) {
      newErrors[oldIndex - 1] = val;
    }
  });
  regexErrors.value = newErrors;
  emit('update:triggers', newTriggers);
}

function addTrigger() {
  const newTrigger = create(ChatbotNodeTriggerSchema, {
    type: DEFAULT_TRIGGER.type,
    value: DEFAULT_TRIGGER.value,
  });
  emit('update:triggers', [...props.triggers, newTrigger]);
}
</script>

<template>
  <div class="space-y-4">
    <template v-if="showHeader">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-medium text-foreground">
          {{ t('features.chatbotNode.trigger.title') }}
        </h3>
        <span class="text-xs text-muted-foreground">
          {{ triggers.length }} {{ t('features.chatbotNode.trigger.count', triggers.length) }}
        </span>
      </div>

      <p class="text-sm text-muted-foreground">
        {{ t('features.chatbotNode.trigger.description') }}
      </p>
    </template>

    <!-- Trigger items -->
    <div class="space-y-3">
      <div
        v-for="(trigger, index) in triggers"
        :key="index"
        class="flex items-start gap-2 p-3 rounded-lg border bg-muted/30"
      >
        <!-- Type select -->
        <div class="w-32 shrink-0">
          <Select
            :model-value="trigger.type"
            :disabled="disabled"
            @update:model-value="updateTriggerType(index, $event as string)"
          >
            <SelectTrigger class="h-9">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="triggerType in TRIGGER_TYPES"
                :key="triggerType"
                :value="triggerType"
              >
                {{ TRIGGER_TYPE_LABELS[triggerType] }}
              </SelectItem>
            </SelectContent>
          </Select>
          <p class="text-xs text-muted-foreground mt-1">
            {{ TRIGGER_TYPE_DESCRIPTIONS[trigger.type as keyof typeof TRIGGER_TYPE_DESCRIPTIONS] }}
          </p>
        </div>

        <!-- Value input -->
        <div class="flex-1 space-y-1">
          <Input
            :model-value="trigger.value"
            :placeholder="t('features.chatbotNode.trigger.valuePlaceholder')"
            :disabled="disabled"
            :maxlength="VALIDATION_LIMITS.triggerValueMaxLength"
            class="h-9"
            @update:model-value="updateTriggerValue(index, $event as string)"
          />
          <p
            v-if="regexErrors[index]"
            class="text-xs text-destructive"
          >
            {{ regexErrors[index] }}
          </p>
          <p class="text-xs text-muted-foreground text-right">
            {{ trigger.value.length }} / {{ VALIDATION_LIMITS.triggerValueMaxLength }}
          </p>
        </div>

        <!-- Remove button - always show since triggers are now optional -->
        <Button
          variant="ghost"
          size="icon"
          class="h-9 w-9 text-destructive hover:text-destructive shrink-0"
          :disabled="disabled"
          @click="removeTrigger(index)"
        >
          <Trash2 class="h-4 w-4" />
        </Button>
      </div>
    </div>

    <!-- Add trigger button -->
    <Button
      variant="outline"
      size="sm"
      class="w-full border-dashed"
      :disabled="disabled"
      @click="addTrigger"
    >
      <Plus class="h-4 w-4 mr-2" />
      {{ t('features.chatbotNode.trigger.add') }}
    </Button>
  </div>
</template>
