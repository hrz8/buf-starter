# Task T46: Chatbot Module Pages & Navigation

**Story Reference:** US11-chatbot-configuration-foundation.md
**Type:** Frontend Integration
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** T45 (Schema-Driven Form Components)

## Objective

Implement dynamic module configuration pages and integrate chatbot navigation into the sidebar with auto-discovered modules from the schema registry.

## Acceptance Criteria

- [ ] Dynamic route `/platform/modules/[name].vue` renders module configs
- [ ] Sidebar shows "Chatbot" menu with module children auto-loaded from schema registry
- [ ] Module enabled/disabled indicators shown in sidebar
- [ ] Module toggle updates enabled status via API
- [ ] Save button persists form changes
- [ ] Toast notifications for success/error feedback
- [ ] Breadcrumb navigation works: Platform > Chatbot > {Module Name}

## Technical Requirements

### Dynamic Module Page

File: `frontend/app/pages/platform/modules/[name].vue`

```vue
<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { toast } from 'vue-sonner';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import SchemaForm from '@/components/features/chatbot/SchemaForm.vue';
import ModuleToggle from '@/components/features/chatbot/ModuleToggle.vue';
import { getModuleSchema, type ModuleName } from '@/components/features/chatbot/schemas';
import { useChatbotService } from '@/composables/services/useChatbotService';
import { useProjectStore } from '@/stores/project';
import { useBreadcrumb } from '@/composables/useBreadcrumb';

definePageMeta({
  layout: 'platform',
});

const route = useRoute();
const router = useRouter();
const { t } = useI18n();
const projectStore = useProjectStore();

// Get module name from route
const moduleName = computed(() => route.params.name as string);

// Validate module name
const schema = computed(() => getModuleSchema(moduleName.value));

// Redirect if invalid module
watch(schema, (val) => {
  if (!val && moduleName.value) {
    router.replace('/platform/modules/llm');
  }
}, { immediate: true });

// Setup chatbot service
const projectId = computed(() => projectStore.currentProject?.id || '');
const { config, loading, fetchConfig, updateModule, getModuleConfig } = useChatbotService(projectId.value);

// Local state
const isLoading = ref(true);
const isSaving = ref(false);
const moduleConfig = ref<Record<string, any>>({});
const schemaFormRef = ref<InstanceType<typeof SchemaForm> | null>(null);

// Module enabled state
const moduleEnabled = computed({
  get: () => moduleConfig.value?.enabled ?? false,
  set: (val: boolean) => {
    moduleConfig.value = { ...moduleConfig.value, enabled: val };
  },
});

// Fetch config on mount
onMounted(async () => {
  if (projectId.value) {
    await fetchConfig();
    moduleConfig.value = getModuleConfig(moduleName.value as ModuleName);
    isLoading.value = false;
  }
});

// Watch for module changes
watch(moduleName, async () => {
  if (projectId.value) {
    moduleConfig.value = getModuleConfig(moduleName.value as ModuleName);
  }
});

// Handle toggle change
async function handleToggleChange(enabled: boolean) {
  isSaving.value = true;
  try {
    const updatedConfig = { ...moduleConfig.value, enabled };
    await updateModule(moduleName.value as ModuleName, updatedConfig);
    moduleConfig.value = updatedConfig;
    toast.success(t('features.chatbot.messages.toggleSuccess'));
  } catch (error) {
    toast.error(t('features.chatbot.messages.toggleError'));
  } finally {
    isSaving.value = false;
  }
}

// Handle form save
async function handleSave() {
  if (!schemaFormRef.value) return;

  const { valid } = await schemaFormRef.value.validate();
  if (!valid) {
    toast.error(t('features.chatbot.messages.validationError'));
    return;
  }

  isSaving.value = true;
  try {
    await updateModule(moduleName.value as ModuleName, moduleConfig.value);
    toast.success(t('features.chatbot.messages.saveSuccess'));
  } catch (error) {
    toast.error(t('features.chatbot.messages.saveError'));
  } finally {
    isSaving.value = false;
  }
}

// Breadcrumb
useBreadcrumb({
  items: computed(() => [
    { label: t('nav.platform'), to: '/platform' },
    { label: t('nav.chatbot.title'), to: '/platform/modules/llm' },
    { label: schema.value?.title || moduleName.value },
  ]),
});
</script>

<template>
  <div class="container max-w-3xl py-8">
    <!-- Loading State -->
    <template v-if="isLoading">
      <Skeleton class="h-8 w-48 mb-6" />
      <Skeleton class="h-24 w-full mb-6" />
      <Skeleton class="h-64 w-full" />
    </template>

    <!-- Content -->
    <template v-else-if="schema">
      <!-- Page Header -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold">{{ schema.title }}</h1>
        <p class="text-muted-foreground">
          {{ t('features.chatbot.page.description') }}
        </p>
      </div>

      <!-- Module Toggle -->
      <div class="mb-6">
        <ModuleToggle
          :title="schema.title"
          :icon="schema.icon"
          :enabled="moduleEnabled"
          :loading="isSaving"
          @update:enabled="handleToggleChange"
        />
      </div>

      <!-- Configuration Form -->
      <Card>
        <CardHeader>
          <CardTitle>{{ t('features.chatbot.page.configTitle') }}</CardTitle>
        </CardHeader>
        <CardContent>
          <SchemaForm
            ref="schemaFormRef"
            :schema="schema"
            v-model="moduleConfig"
            :disabled="!moduleEnabled || isSaving"
            :loading="isSaving"
          >
            <template #actions="{ disabled }">
              <div class="flex justify-end pt-4">
                <Button
                  type="submit"
                  :disabled="disabled || !moduleEnabled"
                  :loading="isSaving"
                  @click.prevent="handleSave"
                >
                  {{ isSaving
                    ? t('features.chatbot.form.saving')
                    : t('features.chatbot.form.save') }}
                </Button>
              </div>
            </template>
          </SchemaForm>
        </CardContent>
      </Card>
    </template>
  </div>
</template>
```

