<script setup lang="ts">
import type { Component } from 'vue';
import type { PropertySchema } from '@/lib/chatbot-modules';
import { ChevronDown, Plus, Trash2 } from 'lucide-vue-next';
import { useFieldArray } from 'vee-validate';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible';
import { Label } from '@/components/ui/label';

const props = defineProps<{
  name: string;
  schema: PropertySchema;
  disabled?: boolean;
}>();

// Use vee-validate field array
const { fields, push, remove } = useFieldArray<Record<string, unknown>>(props.name);

// Track which items are expanded
const expandedItems = ref<Set<number>>(new Set([0])); // First item expanded by default

function toggleItem(index: number) {
  if (expandedItems.value.has(index)) {
    expandedItems.value.delete(index);
  }
  else {
    expandedItems.value.add(index);
  }
  // Trigger reactivity
  expandedItems.value = new Set(expandedItems.value);
}

function isExpanded(index: number): boolean {
  return expandedItems.value.has(index);
}

// Get item title for collapsible header
function getItemTitle(item: { value: Record<string, unknown> }, index: number): string {
  const titleKey = props.schema.items?.titleKey;
  if (titleKey && item.value && item.value[titleKey]) {
    return String(item.value[titleKey]);
  }
  // Fallback: try to find first string value
  if (item.value) {
    for (const [, val] of Object.entries(item.value)) {
      if (typeof val === 'string' && val.trim()) {
        return val;
      }
    }
  }
  return `${props.schema.items?.title || 'Item'} ${index + 1}`;
}

// Add new item with default values
function addItem() {
  const newItem: Record<string, unknown> = {};
  const itemProps = props.schema.items?.properties || {};

  for (const [key, propSchema] of Object.entries(itemProps)) {
    if (propSchema.default !== undefined) {
      newItem[key] = propSchema.default;
    }
    else if (propSchema.type === 'boolean') {
      newItem[key] = false;
    }
    else if (propSchema.type === 'string') {
      newItem[key] = '';
    }
    else if (propSchema.type === 'number') {
      newItem[key] = 0;
    }
    else if (propSchema.type === 'array') {
      newItem[key] = [];
    }
    else if (propSchema.type === 'object') {
      newItem[key] = {};
    }
  }

  push(newItem);

  // Expand the newly added item
  const newIndex = fields.value.length;
  expandedItems.value.add(newIndex);
  expandedItems.value = new Set(expandedItems.value);
}

function removeItem(index: number) {
  remove(index);
  expandedItems.value.delete(index);
  // Re-index expanded items
  const newExpanded = new Set<number>();
  expandedItems.value.forEach((i) => {
    if (i < index) {
      newExpanded.add(i);
    }
    else if (i > index) {
      newExpanded.add(i - 1);
    }
  });
  expandedItems.value = newExpanded;
}

// Get properties to render for each item
const itemProperties = computed(() => {
  return props.schema.items?.properties || {};
});

// Lazy import to avoid circular dependency - resolved at runtime
const SchemaField = defineAsyncComponent(() => import('../SchemaField.vue'));

// Resolve field component based on schema type
function getFieldComponent(propSchema: PropertySchema): Component {
  // JSON field
  if (propSchema.additionalTypeInfo === 'json') {
    return defineAsyncComponent(() => import('./JsonField.vue'));
  }

  // Array of objects (recursive)
  if (propSchema.type === 'array' && propSchema.items?.type === 'object') {
    return defineAsyncComponent(() => import('./ObjectArrayField.vue'));
  }

  // Simple array
  if (propSchema.type === 'array') {
    return defineAsyncComponent(() => import('./ArrayField.vue'));
  }

  // Nested object - render its properties inline
  if (propSchema.type === 'object' && propSchema.properties) {
    // Return a marker - we'll handle this specially in the template
    return markRaw({ __isNestedObject: true }) as unknown as Component;
  }

  // Default to SchemaField for all other types
  return SchemaField;
}

// Check if component is nested object marker
function isNestedObjectMarker(component: Component): boolean {
  return !!(component as { __isNestedObject?: boolean }).__isNestedObject;
}
</script>

<template>
  <div class="space-y-3">
    <div class="flex items-center justify-between">
      <div>
        <Label class="text-sm font-medium">{{ schema.title }}</Label>
        <p v-if="schema.description" class="text-sm text-muted-foreground">
          {{ schema.description }}
        </p>
      </div>
      <Button
        type="button"
        variant="outline"
        size="sm"
        :disabled="disabled"
        @click="addItem"
      >
        <Plus class="h-4 w-4 mr-1" />
        Add {{ schema.items?.title || 'Item' }}
      </Button>
    </div>

    <!-- Empty state -->
    <div
      v-if="fields.length === 0"
      class="border-2 border-dashed border-muted rounded-lg p-6 text-center"
    >
      <p class="text-sm text-muted-foreground">
        No items yet. Click "Add {{ schema.items?.title || 'Item' }}" to create one.
      </p>
    </div>

    <!-- Items -->
    <div class="space-y-3">
      <Card v-for="(field, index) in fields" :key="field.key" class="overflow-hidden">
        <Collapsible :open="isExpanded(index)">
          <CardHeader class="p-0">
            <CollapsibleTrigger
              class="
              flex w-full items-center justify-between p-4
              hover:bg-muted/50 transition-colors
            "
              @click="toggleItem(index)"
            >
              <div class="flex items-center gap-2">
                <ChevronDown
                  class="h-4 w-4 transition-transform duration-200"
                  :class="{ '-rotate-90': !isExpanded(index) }"
                />
                <span class="font-medium text-sm">
                  {{ getItemTitle(field, index) }}
                </span>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                :disabled="disabled"
                @click.stop="removeItem(index)"
              >
                <Trash2 class="h-4 w-4 text-destructive" />
              </Button>
            </CollapsibleTrigger>
          </CardHeader>

          <CollapsibleContent>
            <CardContent class="pt-0 pb-4 px-4 space-y-4">
              <template
                v-for="[propKey, propSchema] in Object.entries(itemProperties)"
                :key="propKey"
              >
                <!-- Nested Object - render properties inline with indentation -->
                <template v-if="isNestedObjectMarker(getFieldComponent(propSchema))">
                  <div class="space-y-3">
                    <div>
                      <Label class="text-sm font-medium">{{ propSchema.title }}</Label>
                      <p v-if="propSchema.description" class="text-sm text-muted-foreground">
                        {{ propSchema.description }}
                      </p>
                    </div>
                    <div class="pl-4 border-l-2 border-muted space-y-4">
                      <template
                        v-for="
                          [nestedKey, nestedSchema]
                            in Object.entries(propSchema.properties || {})"
                        :key="nestedKey"
                      >
                        <component
                          :is="getFieldComponent(nestedSchema)"
                          :name="`${name}[${index}].${propKey}.${nestedKey}`"
                          :schema="nestedSchema"
                          :disabled="disabled"
                        />
                      </template>
                    </div>
                  </div>
                </template>

                <!-- All other field types -->
                <component
                  :is="getFieldComponent(propSchema)"
                  v-else
                  :name="`${name}[${index}].${propKey}`"
                  :schema="propSchema"
                  :disabled="disabled"
                />
              </template>
            </CardContent>
          </CollapsibleContent>
        </Collapsible>
      </Card>
    </div>
  </div>
</template>
