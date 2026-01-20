<script setup lang="ts">
import type { ChatbotNode } from '~~/gen/chatbot/nodes/v1/node_pb';
import type { SidebarProps } from '@/components/ui/sidebar';

import {
  NavIAM,
  NavMenu,
  NavProject,
  NavSettings,
  NavUser,
} from '@/components/custom/nav';
import NodeCreateSheet from '@/components/features/chatbot-node/NodeCreateSheet.vue';
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

// Node create sheet state
const isNodeCreateSheetOpen = ref(false);
const router = useRouter();

function handleNavAction(action: string) {
  if (action === 'createNode') {
    isNodeCreateSheetOpen.value = true;
  }
}

function handleNodeCreated(node: ChatbotNode) {
  // Navigate to the newly created node
  router.push(`/platform/node/${node.id}`);
}
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <NavProject />
    </SidebarHeader>
    <SidebarContent>
      <NavMenu :items="mainNavItems" @action="handleNavAction" />
      <NavIAM :items="data.iam.value" />
      <NavSettings :settings="data.settings.value" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="userData" />
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>

  <!-- Node Create Sheet -->
  <NodeCreateSheet
    v-model:open="isNodeCreateSheetOpen"
    @success="handleNodeCreated"
  />
</template>
