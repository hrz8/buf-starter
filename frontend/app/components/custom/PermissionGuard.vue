<script setup lang="ts">
import { usePermissions } from '~/composables/usePermissions';

const props = defineProps<{
  /**
   * Single permission to check
   */
  permission?: string;
  /**
   * Multiple permissions to check
   */
  permissions?: string[];
  /**
   * Mode for multiple permissions: 'any' (OR) or 'all' (AND)
   */
  mode?: 'any' | 'all';
  /**
   * Fallback behavior: 'hide' (default), 'disabled', 'redirect'
   */
  fallback?: 'hide' | 'disabled' | 'redirect';
  /**
   * Optional project ID for project-scoped permission check
   */
  projectId?: string;
}>();

const { can, canAny, canAll, isMemberOf } = usePermissions();

const hasAccess = computed(() => {
  // If projectId provided, also check membership
  if (props.projectId && !isMemberOf(props.projectId))
    return false;

  // Single permission check
  if (props.permission)
    return can(props.permission);

  // Multiple permissions check
  if (props.permissions && props.permissions.length > 0) {
    return props.mode === 'all'
      ? canAll(props.permissions)
      : canAny(props.permissions);
  }

  // No permission specified = allow
  return true;
});

// Handle redirect fallback
const router = useRouter();
watch(hasAccess, (value) => {
  if (!value && props.fallback === 'redirect')
    router.push('/access-denied');
}, { immediate: true });
</script>

<template>
  <!-- Render slot content if has access -->
  <slot v-if="hasAccess" />

  <!-- Render disabled slot if no access and fallback is 'disabled' -->
  <slot v-else-if="fallback === 'disabled'" name="disabled">
    <!-- Default disabled state -->
    <div class="opacity-50 pointer-events-none">
      <slot />
    </div>
  </slot>

  <!-- Render fallback slot for 'hide' mode (default shows nothing) -->
  <slot v-else-if="fallback !== 'redirect'" name="fallback" />
</template>
