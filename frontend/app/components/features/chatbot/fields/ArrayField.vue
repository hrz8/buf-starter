<script setup lang="ts">
import type { PropertySchema } from '@/lib/chatbot-modules';
import { Plus, X } from 'lucide-vue-next';
import { useFieldArray } from 'vee-validate';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

const props = defineProps<{
  name: string;
  schema: PropertySchema;
  disabled?: boolean;
}>();

// Use vee-validate field array for simple arrays
const { fields, push, remove } = useFieldArray<string>(props.name);

function addItem() {
  push('');
}

function removeItem(index: number) {
  remove(index);
}
</script>

<template>
  <div class="space-y-2">
    <div>
      <Label class="text-sm font-medium">{{ schema.title }}</Label>
      <p v-if="schema.description" class="text-sm text-muted-foreground">
        {{ schema.description }}
      </p>
    </div>

    <div class="space-y-2">
      <!-- Existing items -->
      <div
        v-for="(field, index) in fields"
        :key="field.key"
        class="flex items-center gap-2"
      >
        <Input
          v-model="field.value"
          :placeholder="schema.items?.placeholder || 'Enter value'"
          :disabled="disabled"
          class="flex-1"
        />
        <Button
          type="button"
          variant="ghost"
          size="icon"
          :disabled="disabled"
          @click="removeItem(index)"
        >
          <X class="h-4 w-4" />
        </Button>
      </div>

      <!-- Add button -->
      <Button
        type="button"
        variant="outline"
        size="sm"
        :disabled="disabled"
        @click="addItem"
      >
        <Plus class="h-4 w-4 mr-2" />
        Add Item
      </Button>
    </div>
  </div>
</template>
