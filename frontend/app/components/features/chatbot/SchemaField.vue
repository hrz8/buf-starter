<script setup lang="ts">
import type { PropertySchema } from '@/lib/chatbot-modules';
import { resolveFieldComponent } from './fields';

const props = defineProps<{
  name: string;
  schema: PropertySchema;
  disabled?: boolean;
}>();

const FieldComponent = computed(() => resolveFieldComponent(props.schema));
</script>

<template>
  <component
    :is="FieldComponent"
    :name="name"
    :schema="schema"
    :disabled="disabled"
  />
</template>
