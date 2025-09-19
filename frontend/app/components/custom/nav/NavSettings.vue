<script setup lang="ts">
import type { LucideIcon } from 'lucide-vue-next';

import {
  SidebarGroupLabel,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarGroup,
  SidebarMenu,
} from '@/components/ui/sidebar';
import { useSidebarNavigation } from '@/composables/navigation/useSidebarNavigation';

interface SettingItem {
  name: string;
  url: string;
  icon: LucideIcon;
}

defineProps<{
  settings: SettingItem[];
}>();

const { isItemActive } = useSidebarNavigation();

function toNavItem(item: SettingItem) {
  return {
    title: item.name,
    to: item.url.startsWith('/') ? item.url : `/${item.url}`,
    icon: item.icon,
  };
}
</script>

<template>
  <SidebarGroup class="group-data-[collapsible=icon]:hidden">
    <SidebarGroupLabel>Settings</SidebarGroupLabel>
    <SidebarMenu>
      <SidebarMenuItem
        v-for="item in settings"
        :key="item.name"
      >
        <SidebarMenuButton
          as-child
          :data-active="isItemActive(toNavItem(item))"
        >
          <NuxtLink
            :to="item.url"
            prefetch
          >
            <component :is="item.icon" />
            <span>{{ item.name }}</span>
          </NuxtLink>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  </SidebarGroup>
</template>
