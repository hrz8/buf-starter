package query

type QueryResult[T interface{}] struct {
	Data       []*T
	TotalRows  int32
	TotalPages int32
	Filters    map[string][]string
}
