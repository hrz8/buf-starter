<script setup lang="ts">
import type { OAuthClient } from '~~/gen/altalune/v1/oauth_client_pb';
import { serializeProtoFilters } from '#shared/helpers/serializer';
import { createColumnHelper } from '@tanstack/vue-table';

import {
  DataTable,
  DataTableColumnHeader,
} from '@/components/custom/datatable';
import {
  useDataTableState,
} from '@/components/custom/datatable/utils';
import OAuthClientCreateSheet from '@/components/features/oauth-client/OAuthClientCreateSheet.vue';
import OAuthClientDeleteDialog from '@/components/features/oauth-client/OAuthClientDeleteDialog.vue';
import OAuthClientEditSheet from '@/components/features/oauth-client/OAuthClientEditSheet.vue';
import OAuthClientRevealDialog from '@/components/features/oauth-client/OAuthClientRevealDialog.vue';
import OAuthClientRowActions from '@/components/features/oauth-client/OAuthClientRowActions.vue';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useOAuthClientService } from '@/composables/services/useOAuthClientService';
import { useQueryRequest } from '@/composables/useQueryRequest';

// OAuth clients are GLOBAL entities (not project-scoped)
// Following Keycloak/Auth0 architecture patterns

const { t } = useI18n();
const { locale } = useI18n();

// Services
const { query } = useOAuthClientService();

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
    'oauth-client-table',
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

// Column helper
const columnHelper = createColumnHelper<OAuthClient>();

