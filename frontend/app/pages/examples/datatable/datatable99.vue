<!-- step 99 -->
<script setup lang="ts">
import type { Employee } from '#shared/repository/example';

import type { ColumnFiltersState, PaginationState, SortingState, Table } from '@tanstack/vue-table';

import type { QueryOptions } from '~~/shared/types/query';
import { exampleRepository } from '#shared/repository/example';
import {

  createColumnHelper,

} from '@tanstack/vue-table';

import { DataTable, DataTableFacetedFilter } from '@/components/datatable-sample-only';
import { Input } from '@/components/ui/input';

const example = exampleRepository();

const page = ref(1);
const pageSize = ref(10);
const keyword = ref('');

const columnFilters = ref<ColumnFiltersState>([]);
const debouncedColumnFilters = refDebounced(columnFilters, 500);

const sorting = ref<SortingState>([]);

const queryOptions = computed<QueryOptions>(() => {
  const opts: QueryOptions = {
    pagination: {
      page: page.value,
      pageSize: pageSize.value,
    },
    keyword: keyword.value,
  };

  if (debouncedColumnFilters.value.length > 0) {
    const filters: NonNullable<QueryOptions['filters']> = {};
    for (const f of debouncedColumnFilters.value) {
      filters[f.id] = f.value as string | string[] | number | boolean | null | undefined;
    }
    opts.filters = filters;
  }

  if (sorting.value.length > 0 && sorting.value[0]) {
    opts.sorting = {
      field: String(sorting.value[0].id),
      order: sorting.value[0].desc ? 'desc' : 'asc',
    };
  }

  return opts;
});

function serializeFilters(filters: NonNullable<QueryOptions['filters']>): string {
  return Object.keys(filters)
    .sort()
    .map(key => `${key}:${filters[key]}`)
    .join('|');
}

const asyncDataKey = computed(() => {
  const { pagination, keyword, filters, sorting } = queryOptions.value;
  const { page, pageSize } = pagination;
  const keys = [
    'employee-table',
    page,
    pageSize,
    keyword,
    filters ? serializeFilters(filters) : null,
    sorting ? `${sorting.field}:${sorting.order}` : null,
  ];
  return keys.filter(Boolean).join('-');
});

const { data: response, pending, refresh } = useAsyncData(
  asyncDataKey,
  () => example.query(queryOptions.value),
  {
    server: false,
    watch: [queryOptions],
    immediate: true,
  },
);

const data = computed(() => response.value?.data ?? []);
const rowCount = computed(() => response.value?.meta.rowCount ?? 0);

const columnHelper = createColumnHelper<Employee>();
const columns = [
  columnHelper.accessor('id', {
    header: 'ID',
    cell: info => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('name', {
    header: 'Name',
    cell: info => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('email', {
    header: 'Email',
    cell: info => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('role', {
    header: 'Role',
    cell: info => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('department', {
    header: 'Department',
    cell: info => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('status', {
    header: 'Status',
    cell: (info) => {
      const status = info.getValue();
      return h('span', {
        class: [
          'inline-flex items-center rounded-full px-2 py-1 text-xs font-medium',
          status === 'active'
            ? 'bg-green-100 text-green-800'
            : status === 'inactive'
              ? 'bg-red-100 text-red-800'
              : 'bg-gray-100 text-gray-800',
        ],
      }, status);
    },
    enableSorting: true,
  }),
  columnHelper.accessor('createdAt', {
    header: 'Created At',
    cell: (info) => {
      const date = new Date(info.getValue());
      return date.toLocaleDateString();
    },
    enableSorting: true,
  }),
];

function onDataTableUpdate(p: PaginationState, s: SortingState, cf: ColumnFiltersState) {
  page.value = p.pageIndex + 1;
  pageSize.value = p.pageSize;
  sorting.value = s;
  columnFilters.value = cf;
}

const statusFilterValues = ref<string[]>([]);

function onStatusFilter(table: Table<any>, selected: string[]) {
  table.getColumn('status')?.setFilterValue(selected);
}

function onClearStatusFilter(table?: Table<any>) {
  statusFilterValues.value = [];
  if (table) {
    table.getColumn('status')?.setFilterValue(undefined);
  }
}
</script>

<template>
  <div class="px-5 py-3">
    <DataTable
      :columns="columns"
      :row-count="rowCount"
      :page="page"
      :page-size="pageSize"
      :data="data"
      :pending="pending"
      @update="onDataTableUpdate"
      @refresh="refresh"
      @reset-column-filter="() => onClearStatusFilter()"
    >
      <template #filters="{ table }">
        <Input
          v-model="keyword"
          placeholder="Search employees..."
          class="h-8 w-[150px] lg:w-[250px]"
        />
        <DataTableFacetedFilter
          v-model="statusFilterValues"
          title="Status"
          :options="[
            { label: 'Active', value: 'active' },
            { label: 'Inactive', value: 'inactive' },
          ]"
          @update="(selected) => onStatusFilter(table, selected)"
          @clear="() => onClearStatusFilter(table)"
        />
      </template>
    </DataTable>
  </div>
</template>
