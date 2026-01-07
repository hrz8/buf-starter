<script setup lang="ts">
import type { User } from '~~/gen/altalune/v1/user_pb';
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
import { useUserService } from '@/composables/services/useUserService';
import { useQueryRequest } from '@/composables/useQueryRequest';

import UserCreateSheet from './UserCreateSheet.vue';
import UserDeleteDialog from './UserDeleteDialog.vue';
import UserRowActions from './UserRowActions.vue';
import UserUpdateSheet from './UserUpdateSheet.vue';

const { t, locale } = useI18n();

// Services
const {
  query,
  resetCreateState,
} = useUserService();

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
    'user-table',
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
const columnHelper = createColumnHelper<User>();

// Format date utility
function formatDate(timestamp: any): string {
  if (!timestamp?.seconds)
    return t('features.users.status.unknown');
  const seconds = BigInt(timestamp.seconds);
  const millis = Number(seconds * 1000n);
  const date = new Date(millis);
  return date.toLocaleDateString(locale.value, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

// Table columns (computed for reactivity with i18n)
const columns = computed(() => [
  columnHelper.accessor('email', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.users.columns.email'),
    }),
    cell: info => h('div', { class: 'font-medium' }, info.getValue()),
    enableSorting: true,
  }),
  columnHelper.accessor('firstName', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.users.columns.firstName'),
    }),
    cell: info => h('div', { class: 'text-sm' }, info.getValue()),
    enableSorting: true,
  }),
  columnHelper.accessor('lastName', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.users.columns.lastName'),
    }),
    cell: info => h('div', { class: 'text-sm' }, info.getValue()),
    enableSorting: true,
  }),
  columnHelper.accessor('createdAt', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.users.columns.createdAt'),
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
      return h(UserRowActions, {
        user: row.original,
        onEdit: () => handleEdit(row),
        onDelete: () => handleDelete(row),
      });
    },
  }),
]);

// Create sheet state
const createdUser = ref<User | null>(null);

// Edit and Delete sheet/dialog state
const selectedUser = ref<User | null>(null);
const isEditSheetOpen = ref(false);
const isDeleteDialogOpen = ref(false);

// Event handlers
function handleUserCreated(user: User) {
  resetCreateState();
  if (user) {
    createdUser.value = user;
    toast.success(t('features.users.messages.createSuccess'), {
      description: t('features.users.messages.createSuccessDesc', {
        name: `${user.firstName} ${user.lastName}`,
      }),
    });
    refresh();
  }
}

function handleSheetClose() {
  resetCreateState();
}

// Handle row action events
function handleEdit(row: any) {
  selectedUser.value = row.original as User;
  nextTick(() => {
    isEditSheetOpen.value = true;
  });
}

function handleDelete(row: any) {
  selectedUser.value = row.original as User;
  nextTick(() => {
    isDeleteDialogOpen.value = true;
  });
}

// Handle sheet/dialog completion events
function handleUserUpdated(_user: User) {
  isEditSheetOpen.value = false;
  selectedUser.value = null;
  refresh();
}

function handleUserDeleted() {
  isDeleteDialogOpen.value = false;
  selectedUser.value = null;
  refresh();
}

function closeEditSheet() {
  isEditSheetOpen.value = false;
  selectedUser.value = null;
}

function closeDeleteDialog() {
  isDeleteDialogOpen.value = false;
  selectedUser.value = null;
}

// Reset all filters
function reset() {
  // No filters for users yet
}
</script>

<template>
  <div>
    <div class="space-y-5 px-4 py-3 sm:px-6 lg:px-8">
      <div class="container mx-auto">
        <h2 class="text-2xl font-bold">
          {{ t('features.users.page.title') }}
        </h2>
        <p class="text-muted-foreground">
          {{ t('features.users.page.description') }}
        </p>
      </div>
      <div class="container mx-auto flex justify-end">
        <UserCreateSheet
          @success="handleUserCreated"
          @cancel="handleSheetClose"
        >
          <Button size="sm">
            <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
            {{ t('features.users.actions.create') }}
          </Button>
        </UserCreateSheet>
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
          column-prefix="features.users.columns"
          @refresh="refresh()"
          @reset="reset()"
        >
          <template #filters>
            <Input
              v-model="keyword"
              :placeholder="t('features.users.actions.search')"
              class="h-8 w-[150px] lg:w-[250px]"
            />
          </template>
          <template #loading>
            <div class="flex items-center justify-center py-8">
              <div class="flex items-center space-x-2">
                <Icon name="lucide:loader-2" class="h-4 w-4 animate-spin" />
                <span class="text-sm text-muted-foreground">
                  {{ t('features.users.loading') }}
                </span>
              </div>
            </div>
          </template>
          <template #empty>
            <div class="flex flex-col items-center justify-center space-y-6 py-16">
              <div class="relative">
                <Icon
                  name="lucide:users"
                  class="text-muted-foreground/50"
                  size="2em"
                />
              </div>
              <div class="text-center space-y-2">
                <h3 class="text-lg font-semibold">
                  {{ t('features.users.empty.title') }}
                </h3>
                <p class="text-center text-muted-foreground max-w-md whitespace-normal break-words">
                  {{ t('features.users.empty.description') }}
                </p>
              </div>
              <div class="flex space-x-2">
                <UserCreateSheet
                  @success="handleUserCreated"
                  @cancel="handleSheetClose"
                >
                  <Button>
                    <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
                    {{ t('features.users.actions.create') }}
                  </Button>
                </UserCreateSheet>
              </div>
            </div>
          </template>
        </DataTable>
      </div>
    </div>

    <!-- User Edit Sheet - Outside DataTable to avoid dropdown conflicts -->
    <UserUpdateSheet
      v-if="selectedUser"
      v-model:open="isEditSheetOpen"
      :user="selectedUser"
      @success="handleUserUpdated"
      @cancel="closeEditSheet"
    />

    <!-- User Delete Dialog - Outside DataTable to avoid dropdown conflicts -->
    <UserDeleteDialog
      v-if="selectedUser"
      v-model:open="isDeleteDialogOpen"
      :user="selectedUser"
      @success="handleUserDeleted"
      @cancel="closeDeleteDialog"
    />
  </div>
</template>
