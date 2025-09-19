<script setup lang="ts">
import {
  Smartphone,
  LucideHome,
  Puzzle,
  Key,
} from 'lucide-vue-next';

import type { NavItem } from '@/composables/navigation/useNavigation';
import type { SidebarProps } from '@/components/ui/sidebar';

import {
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
  Sidebar,
} from '@/components/ui/sidebar';
import {
  NavSettings,
  NavProject,
  NavMenu,
  NavUser,
} from '@/components/custom/nav';

const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: 'icon',
});

const data = {
  user: {
    name: 'shadcn',
    email: 'm@example.com',
    avatar: '/avatars/shadcn.jpg',
  },
  settings: [
    {
      name: 'Api Keys',
      url: '/settings/api-keys',
      icon: Key,
    },
  ],
};

const mainNavItems: NavItem[] = [
  {
    title: 'Dashboard',
    to: '/dashboard',
    icon: LucideHome,
  },
  {
    title: 'Devices',
    to: '/devices',
    match: '/devices',
    icon: Smartphone,
    items: [
      {
        title: 'Scan',
        to: '/devices/scan',
      },
      {
        title: 'Chat',
        to: '/devices/chat',
      },
    ],
  },
  {
    title: 'Examples',
    to: '/examples',
    match: '/examples',
    icon: Puzzle,
    items: [
      {
        title: 'Simple Table',
        to: '/examples/simple-table',
      },
      {
        title: 'Datatable',
        to: '/examples/datatable/datatable18',
      },
    ],
  },
];
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
