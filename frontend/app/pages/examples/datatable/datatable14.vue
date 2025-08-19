<!-- step 14 -->
<script setup lang="ts">
import {
  type ColumnFiltersState,
  type VisibilityState,
  createColumnHelper,
  type SortingState,
  getCoreRowModel,
  useVueTable,
  type Table,
} from '@tanstack/vue-table';

import type { Employee } from '#shared/repository/example';
import type { QueryOptions } from '#shared/types/query';

import { exampleRepository } from '#shared/repository/example';

import {
  DataTableFacetedFilter,
  DataTableColumnHeader,
  DataTableRowActions,
  DataTablePagination,
  DataTableToolbar,
  DataTableContent,
} from '@/components/datatable';
import { valueUpdater } from '@/components/ui/table/utils';
import { Input } from '@/components/ui/input';

const example = exampleRepository();

const page = ref(1);
const pageSize = ref(10);
const keyword = ref('');

const columnFilters = ref<ColumnFiltersState>([]);
const debouncedColumnFilters = refDebounced(columnFilters, 500);

const sorting = ref<SortingState>([]);
const columnVisibility = ref<VisibilityState>({});

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
      filters[f.id] = f.value as string | number | boolean | null | undefined;
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
    .map((key) => `${key}:${filters[key]}`)
    .join('|');
}

const asyncDataKey = computed(() => {
  const { pagination, keyword, filters, sorting } = queryOptions.value;
  const {
    page,
    pageSize,
  } = pagination;
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

const {
  data: response,
  pending,
  refresh,
} = useAsyncData(
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
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'ID',
    }),
    cell: (info) => h('div', { class: 'w-20' }, info.getValue()),
    enableSorting: true,
  }),
  columnHelper.accessor('name', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Name',
    }),
    cell: (info) => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('email', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Email',
    }),
    cell: (info) => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('role', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Role',
    }),
    cell: (info) => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('department', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Department',
    }),
    cell: (info) => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('status', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Status',
    }),
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
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Created At',
    }),
    cell: (info) => {
      const date = new Date(info.getValue());
      return date.toLocaleDateString();
    },
    enableSorting: true,
  }),
  columnHelper.display({
    id: 'actions',
    cell: ({ row }) => h(DataTableRowActions, { row }),
  }),
]; ;

const table = useVueTable({
  columns,
  get data() { return data.value; },
  get rowCount() { return rowCount.value; },

  // core row
  getCoreRowModel: getCoreRowModel(),

  // server-side
  manualPagination: true,
  manualSorting: true,
  manualFiltering: true,

  // state (optional if listener used)
  state: {
    get columnFilters() { return debouncedColumnFilters.value; },
    get sorting() { return sorting.value; },
    get columnVisibility() { return columnVisibility.value; },
  },
  onColumnFiltersChange: (updater) => valueUpdater(updater, columnFilters),
  onSortingChange: (updater) => valueUpdater(updater, sorting),
  onColumnVisibilityChange: (updater) => valueUpdater(updater, columnVisibility),
});

// role filters
const roleFilterValues = ref<string[]>([]);
const roleOptions = computed(() => response.value?.meta.filters?.['roles']?.map((role) => ({
  label: role,
  value: role,
})) ?? []);

function onRoleFilterChange(table: Table<any>, selected: string[]) {
  table.getColumn('role')?.setFilterValue(selected);
}

function onRoleFilterClear(table?: Table<any>) {
  roleFilterValues.value = [];
  if (table) {
    table.getColumn('role')?.setFilterValue(undefined);
  }
}

// department filters
const departmentFilterValues = ref<string[]>([]);
const departmentOptions = computed(() => response.value?.meta.filters?.['departments']?.map((dept) => ({
  label: dept,
  value: dept,
})) ?? []);

function onDepartmentFilterChange(table: Table<any>, selected: string[]) {
  table.getColumn('department')?.setFilterValue(selected.length ? selected : undefined);
}

function onDepartmentFilterClear(table?: Table<any>) {
  departmentFilterValues.value = [];
  if (table) {
    table.getColumn('department')?.setFilterValue(undefined);
  }
}

function resetFilters() {
  roleFilterValues.value = [];
  departmentFilterValues.value = [];
}
</script>

<template>
  <div class="px-5 py-3">
    <DataTableToolbar
      :table="table"
      @refresh="refresh()"
      @reset="resetFilters()"
    >
      <template #filters="{ table }">
        <Input
          v-model="keyword"
          placeholder="Search employees..."
          class="h-8 w-[150px] lg:w-[250px]"
        />

        <DataTableFacetedFilter
          v-if="roleOptions.length > 0"
          v-model="roleFilterValues"
          title="Role"
          :options="roleOptions"
          @update="(selected) => onRoleFilterChange(table, selected)"
          @clear="() => onRoleFilterClear(table)"
        />

        <DataTableFacetedFilter
          v-if="departmentOptions.length > 0"
          v-model="departmentFilterValues"
          title="Department"
          :options="departmentOptions!"
          @update="(selected) => onDepartmentFilterChange(table, selected)"
          @clear="() => onDepartmentFilterClear(table)"
        />
      </template>
    </DataTableToolbar>

    <DataTableContent
      :table="table"
      :pending="pending"
      loading-text="Loading employees..."
      no-data-text="No employees found"
      no-data-description="Try adjusting your search or filters"
      no-data-icon="lucide:users"
    >
      <!-- Custom loading state -->
      <template #loading>
        <div class="flex flex-col items-center justify-center space-y-4 py-12">
          <div class="relative">
            <Icon
              name="lucide:loader-2"
              class="w-8 h-8 animate-spin text-primary"
            />
          </div>
          <div class="text-center">
            <p class="font-medium">
              Loading employee data...
            </p>
            <p class="text-sm text-muted-foreground">
              This may take a moment
            </p>
          </div>
        </div>
      </template>

      <!-- Custom empty state -->
      <template #empty>
        <div class="flex flex-col items-center justify-center space-y-6 py-16">
          <div class="relative">
            <Icon
              name="lucide:user-x"
              class="w-16 h-16 text-muted-foreground/50"
            />
          </div>
          <div class="text-center space-y-2">
            <h3 class="text-lg font-semibold">
              No employees found
            </h3>
            <p class="text-muted-foreground max-w-md">
              We couldn't find any employees matching your criteria.
              Try adjusting your filters or search terms.
            </p>
          </div>
          <div class="flex space-x-2">
            <Button size="sm">
              Add Employee
            </Button>
          </div>
        </div>
      </template>
    </DataTableContent>

    <DataTablePagination
      v-model:page="page"
      v-model:page-size="pageSize"
      :row-count="rowCount"
      :pending="pending"
    />
  </div>
</template>
