<script setup lang="ts">
import type { NodeCondition } from '~~/gen/chatbot/nodes/v1/node_pb';
import { Plus, X } from 'lucide-vue-next';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import ConditionNode from './ConditionNode.vue';

/**
 * Recursive condition structure supporting nested groups
 * - type: 'group' for composite conditions (all/any/not)
 * - type: 'atomic' for leaf conditions (fact comparison)
 */
interface EditableCondition {
  id: string; // Unique ID for Vue key
  type: 'group' | 'atomic';
  // For groups (composite conditions)
  logicType?: 'all' | 'any' | 'not';
  children?: EditableCondition[];
  // For atomic conditions
  factType?: string;
  customFact?: string;
  operator?: string;
  value?: string;
}

interface Props {
  condition?: NodeCondition;
  disabled?: boolean;
}

// Recursive component for rendering condition nodes
defineOptions({
  name: 'ConditionBuilder',
});

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
});

const emit = defineEmits<{
  'update:condition': [value: NodeCondition | undefined];
}>();

const { t } = useI18n();

// Generate unique ID
let idCounter = 0;
function generateId(): string {
  return `cond_${Date.now()}_${idCounter++}`;
}

// Available fact types
const factTypeOptions = [
  { value: 'session.mode', label: 'Session Mode' },
  { value: 'context', label: 'Context Variable' },
  { value: 'message.text', label: 'Message Text' },
  { value: 'message.lang', label: 'Message Language' },
  { value: 'custom', label: 'Custom Fact' },
];

// Available operators
const operatorOptions = [
  { value: 'equal', label: 'Equals' },
  { value: 'notEqual', label: 'Not Equals' },
  { value: 'contains', label: 'Contains' },
  { value: 'doesNotContain', label: 'Does Not Contain' },
  { value: 'greaterThan', label: 'Greater Than' },
  { value: 'lessThan', label: 'Less Than' },
  { value: 'in', label: 'In List' },
];

// Logic type options for groups
const logicTypeOptions = [
  { value: 'all', label: 'ALL of (AND)', color: 'bg-blue-500' },
  { value: 'any', label: 'ANY of (OR)', color: 'bg-green-500' },
  { value: 'not', label: 'NONE of (NOT)', color: 'bg-red-500' },
];

// ============ Conversion Functions ============

function determineFaType(fact: string, path: string): { factType: string; customFact: string } {
  // Handle empty fact (custom type with empty customFact)
  if (!fact || fact === '') {
    if (path.startsWith('$.contextLast.')) {
      return { factType: 'context', customFact: path.replace('$.contextLast.', '') };
    }
    return { factType: 'custom', customFact: '' };
  }
  if (fact === 'session' && path.startsWith('$.contextLast.')) {
    return {
      factType: 'context',
      customFact: path.replace('$.contextLast.', ''),
    };
  }
  if (fact === 'session.mode' || fact === 'message.text' || fact === 'message.lang') {
    return { factType: fact, customFact: '' };
  }
  return { factType: 'custom', customFact: fact };
}

function factTypeToFactPath(factType: string, customFact: string): { fact: string; path: string } {
  if (factType === 'context') {
    return { fact: 'session', path: `$.contextLast.${customFact}` };
  }
  if (factType === 'custom') {
    return { fact: customFact, path: '' };
  }
  return { fact: factType, path: '' };
}

// Convert protobuf to editable (recursive)
function protoToEditable(condition: NodeCondition | undefined): EditableCondition | null {
  if (!condition) {
    return null;
  }

  // Check for composite conditions
  if (condition.all && condition.all.length > 0) {
    return {
      id: generateId(),
      type: 'group',
      logicType: 'all',
      children: condition.all.map(c => protoToEditable(c)).filter(Boolean) as EditableCondition[],
    };
  }

  if (condition.any && condition.any.length > 0) {
    return {
      id: generateId(),
      type: 'group',
      logicType: 'any',
      children: condition.any.map(c => protoToEditable(c)).filter(Boolean) as EditableCondition[],
    };
  }

  if (condition.not) {
    const notChild = protoToEditable(condition.not);
    return {
      id: generateId(),
      type: 'group',
      logicType: 'not',
      children: notChild ? [notChild] : [],
    };
  }

  // Atomic condition - check operator (fact can be empty for custom types)
  if (condition.operator) {
    let valueStr = '';
    if (condition.value) {
      // Handle protobuf-es Value structure: { kind: { case: 'stringValue', value: '...' } }
      const v = condition.value as { kind?: { case: string; value: unknown } };
      if (v.kind && v.kind.value !== undefined) {
        valueStr = String(v.kind.value);
      }
    }

    const { factType, customFact } = determineFaType(condition.fact, condition.path || '');

    return {
      id: generateId(),
      type: 'atomic',
      factType,
      customFact,
      operator: condition.operator,
      value: valueStr,
    };
  }

  return null;
}

