<script setup lang="ts">
import { FolderPlus, Plus, Trash2 } from 'lucide-vue-next';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

interface EditableCondition {
  id: string;
  type: 'group' | 'atomic';
  logicType?: 'all' | 'any' | 'not';
  children?: EditableCondition[];
  factType?: string;
  customFact?: string;
  operator?: string;
  value?: string;
  isNullValue?: boolean;
}

interface OptionItem {
  value: string;
  label: string;
  color?: string;
}

interface Props {
  condition: EditableCondition;
  disabled?: boolean;
  isRoot?: boolean;
  depth?: number;
  factTypeOptions: OptionItem[];
  operatorOptions: OptionItem[];
  logicTypeOptions: OptionItem[];
  predefinedValueOptions?: Record<string, OptionItem[]>;
}

// Enable recursive component
defineOptions({
  name: 'ConditionNode',
});

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  isRoot: false,
  depth: 0,
});

const emit = defineEmits<{
  updateLogicType: [condition: EditableCondition, value: string];
  updateAtomicField: [condition: EditableCondition, field: string, value: string];
  addChild: [group: EditableCondition, type: 'atomic' | 'group'];
  removeChild: [group: EditableCondition, index: number];
  removeSelf: [];
}>();

const { t } = useI18n();

function needsCustomInput(factType: string | undefined): boolean {
  return factType === 'context' || factType === 'custom';
}

function hasPredefinedValues(
  factType: string | undefined,
  predefinedOptions?: Record<string, OptionItem[]>,
): boolean {
  if (!factType || !predefinedOptions)
    return false;
  return factType in predefinedOptions;
}

function getPredefinedValues(
  factType: string | undefined,
  predefinedOptions?: Record<string, OptionItem[]>,
): OptionItem[] {
  if (!factType || !predefinedOptions)
    return [];
  return predefinedOptions[factType] || [];
}

function getCustomInputPlaceholder(factType: string | undefined): string {
  if (factType === 'context') {
    return t('features.chatbotNode.condition.contextVarPlaceholder');
  }
  return t('features.chatbotNode.condition.customFactPlaceholder');
}

function getLogicBadgeClass(logicType: string | undefined): string {
  if (logicType === 'all') {
    return 'bg-blue-500 hover:bg-blue-500 text-white';
  }
  if (logicType === 'any') {
    return 'bg-green-500 hover:bg-green-500 text-white';
  }
  if (logicType === 'not') {
    return 'bg-red-500 hover:bg-red-500 text-white';
  }
  return '';
}

function getLogicLabel(logicType: string | undefined): string {
  if (logicType === 'all') {
    return 'AND';
  }
  if (logicType === 'any') {
    return 'OR';
  }
  if (logicType === 'not') {
    return 'NOT';
  }
  return '';
}

function getBorderColor(logicType: string | undefined): string {
  if (logicType === 'all') {
    return 'border-blue-300 dark:border-blue-700';
  }
  if (logicType === 'any') {
    return 'border-green-300 dark:border-green-700';
  }
  if (logicType === 'not') {
    return 'border-red-300 dark:border-red-700';
  }
  return 'border-muted';
}

// Forward events from child components
function handleUpdateLogicType(condition: EditableCondition, value: string) {
  emit('updateLogicType', condition, value);
}

function handleUpdateAtomicField(
  condition: EditableCondition,
  field: string,
  value: string,
) {
  emit('updateAtomicField', condition, field, value);
}

function handleAddChild(group: EditableCondition, type: 'atomic' | 'group') {
  emit('addChild', group, type);
}

function handleRemoveChild(group: EditableCondition, index: number) {
  emit('removeChild', group, index);
}
</script>

