import type { Table } from '@tanstack/vue-table';

export function useDataTableState<T>(
  tableRef: Ref<{ table: Table<T> } | null>,
) {
  const columnFilters = computed(() =>
    tableRef.value?.table.getState().columnFilters,
  );

  const sorting = computed(() =>
    tableRef.value?.table.getState().sorting,
  );

  const columnVisibility = computed(() =>
    tableRef.value?.table.getState().columnVisibility,
  );

  return {
    columnFilters,
    sorting,
    columnVisibility,
  };
}

export function useDataTableFilter<T>(
  table: ComputedRef<Table<T> | undefined>,
  columnId: string,
) {
  const filterValues = ref<string[]>([]);

  const setFilter = (selected: string[]) => {
    const column = table.value?.getColumn(columnId);
    if (column) {
      column.setFilterValue(selected.length ? selected : undefined);
    }
  };

  const clearFilter = () => {
    filterValues.value = [];
    const column = table.value?.getColumn(columnId);
    if (column) {
      column.setFilterValue(undefined);
    }
  };

  return {
    filterValues,
    setFilter,
    clearFilter,
  };
}
