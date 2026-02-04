<script setup lang="ts">
import type { Value } from '@bufbuild/protobuf/wkt';
import type { ChatbotNode, NodeEffect } from '~~/gen/chatbot/nodes/v1/node_pb';
import { ArrowRight, Info, Plus, Trash2, X } from 'lucide-vue-next';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';

// Simplified effect type for editing
interface EditableEffect {
  type: 'set_mode' | 'set_context' | 'goto' | '';
  target: string;
  context: { key: string; value: string }[];
}

interface Props {
  effect?: NodeEffect;
  disabled?: boolean;
  nodes?: ChatbotNode[]; // Available nodes for goto selection
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  nodes: () => [],
});

const emit = defineEmits<{
  'update:effect': [value: NodeEffect | undefined];
}>();

const { t } = useI18n();

// Effect type options
const effectTypeOptions = [
  { value: 'set_mode', label: t('features.chatbotNode.effect.setMode') },
  { value: 'set_context', label: t('features.chatbotNode.effect.setContext') },
  { value: 'goto', label: t('features.chatbotNode.effect.goto') },
];

// Mode options for set_mode effect
const modeOptions = [
  { value: 'assistant', label: 'AI Assistant' },
  { value: 'flow', label: 'Static Flow' },
  { value: 'liveChat', label: 'Live Chat' },
];

// Convert protobuf to editable format
function protoToEditable(effect: NodeEffect | undefined): EditableEffect | null {
  if (!effect || !effect.type)
    return null;

  const contextPairs: { key: string; value: string }[] = [];
  if (effect.context) {
    for (const [key, val] of Object.entries(effect.context)) {
      let valueStr = '';
      if (val) {
        // Handle protobuf-es Value structure: { kind: { case: 'stringValue', value: '...' } }
        const v = val as Value & { kind?: { case: string; value: unknown } };
        if (v.kind && v.kind.value !== undefined) {
          valueStr = String(v.kind.value);
        }
      }
      contextPairs.push({ key, value: valueStr });
    }
  }

  return {
    type: effect.type as EditableEffect['type'],
    target: effect.target || '',
    context: contextPairs.length > 0 ? contextPairs : [{ key: '', value: '' }],
  };
}

// Convert editable format back to protobuf-like structure
function editableToProto(editable: EditableEffect | null): NodeEffect | undefined {
  if (!editable || !editable.type)
    return undefined;

  const context: Record<string, unknown> = {};
  if (editable.type === 'set_context') {
    for (const pair of editable.context) {
      const trimmedKey = pair.key.trim();
      const trimmedValue = pair.value.trim();
      // Only include pairs with non-empty key AND non-empty value
      if (trimmedKey && trimmedValue) {
        // Use protobuf-es Value structure
        if (trimmedValue === 'true' || trimmedValue === 'false') {
          context[trimmedKey] = { kind: { case: 'boolValue', value: trimmedValue === 'true' } };
        }
        else if (!Number.isNaN(Number(trimmedValue))) {
          context[trimmedKey] = { kind: { case: 'numberValue', value: Number(trimmedValue) } };
        }
        else {
          context[trimmedKey] = { kind: { case: 'stringValue', value: trimmedValue } };
        }
      }
    }
  }

  return {
    type: editable.type,
    target: editable.target || '',
    context,
  } as unknown as NodeEffect;
}

// Local editable state
const editableEffect = ref<EditableEffect | null>(null);
const hasEffect = ref(false);

// Flag to skip watch after our own emits
let skipNextWatch = false;

// Initialize from prop
watch(
  () => props.effect,
  (newVal) => {
    if (skipNextWatch) {
      skipNextWatch = false;
      return;
    }
    const converted = protoToEditable(newVal);
    editableEffect.value = converted;
    hasEffect.value = converted !== null;
  },
  { immediate: true },
);

// Emit changes
function emitUpdate() {
  skipNextWatch = true;
  const proto = editableToProto(editableEffect.value);
  emit('update:effect', proto);
}

// Add a new effect
function addEffect() {
  editableEffect.value = {
    type: 'set_mode',
    target: 'assistant',
    context: [{ key: '', value: '' }],
  };
  hasEffect.value = true;
  emitUpdate();
}