### Module Index Page (Optional Redirect)

File: `frontend/app/pages/platform/modules/index.vue`

```vue
<script setup lang="ts">
definePageMeta({
  layout: 'platform',
});

// Redirect to first module (LLM)
const router = useRouter();
onMounted(() => {
  router.replace('/platform/modules/llm');
});
</script>

<template>
  <div class="flex items-center justify-center h-64">
    <Spinner />
  </div>
</template>
```

### Sidebar Navigation Integration

The sidebar needs to show a "Chatbot" menu with auto-discovered module children.

**Option A: Modify useNavigationItems.ts**

File: `frontend/app/composables/navigation/useNavigationItems.ts`

```typescript
import { Bot, Brain, Server, Layout, MessageSquare } from 'lucide-vue-next';
import { getModuleList } from '@/components/features/chatbot/schemas';

// Module icon mapping
const moduleIcons: Record<string, Component> = {
  'lucide:brain': Brain,
  'lucide:server': Server,
  'lucide:layout': Layout,
  'lucide:message-square': MessageSquare,
};

export function useNavigationItems() {
  const { t } = useI18n();

  // Get chatbot modules from schema registry
  const chatbotModules = getModuleList().map((module) => ({
    title: module.title,
    to: `/platform/modules/${module.key}`,
    icon: moduleIcons[module.icon] || Brain,
  }));

  const navigationItems = computed(() => [
    // ... existing items ...

    // Chatbot section
    {
      title: t('nav.chatbot.title'),
      to: '/platform/modules',
      icon: Bot,
      items: chatbotModules,
    },

    // ... remaining items ...
  ]);

  return { navigationItems };
}
```

**Option B: Custom SidebarChatbotMenu Component**

If the sidebar structure doesn't support dynamic children well, create a dedicated component.

File: `frontend/app/components/custom/layout/SidebarChatbotMenu.vue`

```vue
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { Bot, ChevronDown } from 'lucide-vue-next';
import { getModuleList, type ModuleName } from '@/components/features/chatbot/schemas';
import { useChatbotService } from '@/composables/services/useChatbotService';
import { useProjectStore } from '@/stores/project';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
import { Badge } from '@/components/ui/badge';

const route = useRoute();
const projectStore = useProjectStore();
const projectId = computed(() => projectStore.currentProject?.id || '');

const { config, fetchConfig } = useChatbotService(projectId.value);

const isOpen = ref(true);
const modules = getModuleList();

// Check if module is enabled
function isModuleEnabled(key: ModuleName): boolean {
  return config.value[key]?.enabled ?? false;
}

// Check if route is active
function isActive(path: string): boolean {
  return route.path === path;
}

// Fetch config on mount
onMounted(() => {
  if (projectId.value) {
    fetchConfig();
  }
});
</script>

<template>
  <Collapsible v-model:open="isOpen" class="w-full">
    <CollapsibleTrigger class="flex w-full items-center justify-between rounded-lg px-3 py-2 text-sm hover:bg-accent">
      <div class="flex items-center gap-2">
        <Bot class="h-4 w-4" />
        <span>{{ $t('nav.chatbot.title') }}</span>
      </div>
      <ChevronDown
        class="h-4 w-4 transition-transform"
        :class="{ 'rotate-180': isOpen }"
      />
    </CollapsibleTrigger>
    <CollapsibleContent class="pl-4 pt-1">
      <NuxtLink
        v-for="module in modules"
        :key="module.key"
        :to="`/platform/modules/${module.key}`"
        class="flex items-center justify-between rounded-lg px-3 py-2 text-sm hover:bg-accent"
        :class="{ 'bg-accent': isActive(`/platform/modules/${module.key}`) }"
      >
        <span>{{ module.title }}</span>
        <Badge
          :variant="isModuleEnabled(module.key) ? 'default' : 'secondary'"
          class="text-xs"
        >
          {{ isModuleEnabled(module.key) ? 'On' : 'Off' }}
        </Badge>
      </NuxtLink>
    </CollapsibleContent>
  </Collapsible>
</template>
```

