<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';

definePageMeta({
  layout: 'default',
});

const { t } = useI18n();
const authStore = useAuthStore();

// Get first name for greeting, fallback to name or email
const greeting = computed(() => {
  const user = authStore.user;
  if (!user) {
    return '';
  }

  const firstName = user.given_name
    || user.name?.split(' ')[0]
    || user.email?.split('@')[0]
    || 'User';

  return firstName;
});
</script>

<template>
  <div class="flex flex-1 flex-col gap-4 p-4 pt-0">
    <!-- Greeting Section -->
    <div class="space-y-1">
      <h1 class="text-2xl font-semibold tracking-tight">
        {{ t('dashboard.greeting', { name: greeting }) }}
      </h1>
      <p class="text-sm text-muted-foreground">
        {{ t('dashboard.subtitle') }}
      </p>
    </div>

    <!-- Dashboard Content -->
    <div class="grid auto-rows-min gap-4 md:grid-cols-3">
      <div class="bg-muted/50 aspect-video rounded-xl" />
      <div class="bg-muted/50 aspect-video rounded-xl" />
      <div class="bg-muted/50 aspect-video rounded-xl" />
    </div>
    <div class="bg-muted/50 min-h-[100vh] flex-1 rounded-xl md:min-h-min" />
  </div>
</template>
