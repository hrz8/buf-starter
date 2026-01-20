# Task T52: Node Edit Page & Components

**Story Reference:** US12-node-editor.md
**Type:** Frontend UI
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T51 (Node Frontend Foundation)

## Objective

Implement the node creation sheet, edit page, delete dialog, trigger editor, and message editor components.

## Acceptance Criteria

- [ ] NodeCreateSheet.vue - Sheet/dialog for creating new nodes (name, lang, tags)
- [ ] Node edit page (`/platform/node/[id].vue`) loads and displays node data
- [ ] NodeEditForm.vue - Main form with name, tags, enabled fields
- [ ] TriggerEditor.vue - Array editor with add/remove, type select, value input
- [ ] MessageEditor.vue - Array editor with add/remove, reorder, textarea for content
- [ ] NodeDeleteDialog.vue - Confirmation dialog for deletion
- [ ] Form validation shows errors per field
- [ ] Save button persists changes via API
- [ ] After delete, navigates to another node or platform home

## Technical Requirements

### NodeCreateSheet.vue

```vue
<script setup lang="ts">
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { nodeCreateSchema } from './schema';

const emit = defineEmits<{
  created: [node: ChatbotNode];
}>();

const isOpen = ref(false);

const { handleSubmit, resetForm, isSubmitting } = useForm({
  validationSchema: toTypedSchema(nodeCreateSchema),
  initialValues: {
    name: '',
    lang: 'en',
    tags: [],
  },
});

const { createNode, createError } = useNodeService();
const router = useRouter();

const onSubmit = handleSubmit(async (values) => {
  const projectId = useProjectStore().activeProjectId;
  const node = await createNode({
    projectId,
    ...values,
  });
  emit('created', node);
  isOpen.value = false;
  resetForm();
  router.push(`/platform/node/${node.id}`);
});

function open() {
  isOpen.value = true;
}

defineExpose({ open });
</script>
```

**Features:**
- Sheet component with form: name (input), lang (select), tags (tag input)
- Name validation: lowercase_snake_case pattern
- On submit: calls createNode, closes sheet, navigates to new node edit page
- Emits `created` event for sidebar refresh

### Node Edit Page ([id].vue)

```vue
<script setup lang="ts">
definePageMeta({
  layout: 'platform',
});

const route = useRoute();
const nodeId = computed(() => route.params.id as string);

const { getNode, getLoading, getError } = useNodeService();
const projectStore = useProjectStore();

const node = ref<ChatbotNode | null>(null);
const isLoading = ref(true);

onMounted(async () => {
  try {
    node.value = await getNode(projectStore.activeProjectId, nodeId.value);
  } finally {
    isLoading.value = false;
  }
});

const breadcrumb = computed(() => [
  { label: 'Platform', to: '/platform' },
  { label: 'Nodes', to: '/platform' },
  { label: node.value ? `${node.value.name}_${node.value.lang}` : 'Loading...' },
]);

async function handleSave(updated: ChatbotNode) {
  node.value = updated;
  // Show success toast
}

async function handleDeleted() {
  const router = useRouter();
  router.push('/platform');
}
</script>

<template>
  <PageLayout :breadcrumb="breadcrumb">
    <template #actions>
      <NodeDeleteDialog
        v-if="node"
        :node="node"
        @deleted="handleDeleted"
      />
    </template>

    <div v-if="isLoading" class="flex justify-center py-8">
      <Spinner />
    </div>

    <NodeEditForm
      v-else-if="node"
      :node="node"
      @save="handleSave"
    />

    <div v-else class="text-center py-8 text-muted-foreground">
      Node not found
    </div>
  </PageLayout>
</template>
```

### NodeEditForm.vue

