package query

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type SortingParams struct {
	Field string
	Order SortOrder
}
