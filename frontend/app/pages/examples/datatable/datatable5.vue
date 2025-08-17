<!-- step 5 -->
<script setup lang="ts">
import {
  type ColumnFiltersState, createColumnHelper, getCoreRowModel, useVueTable, FlexRender,
} from '@tanstack/vue-table';

import type { QueryOptions } from '#shared/types/query';

import { exampleRepository, type Employee } from '#shared/repository/example';

const example = exampleRepository();

const page = ref(1);
const pageSize = ref(10);
const keyword = ref('');

const columnFilters = ref<ColumnFiltersState>([]);
const debouncedColumnFilters = refDebounced(columnFilters, 500);

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

  return opts;
});

function serializeFilters(filters: NonNullable<QueryOptions['filters']>): string {
  return Object.keys(filters)
    .sort()
    .map((key) => `${key}:${filters[key]}`)
    .join('|');
}

const asyncDataKey = computed(() => {
  const { pagination, keyword, filters } = queryOptions.value;
  const { page, pageSize } = pagination;
  const keys = [
    'employee-table',
    page,
    pageSize,
    keyword,
    filters ? serializeFilters(filters) : null,
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

const {
  pageCount,
  isFirstPage,
  isLastPage,
  prev,
  next,
} = useOffsetPagination({
  total: rowCount,
  page,
  pageSize,
});

const goToPage = (newPage: number) => {
  if (newPage >= 1 && newPage <= pageCount.value) {
    page.value = newPage;
  }
};

const columnHelper = createColumnHelper<Employee>();
const columns = [
  columnHelper.accessor('id', {
    header: 'ID',
    cell: (info) => info.getValue(),
  }),
  columnHelper.accessor('name', {
    header: 'Name',
    cell: (info) => info.getValue(),
  }),
  columnHelper.accessor('email', {
    header: 'Email',
    cell: (info) => info.getValue(),
  }),
  columnHelper.accessor('role', {
    header: 'Role',
    cell: (info) => info.getValue(),
  }),
  columnHelper.accessor('department', {
    header: 'Department',
    cell: (info) => info.getValue(),
  }),
  columnHelper.accessor('status', {
    header: 'Status',
    cell: (info) => info.getValue(),
  }),
  columnHelper.accessor('createdAt', {
    header: 'Created At',
    cell: (info) => info.getValue(),
  }),
];

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
    get columnFilters() { return columnFilters.value; },
  },
  onColumnFiltersChange: (updater) => {
    columnFilters.value = typeof updater === 'function'
      ? updater(columnFilters.value)
      : updater;
  },
});

const departmentFilter = ref((table.getColumn('department')?.getFilterValue() ?? '') as string);
function updateDepartmentFilter(value: string) {
  departmentFilter.value = value;
  table.getColumn('department')?.setFilterValue(value);
}

const roleFilter = ref((table.getColumn('role')?.getFilterValue() ?? '') as string);
function updateRoleFilter(value: string) {
  roleFilter.value = value;
  table.getColumn('role')?.setFilterValue(value);
}
</script>

<template>
  <div>
    <h2>User Management</h2>
    <input
      v-model="keyword"
      type="text"
      placeholder="Search..."
      class="border border-gray-300 rounded px-3 py-2"
    >
    <button
      class="border border-gray-300 rounded px-3 py-2 hover:bg-gray-50"
      @click="refresh()"
    >
      Refresh
    </button>

    <!-- field filters -->
    <div>
      <input
        v-model="roleFilter"
        type="text"
        placeholder="Role..."
        class="border border-gray-300 rounded px-3 py-2"
        @input="(e) => updateRoleFilter((e.target as HTMLInputElement).value)"
      >
    </div>
    <div>
      <input
        v-model="departmentFilter"
        type="text"
        placeholder="Department..."
        class="border border-gray-300 rounded px-3 py-2"
        @input="(e) => updateDepartmentFilter((e.target as HTMLInputElement).value)"
      >
    </div>

    <!-- table -->
    <div>
      <table>
        <thead>
          <tr
            v-for="headerGroup in table.getHeaderGroups()"
            :key="headerGroup.id"
          >
            <th
              v-for="header in headerGroup.headers"
              :key="header.id"
            >
              <div v-if="!header.isPlaceholder">
                <FlexRender
                  :render="header.column.columnDef.header"
                  :props="header.getContext()"
                />
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="pending">
            <td
              :colspan="6"
              class="text-center"
            >
              Loading...
            </td>
          </tr>
          <tr
            v-for="row in table.getRowModel().rows"
            v-else
            :key="row.id"
          >
            <td
              v-for="cell in row.getVisibleCells()"
              :key="cell.id"
            >
              <FlexRender
                :render="cell.column.columnDef.cell"
                :props="cell.getContext()"
              />
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- page info -->
    <div>
      <span v-if="rowCount > 0">
        Showing
        <span>{{ ((page - 1) * pageSize) + 1 }}</span>
        to
        <span>{{ Math.min(page * pageSize, rowCount) }}</span>
        of
        <span>{{ rowCount }}</span>
        results
      </span>
      <span v-else>No results</span>
    </div>

    <!-- page limit selector -->
    <div>
      <label>Show</label>
      <select
        v-model="pageSize"
        class="px-3 py-1 border border-gray-300 rounded-md"
      >
        <option
          v-for="size in [10, 20, 25, 50]"
          :key="size"
          :value="size"
        >
          {{ size }}
        </option>
      </select>
      <span>entries</span>
    </div>

    <!-- page navigator -->
    <div>
      <button
        class="
          border border-gray-300 rounded px-3 py-2 hover:bg-gray-50
          disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-transparent
        "
        :disabled="isFirstPage || pending"
        @click="prev()"
      >
        Previous
      </button>
      <span>Page {{ queryOptions.pagination.page }} of {{ pageCount }}</span>
      <button
        class="
          border border-gray-300 rounded px-3 py-2 hover:bg-gray-50
          disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-transparent
        "
        :disabled="isLastPage || pending"
        @click="next()"
      >
        Next
      </button>
    </div>

    <!-- manual page navigation -->
    <div>
      <span>Go to</span>
      <input
        type="number"
        :value="page"
        min="1"
        :max="pageCount"
        class="px-3 py-1 border border-gray-300 rounded-md"
        @input="(e) => goToPage(Number((e.target as HTMLInputElement)?.value))"
      >
    </div>
  </div>
</template>
