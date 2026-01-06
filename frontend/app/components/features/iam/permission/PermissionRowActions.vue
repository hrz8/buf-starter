<script setup lang="ts">
import type { Permission } from '~~/gen/altalune/v1/permission_pb';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

const props = defineProps<{
  permission: Permission;
}>();

const emit = defineEmits<{
  edit: [];
  delete: [];
}>();

const { t } = useI18n();

// Protected permission that cannot be deleted
const PROTECTED_PERMISSION = 'root';

// Check if this permission is protected
const isProtected = computed(() => props.permission.name === PROTECTED_PERMISSION);

function handleEdit() {
  emit('edit');
}

function handleDelete() {
  if (isProtected.value) {
    return; // Prevent deletion of protected permission
  }
  emit('delete');
}
</script>

<template>
  <!-- Actions dropdown -->
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button
        variant="ghost"
        class="h-8 w-8 p-0"
        aria-label="Actions"
      >
        <Icon name="lucide:more-horizontal" class="h-4 w-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end">
      <DropdownMenuItem
        class="cursor-pointer"
        @click="handleEdit"
      >
        <Icon name="lucide:edit" class="mr-2 h-4 w-4" />
        {{ t('features.permissions.actions.edit') }}
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuItem
        :disabled="isProtected"
        :class="[
          isProtected
            ? 'cursor-not-allowed opacity-50'
            : 'cursor-pointer text-destructive focus:text-destructive',
        ]"
        @click="handleDelete"
      >
        <Icon name="lucide:trash-2" class="mr-2 h-4 w-4" />
        {{ t('features.permissions.actions.delete') }}
        <span v-if="isProtected" class="ml-2 text-xs text-muted-foreground">(Protected)</span>
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
