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

  // Computed: sorted nodes by name, then lang, with default versions first
  const sortedNodes = computed(() => {
    return [...nodes.value].sort((a, b) => {
      // First by name
      const nameCompare = a.name.localeCompare(b.name);
      if (nameCompare !== 0)
        return nameCompare;
      // Then by lang
      const langCompare = a.lang.localeCompare(b.lang);
      if (langCompare !== 0)
        return langCompare;
      // Same name+lang: default (no version) first, then by version name
      if (!a.version && b.version)
        return -1;
      if (a.version && !b.version)
        return 1;
      return (a.version || '').localeCompare(b.version || '');
    });
  });

  // Computed: nodes grouped by name only (for sidebar navigation)
  const nodesByName = computed(() => {
    const groups = new Map<string, ChatbotNode[]>();
    for (const node of nodes.value) {
      const existing = groups.get(node.name) || [];
      existing.push(node);
      groups.set(node.name, existing);
    }
    // Sort each group: by lang, then default version first, then version name
    for (const [key, group] of groups.entries()) {
      groups.set(
        key,
        group.sort((a, b) => {
          const langCompare = a.lang.localeCompare(b.lang);
          if (langCompare !== 0)
            return langCompare;
          if (!a.version && b.version)
            return -1;
          if (a.version && !b.version)
            return 1;
          return (a.version || '').localeCompare(b.version || '');
        }),
      );
    }
    return groups;
  });

  // Computed: nodes grouped by name_lang (for version selection within a language)
  const groupedNodes = computed(() => {
    const groups = new Map<string, ChatbotNode[]>();
    for (const node of nodes.value) {
      const key = `${node.name}_${node.lang}`;
      const existing = groups.get(key) || [];
      existing.push(node);
      groups.set(key, existing);
    }
    // Sort each group: default (no version) first, then by version name
    for (const [key, group] of groups.entries()) {
      groups.set(
        key,
        group.sort((a, b) => {
          if (!a.version && b.version)
            return -1;
          if (a.version && !b.version)
            return 1;
          return (a.version || '').localeCompare(b.version || '');
        }),
      );
    }
    return groups;
  });

  // Getters
  function getNodeById(nodeId: string): ChatbotNode | undefined {
    return nodes.value.find(n => n.id === nodeId);
  }

  function getNodeDisplayName(node: ChatbotNode): string {
    return `${node.name}_${node.lang}`;
  }

  // Get all nodes by name (across all languages)
  function getNodesByName(name: string): ChatbotNode[] {
    return nodesByName.value.get(name) || [];
  }

  // Get all versions of a node by name and lang
  function getNodeVersions(name: string, lang: string): ChatbotNode[] {
    const key = `${name}_${lang}`;
    return groupedNodes.value.get(key) || [];
  }

  // Get all unique versions across all languages for a node name
  function getUniqueVersions(name: string): string[] {
    const allNodes = getNodesByName(name);
    const versions = new Set<string>();
    for (const node of allNodes) {
      versions.add(node.version || ''); // empty string for default
    }
    // Sort: default (empty) first, then alphabetically
    return [...versions].sort((a, b) => {
      if (!a && b)
        return -1;
      if (a && !b)
        return 1;
      return a.localeCompare(b);
    });
  }

  // Check if a node has multiple versions (across all languages)
  function hasMultipleVersions(name: string, lang?: string): boolean {
    if (lang) {
      const versions = getNodeVersions(name, lang);
      return versions.length > 1 || versions.some(v => v.version);
    }
    // Check across all languages
    const uniqueVersions = getUniqueVersions(name);
    return uniqueVersions.length > 1 || uniqueVersions.some(v => v !== '');
  }

  // Find the first node with a specific version (returns the first language that has it)
  function findNodeWithVersion(name: string, version?: string): ChatbotNode | undefined {
    const allNodes = getNodesByName(name);
    if (version) {
      return allNodes.find(n => n.version === version);
    }
    // Return default (no version) or first
    return allNodes.find(n => !n.version) || allNodes[0];
  }

  // Get all languages that have a specific version
  function getLanguagesForVersion(name: string, version?: string): string[] {
    const allNodes = getNodesByName(name);
    const languages = new Set<string>();
    for (const node of allNodes) {
      const nodeVersion = node.version || '';
      const targetVersion = version || '';
      if (nodeVersion === targetVersion) {
        languages.add(node.lang);
      }
    }
    return [...languages].sort();
  }

  // Get a specific node by name, lang, and optional version
  function getNodeByNameVersion(
    name: string,
    lang: string,
    version?: string,
  ): ChatbotNode | undefined {
    const versions = getNodeVersions(name, lang);
    if (version) {
      return versions.find(n => n.version === version);
    }
    // Return default (no version) or first
    return versions.find(n => !n.version) || versions[0];
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
    nodesByName,
    groupedNodes,
    loading: readonly(loading),
    error: readonly(error),
    initialized: readonly(initialized),

    // Getters
    getNodeById,
    getNodeDisplayName,
    getNodesByName,
    getNodeVersions,
    getUniqueVersions,
    hasMultipleVersions,
    findNodeWithVersion,
    getLanguagesForVersion,
    getNodeByNameVersion,

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
