<script setup lang="ts">
import type { Table } from '@tanstack/vue-table';

import DataTableViewOptions from './DataTableViewOptions.vue';

import { Button } from '@/components/ui/button';

interface DataTableToolbarProps {
  table: Table<any>;
}

const props = defineProps<DataTableToolbarProps>();
const emit = defineEmits<{
  refresh: [];
  reset: [];
}>();

const isFiltered = computed(() => props.table.getState().columnFilters.length > 0);

function onClickRefresh() {
  emit('refresh');
}

function onClickReset() {
  props.table.resetColumnFilters();
  emit('reset');
}
</script>

<template>
  <div class="flex items-center justify-between">
    <div class="flex flex-1 items-center space-x-2">
      <slot
        name="filters"
        :table="table"
      />
      <Button
        v-if="isFiltered"
        variant="ghost"
        class="h-8 px-2 lg:px-3"
        @click="onClickReset"
      >
        Reset
        <Icon
          name="radix-icons:cross-2"
          class="ml-2 h-4 w-4"
        />
      </Button>
    </div>
    <div class="flex items-center space-x-2">
      <Button
        variant="outline"
        size="sm"
        @click="onClickRefresh"
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
