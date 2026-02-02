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
  items: IAMItem[];
}>();

const { t } = useI18n();

interface IAMItem {
  name: string;
  url: string;
  icon: LucideIcon;
}

const { isItemActive } = useSidebarNavigation();

function toNavItem(item: IAMItem) {
  return {
    title: item.name,
    to: item.url.startsWith('/') ? item.url : `/${item.url}`,
    icon: item.icon,
  };
}
</script>

<template>
  <SidebarGroup v-if="items.length > 0" class="group-data-[collapsible=icon]:hidden">
    <SidebarGroupLabel>{{ t('nav.iam.title') }}</SidebarGroupLabel>
    <SidebarMenu>
      <SidebarMenuItem
        v-for="item in items"
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