<template>
  <!-- GROUP NODE -->
  <div v-if="condition.type === 'group'" class="space-y-3">
    <!-- Group Header -->
    <div class="flex items-center gap-3 flex-wrap">
      <Badge :class="getLogicBadgeClass(condition.logicType)">
        {{ getLogicLabel(condition.logicType) }}
      </Badge>

      <Select
        :model-value="condition.logicType"
        :disabled="disabled"
        @update:model-value="(v) => emit('updateLogicType', condition, String(v))"
      >
        <SelectTrigger class="w-[160px] h-8">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem
            v-for="opt in logicTypeOptions"
            :key="opt.value"
            :value="opt.value"
          >
            {{ opt.label }}
          </SelectItem>
        </SelectContent>
      </Select>

      <span class="text-sm text-muted-foreground">
        {{ t(`features.chatbotNode.condition.${condition.logicType}Help`) }}
      </span>
    </div>

    <!-- Children -->
    <div
      class="pl-4 border-l-2 space-y-3"
      :class="getBorderColor(condition.logicType)"
    >
      <div
        v-for="(child, index) in condition.children"
        :key="child.id"
        class="relative"
      >
        <!-- Child wrapper with remove button -->
        <div
          class="rounded-lg p-3"
          :class="child.type === 'group' ? 'bg-muted/30' : 'bg-muted/50'"
        >
          <div class="flex items-start gap-2">
            <div class="flex-1 min-w-0">
              <!-- Recursive render for nested groups -->
              <ConditionNode
                v-if="child.type === 'group'"
                :condition="child"
                :disabled="disabled"
                :depth="depth + 1"
                :fact-type-options="factTypeOptions"
                :operator-options="operatorOptions"
                :logic-type-options="logicTypeOptions"
                :predefined-value-options="predefinedValueOptions"
                @update-logic-type="handleUpdateLogicType"
                @update-atomic-field="handleUpdateAtomicField"
                @add-child="handleAddChild"
                @remove-child="handleRemoveChild"
              />

              <!-- Atomic condition -->
              <div v-else class="space-y-2">
                <!-- Row 1: Fact Type + Custom Input -->
                <div class="flex items-center gap-2 flex-wrap">
                  <Select
                    :model-value="child.factType"
                    :disabled="disabled"
                    @update:model-value="
                      (v) => emit('updateAtomicField', child, 'factType', String(v))
                    "
                  >
                    <SelectTrigger class="w-[160px] h-8">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="opt in factTypeOptions"
                        :key="opt.value"
                        :value="opt.value"
                      >
                        {{ opt.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>

                  <Input
                    v-if="needsCustomInput(child.factType)"
                    :model-value="child.customFact"
                    :placeholder="getCustomInputPlaceholder(child.factType)"
                    :disabled="disabled"
                    class="w-[150px] h-8"
                    @update:model-value="
                      (v) => emit('updateAtomicField', child, 'customFact', String(v))
                    "
                  />
                </div>

                <!-- Row 2: Operator + Value -->
                <div class="flex items-center gap-2 flex-wrap">
                  <Select
                    :model-value="child.operator"
                    :disabled="disabled"
                    @update:model-value="
                      (v) => emit('updateAtomicField', child, 'operator', String(v))
                    "
                  >
                    <SelectTrigger class="w-[140px] h-8">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="opt in operatorOptions"
                        :key="opt.value"
                        :value="opt.value"
                      >
                        {{ opt.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>

                  <!-- Predefined value dropdown for certain fact types -->
                  <Select
                    v-if="hasPredefinedValues(child.factType, props.predefinedValueOptions)"
                    :model-value="child.value"
                    :disabled="disabled"
                    @update:model-value="
                      (v) => emit('updateAtomicField', child, 'value', String(v))
                    "
                  >
                    <SelectTrigger class="flex-1 min-w-[160px] h-8">
                      <SelectValue :placeholder="t('features.chatbotNode.condition.selectValue')" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="opt in getPredefinedValues(child.factType, props.predefinedValueOptions)"
                        :key="opt.value"
                        :value="opt.value"
                      >
                        {{ opt.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>

                  <!-- Free text input for other fact types -->
                  <Input
                    v-else
                    :model-value="child.value"
                    :placeholder="t('features.chatbotNode.condition.valuePlaceholder')"
                    :disabled="disabled"
                    class="flex-1 min-w-[120px] h-8"
                    @update:model-value="
                      (v) => emit('updateAtomicField', child, 'value', String(v))
                    "
                  />
                </div>
              </div>
            </div>

            <!-- Remove button - show for groups even if only child, hide only for last atomic -->
            <Button
              v-if="condition.children && (condition.children.length > 1 || child.type === 'group')"
              variant="ghost"
              size="icon"
              class="h-8 w-8 shrink-0"
              :disabled="disabled"
              @click="emit('removeChild', condition, index)"
            >
              <Trash2 class="h-4 w-4 text-destructive" />
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add buttons (not for NOT which can only have one child) -->
    <div
      v-if="condition.logicType !== 'not'"
      class="flex items-center gap-2 pl-4"
    >
      <Button
        variant="outline"
        size="sm"
        :disabled="disabled"
        @click="emit('addChild', condition, 'atomic')"
      >
        <Plus class="h-4 w-4 mr-1" />
        {{ t('features.chatbotNode.condition.addCondition') }}
      </Button>

      <Button
        variant="outline"
        size="sm"
        :disabled="disabled"
        @click="emit('addChild', condition, 'group')"
      >
        <FolderPlus class="h-4 w-4 mr-1" />
        {{ t('features.chatbotNode.condition.addGroup') }}
      </Button>
    </div>
  </div>

  <!-- ATOMIC NODE (shouldn't be rendered at root level normally) -->
  <div v-else class="space-y-2 p-3 bg-muted/50 rounded-lg">
    <div class="flex items-center gap-2 flex-wrap">
      <Select
        :model-value="condition.factType"
        :disabled="disabled"
        @update:model-value="
          (v) => emit('updateAtomicField', condition, 'factType', String(v))
        "
      >
        <SelectTrigger class="w-[160px] h-8">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem
            v-for="opt in factTypeOptions"
            :key="opt.value"
            :value="opt.value"
          >
            {{ opt.label }}
          </SelectItem>
        </SelectContent>
      </Select>

      <Input
        v-if="needsCustomInput(condition.factType)"
        :model-value="condition.customFact"
        :placeholder="getCustomInputPlaceholder(condition.factType)"
        :disabled="disabled"
        class="w-[150px] h-8"
        @update:model-value="
          (v) => emit('updateAtomicField', condition, 'customFact', String(v))
        "
      />
    </div>

    <div class="flex items-center gap-2 flex-wrap">
      <Select
        :model-value="condition.operator"
        :disabled="disabled"
        @update:model-value="
          (v) => emit('updateAtomicField', condition, 'operator', String(v))
        "
      >
        <SelectTrigger class="w-[140px] h-8">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem
            v-for="opt in operatorOptions"
            :key="opt.value"
            :value="opt.value"
          >
            {{ opt.label }}
          </SelectItem>
        </SelectContent>
      </Select>

      <!-- Predefined value dropdown for certain fact types -->
      <Select
        v-if="hasPredefinedValues(condition.factType, props.predefinedValueOptions)"
        :model-value="condition.value"
        :disabled="disabled"
        @update:model-value="
          (v) => emit('updateAtomicField', condition, 'value', String(v))
        "
      >
        <SelectTrigger class="flex-1 min-w-[160px] h-8">
          <SelectValue :placeholder="t('features.chatbotNode.condition.selectValue')" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem
            v-for="opt in getPredefinedValues(condition.factType, props.predefinedValueOptions)"
            :key="opt.value"
            :value="opt.value"
          >
            {{ opt.label }}
          </SelectItem>
        </SelectContent>
      </Select>

      <!-- Free text input for other fact types -->
      <Input
        v-else
        :model-value="condition.value"
        :placeholder="t('features.chatbotNode.condition.valuePlaceholder')"
        :disabled="disabled"
        class="flex-1 min-w-[120px] h-8"
        @update:model-value="
          (v) => emit('updateAtomicField', condition, 'value', String(v))
        "
      />
    </div>
  </div>
</template>
