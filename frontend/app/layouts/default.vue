<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core';
import { LayoutHeader, LayoutSidebar } from '@/components/custom/layout';
import EmailVerificationOverlay from '@/components/features/email-verification/EmailVerificationOverlay.vue';

import {
  SidebarInset,
  SidebarProvider,
} from '@/components/ui/sidebar';
import { Toaster } from '@/components/ui/sonner';
import { useProjectService } from '@/composables/services/useProjectService';
import { useChatbotStore } from '@/stores/chatbot';
import { useChatbotNodeStore } from '@/stores/chatbot-node';
import { useProjectStore } from '@/stores/project';

import 'vue-sonner/style.css';

const route = useRoute();
const sidebarOpen = useLocalStorage('sidebar-state', true);

function handleOpenUpdate(open: boolean) {
  sidebarOpen.value = open;
}

const projectStore = useProjectStore();
const chatbotStore = useChatbotStore();
const nodeStore = useChatbotNodeStore();
const { query } = useProjectService();

async function fetchProjects() {
  projectStore.setLoading(true);
  projectStore.setError(null);

  try {
    const response = await query({
      query: {
        pagination: {
          page: 1,
          pageSize: 50,
        },
      },
    });
    projectStore.setProjects(response?.data ?? []);
  }
  catch (err) {
    projectStore.setError(err as Error);
  }
  finally {
    projectStore.setLoading(false);
  }
}

onMounted(() => {
  if (!projectStore.projects.length) {
    fetchProjects();
  }
});

// Watch for project changes and load chatbot config + nodes
watch(
  () => projectStore.activeProjectId,
  async (projectId) => {
    if (projectId) {
      // Load in parallel for better performance
      await Promise.all([
        chatbotStore.ensureLoaded(),
        nodeStore.ensureLoaded(),
      ]);
    }
  },
  { immediate: true },
);
</script>

<template>
  <SidebarProvider
    :default-open="sidebarOpen"
    @update:open="handleOpenUpdate"
  >
    <Toaster />
    <EmailVerificationOverlay />
    <LayoutSidebar />
    <SidebarInset>
      <LayoutHeader />
      <main
        :key="JSON.stringify(route.query)"
        class="flex-1"
      >
        <slot />
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
