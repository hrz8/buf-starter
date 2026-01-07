<script setup lang="ts">
import type { ApiKey } from '~~/gen/altalune/v1/api_key_pb';
import { serializeProtoFilters } from '#shared/helpers/serializer';
import { createColumnHelper } from '@tanstack/vue-table';
import { toast } from 'vue-sonner';

import {
  DataTable,
  DataTableColumnHeader,
  DataTableFacetedFilter,
} from '@/components/custom/datatable';
import {
  useDataTableFilter,
  useDataTableState,
} from '@/components/custom/datatable/utils';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useApiKeyService } from '@/composables/services/useApiKeyService';
import { useQueryRequest } from '@/composables/useQueryRequest';

import ApiKeyCreateSheet from './ApiKeyCreateSheet.vue';
import ApiKeyDeleteDialog from './ApiKeyDeleteDialog.vue';
import ApiKeyDisplay from './ApiKeyDisplay.vue';
import ApiKeyRowActions from './ApiKeyRowActions.vue';
import ApiKeyUpdateSheet from './ApiKeyUpdateSheet.vue';

const props = defineProps<{
  projectId: string;
}>();

const { t } = useI18n();

// Services
const {
  query,
  resetCreateState,
  activateApiKey,
  deactivateApiKey,
} = useApiKeyService();

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
    'api-key-table',
    props.projectId,
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
    projectId: props.projectId,
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

// Column helper
const columnHelper = createColumnHelper<ApiKey>();

// Format date utility
function formatDate(timestamp: any): string {
  if (!timestamp?.seconds)
    return t('features.api_keys.status.unknown');
  const seconds = BigInt(timestamp.seconds);
  const millis = Number(seconds * 1000n);
  const date = new Date(millis);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

// Combined status utilities
function getCombinedStatus(apiKey: ApiKey): 'active' | 'inactive' | 'expired' | 'expiring_soon' {
  // If status is inactive, determine if expired or just inactive
  if (!apiKey.active) {
    if (!apiKey.expiration) {
      return 'inactive';
    }
    const seconds = BigInt(apiKey.expiration.seconds);
    const millis = Number(seconds * 1000n);
    const expirationDate = new Date(millis);

    // If expiration has passed, it's expired
    const now = new Date();
    if (expirationDate < now) {
      return 'expired';
    }
    // If expiration is still far away but status is inactive, treat as inactive
    // (This is unlikely but can happen, so we treat it as inactive)
    return 'inactive';
  }

  // If status is active, check expiration timing
  if (!apiKey.expiration) {
    return 'active';
  }

  const seconds = BigInt(apiKey.expiration.seconds);
  const millis = Number(seconds * 1000n);
  const expirationDate = new Date(millis);
  const tenDaysFromNow = new Date(Date.now() + 10 * 24 * 60 * 60 * 1000);

  // If expiration is within 10 days, it's expiring soon
  if (expirationDate <= tenDaysFromNow) {
    return 'expiring_soon';
  }

  // Otherwise it's active
  return 'active';
}

function getCombinedStatusDisplay(status: string) {
  switch (status) {
    case 'expired':
      return {
        text: t('features.api_keys.status.expired'),
        class: 'bg-red-100 text-red-800',
      };
    case 'expiring_soon':
      return {
        text: t('features.api_keys.status.expiringSoon'),
        class: 'bg-yellow-100 text-yellow-800',
      };
    case 'inactive':
      return {
        text: t('features.api_keys.status.inactive'),
        class: 'bg-gray-100 text-gray-800',
      };
    case 'active':
    default:
      return {
        text: t('features.api_keys.status.active'),
        class: 'bg-green-100 text-green-800',
      };
  }
}

// Table columns (computed for reactivity with i18n)
const columns = computed(() => [
  columnHelper.accessor('id', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.api_keys.columns.id'),
    }),
    cell: info => h('div', { class: 'font-bold w-40' }, info.getValue()),
    enableSorting: true,
  }),
  columnHelper.accessor('name', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.api_keys.columns.name'),
    }),
    cell: info => h('div', { class: 'font-medium' }, info.getValue()),
    enableSorting: true,
  }),
  columnHelper.accessor('expiration', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.api_keys.columns.expiration'),
    }),
    cell: (info) => {
      const timestamp = info.getValue();
      return h('div', { class: 'text-sm text-muted-foreground' }, formatDate(timestamp));
    },
    enableSorting: true,
  }),
  columnHelper.accessor('createdAt', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.api_keys.columns.createdAt'),
    }),
    cell: (info) => {
      const timestamp = info.getValue();
      return h('div', { class: 'text-sm text-muted-foreground' }, formatDate(timestamp));
    },
    enableSorting: true,
  }),
  columnHelper.display({
    id: 'status',
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.api_keys.columns.status'),
    }),
    cell: ({ row }) => {
      const apiKey = row.original;
      const status = getCombinedStatus(apiKey);
      const statusDisplay = getCombinedStatusDisplay(status);
      return h('span', {
        class: ['inline-flex items-center rounded-full px-2 py-1 text-xs font-medium', statusDisplay.class],
      }, statusDisplay.text);
    },
  }),
  columnHelper.display({
    id: 'actions',
    cell: ({ row }) => {
      return h(ApiKeyRowActions, {
        projectId: props.projectId,
        apiKey: row.original,
        onEdit: () => handleEdit(row),
        onDelete: () => handleDelete(row),
        onToggleStatus: () => handleToggleStatus(row),
      });
    },
  }),
]);

