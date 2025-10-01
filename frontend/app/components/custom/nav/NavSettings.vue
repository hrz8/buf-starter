<script setup lang="ts">
import type { LucideIcon } from 'lucide-vue-next';

import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar';
import { useSidebarNavigation } from '@/composables/navigation/useSidebarNavigation';

defineProps<{
  settings: SettingItem[];
}>();

const { t } = useI18n();

interface SettingItem {
  name: string;
  url: string;
  icon: LucideIcon;
}

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
    <SidebarGroupLabel>{{ t('nav.settings.title') }}</SidebarGroupLabel>
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
