<script setup lang="ts">
import type { ModuleSchema, PropertySchema } from '@/lib/chatbot-modules';
import { isNestedObject, resolveFieldComponent } from './fields';
import SchemaField from './SchemaField.vue';

defineProps<{
  schema: ModuleSchema;
  disabled?: boolean;
}>();

// Get field entries excluding 'enabled' which is handled separately
function getFieldEntries(properties: Record<string, PropertySchema>): [string, PropertySchema][] {
  return Object.entries(properties).filter(([key]) => key !== 'enabled');
}
</script>

<template>
  <div class="space-y-6">
    <!-- Iterate over schema properties (excluding 'enabled') -->
    <template v-for="[key, prop] in getFieldEntries(schema.properties)" :key="key">
      <!-- Nested object - render properties inline with indentation -->
      <div v-if="isNestedObject(prop)" class="space-y-4">
        <div>
          <h4 class="text-sm font-medium">
            {{ prop.title }}
          </h4>
          <p v-if="prop.description" class="text-sm text-muted-foreground">
            {{ prop.description }}
          </p>
        </div>
        <div class="pl-4 border-l-2 border-muted space-y-4">
          <template
            v-for="[nestedKey, nestedProp] in Object.entries(prop.properties || {})"
            :key="nestedKey"
          >
            <!-- Nested fields use SchemaField which auto-resolves type -->
            <SchemaField
              :name="`${key}.${nestedKey}`"
              :schema="nestedProp"
              :disabled="disabled"
            />
          </template>
        </div>
      </div>

      <!-- All other field types (including arrays, json, etc.) -->
      <component
        :is="resolveFieldComponent(prop)"
        v-else
        :name="key"
        :schema="prop"
        :disabled="disabled"
      />
    </template>
  </div>
</template>
