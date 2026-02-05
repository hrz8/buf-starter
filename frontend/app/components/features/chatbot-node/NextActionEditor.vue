<script setup lang="ts">
import type { Value } from '@bufbuild/protobuf/wkt';
import type {
  ChatbotNode,
  NodeCapture,
  NodeCaptureValidation,
  NodeNextAction,
} from '~~/gen/chatbot/nodes/v1/node_pb';
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
import { Separator } from '@/components/ui/separator';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';

// Simplified capture validation for editing
interface EditableCaptureValidation {
  type: 'string' | 'email' | 'phone' | 'number' | 'regex' | 'options' | '';
  pattern: string;
  options: string[];
  minLength: number;
  maxLength: number;
  min: number;
  max: number;
}

// Simplified capture config for editing
interface EditableCapture {
  variableName: string;
  validation: EditableCaptureValidation;
  onFailMessage: string;
  maxRetries: number;
  onSuccessContext: { key: string; value: string }[];
  onSuccessGoto: string;
  onFailGoto: string;
}

// Simplified next action type for editing
interface EditableNextAction {
  type: 'goto' | 'capture' | '';
  target: string;
  capture: EditableCapture;
}

interface Props {
  nextAction?: NodeNextAction;
  disabled?: boolean;
  nodes?: ChatbotNode[]; // Available nodes for goto selection
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  nodes: () => [],
});

const emit = defineEmits<{
  'update:nextAction': [value: NodeNextAction | undefined];
}>();

const { t } = useI18n();

// Action type options
const actionTypeOptions = [
  { value: 'goto', label: t('features.chatbotNode.nextAction.goto') },
  { value: 'capture', label: t('features.chatbotNode.nextAction.capture') },
];

// Validation type options
const validationTypeOptions = [
  { value: 'string', label: t('features.chatbotNode.capture.validationTypeString') },
  { value: 'email', label: t('features.chatbotNode.capture.validationTypeEmail') },
  { value: 'phone', label: t('features.chatbotNode.capture.validationTypePhone') },
  { value: 'number', label: t('features.chatbotNode.capture.validationTypeNumber') },
  { value: 'regex', label: t('features.chatbotNode.capture.validationTypeRegex') },
  { value: 'options', label: t('features.chatbotNode.capture.validationTypeOptions') },
];

// Create default editable capture
function createDefaultCapture(): EditableCapture {
  return {
    variableName: '',
    validation: {
      type: 'string',
      pattern: '',
      options: [],
      minLength: 0,
      maxLength: 0,
      min: 0,
      max: 0,
    },
    onFailMessage: '',
    maxRetries: 0,
    onSuccessContext: [{ key: '', value: '' }],
    onSuccessGoto: '',
    onFailGoto: '',
  };
}

// Convert protobuf validation to editable format
function protoValidationToEditable(
  validation: NodeCaptureValidation | undefined,
): EditableCaptureValidation {
  if (!validation) {
    return {
      type: 'string',
      pattern: '',
      options: [],
      minLength: 0,
      maxLength: 0,
      min: 0,
      max: 0,
    };
  }
  return {
    type: (validation.type || 'string') as EditableCaptureValidation['type'],
    pattern: validation.pattern || '',
    options: [...(validation.options || [])],
    minLength: validation.minLength || 0,
    maxLength: validation.maxLength || 0,
    min: validation.min || 0,
    max: validation.max || 0,
  };
}

// Convert protobuf capture to editable format
function protoCaptureToEditable(capture: NodeCapture | undefined): EditableCapture {
  if (!capture) {
    return createDefaultCapture();
  }

  const contextPairs: { key: string; value: string }[] = [];
  if (capture.onSuccessContext) {
    for (const [key, val] of Object.entries(capture.onSuccessContext)) {
      let valueStr = '';
      if (val) {
        const v = val as Value & { kind?: { case: string; value: unknown } };
        if (v.kind && v.kind.value !== undefined) {
          valueStr = String(v.kind.value);
        }
      }
      contextPairs.push({ key, value: valueStr });
    }
  }

  return {
    variableName: capture.variableName || '',
    validation: protoValidationToEditable(capture.validation),
    onFailMessage: capture.onFailMessage || '',
    maxRetries: capture.maxRetries || 0,
    onSuccessContext: contextPairs.length > 0 ? contextPairs : [{ key: '', value: '' }],
    onSuccessGoto: capture.onSuccessGoto || '',
    onFailGoto: capture.onFailGoto || '',
  };
}

