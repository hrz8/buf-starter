import type { ChatbotNode } from '~~/gen/chatbot/nodes/v1/node_pb';
import { chatbotNodeRepository } from '#shared/repository/chatbot-node';
import { create } from '@bufbuild/protobuf';
import { ListNodesRequestSchema } from '~~/gen/altalune/v1/chatbot_node_pb';
import { useProjectStore } from './project';

export const useChatbotNodeStore = defineStore('chatbot-node', () => {
  const { $chatbotNodeClient } = useNuxtApp();
  const nodeRepo = chatbotNodeRepository($chatbotNodeClient);
  const projectStore = useProjectStore();

  // State
  const nodes = ref<ChatbotNode[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);
  const initialized = ref(false);
  const lastFetchedProjectId = ref<string | null>(null);

  // Computed: sorted nodes by name_lang
  const sortedNodes = computed(() => {
    return [...nodes.value].sort((a, b) => {
      const aKey = `${a.name}_${a.lang}`;
      const bKey = `${b.name}_${b.lang}`;
      return aKey.localeCompare(bKey);
    });
  });

  // Getters
  function getNodeById(nodeId: string): ChatbotNode | undefined {
    return nodes.value.find(n => n.id === nodeId);
  }

  function getNodeDisplayName(node: ChatbotNode): string {
    return `${node.name}_${node.lang}`;
  }

  // Actions
  async function fetchNodes(projectId: string): Promise<void> {
    // Skip if already fetched for this project
    if (initialized.value && lastFetchedProjectId.value === projectId) {
      return;
    }

    loading.value = true;
    error.value = null;

    try {
      const message = create(ListNodesRequestSchema, { projectId });
      const result = await nodeRepo.listNodes(message);
      nodes.value = result.nodes || [];
      initialized.value = true;
      lastFetchedProjectId.value = projectId;
    }
    catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch nodes';
      throw err;
    }
    finally {
      loading.value = false;
    }
  }

  // Force refresh nodes (used after create/delete)
  async function refreshNodes(): Promise<void> {
    const projectId = projectStore.activeProjectId;
    if (!projectId)
      return;

    // Reset initialized to force fetch
    initialized.value = false;
    lastFetchedProjectId.value = null;
    await fetchNodes(projectId);
  }

  // Add a node locally (optimistic update after create)
  function addNode(node: ChatbotNode): void {
    nodes.value = [...nodes.value, node];
  }

  // Remove a node locally (optimistic update after delete)
  function removeNode(nodeId: string): void {
    nodes.value = nodes.value.filter(n => n.id !== nodeId);
  }

  // Update a node locally (optimistic update after update)
  function updateNode(updatedNode: ChatbotNode): void {
    const index = nodes.value.findIndex(n => n.id === updatedNode.id);
    if (index !== -1) {
      nodes.value = [
        ...nodes.value.slice(0, index),
        updatedNode,
        ...nodes.value.slice(index + 1),
      ];
    }
  }

  // Reset store
  function reset(): void {
    nodes.value = [];
    initialized.value = false;
    lastFetchedProjectId.value = null;
    error.value = null;
  }

  // Watch for project changes and reset
  watch(
    () => projectStore.activeProjectId,
    (newProjectId, oldProjectId) => {
      if (newProjectId !== oldProjectId) {
        reset();
      }
    },
  );

  // Ensure nodes are loaded for current project
  async function ensureLoaded(): Promise<void> {
    const projectId = projectStore.activeProjectId;
    if (!projectId)
      return;

    if (!initialized.value || lastFetchedProjectId.value !== projectId) {
      await fetchNodes(projectId);
    }
  }

  return {
    // State (readonly)
    nodes: readonly(nodes),
    sortedNodes,
    loading: readonly(loading),
    error: readonly(error),
    initialized: readonly(initialized),

    // Getters
    getNodeById,
    getNodeDisplayName,

    // Actions
    fetchNodes,
    refreshNodes,
    addNode,
    removeNode,
    updateNode,
    ensureLoaded,
    reset,
  };
});
