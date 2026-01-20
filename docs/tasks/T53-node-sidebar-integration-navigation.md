# Task T53: Node Sidebar Integration & Navigation

**Story Reference:** US12-node-editor.md
**Type:** Frontend Integration
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** T52 (Node Edit Page & Components)

## Objective

Integrate node list into sidebar navigation with dynamic menu items, create node button, and proper navigation/breadcrumb support. Add i18n translations.

## Acceptance Criteria

- [ ] "Node" appears as parent menu item under Platform group in sidebar
- [ ] Each node appears as child menu item formatted as `{name}_{lang}`
- [ ] Nodes sorted alphabetically ascending
- [ ] "Create Node" button (+ icon) in Node section header
- [ ] Clicking create button opens NodeCreateSheet
- [ ] Sidebar updates when nodes are created/deleted
- [ ] Breadcrumb shows: Platform > Nodes > {node_name}_{lang}
- [ ] i18n translations added (en-US, id-ID)

## Technical Requirements

### Sidebar Integration

There are two approaches depending on the existing sidebar architecture:

#### Option A: Extend useNavigationItems.ts

If the sidebar uses a navigation items composable:

```typescript
// In useNavigationItems.ts
import { FileText, Plus } from 'lucide-vue-next';

// Add to mainNavItems after Chatbot section
{
  title: t('nav.node.title'),
  icon: FileText,
  items: computed(() => {
    const nodeStore = useNodeStore();
    return nodeStore.sortedNodes.map((node) => ({
      title: `${node.name}_${node.lang}`,
      to: `/platform/node/${node.id}`,
      icon: FileText,
    }));
  }),
  action: {
    icon: Plus,
    label: t('nav.node.create'),
    onClick: () => {
      // Open create sheet
      const event = new CustomEvent('open-node-create-sheet');
      window.dispatchEvent(event);
    },
  },
}
```

#### Option B: Custom SidebarNodeMenu.vue Component

If sidebar needs custom rendering for action buttons:

```vue
<script setup lang="ts">
import { FileText, Plus } from 'lucide-vue-next';

const nodeStore = useNodeStore();
const { fetchNodes, listLoading } = useNodeService();
const projectStore = useProjectStore();
const { t } = useI18n();

const createSheetRef = ref<InstanceType<typeof NodeCreateSheet> | null>(null);

// Fetch nodes on mount and project change
onMounted(async () => {
  if (projectStore.activeProjectId) {
    await fetchNodes(projectStore.activeProjectId);
  }
});

watch(() => projectStore.activeProjectId, async (newId) => {
  if (newId) {
    await fetchNodes(newId);
  }
});

function openCreateSheet() {
  createSheetRef.value?.open();
}

function handleNodeCreated() {
  // Store automatically updated via service
}
</script>

<template>
  <SidebarMenuItem>
    <Collapsible default-open class="group/collapsible">
      <SidebarMenuButton as-child>
        <CollapsibleTrigger class="flex w-full items-center">
          <FileText class="h-4 w-4" />
          <span class="flex-1 text-left">{{ t('nav.node.title') }}</span>
          <Button
            variant="ghost"
            size="icon"
            class="h-6 w-6"
            @click.stop="openCreateSheet"
          >
            <Plus class="h-3 w-3" />
          </Button>
          <ChevronRight class="h-4 w-4 transition-transform group-data-[state=open]/collapsible:rotate-90" />
        </CollapsibleTrigger>
      </SidebarMenuButton>

      <CollapsibleContent>
        <SidebarMenuSub>
          <SidebarMenuSubItem v-if="listLoading">
            <span class="text-sm text-muted-foreground">{{ t('common.loading') }}</span>
          </SidebarMenuSubItem>

          <SidebarMenuSubItem
            v-for="node in nodeStore.sortedNodes"
            :key="node.id"
          >
            <SidebarMenuSubButton as-child>
              <NuxtLink :to="`/platform/node/${node.id}`">
                <FileText class="h-3 w-3" />
                <span>{{ node.name }}_{{ node.lang }}</span>
              </NuxtLink>
            </SidebarMenuSubButton>
          </SidebarMenuSubItem>

          <SidebarMenuSubItem v-if="!listLoading && nodeStore.sortedNodes.length === 0">
            <span class="text-sm text-muted-foreground">{{ t('nav.node.empty') }}</span>
          </SidebarMenuSubItem>
        </SidebarMenuSub>
      </CollapsibleContent>
    </Collapsible>
  </SidebarMenuItem>

  <NodeCreateSheet ref="createSheetRef" @created="handleNodeCreated" />
</template>
```

### i18n Translations

#### English (en-US.json)

