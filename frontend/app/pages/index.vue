<template>
  <div>Redirecting to dashboard...</div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: () => {
    return navigateTo('/dashboard');
  },
});
</script>