// Combined status filter - matches 'status' column ID
const statusFilter = useDataTableFilter(table, 'status');
const statusOptions = computed(() =>
  filters.value?.statuses?.values?.map((status: string) => ({
    label: status === 'active'
      ? t('features.api_keys.status.active')
      : status === 'inactive'
        ? t('features.api_keys.status.inactive')
        : status === 'expired'
          ? t('features.api_keys.status.expired')
          : status === 'expiring_soon'
            ? t('features.api_keys.status.expiringSoon')
            : status,
    value: status,
  })) ?? [],
);

// Create sheet state
const createdApiKey = ref<ApiKey | null>(null);
const createdKeyValue = ref('');
const isKeyDisplayOpen = ref(false);

// Edit and Delete sheet/dialog state
const selectedApiKey = ref<ApiKey | null>(null);
const isEditSheetOpen = ref(false);
const isDeleteDialogOpen = ref(false);

// Event handlers
function handleApiKeyCreated(result: { apiKey: ApiKey | null; keyValue: string }) {
  resetCreateState();
  if (result.apiKey && result.keyValue) {
    createdApiKey.value = result.apiKey;
    createdKeyValue.value = result.keyValue;
    isKeyDisplayOpen.value = true;
    refresh();
  }
}

function handleSheetClose() {
  resetCreateState();
}

// Handle row action events (following datatable18 pattern)
function handleEdit(row: any) {
  selectedApiKey.value = row.original as ApiKey;
  nextTick(() => {
    isEditSheetOpen.value = true;
  });
}

function handleDelete(row: any) {
  selectedApiKey.value = row.original as ApiKey;
  nextTick(() => {
    isDeleteDialogOpen.value = true;
  });
}

async function handleToggleStatus(row: any) {
  const apiKey = row.original as ApiKey;
  try {
    if (apiKey.active) {
      await deactivateApiKey({
        projectId: props.projectId,
        apiKeyId: apiKey.id,
      });
      toast.success(t('features.api_keys.messages.deactivateSuccess'));
    }
    else {
      await activateApiKey({
        projectId: props.projectId,
        apiKeyId: apiKey.id,
      });
      toast.success(t('features.api_keys.messages.activateSuccess'));
    }
    refresh();
  }
  catch (error) {
    console.error('Failed to toggle API key status:', error);
    toast.error(t('features.api_keys.messages.toggleError'));
  }
}

// Handle sheet/dialog completion events
function handleApiKeyUpdated(_apiKey: ApiKey) {
  isEditSheetOpen.value = false;
  selectedApiKey.value = null;
  refresh();
}

function handleApiKeyDeleted() {
  isDeleteDialogOpen.value = false;
  selectedApiKey.value = null;
  refresh();
}

