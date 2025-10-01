<script setup lang="ts">
import type { ColumnDef, ColumnFiltersState, SortingState, VisibilityState } from '@tanstack/vue-table';
import {
  getCoreRowModel,
  useVueTable,
} from '@tanstack/vue-table';

import {
  DataTableContent,
  DataTablePagination,
  DataTableToolbar,
} from '@/components/custom/datatable';
import { valueUpdater } from '@/components/ui/table/utils';

interface Props {
  columns: ColumnDef<any, any>[];
  data: any[];
  pending?: boolean;
  page?: number;
  pageSize?: number;
  rowCount?: number;
  columnPrefix?: string;
}

interface Emits {
  'update:page': [value: number];
  'update:pageSize': [value: number];
  'refresh': [];
  'reset': [];
}

const props = withDefaults(defineProps<Props>(), {
  pending: false,
  page: 1,
  pageSize: 10,
  rowCount: 0,
  columnPrefix: 'columns',
});

const emit = defineEmits<Emits>();

defineSlots<{
  filters?: (props: { table?: typeof table }) => any;
  loading?: () => any;
  empty?: () => any;
}>();

// Table state
const columnFilters = ref<ColumnFiltersState>([]);
const debouncedColumnFilters = refDebounced(columnFilters, 500);
const sorting = ref<SortingState>([]);
const columnVisibility = ref<VisibilityState>({});

// Table instance
const table = useVueTable({
  get columns() { return props.columns; },
  get data() { return props.data; },
  get rowCount() { return props.rowCount; },

  // Core
  getCoreRowModel: getCoreRowModel(),

  // Server-side configuration
  manualPagination: true,
  manualSorting: true,
  manualFiltering: true,

  // State management
  state: {
    get columnFilters() { return debouncedColumnFilters.value; },
    get sorting() { return sorting.value; },
    get columnVisibility() { return columnVisibility.value; },
  },

  // Change handlers
  onColumnFiltersChange: updater => valueUpdater(updater, columnFilters),
  onSortingChange: updater => valueUpdater(updater, sorting),
  onColumnVisibilityChange: updater => valueUpdater(updater, columnVisibility),
});

// Pagination state
const currentPage = computed({
  get: () => props.page,
  set: (value: number) => emit('update:page', value),
});

const currentPageSize = computed({
  get: () => props.pageSize,
  set: (value: number) => emit('update:pageSize', value),
});

defineExpose({
  table,
});
</script>

<template>
  <div class="space-y-4">
    <DataTableToolbar
      :table="table"
      :column-prefix="columnPrefix"
      @refresh="emit('refresh')"
      @reset="emit('reset')"
    >
      <template #filters>
        <slot
          name="filters"
          :table="table"
        />
      </template>
    </DataTableToolbar>
    <DataTableContent
      :table="table"
      :pending="pending"
      :show-border="true"
    >
      <template #loading>
        <slot name="loading" />
      </template>
      <template #empty>
        <slot name="empty" />
      </template>
    </DataTableContent>
    <DataTablePagination
      v-model:page="currentPage"
      v-model:page-size="currentPageSize"
      :row-count="rowCount"
      :pending="pending"
    />
  </div>
</template>
