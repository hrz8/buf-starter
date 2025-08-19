<!-- step 15 -->
<script setup lang="ts">
import {
  createColumnHelper,
  type Table,
} from '@tanstack/vue-table';

import type { Employee } from '#shared/repository/example';
import type { QueryOptions } from '#shared/types/query';

import { exampleRepository } from '#shared/repository/example';

import {
  DataTableFacetedFilter,
  DataTableColumnHeader,
  DataTableRowActions,
  DataTable,
} from '@/components/datatable';
import { Input } from '@/components/ui/input';

const example = exampleRepository();

const page = ref(1);
const pageSize = ref(10);
const keyword = ref('');
const debouncedKeyword = refDebounced(keyword, 500);

const dataTableRef = ref<InstanceType<typeof DataTable> | null>(null);

const columnFilters = computed(() => dataTableRef.value?.table.getState().columnFilters);
const sorting = computed(() => dataTableRef.value?.table.getState().sorting);

const queryOptions = computed<QueryOptions>(() => {
  const opts: QueryOptions = {
    pagination: {
      page: page.value,
      pageSize: pageSize.value,
    },
    keyword: debouncedKeyword.value,
  };

  if (columnFilters.value && columnFilters.value.length > 0) {
    const filters: NonNullable<QueryOptions['filters']> = {};
    for (const f of columnFilters.value) {
      filters[f.id] = f.value as string | number | boolean | null | undefined;
    }
    opts.filters = filters;
  }

  if (sorting.value && sorting.value.length > 0 && sorting.value[0]) {
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

// role filters
const roleFilterValues = ref<string[]>([]);
const roleOptions = computed(() => response.value?.meta.filters?.['roles']?.map((role) => ({
  label: role,
  value: role,
})) ?? []);

function onRoleFilterChange(table: Table<Employee>, selected: string[]) {
  table.getColumn('role')?.setFilterValue(selected);
}

function onRoleFilterClear(table?: Table<Employee>) {
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

function onDepartmentFilterChange(table: Table<Employee>, selected: string[]) {
  table.getColumn('department')?.setFilterValue(selected.length ? selected : undefined);
}

function onDepartmentFilterClear(table?: Table<Employee>) {
  departmentFilterValues.value = [];
  if (table) {
    table.getColumn('department')?.setFilterValue(undefined);
  }
}

function reset() {
  roleFilterValues.value = [];
  departmentFilterValues.value = [];
}
</script>

<template>
  <div class="px-5 py-3">
    <DataTable
      ref="dataTableRef"
      v-model:page="page"
      v-model:page-size="pageSize"
      :columns="columns"
      :data="data"
      :pending="pending"
      :row-count="rowCount"
      @refresh="refresh()"
      @reset="reset()"
    >
      <template #filters="{ table }">
        <Input
          v-model="keyword"
          placeholder="Search employees..."
          class="h-8 w-[150px] lg:w-[250px]"
        />

        <DataTableFacetedFilter
          v-if="roleOptions.length > 0 && table"
          v-model="roleFilterValues"
          title="Role"
          :options="roleOptions"
          @update="(selected) => onRoleFilterChange(table, selected)"
          @clear="() => onRoleFilterClear(table)"
        />

        <DataTableFacetedFilter
          v-if="departmentOptions.length > 0 && table"
          v-model="departmentFilterValues"
          title="Department"
          :options="departmentOptions!"
          @update="(selected) => onDepartmentFilterChange(table, selected)"
          @clear="() => onDepartmentFilterClear(table)"
        />
      </template>
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
    </DataTable>
  </div>
</template>
