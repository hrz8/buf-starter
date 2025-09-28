<!-- step 9 -->
<script setup lang="ts">
import type { Employee } from '#shared/repository/example';

import type { QueryOptions } from '#shared/types/query';

import type { ColumnFiltersState, SortingState, VisibilityState } from '@tanstack/vue-table';
import { exampleRepository } from '#shared/repository/example';
import {

  createColumnHelper,

  FlexRender,
  getCoreRowModel,
  useVueTable,
} from '@tanstack/vue-table';

import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationFirst,
  PaginationItem,
  PaginationLast,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableEmpty,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { valueUpdater } from '@/components/ui/table/utils';

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

function goToPage(newPage: number) {
  if (newPage >= 1 && newPage <= pageCount.value) {
    page.value = newPage;
  }
}

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
    get sorting() { return sorting.value; },
    get columnVisibility() { return columnVisibility.value; },
  },
  onColumnFiltersChange: updater => valueUpdater(updater, columnFilters),
  onSortingChange: updater => valueUpdater(updater, sorting),
  onColumnVisibilityChange: updater => valueUpdater(updater, columnVisibility),
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

const showSelector = ref(false);
const availableColumns = computed(() => table.getAllColumns()
  .filter(
    column =>
      typeof column.accessorFn !== 'undefined' && column.getCanHide(),
  ));

function onSelectColumnView(e: Event, id: string) {
  table.getColumn(id)?.toggleVisibility((e.target as HTMLInputElement).checked);
}

function getSortIcon(column: any) {
  const sortDirection = column.getIsSorted();
  if (sortDirection === 'asc')
    return 'lucide:arrow-up';
  if (sortDirection === 'desc')
    return 'lucide:arrow-down';
  return 'lucide:arrow-up-down';
}

function getSortIconClass(column: any) {
  const sortDirection = column.getIsSorted();
  return [
    'w-4 h-4 ml-2 transition-colors',
    sortDirection ? 'text-foreground' : 'text-muted-foreground',
  ];
}

// Helper to generate visible page numbers
const visiblePages = computed(() => {
  const current = page.value;
  const total = pageCount.value;
  const delta = 2;

  const range = [];
  const rangeWithDots = [];

  for (let i = Math.max(2, current - delta); i <= Math.min(total - 1, current + delta); i++) {
    range.push(i);
  }

  if (current - delta > 2) {
    rangeWithDots.push(1, '...');
  }
  else {
    rangeWithDots.push(1);
  }

  rangeWithDots.push(...range);

  if (current + delta < total - 1) {
    rangeWithDots.push('...', total);
  }
  else {
    rangeWithDots.push(total);
  }

  return rangeWithDots
    .filter((item, index, arr) => arr.indexOf(item) === index && (item !== 1 || index === 0));
});
</script>

<template>
  <div class="px-5 py-3">
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

    <!-- filters -->
    <div class="flex items-center space-x-4">
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

      <!-- column visibility -->
      <div class="relative">
        <button
          class="border border-gray-300 rounded px-3 py-2 hover:bg-gray-50"
          @click="showSelector = !showSelector"
        >
          <Icon
            name="lucide:columns"
            class="w-4 h-4 mr-2"
          />
          Columns
        </button>

        <div
          v-if="showSelector"
          class="absolute mt-1 w-48 bg-white border border-gray-200 rounded shadow-lg z-10"
        >
          <div
            v-for="col in availableColumns"
            :key="col.id"
            class="px-3 py-2 hover:bg-gray-50 flex items-center gap-2"
          >
            <input
              :id="`column-toggle-${col.id}`"
              type="checkbox"
              :checked="col.getIsVisible()"
              @change="(e) => onSelectColumnView(e, col.id)"
            >
            <label
              :for="`column-toggle-${col.id}`"
              class="flex-1 cursor-pointer"
            >{{ col.id }}</label>
          </div>
        </div>
      </div>
    </div>

    <!-- table -->
    <div class="rounded-md border mb-3">
      <Table>
        <TableHeader>
          <TableRow
            v-for="headerGroup in table.getHeaderGroups()"
            :key="headerGroup.id"
          >
            <TableHead
              v-for="header in headerGroup.headers"
              :key="header.id"
              :class="header.column.getCanSort() ? 'cursor-pointer select-none' : ''"
              @click="header.column.getCanSort() && header.column.toggleSorting()"
            >
              <div
                v-if="!header.isPlaceholder"
                class="flex items-center"
              >
                <FlexRender
                  :render="header.column.columnDef.header"
                  :props="header.getContext()"
                />
                <Icon
                  v-if="header.column.getCanSort()"
                  :name="getSortIcon(header.column)"
                  :class="getSortIconClass(header.column)"
                />
              </div>
            </TableHead>
          </TableRow>
        </TableHeader>

        <TableBody>
          <TableEmpty
            v-if="pending"
            :colspan="table.getVisibleFlatColumns().length"
          >
            <div class="flex items-center justify-center space-x-2">
              <Icon
                name="lucide:loader-2"
                class="w-4 h-4 animate-spin"
              />
              <span>Loading employees...</span>
            </div>
          </TableEmpty>

          <TableEmpty
            v-else-if="!data.length"
            :colspan="table.getVisibleFlatColumns().length"
          >
            <div class="flex flex-col items-center justify-center space-y-2 text-muted-foreground">
              <Icon
                name="lucide:users"
                class="w-8 h-8"
              />
              <span>No employees found</span>
              <span class="text-sm">Try adjusting your search or filters</span>
            </div>
          </TableEmpty>

          <TableRow
            v-for="row in table.getRowModel().rows"
            v-else
            :key="row.id"
          >
            <TableCell
              v-for="cell in row.getVisibleCells()"
              :key="cell.id"
            >
              <FlexRender
                :render="cell.column.columnDef.cell"
                :props="cell.getContext()"
              />
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <!-- page info -->
    <div class="flex items-center justify-between px-2">
      <div class="flex-1 text-sm text-muted-foreground">
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

      <div class="flex items-center space-x-6 lg:space-x-8">
        <!-- page limit selector -->
        <div class="flex items-center space-x-2">
          <p class="text-sm font-medium">
            Rows per page
          </p>
          <Select
            v-model="pageSize"
          >
            <SelectTrigger class="h-8 w-[70px]">
              <SelectValue :placeholder="`${pageSize}`" />
            </SelectTrigger>
            <SelectContent side="top">
              <SelectItem
                v-for="size in [10, 20, 25, 50]"
                :key="size"
                :value="size"
              >
                {{ size }}
              </SelectItem>
            </SelectContent>
          </Select>

          <!-- pagination -->
          <div class="flex items-center space-x-2">
            <Pagination
              v-model:page="page"
              :total="rowCount"
              :items-per-page="pageSize"
              :sibling-count="1"
              show-edges
            >
              <PaginationContent>
                <PaginationFirst
                  v-if="pageCount > 5"
                  @click="goToPage(1)"
                />

                <PaginationPrevious
                  :disabled="isFirstPage || pending"
                  @click="prev()"
                />

                <template
                  v-for="(item, index) in visiblePages"
                  :key="index"
                >
                  <PaginationEllipsis v-if="item === '...'" />
                  <PaginationItem
                    v-else
                    :value="Number(item)"
                    :is-active="item === page"
                    @click="goToPage(Number(item))"
                  >
                    {{ item }}
                  </PaginationItem>
                </template>

                <PaginationNext
                  :disabled="isLastPage || pending"
                  @click="next()"
                />

                <PaginationLast
                  v-if="pageCount > 5"
                  @click="goToPage(pageCount)"
                />
              </PaginationContent>
            </Pagination>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
