<script setup lang="ts">
import type {
  ChatbotNode,
  ChatbotNodeMessage,
  ChatbotNodeTrigger,
  NodeCondition,
  NodeEffect,
  NodeNextAction,
} from '~~/gen/chatbot/nodes/v1/node_pb';
import { Check, ChevronDown, Loader2, Plus, ShieldCheck } from 'lucide-vue-next';
import { toast } from 'vue-sonner';
import ConditionBuilder from '@/components/features/chatbot-node/ConditionBuilder.vue';
import EffectEditor from '@/components/features/chatbot-node/EffectEditor.vue';
import MessageEditor from '@/components/features/chatbot-node/MessageEditor.vue';
import NextActionEditor from '@/components/features/chatbot-node/NextActionEditor.vue';
import NodeAddVersionDialog from '@/components/features/chatbot-node/NodeAddVersionDialog.vue';
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { Switch } from '@/components/ui/switch';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useChatbotNodeService } from '@/composables/services/useChatbotNodeService';
import { usePageTitle } from '@/composables/usePageTitle';
import { useChatbotNodeStore } from '@/stores/chatbot-node';
import { useProjectStore } from '@/stores/project';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const projectStore = useProjectStore();
const nodeStore = useChatbotNodeStore();
const nodeService = useChatbotNodeService();

usePageTitle(computed(() => t('features.chatbotNode.page.title')));

// Get name from route params, lang and version from query
const nodeName = computed(() => route.params.name as string);
const requestedLang = computed(() => (route.query.lang as string) || '');
const requestedVersion = computed(() => route.query.v as string | undefined);
const projectId = computed(() => projectStore.activeProjectId);

// Loading states
const isLoading = ref(true);
const isSaving = computed(() => nodeService.updateLoading.value);

// Node data
const node = ref<ChatbotNode | null>(null);

// Get all languages available for the current version
const availableLanguages = computed(() => {
  if (!nodeName.value)
    return [];
  return nodeStore.getLanguagesForVersion(nodeName.value, requestedVersion.value);
});

// Get current language (from query or first available)
const currentLang = computed(() => {
  if (requestedLang.value && availableLanguages.value.includes(requestedLang.value)) {
    return requestedLang.value;
  }
  return availableLanguages.value[0] || 'en-US';
});

// Get all unique versions across all languages (for version selector)
const uniqueVersions = computed(() => {
  if (!nodeName.value)
    return [];
  return nodeStore.getUniqueVersions(nodeName.value);
});

// Check if there are multiple versions
const hasVersions = computed(() => {
  return uniqueVersions.value.length > 1 || uniqueVersions.value.some(v => v !== '');
});

// Form state (local edits)
const formData = ref({
  name: '',
  enabled: true,
  tags: [] as string[],
  triggers: [] as ChatbotNodeTrigger[],
  messages: [] as ChatbotNodeMessage[],
  priority: 0,
  condition: undefined as NodeCondition | undefined,
  effect: undefined as NodeEffect | undefined,
  nextAction: undefined as NodeNextAction | undefined,
  forceCondition: false,
});

// Check if node is predefined (system node) - based on actual node data
const isPredefined = computed(() => node.value?.isPredefined ?? false);

// Track if form has unsaved changes
const hasChanges = computed(() => {
  if (!node.value) {
    return false;
  }
  return (
    formData.value.name !== node.value.name
    || formData.value.enabled !== node.value.enabled
    || formData.value.priority !== (node.value.priority ?? 0)
    || formData.value.forceCondition !== (node.value.forceCondition ?? false)
    || JSON.stringify(formData.value.triggers) !== JSON.stringify(node.value.triggers)
    || JSON.stringify(formData.value.messages) !== JSON.stringify(node.value.messages)
    || JSON.stringify(formData.value.condition) !== JSON.stringify(node.value.condition)
    || JSON.stringify(formData.value.effect) !== JSON.stringify(node.value.effect)
    || JSON.stringify(formData.value.nextAction) !== JSON.stringify(node.value.nextAction)
  );
});