### i18n Translations

File: `frontend/i18n/locales/en-US.json`

```json
{
  "nav": {
    "chatbot": {
      "title": "Chatbot",
      "modules": "Modules"
    }
  },
  "features": {
    "chatbot": {
      "page": {
        "title": "Module Configuration",
        "description": "Configure the chatbot module settings below.",
        "configTitle": "Settings",
        "enabled": "Enabled",
        "disabled": "Disabled"
      },
      "form": {
        "save": "Save Configuration",
        "saving": "Saving..."
      },
      "messages": {
        "saveSuccess": "Configuration saved successfully",
        "saveError": "Failed to save configuration",
        "toggleSuccess": "Module status updated",
        "toggleError": "Failed to update module status",
        "validationError": "Please fix the validation errors"
      }
    }
  }
}
```

File: `frontend/i18n/locales/id-ID.json`

```json
{
  "nav": {
    "chatbot": {
      "title": "Chatbot",
      "modules": "Modul"
    }
  },
  "features": {
    "chatbot": {
      "page": {
        "title": "Konfigurasi Modul",
        "description": "Konfigurasi pengaturan modul chatbot di bawah.",
        "configTitle": "Pengaturan",
        "enabled": "Aktif",
        "disabled": "Nonaktif"
      },
      "form": {
        "save": "Simpan Konfigurasi",
        "saving": "Menyimpan..."
      },
      "messages": {
        "saveSuccess": "Konfigurasi berhasil disimpan",
        "saveError": "Gagal menyimpan konfigurasi",
        "toggleSuccess": "Status modul diperbarui",
        "toggleError": "Gagal memperbarui status modul",
        "validationError": "Harap perbaiki kesalahan validasi"
      }
    }
  }
}
```

## Files to Create

```
frontend/app/pages/platform/modules/
├── index.vue
└── [name].vue

frontend/app/components/custom/layout/
└── SidebarChatbotMenu.vue  (optional, if sidebar needs custom component)
```

## Files to Modify

- `frontend/app/composables/navigation/useNavigationItems.ts` - Add Chatbot menu with dynamic module children
- `frontend/app/components/custom/layout/LayoutSidebar.vue` - Integrate SidebarChatbotMenu (if using custom component)
- `frontend/i18n/locales/en-US.json` - Add chatbot translations
- `frontend/i18n/locales/id-ID.json` - Add chatbot translations

## Commands to Run

```bash
cd frontend

# Run lint
pnpm lint

# Run dev server
pnpm dev

# Build for production
pnpm build
```

## Validation Checklist

- [ ] Navigate to `/platform/modules/llm` shows LLM config
- [ ] Navigate to `/platform/modules/prompt` shows Prompt config
- [ ] Navigate to `/platform/modules/mcpServer` shows MCP Server config
- [ ] Navigate to `/platform/modules/widget` shows Widget config
- [ ] Invalid module name redirects to LLM
- [ ] Sidebar shows Chatbot menu with all 4 modules
- [ ] Sidebar shows enabled/disabled indicators
- [ ] Toggle enabled updates API and sidebar indicator
- [ ] Save persists changes to database
- [ ] Refresh loads saved data correctly
- [ ] Breadcrumb shows: Platform > Chatbot > {Module Name}
- [ ] Toast notifications appear on success/error
- [ ] Mobile responsive layout works

## Definition of Done

- [ ] Dynamic module page implemented
- [ ] Sidebar navigation integrated with auto-discovered modules
- [ ] Enabled/disabled indicators visible in sidebar
- [ ] Toggle updates API correctly
- [ ] Save button persists form changes
- [ ] Toast notifications working
- [ ] i18n translations added (en-US, id-ID)
- [ ] Breadcrumb navigation working
- [ ] No TypeScript errors
- [ ] No lint errors
- [ ] Responsive design verified

## Dependencies

- T45: Schema-driven form components must be complete
- Sidebar navigation system
- Breadcrumb composable

## Risk Factors

- **Medium Risk**: Sidebar integration may require adapting existing navigation structure
- **Low Risk**: Standard page patterns

## Notes

- All 4 modules appear in sidebar: LLM, Prompt, MCP Server, Widget
- Modules are auto-discovered from `MODULE_SCHEMAS` registry - no hardcoding
- The enabled indicator requires fetching chatbot config on sidebar mount
- Consider using Pinia store to cache chatbot config for sidebar efficiency
- Module titles in sidebar use static English from schema (not i18n translated per user requirement)