// Format date utility
function formatDate(timestamp: any): string {
  if (!timestamp?.seconds)
    return 'N/A';
  const seconds = BigInt(timestamp.seconds);
  const millis = Number(seconds * 1000n);
  const date = new Date(millis);
  return date.toLocaleDateString(locale.value, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

// Mask client ID
function maskClientId(clientId: string): string {
  if (!clientId || clientId.length < 4)
    return '****';
  return `****-****-****-${clientId.slice(-4)}`;
}

// Sheet/Dialog state
const selectedClient = ref<OAuthClient | null>(null);
const isEditSheetOpen = ref(false);
const isDeleteDialogOpen = ref(false);
const isRevealDialogOpen = ref(false);

// Handle create success
function handleClientCreated() {
  refresh();
}

function handleSheetClose() {
  // Sheet closed without creating
}

// Row action handlers
function handleEdit(row: any) {
  selectedClient.value = row.original as OAuthClient;
  nextTick(() => {
    isEditSheetOpen.value = true;
  });
}

function handleRevealSecret(row: any) {
  selectedClient.value = row.original as OAuthClient;
  nextTick(() => {
    isRevealDialogOpen.value = true;
  });
}

function handleDelete(row: any) {
  selectedClient.value = row.original as OAuthClient;
  nextTick(() => {
    isDeleteDialogOpen.value = true;
  });
}

// Handle edit success
function handleEditSuccess() {
  isEditSheetOpen.value = false;
  refresh();
}

// Handle delete success
function handleDeleteSuccess() {
  isDeleteDialogOpen.value = false;
  refresh();
}

// Reset filters (for future use)
function reset() {
  // Reset any filters here when implemented
}

// Columns definition
const columns = [
  columnHelper.accessor('name', {
    header: ({ column }) => h(DataTableColumnHeader, { column, title: 'Client Name' }),
    cell: ({ row }) => {
      const client = row.original;
      return h('div', { class: 'flex items-center gap-2' }, [
        h('span', { class: 'font-medium' }, client.name),
        client.isDefault && h(Badge, { variant: 'secondary', class: 'text-xs' }, 'Default'),
      ]);
    },
  }),
  columnHelper.accessor('clientId', {
    header: ({ column }) => h(DataTableColumnHeader, { column, title: 'Client ID' }),
    cell: ({ row }) => h('code', { class: 'text-sm text-muted-foreground' }, maskClientId(row.original.clientId)),
  }),
  columnHelper.accessor('redirectUris', {
    header: ({ column }) => h(DataTableColumnHeader, { column, title: 'Redirect URIs' }),
    cell: ({ row }) => h('span', { class: 'text-sm' }, `${row.original.redirectUris.length} URI(s)`),
  }),
  columnHelper.accessor('pkceRequired', {
    header: ({ column }) => h(DataTableColumnHeader, { column, title: 'PKCE' }),
    cell: ({ row }) => {
      return row.original.pkceRequired
        ? h(Badge, { variant: 'outline', class: 'text-xs' }, 'Required')
        : h('span', { class: 'text-sm text-muted-foreground' }, 'Optional');
    },
  }),
  columnHelper.accessor('createdAt', {
    header: ({ column }) => h(DataTableColumnHeader, { column, title: 'Created' }),
    cell: ({ row }) => h('span', { class: 'text-sm text-muted-foreground' }, formatDate(row.original.createdAt)),
  }),
  columnHelper.display({
    id: 'actions',
    cell: ({ row }) => {
      return h(OAuthClientRowActions, {
        client: row.original,
        onEdit: () => handleEdit(row),
        onRevealSecret: () => handleRevealSecret(row),
        onDelete: () => handleDelete(row),
      });
    },
  }),
];

defineExpose({
  refresh,
});
</script>

<template>
  <div>
    <div class="space-y-5 px-4 py-3 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="container mx-auto">
        <h2 class="text-2xl font-bold">
          {{ t('features.oauth_clients.page.title') }}
        </h2>
        <p class="text-muted-foreground">
          {{ t('features.oauth_clients.page.description') }}
        </p>
      </div>

      <!-- Create Button -->
      <div class="container mx-auto flex justify-end">
        <OAuthClientCreateSheet
          @success="handleClientCreated"
          @cancel="handleSheetClose"
        >
          <Button size="sm">
            <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
            {{ t('features.oauth_clients.actions.create') }}
          </Button>
        </OAuthClientCreateSheet>
      </div>

      <!-- Data Table -->
      <div class="container mx-auto">
        <DataTable
          ref="dataTableRef"
          v-model:page="page"
          v-model:page-size="pageSize"
          :columns="columns"
          :data="data"
          :pending="pending"
          :row-count="rowCount"
          column-prefix="features.oauth_clients.columns"
          @refresh="refresh()"
          @reset="reset()"
        >
          <template #filters>
            <Input
              v-model="keyword"
              :placeholder="t('features.oauth_clients.actions.search')"
              class="h-8 w-[150px] lg:w-[250px]"
            />
          </template>

          <template #loading>
            <div class="flex items-center justify-center py-8">
              <div class="flex items-center space-x-2">
                <Icon name="lucide:loader-2" class="h-4 w-4 animate-spin" />
                <span class="text-sm text-muted-foreground">
                  {{ t('features.oauth_clients.loading') }}
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
                  {{ t('features.oauth_clients.empty.title') }}
                </h3>
                <p class="text-center text-muted-foreground max-w-md whitespace-normal break-words">
                  {{ t('features.oauth_clients.empty.description') }}
                </p>
              </div>
              <div class="flex space-x-2">
                <OAuthClientCreateSheet
                  @success="handleClientCreated"
                  @cancel="handleSheetClose"
                >
                  <Button>
                    <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
                    {{ t('features.oauth_clients.actions.create') }}
                  </Button>
                </OAuthClientCreateSheet>
              </div>
            </div>
          </template>
        </DataTable>
      </div>
    </div>

    <!-- Edit Sheet - Outside DataTable to avoid dropdown conflicts -->
    <OAuthClientEditSheet
      v-if="selectedClient"
      v-model:open="isEditSheetOpen"
      :client="selectedClient"
      @success="handleEditSuccess"
    />

    <!-- Delete Dialog -->
    <OAuthClientDeleteDialog
      v-if="selectedClient"
      v-model:open="isDeleteDialogOpen"
      :client="selectedClient"
      @success="handleDeleteSuccess"
    />

    <!-- Reveal Dialog -->
    <OAuthClientRevealDialog
      v-if="selectedClient"
      v-model:open="isRevealDialogOpen"
      :client="selectedClient"
    />
  </div>
</template>