// Display name for the header (just node name, no lang)
const displayName = computed(() => {
  return node.value?.name || nodeName.value || '';
});

// Find the current node based on name, lang, and version query params
function findCurrentNode(): ChatbotNode | undefined {
  if (!nodeName.value)
    return undefined;

  const nodeVersions = nodeStore.getNodeVersions(nodeName.value, currentLang.value);

  if (nodeVersions.length === 0) {
    // Try to find any node with this name (for first load when lang might not be in URL)
    return nodeStore.findNodeWithVersion(nodeName.value, requestedVersion.value);
  }

  // If version is requested, try to find it
  if (requestedVersion.value) {
    const found = nodeVersions.find(n => n.version === requestedVersion.value);
    if (found)
      return found;
    // Fallback to default if requested version not found
  }

  // Return default (no version) or first
  return nodeVersions.find(n => !n.version) || nodeVersions[0];
}

// Load node data
async function loadNode() {
  if (!projectId.value || !nodeName.value) {
    return;
  }

  isLoading.value = true;

  try {
    // Ensure nodes are loaded
    await nodeStore.ensureLoaded();

    // Find the node from store
    const foundNode = findCurrentNode();

    if (foundNode) {
      // Fetch fresh data from API using the node's public ID
      const result = await nodeService.getNode(projectId.value, foundNode.id);
      if (result) {
        node.value = result;
        // Initialize form data from node
        formData.value = {
          name: result.name,
          enabled: result.enabled,
          tags: [...result.tags],
          triggers: [...result.triggers],
          messages: [...result.messages],
          priority: result.priority ?? 0,
          condition: result.condition,
          effect: result.effect,
          nextAction: result.nextAction,
          forceCondition: result.forceCondition ?? false,
        };

        // Update URL if lang is not in query params (ensure URL has lang)
        if (!requestedLang.value && result.lang) {
          router.replace({
            query: {
              ...route.query,
              lang: result.lang,
            },
          });
        }
      }
      else {
        node.value = null;
      }
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
  if (!projectId.value || !node.value) {
    return;
  }

  try {
    const result = await nodeService.updateNode(projectId.value, node.value.id, {
      name: formData.value.name,
      enabled: formData.value.enabled,
      tags: formData.value.tags,
      triggers: formData.value.triggers,
      messages: formData.value.messages,
      priority: formData.value.priority,
      condition: formData.value.condition,
      effect: formData.value.effect,
      nextAction: formData.value.nextAction,
      forceCondition: formData.value.forceCondition,
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
  const remainingNodes = nodeStore.sortedNodes.filter(n => n.id !== node.value?.id);
  if (remainingNodes.length > 0 && remainingNodes[0]) {
    const firstNode = remainingNodes[0];
    router.push(`/platform/node/${firstNode.name}?lang=${firstNode.lang}`);
  }
  else {
    router.push('/platform');
  }
}

// Special value for default version (empty string not allowed by SelectItem)
const DEFAULT_VERSION_VALUE = '__default__';

// Handle version selection - find the language that has this version
function selectVersion(version: string) {
  const actualVersion = version === DEFAULT_VERSION_VALUE ? undefined : version;

  // Find a node with this version to get its language
  const nodeWithVersion = nodeStore.findNodeWithVersion(nodeName.value, actualVersion);
  const lang = nodeWithVersion?.lang || currentLang.value;

  if (actualVersion) {
    router.push({
      query: { lang, v: actualVersion },
    });
  }
  else {
    // Navigate to default version (remove v query param)
    router.push({
      query: { lang },
    });
  }
}

// Handle language selection - keep version if it exists in that language
function selectLanguage(lang: string) {
  // Check if current version exists in the selected language
  const nodeInLang = nodeStore.getNodeByNameVersion(nodeName.value, lang, requestedVersion.value);

  if (nodeInLang) {
    // Version exists in this language, keep it
    if (requestedVersion.value) {
      router.push({ query: { lang, v: requestedVersion.value } });
    }
    else {
      router.push({ query: { lang } });
    }
  }
  else {
    // Version doesn't exist, navigate to default in that language
    router.push({ query: { lang } });
  }
}

// Get the select value for a version (handles default)
function getVersionSelectValue(version: string | undefined): string {
  return version || DEFAULT_VERSION_VALUE;
}

// Add version dialog state
const showAddVersionDialog = ref(false);

// Handle version creation
function handleVersionCreated(newNode: ChatbotNode) {
  showAddVersionDialog.value = false;
  // Navigate to the new version with its language
  router.push({
    query: {
      lang: newNode.lang,
      ...(newNode.version ? { v: newNode.version } : {}),
    },
  });
}

// Watch for route changes and reload node
watch(
  [() => route.params.name, () => route.query.lang, () => route.query.v],
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
      <div class="space-y-4">
        <!-- Title row -->
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 min-w-0">
            <h1 class="text-2xl font-bold truncate">
              {{ displayName }}
            </h1>
            <Badge :variant="formData.enabled ? 'default' : 'secondary'" class="shrink-0">
              {{ formData.enabled ? t('common.enabled') : t('common.disabled') }}
            </Badge>
            <!-- System badge for predefined nodes - based on actual node data -->
            <TooltipProvider v-if="isPredefined">
              <Tooltip>
                <TooltipTrigger as-child>
                  <Badge variant="outline" class="gap-1 shrink-0">
                    <ShieldCheck class="h-3 w-3" />
                    {{ t('features.chatbotNode.badge.system') }}
                  </Badge>
                </TooltipTrigger>
                <TooltipContent>
                  <p>{{ t('features.chatbotNode.badge.systemTooltip') }}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </div>

          <!-- Action buttons -->
          <div class="flex items-center gap-2 shrink-0">
            <Button
              :disabled="!hasChanges || isSaving"
              @click="handleSave"
            >
              <Check v-if="!isSaving" class="h-4 w-4 mr-2" />
              <Loader2 v-else class="h-4 w-4 mr-2 animate-spin" />
              {{ t('common.save') }}
            </Button>
            <NodeDeleteDialog :node="node" @deleted="handleDeleted" />
          </div>
        </div>

        <!-- Language and Version row -->
        <div class="flex items-center justify-between">
          <!-- Language selector -->
          <div class="flex items-center gap-2">
            <Label class="text-sm text-muted-foreground">
              {{ t('features.chatbotNode.form.lang') }}:
            </Label>
            <Select
              :model-value="node.lang"
              @update:model-value="(v) => selectLanguage(String(v))"
            >
              <SelectTrigger class="w-[120px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="lang in availableLanguages"
                  :key="lang"
                  :value="lang"
                >
                  {{ lang }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div class="flex items-center gap-2">
            <!-- Version selector (if multiple versions exist) -->
            <Select
              v-if="hasVersions"
              :model-value="getVersionSelectValue(node.version)"
              @update:model-value="(v) => selectVersion(String(v))"
            >
              <SelectTrigger class="w-[140px]">
                <SelectValue :placeholder="t('features.chatbotNode.form.version')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="version in uniqueVersions"
                  :key="version || '__default__'"
                  :value="getVersionSelectValue(version || undefined)"
                >
                  {{ version || t('features.chatbotNode.defaultVersion') }}
                </SelectItem>
              </SelectContent>
            </Select>

            <!-- Add Version button -->
            <Button
              variant="outline"
              size="sm"
              @click="showAddVersionDialog = true"
            >
              <Plus class="h-4 w-4 mr-1" />
              {{ t('features.chatbotNode.addVersion') }}
            </Button>
          </div>
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
                :disabled="isSaving || isPredefined"
                class="max-w-xs"
              />
            </div>

            <div class="flex items-center gap-6">
              <!-- Priority field -->
              <div class="space-y-1">
                <Label for="node-priority">{{ t('features.chatbotNode.form.priority') }}</Label>
                <Input
                  id="node-priority"
                  v-model.number="formData.priority"
                  type="number"
                  min="0"
                  max="1000"
                  :placeholder="t('features.chatbotNode.form.priorityPlaceholder')"
                  :disabled="isSaving"
                  class="w-24"
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

          <!-- Conditions section (collapsible) -->
          <Collapsible class="space-y-4">
            <CollapsibleTrigger class="flex items-center justify-between w-full group">
              <div class="flex items-center gap-2">
                <h3 class="text-lg font-semibold">
                  {{ t('features.chatbotNode.condition.title') }}
                </h3>
                <Badge v-if="formData.condition" variant="secondary" class="text-xs">
                  Active
                </Badge>
                <Badge v-if="formData.forceCondition" variant="outline" class="text-xs">
                  {{ t('features.chatbotNode.form.forceCondition') }}
                </Badge>
              </div>
              <ChevronDown
                class="h-4 w-4 text-muted-foreground transition-transform
                  group-data-[state=open]:rotate-180"
              />
            </CollapsibleTrigger>
            <CollapsibleContent class="space-y-4">
              <!-- Force Condition toggle -->
              <div class="flex items-center justify-between p-3 rounded-lg border bg-muted/50">
                <div class="space-y-0.5">
                  <Label for="force-condition" class="text-sm font-medium">
                    {{ t('features.chatbotNode.form.forceCondition') }}
                  </Label>
                  <p class="text-xs text-muted-foreground">
                    {{ t('features.chatbotNode.form.forceConditionHelp') }}
                  </p>
                </div>
                <Switch
                  id="force-condition"
                  :model-value="formData.forceCondition"
                  :disabled="isSaving"
                  @update:model-value="formData.forceCondition = $event"
                />
              </div>

              <ConditionBuilder
                :condition="formData.condition"
                :disabled="isSaving"
                @update:condition="formData.condition = $event"
              />
            </CollapsibleContent>
          </Collapsible>

          <Separator />

          <!-- Messages section (chat bubbles) -->
          <MessageEditor
            v-model:messages="formData.messages"
            :disabled="isSaving"
          />

          <Separator />

          <!-- Effect section (collapsible) - Immediate actions after response -->
          <Collapsible class="space-y-4">
            <CollapsibleTrigger class="flex items-center justify-between w-full group">
              <div class="flex items-center gap-2">
                <h3 class="text-lg font-semibold">
                  {{ t('features.chatbotNode.effect.title') }}
                </h3>
                <Badge v-if="formData.effect" variant="secondary" class="text-xs">
                  {{ formData.effect.type }}
                </Badge>
              </div>
              <ChevronDown
                class="h-4 w-4 text-muted-foreground transition-transform
                  group-data-[state=open]:rotate-180"
              />
            </CollapsibleTrigger>
            <CollapsibleContent>
              <EffectEditor
                :effect="formData.effect"
                :disabled="isSaving"
                :nodes="nodeStore.sortedNodes"
                @update:effect="formData.effect = $event"
              />
            </CollapsibleContent>
          </Collapsible>

          <Separator />

          <!-- Next Action section (collapsible) - Deferred actions on user reply -->
          <Collapsible class="space-y-4">
            <CollapsibleTrigger class="flex items-center justify-between w-full group">
              <div class="flex items-center gap-2">
                <h3 class="text-lg font-semibold">
                  {{ t('features.chatbotNode.nextAction.title') }}
                </h3>
                <Badge v-if="formData.nextAction" variant="secondary" class="text-xs">
                  {{ formData.nextAction.type }}
                </Badge>
              </div>
              <ChevronDown
                class="h-4 w-4 text-muted-foreground transition-transform
                  group-data-[state=open]:rotate-180"
              />
            </CollapsibleTrigger>
            <CollapsibleContent>
              <NextActionEditor
                :next-action="formData.nextAction"
                :disabled="isSaving"
                :nodes="nodeStore.sortedNodes"
                @update:next-action="formData.nextAction = $event"
              />
            </CollapsibleContent>
          </Collapsible>
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

      <!-- Add Version Dialog -->
      <NodeAddVersionDialog
        v-if="node"
        v-model:open="showAddVersionDialog"
        :node-name="node.name"
        :node-lang="node.lang"
        @success="handleVersionCreated"
        @cancel="showAddVersionDialog = false"
      />
    </div>
  </div>
</template>
