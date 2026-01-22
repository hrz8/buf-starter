<script setup lang="ts">
import type { NavItem } from '@/composables/navigation/useNavigation';

import { ChevronRight, Plus } from 'lucide-vue-next';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from '@/components/ui/sidebar';
import { useSidebarNavigation } from '@/composables/navigation/useSidebarNavigation';

interface Props {
  items: NavItem[];
  label?: string;
}

const props = withDefaults(defineProps<Props>(), {
  label: 'Platform',
});

const emit = defineEmits<{
  action: [action: string];
}>();

const {
  isItemActive,
  hasActiveSubItem,
  toggleExpanded,
  isManuallyExpanded,
} = useSidebarNavigation();

const openStates = reactive<Record<string, boolean>>({});

// Initialize open states for items
function initializeOpenStates() {
  props.items.forEach((item) => {
    if (item.items?.length) {
      // Skip if already initialized (unless defaultExpanded)
      if (openStates[item.title] !== undefined && !item.defaultExpanded) {
        return;
      }

      // Items with defaultExpanded always start expanded
      if (item.defaultExpanded) {
        openStates[item.title] = true;
        return;
      }

      const shouldBeOpen = hasActiveSubItem(item) || isItemActive(item);
      const hasManualState = isManuallyExpanded(item.title);

      if (shouldBeOpen && !hasManualState) {
        openStates[item.title] = true;
      }
      else if (hasManualState) {
        openStates[item.title] = true;
      }
      else {
        openStates[item.title] = false;
      }
    }
  });
}

onMounted(() => {
  initializeOpenStates();
});

// Watch for items changes (e.g., when nodes load async)
watch(
  () => props.items,
  () => {
    initializeOpenStates();
  },
  { deep: true, immediate: true },
);

watch(() => useRoute().path, () => {
  props.items.forEach((item) => {
    if (item.items?.length) {
      const shouldBeOpen = hasActiveSubItem(item) || isItemActive(item);

      if (shouldBeOpen && !openStates[item.title]) {
        openStates[item.title] = true;
      }
    }
  });
});

function handleToggle(item: NavItem) {
  // Note: Don't manually toggle openStates here - CollapsibleTrigger handles
  // the toggle via v-model binding. We only track manual expansion state.
  toggleExpanded(item.title);
}

function handleAction(action: string, event: Event) {
  event.stopPropagation();
  emit('action', action);
}
</script>

<template>
  <SidebarGroup>
    <SidebarGroupLabel v-if="label">
      {{ label }}
    </SidebarGroupLabel>
    <SidebarMenu>
      <template
        v-for="item in items"
        :key="item.title"
      >
        <!-- Collapsible menu item (has sub-items) -->
        <Collapsible
          v-if="item.items?.length"
          v-model:open="openStates[item.title]"
          class="group/collapsible"
        >
          <SidebarMenuItem>
            <!--
              Note: Don't use tooltip on SidebarMenuButton when used with CollapsibleTrigger as-child
              The Tooltip wrapper breaks the as-child chain and prevents clicks from working
            -->
            <CollapsibleTrigger as-child>
              <SidebarMenuButton @click="handleToggle(item)">
                <component
                  :is="item.icon"
                  v-if="item.icon"
                />
                <span>{{ item.title }}</span>
                <!-- Action button (e.g., + for creating nodes) -->
                <Button
                  v-if="item.action"
                  variant="ghost"
                  size="icon"
                  class="ml-auto h-5 w-5 hover:bg-accent"
                  @click.stop="handleAction(item.action, $event)"
                >
                  <Plus class="h-3 w-3" />
                </Button>
                <ChevronRight
                  v-else
                  class="
                    ml-auto transition-transform duration-200
                    group-data-[state=open]/collapsible:rotate-90
                  "
                />
              </SidebarMenuButton>
            </CollapsibleTrigger>
            <CollapsibleContent>
              <SidebarMenuSub>
                <SidebarMenuSubItem
                  v-for="subItem in item.items"
                  :key="subItem.to"
                >
                  <SidebarMenuSubButton
                    as-child
                    :data-active="isItemActive(subItem)"
                  >
                    <NuxtLink
                      :to="subItem.to"
                      prefetch
                      class="flex items-center justify-between w-full min-w-0"
                    >
                      <span class="flex items-center gap-2 min-w-0 flex-1">
                        <component
                          :is="subItem.icon"
                          v-if="subItem.icon"
                          class="h-4 w-4 shrink-0"
                        />
                        <span class="truncate">{{ subItem.title }}</span>
                      </span>
                      <Badge
                        v-if="subItem.badge"
                        :variant="subItem.badgeVariant || 'secondary'"
                        class="ml-auto text-[10px] px-1.5 py-0"
                      >
                        {{ subItem.badge }}
                      </Badge>
                    </NuxtLink>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              </SidebarMenuSub>
            </CollapsibleContent>
          </SidebarMenuItem>
        </Collapsible>

        <!-- Non-collapsible menu item (no sub-items) -->
        <SidebarMenuItem v-else>
          <SidebarMenuButton
            as-child
            :tooltip="item.title"
            :data-active="isItemActive(item)"
          >
            <NuxtLink
              :to="item.to"
              prefetch
            >
              <component
                :is="item.icon"
                v-if="item.icon"
              />
              <span>{{ item.title }}</span>
            </NuxtLink>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </template>
    </SidebarMenu>
  </SidebarGroup>
</template>
