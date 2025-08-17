<script setup lang="ts">
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar';
import AppSidebar from '@/components/AppSidebar.vue';
import AppHeader from '@/components/AppHeader.vue';

</script>

<template>
  <SidebarProvider>
    <AppSidebar />
    <SidebarInset>
      <AppHeader />
      <main class="flex-1">
        <slot />
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
