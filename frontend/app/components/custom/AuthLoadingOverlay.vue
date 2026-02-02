<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';
import { useProjectStore } from '@/stores/project';

const authStore = useAuthStore();
const projectStore = useProjectStore();

/**
 * Randomized loading messages for a bit of fun
 */
const loadingMessages = [
  'Loading...',
  'Preparing your workspace...',
  'Brewing some coffee...',
  'Waking up the servers...',
  'Teaching hamsters to run faster...',
  'Convincing pixels to align...',
  'Summoning the data spirits...',
  'Polishing the dashboard...',
  'Herding digital cats...',
  'Assembling the bits and bytes...',
  'Warming up the flux capacitor...',
  'Consulting the cloud oracle...',
  'Untangling the network cables...',
  'Feeding the code monkeys...',
  'Reticulating splines...',
];

function getRandomMessage(): string {
  return loadingMessages[Math.floor(Math.random() * loadingMessages.length)] ?? loadingMessages[0]!;
}

const loadingMessage = ref(getRandomMessage());

/**
 * Show loading overlay when:
 * 1. Auth is not yet initialized (first page load)
 * 2. User is authenticated but projects haven't loaded yet
 *
 * This prevents user interaction with incomplete data state and
 * stops components from making API calls before auth is ready.
 */
const isLoading = computed(() => {
  // Show loading if auth hasn't been initialized yet
  if (!authStore.isInitialized)
    return true;

  // Only show project loading when authenticated
  if (!authStore.isAuthenticated)
    return false;

  // Show loading while projects are being fetched for the first time
  if (projectStore.pending && projectStore.projects.length === 0)
    return true;

  return false;
});

// Pick a new random message each time loading starts
watch(isLoading, (newVal, oldVal) => {
  if (newVal && !oldVal) {
    loadingMessage.value = getRandomMessage();
  }
});
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-200"
      leave-active-class="transition-opacity duration-200"
      enter-from-class="opacity-0"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isLoading"
        class="
          fixed inset-0 z-[100] flex items-center justify-center
          bg-background/80 backdrop-blur-sm
        "
      >
        <div class="flex flex-col items-center gap-4">
          <Icon
            name="lucide:loader-circle"
            class="size-8 animate-spin text-primary"
          />
          <p class="text-sm text-muted-foreground">
            {{ loadingMessage }}
          </p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
