<script setup lang="ts">
import type { NavItem } from '@/composables/navigation/useNavigation';

import { ChevronRight } from 'lucide-vue-next';

import { Badge } from '@/components/ui/badge';
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

const {
  isItemActive,
  hasActiveSubItem,
  toggleExpanded,
  isManuallyExpanded,
} = useSidebarNavigation();

const openStates = reactive<Record<string, boolean>>({});

onMounted(() => {
  props.items.forEach((item) => {
    if (item.items?.length) {
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
});

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
  openStates[item.title] = !openStates[item.title];
  toggleExpanded(item.title);
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
            <CollapsibleTrigger
              as-child
              @click.prevent="handleToggle(item)"
            >
              <SidebarMenuButton :tooltip="item.title">
                <component
                  :is="item.icon"
                  v-if="item.icon"
                />
                <span>{{ item.title }}</span>
                <ChevronRight
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
                  :key="subItem.title"
                >
                  <SidebarMenuSubButton
                    as-child
                    :data-active="isItemActive(subItem)"
                  >
                    <NuxtLink
                      :to="subItem.to"
                      prefetch
                      class="flex items-center justify-between w-full"
                    >
                      <span class="flex items-center gap-2">
                        <component
                          :is="subItem.icon"
                          v-if="subItem.icon"
                          class="h-4 w-4"
                        />
                        <span>{{ subItem.title }}</span>
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
