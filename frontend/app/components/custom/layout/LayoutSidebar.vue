<script setup lang="ts">
import type { SidebarProps } from '@/components/ui/sidebar';

import type { NavItem } from '@/composables/navigation/useNavigation';
import {
  Key,
  LucideHome,
  Puzzle,
  Smartphone,
} from 'lucide-vue-next';

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
