<script setup lang="ts">
import type { PropertySchema } from '@/lib/chatbot-modules';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';

defineProps<{
  name: string;
  schema: PropertySchema;
  disabled?: boolean;
}>();
</script>

<template>
  <FormField v-slot="{ componentField }" :name="name">
    <FormItem>
      <FormLabel>{{ schema.title }}</FormLabel>
      <FormControl>
        <Input
          v-bind="componentField"
          type="number"
          :min="schema.minimum"
          :max="schema.maximum"
          :step="schema.step || 1"
          :placeholder="schema.placeholder"
          :disabled="disabled"
        />
      </FormControl>
      <FormDescription v-if="schema.description">
        {{ schema.description }}
      </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
</template>
