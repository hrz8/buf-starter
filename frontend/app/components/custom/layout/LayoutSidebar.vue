<script setup lang="ts">
import type { SidebarProps } from '@/components/ui/sidebar';

import {
  NavMenu,
  NavProject,
  NavSettings,
  NavUser,
} from '@/components/custom/nav';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '@/components/ui/sidebar';
import { mainNavItems, settingsNavItems } from '@/config/navigation';

const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: 'icon',
});

const data = {
  user: {
    name: 'shadcn',
    email: 'm@example.com',
    avatar: '/avatars/shadcn.jpg',
  },
  settings: settingsNavItems,
};
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <NavProject />
    </SidebarHeader>
    <SidebarContent>
      <NavMenu :items="mainNavItems" />
      <NavSettings :settings="data.settings" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="data.user" />
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>
</template>
