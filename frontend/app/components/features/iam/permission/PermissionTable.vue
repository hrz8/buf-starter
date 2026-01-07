<script setup lang="ts">
import type { Permission } from '~~/gen/altalune/v1/permission_pb';
import { serializeProtoFilters } from '#shared/helpers/serializer';
import { createColumnHelper } from '@tanstack/vue-table';
import { toast } from 'vue-sonner';

import {
  DataTable,
  DataTableColumnHeader,
} from '@/components/custom/datatable';
import {
  useDataTableState,
} from '@/components/custom/datatable/utils';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { usePermissionService } from '@/composables/services/usePermissionService';
import { useQueryRequest } from '@/composables/useQueryRequest';

import PermissionCreateSheet from './PermissionCreateSheet.vue';
import PermissionDeleteDialog from './PermissionDeleteDialog.vue';
import PermissionRowActions from './PermissionRowActions.vue';
import PermissionUpdateSheet from './PermissionUpdateSheet.vue';

const { t } = useI18n();

// Services
const {
  query,
  resetCreateState,
} = usePermissionService();

// Table state
const page = ref(1);
const pageSize = ref(10);
const keyword = ref('');

const dataTableRef = ref<InstanceType<typeof DataTable> | null>(null);
const table = computed(() => dataTableRef.value?.table);

const { columnFilters, sorting } = useDataTableState(dataTableRef);

// Query request composable
const { queryRequest } = useQueryRequest({
  page,
  pageSize,
  keyword,
  columnFilters,
  sorting,
});

// Data fetching
const asyncDataKey = computed(() => {
  const { pagination, keyword, filters, sorting } = queryRequest.value;
  const { page, pageSize } = pagination!;
  const keys = [
    'permission-table',
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
  () => query({ query: queryRequest.value }),
  {
    server: false,
    watch: [queryRequest],
    immediate: false,
  },
);

const data = computed(() => response.value?.data ?? []);
const rowCount = computed(() => response.value?.meta?.rowCount ?? 0);

// Column helper
const columnHelper = createColumnHelper<Permission>();

// Format date utility
function formatDate(timestamp: any): string {
  if (!timestamp?.seconds)
    return t('features.permissions.status.unknown');
  const seconds = BigInt(timestamp.seconds);
  const millis = Number(seconds * 1000n);
  const date = new Date(millis);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

// Table columns (computed for reactivity with i18n)
const columns = computed(() => [
  columnHelper.accessor('name', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.permissions.columns.name'),
    }),
    cell: info => h('div', { class: 'font-medium' }, info.getValue()),
    enableSorting: true,
  }),
  columnHelper.accessor('description', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.permissions.columns.description'),
    }),
    cell: (info) => {
      const desc = info.getValue();
      return h('div', { class: 'text-sm text-muted-foreground max-w-md truncate' }, desc || '-');
    },
    enableSorting: false,
  }),
  columnHelper.accessor('createdAt', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.permissions.columns.createdAt'),
    }),
    cell: (info) => {
      const timestamp = info.getValue();
      return h('div', { class: 'text-sm text-muted-foreground' }, formatDate(timestamp));
    },
    enableSorting: true,
  }),
  columnHelper.display({
    id: 'actions',
    cell: ({ row }) => {
      return h(PermissionRowActions, {
        permission: row.original,
        onEdit: () => handleEdit(row),
        onDelete: () => handleDelete(row),
      });
    },
  }),
]);

// Create sheet state
const createdPermission = ref<Permission | null>(null);

// Edit and Delete sheet/dialog state
const selectedPermission = ref<Permission | null>(null);
const isEditSheetOpen = ref(false);
const isDeleteDialogOpen = ref(false);

// Event handlers
function handlePermissionCreated(permission: Permission) {
  resetCreateState();
  if (permission) {
    createdPermission.value = permission;
    toast.success(t('features.permissions.messages.createSuccess'), {
      description: t('features.permissions.messages.createSuccessDesc', { name: permission.name }),
    });
    refresh();
  }
}

