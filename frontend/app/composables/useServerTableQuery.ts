import type { ColumnFiltersState, SortingState } from '@tanstack/vue-table';
import type { QueryOptions } from '#shared/types/query';

export function useServerTableQuery(options: {
  page: Ref<number>;
  pageSize: Ref<number>;
  keyword: Ref<string>;
  columnFilters?: ComputedRef<ColumnFiltersState | undefined>;
  sorting?: ComputedRef<SortingState | undefined>;
}) {
  const { page, pageSize, keyword, columnFilters, sorting } = options;

  const debouncedKeyword = refDebounced(keyword, 500);

  const paginationState = computed(() => ({
    page: page.value,
    pageSize: pageSize.value,
  }));

  const queryOptions = computed<QueryOptions>(() => {
    const opts: QueryOptions = {
      pagination: {
        page: paginationState.value.page,
        pageSize: paginationState.value.pageSize,
      },
      keyword: debouncedKeyword.value,
    };

    if (columnFilters?.value && columnFilters.value.length > 0) {
      const filters: NonNullable<QueryOptions['filters']> = {};
      for (const f of columnFilters.value) {
        filters[f.id] = f.value as string | number | boolean | null | undefined;
      }
      opts.filters = filters;
    }

    if (sorting?.value && sorting.value.length > 0 && sorting.value[0]) {
      opts.sorting = {
        field: String(sorting.value[0].id),
        order: sorting.value[0].desc ? 'desc' : 'asc',
      };
    }

    return opts;
  });

  return {
    queryOptions,
  };
}
