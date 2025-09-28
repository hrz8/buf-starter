<script setup lang="ts">
import type { ColumnDef, Table as TanstackTable } from '@tanstack/vue-table';
import {
  FlexRender,
} from '@tanstack/vue-table';

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

interface DataTableContentProps {
  columns: ColumnDef<any, any>[];
  table: TanstackTable<any>;
  pending: boolean;
}
const props = defineProps<DataTableContentProps>();
</script>

<template>
  <div class="rounded-md border">
    <Table>
      <TableHeader>
        <TableRow
          v-for="headerGroup in table.getHeaderGroups()"
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
          v-if="props.pending"
          :colspan="table.getVisibleFlatColumns().length"
        >
          <div class="flex items-center justify-center space-x-2">
            <Icon
              name="lucide:loader-2"
              class="w-4 h-4 animate-spin"
            />
            <span>Loading...</span>
          </div>
        </TableEmpty>

        <template v-else-if="table.getRowModel().rows?.length">
          <TableRow
            v-for="row in table.getRowModel().rows"
            :key="row.id"
            :data-state="row.getIsSelected() && 'selected'"
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
        </template>

        <TableEmpty
          v-else
          :colspan="table.getVisibleFlatColumns().length"
        >
          <div class="flex flex-col items-center justify-center space-y-2 text-muted-foreground">
            <Icon
              name="lucide:users"
              class="w-8 h-8"
            />
            <span>No results</span>
            <span class="text-sm">Try adjusting your search or filters</span>
          </div>
        </TableEmpty>
      </TableBody>
    </Table>
  </div>
</template>
