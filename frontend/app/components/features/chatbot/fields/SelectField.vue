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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

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
        <Select v-bind="componentField" :disabled="disabled">
          <SelectTrigger>
            <SelectValue :placeholder="schema.placeholder || 'Select an option'" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="option in schema.enum" :key="option" :value="option">
              {{ schema.enumLabels?.[option] || option }}
            </SelectItem>
          </SelectContent>
        </Select>
      </FormControl>
      <FormDescription v-if="schema.description">
        {{ schema.description }}
      </FormDescription>
      <FormMessage />
    </FormItem>
  </FormField>
</template>
