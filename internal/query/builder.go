package query

type QueryParamsBuilder struct {
	pagination PaginationParams
	keyword    string
	filters    map[string][]string
	sorting    *SortingParams
}

func NewQueryParamsBuilder() *QueryParamsBuilder {
	return &QueryParamsBuilder{
		pagination: PaginationParams{Page: 1, PageSize: 10},
		filters:    make(map[string][]string),
	}
}

func (b *QueryParamsBuilder) WithPagination(page, pageSize int32) *QueryParamsBuilder {
	b.pagination = PaginationParams{Page: page, PageSize: pageSize}
	return b
}

func (b *QueryParamsBuilder) WithKeyword(keyword string) *QueryParamsBuilder {
	b.keyword = keyword
	return b
}

func (b *QueryParamsBuilder) WithFilter(field string, values []string) *QueryParamsBuilder {
	if len(values) > 0 {
		b.filters[field] = values
	}
	return b
}

func (b *QueryParamsBuilder) WithSorting(field string, order SortOrder) *QueryParamsBuilder {
	if field != "" {
		b.sorting = &SortingParams{Field: field, Order: order}
	}
	return b
}

func (b *QueryParamsBuilder) Build() *QueryParams {
	return &QueryParams{
		Pagination: b.pagination,
		Keyword:    b.keyword,
		Filters:    b.filters,
		Sorting:    b.sorting,
	}
}
