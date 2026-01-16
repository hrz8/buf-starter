<script setup lang="ts">
import type { SidebarProps } from '@/components/ui/sidebar';

import {
  NavIAM,
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
import { useNavigationItems } from '@/composables/navigation/useNavigationItems';
import { useAuthStore } from '@/stores/auth';

const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: 'icon',
});

const authStore = useAuthStore();
const { mainNavItems, settingsNavItems, iamNavItems } = useNavigationItems();

const userData = computed(() => {
  const user = authStore.user;
  if (!user) {
    return {
      name: 'Guest',
      email: '',
      avatar: '',
    };
  }

  const displayName = user.name
    || [user.given_name, user.family_name].filter(Boolean).join(' ')
    || user.email
    || 'User';

  return {
    name: displayName,
    email: user.email || '',
    avatar: '', // TODO: Avatar could be added later if needed
  };
});

const data = {
  settings: settingsNavItems,
  iam: iamNavItems,
};
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <NavProject />
    </SidebarHeader>
    <SidebarContent>
      <NavMenu :items="mainNavItems" />
      <NavIAM :items="data.iam.value" />
      <NavSettings :settings="data.settings.value" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="userData" />
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>
</template>
