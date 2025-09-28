<script setup lang="ts">
definePageMeta({
  middleware: () => {
    return navigateTo('/dashboard');
  },
});
</script>

<template>
  <div>Redirecting to dashboard...</div>
</template>
