<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core';

import {
  SidebarProvider,
  SidebarInset,
} from '@/components/ui/sidebar';
import AppSidebar from '@/components/AppSidebar.vue';
import AppHeader from '@/components/AppHeader.vue';

const sidebarOpen = useLocalStorage('sidebar-state', true);

function handleOpenUpdate(open: boolean) {
  sidebarOpen.value = open;
}
</script>

<template>
  <SidebarProvider
    :default-open="sidebarOpen"
    @update:open="handleOpenUpdate"
  >
    <AppSidebar />
    <SidebarInset>
      <AppHeader />
      <main class="flex-1">
        <slot />
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