// Convert editable to protobuf (recursive)
function editableToProto(editable: EditableCondition | null): NodeCondition | undefined {
  if (!editable) {
    return undefined;
  }

  if (editable.type === 'group') {
    const children = (editable.children || [])
      .map(c => editableToProto(c))
      .filter(Boolean) as NodeCondition[];

    if (editable.logicType === 'all') {
      return {
        all: children,
        any: [],
        fact: '',
        operator: '',
        path: '',
      } as unknown as NodeCondition;
    }
    if (editable.logicType === 'any') {
      return {
        all: [],
        any: children,
        fact: '',
        operator: '',
        path: '',
      } as unknown as NodeCondition;
    }
    if (editable.logicType === 'not' && children.length > 0) {
      return {
        all: [],
        any: [],
        not: children[0],
        fact: '',
        operator: '',
        path: '',
      } as unknown as NodeCondition;
    }
    return undefined;
  }

  // Atomic condition
  if (editable.factType && editable.operator) {
    const { fact, path } = factTypeToFactPath(editable.factType, editable.customFact || '');

    // Only create valueObj if value is non-empty
    // Empty string causes "google.protobuf.Value must have a value" error
    // Use protobuf-es Value structure: { kind: { case: '...', value: ... } }
    let valueObj: unknown;
    const trimmedValue = (editable.value || '').trim();
    if (trimmedValue) {
      if (trimmedValue === 'true' || trimmedValue === 'false') {
        valueObj = { kind: { case: 'boolValue', value: trimmedValue === 'true' } };
      }
      else if (!Number.isNaN(Number(trimmedValue))) {
        valueObj = { kind: { case: 'numberValue', value: Number(trimmedValue) } };
      }
      else {
        valueObj = { kind: { case: 'stringValue', value: trimmedValue } };
      }
    }

    return {
      all: [],
      any: [],
      fact,
      operator: editable.operator,
      value: valueObj,
      path: path || '',
    } as unknown as NodeCondition;
  }

  return undefined;
}

// ============ State Management ============

const rootCondition = ref<EditableCondition | null>(null);
const hasCondition = computed(() => rootCondition.value !== null);

// Flag to skip watch after our own emits
let skipNextWatch = false;

// Initialize from prop - only sync from external changes
watch(
  () => props.condition,
  (newVal) => {
    if (skipNextWatch) {
      skipNextWatch = false;
      return;
    }
    const converted = protoToEditable(newVal);
    rootCondition.value = converted;
  },
  { immediate: true },
);

// Emit changes
function emitUpdate() {
  skipNextWatch = true;
  const proto = editableToProto(rootCondition.value);
  emit('update:condition', proto);
}

// ============ Actions ============

function addRootCondition() {
  rootCondition.value = {
    id: generateId(),
    type: 'group',
    logicType: 'all',
    children: [createAtomicCondition()],
  };
  emitUpdate();
}

function removeRootCondition() {
  rootCondition.value = null;
  emit('update:condition', undefined);
}

function createAtomicCondition(): EditableCondition {
  return {
    id: generateId(),
    type: 'atomic',
    factType: 'session.mode',
    customFact: '',
    operator: 'equal',
    value: '',
  };
}

function createGroupCondition(logicType: 'all' | 'any' | 'not' = 'all'): EditableCondition {
  return {
    id: generateId(),
    type: 'group',
    logicType,
    children: [createAtomicCondition()],
  };
}

function addChildToGroup(group: EditableCondition, childType: 'atomic' | 'group') {
  if (!group.children) {
    group.children = [];
  }
  if (childType === 'atomic') {
    group.children.push(createAtomicCondition());
  }
  else {
    group.children.push(createGroupCondition());
  }
  emitUpdate();
}

function removeChildFromGroup(group: EditableCondition, index: number) {
  if (group.children) {
    group.children.splice(index, 1);
    emitUpdate();
  }
}

function updateGroupLogicType(group: EditableCondition, logicType: string) {
  group.logicType = logicType as 'all' | 'any' | 'not';
  // NOT can only have one child
  if (logicType === 'not' && group.children && group.children.length > 1) {
    group.children = [group.children[0]!];
  }
  emitUpdate();
}

function updateAtomicField(condition: EditableCondition, field: string, value: string) {
  (condition as unknown as Record<string, unknown>)[field] = value;
  emitUpdate();
}
</script>

<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <Label class="text-base font-medium">
        <!-- {{ t('features.chatbotNode.condition.title') }} -->
      </Label>
      <Button
        v-if="!hasCondition"
        variant="outline"
        size="sm"
        :disabled="disabled"
        @click="addRootCondition"
      >
        <Plus class="h-4 w-4 mr-2" />
        {{ t('features.chatbotNode.condition.add') }}
      </Button>
      <Button
        v-else
        variant="ghost"
        size="sm"
        :disabled="disabled"
        @click="removeRootCondition"
      >
        <X class="h-4 w-4 mr-2" />
        {{ t('features.chatbotNode.condition.remove') }}
      </Button>
    </div>

    <p class="text-sm text-muted-foreground">
      {{ t('features.chatbotNode.condition.description') }}
    </p>

    <!-- No condition -->
    <div
      v-if="!hasCondition"
      class="text-sm text-muted-foreground italic py-4 text-center"
    >
      {{ t('features.chatbotNode.condition.noCondition') }}
    </div>

    <!-- Condition Editor -->
    <Card v-else-if="rootCondition">
      <CardContent class="pt-4">
        <!-- Recursive condition renderer -->
        <ConditionNode
          :condition="rootCondition"
          :disabled="disabled"
          :is-root="true"
          :fact-type-options="factTypeOptions"
          :operator-options="operatorOptions"
          :logic-type-options="logicTypeOptions"
          @update-logic-type="
            (c: EditableCondition, v: string) => updateGroupLogicType(c, v)
          "
          @update-atomic-field="
            (c: EditableCondition, f: string, v: string) => updateAtomicField(c, f, v)
          "
          @add-child="
            (g: EditableCondition, t: 'atomic' | 'group') => addChildToGroup(g, t)
          "
          @remove-child="(g: EditableCondition, i: number) => removeChildFromGroup(g, i)"
        />
      </CardContent>
    </Card>
  </div>
</template>
