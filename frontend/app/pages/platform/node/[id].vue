<script setup lang="ts">
import type { ChatbotNode, ChatbotNodeMessage, ChatbotNodeTrigger } from '~~/gen/chatbot/nodes/v1/node_pb';
import { Check, ChevronDown, Loader2 } from 'lucide-vue-next';
import { toast } from 'vue-sonner';
import MessageEditor from '@/components/features/chatbot-node/MessageEditor.vue';
import NodeDeleteDialog from '@/components/features/chatbot-node/NodeDeleteDialog.vue';
import TriggerEditor from '@/components/features/chatbot-node/TriggerEditor.vue';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { Switch } from '@/components/ui/switch';
import { useChatbotNodeService } from '@/composables/services/useChatbotNodeService';
import { useChatbotNodeStore } from '@/stores/chatbot-node';
import { useProjectStore } from '@/stores/project';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const projectStore = useProjectStore();
const nodeStore = useChatbotNodeStore();
const nodeService = useChatbotNodeService();

// Get node ID from route params
const nodeId = computed(() => route.params.id as string);
const projectId = computed(() => projectStore.activeProjectId);

// Loading states
const isLoading = ref(true);
const isSaving = computed(() => nodeService.updateLoading.value);

// Node data
const node = ref<ChatbotNode | null>(null);

// Form state (local edits)
const formData = ref({
  name: '',
  enabled: true,
  tags: [] as string[],
  triggers: [] as ChatbotNodeTrigger[],
  messages: [] as ChatbotNodeMessage[],
});

// Track if form has unsaved changes
const hasChanges = computed(() => {
  if (!node.value) {
    return false;
  }
  return (
    formData.value.name !== node.value.name
    || formData.value.enabled !== node.value.enabled
    || JSON.stringify(formData.value.triggers) !== JSON.stringify(node.value.triggers)
    || JSON.stringify(formData.value.messages) !== JSON.stringify(node.value.messages)
  );
});

// Display name for the node
const nodeName = computed(() => {
  if (!node.value) {
    return '';
  }
  return `${node.value.name}_${node.value.lang}`;
});

// Load node data
async function loadNode() {
  if (!projectId.value || !nodeId.value) {
    return;
  }

  isLoading.value = true;

  try {
    const result = await nodeService.getNode(projectId.value, nodeId.value);
    if (result) {
      node.value = result;
      // Initialize form data from node
      formData.value = {
        name: result.name,
        enabled: result.enabled,
        tags: [...result.tags],
        triggers: [...result.triggers],
        messages: [...result.messages],
      };
    }
    else {
      node.value = null;
    }
  }
  catch {
    node.value = null;
  }
  finally {
    isLoading.value = false;
  }
}

// Save changes
async function handleSave() {
  if (!projectId.value || !nodeId.value || !node.value) {
    return;
  }

  try {
    const result = await nodeService.updateNode(projectId.value, nodeId.value, {
      name: formData.value.name,
      enabled: formData.value.enabled,
      tags: formData.value.tags,
      triggers: formData.value.triggers,
      messages: formData.value.messages,
    });

    if (result) {
      node.value = result;
      nodeStore.updateNode(result);
      toast.success(t('common.success'), {
        description: t('features.chatbotNode.messages.saveSuccess'),
      });
    }
  }
  catch {
    toast.error(t('common.error'), {
      description: nodeService.updateError.value || t('features.chatbotNode.messages.saveError'),
    });
  }
}

// Handle node deletion
function handleDeleted() {
  // Navigate to another node or the nodes list
  const remainingNodes = nodeStore.sortedNodes.filter(n => n.id !== nodeId.value);
  if (remainingNodes.length > 0 && remainingNodes[0]) {
    router.push(`/platform/node/${remainingNodes[0].id}`);
  }
  else {
    router.push('/platform');
  }
}

// Watch for route changes and reload node
watch(
  () => route.params.id,
  () => {
    loadNode();
  },
);

// Initial load
onMounted(() => {
  loadNode();
});
</script>

