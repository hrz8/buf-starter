<!-- step 16 -->
<script setup lang="ts">
import type { Employee } from '#shared/repository/example';

import { serializeFilters } from '#shared/helpers/serializer';

import { exampleRepository } from '#shared/repository/example';
import { createColumnHelper } from '@tanstack/vue-table';

import {
  DataTable,
  DataTableBasicRowActions,
  DataTableColumnHeader,
  DataTableFacetedFilter,
} from '@/components/custom/datatable';
import {
  useDataTableFilter,
  useDataTableState,
} from '@/components/custom/datatable/utils';
import { Input } from '@/components/ui/input';
import { useServerTableQuery } from '@/composables/useServerTableQuery';

const example = exampleRepository();

const page = ref(1);
const pageSize = ref(10);
const keyword = ref('');

const dataTableRef = ref<InstanceType<typeof DataTable> | null>(null);
const table = computed(() => dataTableRef.value?.table);

const { columnFilters, sorting } = useDataTableState(dataTableRef);

const { queryOptions } = useServerTableQuery({
  page,
  pageSize,
  keyword,
  columnFilters,
  sorting,
});

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
} = useLazyAsyncData(
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
const filters = computed(() => response.value?.meta.filters);

const columnHelper = createColumnHelper<Employee>();
const columns = [
  columnHelper.accessor('id', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'ID',
    }),
    cell: info => h('div', { class: 'w-20' }, info.getValue()),
    enableSorting: true,
  }),
  columnHelper.accessor('name', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Name',
    }),
    cell: info => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('email', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Email',
    }),
    cell: info => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('role', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Role',
    }),
    cell: info => info.getValue(),
    enableSorting: true,
  }),
  columnHelper.accessor('department', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Department',
    }),
    cell: info => info.getValue(),
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
    cell: ({ row }) => h(DataTableBasicRowActions, { row }),
  }),
]; ;

// role filter
const roleFilter = useDataTableFilter(table, 'role');
const roleOptions = computed(() =>
  filters.value?.roles?.map((role: string) => ({
    label: role,
    value: role,
  })) ?? [],
);

// department filter
const departmentFilter = useDataTableFilter(table, 'department');
const departmentOptions = computed(() =>
  filters.value?.departments?.map((dept: string) => ({
    label: dept,
    value: dept,
  })) ?? [],
);

// Reset all filters
function reset() {
  roleFilter.clearFilter();
  departmentFilter.clearFilter();
}
</script>

<template>
  <div class="space-y-5 px-4 py-3 sm:px-6 lg:px-8">
    <div class="container mx-auto flex justify-end">
      <Button size="sm">
        Add Employee
      </Button>
    </div>
    <div class="container mx-auto">
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
        <template #filters>
          <Input
            v-model="keyword"
            placeholder="Search employees..."
            class="h-8 w-[150px] lg:w-[250px]"
          />

          <DataTableFacetedFilter
            v-if="roleOptions.length > 0"
            v-model="roleFilter.filterValues.value"
            title="Role"
            :options="roleOptions"
            @update="roleFilter.setFilter"
            @clear="roleFilter.clearFilter"
          />

          <DataTableFacetedFilter
            v-if="departmentOptions.length > 0"
            v-model="departmentFilter.filterValues.value"
            title="Department"
            :options="departmentOptions!"
            @update="departmentFilter.setFilter"
            @clear="departmentFilter.clearFilter"
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
  </div>
</template>
