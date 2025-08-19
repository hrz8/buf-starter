package query

import altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"

type PaginationParams struct {
	Page     int32
	PageSize int32
}

type QueryParams struct {
	Pagination PaginationParams
	Keyword    string
	Filters    map[string][]string
	Sorting    *SortingParams
}

func DefaultQueryParams(req *altalunev1.QueryRequest) *QueryParams {
	queryParamsBuilder := NewQueryParamsBuilder().WithKeyword(req.Keyword)
	if req.Pagination != nil {
		queryParamsBuilder.WithPagination(req.Pagination.Page, req.Pagination.PageSize)
	}
	if req.Filters != nil {
		for field, stringList := range req.Filters {
			queryParamsBuilder.WithFilter(field, stringList.Values)
		}
	}
	if req.Sorting != nil {
		var order SortOrder
		switch req.Sorting.Order {
		case altalunev1.SortOrder_SORT_ORDER_DESC:
			order = SortOrderDesc
		case altalunev1.SortOrder_SORT_ORDER_ASC:
			order = SortOrderAsc
		default:
			order = SortOrderAsc
		}
		queryParamsBuilder.WithSorting(req.Sorting.Field, order)
	}

	return queryParamsBuilder.Build()
}