```vue
<script setup lang="ts">
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { nodeEditSchema } from './schema';

const props = defineProps<{
  node: ChatbotNode;
}>();

const emit = defineEmits<{
  save: [node: ChatbotNode];
}>();

const { updateNode, updateLoading, updateError } = useNodeService();
const projectStore = useProjectStore();
const { t } = useI18n();

const { handleSubmit, values, setFieldValue } = useForm({
  validationSchema: toTypedSchema(nodeEditSchema),
  initialValues: {
    name: props.node.name,
    tags: props.node.tags || [],
    enabled: props.node.enabled,
    triggers: props.node.triggers,
    messages: props.node.messages,
  },
});

const onSubmit = handleSubmit(async (formValues) => {
  const updated = await updateNode({
    projectId: projectStore.activeProjectId,
    nodeId: props.node.id,
    ...formValues,
  });
  emit('save', updated);
});
</script>

<template>
  <form @submit="onSubmit" class="space-y-6">
    <!-- Read-only language display -->
    <div class="flex items-center gap-2 text-sm text-muted-foreground">
      <span>{{ t('features.node.form.lang') }}:</span>
      <Badge variant="secondary">{{ node.lang }}</Badge>
    </div>

    <!-- Name field -->
    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>{{ t('features.node.form.name') }}</FormLabel>
        <FormControl>
          <Input v-bind="componentField" />
        </FormControl>
        <FormDescription>{{ t('features.node.form.nameHelp') }}</FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- Tags field -->
    <FormField v-slot="{ componentField }" name="tags">
      <FormItem>
        <FormLabel>{{ t('features.node.form.tags') }}</FormLabel>
        <FormControl>
          <TagInput v-bind="componentField" />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- Enabled toggle -->
    <FormField v-slot="{ value, handleChange }" name="enabled">
      <FormItem class="flex items-center gap-2">
        <FormControl>
          <Switch :checked="value" @update:checked="handleChange" />
        </FormControl>
        <FormLabel>{{ t('features.node.form.enabled') }}</FormLabel>
      </FormItem>
    </FormField>

    <!-- Triggers section -->
    <div class="space-y-2">
      <Label>{{ t('features.node.form.triggers') }}</Label>
      <TriggerEditor
        :triggers="values.triggers"
        @update="(triggers) => setFieldValue('triggers', triggers)"
      />
      <FormField name="triggers">
        <FormMessage />
      </FormField>
    </div>

    <!-- Messages section -->
    <div class="space-y-2">
      <Label>{{ t('features.node.form.messages') }}</Label>
      <MessageEditor
        :messages="values.messages"
        @update="(messages) => setFieldValue('messages', messages)"
      />
      <FormField name="messages">
        <FormMessage />
      </FormField>
    </div>

    <!-- Submit button -->
    <Button type="submit" :disabled="updateLoading">
      <Spinner v-if="updateLoading" class="mr-2" />
      {{ t('features.node.form.save') }}
    </Button>
  </form>
</template>
```

### TriggerEditor.vue

```vue
<script setup lang="ts">
import { TRIGGER_TYPE_OPTIONS, DEFAULT_TRIGGER } from './constants';
import type { TriggerInput } from './schema';

const props = defineProps<{
  triggers: TriggerInput[];
}>();

const emit = defineEmits<{
  update: [triggers: TriggerInput[]];
}>();

const { t } = useI18n();

function addTrigger() {
  emit('update', [...props.triggers, { ...DEFAULT_TRIGGER }]);
}

function removeTrigger(index: number) {
  const updated = props.triggers.filter((_, i) => i !== index);
  emit('update', updated);
}

function updateTrigger(index: number, field: keyof TriggerInput, value: string) {
  const updated = [...props.triggers];
  updated[index] = { ...updated[index], [field]: value };
  emit('update', updated);
}

// Regex validation preview
function isValidRegex(pattern: string): boolean {
  try {
    new RegExp(pattern);
    return true;
  } catch {
    return false;
  }
}
</script>

<template>
  <div class="space-y-2">
    <div
      v-for="(trigger, index) in triggers"
      :key="index"
      class="flex items-center gap-2"
    >
      <Select
        :model-value="trigger.type"
        @update:model-value="(v) => updateTrigger(index, 'type', v)"
      >
        <SelectTrigger class="w-32">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem
            v-for="option in TRIGGER_TYPE_OPTIONS"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </SelectItem>
        </SelectContent>
      </Select>

      <div class="flex-1 relative">
        <Input
          :model-value="trigger.value"
          @update:model-value="(v) => updateTrigger(index, 'value', v)"
          :placeholder="t('features.node.trigger.value')"
        />
        <!-- Regex validation indicator -->
        <div
          v-if="trigger.type === 'regex' && trigger.value && !isValidRegex(trigger.value)"
          class="absolute right-2 top-1/2 -translate-y-1/2 text-destructive"
        >
          <AlertCircle class="h-4 w-4" />
        </div>
      </div>

      <Button
        type="button"
        variant="ghost"
        size="icon"
        @click="removeTrigger(index)"
        :disabled="triggers.length <= 1"
      >
        <Trash2 class="h-4 w-4" />
      </Button>
    </div>

    <Button type="button" variant="outline" size="sm" @click="addTrigger">
      <Plus class="h-4 w-4 mr-1" />
      {{ t('features.node.trigger.add') }}
    </Button>
  </div>
</template>
```

### MessageEditor.vue