function handleKeyDisplayClose() {
  isKeyDisplayOpen.value = false;
  createdApiKey.value = null;
  createdKeyValue.value = '';
}

function closeEditSheet() {
  isEditSheetOpen.value = false;
  selectedApiKey.value = null;
}

function closeDeleteDialog() {
  isDeleteDialogOpen.value = false;
  selectedApiKey.value = null;
}

// Reset all filters
function reset() {
  statusFilter.clearFilter();
}
</script>

<template>
  <div>
    <div class="space-y-5 px-4 py-3 sm:px-6 lg:px-8">
      <div class="container mx-auto">
        <h2 class="text-2xl font-bold">
          {{ t('features.api_keys.page.title') }}
        </h2>
        <p class="text-muted-foreground">
          {{ t('features.api_keys.page.description') }}
        </p>
      </div>
      <div class="container mx-auto flex justify-end">
        <ApiKeyCreateSheet
          :project-id="props.projectId"
          @success="handleApiKeyCreated"
          @cancel="handleSheetClose"
        >
          <Button size="sm">
            <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
            {{ t('features.api_keys.actions.create') }}
          </Button>
        </ApiKeyCreateSheet>
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
          column-prefix="features.api_keys.columns"
          @refresh="refresh()"
          @reset="reset()"
        >
          <template #filters>
            <Input
              v-model="keyword"
              :placeholder="t('features.api_keys.actions.search')"
              class="h-8 w-[150px] lg:w-[250px]"
            />

            <DataTableFacetedFilter
              v-if="statusOptions.length > 0"
              v-model="statusFilter.filterValues.value"
              :title="t('features.api_keys.filter.status')"
              :options="statusOptions"
              @update="statusFilter.setFilter"
              @clear="statusFilter.clearFilter"
            />
          </template>
          <template #loading>
            <div class="flex items-center justify-center py-8">
              <div class="flex items-center space-x-2">
                <Icon name="lucide:loader-2" class="h-4 w-4 animate-spin" />
                <span class="text-sm text-muted-foreground">
                  {{ t('features.api_keys.loading') }}
                </span>
              </div>
            </div>
          </template>
          <template #empty>
            <div class="flex flex-col items-center justify-center space-y-6 py-16">
              <div class="relative">
                <Icon
                  name="lucide:key"
                  class="text-muted-foreground/50"
                  size="2em"
                />
              </div>
              <div class="text-center space-y-2">
                <h3 class="text-lg font-semibold">
                  {{ t('features.api_keys.empty.title') }}
                </h3>
                <p class="text-center text-muted-foreground max-w-md whitespace-normal break-words">
                  {{ t('features.api_keys.empty.description') }}
                </p>
              </div>
              <div class="flex space-x-2">
                <ApiKeyCreateSheet
                  :project-id="props.projectId"
                  @success="handleApiKeyCreated"
                  @cancel="handleSheetClose"
                >
                  <Button>
                    <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
                    {{ t('features.api_keys.actions.create') }}
                  </Button>
                </ApiKeyCreateSheet>
              </div>
            </div>
          </template>
        </DataTable>
      </div>
    </div>

    <!-- API Key Edit Sheet - Outside DataTable to avoid dropdown conflicts -->
    <ApiKeyUpdateSheet
      v-if="selectedApiKey"
      v-model:open="isEditSheetOpen"
      :project-id="props.projectId"
      :api-key="selectedApiKey"
      @success="handleApiKeyUpdated"
      @cancel="closeEditSheet"
    />

    <!-- API Key Delete Dialog - Outside DataTable to avoid dropdown conflicts -->
    <ApiKeyDeleteDialog
      v-if="selectedApiKey"
      v-model:open="isDeleteDialogOpen"
      :project-id="props.projectId"
      :api-key="selectedApiKey"
      @success="handleApiKeyDeleted"
      @cancel="closeDeleteDialog"
    />

    <!-- Key Display Modal -->
    <ApiKeyDisplay
      v-model:open="isKeyDisplayOpen"
      :api-key="createdApiKey"
      :key-value="createdKeyValue"
      @close="handleKeyDisplayClose"
    />
  </div>
</template>
