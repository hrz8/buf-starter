<script setup lang="ts">
import type {
  ColumnDef,
  ColumnFiltersState,
  PaginationState,
  SortingState,
  VisibilityState,
} from '@tanstack/vue-table';

import {
  getCoreRowModel,
  useVueTable,
} from '@tanstack/vue-table';

import { valueUpdater } from '@/components/ui/table/utils';
import DataTableContent from './DataTableContent.vue';
import DataTablePagination from './DataTablePagination.vue';

import DataTableToolbar from './DataTableToolbar.vue';

interface DataTableProps {
  columns: ColumnDef<any, any>[];
  rowCount: number;
  page: number;
  pageSize: number;
  pending: boolean;
  data: any[];
}
const props = defineProps<DataTableProps>();
const emit = defineEmits<{
  update: [
    pagination: PaginationState,
    sorting: SortingState,
    columnFilters: ColumnFiltersState,
  ];
  refresh: [];
  resetColumnFilter: [];
}>();

const sorting = ref<SortingState>([]);
const columnFilters = ref<ColumnFiltersState>([]);
const columnVisibility = ref<VisibilityState>({});
const pagination = ref<PaginationState>({
  pageIndex: props.page - 1,
  pageSize: props.pageSize,
});

const table = useVueTable({
  get data() { return props.data; },
  get columns() { return props.columns; },
  get rowCount() { return props.rowCount; },
  state: {
    get sorting() { return sorting.value; },
    get columnFilters() { return columnFilters.value; },
    get columnVisibility() { return columnVisibility.value; },
    get pagination() { return pagination.value; },
  },
  onSortingChange: (updaterOrValue) => {
    const res = valueUpdater(updaterOrValue, sorting);
    emit('update', pagination.value, sorting.value, columnFilters.value);
    return res;
  },
  onColumnFiltersChange: (updaterOrValue) => {
    const res = valueUpdater(updaterOrValue, columnFilters);
    emit('update', pagination.value, sorting.value, columnFilters.value);
    return res;
  },
  onColumnVisibilityChange: updaterOrValue => valueUpdater(updaterOrValue, columnVisibility),
  onPaginationChange: (updaterOrValue) => {
    const res = valueUpdater(updaterOrValue, pagination);
    emit('update', pagination.value, sorting.value, columnFilters.value);
    return res;
  },
  getCoreRowModel: getCoreRowModel(),

  // server-side
  manualPagination: true,
  manualSorting: true,
  manualFiltering: true,
});

watch(() => props.rowCount, () => {
  nextTick(() => {
    const currentPageIndex = table.getState().pagination.pageIndex;
    const totalPages = table.getPageCount();
    const isOverflow = currentPageIndex >= totalPages && totalPages > 0;

    if (isOverflow) {
      table.setPageIndex(0);
    }
  });
}, { immediate: true });
</script>

<template>
  <div class="space-y-4">
    <DataTableToolbar
      :table="table"
      @refresh="emit('refresh')"
      @reset="emit('resetColumnFilter')"
    >
      <template #filters>
        <slot
          name="filters"
          :table="table"
        />
      </template>
    </DataTableToolbar>
    <DataTableContent
      :columns="props.columns"
      :table="table"
      :pending="props.pending"
    />
    <DataTablePagination :table="table" />
  </div>
</template>
