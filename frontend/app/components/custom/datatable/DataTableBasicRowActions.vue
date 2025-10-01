<script setup lang="ts">
import type { Row } from '@tanstack/vue-table';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

const props = withDefaults(defineProps<DataTableBasicRowActionsProps>(), {
  actions: () => ({
    edit: true,
    duplicate: true,
    delete: true,
  }),
});

const emit = defineEmits<{
  edit: [row: Row<any>];
  duplicate: [row: Row<any>];
  delete: [row: Row<any>];
}>();

const { t } = useI18n();

interface DataTableBasicRowActionsProps {
  row: Row<any>;
  actions?: {
    edit?: boolean;
    duplicate?: boolean;
    delete?: boolean;
  };
}

function handleEdit() {
  emit('edit', props.row);
}

function handleDuplicate() {
  emit('duplicate', props.row);
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
        <span class="sr-only">{{ t('datatable.openMenu') }}</span>
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
        <Icon name="lucide:edit" class="mr-2 h-4 w-4" />
        {{ t('common.btn.edit') }}
      </DropdownMenuItem>

      <DropdownMenuItem
        v-if="props.actions.duplicate"
        @click="handleDuplicate"
      >
        <Icon name="lucide:copy" class="mr-2 h-4 w-4" />
        {{ t('datatable.makeCopy') }}
      </DropdownMenuItem>

      <DropdownMenuSeparator v-if="props.actions.delete" />

      <DropdownMenuItem
        v-if="props.actions.delete"
        class="text-destructive focus:text-destructive-foreground"
        @click="handleDelete"
      >
        <Icon name="lucide:trash-2" class="mr-2 h-4 w-4" />
        {{ t('common.btn.delete') }}
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
