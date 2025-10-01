<script setup lang="ts">
import type { Table as TanstackTable } from '@tanstack/vue-table';

import { FlexRender } from '@tanstack/vue-table';

import {
  Table,
  TableBody,
  TableCell,
  TableEmpty,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

const props = withDefaults(defineProps<Props>(), {
  pending: false,
  showBorder: true,
});

defineSlots<{
  loading?: () => any;
  empty?: () => any;
}>();

const { t } = useI18n();

interface Props {
  table: TanstackTable<any>;
  pending?: boolean;
  showBorder?: boolean;
}

const headerGroups = computed(() => props.table.getHeaderGroups());
const rows = computed(() => props.table.getRowModel().rows);
const visibleColumns = computed(() => props.table.getVisibleFlatColumns());
const hasData = computed(() => rows.value.length > 0);
</script>

<template>
  <div :class="[showBorder ? 'rounded-md border mb-3' : '']">
    <Table>
      <TableHeader>
        <TableRow
          v-for="headerGroup in headerGroups"
          :key="headerGroup.id"
        >
          <TableHead
            v-for="header in headerGroup.headers"
            :key="header.id"
          >
            <FlexRender
              v-if="!header.isPlaceholder"
              :render="header.column.columnDef.header"
              :props="header.getContext()"
            />
          </TableHead>
        </TableRow>
      </TableHeader>

      <TableBody>
        <TableEmpty
          v-if="pending"
          :colspan="visibleColumns.length"
        >
          <slot
            name="loading"
            :colspan="visibleColumns.length"
          >
            <div class="flex items-center justify-center space-x-2 py-8">
              <Icon
                name="lucide:loader-2"
                class="w-4 h-4 animate-spin"
              />
              <span>{{ t('common.status.loading') }}</span>
            </div>
          </slot>
        </TableEmpty>

        <TableEmpty
          v-else-if="!hasData"
          :colspan="visibleColumns.length"
        >
          <slot
            name="empty"
            :colspan="visibleColumns.length"
          >
            <div
              class="
                flex flex-col items-center justify-center space-y-2 text-muted-foreground py-8
              "
            >
              <Icon
                name="lucide:search-x"
                class="w-8 h-8"
              />
              <span class="font-medium">{{ t('datatable.noDataFound') }}</span>
              <span class="text-sm">{{ t('datatable.tryAdjusting') }}</span>
            </div>
          </slot>
        </TableEmpty>

        <TableRow
          v-for="row in rows"
          v-else
          :key="row.id"
          class="group"
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
</template>
