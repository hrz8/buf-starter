<script setup lang="ts">
import type { PropertySchema } from '@/lib/chatbot-modules';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
} from '@/components/ui/form';
import { Switch } from '@/components/ui/switch';

defineProps<{
  name: string;
  schema: PropertySchema;
  disabled?: boolean;
}>();
</script>

<template>
  <FormField v-slot="{ componentField, value }" :name="name">
    <FormItem class="flex items-center justify-between gap-4">
      <div class="space-y-0.5">
        <FormLabel>{{ schema.title }}</FormLabel>
        <FormDescription v-if="schema.description">
          {{ schema.description }}
        </FormDescription>
      </div>
      <FormControl>
        <Switch
          :model-value="value"
          :disabled="disabled"
          @update:model-value="componentField['onUpdate:modelValue']"
        />
      </FormControl>
    </FormItem>
  </FormField>
</template>
