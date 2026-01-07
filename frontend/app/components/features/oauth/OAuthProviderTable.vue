<script setup lang="ts">
import type { OAuthProvider } from '~~/gen/altalune/v1/oauth_provider_pb';
import { serializeProtoFilters } from '#shared/helpers/serializer';
import { createColumnHelper } from '@tanstack/vue-table';

import {
  DataTable,
  DataTableColumnHeader,
  DataTableFacetedFilter,
} from '@/components/custom/datatable';
import {
  useDataTableFilter,
  useDataTableState,
} from '@/components/custom/datatable/utils';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useOAuthProviderService } from '@/composables/services/useOAuthProviderService';
import { useQueryRequest } from '@/composables/useQueryRequest';

import { getProviderMetadata } from './constants';
import OAuthProviderCreateSheet from './OAuthProviderCreateSheet.vue';
import OAuthProviderDeleteDialog from './OAuthProviderDeleteDialog.vue';
import OAuthProviderRowActions from './OAuthProviderRowActions.vue';
import OAuthProviderUpdateSheet from './OAuthProviderUpdateSheet.vue';

const { t, locale } = useI18n();

// Services
const { query, resetCreateState } = useOAuthProviderService();

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
    'oauth-provider-table',
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
const columnHelper = createColumnHelper<OAuthProvider>();

// Format date utility (handles Protobuf Timestamp)
function formatDate(timestamp: any): string {
  if (!timestamp?.seconds)
    return t('features.oauth.table.noScopes');
  const seconds = BigInt(timestamp.seconds);
  const millis = Number(seconds * 1000n);
  const date = new Date(millis);
  return date.toLocaleDateString(locale.value, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

// Provider type filter
const providerTypeFilter = useDataTableFilter(table, 'providerType');
const providerTypeOptions = computed(() => {
  const options = filters.value?.providerType?.values || [];
  return options.map((value: string) => {
    const numValue = Number.parseInt(value, 10);
    const metadata = getProviderMetadata(numValue);
    return {
      label: metadata?.name || value,
      value,
    };
  });
});

// Enabled filter
const enabledFilter = useDataTableFilter(table, 'enabled');
const enabledOptions = computed(() => {
  const options = filters.value?.enabled?.values || [];
  return options.map((value: string) => ({
    label: value === 'true'
      ? t('features.oauth.status.enabled')
      : t('features.oauth.status.disabled'),
    value,
  }));
});

// Define columns
const columns = computed(() => [
  // Provider Type column with icon
  columnHelper.accessor('providerType', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.oauth.table.providerType'),
    }),
    cell: ({ row }) => {
      const provider = row.original;
      const metadata = getProviderMetadata(provider.providerType);

      return h('div', { class: 'flex items-center gap-2' }, [
        metadata?.icon
          ? h(resolveComponent('Icon'), {
              name: metadata.icon,
              class: 'h-4 w-4',
            })
          : null,
        h('span', { class: 'font-medium' }, metadata?.name || 'Unknown'),
      ]);
    },
    enableSorting: false,
    enableHiding: false,
  }),

  // Client ID column
  columnHelper.accessor('clientId', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.oauth.table.clientId'),
    }),
    cell: ({ row }) => {
      const clientId = row.original.clientId;
      return h('code', { class: 'text-xs font-mono' }, clientId);
    },
    enableSorting: false,
  }),

  // Redirect URL column (truncated)
  columnHelper.accessor('redirectUrl', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.oauth.table.redirectUrl'),
    }),
    cell: ({ row }) => {
      const url = row.original.redirectUrl;
      const truncated = url.length > 40 ? `${url.slice(0, 40)}...` : url;
      return h('span', { class: 'text-sm', title: url }, truncated);
    },
    enableSorting: false,
  }),

  // Scopes column (badges)
  columnHelper.accessor('scopes', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.oauth.table.scopes'),
    }),
    cell: ({ row }) => {
      const scopes = row.original.scopes.split(',').map(s => s.trim()).filter(Boolean);
      return h(
        'div',
        { class: 'flex flex-wrap gap-1' },
        scopes.slice(0, 3).map(scope =>
          h(Badge, { variant: 'outline', class: 'text-xs' }, () => scope),
        ).concat(
          scopes.length > 3
            ? [h(Badge, { variant: 'secondary', class: 'text-xs' }, () => `+${scopes.length - 3}`)]
            : [],
        ),
      );
    },
    enableSorting: false,
  }),

  // Status column (Enabled/Disabled badge)
  columnHelper.accessor('enabled', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.oauth.table.status'),
    }),
    cell: ({ row }) => {
      const enabled = row.original.enabled;
      const variant = enabled ? 'success' : 'destructive';
      return h(
        Badge,
        { variant },
        () => enabled ? t('features.oauth.status.enabled') : t('features.oauth.status.disabled'),
      );
    },
    enableSorting: false,
  }),

  // Created At column (sortable)
  columnHelper.accessor('createdAt', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.oauth.table.createdAt'),
    }),
    cell: ({ row }) => {
      const date = row.original.createdAt;
      return h('span', { class: 'text-sm text-muted-foreground' }, formatDate(date));
    },
    enableSorting: true,
  }),

  // Actions column
  columnHelper.display({
    id: 'actions',
    header: () => h('span', { class: 'sr-only' }, t('features.oauth.table.actions')),
    cell: ({ row }) => {
      const provider = row.original;
      return h(OAuthProviderRowActions, {
        provider,
        onEdit: () => handleEdit(provider),
        onDelete: () => handleDelete(provider),
      });
    },
    enableHiding: false,
  }),
]);

