<!-- step 17 -->
<script setup lang="ts">
import { EmployeeStatus } from '~~/gen/altalune/v1/employee_pb';
import { createColumnHelper } from '@tanstack/vue-table';

import type { Employee } from '~~/gen/altalune/v1/employee_pb';

import { serializeProtoFilters } from '#shared/helpers/serializer';

import {
  DataTableBasicRowActions,
  DataTableFacetedFilter,
  DataTableColumnHeader,
  DataTable,
} from '@/components/custom/datatable';
import {
  SheetDescription,
  SheetTrigger,
  SheetContent,
  SheetHeader,
  SheetTitle,
  Sheet,
} from '@/components/ui/sheet';
import {
  useDataTableFilter,
  useDataTableState,
} from '@/components/custom/datatable/utils';
import { useEmployeeService } from '@/composables/services/useEmployeeService';
import { useQueryRequest } from '@/composables/useQueryRequest';
import { useProjectStore } from '@/stores/project';
import { Input } from '@/components/ui/input';

const { activeProjectId } = useProjectStore();

const { query } = useEmployeeService();

const page = ref(1);
const pageSize = ref(10);
const keyword = ref('');

const dataTableRef = ref<InstanceType<typeof DataTable> | null>(null);
const table = computed(() => dataTableRef.value?.table);

const { columnFilters, sorting } = useDataTableState(dataTableRef);

const { queryRequest } = useQueryRequest({
  page,
  pageSize,
  keyword,
  columnFilters,
  sorting,
});

const asyncDataKey = computed(() => {
  const { pagination, keyword, filters, sorting } = queryRequest.value;
  const { page, pageSize } = pagination!;
  const keys = [
    'employee-table',
    page,
    pageSize,
    keyword,
    filters ? serializeProtoFilters(filters) : null,
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
  () => query({
    projectId: activeProjectId ?? undefined,
    query: queryRequest.value,
  }),
  {
    server: false,
    watch: [queryRequest],
    immediate: false,
  },
);

const data = computed(() => response.value?.data ?? []);
const rowCount = computed(() => response.value?.meta?.rowCount ?? 0);
const filters = computed(() => response.value?.meta?.filters);

const getStatusDisplay = (status: EmployeeStatus) => {
  switch (status) {
    case EmployeeStatus.ACTIVE:
      return {
        text: 'Active',
        class: 'bg-green-100 text-green-800',
      };
    case EmployeeStatus.INACTIVE:
      return {
        text: 'Inactive',
        class: 'bg-red-100 text-red-800',
      };
    default:
      return {
        text: 'Unknown',
        class: 'bg-gray-100 text-gray-800',
      };
  }
};

const columnHelper = createColumnHelper<Employee>();
const columns = [
  columnHelper.accessor('id', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'ID',
    }),
    cell: (info) => h('div', { class: 'font-bold w-40' }, info.getValue()),
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
      const status = getStatusDisplay(info.getValue() as EmployeeStatus);
      return h('span', {
        class: ['inline-flex items-center rounded-full px-2 py-1 text-xs font-medium', status.class],
      }, status.text);
    },
    enableSorting: true,
  }),
  columnHelper.accessor('createdAt', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Created At',
    }),
    cell: (info) => {
      const seconds = BigInt(info.getValue()?.seconds ?? 0n);
      const millis = Number(seconds * 1000n);
      const date = new Date(millis); ;
      return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
    },
    enableSorting: true,
  }),
  columnHelper.accessor('updatedAt', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: 'Updated At',
    }),
    cell: (info) => {
      const seconds = BigInt(info.getValue()?.seconds ?? 0n);
      const millis = Number(seconds * 1000n);
      const date = new Date(millis); ;
      return `${date.toLocaleDateString()} ${date.toLocaleTimeString()}`;
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
  filters.value?.['roles']?.values?.map((role: string) => ({
    label: role,
    value: role,
  })) ?? [],
);

// department filter
const departmentFilter = useDataTableFilter(table, 'department');
const departmentOptions = computed(() =>
  filters.value?.['departments']?.values?.map((dept: string) => ({
    label: dept,
    value: dept,
  })) ?? [],
);

// status filter
const statusFilter = useDataTableFilter(table, 'status');
const statusOptions = computed(() =>
  filters.value?.['statuses']?.values?.map((status: string) => ({
    label: status === 'active' ? 'Active' : 'Inactive',
    value: status,
  })) ?? [],
);

// Reset all filters
function reset() {
  roleFilter.clearFilter();
  departmentFilter.clearFilter();
  statusFilter.clearFilter();
}
</script>

<template>
  <div class="space-y-5 px-4 py-3 sm:px-6 lg:px-8">
    <div class="container mx-auto flex justify-end">
      <Sheet>
        <SheetTrigger as-child>
          <Button size="sm">
            Add Employee
          </Button>
        </SheetTrigger>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Add new Employee</SheetTitle>
            <SheetDescription>
              <!-- Add better description here using better UX and wording -->
            </SheetDescription>
          </SheetHeader>
          <div>Body</div>
        </SheetContent>
      </Sheet>
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

          <DataTableFacetedFilter
            v-if="statusOptions.length > 0"
            v-model="statusFilter.filterValues.value"
            title="Status"
            :options="statusOptions"
            @update="statusFilter.setFilter"
            @clear="statusFilter.clearFilter"
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
                class="text-muted-foreground/50"
                size="2em"
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
              <Sheet>
                <SheetTrigger as-child>
                  <Button size="sm">
                    Add Employee
                  </Button>
                </SheetTrigger>
                <SheetContent>
                  <SheetHeader>
                    <SheetTitle>Add new Employee</SheetTitle>
                    <SheetDescription>
                      <!-- Add better description here using better UX and wording -->
                    </SheetDescription>
                  </SheetHeader>
                  <div>Body</div>
                </SheetContent>
              </Sheet>
            </div>
          </div>
        </template>
      </DataTable>
    </div>
  </div>
</template>
