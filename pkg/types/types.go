package types

// Pagination holds pagination state.
type Pagination struct {
	Page     int
	PageSize int
	Total    int64
}

// SortDirection represents SQL sort order.
type SortDirection string

const (
	SortASC  SortDirection = "ASC"
	SortDESC SortDirection = "DESC"
)
