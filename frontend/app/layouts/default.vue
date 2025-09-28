<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core';

import { LayoutHeader, LayoutSidebar } from '@/components/custom/layout';
import {
  SidebarInset,
  SidebarProvider,
} from '@/components/ui/sidebar';
import { useProjectService } from '@/composables/services/useProjectService';
import { useProjectStore } from '@/stores/project';

const route = useRoute();
const sidebarOpen = useLocalStorage('sidebar-state', true);

function handleOpenUpdate(open: boolean) {
  sidebarOpen.value = open;
}

const projectStore = useProjectStore();
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
</script>

<template>
  <SidebarProvider
    :default-open="sidebarOpen"
    @update:open="handleOpenUpdate"
  >
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
