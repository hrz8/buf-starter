<script setup lang="ts">
import type { Role } from '~~/gen/altalune/v1/role_pb';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

const props = defineProps<{
  role: Role;
}>();

const emit = defineEmits<{
  edit: [];
  delete: [];
}>();

const { t } = useI18n();

// Protected role that cannot be deleted
const PROTECTED_ROLE = 'superadmin';

// Check if this role is protected
const isProtected = computed(() => props.role.name === PROTECTED_ROLE);

function handleEdit() {
  emit('edit');
}

function handleDelete() {
  if (isProtected.value) {
    return; // Prevent deletion of protected role
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
        {{ t('features.roles.actions.edit') }}
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
        {{ t('features.roles.actions.delete') }}
        <span v-if="isProtected" class="ml-2 text-xs text-muted-foreground">(Protected)</span>
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
