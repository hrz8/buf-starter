<script setup lang="ts">
import type { ModuleSchema } from '@/lib/chatbot-modules';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';

defineProps<{
  schema: ModuleSchema;
  enabled: boolean;
  disabled?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:enabled', value: boolean): void;
}>();
</script>

<template>
  <div class="flex items-center justify-between p-4 bg-muted/50 rounded-lg">
    <div class="flex items-center gap-3">
      <div>
        <div class="flex items-center gap-2">
          <h3 class="font-semibold">
            {{ schema.title }}
          </h3>
          <Badge :variant="enabled ? 'default' : 'secondary'">
            {{ enabled ? 'Enabled' : 'Disabled' }}
          </Badge>
        </div>
        <p class="text-sm text-muted-foreground">
          {{ schema.description }}
        </p>
      </div>
    </div>
    <Switch
      :model-value="enabled"
      :disabled="disabled"
      @update:model-value="emit('update:enabled', $event)"
    />
  </div>
</template>
