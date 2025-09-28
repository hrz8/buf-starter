<script setup lang="ts">
import type { Row } from '@tanstack/vue-table';

import {
  DropdownMenuSeparator,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownMenuItem,
  DropdownMenu,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';

interface DataTableBasicRowActionsProps {
  row: Row<any>;
  actions?: {
    edit?: boolean;
    duplicate?: boolean;
    favorite?: boolean;
    delete?: boolean;
  };
}

const props = withDefaults(defineProps<DataTableBasicRowActionsProps>(), {
  actions: () => ({
    edit: true,
    duplicate: true,
    favorite: true,
    delete: true,
  }),
});

const emit = defineEmits<{
  edit: [row: Row<any>];
  duplicate: [row: Row<any>];
  favorite: [row: Row<any>];
  delete: [row: Row<any>];
}>();

function handleEdit() {
  emit('edit', props.row);
}

function handleDuplicate() {
  emit('duplicate', props.row);
}

function handleFavorite() {
  emit('favorite', props.row);
}

function handleDelete() {
  emit('delete', props.row);
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button
        variant="ghost"
        class="flex h-8 w-8 p-0 data-[state=open]:bg-muted"
      >
        <Icon
          name="radix-icons:dots-horizontal"
          class="h-4 w-4"
        />
        <span class="sr-only">Open menu</span>
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent
      align="end"
      class="w-[160px]"
    >
      <DropdownMenuItem
        v-if="props.actions.edit"
        @click="handleEdit"
      >
        Edit
      </DropdownMenuItem>

      <DropdownMenuItem
        v-if="props.actions.duplicate"
        @click="handleDuplicate"
      >
        Make a copy
      </DropdownMenuItem>

      <DropdownMenuItem
        v-if="props.actions.favorite"
        @click="handleFavorite"
      >
        Favorite
      </DropdownMenuItem>

      <DropdownMenuSeparator v-if="props.actions.delete" />

      <DropdownMenuItem
        v-if="props.actions.delete"
        class="text-destructive focus:text-destructive-foreground"
        @click="handleDelete"
      >
        Delete
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
