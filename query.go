package altalune

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type SortingParams struct {
	Field string
	Order SortOrder
}

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

type QueryResult[T interface{}] struct {
	Data       []*T
	TotalRows  int32
	TotalPages int32
	Filters    map[string][]string
}