// Remove effect entirely
function removeEffect() {
  editableEffect.value = null;
  hasEffect.value = false;
  emit('update:effect', undefined);
}

// Update effect type
function updateEffectType(value: string | number | bigint | Record<string, unknown> | null) {
  if (!editableEffect.value || value === null)
    return;
  editableEffect.value.type = String(value) as EditableEffect['type'];

  // Set sensible defaults based on type
  const strValue = String(value);
  if (strValue === 'set_mode') {
    editableEffect.value.target = 'assistant';
    editableEffect.value.context = [];
  }
  else if (strValue === 'goto') {
    editableEffect.value.target = '';
    editableEffect.value.context = [];
  }
  else if (strValue === 'set_context') {
    editableEffect.value.target = '';
    editableEffect.value.context = [{ key: '', value: '' }];
  }

  emitUpdate();
}

// Update target
function updateTarget(value: string | number | bigint | Record<string, unknown> | null) {
  if (!editableEffect.value || value === null)
    return;
  editableEffect.value.target = String(value);
  emitUpdate();
}

// Add context pair
function addContextPair() {
  if (!editableEffect.value)
    return;
  editableEffect.value.context.push({ key: '', value: '' });
}

// Remove context pair
function removeContextPair(index: number) {
  if (!editableEffect.value)
    return;
  editableEffect.value.context.splice(index, 1);
  if (editableEffect.value.context.length === 0) {
    editableEffect.value.context.push({ key: '', value: '' });
  }
  emitUpdate();
}

// Update context pair
function updateContextPair(index: number, field: 'key' | 'value', value: string) {
  if (!editableEffect.value || !editableEffect.value.context[index])
    return;
  editableEffect.value.context[index][field] = value;
  emitUpdate();
}

// Get available nodes deduped by name_lang (only show one per unique name+lang)
// This is because goto targets use node names, not IDs, and versioned nodes
// should appear as a single entry (the version is resolved at runtime)
const uniqueNodesByNameLang = computed(() => {
  const seen = new Set<string>();
  const result: ChatbotNode[] = [];

  for (const node of props.nodes) {
    if (!node.enabled)
      continue;

    const key = `${node.name}_${node.lang}`;
    if (!seen.has(key)) {
      seen.add(key);
      result.push(node);
    }
  }

  return result;
});

