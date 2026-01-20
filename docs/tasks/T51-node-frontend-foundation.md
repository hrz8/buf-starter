# Task T51: Node Frontend Foundation

**Story Reference:** US12-node-editor.md
**Type:** Frontend Foundation
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T50 (Node Backend Domain)

## Objective

Implement repository, service composable, Zod schemas, constants, and error utilities for chatbot nodes. Follow the pattern established in `frontend/app/lib/chatbot-modules/` for consistency.

## Acceptance Criteria

- [ ] Repository implements all CRUD API calls (listNodes, createNode, getNode, updateNode, deleteNode)
- [ ] Service composable (useNodeService) provides reactive state and methods
- [ ] Zod schemas define validation for node form, trigger, message
- [ ] Constants define trigger types, language options, default values
- [ ] Error utilities handle ConnectRPC errors for node operations
- [ ] Types exported for Node, Trigger, Message
- [ ] Node store (Pinia) manages node list state for sidebar

## Technical Requirements

### Repository Pattern (chatbot-node.ts)

```typescript
import type { Client } from '@connectrpc/connect';
import type { ChatbotNodeService } from '~~/gen/altalune/v1/chatbot_node_pb';

export function chatbotNodeRepository(client: Client<typeof ChatbotNodeService>) {
  return {
    async listNodes(projectId: string) {
      return client.listNodes({ projectId });
    },

    async createNode(input: {
      projectId: string;
      name: string;
      lang: string;
      tags?: string[];
    }) {
      return client.createNode(input);
    },

    async getNode(projectId: string, nodeId: string) {
      return client.getNode({ projectId, nodeId });
    },

    async updateNode(input: {
      projectId: string;
      nodeId: string;
      name?: string;
      tags?: string[];
      enabled?: boolean;
      triggers?: Array<{ type: string; value: string }>;
      messages?: Array<{ role: string; content: string }>;
    }) {
      return client.updateNode(input);
    },

    async deleteNode(projectId: string, nodeId: string) {
      return client.deleteNode({ projectId, nodeId });
    },
  };
}
```

### Service Composable Pattern (useNodeService.ts)

```typescript
export function useNodeService() {
  const nodeStore = useNodeStore();
  const { $chatbotNodeClient } = useNuxtApp();
  const repository = chatbotNodeRepository($chatbotNodeClient);
  const { parseError } = useErrorMessage();

  const listState = reactive({ loading: false, error: '', success: false });
  const createState = reactive({ loading: false, error: '', success: false });
  const getState = reactive({ loading: false, error: '', success: false });
  const updateState = reactive({ loading: false, error: '', success: false });
  const deleteState = reactive({ loading: false, error: '', success: false });

  async function fetchNodes(projectId: string) {
    listState.loading = true;
    listState.error = '';
    try {
      const result = await repository.listNodes(projectId);
      nodeStore.setNodes(result.nodes);
      listState.success = true;
      return result.nodes;
    } catch (err) {
      listState.error = parseError(err);
      throw err;
    } finally {
      listState.loading = false;
    }
  }

  // ... other methods: createNode, getNode, updateNode, deleteNode

  return {
    fetchNodes,
    createNode,
    getNode,
    updateNode,
    deleteNode,
    listLoading: computed(() => listState.loading),
    listError: computed(() => listState.error),
    // ... other computed states
  };
}
```

### Pinia Store Pattern (useNodeStore.ts)

```typescript
export const useNodeStore = defineStore('chatbotNode', () => {
  const nodes = ref<Map<string, ChatbotNode>>(new Map());
  const loading = ref(false);
  const error = ref<string | null>(null);
  const initialized = ref(false);
  const lastFetchedProjectId = ref<string | null>(null);

  const nodesList = computed(() => Array.from(nodes.value.values()));

  const sortedNodes = computed(() =>
    [...nodesList.value].sort((a, b) =>
      `${a.name}_${a.lang}`.localeCompare(`${b.name}_${b.lang}`)
    )
  );

  function setNodes(nodeList: ChatbotNode[]) {
    nodes.value.clear();
    nodeList.forEach(node => nodes.value.set(node.id, node));
  }

  function addNode(node: ChatbotNode) {
    nodes.value.set(node.id, node);
  }

  function updateNode(nodeId: string, updates: Partial<ChatbotNode>) {
    const existing = nodes.value.get(nodeId);
    if (existing) {
      nodes.value.set(nodeId, { ...existing, ...updates });
    }
  }

  function removeNode(nodeId: string) {
    nodes.value.delete(nodeId);
  }

  function reset() {
    nodes.value.clear();
    initialized.value = false;
    lastFetchedProjectId.value = null;
  }

  // Watch project changes to reset store
  const projectStore = useProjectStore();
  watch(() => projectStore.activeProjectId, (newId, oldId) => {
    if (newId !== oldId) reset();
  });

  return {
    nodesList: readonly(nodesList),
    sortedNodes: readonly(sortedNodes),
    loading: readonly(loading),
    error: readonly(error),
    initialized: readonly(initialized),
    setNodes,
    addNode,
    updateNode,
    removeNode,
    reset,
  };
});
```

### Constants (constants.ts)