<template>
  <div class="container mx-auto px-2 py-3">
    <!-- No project selected -->
    <div v-if="!projectId" class="text-center py-8">
      <p class="text-muted-foreground">
        {{ t('features.chatbotNode.page.noProjectSelected') }}
      </p>
    </div>

    <!-- Loading state -->
    <div v-else-if="isLoading" class="flex items-center justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <!-- Node not found -->
    <Alert v-else-if="!node" variant="destructive">
      <AlertTitle>{{ t('features.chatbotNode.page.notFoundTitle') }}</AlertTitle>
      <AlertDescription>
        {{ t('features.chatbotNode.page.notFoundDesc') }}
      </AlertDescription>
    </Alert>

    <!-- Node editor -->
    <div v-else class="max-w-3xl w-full pl-4 sm:pl-6 space-y-6">
      <!-- Header with node name and actions -->
      <div class="flex items-start justify-between gap-4">
        <div class="space-y-1">
          <div class="flex items-center gap-2">
            <h1 class="text-2xl font-bold">
              {{ nodeName }}
            </h1>
            <Badge :variant="formData.enabled ? 'default' : 'secondary'">
              {{ formData.enabled ? t('common.enabled') : t('common.disabled') }}
            </Badge>
          </div>
          <p class="text-sm text-muted-foreground">
            {{ t('features.chatbotNode.page.subtitle', { lang: node.lang }) }}
          </p>
        </div>

        <div class="flex items-center gap-2">
          <!-- Save button -->
          <Button
            :disabled="!hasChanges || isSaving"
            @click="handleSave"
          >
            <Check v-if="!isSaving" class="h-4 w-4 mr-2" />
            <Loader2 v-else class="h-4 w-4 mr-2 animate-spin" />
            {{ t('common.save') }}
          </Button>

          <!-- Delete button -->
          <NodeDeleteDialog :node="node" @deleted="handleDeleted" />
        </div>
      </div>

      <!-- Content card (like KeyReply) -->
      <Card>
        <CardHeader class="pb-4">
          <div class="flex items-center justify-between">
            <div class="space-y-1">
              <Label for="node-name">{{ t('features.chatbotNode.form.name') }}</Label>
              <Input
                id="node-name"
                v-model="formData.name"
                :placeholder="t('features.chatbotNode.form.namePlaceholder')"
                :disabled="isSaving"
                class="max-w-xs"
              />
            </div>

            <div class="flex items-center gap-2">
              <Label for="node-enabled">{{ t('features.chatbotNode.form.enabled') }}</Label>
              <Switch
                id="node-enabled"
                :model-value="formData.enabled"
                :disabled="isSaving"
                @update:model-value="formData.enabled = $event"
              />
            </div>
          </div>
        </CardHeader>

        <Separator />

        <CardContent class="pt-6 space-y-6">
          <!-- Triggers section (collapsible, collapsed by default) -->
          <Collapsible class="space-y-4">
            <CollapsibleTrigger class="flex items-center justify-between w-full group">
              <div class="flex items-center gap-2">
                <h3 class="text-lg font-semibold">
                  {{ t('features.chatbotNode.trigger.title') }}
                </h3>
                <Badge variant="secondary" class="text-xs">
                  {{ formData.triggers.length }}
                  {{ formData.triggers.length === 1 ? 'trigger' : 'triggers' }}
                </Badge>
              </div>
              <ChevronDown
                class="h-4 w-4 text-muted-foreground transition-transform
                  group-data-[state=open]:rotate-180"
              />
            </CollapsibleTrigger>
            <CollapsibleContent>
              <TriggerEditor
                v-model:triggers="formData.triggers"
                :disabled="isSaving"
                :show-header="false"
              />
            </CollapsibleContent>
          </Collapsible>

          <Separator />

          <!-- Messages section (chat bubbles) -->
          <MessageEditor
            v-model:messages="formData.messages"
            :disabled="isSaving"
          />
        </CardContent>
      </Card>

      <!-- Unsaved changes indicator -->
      <div
        v-if="hasChanges"
        class="fixed bottom-4 left-1/2 -translate-x-1/2 bg-background border
          rounded-lg shadow-lg px-4 py-2 flex items-center gap-3"
      >
        <span class="text-sm text-muted-foreground">
          {{ t('features.chatbotNode.page.unsavedChanges') }}
        </span>
        <Button size="sm" :disabled="isSaving" @click="handleSave">
          <Check v-if="!isSaving" class="h-4 w-4 mr-1" />
          <Loader2 v-else class="h-4 w-4 mr-1 animate-spin" />
          {{ t('common.save') }}
        </Button>
      </div>
    </div>
  </div>
</template>
