<script setup lang="ts">
import type { Role } from '~~/gen/altalune/v1/role_pb';
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
import { useRoleService } from '@/composables/services/useRoleService';
import { useQueryRequest } from '@/composables/useQueryRequest';

import RoleCreateSheet from './RoleCreateSheet.vue';
import RoleDeleteDialog from './RoleDeleteDialog.vue';
import RoleRowActions from './RoleRowActions.vue';
import RoleUpdateSheet from './RoleUpdateSheet.vue';

const { t } = useI18n();

// Services
const {
  query,
  resetCreateState,
} = useRoleService();

// Table state
const page = ref(1);
const pageSize = ref(10);
const keyword = ref('');

const dataTableRef = ref<InstanceType<typeof DataTable> | null>(null);

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
    'role-table',
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
const columnHelper = createColumnHelper<Role>();

// Format date utility
function formatDate(timestamp: any): string {
  if (!timestamp?.seconds)
    return t('features.roles.status.unknown');
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
      title: t('features.roles.columns.name'),
    }),
    cell: info => h('div', { class: 'font-medium' }, info.getValue()),
    enableSorting: true,
  }),
  columnHelper.accessor('description', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.roles.columns.description'),
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
      title: t('features.roles.columns.createdAt'),
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
      return h(RoleRowActions, {
        role: row.original,
        onEdit: () => handleEdit(row),
        onDelete: () => handleDelete(row),
      });
    },
  }),
]);

// Create sheet state
const createdRole = ref<Role | null>(null);

// Edit and Delete sheet/dialog state
const selectedRole = ref<Role | null>(null);
const isEditSheetOpen = ref(false);
const isDeleteDialogOpen = ref(false);

// Event handlers
function handleRoleCreated(role: Role) {
  resetCreateState();
  if (role) {
    createdRole.value = role;
    toast.success(t('features.roles.messages.createSuccess'), {
      description: t('features.roles.messages.createSuccessDesc', { name: role.name }),
    });
    refresh();
  }
}

function handleSheetClose() {
  resetCreateState();
}

// Handle row action events
function handleEdit(row: any) {
  selectedRole.value = row.original as Role;
  nextTick(() => {
    isEditSheetOpen.value = true;
  });
}

function handleDelete(row: any) {
  selectedRole.value = row.original as Role;
  nextTick(() => {
    isDeleteDialogOpen.value = true;
  });
}

// Handle sheet/dialog completion events
function handleRoleUpdated(_role: Role) {
  isEditSheetOpen.value = false;
  selectedRole.value = null;
  refresh();
}

function handleRoleDeleted() {
  isDeleteDialogOpen.value = false;
  selectedRole.value = null;
  refresh();
}

function closeEditSheet() {
  isEditSheetOpen.value = false;
  selectedRole.value = null;
}

function closeDeleteDialog() {
  isDeleteDialogOpen.value = false;
  selectedRole.value = null;
}

// Reset all filters
function reset() {
  // No filters for roles yet
}
</script>

<template>
  <div>
    <div class="space-y-5 px-4 py-3 sm:px-6 lg:px-8">
      <div class="container mx-auto flex justify-end">
        <RoleCreateSheet
          @success="handleRoleCreated"
          @cancel="handleSheetClose"
        >
          <Button size="sm">
            <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
            {{ t('features.roles.actions.create') }}
          </Button>
        </RoleCreateSheet>
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
          column-prefix="features.roles.columns"
          @refresh="refresh()"
          @reset="reset()"
        >
          <template #filters>
            <Input
              v-model="keyword"
              :placeholder="t('features.roles.actions.search')"
              class="h-8 w-[150px] lg:w-[250px]"
            />
          </template>
          <template #loading>
            <div class="flex items-center justify-center py-8">
              <div class="flex items-center space-x-2">
                <Icon name="lucide:loader-2" class="h-4 w-4 animate-spin" />
                <span class="text-sm text-muted-foreground">
                  {{ t('features.roles.loading') }}
                </span>
              </div>
            </div>
          </template>
          <template #empty>
            <div class="flex flex-col items-center justify-center space-y-6 py-16">
              <div class="relative">
                <Icon
                  name="lucide:user-cog"
                  class="text-muted-foreground/50"
                  size="2em"
                />
              </div>
              <div class="text-center space-y-2">
                <h3 class="text-lg font-semibold">
                  {{ t('features.roles.empty.title') }}
                </h3>
                <p class="text-center text-muted-foreground max-w-md whitespace-normal break-words">
                  {{ t('features.roles.empty.description') }}
                </p>
              </div>
              <div class="flex space-x-2">
                <RoleCreateSheet
                  @success="handleRoleCreated"
                  @cancel="handleSheetClose"
                >
                  <Button>
                    <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
                    {{ t('features.roles.actions.create') }}
                  </Button>
                </RoleCreateSheet>
              </div>
            </div>
          </template>
        </DataTable>
      </div>
    </div>

    <!-- Role Edit Sheet - Outside DataTable to avoid dropdown conflicts -->
    <RoleUpdateSheet
      v-if="selectedRole"
      v-model:open="isEditSheetOpen"
      :role="selectedRole"
      @success="handleRoleUpdated"
      @cancel="closeEditSheet"
    />

    <!-- Role Delete Dialog - Outside DataTable to avoid dropdown conflicts -->
    <RoleDeleteDialog
      v-if="selectedRole"
      v-model:open="isDeleteDialogOpen"
      :role="selectedRole"
      @success="handleRoleDeleted"
      @cancel="closeDeleteDialog"
    />
  </div>
</template>