// Sheet and dialog state management
const createSheetOpen = ref(false);
const updateSheetOpen = ref(false);
const deleteDialogOpen = ref(false);
const selectedProvider = ref<OAuthProvider | null>(null);

// Row action handlers
function handleEdit(provider: OAuthProvider) {
  selectedProvider.value = provider;
  updateSheetOpen.value = true;
}

function handleDelete(provider: OAuthProvider) {
  selectedProvider.value = provider;
  deleteDialogOpen.value = true;
}

// Refresh handlers
async function handleProviderCreated() {
  createSheetOpen.value = false;
  await refresh();
}

async function handleProviderUpdated() {
  updateSheetOpen.value = false;
  selectedProvider.value = null;
  await refresh();
}

async function handleProviderDeleted() {
  deleteDialogOpen.value = false;
  selectedProvider.value = null;
  await refresh();
}

function handleUpdateCancel() {
  updateSheetOpen.value = false;
  selectedProvider.value = null;
}

function handleDeleteCancel() {
  deleteDialogOpen.value = false;
  selectedProvider.value = null;
}

function handleCreateCancel() {
  createSheetOpen.value = false;
  resetCreateState();
}

// Reset all filters
function reset() {
  providerTypeFilter.clearFilter();
  enabledFilter.clearFilter();
}
</script>

<template>
  <div>
    <div class="space-y-5 px-4 py-3 sm:px-6 lg:px-8">
      <div class="container mx-auto">
        <h2 class="text-2xl font-bold">
          {{ t('features.oauth.page.title') }}
        </h2>
        <p class="text-muted-foreground">
          {{ t('features.oauth.page.description') }}
        </p>
      </div>
      <div class="container mx-auto flex justify-end">
        <Button
          size="sm"
          @click="createSheetOpen = true"
        >
          <Icon
            name="lucide:plus"
            class="mr-2 h-4 w-4"
          />
          {{ t('features.oauth.actions.create') }}
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
          column-prefix="features.oauth.table"
          @refresh="refresh()"
          @reset="reset()"
        >
          <template #filters>
            <Input
              v-model="keyword"
              :placeholder="t('features.oauth.actions.search')"
              class="h-8 w-[150px] lg:w-[250px]"
            />

            <DataTableFacetedFilter
              v-if="providerTypeOptions.length > 0"
              v-model="providerTypeFilter.filterValues.value"
              :title="t('features.oauth.filters.providerType')"
              :options="providerTypeOptions"
              @update="providerTypeFilter.setFilter"
              @clear="providerTypeFilter.clearFilter"
            />

            <DataTableFacetedFilter
              v-if="enabledOptions.length > 0"
              v-model="enabledFilter.filterValues.value"
              :title="t('features.oauth.filters.enabled')"
              :options="enabledOptions"
              @update="enabledFilter.setFilter"
              @clear="enabledFilter.clearFilter"
            />
          </template>
          <template #loading>
            <div class="flex items-center justify-center py-8">
              <div class="flex items-center space-x-2">
                <Icon
                  name="lucide:loader-2"
                  class="h-4 w-4 animate-spin"
                />
                <span class="text-sm text-muted-foreground">
                  {{ t('features.oauth.loading') }}
                </span>
              </div>
            </div>
          </template>
          <template #empty>
            <div class="flex flex-col items-center justify-center space-y-6 py-16">
              <Icon
                name="lucide:key-round"
                class="h-16 w-16 text-muted-foreground/30"
              />
              <div class="space-y-2 text-center">
                <h3 class="text-lg font-semibold">
                  {{ t('features.oauth.empty.title') }}
                </h3>
                <p class="text-sm text-muted-foreground max-w-md">
                  {{ t('features.oauth.empty.description') }}
                </p>
              </div>
              <Button
                size="sm"
                @click="createSheetOpen = true"
              >
                <Icon
                  name="lucide:plus"
                  class="mr-2 h-4 w-4"
                />
                {{ t('features.oauth.actions.createFirst') }}
              </Button>
            </div>
          </template>
        </DataTable>
      </div>
    </div>

    <!-- Create sheet -->
    <OAuthProviderCreateSheet
      v-model:open="createSheetOpen"
      @success="handleProviderCreated"
      @cancel="handleCreateCancel"
    />

    <!-- Update sheet -->
    <OAuthProviderUpdateSheet
      v-if="selectedProvider"
      v-model:open="updateSheetOpen"
      :provider="selectedProvider"
      @success="handleProviderUpdated"
      @cancel="handleUpdateCancel"
    />

    <!-- Delete dialog -->
    <OAuthProviderDeleteDialog
      v-if="selectedProvider"
      v-model:open="deleteDialogOpen"
      :provider="selectedProvider"
      @success="handleProviderDeleted"
      @cancel="handleDeleteCancel"
    />
  </div>
</template>