// Convert protobuf to editable format
function protoToEditable(action: NodeNextAction | undefined): EditableNextAction | null {
  if (!action || !action.type)
    return null;

  return {
    type: action.type as EditableNextAction['type'],
    target: action.target || '',
    capture: protoCaptureToEditable(action.capture),
  };
}

// Convert editable validation to protobuf format
function editableValidationToProto(
  validation: EditableCaptureValidation,
): NodeCaptureValidation | undefined {
  if (!validation.type)
    return undefined;

  return {
    type: validation.type,
    pattern: validation.pattern || '',
    options: validation.options.filter(o => o.trim() !== ''),
    minLength: validation.minLength || 0,
    maxLength: validation.maxLength || 0,
    min: validation.min || 0,
    max: validation.max || 0,
  } as unknown as NodeCaptureValidation;
}

// Convert editable capture to protobuf format
function editableCaptureToProto(capture: EditableCapture): NodeCapture | undefined {
  if (!capture.variableName.trim())
    return undefined;

  const onSuccessContext: Record<string, unknown> = {};
  for (const pair of capture.onSuccessContext) {
    const trimmedKey = pair.key.trim();
    const trimmedValue = pair.value.trim();
    if (trimmedKey && trimmedValue) {
      if (trimmedValue === 'true' || trimmedValue === 'false') {
        onSuccessContext[trimmedKey] = {
          kind: { case: 'boolValue', value: trimmedValue === 'true' },
        };
      }
      else if (!Number.isNaN(Number(trimmedValue))) {
        onSuccessContext[trimmedKey] = {
          kind: { case: 'numberValue', value: Number(trimmedValue) },
        };
      }
      else {
        onSuccessContext[trimmedKey] = {
          kind: { case: 'stringValue', value: trimmedValue },
        };
      }
    }
  }

  return {
    variableName: capture.variableName.trim(),
    validation: editableValidationToProto(capture.validation),
    onFailMessage: capture.onFailMessage || '',
    maxRetries: capture.maxRetries || 0,
    onSuccessContext,
    onSuccessGoto: capture.onSuccessGoto || '',
    onFailGoto: capture.onFailGoto || '',
  } as unknown as NodeCapture;
}

// Convert editable format back to protobuf-like structure
function editableToProto(editable: EditableNextAction | null): NodeNextAction | undefined {
  if (!editable || !editable.type)
    return undefined;

  const result: Partial<NodeNextAction> = {
    type: editable.type,
    target: editable.target || '',
  };

  if (editable.type === 'capture') {
    result.capture = editableCaptureToProto(editable.capture);
  }

  return result as unknown as NodeNextAction;
}

// Local editable state
const editableAction = ref<EditableNextAction | null>(null);
const hasAction = ref(false);

// Flag to skip watch after our own emits
let skipNextWatch = false;

// Initialize from prop
watch(
  () => props.nextAction,
  (newVal) => {
    if (skipNextWatch) {
      skipNextWatch = false;
      return;
    }
    const converted = protoToEditable(newVal);
    editableAction.value = converted;
    hasAction.value = converted !== null;
  },
  { immediate: true },
);

// Emit changes
function emitUpdate() {
  skipNextWatch = true;
  const proto = editableToProto(editableAction.value);
  emit('update:nextAction', proto);
}

// Add a new action
function addAction() {
  editableAction.value = {
    type: 'goto',
    target: '',
    capture: createDefaultCapture(),
  };
  hasAction.value = true;
  emitUpdate();
}

// Remove action entirely
function removeAction() {
  editableAction.value = null;
  hasAction.value = false;
  emit('update:nextAction', undefined);
}

// Update action type
function updateActionType(value: string | number | bigint | Record<string, unknown> | null) {
  if (!editableAction.value || value === null)
    return;
  const strValue = String(value);
  editableAction.value.type = strValue as EditableNextAction['type'];

  // Set sensible defaults based on type
  if (strValue === 'goto') {
    editableAction.value.target = '';
    editableAction.value.capture = createDefaultCapture();
  }
  else if (strValue === 'capture') {
    editableAction.value.target = '';
    editableAction.value.capture = createDefaultCapture();
  }

  emitUpdate();
}