function handleSheetClose() {
  resetCreateState();
}

// Handle row action events
function handleEdit(row: any) {
  selectedPermission.value = row.original as Permission;
  nextTick(() => {
    isEditSheetOpen.value = true;
  });
}

function handleDelete(row: any) {
  selectedPermission.value = row.original as Permission;
  nextTick(() => {
    isDeleteDialogOpen.value = true;
  });
}

// Handle sheet/dialog completion events
function handlePermissionUpdated(_permission: Permission) {
  isEditSheetOpen.value = false;
  selectedPermission.value = null;
  refresh();
}

function handlePermissionDeleted() {
  isDeleteDialogOpen.value = false;
  selectedPermission.value = null;
  refresh();
}

function closeEditSheet() {
  isEditSheetOpen.value = false;
  selectedPermission.value = null;
}

function closeDeleteDialog() {
  isDeleteDialogOpen.value = false;
  selectedPermission.value = null;
}

// Reset all filters
function reset() {
  // No filters currently active
}
</script>

<template>
  <div>
    <div class="space-y-5 px-4 py-3 sm:px-6 lg:px-8">
      <div class="container mx-auto">
        <h2 class="text-2xl font-bold">
          {{ t('features.permissions.page.title') }}
        </h2>
        <p class="text-muted-foreground">
          {{ t('features.permissions.page.description') }}
        </p>
      </div>
      <div class="container mx-auto flex justify-end">
        <PermissionCreateSheet
          @success="handlePermissionCreated"
          @cancel="handleSheetClose"
        >
          <Button size="sm">
            <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
            {{ t('features.permissions.actions.create') }}
          </Button>
        </PermissionCreateSheet>
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
          column-prefix="features.permissions.columns"
          @refresh="refresh()"
          @reset="reset()"
        >
          <template #filters>
            <Input
              v-model="keyword"
              :placeholder="t('features.permissions.actions.search')"
              class="h-8 w-[150px] lg:w-[250px]"
            />
          </template>
          <template #loading>
            <div class="flex items-center justify-center py-8">
              <div class="flex items-center space-x-2">
                <Icon name="lucide:loader-2" class="h-4 w-4 animate-spin" />
                <span class="text-sm text-muted-foreground">
                  {{ t('features.permissions.loading') }}
                </span>
              </div>
            </div>
          </template>
          <template #empty>
            <div class="flex flex-col items-center justify-center space-y-6 py-16">
              <div class="relative">
                <Icon
                  name="lucide:shield-check"
                  class="text-muted-foreground/50"
                  size="2em"
                />
              </div>
              <div class="text-center space-y-2">
                <h3 class="text-lg font-semibold">
                  {{ t('features.permissions.empty.title') }}
                </h3>
                <p class="text-center text-muted-foreground max-w-md whitespace-normal break-words">
                  {{ t('features.permissions.empty.description') }}
                </p>
              </div>
              <div class="flex space-x-2">
                <PermissionCreateSheet
                  @success="handlePermissionCreated"
                  @cancel="handleSheetClose"
                >
                  <Button>
                    <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
                    {{ t('features.permissions.actions.create') }}
                  </Button>
                </PermissionCreateSheet>
              </div>
            </div>
          </template>
        </DataTable>
      </div>
    </div>

    <!-- Permission Edit Sheet - Outside DataTable to avoid dropdown conflicts -->
    <PermissionUpdateSheet
      v-if="selectedPermission"
      v-model:open="isEditSheetOpen"
      :permission="selectedPermission"
      @success="handlePermissionUpdated"
      @cancel="closeEditSheet"
    />

    <!-- Permission Delete Dialog - Outside DataTable to avoid dropdown conflicts -->
    <PermissionDeleteDialog
      v-if="selectedPermission"
      v-model:open="isDeleteDialogOpen"
      :permission="selectedPermission"
      @success="handlePermissionDeleted"
      @cancel="closeDeleteDialog"
    />
  </div>
</template>
