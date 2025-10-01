<script setup lang="ts">
import type { Table } from '@tanstack/vue-table';
import type { Data } from '.';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

const props = withDefaults(defineProps<Props>(), {
  label: undefined,
  columnPrefix: 'columns',
});

const { t } = useI18n();

interface Props {
  table: Table<Data>;
  label?: string;
  columnPrefix?: string;
}

const availableColumns = computed(() => props.table
  .getAllColumns()
  .filter(
    column =>
      typeof column.accessorFn !== 'undefined' && column.getCanHide(),
  ),
);

/**
 * Get translated column name using column ID with prefix as translation key
 */
function getColumnLabel(column: any): string {
  return t(`${props.columnPrefix}.${column.id}`);
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button
        variant="outline"
        size="sm"
        class="ml-auto h-8"
      >
        <Icon
          name="lucide:settings-2"
          class="mr-2 h-4 w-4"
        />
        {{ label ?? t('datatable.view') }}
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent
      align="end"
      class="w-[150px]"
    >
      <DropdownMenuLabel>{{ t('datatable.toggleColumns') }}</DropdownMenuLabel>
      <DropdownMenuSeparator />
      <DropdownMenuCheckboxItem
        v-for="column in availableColumns"
        :key="column.id"
        class="capitalize"
        :model-value="column.getIsVisible()"
        @update:model-value="(value) => column.toggleVisibility(!!value)"
      >
        {{ getColumnLabel(column) }}
      </DropdownMenuCheckboxItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