// Update target
function updateTarget(value: string | number | bigint | Record<string, unknown> | null) {
  if (!editableAction.value || value === null)
    return;
  editableAction.value.target = String(value);
  emitUpdate();
}

// Update capture field
function updateCaptureField<K extends keyof EditableCapture>(field: K, value: EditableCapture[K]) {
  if (!editableAction.value)
    return;
  editableAction.value.capture[field] = value;
  emitUpdate();
}

// Update validation field
function updateValidationField<K extends keyof EditableCaptureValidation>(
  field: K,
  value: EditableCaptureValidation[K],
) {
  if (!editableAction.value)
    return;
  editableAction.value.capture.validation[field] = value;
  emitUpdate();
}

// Add option (for options validation type)
function addOption() {
  if (!editableAction.value)
    return;
  editableAction.value.capture.validation.options.push('');
}

// Remove option
function removeOption(index: number) {
  if (!editableAction.value)
    return;
  editableAction.value.capture.validation.options.splice(index, 1);
  emitUpdate();
}

// Update option
function updateOption(index: number, value: string) {
  if (!editableAction.value)
    return;
  editableAction.value.capture.validation.options[index] = value;
  emitUpdate();
}

// Add success context pair
function addSuccessContextPair() {
  if (!editableAction.value)
    return;
  editableAction.value.capture.onSuccessContext.push({ key: '', value: '' });
}

// Remove success context pair
function removeSuccessContextPair(index: number) {
  if (!editableAction.value)
    return;
  editableAction.value.capture.onSuccessContext.splice(index, 1);
  if (editableAction.value.capture.onSuccessContext.length === 0) {
    editableAction.value.capture.onSuccessContext.push({ key: '', value: '' });
  }
  emitUpdate();
}

// Update success context pair
function updateSuccessContextPair(index: number, field: 'key' | 'value', value: string) {
  if (!editableAction.value || !editableAction.value.capture.onSuccessContext[index])
    return;
  editableAction.value.capture.onSuccessContext[index][field] = value;
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
  // Find the node to get its lang and version
  const node = props.nodes.find(n => n.name === nodeName);
  if (node) {
    const params = new URLSearchParams();
    if (node.lang)
      params.set('lang', node.lang);
    if (node.version)
      params.set('v', node.version);
    const queryString = params.toString();
    return `/platform/node/${node.name}${queryString ? `?${queryString}` : ''}`;
  }
  return '';
}

