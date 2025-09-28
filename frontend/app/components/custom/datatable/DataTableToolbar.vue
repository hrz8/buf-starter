<script setup lang="ts">
import type { Table } from '@tanstack/vue-table';
import type { Data } from '.';

import { DataTableViewOptions } from '@/components/custom/datatable';
import { Button } from '@/components/ui/button';

interface Props {
  table: Table<Data>;
}

interface Emits {
  refresh: [];
  reset: [];
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const isFiltered = computed(() => props.table.getState().columnFilters.length > 0);

function onRefresh() {
  emit('refresh');
}

function onReset() {
  props.table.resetColumnFilters();
  emit('reset');
}
</script>

<template>
  <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between mb-3">
    <div class="flex flex-1 flex-wrap items-center gap-2">
      <slot
        name="filters"
        :table="table"
      />
      <Button
        v-if="isFiltered"
        variant="ghost"
        size="sm"
        class="h-8 px-2 lg:px-3"
        @click="onReset"
      >
        Reset
        <Icon
          name="radix-icons:cross-2"
          class="ml-2 h-4 w-4"
        />
      </Button>
    </div>
    <div class="flex items-center gap-2">
      <Button
        variant="outline"
        size="sm"
        @click="onRefresh"
      >
        <Icon
          name="lucide:refresh-cw"
          class="h-4 w-4"
        />
      </Button>
      <DataTableViewOptions :table="table" />
    </div>
  </div>
</template>