// Get URL to navigate to a node by name
function getNodeUrl(nodeName: string): string {
  // Find the node to get its lang
  const node = props.nodes.find(n => n.name === nodeName);
  if (node) {
    return `/platform/node/${node.name}_${node.lang}`;
  }
  return '';
}
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <Label class="text-base font-medium">
        <!-- {{ t('features.chatbotNode.effect.title') }} -->
      </Label>
      <div v-if="!hasEffect">
        <Button
          variant="outline"
          size="sm"
          :disabled="disabled"
          @click="addEffect"
        >
          <Plus class="h-4 w-4 mr-2" />
          {{ t('features.chatbotNode.effect.add') }}
        </Button>
      </div>
      <div v-else>
        <Button
          variant="ghost"
          size="sm"
          :disabled="disabled"
          @click="removeEffect"
        >
          <X class="h-4 w-4 mr-2" />
          {{ t('features.chatbotNode.effect.remove') }}
        </Button>
      </div>
    </div>

    <Alert
      variant="default"
      class="bg-blue-50 dark:bg-blue-950 border-blue-200 dark:border-blue-800"
    >
      <Info class="h-4 w-4 text-blue-600 dark:text-blue-400" />
      <AlertDescription class="text-blue-700 dark:text-blue-300">
        {{ t('features.chatbotNode.effect.description') }}
      </AlertDescription>
    </Alert>

    <!-- No effect set -->
    <div v-if="!hasEffect" class="text-sm text-muted-foreground italic py-4 text-center">
      {{ t('features.chatbotNode.effect.noEffect') }}
    </div>

    <!-- Effect editor -->
    <Card v-else-if="editableEffect">
      <CardContent class="pt-4 space-y-4">
        <!-- Effect type selector -->
        <div class="space-y-2">
          <Label>{{ t('features.chatbotNode.effect.type') }}</Label>
          <Select
            :model-value="editableEffect.type"
            :disabled="disabled"
            @update:model-value="updateEffectType"
          >
            <SelectTrigger class="w-[200px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="opt in effectTypeOptions"
                :key="opt.value"
                :value="opt.value"
              >
                {{ opt.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <!-- set_mode: Mode selector -->
        <div v-if="editableEffect.type === 'set_mode'" class="space-y-2">
          <Label>{{ t('features.chatbotNode.effect.targetMode') }}</Label>
          <Select
            :model-value="editableEffect.target"
            :disabled="disabled"
            @update:model-value="updateTarget"
          >
            <SelectTrigger class="w-[200px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="opt in modeOptions"
                :key="opt.value"
                :value="opt.value"
              >
                {{ opt.label }}
              </SelectItem>
            </SelectContent>
          </Select>
          <p class="text-xs text-muted-foreground">
            {{ t('features.chatbotNode.effect.setModeHelp') }}
          </p>
        </div>

        <!-- goto: Node selector -->
        <div v-if="editableEffect.type === 'goto'" class="space-y-2">
          <Label>{{ t('features.chatbotNode.effect.targetNode') }}</Label>
          <div class="flex items-center gap-2">
            <Select
              :model-value="editableEffect.target"
              :disabled="disabled"
              class="flex-1"
              @update:model-value="updateTarget"
            >
              <SelectTrigger class="w-full">
                <SelectValue :placeholder="t('features.chatbotNode.effect.selectNode')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="node in uniqueNodesByNameLang"
                  :key="`${node.name}_${node.lang}`"
                  :value="node.name"
                >
                  {{ node.name }}_{{ node.lang }}
                </SelectItem>
              </SelectContent>
            </Select>
            <Tooltip v-if="editableEffect.target && getNodeUrl(editableEffect.target)">
              <TooltipTrigger as-child>
                <Button
                  variant="outline"
                  size="icon"
                  as-child
                >
                  <NuxtLink :to="getNodeUrl(editableEffect.target)">
                    <ArrowRight class="h-4 w-4" />
                  </NuxtLink>
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                {{ t('features.chatbotNode.effect.openNode') }}
              </TooltipContent>
            </Tooltip>
          </div>
          <p class="text-xs text-muted-foreground">
            {{ t('features.chatbotNode.effect.gotoHelp') }}
          </p>
        </div>

        <!-- set_context: Key-value pairs -->
        <div v-if="editableEffect.type === 'set_context'" class="space-y-3">
          <Label>{{ t('features.chatbotNode.effect.contextPairs') }}</Label>
          <p class="text-xs text-muted-foreground">
            {{ t('features.chatbotNode.effect.setContextHelp') }}
          </p>

          <div class="space-y-2">
            <div
              v-for="(pair, index) in editableEffect.context"
              :key="index"
              class="flex items-center gap-2"
            >
              <Input
                :model-value="pair.key"
                :placeholder="t('features.chatbotNode.effect.keyPlaceholder')"
                :disabled="disabled"
                class="flex-1"
                @update:model-value="
                  (v: string | number) => updateContextPair(index, 'key', String(v))
                "
              />
              <span class="text-muted-foreground">=</span>
              <Input
                :model-value="pair.value"
                :placeholder="t('features.chatbotNode.effect.valuePlaceholder')"
                :disabled="disabled"
                class="flex-1"
                @update:model-value="
                  (v: string | number) => updateContextPair(index, 'value', String(v))
                "
              />
              <Button
                variant="ghost"
                size="icon"
                :disabled="disabled"
                @click="removeContextPair(index)"
              >
                <Trash2 class="h-4 w-4 text-destructive" />
              </Button>
            </div>
          </div>

          <Button
            variant="outline"
            size="sm"
            :disabled="disabled"
            @click="addContextPair"
          >
            <Plus class="h-4 w-4 mr-2" />
            {{ t('features.chatbotNode.effect.addPair') }}
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