```json
{
  "nav": {
    "node": {
      "title": "Nodes",
      "create": "Create Node",
      "empty": "No nodes yet"
    }
  },
  "features": {
    "node": {
      "page": {
        "title": "Edit Node",
        "createTitle": "Create Node"
      },
      "form": {
        "name": "Node Name",
        "namePlaceholder": "e.g., greeting, faq_pricing",
        "nameHelp": "Use lowercase letters, numbers, and underscores only",
        "lang": "Language",
        "tags": "Tags",
        "tagsPlaceholder": "Add tags...",
        "enabled": "Enabled",
        "triggers": "Triggers",
        "triggersHelp": "At least one trigger required",
        "messages": "Messages",
        "messagesHelp": "At least one message required",
        "save": "Save Node",
        "saving": "Saving..."
      },
      "trigger": {
        "type": "Type",
        "value": "Value",
        "types": {
          "keyword": "Keyword",
          "contains": "Contains",
          "regex": "Regex",
          "equals": "Equals"
        },
        "add": "Add Trigger",
        "remove": "Remove"
      },
      "message": {
        "content": "Message Content",
        "add": "Add Message",
        "remove": "Remove",
        "moveUp": "Move Up",
        "moveDown": "Move Down"
      },
      "delete": {
        "title": "Delete Node",
        "description": "Are you sure you want to delete the node \"{name}\"? This action cannot be undone.",
        "confirm": "Delete",
        "cancel": "Cancel"
      },
      "messages": {
        "createSuccess": "Node created successfully",
        "saveSuccess": "Node saved successfully",
        "deleteSuccess": "Node deleted successfully",
        "createError": "Failed to create node",
        "saveError": "Failed to save node",
        "deleteError": "Failed to delete node",
        "invalidRegex": "Invalid regex pattern"
      }
    }
  }
}
```

#### Indonesian (id-ID.json)

```json
{
  "nav": {
    "node": {
      "title": "Node",
      "create": "Buat Node",
      "empty": "Belum ada node"
    }
  },
  "features": {
    "node": {
      "page": {
        "title": "Edit Node",
        "createTitle": "Buat Node"
      },
      "form": {
        "name": "Nama Node",
        "namePlaceholder": "contoh: greeting, faq_pricing",
        "nameHelp": "Gunakan huruf kecil, angka, dan garis bawah saja",
        "lang": "Bahasa",
        "tags": "Tag",
        "tagsPlaceholder": "Tambah tag...",
        "enabled": "Aktif",
        "triggers": "Pemicu",
        "triggersHelp": "Minimal satu pemicu diperlukan",
        "messages": "Pesan",
        "messagesHelp": "Minimal satu pesan diperlukan",
        "save": "Simpan Node",
        "saving": "Menyimpan..."
      },
      "trigger": {
        "type": "Tipe",
        "value": "Nilai",
        "types": {
          "keyword": "Kata Kunci",
          "contains": "Mengandung",
          "regex": "Regex",
          "equals": "Sama Dengan"
        },
        "add": "Tambah Pemicu",
        "remove": "Hapus"
      },
      "message": {
        "content": "Isi Pesan",
        "add": "Tambah Pesan",
        "remove": "Hapus",
        "moveUp": "Naik",
        "moveDown": "Turun"
      },
      "delete": {
        "title": "Hapus Node",
        "description": "Apakah Anda yakin ingin menghapus node \"{name}\"? Tindakan ini tidak dapat dibatalkan.",
        "confirm": "Hapus",
        "cancel": "Batal"
      },
      "messages": {
        "createSuccess": "Node berhasil dibuat",
        "saveSuccess": "Node berhasil disimpan",
        "deleteSuccess": "Node berhasil dihapus",
        "createError": "Gagal membuat node",
        "saveError": "Gagal menyimpan node",
        "deleteError": "Gagal menghapus node",
        "invalidRegex": "Pola regex tidak valid"
      }
    }
  }
}
```

### Sidebar Data Flow

```
Project Change
     │
     ▼
useNodeStore.reset()
     │
     ▼
fetchNodes(projectId)
     │
     ▼
nodeStore.setNodes()
     │
     ▼
sortedNodes (computed)
     │
     ▼
SidebarNodeMenu re-renders
```

## Files to Create

```
frontend/app/components/custom/sidebar/
└── SidebarNodeMenu.vue    # Custom sidebar section (if needed)
```

## Files to Modify

- `frontend/app/composables/navigation/useNavigationItems.ts` - Add Node menu with dynamic children (Option A)
- `frontend/app/components/custom/layout/LayoutSidebar.vue` - Integrate SidebarNodeMenu (Option B)
- `frontend/i18n/locales/en-US.json` - Add node translations
- `frontend/i18n/locales/id-ID.json` - Add node translations (Indonesian)

## Commands to Run

```bash
cd frontend

# Type check
pnpm typecheck

# Lint
pnpm lint

# Dev server to test
pnpm dev
```

## Validation Checklist

- [ ] "Node" menu appears in sidebar
- [ ] Node items display as `{name}_{lang}` format
- [ ] Nodes sorted alphabetically (ascending)
- [ ] Create button opens NodeCreateSheet
- [ ] New nodes appear in sidebar after creation
- [ ] Deleted nodes removed from sidebar
- [ ] Navigation to `/platform/node/{id}` works
- [ ] Breadcrumb displays correctly
- [ ] All i18n keys work in both languages
- [ ] Mobile responsive sidebar works
- [ ] Sidebar updates on project change

## Definition of Done

- [ ] Sidebar shows dynamic Node menu
- [ ] Create button functional
- [ ] Navigation working
- [ ] Store reactivity verified
- [ ] i18n translations complete (en-US, id-ID)
- [ ] Breadcrumb working
- [ ] Mobile responsive
- [ ] No console errors

## Dependencies

- T52: All node components must exist
- T51: Node store and service must be implemented
- Existing sidebar infrastructure

## Risk Factors

- **Medium Risk**: Sidebar architecture may need custom component approach
- **Low Risk**: i18n translations are straightforward

## Notes

- Choose Option A or B based on existing sidebar implementation
- Store reactivity ensures sidebar updates automatically
- Consider debouncing node fetch on rapid project switches
- Create button should be easily accessible but not disruptive
- Test with many nodes (10+) for scrolling behavior
