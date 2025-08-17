<script setup lang="ts">
import type { Table } from '@tanstack/vue-table';

import {
  DropdownMenuCheckboxItem,
  DropdownMenuSeparator,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownMenuLabel,
  DropdownMenu,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';

interface DataTableViewOptionsProps {
  table: Table<any>;
}

const props = defineProps<DataTableViewOptionsProps>();

const columns = computed(() => props.table.getAllColumns()
  .filter(
    (column) =>
      typeof column.accessorFn !== 'undefined' && column.getCanHide(),
  ));
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button
        variant="outline"
        size="sm"
        class="ml-auto hidden h-8 lg:flex"
      >
        <Icon
          name="radix-icons:mixer-horizontal"
          class="mr-2 h-4 w-4"
        />
        View
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent
      align="end"
      class="w-[150px]"
    >
      <DropdownMenuLabel>Toggle columns</DropdownMenuLabel>
      <DropdownMenuSeparator />

      <DropdownMenuCheckboxItem
        v-for="column in columns"
        :key="column.id"
        class="capitalize"
        :model-value="column.getIsVisible()"
        @update:model-value="(value) => column.toggleVisibility(!!value)"
      >
        {{ column.id }}
      </DropdownMenuCheckboxItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
