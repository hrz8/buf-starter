import { SortOrder } from '~~/gen/altalune/v1/common_pb';

import type {
  QueryRequestSchema, PaginationSchema, StringListSchema,
  SortingSchema,
} from '~~/gen/altalune/v1/common_pb';
import type { ColumnFiltersState, SortingState } from '@tanstack/vue-table';
import type { MessageInitShape } from '@bufbuild/protobuf';

export function useQueryRequest(options: {
  page: Ref<number>;
  pageSize: Ref<number>;
  keyword: Ref<string>;
  columnFilters?: ComputedRef<ColumnFiltersState | undefined>;
  sorting?: ComputedRef<SortingState | undefined>;
}) {
  const { page, pageSize, keyword, columnFilters, sorting } = options;

  const debouncedKeyword = refDebounced(keyword, 500);

  const queryRequest = computed<MessageInitShape<typeof QueryRequestSchema>>(() => {
    // Build pagination
    const pagination: MessageInitShape<typeof PaginationSchema> = {
      page: page.value,
      pageSize: pageSize.value,
    };

    // Build filters from column filters
    let filters: Record<string, MessageInitShape<typeof StringListSchema>> | undefined;

    if (columnFilters?.value && columnFilters.value.length > 0) {
      const protoFilters: Record<string, MessageInitShape<typeof StringListSchema>> = {};

      for (const filter of columnFilters.value) {
        const value = filter.value;

        if (value === null || value === undefined || value === '') {
          continue;
        }

        if (Array.isArray(value)) {
          if (value.length > 0) {
            protoFilters[filter.id] = {
              values: value.map((v) => String(v)),
            };
          }
        } else {
          protoFilters[filter.id] = {
            values: [String(value)],
          };
        }
      }

      if (Object.keys(protoFilters).length > 0) {
        filters = protoFilters;
      }
    }

    // Build sorting
    let sortingProto: MessageInitShape<typeof SortingSchema> | undefined;

    if (sorting?.value && sorting.value.length > 0 && sorting.value[0]) {
      sortingProto = {
        field: String(sorting.value[0].id),
        order: sorting.value[0].desc ? SortOrder.DESC : SortOrder.ASC,
      };
    }

    const request: MessageInitShape<typeof QueryRequestSchema> = {
      pagination,
      keyword: debouncedKeyword.value,
      filters,
      sorting: sortingProto,
    };

    return request;
  });

  return {
    queryRequest,
  };
}
