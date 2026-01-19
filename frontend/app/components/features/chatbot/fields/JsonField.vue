<script setup lang="ts">
import type { PropertySchema } from '@/lib/chatbot-modules';
import { useField } from 'vee-validate';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

const props = defineProps<{
  name: string;
  schema: PropertySchema;
  disabled?: boolean;
}>();

// Use vee-validate field
const { value, errorMessage, setValue } = useField<string>(() => props.name);

// Local state for JSON validation
const jsonError = ref<string | null>(null);

// Validate JSON and update field
function handleInput(event: Event) {
  const target = event.target as HTMLTextAreaElement;
  const newValue = target.value;

  // Update the field value
  setValue(newValue);

  // Validate JSON
  if (!newValue || newValue.trim() === '') {
    jsonError.value = null;
    return;
  }

  try {
    JSON.parse(newValue);
    jsonError.value = null;
  }
  catch (e) {
    jsonError.value = e instanceof Error ? e.message : 'Invalid JSON';
  }
}

// Format/pretty-print JSON
function formatJson() {
  if (!value.value || value.value.trim() === '') {
    return;
  }

  try {
    const parsed = JSON.parse(value.value);
    const formatted = JSON.stringify(parsed, null, 2);
    setValue(formatted);
    jsonError.value = null;
  }
  catch (e) {
    jsonError.value = e instanceof Error ? e.message : 'Invalid JSON';
  }
}

// Minify JSON
function minifyJson() {
  if (!value.value || value.value.trim() === '') {
    return;
  }

  try {
    const parsed = JSON.parse(value.value);
    const minified = JSON.stringify(parsed);
    setValue(minified);
    jsonError.value = null;
  }
  catch (e) {
    jsonError.value = e instanceof Error ? e.message : 'Invalid JSON';
  }
}

// Combined error message
const displayError = computed(() => jsonError.value || errorMessage.value);
</script>

<template>
  <div class="space-y-2">
    <div class="flex items-center justify-between">
      <div>
        <Label class="text-sm font-medium">{{ schema.title }}</Label>
        <p v-if="schema.description" class="text-sm text-muted-foreground">
          {{ schema.description }}
        </p>
      </div>
      <div class="flex gap-1">
        <Button
          type="button"
          variant="ghost"
          size="sm"
          :disabled="disabled || !value"
          @click="formatJson"
        >
          Format
        </Button>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          :disabled="disabled || !value"
          @click="minifyJson"
        >
          Minify
        </Button>
      </div>
    </div>

    <Textarea
      :model-value="value || ''"
      :placeholder="schema.placeholder || '{\n  \n}'"
      :disabled="disabled"
      class="font-mono text-sm min-h-[150px]"
      @input="handleInput"
    />

    <p v-if="displayError" class="text-sm text-destructive">
      {{ displayError }}
    </p>
  </div>
</template>
