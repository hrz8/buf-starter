<script setup lang="ts">
import { ChevronLeft, ChevronRight, Info, Plus } from 'lucide-vue-next';
import { ref } from 'vue';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import InlinePermissionCreateDialog from './InlinePermissionCreateDialog.vue';

interface TransferListItem {
  id: string;
  [key: string]: any;
}

interface Props {
  availableItems: TransferListItem[];
  assignedItems: TransferListItem[];
  label: string; // e.g., "Roles", "Permissions"
  singularLabel: string; // e.g., "Role", "Permission"
  labelKey?: string; // Field to use for display (default: "name")
  isLoading?: boolean;
  allowInlineCreate?: boolean; // Only true for permissions
  showTooltip?: boolean; // Show description tooltip
}

const props = withDefaults(defineProps<Props>(), {
  labelKey: 'name',
  isLoading: false,
  allowInlineCreate: false,
  showTooltip: false,
});

const emit = defineEmits<{
  assign: [ids: string[]];
  remove: [ids: string[]];
}>();

const availableSelected = ref<string[]>([]);
const assignedSelected = ref<string[]>([]);
const inlineCreateDialogOpen = ref(false);

function toggleAvailableSelection(id: string) {
  const index = availableSelected.value.indexOf(id);
  if (index > -1) {
    availableSelected.value.splice(index, 1);
  }
  else {
    availableSelected.value.push(id);
  }
}

function toggleAssignedSelection(id: string) {
  const index = assignedSelected.value.indexOf(id);
  if (index > -1) {
    assignedSelected.value.splice(index, 1);
  }
  else {
    assignedSelected.value.push(id);
  }
}

function assignSelected() {
  if (availableSelected.value.length > 0) {
    emit('assign', [...availableSelected.value]);
    availableSelected.value = [];
  }
}

function removeSelected() {
  if (assignedSelected.value.length > 0) {
    emit('remove', [...assignedSelected.value]);
    assignedSelected.value = [];
  }
}

function getItemLabel(item: TransferListItem): string {
  return item[props.labelKey] || item.id;
}

function openInlineCreateDialog() {
  inlineCreateDialogOpen.value = true;
}

function handleInlineCreated(permission: TransferListItem) {
  // Auto-assign newly created permission
  emit('assign', [permission.id]);
}
</script>

<template>
  <div class="grid grid-cols-[1fr_auto_1fr] gap-4 w-full">
    <!-- Available List -->
    <div class="space-y-2">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-medium">
          Available {{ label }}
        </h3>
        <Badge variant="secondary">
          {{ availableItems.length }}
        </Badge>
      </div>

      <Command class="border rounded-lg">
        <CommandInput :placeholder="`Search ${label.toLowerCase()}...`" />
        <CommandList>
          <CommandEmpty>No {{ label.toLowerCase() }} found.</CommandEmpty>
          <CommandGroup>
            <CommandItem
              v-for="item in availableItems"
              :key="item.id"
              :value="item.id"
              @select="toggleAvailableSelection(item.id)"
            >
              <Checkbox
                :checked="availableSelected.includes(item.id)"
                class="mr-2"
              />
              <div class="flex-1 flex items-center">
                <span>{{ getItemLabel(item) }}</span>
                <!-- Description Tooltip for Permissions -->
                <TooltipProvider v-if="showTooltip && item.description">
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <Button type="button" variant="ghost" size="icon" class="h-4 w-4 ml-1">
                        <Info class="h-3 w-3" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p class="text-xs">
                        {{ item.description }}
                      </p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </div>
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>

      <!-- Inline Create Button (for permissions only) -->
      <Button
        v-if="allowInlineCreate"
        type="button"
        variant="outline"
        size="sm"
        class="w-full"
        @click="openInlineCreateDialog"
      >
        <Plus class="h-4 w-4 mr-2" />
        Create New {{ singularLabel }}
      </Button>
    </div>

    <!-- Arrow Buttons -->
    <div class="flex flex-col justify-center gap-2">
      <Button
        type="button"
        variant="outline"
        size="icon"
        :disabled="availableSelected.length === 0 || isLoading"
        @click="assignSelected"
      >
        <ChevronRight class="h-4 w-4" />
      </Button>
      <Button
        type="button"
        variant="outline"
        size="icon"
        :disabled="assignedSelected.length === 0 || isLoading"
        @click="removeSelected"
      >
        <ChevronLeft class="h-4 w-4" />
      </Button>
    </div>

    <!-- Assigned List -->
    <div class="space-y-2">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-medium">
          Assigned {{ label }}
        </h3>
        <Badge variant="secondary">
          {{ assignedItems.length }}
        </Badge>
      </div>

      <Command class="border rounded-lg">
        <CommandInput :placeholder="`Search ${label.toLowerCase()}...`" />
        <CommandList>
          <CommandEmpty>No {{ label.toLowerCase() }} assigned.</CommandEmpty>
          <CommandGroup>
            <CommandItem
              v-for="item in assignedItems"
              :key="item.id"
              :value="item.id"
              @select="toggleAssignedSelection(item.id)"
            >
              <Checkbox
                :checked="assignedSelected.includes(item.id)"
                class="mr-2"
              />
              <div class="flex-1 flex items-center">
                <span>{{ getItemLabel(item) }}</span>
                <!-- Description Tooltip for Permissions -->
                <TooltipProvider v-if="showTooltip && item.description">
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <Button type="button" variant="ghost" size="icon" class="h-4 w-4 ml-1">
                        <Info class="h-3 w-3" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                      <p class="text-xs">
                        {{ item.description }}
                      </p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </div>
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>
    </div>
  </div>

  <!-- Inline Create Dialog -->
  <InlinePermissionCreateDialog
    v-if="allowInlineCreate"
    v-model:open="inlineCreateDialogOpen"
    @created="handleInlineCreated"
  />
</template>
