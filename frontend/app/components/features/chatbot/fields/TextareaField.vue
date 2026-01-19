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
import { Textarea } from '@/components/ui/textarea';

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
        <Textarea
          v-bind="componentField"
          :placeholder="schema.placeholder"
          :maxlength="schema.maxLength"
          :disabled="disabled"
          class="min-h-[100px]"
        />
      </FormControl>
      <FormDescription v-if="schema.description">
        {{ schema.description }}
      </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
</template>