```vue
<script setup lang="ts">
import type { MessageInput } from './schema';
import { DEFAULT_MESSAGE } from './constants';

const props = defineProps<{
  messages: MessageInput[];
}>();

const emit = defineEmits<{
  update: [messages: MessageInput[]];
}>();

const { t } = useI18n();

function addMessage() {
  emit('update', [...props.messages, { ...DEFAULT_MESSAGE }]);
}

function removeMessage(index: number) {
  const updated = props.messages.filter((_, i) => i !== index);
  emit('update', updated);
}

function updateContent(index: number, content: string) {
  const updated = [...props.messages];
  updated[index] = { ...updated[index], content };
  emit('update', updated);
}

function moveUp(index: number) {
  if (index === 0) return;
  const updated = [...props.messages];
  [updated[index - 1], updated[index]] = [updated[index], updated[index - 1]];
  emit('update', updated);
}

function moveDown(index: number) {
  if (index === props.messages.length - 1) return;
  const updated = [...props.messages];
  [updated[index], updated[index + 1]] = [updated[index + 1], updated[index]];
  emit('update', updated);
}
</script>

<template>
  <div class="space-y-3">
    <div
      v-for="(message, index) in messages"
      :key="index"
      class="border rounded-lg p-3 space-y-2"
    >
      <div class="flex items-center justify-between">
        <span class="text-sm text-muted-foreground">
          {{ t('features.node.message.content') }} #{{ index + 1 }}
        </span>
        <div class="flex items-center gap-1">
          <Button
            type="button"
            variant="ghost"
            size="icon"
            @click="moveUp(index)"
            :disabled="index === 0"
          >
            <ChevronUp class="h-4 w-4" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            @click="moveDown(index)"
            :disabled="index === messages.length - 1"
          >
            <ChevronDown class="h-4 w-4" />
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            @click="removeMessage(index)"
            :disabled="messages.length <= 1"
          >
            <Trash2 class="h-4 w-4" />
          </Button>
        </div>
      </div>

      <Textarea
        :model-value="message.content"
        @update:model-value="(v) => updateContent(index, v)"
        :placeholder="t('features.node.message.content')"
        rows="3"
      />

      <div class="text-xs text-muted-foreground text-right">
        {{ message.content.length }} / 5000
      </div>
    </div>

    <Button type="button" variant="outline" size="sm" @click="addMessage">
      <Plus class="h-4 w-4 mr-1" />
      {{ t('features.node.message.add') }}
    </Button>
  </div>
</template>
```

### NodeDeleteDialog.vue

```vue
<script setup lang="ts">
const props = defineProps<{
  node: ChatbotNode;
}>();

const emit = defineEmits<{
  deleted: [];
}>();

const { deleteNode, deleteLoading } = useNodeService();
const projectStore = useProjectStore();
const nodeStore = useNodeStore();
const { t } = useI18n();

const isOpen = ref(false);

async function handleDelete() {
  await deleteNode(projectStore.activeProjectId, props.node.id);
  nodeStore.removeNode(props.node.id);
  isOpen.value = false;
  emit('deleted');
}
</script>

<template>
  <AlertDialog v-model:open="isOpen">
    <AlertDialogTrigger as-child>
      <Button variant="destructive" size="sm">
        <Trash2 class="h-4 w-4 mr-1" />
        {{ t('features.node.delete.title') }}
      </Button>
    </AlertDialogTrigger>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ t('features.node.delete.title') }}</AlertDialogTitle>
        <AlertDialogDescription>
          {{ t('features.node.delete.description', { name: `${node.name}_${node.lang}` }) }}
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>{{ t('features.node.delete.cancel') }}</AlertDialogCancel>
        <AlertDialogAction
          @click="handleDelete"
          :disabled="deleteLoading"
          class="bg-destructive text-destructive-foreground"
        >
          <Spinner v-if="deleteLoading" class="mr-2" />
          {{ t('features.node.delete.confirm') }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
```

## Files to Create

```
frontend/app/components/features/chatbot-node/
├── NodeCreateSheet.vue    # Create node sheet/dialog
├── NodeEditForm.vue       # Main edit form
├── NodeDeleteDialog.vue   # Delete confirmation
├── TriggerEditor.vue      # Trigger array editor
└── MessageEditor.vue      # Message array editor

frontend/app/pages/platform/node/
└── [id].vue               # Node edit page
```

## Files to Modify

- `frontend/app/components/features/chatbot-node/index.ts` - Export new components

## Validation Checklist

- [ ] Create sheet validates and creates node
- [ ] Edit page loads node data correctly
- [ ] TriggerEditor adds/removes triggers
- [ ] TriggerEditor validates regex patterns
- [ ] MessageEditor adds/removes/reorders messages
- [ ] Form save persists all changes
- [ ] Delete dialog confirms and removes node
- [ ] All forms follow vee-validate best practices
- [ ] Forms are responsive on mobile

## Definition of Done

- [ ] All 5 components implemented
- [ ] Node edit page functional
- [ ] Trigger editor supports all trigger types
- [ ] Message editor with reordering
- [ ] Form validation working
- [ ] Delete flow complete
- [ ] Toast notifications for success/error
- [ ] Responsive design verified

## Dependencies

- T51: Frontend foundation (repository, service, store, schemas)
- shadcn-vue components (Sheet, Dialog, Form, etc.)

## Risk Factors

- **Medium Risk**: vee-validate integration with array fields
- **Low Risk**: Standard component implementation

## Notes

- Follow vee-validate FormField best practices from CLAUDE.md
- Use isLoading = ref(true) pattern for form loading
- NO :key attributes on FormField components
- Regex validation shows inline error indicator
- Message reordering uses simple swap logic
