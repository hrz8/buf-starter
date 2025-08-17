<script setup lang="ts">
import ExampleLanguageSelector from '@/components/example/ExampleLanguageSelector.vue';
import ExampleTableName from '@/components/example/ExampleTableName.vue';
import ExampleSayHello from '@/components/example/ExampleSayHello.vue';

definePageMeta({
  layout: 'default',
});
</script>

<template>
  <div class="max-w-lg mx-auto mt-10 p-6 bg-white shadow-md rounded-md space-y-8">
    <div class="flex justify-end">
      <ExampleLanguageSelector />
    </div>
    <ExampleSayHello />
    <ExampleTableName />
  </div>
</template>