```typescript
export const TRIGGER_TYPES = ['keyword', 'contains', 'regex', 'equals'] as const;
export type TriggerType = (typeof TRIGGER_TYPES)[number];

export const TRIGGER_TYPE_OPTIONS = [
  { value: 'keyword', label: 'Keyword' },
  { value: 'contains', label: 'Contains' },
  { value: 'regex', label: 'Regex' },
  { value: 'equals', label: 'Equals' },
] as const;

export const LANGUAGES = ['en-US', 'id-ID'] as const;
export type Language = (typeof LANGUAGES)[number];

export const LANGUAGE_OPTIONS = [
  { value: 'en-US', label: 'English (US)' },
  { value: 'id-ID', label: 'Indonesian (ID)' },
] as const;

export const DEFAULT_TRIGGER = { type: 'keyword' as TriggerType, value: '' };
export const DEFAULT_MESSAGE = { role: 'assistant', content: '' };
```

### Zod Schemas (schema.ts)

```typescript
import { z } from 'zod';
import { TRIGGER_TYPES, LANGUAGES } from './constants';

export const triggerSchema = z.object({
  type: z.enum(TRIGGER_TYPES),
  value: z.string().min(1, 'Value is required').max(500, 'Max 500 characters'),
});

export const messageSchema = z.object({
  role: z.literal('assistant'),
  content: z.string().min(1, 'Content is required').max(5000, 'Max 5000 characters'),
});

export const nodeCreateSchema = z.object({
  name: z.string()
    .min(2, 'Min 2 characters')
    .max(100, 'Max 100 characters')
    .regex(/^[a-z][a-z0-9_]*$/, 'Must be lowercase_snake_case'),
  lang: z.enum(['en-US', 'id-ID']),
  tags: z.array(z.string().min(1).max(50)).optional(),
});

export const nodeEditSchema = z.object({
  name: z.string()
    .min(2, 'Min 2 characters')
    .max(100, 'Max 100 characters')
    .regex(/^[a-z][a-z0-9_]*$/, 'Must be lowercase_snake_case'),
  tags: z.array(z.string().min(1).max(50)).optional(),
  enabled: z.boolean(),
  triggers: z.array(triggerSchema).min(1, 'At least one trigger required'),
  messages: z.array(messageSchema).min(1, 'At least one message required'),
});

export type NodeCreateInput = z.infer<typeof nodeCreateSchema>;
export type NodeEditInput = z.infer<typeof nodeEditSchema>;
export type TriggerInput = z.infer<typeof triggerSchema>;
export type MessageInput = z.infer<typeof messageSchema>;
```

### Error Utilities (error.ts)

```typescript
import { ConnectError } from '@connectrpc/connect';

export const NODE_ERROR_CODES = {
  NOT_FOUND: 62001,
  INVALID_NAME: 62002,
  INVALID_LANGUAGE: 62003,
  DUPLICATE_NAME: 62004,
  NO_TRIGGERS: 62005,
  NO_MESSAGES: 62006,
  INVALID_TRIGGER_TYPE: 62007,
  INVALID_REGEX: 62008,
} as const;

export function getNodeError(error: unknown): { code: number; message: string } | null {
  if (error instanceof ConnectError) {
    const match = error.message.match(/\[(\d+)\]/);
    if (match) {
      return {
        code: parseInt(match[1]),
        message: error.message.replace(/\[\d+\]\s*/, ''),
      };
    }
  }
  return null;
}

export function hasNodeError(error: unknown, code: number): boolean {
  const nodeError = getNodeError(error);
  return nodeError?.code === code;
}
```

## Files to Create

```
frontend/shared/repository/chatbot-node.ts

frontend/app/components/features/chatbot-node/
├── schema.ts              # Zod validation schemas
├── error.ts               # Error utilities
├── constants.ts           # Trigger types, language options, defaults
├── types.ts               # TypeScript types for Node, Trigger, Message
└── index.ts               # Barrel exports

frontend/app/composables/services/useNodeService.ts
frontend/app/stores/useNodeStore.ts
```

## Files to Modify

- `frontend/app/plugins/connect.ts` - Add chatbotNodeClient creation
- `frontend/shared/repository/index.ts` - Export chatbot-node repository

## Commands to Run

```bash
cd frontend

# Generate proto types (should already exist from T50)
pnpm generate

# Type check
pnpm typecheck

# Lint
pnpm lint
```

## Validation Checklist

- [ ] Repository methods call API correctly
- [ ] Service composable integrates with store
- [ ] Zod schemas validate form data correctly
- [ ] Store maintains sorted node list
- [ ] Error utilities extract error messages properly
- [ ] TypeScript types compile without errors
- [ ] Nuxt plugin provides client via $chatbotNodeClient

## Definition of Done

- [ ] Repository file created with all CRUD methods
- [ ] Service composable provides reactive state
- [ ] Pinia store manages node list with sorting
- [ ] Zod schemas match backend validation
- [ ] Constants defined for trigger types and languages
- [ ] Error utilities handle all node error codes
- [ ] Nuxt plugin configured for chatbot node client
- [ ] All files type-check successfully

## Dependencies

- T50: Backend must provide API endpoints
- Generated TypeScript types from proto

## Risk Factors

- **Low Risk**: Following established patterns from chatbot module
- **Medium Risk**: Store reactivity for sidebar updates

## Notes

- Store uses Map for O(1) lookups by node ID
- sortedNodes computed property for sidebar display
- Reset store on project change to avoid stale data
- Service composable manages loading/error states per operation
