<script setup lang="ts">
import {
  GalleryVerticalEnd,
  AudioWaveform,
  Smartphone,
  LucideHome,
  Command,
  Puzzle,
  Key,
} from 'lucide-vue-next';

import type { NavItem } from '~/composables/navigation/useNavigation';
import type { SidebarProps } from '@/components/ui/sidebar';

import {
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
  Sidebar,
} from '@/components/ui/sidebar';
import TeamSwitcher from '@/components/TeamSwitcher.vue';
import NavSettings from '@/components/NavSettings.vue';
import NavMain from '@/components/NavMain.vue';
import NavUser from '@/components/NavUser.vue';

const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: 'icon',
});

const data = {
  user: {
    name: 'shadcn',
    email: 'm@example.com',
    avatar: '/avatars/shadcn.jpg',
  },
  teams: [
    {
      name: 'Acme Inc',
      logo: GalleryVerticalEnd,
      plan: 'Enterprise',
    },
    {
      name: 'Acme Corp.',
      logo: AudioWaveform,
      plan: 'Startup',
    },
    {
      name: 'Evil Corp.',
      logo: Command,
      plan: 'Free',
    },
  ],
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
        title: 'Datatable 1',
        to: '/examples/datatable/datatable1',
      },
      {
        title: 'Datatable 2',
        to: '/examples/datatable/datatable2',
      },
      {
        title: 'Datatable 3',
        to: '/examples/datatable/datatable3',
      },
      {
        title: 'Datatable 4',
        to: '/examples/datatable/datatable4',
      },
      {
        title: 'Datatable 5',
        to: '/examples/datatable/datatable5',
      },
      {
        title: 'Datatable 6',
        to: '/examples/datatable/datatable6',
      },
      {
        title: 'Datatable 7',
        to: '/examples/datatable/datatable7',
      },
      {
        title: 'Datatable 8',
        to: '/examples/datatable/datatable8',
      },
      {
        title: 'Datatable 9',
        to: '/examples/datatable/datatable9',
      },
      {
        title: 'Datatable 10',
        to: '/examples/datatable/datatable10',
      },
      {
        title: 'Datatable 11',
        to: '/examples/datatable/datatable11',
      },
      {
        title: 'Datatable 12',
        to: '/examples/datatable/datatable12',
      },
      {
        title: 'Datatable 13',
        to: '/examples/datatable/datatable13',
      },
      {
        title: 'Datatable 14',
        to: '/examples/datatable/datatable14',
      },
      {
        title: 'Datatable 15',
        to: '/examples/datatable/datatable15',
      },
      {
        title: 'Datatable 16',
        to: '/examples/datatable/datatable16',
      },
    ],
  },
];
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <TeamSwitcher :teams="data.teams" />
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="mainNavItems" />
      <NavSettings :settings="data.settings" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="data.user" />
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>
</template>