// Check if current validation type needs additional fields
const showPatternField = computed(() => editableAction.value?.capture.validation.type === 'regex');
const showOptionsField = computed(
  () => editableAction.value?.capture.validation.type === 'options',
);
const showLengthFields = computed(() => editableAction.value?.capture.validation.type === 'string');
const showNumberFields = computed(() => editableAction.value?.capture.validation.type === 'number');
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <Label class="text-base font-medium">
        <!-- {{ t('features.chatbotNode.nextAction.title') }} -->
      </Label>
      <div v-if="!hasAction">
        <Button
          variant="outline"
          size="sm"
          :disabled="disabled"
          @click="addAction"
        >
          <Plus class="h-4 w-4 mr-2" />
          {{ t('features.chatbotNode.nextAction.add') }}
        </Button>
      </div>
      <div v-else>
        <Button
          variant="ghost"
          size="sm"
          :disabled="disabled"
          @click="removeAction"
        >
          <X class="h-4 w-4 mr-2" />
          {{ t('features.chatbotNode.nextAction.remove') }}
        </Button>
      </div>
    </div>

    <Alert
      variant="default"
      class="bg-amber-50 dark:bg-amber-950 border-amber-200 dark:border-amber-800"
    >
      <Info class="h-4 w-4 text-amber-600 dark:text-amber-400" />
      <AlertDescription class="text-amber-700 dark:text-amber-300">
        {{ t('features.chatbotNode.nextAction.description') }}
      </AlertDescription>
    </Alert>

    <!-- No action set -->
    <div v-if="!hasAction" class="text-sm text-muted-foreground italic py-4 text-center">
      {{ t('features.chatbotNode.nextAction.noAction') }}
    </div>

    <!-- Action editor -->
    <Card v-else-if="editableAction">
      <CardContent class="pt-4 space-y-4">
        <!-- Action type selector -->
        <div class="space-y-2">
          <Label>{{ t('features.chatbotNode.nextAction.actionType') }}</Label>
          <Select
            :model-value="editableAction.type"
            :disabled="disabled"
            @update:model-value="updateActionType"
          >
            <SelectTrigger class="w-[200px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="opt in actionTypeOptions"
                :key="opt.value"
                :value="opt.value"
              >
                {{ opt.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <!-- goto: Node selector -->
        <div v-if="editableAction.type === 'goto'" class="space-y-2">
          <Label>{{ t('features.chatbotNode.nextAction.targetNode') }}</Label>
          <div class="flex items-center gap-2">
            <Select
              :model-value="editableAction.target"
              :disabled="disabled"
              class="flex-1"
              @update:model-value="updateTarget"
            >
              <SelectTrigger class="w-full">
                <SelectValue :placeholder="t('features.chatbotNode.nextAction.selectNode')" />
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
            <Tooltip v-if="editableAction.target && getNodeUrl(editableAction.target)">
              <TooltipTrigger as-child>
                <Button
                  variant="outline"
                  size="icon"
                  as-child
                >
                  <NuxtLink :to="getNodeUrl(editableAction.target)">
                    <ArrowRight class="h-4 w-4" />
                  </NuxtLink>
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                {{ t('features.chatbotNode.nextAction.openNode') }}
              </TooltipContent>
            </Tooltip>
          </div>
          <p class="text-xs text-muted-foreground">
            {{ t('features.chatbotNode.nextAction.gotoHelp') }}
          </p>
        </div>

        <!-- capture: Full capture configuration -->
        <div v-if="editableAction.type === 'capture'" class="space-y-4">
          <!-- Variable Name -->
          <div class="space-y-2">
            <Label>{{ t('features.chatbotNode.capture.variableName') }} *</Label>
            <Input
              :model-value="editableAction.capture.variableName"
              :placeholder="t('features.chatbotNode.capture.variableNamePlaceholder')"
              :disabled="disabled"
              @update:model-value="
                (v: string | number) => updateCaptureField('variableName', String(v))
              "
            />
            <p class="text-xs text-muted-foreground">
              {{ t('features.chatbotNode.capture.variableNameHelp') }}
            </p>
          </div>

          <Separator />

          <!-- Validation Section -->
          <div class="space-y-3">
            <Label class="text-sm font-medium">
              {{ t('features.chatbotNode.capture.validation') }}
            </Label>

            <!-- Validation Type -->
            <div class="space-y-2">
              <Label class="text-xs">{{ t('features.chatbotNode.capture.validationType') }}</Label>
              <Select
                :model-value="editableAction.capture.validation.type"
                :disabled="disabled"
                @update:model-value="
                  (v: string | number | bigint | Record<string, unknown> | null) => {
                    if (v !== null)
                      updateValidationField('type', String(v) as EditableCaptureValidation['type'])
                  }
                "
              >
                <SelectTrigger class="w-[200px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem
                    v-for="opt in validationTypeOptions"
                    :key="opt.value"
                    :value="opt.value"
                  >
                    {{ opt.label }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <!-- Regex pattern field -->
            <div v-if="showPatternField" class="space-y-2">
              <Label class="text-xs">{{ t('features.chatbotNode.capture.pattern') }}</Label>
              <Input
                :model-value="editableAction.capture.validation.pattern"
                :placeholder="t('features.chatbotNode.capture.patternPlaceholder')"
                :disabled="disabled"
                @update:model-value="
                  (v: string | number) => updateValidationField('pattern', String(v))
                "
              />
              <p class="text-xs text-muted-foreground">
                {{ t('features.chatbotNode.capture.patternHelp') }}
              </p>
            </div>

            <!-- Options field -->
            <div v-if="showOptionsField" class="space-y-2">
              <Label class="text-xs">{{ t('features.chatbotNode.capture.options') }}</Label>
              <div class="space-y-2">
                <div
                  v-for="(opt, index) in editableAction.capture.validation.options"
                  :key="index"
                  class="flex items-center gap-2"
                >
                  <Input
                    :model-value="opt"
                    :placeholder="t('features.chatbotNode.capture.optionPlaceholder')"
                    :disabled="disabled"
                    class="flex-1"
                    @update:model-value="(v: string | number) => updateOption(index, String(v))"
                  />
                  <Button
                    variant="ghost"
                    size="icon"
                    :disabled="disabled"
                    @click="removeOption(index)"
                  >
                    <Trash2 class="h-4 w-4 text-destructive" />
                  </Button>
                </div>
              </div>
              <Button
                variant="outline"
                size="sm"
                :disabled="disabled"
                @click="addOption"
              >
                <Plus class="h-4 w-4 mr-2" />
                {{ t('features.chatbotNode.capture.addOption') }}
              </Button>
            </div>

            <!-- String length fields -->
            <div v-if="showLengthFields" class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label class="text-xs">{{ t('features.chatbotNode.capture.minLength') }}</Label>
                <Input
                  :model-value="editableAction.capture.validation.minLength"
                  type="number"
                  min="0"
                  :disabled="disabled"
                  @update:model-value="
                    (v: string | number) => updateValidationField('minLength', Number(v))
                  "
                />
              </div>
              <div class="space-y-2">
                <Label class="text-xs">{{ t('features.chatbotNode.capture.maxLength') }}</Label>
                <Input
                  :model-value="editableAction.capture.validation.maxLength"
                  type="number"
                  min="0"
                  :disabled="disabled"
                  @update:model-value="
                    (v: string | number) => updateValidationField('maxLength', Number(v))
                  "
                />
              </div>
            </div>

            <!-- Number range fields -->
            <div v-if="showNumberFields" class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label class="text-xs">{{ t('features.chatbotNode.capture.min') }}</Label>
                <Input
                  :model-value="editableAction.capture.validation.min"
                  type="number"
                  :disabled="disabled"
                  @update:model-value="
                    (v: string | number) => updateValidationField('min', Number(v))
                  "
                />
              </div>
              <div class="space-y-2">
                <Label class="text-xs">{{ t('features.chatbotNode.capture.max') }}</Label>
                <Input
                  :model-value="editableAction.capture.validation.max"
                  type="number"
                  :disabled="disabled"
                  @update:model-value="
                    (v: string | number) => updateValidationField('max', Number(v))
                  "
                />
              </div>
            </div>
          </div>

          <Separator />

          <!-- On Failure Section -->
          <div class="space-y-3">
            <Label class="text-sm font-medium">
              {{ t('features.chatbotNode.capture.onFailureSection') }}
            </Label>

            <div class="space-y-2">
              <Label class="text-xs">{{ t('features.chatbotNode.capture.onFailMessage') }}</Label>
              <Input
                :model-value="editableAction.capture.onFailMessage"
                :placeholder="t('features.chatbotNode.capture.onFailMessagePlaceholder')"
                :disabled="disabled"
                @update:model-value="
                  (v: string | number) => updateCaptureField('onFailMessage', String(v))
                "
              />
              <p class="text-xs text-muted-foreground">
                {{ t('features.chatbotNode.capture.onFailMessageHelp') }}
              </p>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label class="text-xs">{{ t('features.chatbotNode.capture.maxRetries') }}</Label>
                <Input
                  :model-value="editableAction.capture.maxRetries"
                  type="number"
                  min="0"
                  :disabled="disabled"
                  @update:model-value="
                    (v: string | number) => updateCaptureField('maxRetries', Number(v))
                  "
                />
                <p class="text-xs text-muted-foreground">
                  {{ t('features.chatbotNode.capture.maxRetriesHelp') }}
                </p>
              </div>

              <div class="space-y-2">
                <Label class="text-xs">{{ t('features.chatbotNode.capture.onFailGoto') }}</Label>
                <div class="flex items-center gap-2">
                  <Select
                    :model-value="editableAction.capture.onFailGoto || '__none__'"
                    :disabled="disabled"
                    class="flex-1"
                    @update:model-value="
                      (v: string | number | bigint | Record<string, unknown> | null) => {
                        const str = String(v ?? '')
                        updateCaptureField('onFailGoto', str === '__none__' ? '' : str)
                      }
                    "
                  >
                    <SelectTrigger class="w-full">
                      <SelectValue :placeholder="t('features.chatbotNode.capture.selectNode')" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="__none__">
                        {{ t('features.chatbotNode.capture.noNode') }}
                      </SelectItem>
                      <SelectItem
                        v-for="node in uniqueNodesByNameLang"
                        :key="`${node.name}_${node.lang}`"
                        :value="node.name"
                      >
                        {{ node.name }}_{{ node.lang }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <Tooltip
                    v-if="editableAction.capture.onFailGoto
                      && getNodeUrl(editableAction.capture.onFailGoto)"
                  >
                    <TooltipTrigger as-child>
                      <Button
                        variant="outline"
                        size="icon"
                        as-child
                      >
                        <NuxtLink :to="getNodeUrl(editableAction.capture.onFailGoto)">
                          <ArrowRight class="h-4 w-4" />
                        </NuxtLink>
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      {{ t('features.chatbotNode.capture.openNode') }}
                    </TooltipContent>
                  </Tooltip>
                </div>
                <p class="text-xs text-muted-foreground">
                  {{ t('features.chatbotNode.capture.onFailGotoHelp') }}
                </p>
              </div>
            </div>
          </div>

          <Separator />

          <!-- On Success Section -->
          <div class="space-y-3">
            <Label class="text-sm font-medium">
              {{ t('features.chatbotNode.capture.onSuccessSection') }}
            </Label>

            <div class="space-y-2">
              <Label class="text-xs">
                {{ t('features.chatbotNode.capture.onSuccessContext') }}
              </Label>
              <p class="text-xs text-muted-foreground">
                {{ t('features.chatbotNode.capture.onSuccessContextHelp') }}
              </p>

              <div class="space-y-2">
                <div
                  v-for="(pair, index) in editableAction.capture.onSuccessContext"
                  :key="index"
                  class="flex items-center gap-2"
                >
                  <Input
                    :model-value="pair.key"
                    :placeholder="t('features.chatbotNode.capture.keyPlaceholder')"
                    :disabled="disabled"
                    class="flex-1"
                    @update:model-value="
                      (v: string | number) => updateSuccessContextPair(index, 'key', String(v))
                    "
                  />
                  <span class="text-muted-foreground">=</span>
                  <Input
                    :model-value="pair.value"
                    :placeholder="t('features.chatbotNode.capture.valuePlaceholder')"
                    :disabled="disabled"
                    class="flex-1"
                    @update:model-value="
                      (v: string | number) => updateSuccessContextPair(index, 'value', String(v))
                    "
                  />
                  <Button
                    variant="ghost"
                    size="icon"
                    :disabled="disabled"
                    @click="removeSuccessContextPair(index)"
                  >
                    <Trash2 class="h-4 w-4 text-destructive" />
                  </Button>
                </div>
              </div>

              <Button
                variant="outline"
                size="sm"
                :disabled="disabled"
                @click="addSuccessContextPair"
              >
                <Plus class="h-4 w-4 mr-2" />
                {{ t('features.chatbotNode.capture.addContextPair') }}
              </Button>
            </div>

            <div class="space-y-2">
              <Label class="text-xs">{{ t('features.chatbotNode.capture.onSuccessGoto') }}</Label>
              <div class="flex items-center gap-2">
                <Select
                  :model-value="editableAction.capture.onSuccessGoto || '__none__'"
                  :disabled="disabled"
                  class="flex-1"
                  @update:model-value="
                    (v: string | number | bigint | Record<string, unknown> | null) => {
                      const str = String(v ?? '')
                      updateCaptureField('onSuccessGoto', str === '__none__' ? '' : str)
                    }
                  "
                >
                  <SelectTrigger class="w-full">
                    <SelectValue :placeholder="t('features.chatbotNode.capture.selectNode')" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="__none__">
                      {{ t('features.chatbotNode.capture.noNode') }}
                    </SelectItem>
                    <SelectItem
                      v-for="node in uniqueNodesByNameLang"
                      :key="`${node.name}_${node.lang}`"
                      :value="node.name"
                    >
                      {{ node.name }}_{{ node.lang }}
                    </SelectItem>
                  </SelectContent>
                </Select>
                <Tooltip
                  v-if="editableAction.capture.onSuccessGoto
                    && getNodeUrl(editableAction.capture.onSuccessGoto)"
                >
                  <TooltipTrigger as-child>
                    <Button
                      variant="outline"
                      size="icon"
                      as-child
                    >
                      <NuxtLink :to="getNodeUrl(editableAction.capture.onSuccessGoto)">
                        <ArrowRight class="h-4 w-4" />
                      </NuxtLink>
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    {{ t('features.chatbotNode.capture.openNode') }}
                  </TooltipContent>
                </Tooltip>
              </div>
              <p class="text-xs text-muted-foreground">
                {{ t('features.chatbotNode.capture.onSuccessGotoHelp') }}
              </p>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
